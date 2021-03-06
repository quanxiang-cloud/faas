package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

// Errcode
const (
	InvalidURI              = 90014000000
	InvalidParams           = 90014000001
	InvalidTimestamp        = 90014000002
	ErrDataExist            = 90014000003
	ErrDataNotExist         = 90014000004
	ErrFunctionExist        = 90014000005
	ErrDataIllegal          = 90014000006
	ErrNotSupportedLanguage = 90014000007
)

// CodeTable 码表
var CodeTable = map[int64]string{
	InvalidURI:              "无效的URI.",
	InvalidParams:           "无效的参数.",
	InvalidTimestamp:        "无效的时间格式.",
	ErrDataExist:            "数据已经存在",
	ErrDataNotExist:         "数据不存在",
	ErrFunctionExist:        "函数已创建",
	ErrDataIllegal:          "数据不合法",
	ErrNotSupportedLanguage: "不支持的语言(%s:%s)",
}
