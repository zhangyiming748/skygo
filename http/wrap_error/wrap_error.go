package wrap_error

type WrapError struct{}

func (this WrapError) GetWapErrorCode(rpcErrCode int64) int64 {
	if val, exist := CustomWrapErrMap[rpcErrCode]; exist {
		return val
	} else {
		return DEFAULT_ERROR_CODE
	}
}

const DEFAULT_ERROR_CODE = -1 //默认错误返回码

// 自定义http接口返回错误码
// key：RPC错误码
// val：HTTP错误码
var CustomWrapErrMap = map[int64]int64{
	-2001100001: -1,
	-3001100002: -2, //没传递token
	3001100012:  -2, //token验证失败
	2001400052:  -3, //权限问题
	2001400053:  -3, //权限问题
	2001400054:  -3, //权限问题
	2001400055:  -3, //权限问题
}
