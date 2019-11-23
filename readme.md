# 构建镜像

    $ docker build -t go-demo:v1 .

# 运行

    $ docker run -it go-demo:v1

# 压力测试
    
    $ wrk -t 8 -d 2m -c 1000 --timeout 2 --latency http://localhost:8080/hello
    Running 2m test @ http://localhost:8080/hello
      8 threads and 1000 connections
      Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     1.16s   571.43ms   2.00s    63.81%
        Req/Sec    37.57     32.63   380.00     89.07%
      Latency Distribution
         50%    1.30s 
         75%    1.62s 
         90%    1.82s 
         99%    1.98s 
      31987 requests in 2.00m, 7.02GB read
      Socket errors: connect 0, read 0, write 0, timeout 30912
    Requests/sec:    266.34
    Transfer/sec:     59.87MB
    
    
    当压力测试之后，停止程序
    2019/11/22 21:23:07 exit signal:  interrupt
    2019/11/22 21:23:07 http: Server closed
    2019/11/22 21:23:07 profile: caught interrupt, stopping profiles
    2019/11/22 21:23:07 profile: cpu profiling disabled, /tmp/profile028336771/cpu.pprof
    
    对于/tmp/profile028336771/cpu.pprof
    做性能分析
    $ go tool pprof -http=:6060 /tmp/profile028336771/cpu.pprof
    Serving web UI on http://localhost:6060

    就可以在浏览器中看火焰图和cpu,gc等情况
    通过分析cpu,gc发现对于map[xxx]*xxx这样的类型，gc非常频繁
    runtime.scanobject
    /usr/local/go/src/runtime/mgcmark.go
    Total:      10.13s     16.43s (flat, cum) 16.54%
    
    runtime.mallocgc
    /usr/local/go/src/runtime/malloc.go
    
    Total:       740ms      8.47s (flat, cum)  8.53%
    
    runtime.gcBgMarkWorker
    /usr/local/go/src/runtime/mgc.go
    
    Total:           0     16.83s (flat, cum) 16.94%
    查看源码发现gc worker特别频繁
       1895            .          .           		if decnwait == work.nproc { 
       1896            .          .           			println("runtime: work.nwait=", decnwait, "work.nproc=", work.nproc) 
       1897            .          .           			throw("work.nwait was > work.nproc") 
       1898            .          .           		} 
       1899            .          .            
       1900            .     16.79s           		systemstack(func() { 
   
    减少map gc，因为map指针类型，底层采用桶机制存放数据
    app.go 函数中  getUserInfo 返回值改成  map[string]User
    runtime.mallocgc
    /usr/local/go/src/runtime/malloc.go
    
      Total:       2.74s     19.60s (flat, cum)  7.85%
    
    runtime.mallocgc.func1
    /usr/local/go/src/runtime/malloc.go
    
      Total:        20ms      1.65s (flat, cum)  0.66%
      
    压力测试过程中发现，log打印到终端也是枷锁的
    log.(*Logger).Output
    /usr/local/go/src/log/log.go
    
    Total:       7.71s      8.31s (flat, cum) 26.51%
    153        640ms      640ms           	l.mu.Lock() 
    154         80ms       80ms           	defer l.mu.Unlock() 
    所以对于线上来说，一般建议把log.Println这样的语句注释掉
