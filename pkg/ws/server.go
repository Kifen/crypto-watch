package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/websocket"
)

func Connect(url string, dialer *websocket.Dialer) (*websocket.Conn, error) {
	conn, _, err := dialer.Dial(url, http.Header{})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Streamticker(conn *websocket.Conn) {
	for {
		msgByte, err := WsRead(conn)
		if err != nil {
			log.Fatal(err)
		}

		//log.Println("Logging data:\n", string(msgByte))
		j, err := simplejson.NewJson(msgByte)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("#PRICE SKYBTC: ", j.Get("c"))
		//os.Exit(1)
	}
}

func WsWrite(conn *websocket.Conn, data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, msg)
}

func WsRead(conn *websocket.Conn) ([]byte, error) {
	_, p, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return p, nil
}
