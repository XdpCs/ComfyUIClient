package main

import (
	"fmt"

	"github.com/XdpCs/comfyUIclient"
)

func main() {
	client := comfyUIclient.NewDefaultClient("hz-t3.matpool.com", "26211")
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
