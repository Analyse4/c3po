package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"hedwig/pbprotocol"
	"log"
	"net/url"
	"os"
	"os/signal"
)

var addr = flag.String("addr", "localhost:2244", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	req := new(pbprotocol.LoginReq)
	req.Username = "jack"
	req.Password = "123456"
	req.Type = pbprotocol.LoginReq_Visitor
	r, _ := proto.Marshal(req)
	mp := new(pbprotocol.MsgPack)
	mp.Router = "digimon.login"
	mp.Data = r
	data, _ := proto.Marshal(mp)
	err = c.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case <-done:
			return
		}
	}
	//ticker := time.NewTicker(time.Second)
	//defer ticker.Stop()

	//for {
	//	select {
	//	case <-done:
	//		return
	//	case t := <-ticker.C:
	//		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	//		if err != nil {
	//			log.Println("write:", err)
	//			return
	//		}
	//	case <-interrupt:
	//		log.Println("interrupt")
	//
	//		// Cleanly close the connection by sending a close message and then
	//		// waiting (with timeout) for the server to close the connection.
	//		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	//		if err != nil {
	//			log.Println("write close:", err)
	//			return
	//		}
	//		select {
	//		case <-done:
	//		case <-time.After(time.Second):
	//		}
	//		return
	//	}
	//}

}
