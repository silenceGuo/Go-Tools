package redis

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
)

func Client(host, port, pwd string, db int) *redis.Client {
	//ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})
	return rdb
}

func PutChan(ch chan int) {
	for i := 1001; i <= 2000; i++ {
		//err := cli.Set(ctx, fmt.Sprintf("key-%d",i), fmt.Sprintf("value-%d",i), 0).Err
		//m1[fmt.Sprintf("key-%d",i)] = fmt.Sprintf("key-%d",i)
		//fmt.Println(fmt.Sprintf("key-%d",i))
		ch <- i
	}
	//quit <- 0
	defer close(ch)
}
func PutRedis(client *redis.Client, ctx context.Context, ch chan int, quit chan int, wg sync.WaitGroup) {

	for i := range ch {
		err := client.Set(ctx, fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), 0).Err()
		if err != nil {
			//panic(err)
			fmt.Println(fmt.Sprintf("key-%d", i))
		}

	}
	defer wg.Done()
	//for {
	//	select {
	//	case i :=<- ch:
	//		err := client.Set(ctx, fmt.Sprintf("key-%d",i), fmt.Sprintf("value-%d",i), 0).Err()
	//		if err != nil {
	//			//panic(err)
	//			fmt.Println(fmt.Sprintf("key-%d",i))
	//		}
	//
	//	case <-quit:
	//		fmt.Println("quit")
	//		return
	//	}
	//}
	//i:= <- ch

}

func RedisTest() {
	ctx := context.Background()
	runtime.GOMAXPROCS(1)
	//runtime.GOMAXPROCS(runtime.NumCPU())
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	//cli := Client("192.168.254.11", "30001", "q8sXu]%cjR.UWK>o", 2)
	cli := Client("192.168.254.11", "30001", "123456", 2)
	//for i := 1; i <= 100000; i++ {
	//	//err := cli.Set(ctx, fmt.Sprintf("key-%d",i), fmt.Sprintf("value-%d",i), 0).Err()
	//	err:= cli.Get(ctx,fmt.Sprintf("key-%d",i)).Err()
	//	if err != nil {
	//		//panic(err)
	//		fmt.Println(fmt.Sprintf("key-%d",i))
	//	}
	//}
	quit := make(chan int)
	c := make(chan int, 1000)
	PutChan(c)

	go PutRedis(cli, ctx, c, quit, waitGroup)
	go PutRedis(cli, ctx, c, quit, waitGroup)
	//go PutRedis(cli,ctx,c,quit,waitGroup)
	//go PutRedis(cli,ctx,c,quit,waitGroup)
	waitGroup.Wait()
	fmt.Println("=== end ===")
}
