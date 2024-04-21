package entity

type ExecuteMessage struct {
	ExitValue    int
	Message      string
	ErrorMessage string
	Time         int64
	Memory       int64
}
