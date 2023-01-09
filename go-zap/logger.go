package go_zap

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
)

func NewLoggerFromConfig() {
	var cfg zap.Config
	rawJson := []byte(`{
		"level": "debug",
	  	"encoding": "json",
	  	"outputPaths": ["stdout", "/Users/tianzhiwei/Code/golang/learn-packages/go_log/go-zap/zap.log"],
	  	"errorOutputPaths": ["stderr"],
	  	"initialFields": {"foo": "bar"},
	  	"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
	  	}
	}`)
	if err := json.Unmarshal(rawJson, &cfg);err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

func NewLoggerFromCore()  {
	// 注意，这里是把一个匿名函数转为了 LevelEnablerFunc 类型，从而实现 LevelEnabler 接口
	// 所以，可以这样调用来判定log是否符合写入级别：highPriority.Enabled()
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// 关于：ioutil.Discard
	// 实现了io.Writer，但是所有的 Write 方法不会做任何操作，直接返回成功。有点类似 /dev/null

	// AddSync 把一个 io.Writer 接口值转换为了 zapcore.WriteSyncer 接口值
	// 如果 io.Writer 底层值已经实现了 WriteSyncer，那么 Sync 方法就使用该底层值提供的 Sync 方法
	// 否则，添加一个没有任何操作的 Sync 方法
	// Sync 方法的作用应该是类似 flush 的作用，也就是 defer logger.Sync 实际调用的方法

	// 创建了两个 WriteSyncer，在调用 zapcore.NewCore 时使用
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	// zapcore.Lock 方法其实是对 zapcore.WriteSyncer 接口值的升级，通过加锁来保证写入的并发安全性，但是仍然返回一个 WriteSyncer 接口值
	// 但是调用 write 方法时，实际调用的是 lockedWriteSyncer.Write
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// 这里没啥好说的，针对不同的output，使用不同的 encoder 配置
	// 重点是看 zapcore.EncoderConfig 的配置项
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// 必须使用 New 来实例化一个 logger
	logger := zap.New(core)
	// 这个应该和 DB 类似，在主程退出时调用，以确保日志持久化
	defer logger.Sync()

	logger.Info("constructed a logger")
}

func NewLoggerFromTest()  {
	encoderCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
	}
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, zapcore.DebugLevel)
	logger := zap.New(core)
	defer logger.Sync()

	logger.Debug("这是一个测试项")
}
