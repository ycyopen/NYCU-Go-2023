package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/gorilla/websocket"
	"github.com/reactivex/rxgo/v2"
)

type client chan<- string // an outgoing message channel

var (
	entering      = make(chan client)
	leaving       = make(chan client)
	messages      = make(chan rxgo.Item) // all incoming client messages
	ObservableMsg = rxgo.FromChannel(messages)
)

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	MessageBroadcast := ObservableMsg.Observe()
	for {
		select {
		case msg := <-MessageBroadcast:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg.V.(string)
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn *websocket.Conn, ch <-chan string) {
	for msg := range ch {
		conn.WriteMessage(1, []byte(msg))
	}
}

func wshandle(w http.ResponseWriter, r *http.Request) {
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "你是 " + who + "\n"
	messages <- rxgo.Of(who + " 來到了現場" + "\n")
	entering <- ch

	defer func() {
		log.Println("disconnect !!")
		leaving <- ch
		messages <- rxgo.Of(who + " 離開了" + "\n")
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		messages <- rxgo.Of(who + " 表示: " + string(msg))
	}
}

func InitObservable() {
	// TODO: Please create an Observable to handle the messages
	/*
		ObservableMsg = ObservableMsg.Filter(...) ... {
		}).Map(...) {
			...
		})
	*/
	swearWords, _ := os.ReadFile("swear_word.txt")
	sensitiveNames, _ := os.ReadFile("sensitive_name.txt")

	swearWordsSlice := strings.Split(string(swearWords), "\n")
	sensitiveNamesSlice := strings.Split(string(sensitiveNames), "\n")

	ObservableMsg = ObservableMsg.Filter(func(item interface{}) bool {
		msg, ok := item.(string)
		if !ok {
			return false
		}
		for _, swearWord := range swearWordsSlice {
			swearWord = strings.TrimSpace(swearWord)
			if swearWord != "" && strings.Contains(msg, swearWord) {
				return false
			}
		}
		return true
	}).Map(func(_ context.Context, item interface{}) (interface{}, error) {
		msg, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("item is not a string")
		}
		if !utf8.ValidString(msg) {
			return nil, fmt.Errorf("invalid UTF-8 string")
		}
		for _, sensitiveName := range sensitiveNamesSlice {
			sensitiveName = strings.TrimSpace(sensitiveName)
			if sensitiveName != "" && strings.Contains(msg, sensitiveName) {
				/*if len(sensitiveName) <= 2 {
					msg = strings.ReplaceAll(msg, sensitiveName, sensitiveName[:3]+"*")
				} else {*/
				msg = strings.ReplaceAll(msg, sensitiveName, sensitiveName[:3]+"*"+sensitiveName[6:])
				//}
			}
		}
		return msg, nil
	})
}

func main() {
	InitObservable()
	go broadcaster()
	http.HandleFunc("/wschatroom", wshandle)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("server start at :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
