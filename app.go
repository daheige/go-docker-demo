package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/daheige/thinkgo/monitor"
	"github.com/pkg/profile"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	promeHandler()
}

func main() {
	defer profile.Start().Stop()

	//启动http服务
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("/hello", hello)

	port := 8080
	server := &http.Server{
		Handler:      mux,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		ReadTimeout:  5 * time.Second,  //read request timeout
		WriteTimeout: 10 * time.Second, //write timeout
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	//平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig := <-ch

	log.Println("exit signal: ", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx) //在独立的携程中关闭服务器
	<-ctx.Done()

	log.Println("shutting down")
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println(111)
	userInfo := getUserInfo()

	for k, _ := range userInfo {
		log.Println(userInfo[k].Id, userInfo[k].Name, userInfo[k].Age)
	}

	b, _ := json.Marshal(userInfo)

	time.Sleep(10 * time.Millisecond)

	w.Write(b)
}

type User struct {
	Id      int64
	Name    string
	Age     int
	Content string
}

func getUserInfo() map[string]User {
	user := make(map[string]User, 500)
	str := `What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
	What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
What to profile is controlled by config value passed to profile.Start. By default CPU profiling is enabled.
`

	for i := 0; i < 100; i++ {
		nick := "hello_" + strconv.Itoa(i)
		user[nick] = User{
			Id:      int64(i),
			Name:    nick,
			Age:     i + 10,
			Content: str,
		}
	}

	return user
}

func promeHandler() {
	//注册监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	//性能监控的端口port+1000,只能在内网访问
	go func() {
		pprof_port := 2338
		log.Println("server pprof run on: ", pprof_port)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", HealthCheckHandler)

		//metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", pprof_port), httpMux); err != nil {
			log.Println(err)
		}
	}()
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	w.Write([]byte(`{"alive": true}`))
}
