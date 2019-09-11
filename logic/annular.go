package logic

import (
	"AnnularQueue/models"
	"fmt"
	"time"
)

const DEFAULT_COUNT = 60

type Queue struct {
	CurrentCirclePosition int
	Capacity              int
	List                  [][]*models.Task
}

func New() Queue {
	queue := Queue{}
	if queue.Capacity == 0 {
		queue.Capacity = DEFAULT_COUNT
	}
	queue.List = make([][]*models.Task, queue.Capacity, queue.Capacity)
	for i, _ := range queue.List {
		queue.List[i] = make([]*models.Task, 0)
	}
	return queue
}

//func (this Queue) SetCapacity(capacity int) {
//	this.Capacity = capacity
//}

func (this Queue) AddTask(fun func(), seconds int) {
	task := models.Task{}
	task.Run = fun
	//初始化当前的圈数
	task.CurrentCircleCount = 1
	//所需要的圈数
	circleCount := seconds/this.Capacity + 1
	task.CircleCount = circleCount
	var circlePosition = 1
	//圈位置
	if this.CurrentCirclePosition+seconds < this.Capacity {
		circlePosition = this.CurrentCirclePosition + seconds
	} else {
		circlePosition = (this.CurrentCirclePosition + seconds) % this.Capacity
	}
	this.List[circlePosition-1] = append(this.List[circlePosition-1], &task)
}

func (this Queue) Run() {
	for i := 0; i < this.Capacity; i++ {
		time.Sleep(time.Second)
		this.CurrentCirclePosition = i + 1
		//fmt.Println("第" + strconv.Itoa(i+1) + "位")
		for j, tempTask := range this.List[i] {
			go func(task *models.Task, x int, y int) {
				if task.CircleCount == task.CurrentCircleCount {
					go func() {
						task.Run()
						this.List[x] = append(this.List[x][:y], this.List[x][y+1:]...)
						fmt.Println()
					}()
				} else {
					tempTask.CurrentCircleCount = tempTask.CurrentCircleCount + 1
				}
			}(tempTask, i, j)
		}
		if i == this.Capacity-1 {
			i = -1
		}
	}
	select {}
}
