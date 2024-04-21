package response

type ExecuteCodeResponse struct {
	InputList []string
	Code      string
	Language  string
}

type JudgeInfo struct {
	Message string
	Memory  int64
	Time    int64
}
