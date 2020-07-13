## zap 的使用文档

先放一张网上关于zap的架构图

![zap架构图](https://tva1.sinaimg.cn/large/007S8ZIlly1ggn3ptnwwoj30nh0h7q55.jpg)

#### 基础概念

`zap` : 是一个对外的接口，提供了简单易用的函数，初级使用zap的话，只了解这一层的使用方式即可

`zapcore` : 核心层，是zap强大功能的基础，如果需要定制使用zap，那么这一层必须了解

`encoder` : 编码器，将用户提供的日志信息，以指定的格式编码

#### Sugar

如果对性能要求不是那么高，但是需要一种易用的方式来使用zap，那么可以使用zap提供的Sugar：

```go
// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func (log *Logger) Sugar() *SugaredLogger {
	core := log.clone()
	core.callerSkip += 2
	return &SugaredLogger{core}
}
```

来看一下 SugaredLogger：

SugaredLogger 提供了  `structured` 和  `printf-style` logging。

```go
// SugaredLogger 组合了 Logger
type SugaredLogger struct {
	base *Logger
}

// 总的来看，SugaredLogger 提供了三种写日志的方式,这三种方式都是调用的 SugaredLogger.log

// Debug uses fmt.Sprint to construct and log a message.
func (s *SugaredLogger) Debug(args ...interface{})

// Debugf uses fmt.Sprintf to log a templated message.
func (s *SugaredLogger) Debugf(template string, args ...interface{})

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{})
```

对于 structured logging，可以使用 `SugaredLogger.With` 来添加 context。

要注意的是：context 的 key 必须是一个字符串，否则在 development 会 panic，在production 会返回 error。

#### Logger

如果对性能和内存分配有严格的要求，那么应该使用更为底层的 logger，但是它只支持强类型、结构化的日志记录。

```go
logger.Info("failed to fetch URL",
  zap.String("url", "http://example.com"),
  zap.Int("attempt", 3),
  zap.Duration("backoff", time.Second),
)
```

此外，官方文档中有这么一段描述：

> The zap package itself is a relatively thin wrapper around the interfaces in go.uber.org/zap/zapcore. Extending zap to support a new encoding (e.g., BSON), a new log sink (e.g., Kafka), or something more exotic (perhaps an exception aggregation service, like Sentry or Rollbar) typically requires implementing the zapcore.Encoder, zapcore.WriteSyncer, or zapcore.Core interfaces. See the zapcore documentation for details.
>
> Similarly, package authors can use the high-performance Encoder and Core implementations in the zapcore package to build their own loggers.

有两种方式创建 logger：

```go
// logger.go
func New(core zapcore.Core, options ...Option) *Logger

// config.go
// 实际上仍是调用了 New 来创建logger
func (cfg Config) Build(opts ...Option) (*Logger, error) 
```

#### zap.Config

提供了一种构造 logger 的声明性的方法。相比直接调用 zapcore，config 更加的简单。但是，config 只提供比较通用的选项，更加定制化的设置还是需要调用 zapcore 来完成。

```go
type Config struct {
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	// 支持的最低级别，这是个动态的级别，通过调用Config.Level.SetLevel可以原子性、运行时更改包括所有子logger的级别
	Level AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool `json:"disableStacktrace" yaml:"disableStacktrace"`
	// Sampling sets a sampling policy. A nil SamplingConfig disables sampling.
	// 没太看明白这个是干啥的
	Sampling *SamplingConfig `json:"sampling" yaml:"sampling"`
	// Encoding sets the logger's encoding. Valid values are "json" and
	// "console", as well as any third-party encodings registered via
	// RegisterEncoder.
	// 编码方式：json || console
	Encoding string `json:"encoding" yaml:"encoding"`
	// EncoderConfig sets options for the chosen encoder. See
	// zapcore.EncoderConfig for details.
	// 编码设置项，这个是后面需要重点看的一个设置项
	EncoderConfig zapcore.EncoderConfig `json:"encoderConfig" yaml:"encoderConfig"`
	// OutputPaths is a list of URLs or file paths to write logging output to.
	// See Open for details.
	// 日志输出的 urls 或 files，也就是日志写到哪里去
	OutputPaths []string `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of URLs to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	// 输出logger内部错误的地方
	ErrorOutputPaths []string `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	// InitialFields is a collection of fields to add to the root logger.
	// 添加到 root logger 中的字段，应该是每条日志都会添加这些字段
	InitialFields map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}
```

#### zapcore

先看官方手册的描述：

```go
// The bundled Config struct only supports the most common configuration
// options. More complex needs, like splitting logs between multiple files
// or writing to non-file outputs, require use of the zapcore package.
```

意思就是 zap.Config 只是提供了常用项，如果定制化要求高，那么就要使用 zapcore。

具体来看官方提供的示例：

```go
// 首先，LevelEnablerFunc 是一个函数类型
// 其次，LevelEnablerFunc 实现了 zapcore.LevelEnabler 接口。
// 在分割不同级别的log进入不同的outputs时，该函数非常有用
type LevelEnablerFunc func(zapcore.Level) bool

// 该接口的作用是：判定一条log是否符合写入的级别
// 注意：zapcore.DebugLevel 和 zap.LevelEnablerFunc 都实现了该接口
type LevelEnabler interface {
	Enabled(Level) bool
}

// 官方示例：
// 注意，这里是把一个匿名函数转为了 LevelEnablerFunc 类型，从而实现 LevelEnabler 接口
// 所以，可以这样调用来判定log是否符合写入级别：highPriority.Enabled()
highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
  return lvl >= zapcore.ErrorLevel
})
lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
  return lvl < zapcore.ErrorLevel
})
```

```go
// AddSync 把一个 io.Writer 接口值转换为了 zapcore.WriteSyncer 接口值
// 如果 io.Writer 底层值已经实现了 WriteSyncer，那么 Sync 方法就使用该底层值提供的 Sync 方法
// 否则，添加一个没有任何操作的 Sync 方法
// Sync 方法的作用应该是类似 flush 的作用，也就是 defer logger.Sync 实际调用的方法
func AddSync(w io.Writer) WriteSyncer


// 关于：ioutil.Discard
// 实现了io.Writer，但是所有的 Write 方法不会做任何操作，直接返回成功。有点类似 /dev/null

// 官方示例：
// 创建了两个 WriteSyncer，在调用 zapcore.NewCore 时使用
topicDebugging := zapcore.AddSync(ioutil.Discard)
topicErrors := zapcore.AddSync(ioutil.Discard)
```

```go
// 这里要关注的是：结构体的类型组合
type lockedWriteSyncer struct {
	sync.Mutex
	ws WriteSyncer
}

// 这里要关注的是：接口的类型断言和Lock方法的返回值类型
// Lock wraps a WriteSyncer in a mutex to make it safe for concurrent use. In
// particular, *os.Files must be locked before use.
func Lock(ws WriteSyncer) WriteSyncer {
	if _, ok := ws.(*lockedWriteSyncer); ok {
		// no need to layer on another lock
		return ws
	}
	return &lockedWriteSyncer{ws: ws}
}


// 官方示例：
// zapcore.Lock 方法其实是对 zapcore.WriteSyncer 接口值的升级，通过加锁来保证写入的并发安全性，但是仍然返回一个 WriteSyncer 接口值
// 但是调用 write 方法时，实际调用的是 lockedWriteSyncer.Write 
consoleDebugging := zapcore.Lock(os.Stdout)
consoleErrors := zapcore.Lock(os.Stderr)
```

**关于编码**

内置有 `console` 和 `json` 这两种方式

```go
// 这个是重点，是定制化日志内容和格式的重要配置

// An EncoderConfig allows users to configure the concrete encoders supplied by
// zapcore.
type EncoderConfig struct {
	
	// 如果某一项为空则写日志时会被忽略，MessageKey不能省略
	// 注意：如果配置了 LevelKey，则必须同时配置 EncodeLevel，否则可能会报错，其它的 key 也是一样
	// Set the keys used for each log entry. If any key is empty, that portion
	// of the entry is omitted.
	MessageKey    string `json:"messageKey" yaml:"messageKey"` 	// 消息体的key
	LevelKey      string `json:"levelKey" yaml:"levelKey"`			// 消息级别的key
	TimeKey       string `json:"timeKey" yaml:"timeKey"`				// 时间的key
	NameKey       string `json:"nameKey" yaml:"nameKey"`			  // 不知道干啥的
	CallerKey     string `json:"callerKey" yaml:"callerKey"`		// 调用者的key
	StacktraceKey string `json:"stacktraceKey" yaml:"stacktraceKey"`	// 调用栈的key
	LineEnding    string `json:"lineEnding" yaml:"lineEnding"`	// 日志结尾的字符
  
	// 配置常见复杂类型的基本表示，白话点说就是级别、时间等该显示为什么格式。
	// 这几个成员都是函数类型，可以直接使用zap内置的，也可以定制。
	// 看函数定义的话，有个 PrimitiveArrayEncoder 接口，见下面的解析
	// Configure the primitive representations of common complex types. For
	// example, some users may want all time.Times serialized as floating-point
	// seconds since epoch, while others may prefer ISO8601 strings.
	EncodeLevel    LevelEncoder    `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime     TimeEncoder     `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration DurationEncoder `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller   CallerEncoder   `json:"callerEncoder" yaml:"callerEncoder"`
  
	// 暂时没看明白这个是干嘛的
	// Unlike the other primitive type encoders, EncodeName is optional. The
	// zero value falls back to FullNameEncoder.
	EncodeName NameEncoder `json:"nameEncoder" yaml:"nameEncoder"`
}


