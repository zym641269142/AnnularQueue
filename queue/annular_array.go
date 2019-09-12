package queue

import (
	"AnnularQueue/interface"
	"AnnularQueue/models"
	"fmt"
	"strconv"
	"time"
)
/*********
	单任务执行，多任务将会出现任务删除时的错误
*****/

const DEFAULT_COUNT = 60


type QueueByArray struct {
	CurrentCirclePosition int
	Capacity              int
	List                  [][] *models.Task
}


func New(capacity int) _interface.Queue {
	queue := &QueueByArray{}
	if capacity == 0 {
		queue.Capacity = DEFAULT_COUNT
	} else {
		queue.Capacity = capacity
	}
	queue.List = make([][]*models.Task, queue.Capacity, queue.Capacity)
	for i, _ := range queue.List {
		queue.List[i] = make([]*models.Task, 0)
	}
	return queue
}


func (this *QueueByArray) AddTask(fun func() error, seconds int) {
	fmt.Println("------加入任务:", fun)
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

func (this *QueueByArray) Run() {
	for i := 0; i < this.Capacity; i++ {
		time.Sleep(time.Second)
		this.CurrentCirclePosition = i + 1
		fmt.Println("第" + strconv.Itoa(i+1) + "段")
		for j, tempTask := range this.List[i] {
			go func(task *models.Task, x int, y int) {
				//如任务当前圈数==任务触发的圈数 本圈执行
				if task.CircleCount == task.CurrentCircleCount {
					err := task.Run()
					if err != nil {
						fmt.Printf("任务执行错误：%v,等待再次执行\n", err)
						//任务执行失败，重置圈数等待执行
						task.CurrentCircleCount = 1
						return
					} else {
						//任务执行成功 删除任务
						this.List[x] = append(this.List[x][:y], this.List[x][y+1:]...)
					}
				} else {
					//如本圈不执行，圈数加一
					tempTask.CurrentCircleCount = tempTask.CurrentCircleCount + 1
				}
			}(tempTask, i, j)
		}
		//最后一段循环完之后跳转回第一段，完成一圈环形链
		if i == this.Capacity-1 {
			i = -1
		}
	}
	select {}
}
