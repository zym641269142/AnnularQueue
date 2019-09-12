package main

import (
	"AnnularQueue/queue"
	"errors"
	"fmt"
	"time"
)

func main() {

	queue := queue.New(60)
	go func() {
		time.Sleep(time.Second*10)
		queue.AddTask(Fmt4,5)
	}()
	for i:=0;i<100;i++{
		queue.AddTask(Fmt2, 2)
	}
	queue.AddTask(Fmt, 3)
	queue.AddTask(Fmt3, 61)
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