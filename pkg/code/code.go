package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// InvalidURI 无效的URI
	InvalidURI = 90014000000
	// InvalidParams 无效的参数
	InvalidParams = 90014000001
	// InvalidTimestamp 无效的时间格式
	InvalidTimestamp = 90014000002
	// ErrDataExist 名字已经存在
	ErrDataExist = 90014000003
)

// CodeTable 码表
var CodeTable = map[int64]string{
	InvalidURI:       "无效的URI.",
	InvalidParams:    "无效的参数.",
	InvalidTimestamp: "无效的时间格式.",
	ErrDataExist:     "数据不存在",
}
