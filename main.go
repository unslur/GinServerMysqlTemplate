// MyServerTemplate project main.go
package main

import (
	"context"
	"fmt"
	"myfunc"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

//Server 读取服务器配置文件
var Server myfunc.ServerConfig
var cylog seelog.LoggerInterface
var db *sqlx.DB

//ConfigPath 必要配置文件存放路径
var ConfigPath = "./config/config.ini"

// SeelogConfigPath 日志文件设置存放路径
var SeelogConfigPath = "./config/cy_seelog.xml"

//GinLogPath 日志输出文件路径
var GinLogPath = "./log/cy_seelog.log"

func main() {
	/////////////////初始化日志输出以及数据库各项配置
	var err error
	cylog, err = seelog.LoggerFromConfigAsFile(SeelogConfigPath)
	if err != nil {
		fmt.Println("读取 seelog.xml 日志配置文件错误,请确定格式:", err.Error())
		return

	}

	Server = myfunc.ServerConfig{}
	myfunc.InitConfig(&Server, ConfigPath)
	cylog.Info("=======服务器端口", Server.ServerPort)
	cylog.Info("=======数据库地址:", Server.DBIp)
	cylog.Info("=======数据库端口:", Server.DBPort)
	cylog.Info("=======数据库用户:", Server.DBUser)
	cylog.Info("=======数据库密码:", Server.DBPassWd)
	cylog.Info("=======数据库名称:", Server.DBName)
	//////////////////连接数据库

	mysqlinfo := fmt.Sprintf(`%s:%s@tcp(%s:%s)/%s?charset=utf8`, Server.DBUser, Server.DBPassWd, Server.DBIp, Server.DBPort, Server.DBName)

	db = sqlx.MustConnect("mysql", mysqlinfo)
	db.DB.SetMaxIdleConns(10)
	db.DB.SetMaxOpenConns(40)
	err = db.Ping()
	if err != nil {
		cylog.Error("连接数据库错误:", err.Error())
		return
	}
	///////////////启动http服务
	httpInterface()
}
func httpInterface() {
	gin.SetMode(gin.ReleaseMode)

	if runtime.GOOS == "windows" {

	} else if os.Getenv("GOEDIT") == "cyy" {

	} else {

		gin.DefaultWriter = &lumberjack.Logger{
			Filename:   GinLogPath,
			MaxAge:     28,
			MaxBackups: 3,
			MaxSize:    5000, // megabytes
			Compress:   true,
		}
		gin.DefaultErrorWriter = &lumberjack.Logger{
			Filename:   GinLogPath,
			MaxAge:     28,
			MaxBackups: 3,
			MaxSize:    5000, // megabytes
			Compress:   true,
		}
	}
	router := gin.Default()
	router.Use(myfunc.Middleware())
	V2 := router.Group("/cry")
	{
		V2.Any("/QueryLastPlace", queryLastPlace)
	}
	srv := &http.Server{
		Addr:    Server.ServerPort,
		Handler: router,
	}

	go func() {

		if err := srv.ListenAndServe(); err != nil {
			cylog.Errorf("listen: %s\n", err)
		}
		cylog.Info("server shutdown")
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	sig := <-quit
	cylog.Info("Shutdown Server by signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		cylog.Error("Server Shutdown:", err)
	}
	cylog.Info("Server exiting")

}
