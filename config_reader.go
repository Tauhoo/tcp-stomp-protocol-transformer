package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Configuration is configured struct
type Configuration struct {
	StompHost            string `json:"stomp_host"`
	TCPPort              string `json:"tcp_port"`
	TCPBuffer            int    `json:"tcp_buffer"`
	TCPToStomp           string `json:"tcp_to_stomp"`
	StompToTCP           string `json:"stomp_to_tcp"`
	ReferenceStompHeader string `json:"reference_stomp_header"`
}

func getConfig() (Configuration, error) {
	jsonFile, err := os.Open("config.json")
	if err != nil {
		return Configuration{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result Configuration
	json.Unmarshal([]byte(byteValue), &result)
	return result, nil

}
