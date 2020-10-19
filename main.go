package main

import (
	"fmt"
)

func main() {
	config, err := getConfig()
	if err != nil {
		fmt.Println("read conffix fail")
		fmt.Println(err.Error())
	}

	subscribeRequestChannel := make(chan SubscribeRequestData)
	sendRequestChannel := make(chan SendRequestData)

	go connectionCreatorRoutine(config, subscribeRequestChannel, sendRequestChannel)
	go StompConnectionCreatorRoutine(config, subscribeRequestChannel, sendRequestChannel)

	for {
	}
}
