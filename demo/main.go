package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/lithammer/shortuuid/v4"
	"log"
	"sync"
	"time"
)

const (
	maxMessageSize = 64 * 1024 // 64 KB
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	waitTimeout    = 30 * time.Second
)
const fnindex = 18

func main() {
	//url := "academic.jinresearch.site"
	//url := "xueshu.52apikey.cn"
	url := "nsgzsupr.bja.sealos.run"
	content := "你好"
	client, err := NewClient(fmt.Sprintf("wss://%s/queue/join", url))
	if err != nil {
		panic(err)
	}

	messageChan, err := client.SendMessage(context.TODO(), content, "gpt-4o")
	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				log.Println("Message channel closed")
				return
			}
			log.Printf("Message: %s", msg)
		}
	}
}

type Client struct {
	conn        *websocket.Conn
	addr        string
	sessionHash string
	sendChan    chan []byte
	receiveChan chan []byte
	mu          sync.Mutex
	isConnected bool
}

func NewClient(addr string) (*Client, error) {
	c := &Client{
		addr:        addr,
		sessionHash: shortuuid.New(),
		sendChan:    make(chan []byte, 256),
		receiveChan: make(chan []byte, 256),
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	go c.readPump()
	go c.writePump()

	return c, nil
}

func (c *Client) connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}
	conn, _, err := dialer.Dial(c.addr, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	c.isConnected = true
	return nil
}

func (c *Client) SendMessage(ctx context.Context, message string, model string) (<-chan string, error) {
	responseChan := make(chan string)

	go c.handleMessageExchange(ctx, message, model, responseChan)

	return responseChan, nil
}

func (c *Client) handleMessageExchange(ctx context.Context, message string, model string, responseChan chan<- string) {
	defer close(responseChan)

	if err := c.performHandshake(); err != nil {
		responseChan <- fmt.Sprintf("Error: %v", err)
		return
	}

	if err := c.sendAiRequest(message, model); err != nil {
		responseChan <- fmt.Sprintf("Error: %v", err)
		return
	}

	c.processResponses(ctx, responseChan)
}

func (c *Client) performHandshake() error {
	if err := c.waitForMessage("send_hash"); err != nil {
		return fmt.Errorf("waiting for send_hash: %w", err)
	}

	if err := c.sendJSON(map[string]interface{}{
		"fn_index":     fnindex,
		"session_hash": c.sessionHash,
	}); err != nil {
		return fmt.Errorf("sending session hash: %w", err)
	}

	if err := c.waitForMessage("estimation"); err != nil {
		return fmt.Errorf("waiting for estimation: %w", err)
	}

	if err := c.waitForMessage("send_data"); err != nil {
		return fmt.Errorf("waiting for send_data: %w", err)
	}

	return nil
}

func (c *Client) sendAiRequest(message string, model string) error {
	req := AiRequest{
		Data: []interface{}{
			nil, 4096, model, message, "", 1, 1, [][]string{},
			nil, "Serve me as a writing and programming assistant.", "", nil,
		},
		FnIndex:     fnindex,
		SessionHash: c.sessionHash,
	}

	return c.sendJSON(req)
}

func (c *Client) processResponses(ctx context.Context, responseChan chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-c.receiveChan:
			var response AiResponse
			if err := json.Unmarshal(msg, &response); err != nil {
				responseChan <- fmt.Sprintf("Error unmarshalling response: %v", err)
				return
			}

			switch response.Msg {
			case "process_starts":
				// Process has started, wait for generating messages
			case "process_generating", "process_completed":
				if !response.Success {
					responseChan <- fmt.Sprintf("Error from server: %v", response.Output)
					return
				}
				if latestResponse := c.extractLatestResponse(response); latestResponse != "" {
					responseChan <- latestResponse
				}
				if response.Msg == "process_completed" {
					return
				}
			default:
				log.Printf("Unexpected message type: %s", response.Msg)
			}
		}
	}
}

func (c *Client) extractLatestResponse(response AiResponse) string {
	if len(response.Output.Data) <= 1 {
		return ""
	}

	conversations, ok := response.Output.Data[1].([]interface{})
	if !ok || len(conversations) == 0 {
		return ""
	}

	var latestResponse string
	if len(conversations) == 1 {
		conversation := conversations[0].([]interface{})
		if len(conversation) > 1 {
			latestResponse = conversation[1].(string)
		}
	} else {
		latestConversation := conversations[len(conversations)-1].([]interface{})
		if len(latestConversation) > 1 {
			latestResponse = latestConversation[1].(string)
		}
	}

	return latestResponse
}

func (c *Client) waitForMessage(expectedMsg string) error {
	select {
	case msg := <-c.receiveChan:
		var response AiResponse
		if err := json.Unmarshal(msg, &response); err != nil {
			return fmt.Errorf("unmarshalling response: %w", err)
		}
		if response.Msg != expectedMsg {
			return fmt.Errorf("unexpected message: expected %s, got %s", expectedMsg, response.Msg)
		}
		return nil
	case <-time.After(waitTimeout):
		return errors.New("timeout waiting for server message")
	}
}

func (c *Client) sendJSON(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	marshal, err := json.Marshal(v)
	if err != nil {
		return err
	}
	log.Printf("Sending: %s", marshal)
	return c.conn.WriteMessage(websocket.TextMessage, marshal)
	//return c.conn.WriteJSON(v)
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
		c.isConnected = false
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.receiveChan <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type AiRequest struct {
	Data        []interface{} `json:"data"`
	EventData   interface{}   `json:"event_data"`
	FnIndex     int           `json:"fn_index"`
	SessionHash string        `json:"session_hash"`
}

type AiResponse struct {
	Msg    string `json:"msg"`
	Output struct {
		Data         []interface{} `json:"data"`
		IsGenerating bool          `json:"is_generating"`
	} `json:"output"`
	Success bool `json:"success"`
}
