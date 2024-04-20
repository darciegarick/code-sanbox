package global

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError       CustomError
	ValidateError       CustomError
	TokenError          CustomError
	AnswerError         CustomError
	CompileError        CustomError
	MemoryLimitExceeded CustomError
	TimeLimitExceeded   CustomError
	PresentationError   CustomError
	OutputLimitExceeded CustomError
	DangerousOperation  CustomError
	RuntimeError        CustomError
	SystemError         CustomError
}

var Errors = CustomErrors{
	BusinessError:       CustomError{40000, "业务错误"},
	ValidateError:       CustomError{42200, "请求参数错误"},
	TokenError:          CustomError{40100, "登录授权失效"},
	AnswerError:         CustomError{50001, "答案错误"},
	CompileError:        CustomError{50002, "编译错误"},
	MemoryLimitExceeded: CustomError{50003, "内存溢出"},
	TimeLimitExceeded:   CustomError{50004, "超时"},
	PresentationError:   CustomError{50005, "格式错误"},
	OutputLimitExceeded: CustomError{50006, "输出溢出"},
	DangerousOperation:  CustomError{50007, "危险操作"},
	RuntimeError:        CustomError{50008, "用户程序运行错误"},
	SystemError:         CustomError{50009, "系统错误"},
}
