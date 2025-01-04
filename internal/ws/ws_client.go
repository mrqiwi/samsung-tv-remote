package ws

import (
	"crypto/tls"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	Conn *websocket.Conn
}

func NewWebSocketClient(url string) (*WebSocketClient, error) {
	dialer := websocket.Dialer{
			TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &WebSocketClient{
		Conn: conn,
	}, nil
}

func (client *WebSocketClient) SendMessage(msg []byte) error {
	return client.Conn.WriteMessage(websocket.TextMessage, msg)
}

func (client *WebSocketClient) ReadMessage() ([]byte, error) {
	_, message, err := client.Conn.ReadMessage()
	return message, err
}

func (client *WebSocketClient) Close() {
	client.Conn.Close()
}

