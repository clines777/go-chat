package consumer

import (
	gonats "github.com/nats-io/nats.go"
	infranats "gochat/internal/infra/nats"
)

func Register() {
	infranats.RegisterConsumer(func() error {
		nc, err := infranats.GetNats()
		if err != nil {
			return err
		}
		chatOpts := []gonats.SubOpt{
			gonats.AckNone(),
			gonats.DeliverNew(),
		}
		if err := nc.SubscribeSubject(infranats.SubjectGroupChat, GroupChat, chatOpts...); err != nil {
			return err
		}
		groupUpdOpts := []gonats.SubOpt{
			gonats.AckNone(),
			gonats.DeliverNew(),
		}
		return nc.SubscribeSubject(infranats.SubjectGroupUpdate, GroupUpdate, groupUpdOpts...)
	})
}
