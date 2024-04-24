package entity

type ExecuteMessage struct {
	Language string
	Code     string
	testCase []string
}

// type ExecuteMessage struct {
// 	ExitValue    int
// 	Message      string
// 	ErrorMessage string
// 	Time         int64
// 	Memory       int64
// }
