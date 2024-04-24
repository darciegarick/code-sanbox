package request

type ExecuteCodeRequest struct {
	TestCase []string `form:"test_case" json:"test_case"`
	Code     string   `form:"code" json:"code" binding:"required"`
	Language string   `form:"language" json:"language"`
}

func (codeRquest ExecuteCodeRequest) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"code.required": "用户代码不能为空",
	}
}
