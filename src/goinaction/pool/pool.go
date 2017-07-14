package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

// Pool管理一组可以安全地在多个goroutine间共享的资源。被管理的资源必须实现io.Closer接口
type Pool struct {
	m         sync.Mutex
	resources chan io.Closer
	factory   func() (io.Closer, error)
	closed    bool
}

// ErrPoolClosed表示请求了一个已经关闭了的池
var ErrPoolClosed error = errors.New("Pool has been closed")

func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("Size value too small")
	}

	return &Pool{
		factory:   fn,
		resources: make(chan io.Closer, size),
	}, nil
}

// Acquire从池中获取资源
func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		log.Println("Acquire:", "Shared Resource")
		if !ok {
			return nil, ErrPoolClosed
		}

		return r, nil

		// 无资源可用，创建一个资源
	default:
		log.Println("Acquire:", "New Resource")
		return p.factory()
	}
}

// Release释放资源
func (p *Pool) Release(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	// 如果池已经关闭，销毁这个资源
	if p.closed {
		r.Close()
		return
	}

	select {
	// 释放进入队列
	case p.resources <- r:
		log.Println("Release:", "In Queue")

	// 如果队列已满，则关闭这个资源
	default:
		log.Println("Release:", "Closing")
		r.Close()
	}
}

// 资源池停止工作，并关闭所有现有的资源
func (p *Pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	// 关闭池
	p.closed = true

	//关闭通道
	close(p.resources)

	// 关闭资源，清空通道里的内容
	for r := range p.resources {
		r.Close()
	}
}
