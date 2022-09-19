package field

import "github.com/bearki/belog/v2/internal/convert"

// Append 将字段拼接为行格式
//
// 注意：为防止dst发生扩容，请在外部提前申请好容量
//
// 拼接格式：k1: v1, k2: v2, ..., message
//
// @params dst 已申请好容量的切片
//
// @params message 消息内容
//
// @params val 字段列表
//
// @return 组装好的JSON数据
func Append(dst []byte, message string, val ...Field) []byte {
	// 遍历所有字段
	for _, v := range val {
		// 追加字段并序列化
		dst = append(dst, v.KeyBytes...)
		dst = append(dst, ':', ' ')
		dst = append(dst, v.ValBytes...)
		dst = append(dst, ',', ' ')
		// 回收到复用池
		v.Put()
	}

	// 追加message内容
	dst = append(dst, convert.StringToBytes(message)...)

	// 返回组装好的数据
	return dst
}

// AppendJSON 将字段拼接为json格式
//
// 注意：为防止dst发生扩容，请在外部提前申请好容量
//
// [startEndSym == true]: {"fields": {"k1": "v1", ...}, "msg": "message"}
//
// [startEndSym == false]: "fields": {"k1": "v1", ...}, "msg": "message"
//
// @params dst 已申请好容量的切片
//
// @params fieldsKey 包裹所有字段的键名
//
// @params messageKey 消息的键名
//
// @params startEndSym 是否需要开始和结束符号 `{` 和 `}`
//
// @params message 消息内容
//
// @params val 字段列表
//
// @return 组装好的JSON数据
func AppendJSON(dst []byte, fieldsKey string, messageKey string, startEndSym bool, message string, val ...Field) []byte {
	// 是否需要追加开始和结束符号
	if startEndSym {
		dst = append(dst, '{')
	}
	// 追加字段集字段
	dst = append(dst, '"')
	dst = append(dst, fieldsKey...)
	dst = append(dst, `": {`...)
	// 是否需要追加分隔符了
	appendDelimiter := false
	// 遍历所有字段
	for _, v := range val {
		// 从第二个有效字段开始追加分隔符号
		if appendDelimiter {
			dst = append(dst, ',', ' ')
		}

		// 追加字段并序列化
		dst = append(dst, '"')
		dst = append(dst, v.KeyBytes...)
		dst = append(dst, '"')
		dst = append(dst, ':', ' ')
		dst = append(dst, v.ValPrefixBytes...)
		dst = append(dst, v.ValBytes...)
		dst = append(dst, v.ValSuffixBytes...)

		// 已经填充了一个有效字段了
		if !appendDelimiter {
			// 下一次需要追加分隔符
			appendDelimiter = true
		}

		// 回收到复用池
		v.Put()
	}
	// 追加字段结束括号
	dst = append(dst, `}, `...)

	// 追加message字段及其内容
	dst = append(dst, '"')
	dst = append(dst, messageKey...)
	dst = append(dst, `": "`...)
	dst = append(dst, convert.StringToBytes(message)...)
	dst = append(dst, '"')

	// 是否需要追加开始和结束符号
	if startEndSym {
		dst = append(dst, '}')
	}

	// 返回组装好的数据
	return dst
}
