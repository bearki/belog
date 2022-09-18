package test

import (
	"fmt"
	"testing"
)

// TestDefultBelog 默认方式输出日志
func TestDefultBelog(t *testing.T) {
	type ab struct{}
	var x interface{} = complex64(3)
	switch x.(type) {
	case bool:
		fmt.Println("bool")
	case string:
		fmt.Println("string")
	case int8:
		fmt.Println("int8")
	case int16:
		fmt.Println("int16")
	case int32:
		fmt.Println("int32")
	case int64:
		fmt.Println("int64")
	case int:
		fmt.Println("int")
	case uint8:
		fmt.Println("uint8")
	case uint16:
		fmt.Println("uint16")
	case uint32:
		fmt.Println("uint32")
	case uint64:
		fmt.Println("uint64")
	case uint:
		fmt.Println("uint")
	case float32:
		fmt.Println("float32")
	case float64:
		fmt.Println("float64")
	case complex64:
		fmt.Println("complex64")
	case complex128:
		fmt.Println("complex128")
	case uintptr:
		fmt.Println("uintptr")
	default:

	}
}
