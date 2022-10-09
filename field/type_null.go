/**
 * @Title 空类型处理
 * @Desc 输出null
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

// 创建null值字段
//
// @params name 字段名称
func nullField(name string) Field {
	return Field{Key: name, ValType: TypeNull, String: "null"}
}
