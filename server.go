package main

import (
	"fmt"
	"runtime"
	"time"
	"os"
	"os/signal"
	"./mypool"
)

var (
	killed = false
)

func test(args ...interface{}) {
	param := args[0].([]interface{})
	idx := param[0].(int)
	t := param[1].(string)
	i := args[1].(int)
	time.Sleep(time.Second*1)
	fmt.Println("runtime num:", runtime.NumGoroutine(), idx, t, i)
}

func signalListen(pool *mypool.Pool) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs)
	sig := <-sigs
	fmt.Println("Got signal:", sig, pool.Length())
	killed = true
}

func main() {
	pool := mypool.New(3)
	go signalListen(pool)		//安全退出
	for i := 0; i<10 && !killed; i++ {
		pool.Add(test, i, "test")
	}
	pool.Wait()
	fmt.Println("done")
}