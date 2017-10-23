package mypool
import (
	"testing"
	"runtime"
	"os"
	"os/signal"
	"time"
)

var (
	killed = false
)

func TestNew0(t *testing.T) {
	pool := New(-1)

	t.Log("routine数", runtime.NumGoroutine())
	for i := 0; i < 10 && !killed; i++ {
		pool.Add(func(args ...interface{}) {
			param := args[0].([]interface{})
			idx := param[0].(int)
			str := param[1].(string)
			i := args[1].(int)
			
			t.Logf("序号%d,channel:%d,参数:%s,routine数:%d,quene len:%d", idx, i, str, runtime.NumGoroutine(), pool.Length())
			
		}, i, "test")
	}
	pool.Wait()
	t.Log("通过")
}

func TestNew(t *testing.T) {
	pool := New(3)
	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs)
		sig := <-sigs
		t.Log("Got signal:", sig, pool.Length())
		killed = true
	}()
	t.Log("routine数", runtime.NumGoroutine())
	for i := 0; i < 10 && !killed; i++ {
		pool.Add(func(args ...interface{}) {
			param := args[0].([]interface{})
			idx := param[0].(int)
			str := param[1].(string)
			i := args[1].(int)
			
			t.Logf("序号%d,channel:%d,参数:%s,routine数:%d,quene len:%d", idx, i, str, runtime.NumGoroutine(), pool.Length())
			time.Sleep(time.Second*1)
		}, i, "test")
	}
	pool.Wait()
	t.Log("通过")
}