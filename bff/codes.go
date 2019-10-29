package bff

import "fmt"

// service inner error code list
const (
	CodeSuccess                = 0
	CodeVersionError           = 1
	CodeUpdating               = 101
	CodeNotFound               = 1000
	CodeRequestUrlParamError   = 1001
	CodeRequestQueryParamError = 1002
	CodeRequestBodyError       = 1004
	CodeNotLogin               = 1005
	CodeServerError            = 1006
	CodeNotAllow               = 1007
)

type Codes map[int]string

var (
	_codes = Codes{
		CodeSuccess:                "请求成功",
		CodeVersionError:           "客户端版本错误，请升级客户端",
		CodeUpdating:               "服务正在升级",
		CodeNotFound:               "找不到对于系统&模块",
		CodeRequestUrlParamError:   "URL参数错误",
		CodeRequestQueryParamError: "查询参数错误",
		CodeNotLogin:               "用户没有登录",
		CodeRequestBodyError:       "请求结构错误",
		CodeServerError:            "服务器错误",
		CodeNotAllow:               "权限校验失败",
	}
)

func GetMessage(code int) string {
	return _codes[code]
}

func CodeIter(fn func(key int, value string) bool) {
	for key, value := range _codes {
		if !fn(key, value) {
			return
		}
	}
}

// 请在初始化阶段使用该函数，运行时使用该函数可能导致同步问题
func MergeCodes(codes Codes) {

	for code, describe := range codes {

		_, exists := _codes[code]

		if exists {
			panic(
				fmt.Sprintf("code %d[%s] already exists", code, GetMessage(code)),
			)
		}

		_codes[code] = describe

	}
}
