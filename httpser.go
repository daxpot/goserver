package main

import (
	"fmt"
	"net/http"
    "net/http/httputil"
    "net/url"
	"log"
	"os"
	"io"
	"github.com/VividCortex/godaemon"
	"strconv"
	"regexp"
    "path/filepath"
    "syscall"
    "bufio"
)

var (
	port = 9090
	path = "."
    blog = false
    proxy = ""
)

func staticRun(w http.ResponseWriter, r *http.Request) {
    var (
        file *os.File
        err error
    )
    r.ParseForm()  //解析参数，默认是不会解析的
    if blog {
        fmt.Println(r.URL.String())
    }

    if proxy != "" {
        remote, _ := url.Parse(proxy)
        proxy_pass := httputil.NewSingleHostReverseProxy(remote)
        proxy_pass.ServeHTTP(w, r)
        return
    }

    p := fmt.Sprintf("%s%s", path, r.URL.Path)
    if p == "." {
    	p = "./"
    }
    if m, _ := regexp.MatchString(".*/$", p); m {
        defaults := [...]string{"index.htm", "index.html", "index"}
        for i := range defaults {
            tp := fmt.Sprintf("%s%s", p, defaults[i])
            file, err = os.Open(tp)
            if err == nil {
                break
            }
        }
    } else {
        file, err = os.Open(p)
    }

    defer file.Close()
    if err != nil {
        w.WriteHeader(404)
    	w.Write([]byte("404 not found"))
    	return
    }
    io.Copy(w, file)
}

func initparam(args []string) {
    if len(args) > 1 {
        path, _ = filepath.Abs(filepath.Dir(path))
        daemon := false
        for i := 1; i < len(args); i++ {
            cmd := args[i]
            if cmd == "-h" || cmd == "--help" {
                help := "Usage: ./httpser [daemon|kill] [-p port] [-r path] [-l] [--proxy url]\n" +
                        "Options:\n daemon\t: daemon run\n" +
                        " kill\t: kill last daemon run\n" +
                        " -p\t: set http listen port default 9090\n" +
                        " -r\t: set root path default ./\n" +
                        " -l\t: set be log" +
                        " --proxy\t: set proxy pass example http://op.kagirl.cn:80"
                fmt.Println(help)
                os.Exit(1)
            } else if cmd == "-p" && len(args) > i {
                p, err:=strconv.Atoi(args[i+1])
                if err == nil {
                    port = p
                }
                i++
            } else if cmd == "daemon" {
                daemon = true
            } else if cmd == "-r" && len(args) > i {
                path = args[i+1]
                i++
            } else if cmd == "--proxy" && len(args) > i {
                proxy = args[i+1]
                i++
            } else if cmd == "-l" {
                blog = true
            } else if cmd == "kill" {
                var (
                    f *os.File
                    pid int
                    pids string
                    err error
                    process *os.Process
                )
                f, err = os.Open(path + "/.httpserpid")
                if err != nil {
                    fmt.Println("该目录未启动过httpser")
                }
                reader := bufio.NewReader(f)
                pids, _ = reader.ReadString('\n')
                pid, err = strconv.Atoi(pids)
                if err != nil {
                    fmt.Println("读取.httpserpid错误")
                }
                process, err = os.FindProcess(pid)
                if err != nil {
                    fmt.Println("进程", pid, "不存在")
                }
                process.Kill()
                fmt.Println("进程", pid, "已杀死")
                os.Remove(path + "/.httpserpid")
                os.Exit(1)
            }
        }

        fmt.Println(path)

        if daemon {
            var (
                f   *os.File
                err error
            )
            if godaemon.Stage() == godaemon.StageParent {
                f, err = os.OpenFile(path + "/.httpserpid", os.O_WRONLY|os.O_CREATE, 0666)
                if err != nil {
                    os.Exit(1)
                }
                err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX)
                if err != nil {
                    os.Exit(1)
                }
            }

            _, _, err = godaemon.MakeDaemon(&godaemon.DaemonAttr{
                Files: []**os.File{&f},
            })

            pid := os.Getpid()
            f.WriteString(fmt.Sprintf("%d", pid))
            f.Close()

        }
    }
}

func main() {
    initparam(os.Args)

	http.HandleFunc("/", staticRun) //设置访问的路由

    err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}