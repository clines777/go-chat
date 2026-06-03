package nats

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	gonats "github.com/nats-io/nats.go"
	"gochat/internal/infra"
)

const StreamChat = "CHAT"
const SubjectGroupChat = "group.chat."
const SubjectGroupUpdate = "group.update."

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
	cfg := &gonats.StreamConfig{
		Name:     StreamChat,
		Subjects: []string{SubjectGroupChat + "*", SubjectGroupUpdate + "*"},
		Storage:  gonats.MemoryStorage,
		MaxAge:   5 * time.Minute,
	}
	_, err := c.JS.StreamInfo(StreamChat)
	if errors.Is(err, gonats.ErrStreamNotFound) {
		_, err = c.JS.AddStream(cfg)
		return err
	}
	if err != nil {
		return err
	}
	_, err = c.JS.UpdateStream(cfg)
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

func (c *Client) publish(subject string, data []byte) error {
	_, err := c.JS.Publish(subject, data)
	return err
}

func Publish(subject string, v any) error {
	nc, err := GetNats()
	if err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return nc.publish(subject, data)
}

// SubscribeSubject 訂閱指定subject，選項由caller傳入。
func (c *Client) SubscribeSubject(subject string, handler gonats.MsgHandler, opts ...gonats.SubOpt) error {
	_, err := c.JS.Subscribe(subject, handler, opts...)
	return err
}
