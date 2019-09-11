package main

import (
	"AnnularQueue/logic"
	"fmt"
)

func main() {
	queue := logic.New()
	queue.AddTask(Fmt, 61)
	queue.Run()
}

func Fmt() {
	fmt.Println("hahahaha")
}
