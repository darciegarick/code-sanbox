package response

type CreateContainerResponse struct {
	ContainerId     string
}

type ExecuteCodeResponse struct {
	Language     string
	Result       string
	TestCaseList []string
}

type JudgeInfo struct {
	Message string
	Memory  int64
	Time    int64
}
