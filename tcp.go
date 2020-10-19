package main

import (
	"log"
	"net"

	"github.com/pborman/uuid"
)

func readConnectionHandlerRoutine(conn net.Conn, id uuid.UUID, sendRequestChannel chan SendRequestData, config Configuration) {
	for {
		var buffer = make([]byte, config.TCPBuffer)
		n, _ := conn.Read(buffer)
		sendRequestChannel <- SendRequestData{
			contentType:  "application/octet-stream",
			data:         buffer[:n],
			destination:  config.TCPToStomp,
			connectionID: id.String(),
		}

		log.Println("tcp read data | local address : " + conn.LocalAddr().String() + ", remote address : " + conn.RemoteAddr().String())
	}
}

func writeConnectionHandlerRoutine(conn net.Conn, writeRequestChannel chan WriteRequestData) {
	for {
		writeRequestData := <-writeRequestChannel
		conn.Write(writeRequestData.data)
		log.Println("write data | local address : " + conn.LocalAddr().String() + ", remote address : " + conn.RemoteAddr().String())
	}
}

//WriteRequestData is struct that writeConnectionHandlerRoutine require to write data
type WriteRequestData struct {
	data []byte
}

func connectionCreatorRoutine(config Configuration, subscribeChannel chan SubscribeRequestData, sendRequestChannel chan SendRequestData) {

	ln, _ := net.Listen("tcp", ":"+config.TCPPort)
	log.Println("TCP server start at " + config.TCPPort)
	for {
		conn, _ := ln.Accept()
		id := uuid.NewUUID()
		writeRequestChannel := make(chan WriteRequestData)
		subscribeChannel <- SubscribeRequestData{
			writeRequestChannel: writeRequestChannel,
			path:                config.StompToTCP,
			selector:            config.ReferenceStompHeader + "='" + id.String() + "'",
		}
		log.Println("new client local address = " + conn.LocalAddr().String() + " | remote address = " + conn.RemoteAddr().String())
		go readConnectionHandlerRoutine(conn, id, sendRequestChannel, config)
		go writeConnectionHandlerRoutine(conn, writeRequestChannel)

	}
}
