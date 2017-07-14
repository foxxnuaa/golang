package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

// Runner在给定的超时时间内执行一组任务
// 并且在操作系统发送中断信号时结束这些任务
type Runner struct {
	// interrupt 通道报告从操作系统发送的信号
	interrupt chan os.Signal

	// complete 通道报告处理任务已完成
	complete chan error

	// timeout报告处理任务已超时
	timeout <-chan time.Time

	// 一组任务
	tasks []func(int)
}

// ErrTimeout会在任务执行超时时返回
var ErrTimeout error = errors.New("received timeout")

// ErrInterrupt会在接收到操作系统的事件时返回
var ErrInterrupt error = errors.New("received interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		//有缓冲的通道
		interrupt: make(chan os.Signal, 1),
		complete:  make(chan error),
		timeout:   time.After(d),
	}
}

// Add添加一个任务
func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

// Start执行所有任务，并监视通道事件
func (r *Runner) Start() error {
	// 接收中断信号
	signal.Notify(r.interrupt, os.Interrupt)

	// 用不同的goroutine执行不同的任务
	go func() {
		r.complete <- r.run()
	}()

	select {
	// 当任务处理完成时
	case err := <-r.complete:
		return err
		//当任务处理超时时
	case <-r.timeout:
		return ErrTimeout
	}
}

// run执行每一个已注册的任务
func (r *Runner) run() error {
	for id, task := range r.tasks {
		// 检测操作系统的中断信号
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		//执行任务
		task(id)
	}

	return nil
}

// got Interrupt验证是否接收到了中断信号
func (r *Runner) gotInterrupt() bool {
	select {
	//当中断事件发生时
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true

	default:
		return false
	}
}
