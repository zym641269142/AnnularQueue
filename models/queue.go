package models

type Task struct {
	CurrentCircleCount int
	CurrentReplayCount int
	CircleCount        int
	ReplayCount        int
	ReplayTime         int //second
	Run                func() error
}
