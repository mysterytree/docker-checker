package main

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	//docker-client连接
	client *docker.Client
	//socket位置
	endPoint = "tcp://192.168.1.56:2375" //tcp://192.168.1.88:2375
)

func main() {
	//err
	var err error
	//初始化连接
	client, err = docker.NewClient(endPoint)

	if err != nil {
		fmt.Println(err)
	}

	// event channel
	events := make(chan *docker.APIEvents)
	// 监听 events
	client.AddEventListener(events)

	fmt.Println("Listening for Docker events ...")

	select {
	case <-quit:
		return
	default:
	}

}
