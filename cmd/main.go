package main

import (
	"github.com/bearki/belog/v2"
	"github.com/bearki/belog/v2/adapter/file"
	"github.com/bearki/belog/v2/logger/field"
)

func main() {
	// 初始化文件日志适配器
	fielAdapter, err := file.New(file.Options{
		LogPath:      "../logs/test_new.log", // 日志储存路径
		MaxSize:      200,                    // 日志单文件大小
		MaxLines:     1000000,                // 单文件最大行数
		SaveDay:      7,                      // 日志保存天数
		Async:        true,                   // 开启异步写入(main函数提前结束会导致日志未写入)
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
	mylog.Tracef("this is a trace log", field.Intn("value", 3))
}
