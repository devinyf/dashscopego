package httpclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Timeout for establishing the connection and for reading/writing messages
	writeWait = 30 * time.Second

	pongWait   = 20 * time.Second
	pingPeriod = (pongWait * 8) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 1024
	// maxMessageSize = 512
)

type WsMessage struct {
	// ws data type, e.g. websocket.TextMessage, websocket.BinaryMessage...
	Type int
	// ws data body
	Data []byte
}

// Client represents a websocket client.
type WsClient struct {
	Url        string
	Headers    http.Header
	Conn       *websocket.Conn
	inputChan  chan WsMessage
	outputChan chan WsMessage
	errChan    chan error
}

func NewWsClient(url string, headers http.Header) *WsClient {
	return &WsClient{
		Url:     url,
		Headers: headers,
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *WsClient) readPump() {
	defer func() {
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				c.errChan <- err
			}
			break
		}

		fmt.Println("message: ", string(message))
		c.outputChan <- WsMessage{
			Type: websocket.TextMessage,
			Data: message,
		}
		// Process the message (this part needs to be implemented based on your application logic).
	}
}

// writePump pumps messages from the write channel to the websocket connection.
func (c *WsClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.inputChan:
			if !ok {
				// The write channel is closed.
				c.errChan <- fmt.Errorf("write channel is closed")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if message.Type == websocket.TextMessage {
				fmt.Printf("send start data-- err: %v\n", string(message.Data))
			}

			if err := c.Conn.WriteMessage(message.Type, message.Data); err != nil {
				fmt.Println("err in write message: ", err)
				c.errChan <- err
				return
			}

			c.errChan <- nil
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.errChan <- err
				return
			}
		}
	}
}

// connect initializes the websocket connection and starts the read and write pumps.
func (c *WsClient) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.Url, c.Headers)
	if err != nil {
		return err
	}
	c.Conn = conn
	c.inputChan = make(chan WsMessage, 100)
	c.outputChan = make(chan WsMessage, 100)
	c.errChan = make(chan error, 1)
	go c.writePump()
	go c.readPump()
	return nil
}

// StartClient starts the client operation.
func (c *WsClient) ConnClient(req interface{}) error {
	if err := c.connect(); err != nil {
		// log.Fatal("dial:", err)
		return err
	}

	reqJson, _ := json.Marshal(req)
	reqInput := WsMessage{
		Type: websocket.TextMessage,
		Data: reqJson,
	}

	c.inputChan <- reqInput

	err, ok := <-c.errChan
	if ok && err != nil {
		fmt.Println("error: ", err)
	}
	return nil
}

func (c *WsClient) CloseClient() error {
	close(c.inputChan)
	close(c.outputChan)
	close(c.errChan)
	c.Conn.Close()
	return nil
}

func (c *WsClient) SendBinaryDates(data []byte) {
	streamInput := WsMessage{
		Type: websocket.BinaryMessage,
		Data: data,
	}

	c.inputChan <- streamInput
}

func (c *WsClient) ResultChans() (<-chan WsMessage, <-chan error) {
	return c.outputChan, c.errChan
}
