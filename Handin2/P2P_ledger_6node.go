package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Message struct {
	MsgType    string
	MsgContent string
}

type PeerInfo struct {
	Addr string
	Port int
}

type Peer struct {
	Addr  string
	Port  int
	Peers []PeerInfo
}

func (p *Peer) HandleConnection(conn net.Conn) {
	defer conn.Close()
	otherEnd := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			p.log("Ending session with " + otherEnd)
			return
		} else {
			receivedMessage := new(Message)
			json.Unmarshal(msg, receivedMessage)
			p.log("received : " + string(msg))
			responseMessage := new(Message)
			switch receivedMessage.MsgType {
			case "askPeersInfo":
				selfPeers, _ := json.Marshal(p.Peers)
				responseMessage.MsgType = "askPeersInfoResponse"
				responseMessage.MsgContent = string(selfPeers)
				responseData, _ := json.Marshal(responseMessage)
				//reply
				writer := bufio.NewWriter(conn)
				writer.Write(append(responseData, '\n'))
				writer.Flush()
			case "joinMessage":
				receivedJoinPeer := new(PeerInfo)
				json.Unmarshal([]byte(receivedMessage.MsgContent), receivedJoinPeer)
				var receivedJoin []PeerInfo
				receivedJoin = append(receivedJoin, *receivedJoinPeer)
				p.recordPeers(receivedJoin)
			}

		}
	}
}
func (p *Peer) Listen() {
	listen_port := ":" + strconv.Itoa(p.Port)
	ln, _ := net.Listen("tcp", listen_port)
	defer ln.Close()
	p.log("now listening for connection...")
	for {
		conn, _ := ln.Accept()
		p.log("Got a connection...")
		go p.HandleConnection(conn)
	}
}

func (p *Peer) Connect(addr string, port int) {
	// adding self to peers
	p.Peers = append(p.Peers, PeerInfo{p.Addr, p.Port})
	// bind local port and connect to target
	target := addr + ":" + strconv.Itoa(port)
	// netAddr := &net.TCPAddr{Port: p.Port}
	// d := net.Dialer{LocalAddr: netAddr}
	conn, err := net.Dial("tcp", target)
	if err != nil {
		// address or port invalid
		// start its own network
		p.log("address or port invalid, starting own network...")
		go p.Listen()
		return
	}
	// defer conn.Close()
	p.askPeersInfo(conn)
	conn.Close()
	// time.Sleep(70 * time.Second)

	p.floodJoin()
	go p.Listen()
}

func (p *Peer) floodJoin() {
	for _, selfPeerInfo := range p.Peers {
		if p.Addr == selfPeerInfo.Addr && p.Port == selfPeerInfo.Port {
			continue
		}
		target := selfPeerInfo.Addr + ":" + strconv.Itoa(selfPeerInfo.Port)
		p.log("target: " + target)
		// netAddr := &net.TCPAddr{Port: p.Port}
		// d := net.Dialer{LocalAddr: netAddr}
		conn, err := net.Dial("tcp", target)
		if err != nil {
			// address or port invalid
			// start its own network
			p.log("floodJoin falied" + err.Error())

			continue
		}
		defer conn.Close()

		joinPeerInfo := PeerInfo{p.Addr, p.Port}
		join_data, _ := json.Marshal(joinPeerInfo)
		joinInfoMsg := Message{"joinMessage", string(join_data)}

		data, _ := json.Marshal(joinInfoMsg)
		writer := bufio.NewWriter(conn)
		writer.Write(append(data, '\n'))
		writer.Flush()
		p.log("send joinMessage")
	}
}

func (p *Peer) log(content string) {
	if content == "updated peers info" {
		fmt.Println("IP:", p.Addr, "port:", p.Port, ": ", content, " ", p.Peers)
	} else {
		fmt.Println("IP:", p.Addr, "port:", p.Port, ": ", content)
	}
}

func (p *Peer) askPeersInfo(conn net.Conn) {
	// send askPeersInfo msg to get peers set
	// use json to marshal
	askInfoMsg := Message{"askPeersInfo", "hello"}
	data, _ := json.Marshal(askInfoMsg)
	writer := bufio.NewWriter(conn)
	writer.Write(append(data, '\n'))
	writer.Flush()
	p.log("send askPeersInfo msg")
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadBytes('\n')
	if err != nil {
		p.log("askPeersInfo failed")
		return
	} else {
		p.log("askPeersInfo success")
		askPeersInfoResponse := new(Message)
		json.Unmarshal(msg, askPeersInfoResponse)

		var receivedPeers []PeerInfo
		json.Unmarshal([]byte(askPeersInfoResponse.MsgContent), &receivedPeers)
		p.recordPeers(receivedPeers)
		p.log("updated peers info")
	}
}

func (p *Peer) recordPeers(receivedPeers []PeerInfo) {
	// add received peers info and avoid repeat
	for _, receivedPeerInfo := range receivedPeers {
		flag := 0
		for _, selfPeerInfo := range p.Peers {
			if receivedPeerInfo.Addr == selfPeerInfo.Addr && receivedPeerInfo.Port == selfPeerInfo.Port {
				flag = 1
			}
		}
		if flag == 0 { // new peer info
			p.Peers = append(p.Peers, receivedPeerInfo)
		}
	}
}

func main() {
	p1 := Peer{"127.0.0.1", 50001, []PeerInfo{}}
	p1.Connect("127.0.0.1", 99999)

	// time.Sleep(2 * time.Second)

	p2 := Peer{"127.0.0.1", 50002, []PeerInfo{}}
	p2.Connect("127.0.0.1", 50001)

	// time.Sleep(2 * time.Second)

	p3 := Peer{"127.0.0.1", 50003, []PeerInfo{}}
	p3.Connect("127.0.0.1", 50002)

	p4 := Peer{"127.0.0.1", 50004, []PeerInfo{}}
	p4.Connect("127.0.0.1", 50002)

	p5 := Peer{"127.0.0.1", 50005, []PeerInfo{}}
	p5.Connect("127.0.0.1", 50003)

	p6 := Peer{"127.0.0.1", 50006, []PeerInfo{}}
	p6.Connect("127.0.0.1", 50001)

	time.Sleep(1 * time.Second)
	p1.log("updated peers info")
	p2.log("updated peers info")
	p3.log("updated peers info")
	p4.log("updated peers info")
	p5.log("updated peers info")
	p6.log("updated peers info")

	select {}
}
