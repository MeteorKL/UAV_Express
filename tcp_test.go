package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"testing"
)

// go test -v tcp_test.go tcp.go db.go index.go api.go uav.go
func Test_tcp(t *testing.T) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "localhost:2018")
	checkError(err, "ResolveTCPAddr")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err, "DialTCP")

	_msg := map[string]interface{}{
		"type":     "login",
		"id":       "1",
		"name":     "S2Meteor",
		"password": "MeteorS2",
	}
	msg, err := json.Marshal(_msg)
	if err != nil {
		t.Log(err.Error())
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, int16(len(msg)))
	msgHead := bytesBuffer.Bytes()
	_, err = conn.Write(msgHead)
	conn.Write(msg)

	for {
		_length := make([]byte, 2)
		_, err := conn.Read(_length)
		if checkError(err, "Connection") == false {
			conn.Close()
			fmt.Println("Server is dead ...ByeBye")
			os.Exit(0)
		}
		bytesBuffer := bytes.NewBuffer(_length)
		var length int16
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		buf := make([]byte, length)
		l, err := conn.Read(buf)
		if checkError(err, "Connection") == false {
			conn.Close()
			fmt.Println("Server is dead ...ByeBye")
			os.Exit(0)
		}
		if l == 0 {
			continue
		}
		fmt.Println(string(buf))
	}
}
