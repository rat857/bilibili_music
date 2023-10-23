package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建一个背景context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	timeStart := time.Now() // 记录启动时间
	go func() {
		for {
			select {
			case <-ctx.Done():
				// context被取消，退出goroutine
				return
			default:
				// 检查时间是否超过10秒
				if time.Since(timeStart) > 5*time.Second {
					fmt.Println("Time exceeded 5 seconds in goroutine.")
					cancel() // 取消context，触发主函数退出
				}
			}
		}
	}()

	//timeStart := time.Now() // 记录启动时间
	i := 1
	for {
		select {
		case <-ctx.Done():
			// 主函数会在goroutine中取消context后退出
			fmt.Println("Main function exiting.")
			return
		default:
			fmt.Println(i, "helloworld")
			i = i + 1
			time.Sleep(500 * time.Millisecond)
		}
	}
}
