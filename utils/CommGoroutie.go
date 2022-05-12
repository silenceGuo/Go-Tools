package utils

import (
	"fmt"
	"sync"
)

type Mytype any

type MyGoroutie struct {
	//WorkFunc func(...Mytype) Mytype
	Num      int
	WorkerCh chan Mytype
	Wg       *sync.WaitGroup
	Args     Mytype
}

func (this *MyGoroutie) MyRun(arg any) error {
	//for _, v := range args {
	//	fmt.Println(v)
	//	fmt.Printf("%T\n", v)
	//}
	//fmt.Println("ss")
	//this.WorkFunc(args)
	return nil
}

func Work(args ...any) any {
	for _, v := range args {
		fmt.Println("workeing...", v)
	}
	return nil
}
func MyInit() {
	g := MyGoroutie{}
	var wg sync.WaitGroup
	wch := make(chan Mytype, 10)
	g.Wg = &wg
	g.Num = 3
	g.WorkerCh = wch
	//g.WorkFunc = Work

	go g.MyRun(Work("arglist", "1", 1, 1.3))
	fmt.Println("sss")
}
