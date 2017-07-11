package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + ":  " + err.Error())
		return false
	}
	return true
}

const (
	STAUTS_TCP_USER_UNLOGIN = iota
	STAUTS_TCP_USER_LOGIN   = iota
)

type TcpUser struct {
	conn   net.Conn
	u      *User
	msg    chan map[string]interface{}
	status int
}

var conns = make(map[string]*TcpUser)
var loginedConns = make(map[int]*TcpUser)

func tcpServer(port string) {
	service := ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err, "ResolveTCPAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	//启动服务器广播线程
	// go echoHandler(&conns)

	for {
		fmt.Println("Listening ...")
		conn, err := l.Accept()
		checkError(err, "Accept")
		fmt.Println("Accepting ...")
		tcpUser := &TcpUser{
			conn:   conn,
			msg:    make(chan map[string]interface{}, 10),
			status: STAUTS_TCP_USER_UNLOGIN,
		}
		conns[conn.RemoteAddr().String()] = tcpUser
		go tcpUser.write()
		go tcpUser.read()
		//启动一个新线程
		// go Handler(conn, messages)
	}
}

func (tcpUser *TcpUser) write() {
	for {
		_msg := <-tcpUser.msg
		if tcpUser.u != nil {
			fmt.Printf("send msg to %s: %s\n", tcpUser.u.User_name, _msg)
		}
		msgBody, _ := json.Marshal(_msg)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, int16(len(msgBody)))
		msgHead := bytesBuffer.Bytes()
		_, err := tcpUser.conn.Write(msgHead)
		if err != nil {
			break
		}
		_, err = tcpUser.conn.Write(msgBody)
		if err != nil {
			break
		}
	}
	delete(conns, tcpUser.conn.RemoteAddr().String())
	if tcpUser.u != nil {
		if _, ok := loginedConns[tcpUser.u.User_id]; ok {
			loginedConns[tcpUser.u.User_id] = nil
		}
	}
}

func sendMsg2User(user_id int, msg string) {
	if c, ok := loginedConns[user_id]; ok && c != nil {
		c.msg <- map[string]interface{}{
			"type": "msg",
			"data": msg,
		}
	}
}

func (tcpUser *TcpUser) read() {
	for {
		_length := make([]byte, 2)
		_, err := tcpUser.conn.Read(_length)
		if checkError(err, "Connection") == false {
			break
		}
		bytesBuffer := bytes.NewBuffer(_length)
		var length int16
		binary.Read(bytesBuffer, binary.BigEndian, &length)
		buf := make([]byte, length)
		l, err := tcpUser.conn.Read(buf)
		if checkError(err, "Connection") == false {
			break
		}
		if l == 0 {
			continue
		}
		fmt.Printf("read msg: %s\n", buf)
		var msg map[string]interface{}
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			fmt.Println("json: ", err)
			continue
		}
		if t, ok := msg["type"]; !ok {
			log.Println("msg json should have type")
		} else {
			switch t {
			case "login":
				id, err := strconv.Atoi(msg["id"].(string))
				if err != nil {
					tcpUser.msg <- map[string]interface{}{
						"type": "error",
						"data": err.Error(),
					}
					continue
				}
				u := getUserById(id)
				// if u.User_name == msg["name"] && u.User_password == msg["password"] {
				tcpUser.u = u
				loginedConns[u.User_id] = tcpUser
				tcpUser.msg <- map[string]interface{}{
					"type": "login",
					"data": "login success",
				}
				if u.User_balance < 0 {
					tcpUser.msg <- map[string]interface{}{
						"type": "msg",
						"data": "当前账户余额不足，请及时充值",
					}
				}
				// } else {
				// 	tcpUser.msg <- map[string]interface{}{
				// 		"type": "login",
				// 		"data": "login failed",
				// 	}
				// }
			default:
			}
		}
		// fmt.Println("Rec[",conn.RemoteAddr().String(),"] Say :" ,string(buf[0:lenght]))
		// reciveStr := string(buf[0:lenght])
		// messages <- reciveStr
	}
	close(tcpUser.msg)
	tcpUser.conn.Close()
}
