package utils

import (
	"falcon-mail-transmit/lib/log"
)

func HandlePanic(tasks ...func()) {

	for _, t := range tasks {
		go func(f func()) { // 匿名函数的参数为业务逻辑函数
			defer func() {
				// 在每个协程内部接收该协程自身抛出来的 panic
				if err := recover(); err != nil {
					log.Logger.Error("defer", err)
				}
			}()
			f() // 业务函数调用执行
		}(t) // 将当前的业务函数名传递给协程
	}
}
