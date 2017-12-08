#!/bin/bash
#start $ClosePro
#2013/1/23 李秀旭 
#QQ:76224980
  
ClosePro=./MonitoringSys

pid=`ps -ef|grep $ClosePro |grep -v grep|awk '{print $2}'`
if [ "x$pid" = "x" ] ; then
    echo "$ClosePro 没有启动！"
else
    kill -9 $pid
    sleep 1
    pid1=`ps -ef|grep $ClosePro|grep -v grep|awk '{print $2}'`
    if [ "x$pid1" = "x" ] ; then
            echo "成功杀死$ClosePro进程：" $pid
        else
            echo "$ClosePro进程杀死失败！"
            exit 1
    fi
fi
$ClosePro  1>> log.out 2>&1 &
pid2=`ps -ef|grep $ClosePro|grep -v grep|awk '{print $2}'`
if [ "x$pid2" = "x" ] ; then
        echo "$ClosePro服务启动失败！"
    else
        echo "$ClosePro服务成功启动:" $pid2
fi

