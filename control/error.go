package control

var errorChan = make(chan string, 100)

// 工作过程时的错误处理
func addError(e string) {
	errorChan <- e
}
