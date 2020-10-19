package main

import (
	"log"

	"github.com/go-stomp/stomp"
)

//SubscribeRequestData is data that require for stomp subscriber use for subscribe
type SubscribeRequestData struct {
	path                string
	selector            string
	writeRequestChannel chan WriteRequestData
}

//SendRequestData is data required for send data to stomp server
type SendRequestData struct {
	data         []byte
	destination  string
	contentType  string
	connectionID string
}

func stompSubscriptionHandlerRoutine(subscription *stomp.Subscription, writeRequestChannel chan WriteRequestData) {
	defer subscription.Unsubscribe()
	log.Println("new stomp subscription handler routine is up")
	for {
		stompMessage, _ := <-subscription.C
		writeRequestChannel <- WriteRequestData{
			data: stompMessage.Body,
		}
		log.Println("send write request")
	}
}

func stompSendingHandlerRoutine(conn *stomp.Conn, sendRequestChannel chan SendRequestData, config Configuration) {
	log.Println("stomp sending handler routine is up")
	for {
		sendRequestData := <-sendRequestChannel
		conn.Send(sendRequestData.destination, sendRequestData.contentType, sendRequestData.data, stomp.SendOpt.Header(config.ReferenceStompHeader, sendRequestData.connectionID))
		log.Println("send data to stomp server")
	}
}

//StompConnectionCreatorRoutine is used to create stomp subscribe handler
func StompConnectionCreatorRoutine(config Configuration, subscribeRequestChannel chan SubscribeRequestData, sendRequestChannel chan SendRequestData) {
	conn, _ := stomp.Dial("tcp", config.StompHost)
	log.Println("Stomp connected at " + config.StompHost)
	defer conn.Disconnect()
	go stompSendingHandlerRoutine(conn, sendRequestChannel, config)
	for {
		subscribeRequestData := <-subscribeRequestChannel
		log.Println("new subscription", subscribeRequestData)
		headerSelector := stomp.SubscribeOpt.Header("selector", subscribeRequestData.selector)
		subscription, err := conn.Subscribe(subscribeRequestData.path, stomp.AckClientIndividual, headerSelector)
		if err != nil {
			log.Println("subscribe fail \n" + err.Error())
			continue
		}
		go stompSubscriptionHandlerRoutine(subscription, subscribeRequestData.writeRequestChannel)
	}
}
