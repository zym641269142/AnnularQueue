package models

type Task struct {
	CurrentCircleCount int
	CircleCount        int
	Run                func()
}
 