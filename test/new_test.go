package test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	"github.com/bearki/belog/v2/internal/convert"
	"github.com/bearki/belog/v2/logger/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestXxx(tt *testing.T) {
	tmp := make([]byte, 30, 30)
	t := time.Now()
	for i := 4; i < 27; i++ {
		tmp[i] = ' '
	}
	e := convert.TimeToBytes(tmp[4:27], t)
	if e != nil {
		tt.Fatal(e)
	}
	fmt.Println(tmp)
}

// BenchmarkNewFileBelog 实例方式输出文件日志
func BenchmarkNewFileBelog(b *testing.B) {
	// 初始化文件日志适配器
	fielAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      200,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        false,                  // 开启异步写入(main函数提前结束会导致日志未写入)
		AsyncChanCap: 100,                    // 异步缓存管道容量
	})
	if err != nil {
		panic("file adapter create failed, " + err.Error())
	}
	// 初始化一个实例(可实例化任意适配器)
	mylog, err := belog.New(fielAdapter)
	if err != nil {
		panic("belog logger create failed, " + err.Error())
	}
	for i := 0; i < b.N; i++ {
		mylog.Tracef(
			"this is a trace log",
			field.Intn("index", int64(i)),
		)
		// mylog.Trace("this is a trace log, value: %d, index: %v, bba: %v", i, i, i)
	}
}

func BenchmarkZap(b *testing.B) {
	file, _ := os.Create("../logs/zap.log")
	writeSyncer := zapcore.AddSync(file)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	defer sugarLogger.Sync()
	for i := 0; i < b.N; i++ {
		sugarLogger.Error(
			"this is a trace log",
			zap.Int("index", i),
		)
	}
}
