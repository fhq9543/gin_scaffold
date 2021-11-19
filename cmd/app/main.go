package main

import (
	"baseFrame/app"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	injector, injectorCleanFunc, err := app.BuildInjector()
	if err != nil {
		fmt.Printf("injector error:" + err.Error())
		panic(err)
	}

	cleanFunc, err := injector.Init()
	if err != nil {
		fmt.Printf("init failed, err:%v\n", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 监听服务
	go injector.Start(ctx)

	// 等待结束
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 清理
	{
		injectorCleanFunc()
		cleanFunc()
	}

	// 结束服务
	injector.Stop(ctx)
}
