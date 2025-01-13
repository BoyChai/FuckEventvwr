package control

import (
	"fmt"
	"sync"
)

var errorLock sync.Mutex

var Error []string

// 工作过程时的错误处理
func addError(e string) {
	errorLock.Lock()
	defer errorLock.Unlock()
	fmt.Println(e)
	Error = append(Error, e)
}
