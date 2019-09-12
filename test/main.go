package main

import (
	"AnnularQueue/queue"
	"errors"
	"fmt"
)

func main() {

	queue := queue.NewByChannel(60, 0)
	go func() {
		queue.AddTask(Fmt, 3,2,0)
		//for i := 0; i < 5; i++ {
		//	queue.AddTask(Fmt2, 2,1)
		//}
		//time.Sleep(time.Second*10)
		//queue.AddTask(Fmt4,5)
	}()
	//queue.AddTask(Fmt2, 3)
	//queue.AddTask(Fmt3, 61)
	queue.Run()
}

func Fmt() error {
	fmt.Println("第一个任务输出了，哈哈哈哈哈")
	return errors.New("第一个任务报错了")
}

func Fmt2() error {
	fmt.Println("第二个任务输出了，哈哈哈哈哈")
	return nil
}

func Fmt3() error {
	fmt.Println("第三个任务输出了，哈哈哈哈哈")
	return nil
}

func Fmt4() error {
	fmt.Println("运行中插入的任务4输出了，哈哈哈哈哈")
	return nil
}
