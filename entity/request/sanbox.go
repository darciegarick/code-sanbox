package request

type ExecuteCodeRequest struct {
	InputList []string
	Code      string
	Language  string
}
