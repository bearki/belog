package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bearki/belog/v3"
	"github.com/bearki/belog/v3/adapter/discard"
	"github.com/bearki/belog/v3/field"
	"github.com/bearki/belog/v3/logger"
)

// BenchmarkBelogLoggerFormatStatic 测试belog标准记录器序列化静态字符串
func BenchmarkBelogLoggerFormatStatic(b *testing.B) {
	// 初始化一个实例(无输出)
	l, err := belog.New(logger.Option{}, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}

	// 重置测试参数
	b.ReportAllocs()
	b.StartTimer()

	// 执行测试
	for i := 0; i < b.N; i++ {
		l.Info("this is a info log")
	}
}

// BenchmarkBelogLoggerFormatFiveFields 测试belog标准记录器序列化5个字段
func BenchmarkBelogLoggerFormatFiveFields(b *testing.B) {
	// 初始化一个实例(无输出)
	l, err := belog.New(logger.Option{}, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}

	// 重置测试参数
	b.ReportAllocs()
	b.StartTimer()

	// 执行测试
	for i := 0; i < b.N; i++ {
		tb := i%2 == 0
		ts := "value"
		tf := 3.1415926
		tt := time.Now()
		l.Info(
			"this is a info log",
			field.Int("key1", i),
			field.Bool("key2", tb),
			field.String("key3", ts),
			field.Float64("key4", tf),
			field.Time("key5", tt),
		)
	}
}

// BenchmarkBelogLoggerFormatTenFields 测试belog标准记录器序列化10个字段
func BenchmarkBelogLoggerFormatTenFields(b *testing.B) {
	// 初始化一个实例(无输出)
	l, err := belog.New(logger.Option{}, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}

	// 重置测试参数
	b.ReportAllocs()
	b.StartTimer()

	// 执行测试
	for i := 0; i < b.N; i++ {
		tb := i%2 == 0
		ts := "value"
		tf := 3.1415926
		tt := time.Now()
		l.Info(
			"this is a info log",
			field.Int("key1", i),
			field.Bool("key2", tb),
			field.String("key3", ts),
			field.Float64("key4", tf),
			field.Time("key5", tt),
			field.Intp("key6", &i),
			field.Boolp("key7", &tb),
			field.Stringp("key8", &ts),
			field.Float64p("key9", &tf),
			field.Timep("key10", &tt),
		)
	}
}

// BenchmarkBelogLoggerFormatSlice 测试belog标准记录器序列化切片
func BenchmarkBelogLoggerFormatSlice(b *testing.B) {
	// 初始化一个实例(无输出)
	l, err := belog.New(logger.Option{}, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}

	// 重置测试参数
	b.ReportAllocs()
	b.StartTimer()

	// 执行测试
	for i := 0; i < b.N; i++ {
		l.Info(
			"this is a info log",
			field.Ints("key", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}),
		)
	}
}

// BenchmarkBelogLoggerFormatInterface 测试belog标准记录器序列化反射类型
func BenchmarkBelogLoggerFormatInterface(b *testing.B) {
	// 初始化一个实例(无输出)
	l, err := belog.New(logger.Option{}, discard.New())
	if err != nil {
		fmt.Printf("belog logger create failed, %s\r\n", err)
		return
	}

	// 重置测试参数
	b.ReportAllocs()
	b.StartTimer()

	// 执行测试
	for i := 0; i < b.N; i++ {
		l.Info(
			"this is a info log",
			field.Interface("key", map[string]string{"key1": "key1", "key2": "key2", "key3": "key3", "key4": "key4", "key5": "key5"}),
		)
	}
}
