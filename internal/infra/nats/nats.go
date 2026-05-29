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

const StreamChat = "CHAT"
const SubjectGroupChat = "group.chat.*"

var conn *Client
var consumers []func() error

func RegisterConsumer(fn func() error) {
	consumers = append(consumers, fn)
}

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

func (c *Client) EnsureStreams() error {
	_, err := c.JS.StreamInfo(StreamChat)
	if errors.Is(err, gonats.ErrStreamNotFound) {
		_, err = c.JS.AddStream(&gonats.StreamConfig{
			Name:     StreamChat,
			Subjects: []string{SubjectGroupChat},
			Storage:  gonats.MemoryStorage,
			MaxAge:   5 * time.Minute,
		})
	}
	return err
}

func (c *Client) Init() error {
	if err := c.EnsureStreams(); err != nil {
		return err
	}
	for _, fn := range consumers {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Publish(subject string, data []byte) error {
	_, err := c.JS.Publish(subject, data)
	return err
}

// SubscribeChatStream 訂閱group chat.
// The handler receives (subjectName, bytes) for each message.
func (c *Client) SubscribeChatStream(serverName string, handler func(subject string, data []byte)) error {
	_, err := c.JS.Subscribe(
		SubjectGroupChat,
		func(msg *gonats.Msg) { handler(msg.Subject, msg.Data) },
		gonats.Durable("chat-"+invalidConsumerChar.ReplaceAllString(serverName, "_")),
		gonats.AckNone(),
		gonats.DeliverNew(),
	)
	return err
}
