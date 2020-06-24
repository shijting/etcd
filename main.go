package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shijting/etcd/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main()  {
	// 第三方路由库
	route := mux.NewRouter()
	//route.HandleFunc("/prod/{id:\\d+}", func(writer http.ResponseWriter, request *http.Request) {
	//	vars := mux.Vars(request)
	//	str := "get product by id:" + vars["id"]
	//	writer.Write([]byte(str))
	//})

	service := util.NewService()
	serviceId := "s1"
	serviceName := "product1"
	port := 8080
	serviceAddr := "192.168.0.1"
	// 启动http 服务
	httpServer := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           route,
	}
	// 用于通知错误处理
	errNotify := make(chan error)
	go func() {
		// 注册服务
		err := service.RegService(serviceId, serviceName, serviceAddr +":" + strconv.Itoa(port))
		if err != nil {
			errNotify <- err
			return
		}
		err = httpServer.ListenAndServe()
		if err != nil {
			errNotify <- err
			return
		}
	}()
	// 监听系统信号 如：ctrl + c ， kill -9
	go func() {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		// 监听到信号
		errNotify <- fmt.Errorf("%s", <- sig)

	}()

	getErr := <- errNotify
	log.Println("服务停止...")
	// 服务出现异常，进行反注册
	err := service.UnRegService(serviceId)
	if err != nil {
		log.Fatal(err)
	}

	// 关闭http服务器，关闭之前可以执行一些回收工作，如关闭数据库...
	err = httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(getErr)
}
