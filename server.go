package main

import (
	"fmt"
	"runtime"
	"time"
	"./mypool"
)

func test(args ...interface{}) {
	time.Sleep(time.Second*1)
	param := args[0].([]interface{})
	idx := param[0].(int)
	t := param[1].(string)
	fmt.Println("runtime num:", runtime.NumGoroutine(), idx, t)
	// for _, arg := range args {
	// 	fmt.Println(arg.(type))
	// }
}

func main() {
	pool := mypool.New(2)
	for i := 0; i<10; i++ {
		pool.Add(test, i, "test")
	}
	pool.Wait()
	fmt.Println("done")
}