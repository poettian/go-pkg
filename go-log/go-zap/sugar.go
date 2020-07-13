package go_zap

import (
	"go.uber.org/zap"
	"time"
)

func SugarExample() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()

	// example 提供了一个仅用于了解zap的logger，它会将 debug 及以上级别的日志以json格式写入到标准输出中
	// 但是会忽略时间和调用栈的信息，以保证输出日志的简短性
	// output：{"level":"debug","msg":"http://www.baidu.com1 2020-07-11 17:49:38.981799 +0800 CST m=+0.000524096"}
	sugar.Debug("http://www.baidu.com", 1, time.Now())
}

func SugarDevelopment()  {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic("zap NewDevelopment fail.")
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	// development 的 encoder 是 console
	// output：2020-07-11T17:49:38.982+0800    DEBUG   go-zap/sugar.go:23      failed to fetch URL: http://example.com
	sugar.Debugf("failed to fetch URL: %s", "http://example.com")
}

func SugarProduction()  {
	logger, err := zap.NewProduction()
	if err != nil {
		panic("zap NewProduction fail.")
	}
	sugar := logger.Sugar()
	defer sugar.Sync()

	// 需要注意的是：production 只记录 info 及以上级别的 log
	// output：{"level":"info","ts":1594460978.982084,"caller":"go-zap/sugar.go:34","msg":"failed to fetch URL: http://example.com"}
	sugar.Infof("failed to fetch URL: %s", "http://example.com")
}