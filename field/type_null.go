/**
 * @Title 对任意指针类型进行判断
 * @Desc 支持所有指针检查，为空时将输出null
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

import (
	"unsafe"

	"github.com/bearki/belog/v2/internal/convert"
)

// null常量数组，用于切片的底层引用
var nullConstArray = [4]byte{'n', 'u', 'l', 'l'}

// CheckPtr 检查任意字段的值指针
//
// 注意: 指针为空时将会输出null
//
// @params name 字段名称
//
// @params value 任意类型指针
//
// @return 当指针为空时会创建null的字段数据
//
// @return 指针是否为空
func CheckPtr(name string, value unsafe.Pointer) (Field, bool) {
	// 检查指针是否为空
	if value == nil {
		// 指针为空，检查不通过
		return Field{
			Key:     convert.StringToBytes(name),
			ValType: TypeNull,
			Bytes:   nullConstArray[:],
		}, false
	}
	// 指针不为空，检查通过
	return Field{}, true
}
