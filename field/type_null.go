/**
 * @Title 空类型处理
 * @Desc 输出null
 * @Author Bearki
 * @DateTime 2022/09/17 21:46
 */

package field

// 创建null值字段
//
//	@var key 字段名称
func nullField(key string) Field {
	return Field{Key: key, Type: TypeNull, String: "null"}
}
