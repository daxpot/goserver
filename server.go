package main

import (
	"fmt"
	"runtime"
	"time"
	"./mypool"
)

func test(args ...interface{}) {
	param := args[0].([]interface{})
	idx := param[0].(int)
	t := param[1].(string)
	i := args[1].(int)
	time.Sleep(time.Second*1)
	fmt.Println("runtime num:", runtime.NumGoroutine(), idx, t, i)
}

func main() {
	pool := mypool.New(3)
	for i := 0; i<10; i++ {
		pool.Add(test, i, "test")
	}
	pool.Wait()
	fmt.Println("done")
}