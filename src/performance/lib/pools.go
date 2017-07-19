package lib

import (
	"errors"
	"fmt"
	"log"
)

/*Goroutine池的默认实现*/
type Pools struct {
	total  uint32    //总数
	ch     chan byte //容器
	active bool      //是否激活
}

func NewPools(total uint32) (Pooler, error) {
	pl := Pools{}

	if !pl.init(total) {
		errMsg := fmt.Sprintf("The goroutine ticket pool can not be initialized!total = %d\n", total)
		return nil, errors.New(errMsg)
	}

	return &pl, nil
}

// goroutine池的初始化
func (pl *Pools) init(total uint32) bool {
	if pl.active {
		log.Printf("init:actived,init failure")
		return false
	}

	if total == 0 {
		log.Printf("init:total = %d empty,init failure", total)
		return false
	}

	//初始化缓冲区通道
	ch := make(chan byte, total)
	n := int(total)
	//将通道打满
	for i := 0; i < n; i++ {
		ch <- 1
	}

	pl.ch = ch
	pl.total = total
	pl.active = true

	return true
}

/*实现Pooler接口*/
func (lp *Pools) Take() {
	if !lp.Active() {
		return
	}

	<-lp.ch
}

func (lp *Pools) Return() {
	if !lp.Active() {
		return
	}

	lp.ch <- 1
}

func (lp *Pools) Active() bool {
	return lp.active
}

func (lp *Pools) Total() uint32 {
	return lp.total
}

func (lp *Pools) Remainder() uint32 {
	return uint32(len(lp.ch))
}
