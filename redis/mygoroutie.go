package redis

import (
	"context"
	"fmt"

	"sync"
)

// goroutine 的例子
type GroutinTools struct {
	Num      int
	WorkerCh chan interface{}
	Wg       *sync.WaitGroup
	WorkerF  func(interface{})
	ctx      context.Context
	cancle   context.CancelFunc
	//ProudF func()
}

func (this *GroutinTools) Produf() {
	for i := 0; i < 1000; i++ {
		fmt.Println("生产数据：", i)
		this.WorkerCh <- i
	}
	close(this.WorkerCh)
	fmt.Println("生产数据完成。")
	this.cancle()
}

func (this *GroutinTools) Init() {
	ch := make(chan interface{}, 10)
	var wg sync.WaitGroup
	this.Wg = &wg
	this.WorkerCh = ch
	this.ctx, this.cancle = context.WithCancel(context.Background())
	//d.Done()
}

func (this *GroutinTools) Workerf(i int) {
	defer this.Wg.Done()
	for {
		v, ok := <-this.WorkerCh
		if !ok {
			fmt.Println("完成2")
			break
		}
		fmt.Printf("%d get %d \n", i, v)
	}

}

func (this *GroutinTools) GoroutinTools() {

	for i := 0; i < this.Num; i++ {
		fmt.Println("协成启动：", i)
		this.Wg.Add(1)
		//go this.Workerf(i)
		go this.Worke(i)
	}
	//this.cancle()
	this.Wg.Wait()
}

func (this *GroutinTools) Worke(i int) {
	defer this.Wg.Done()
	//Lable1:
	for {
		v, ok := <-this.WorkerCh
		//fmt.Printf("%d get %d \n", i, v)
		//select {
		//case <-this.ctx.Done():
		//	break Lable1
		////default:
		////case !ok:
		//	this.cancle()
		//}
		if !ok {
			fmt.Println("完成2")
			this.cancle()
			break
		}

		fmt.Printf("%d get %d \n", i, v)

	}

}

func Testg() {
	// 测试函数
	G := GroutinTools{
		Num:      8,
		WorkerCh: nil,
		Wg:       nil,
	}
	//ctx, _ := context.WithCancel(context.Background())
	//runtime.GOMAXPROCS(2)
	//this.WorkerF=TestWork
	G.Init()
	//go G.Produf(ctx)
	go G.Produf()
	G.GoroutinTools()

}
