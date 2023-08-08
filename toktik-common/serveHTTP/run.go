package serveHTTP

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(r *gin.Engine, srvName string, addr string, stop func()) {
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// 保证下面的优雅启停
	go func() {
		log.Printf("%s server running in %s \n", srvName, server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("hah")
			log.Fatalln(err)
		}
	}()

	fmt.Println("----------start----------")
	quit := make(chan os.Signal)
	// SIGINT 用户发送INTR字符(Ctrl+C)触发 kill -2
	// SIGTERM 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutting Down project %s...\n", srvName)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if stop != nil {
		stop()
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("%s server Shutdown,causer by %v: \n", srvName, err)
	} else if err == context.DeadlineExceeded {
		log.Fatalln("服务器关闭超时")
	}
	select {
	case <-ctx.Done():
		log.Println("wait timeout...")
	}
	log.Printf("%s server stop success...\n", srvName)
}
