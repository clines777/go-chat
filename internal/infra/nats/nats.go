package nats

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	gonats "github.com/nats-io/nats.go"
	"gochat/internal/infra"
)

var invalidConsumerChar = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)

const StreamName = "CHAT"

var conn *Client

type Client struct {
	nc *gonats.Conn
	JS gonats.JetStreamContext
}

func GetNats() (*Client, error) {
	if conn != nil {
		return conn, nil
	}
	cfg := infra.GetEnvConfig()
	nc, err := gonats.Connect(cfg.NatsURL)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats jetstream: %w", err)
	}
	conn = &Client{nc: nc, JS: js}
	return conn, nil
}

func (c *Client) Ping() error {
	if !c.nc.IsConnected() {
		return fmt.Errorf("nats not connected")
	}
	return nil
}

func (c *Client) EnsureStream() error {
	_, err := c.JS.StreamInfo(StreamName)
	if errors.Is(err, gonats.ErrStreamNotFound) {
		_, err = c.JS.AddStream(&gonats.StreamConfig{
			Name:     StreamName,
			Subjects: []string{"chat.group.*"},
			Storage:  gonats.MemoryStorage,
			MaxAge:   5 * time.Minute,
		})
	}
	return err
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.JS.Publish(subject, data)
	return err
}

// SubscribeGroupChat starts a durable push consumer for this server instance.
// The handler receives (subject, rawJSON) for each message.
func (c *Client) SubscribeGroupChat(serverName string, handler func(subject string, data []byte)) error {
	_, err := c.JS.Subscribe(
		"chat.group.*",
		func(msg *gonats.Msg) { handler(msg.Subject, msg.Data) },
		gonats.Durable("chat-"+invalidConsumerChar.ReplaceAllString(serverName, "_")),
		gonats.AckNone(),
		gonats.DeliverNew(),
	)
	return err
}
