package main

import (
	"bytes"
	"cstructs"
	"fmt"
	"github.com/franela/goreq"
	docker "github.com/fsouza/go-dockerclient"
	"strings"
	"time"
)

var (
	//docker-client连接
	client *docker.Client
	//socket位置
	endPoint = "unix:///tmp/docker.sock" //tcp://192.168.1.88:2375
	//serverhost
	serverUrl string
)

func main() {
	//初始化
	serverUrl = "http://" + cstructs.GetEnv("SERVER_URL", "")
	//err
	var err error
	//初始化连接
	client, err = docker.NewClient(endPoint)

	if err != nil {
		panic(err)
	}

	// event channel
	events := make(chan *docker.APIEvents)
	// 监听 events
	client.AddEventListener(events)

	fmt.Println("Listening for Docker events ...")

	//接收并转发消息
	for msg := range events {
		sendMessage(msg)
	}
}

// 发送消息给 Server
func sendMessage(msg *docker.APIEvents) {
	//获取name
	inspectContainer, err := client.InspectContainer(msg.ID)

	if err != nil {
		fmt.Println(err)
		return
	}

	var out bytes.Buffer
	var errout bytes.Buffer

	//配置log  opinion
	logOpinion := docker.LogsOptions{
		Container:    msg.ID,
		OutputStream: &out,
		ErrorStream:  &errout,
		Follow:       false,
		Stdout:       true,
		Stderr:       true,
		Tail:         "5",
	}

	//查看log
	logError := client.Logs(logOpinion)

	if logError != nil {
		fmt.Println(logError)
		return
	}

	//发送
	send(strings.Split(inspectContainer.Name, "/")[1], msg.Status, getTime(msg.Time), out.String()+errout.String())
}

// server 通信
func send(name string, status string, time string, logs string) {
	//组装发送的container消息
	container := cstructs.Container{
		Name:   name,
		Status: status,
		Time:   time,
		Log:    logs,
	}

	//发送
	res, err := goreq.Request{
		Method: "POST",
		Uri:    serverUrl,
		Body:   container,
	}.Do()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}

func getTime(timeStamp int64) string {
	tm := time.Unix(timeStamp, 0)
	return tm.Format("2006-01-02 03:04:05 PM")
}
