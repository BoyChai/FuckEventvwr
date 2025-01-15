package control

import (
	"FuckEventvwr/output"
	"fmt"
)

// 工作过程时的错误处理
func addError(e string) {
	fmt.Println(e)
	output.Output.WriteError(e)
}
