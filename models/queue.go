package models

type Task struct {
	CurrentCircleCount int
	CurrentReplayCount int
	CircleCount        int
	ReplayCount        int
	Run                func() error
}
