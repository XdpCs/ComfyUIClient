package main

import (
	"fmt"

	"github.com/XdpCs/comfyUIclient"
)

func main() {
	endPoint := comfyUIclient.NewEndPoint("https", "serverAddress", "port")
	client := comfyUIclient.NewDefaultClient(endPoint)
	client.ConnectAndListen()
	for !client.IsInitialized() {
	}

	info, err := client.GetQueueInfo()
	if err != nil {
		panic(err)
	}
	if len(info.QueueRunning) != 0 {
		fmt.Println(info.QueueRunning[0])
	}
	if len(info.QueuePending) != 0 {
		fmt.Println(info.QueuePending[0])
	}
}
