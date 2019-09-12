package queue

import (
	"AnnularQueue/models"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const DEFAULT_CHANNEL_BUF = 10

type QueueByChannel struct {
	CurrentCirclePosition int
	Capacity              int
	List                  []chan *models.Task
}

func NewByChannel(capacity int, meanwhile int) *QueueByChannel {
	queue := &QueueByChannel{}
	if capacity == 0 {
		queue.Capacity = DEFAULT_COUNT
	} else {
		queue.Capacity = capacity
	}

	queue.List = make([]chan *models.Task, queue.Capacity, queue.Capacity)
	for i, _ := range queue.List {
		if meanwhile == 0 {
			queue.List[i] = make(chan *models.Task, DEFAULT_CHANNEL_BUF)
		} else {
			queue.List[i] = make(chan *models.Task, meanwhile)
		}
	}
	return queue
}
func (this *QueueByChannel) AddTask(fun func() error, seconds int) {
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
	this.List[circlePosition-1] <- &task
}

func (this *QueueByChannel) Run() {
	for i := 0; i < this.Capacity; i++ {
		time.Sleep(time.Second)
		this.CurrentCirclePosition = i + 1
		fmt.Println("第" + strconv.Itoa(i+1) + "段")
		go ReadChan(this.List[i])
		//最后一段循环完之后跳转回第一段，完成一圈环形链
		if i == this.Capacity-1 {
			i = -1
		}
	}
	select {}
}

func ReadChan(channel chan *models.Task) {
	for {
		select {
		case task := <-channel:
			go func() {
				err := Execute(task)
				if err != nil {
					channel <- task
				}
			}()
		default:
			return
		}
	}
}

func Execute(task *models.Task) error {
	//如任务当前圈数==任务触发的圈数 本圈执行
	if task.CircleCount == task.CurrentCircleCount {
		err := task.Run()
		if err != nil {
			fmt.Printf("任务执行错误：%v,等待再次执行\n", err)
			//任务执行失败，重置圈数等待执行
			task.CurrentCircleCount = 1
			return err
		}
	} else {
		//如本圈不执行，圈数加一
		task.CurrentCircleCount = task.CurrentCircleCount + 1
		return errors.New("本圈不执行")
	}
	return nil
}