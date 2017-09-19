package main

import (
	"menteslibres.net/gosexy/redis"
	"github.com/VividCortex/godaemon"
	"./gpool"
	"runtime"
	"sync"
	"net"
	"net/http"
	"fmt"
	"io"
	"os"
	"time"
	"strconv"
)

var (
	rdb *redis.Client
	
	httpClient http.Client
	logFile *os.File
	m *sync.Mutex
)

func init() {
	rdb = redis.New()
	rdb.Connect("localhost", 6379)
	rdb.Auth("redispwd")
	rdb.Select(0)

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   8 * time.Second,
			KeepAlive: 8 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   100,		//避免大量TIME_WAIT
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   8 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,		//设置请求5秒超时
	}
	httpClient = http.Client{
	    Transport: transport,
	}

	logFile, _ = os.OpenFile("logs/fetch.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	m = new(sync.Mutex)
}

func main() {
	if len(os.Args) > 1 {
        command := os.Args[1]
        if command == "daemon" {
        	godaemon.MakeDaemon(&godaemon.DaemonAttr{})
        } else if command == "test" {
        	url := "http://mmbiz.qpic.cn/mmbiz_jpg/ZPL7qASgObSkvzr2shDgDxr6KzurtX2lG6Cxo52NC9giaXSRdodkYvBOcn2LR8ydzssaEaib8PaXDnN1clQLFufg/0"
        	key := "files/test.jpg"
        	download(url, key)
        	return
        }
    }
    rcount := 100
    if len(os.Args) > 2 {
    	c, err := strconv.Atoi(os.Args[2])
    	if err == nil {
    		rcount = c
    	}
    }
    fmt.Println(rcount)
    pool := gpool.New(rcount)
    idx := 0
	for {

		ret1, ret2 := rdb.BRPop(3600, "downlist")
		if ret2 == nil {
			url := ret1[1]
			pool.Add(1)
			i, _ := rdb.Incr("filename")
			key :=  fmt.Sprintf("files/%d.jpg", i)
			go func(url, key string) {
				fmt.Println("runtime num:", runtime.NumGoroutine())
				download(url, key)
				pool.Done()
			}(url, key)
			if idx % 10000 == 0 {		//下载一万个文件后打一个标记
				log(strconv.Itoa(runtime.NumGoroutine()))
				log(url)
			}
			idx += 1
		}
	}
}

func download(url, key string) {
	var (
		f *os.File
	)
	start := time.Now()
	resp, err := httpClient.Get(url)
	if err != nil {
		mark("down-err")
		return
	}
	defer resp.Body.Close()
	f, err = os.OpenFile(key, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		mark("down-err")
		return
	}
	io.Copy(f, resp.Body)
	mark("down-ok")
	end := time.Now()
	delta := end.Sub(start).Seconds()*1000
	fmt.Println("down file done:", delta)
}

func mark(key string) {
	today := time.Now().Format("20060102")
	m.Lock()
	rdb.HIncrBy(key, today, 1)
	m.Unlock()
}

func log(msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "[%s]\n%s\n", now, msg)
}