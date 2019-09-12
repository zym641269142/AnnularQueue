package queue

import (
	"AnnularQueue/interface"
	"AnnularQueue/models"
	"errors"
	"fmt"
	"strconv"
	"time"
)

/*******
  管道式消费队列，需要预估业务量
	20190912 增加重试次数设置、增加
	默认重试策略：仅需要重试次数设置，会在当前圈数再次轮到时重试
	自定义重试策略：当重试时间不为0时启用自定义重试策略，设置重试次数和时间，将在经过重试时间后再次重试
 */

const DEFAULT_CHANNEL_BUF = 10

type replay_policy int

const (
	DEFAULT replay_policy = iota
	TIME_REPLAY
)

type QueueByChannel struct {
	CurrentCirclePosition int
	Capacity              int
	//ReplayPolicy          replay_policy
	List []chan *models.Task
}

func NewByChannel(capacity int, meanwhile int) _interface.Queue {
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

//func (this *QueueByChannel) SetReplayPolicy(replayPolicy replay_policy) {
//	this.ReplayPolicy = replayPolicy
//}

func (this *QueueByChannel) AddTask(fun func() error, seconds int, replayCount int, replayTime int) {
	fmt.Println("------加入任务:", fun)
	task := models.Task{}
	task.Run = fun
	task.ReplayCount = replayCount
	task.ReplayTime = replayTime
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
		go this.ReadChan(this.List[i])
		//最后一段循环完之后跳转回第一段，完成一圈环形链
		if i == this.Capacity-1 {
			i = -1
		}
	}
	select {}
}

func (this *QueueByChannel) ReadChan(channel chan *models.Task) {
	//先计算出管道中的任务数量，不能直接使用len(channel)，因为len(channel)是可变量,也不能使用range循环，与数组逻辑不通管道不会先把任务复制一份，是实时提取
	taskCount := len(channel)
	for i := 0; i < taskCount; i++ {
		task := <-channel
		go func() {
			err := this.Execute(task)
			if err != nil {
				this.Replay(task, channel)
			}
		}()
	}
	return
}

func (this *QueueByChannel) Execute(task *models.Task) error {
	//如任务当前圈数==任务触发的圈数 本圈执行
	if task.CircleCount == task.CurrentCircleCount {
		err := task.Run()
		if err != nil {
			fmt.Printf("任务执行错误：%v,等待再次执行\n", err)
			//当前重试次数等于预设重试次数，则返回正确完成任务
			if task.CurrentReplayCount == task.ReplayCount {
				return nil
			}
			//重试次数++
			task.CurrentReplayCount++
			//任务执行失败，重置圈数等待执行
			task.CurrentCircleCount = 1
			return err
		}
	} else {
		//如本圈不执行，圈数加一
		task.CurrentCircleCount ++
		return errors.New("本圈不执行")
	}
	return nil
}

func (this *QueueByChannel) Replay(task *models.Task, channel chan *models.Task) {
	if task.ReplayTime > 0 {
		this.TimeReplay(task)
	} else {
		this.DefaultReplay(task, channel)
	}
}

func (this *QueueByChannel) DefaultReplay(task *models.Task, channel chan *models.Task) {
	channel <- task
}

func (this *QueueByChannel) TimeReplay(task *models.Task) {
	this.AddTask(task.Run, task.ReplayTime, task.ReplayCount-1, task.ReplayTime)
}
