package _interface

type Queue interface {
	AddTask(fun func() error, seconds int)
	Run()
}
