package custom_error

const (
	// 模块通用错误码

	// ErrCodeCommonParamMiss 缺少参数错误
	ErrCodeCommonParamMiss = -900001
	// ErrCodeCommonParamInvalid 参数格式无效错误
	ErrCodeCommonParamInvalid = -900002
	// ErrCodePageParam 页数错误
	ErrCodePageParam = -900004
	// ErrCodePageLimitParam 页条数错误
	ErrCodePageLimitParam = -900005
	// ErrCodePageLimitGtMax 页条数超过最大值
	ErrCodePageLimitGtMax = -900006
	// ErrCodeCommonGetRowMiss 数据获取不存在
	ErrCodeCommonGetRowMiss = -900101
	// ErrCodeCommonGetDbDataFail 获取数据失败
	ErrCodeCommonGetDbDataFail = -900102

	// ErrCodeCommonCreateDbDataFail 创建数据失败
	ErrCodeCommonCreateDbDataFail = -900501
	// ErrCodeCommonUpdateDbDataFail 更新数据失败
	ErrCodeCommonUpdateDbDataFail = -900502
	// ErrCodeCommonDeleteDbDataFail 删除数据失败
	ErrCodeCommonDeleteDbDataFail = -900503
	// ErrCodeCommonOperaDbFail 操作数据库失败
	ErrCodeCommonOperaDbFail = -900504

	// ErrCodeCommonCustomFail 自定义失败
	ErrCodeCommonCustomFail = -900605

	// 上传文件失败
	ErrCodeFileUploadFail = -800001
	ErrCodeInvalidRefer   = -800002
)

func errorCode(code int64) int64 {
	return code
}

type SysError struct {
	Code int64
	EMsg string
}

func (this SysError) Error() string {
	return this.EMsg
}

func NewError(code int, msg string) SysError {
	return SysError{
		errorCode(int64(code)),
		msg,
	}
}

func NewParamMissError(msg string) SysError {
	return SysError{
		errorCode(ErrCodeCommonParamMiss),
		msg,
	}
}

func NewParamInvalidError(msg string) SysError {
	return SysError{
		errorCode(ErrCodeCommonParamInvalid),
		msg,
	}
}

func NewOperaDbFailError(msg string) SysError {
	return SysError{
		errorCode(ErrCodeCommonOperaDbFail),
		msg,
	}
}

func NewCustomFailError(msg string) SysError {
	return SysError{
		errorCode(ErrCodeCommonCustomFail),
		msg,
	}
}

var (
	ErrFileUploadFail = SysError{errorCode(ErrCodeFileUploadFail), "文件上传失败"}
	ErrInvalidRefer   = SysError{errorCode(ErrCodeInvalidRefer), "无效refer"}

	/* 通用错误 */
	ErrPageParam      = SysError{errorCode(ErrCodePageParam), "页码错误"}
	ErrPageLimitParam = SysError{errorCode(ErrCodePageLimitParam), "页尺寸错误"}
	ErrPageLimitGtMax = SysError{errorCode(ErrCodePageLimitGtMax), "页尺寸超过最大限制"}

	ErrGetRowMiss       = SysError{errorCode(ErrCodeCommonGetRowMiss), "没有发现数据"}
	ErrGetDbDataFail    = SysError{errorCode(ErrCodeCommonGetDbDataFail), "获取数据失败"}
	ErrCreateDbDataFail = SysError{errorCode(ErrCodeCommonCreateDbDataFail), "创建数据失败"}
	ErrUpdateDbDataFail = SysError{errorCode(ErrCodeCommonUpdateDbDataFail), "更新数据失败"}
	ErrDeleteDbDataFail = SysError{errorCode(ErrCodeCommonDeleteDbDataFail), "删除数据失败"}
	ErrOperaDbFail      = SysError{errorCode(ErrCodeCommonOperaDbFail), "操作数据库失败"}
)