// 看注释不是很好理解，其实这个接口就是定义了一系列以指定类型写入编码值的方法
// 比如我们上面定义了 EncodeTime 函数，那么该函数最终生成的值要返回给 zap。如何返回？就是通过该接口定义的方法来返回。
// 明确地告诉 zap 编码值的类型
// PrimitiveArrayEncoder is the subset of the ArrayEncoder interface that deals
// only in Go's built-in types. It's included only so that Duration- and
// TimeEncoders cannot trigger infinite recursion.
type PrimitiveArrayEncoder interface {
	// Built-in types.
	AppendBool(bool)
	AppendByteString([]byte) // for UTF-8 encoded bytes
	AppendComplex128(complex128)
	AppendComplex64(complex64)
	AppendFloat64(float64)
	AppendFloat32(float32)
	AppendInt(int)
	AppendInt64(int64)
	AppendInt32(int32)
	AppendInt16(int16)
	AppendInt8(int8)
	AppendString(string)
	AppendUint(uint)
	AppendUint64(uint64)
	AppendUint32(uint32)
	AppendUint16(uint16)
	AppendUint8(uint8)
	AppendUintptr(uintptr)
}
```

```go
// 这个没啥好说的，类似于 linux 下的 tee 命令
// NewTee creates a Core that duplicates log entries into two or more
// underlying Cores.
//
// Calling it with a single Core returns the input unchanged, and calling
// it with no input returns a no-op Core.
func NewTee(cores ...Core) Core
```

```go
// 最后要注意的是写 field 时，zap 提供的方法：
// 这些可用的方法在 field.go 中
logger.Info("failed to fetch URL",
  zap.String("url", "http://example.com"),
  zap.Int("attempt", 3),
  zap.Duration("backoff", time.Second),
)
```

