package consumer

import (
	"fmt"
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
		if err := nc.SubscribeSubject(infranats.SubjectGroupChat+"*", CastSendChat, chatOpts...); err != nil {
			fmt.Printf("Failed to subscribe subject: %v\n", err)
			return err
		}
		groupUpdOpts := []gonats.SubOpt{
			gonats.AckNone(),
			gonats.DeliverNew(),
		}

		if err = nc.SubscribeSubject(infranats.SubjectGroupUpdate+"*", CastUpdateGroup, groupUpdOpts...); err != nil {
			fmt.Printf("Failed to subscribe subject: %v\n", err)
			return err
		}

		leaveGroupOpt := []gonats.SubOpt{
			gonats.AckNone(),
			gonats.DeliverNew(),
		}
		if err = nc.SubscribeSubject(infranats.SubjectGroupLeave+"*", CastLeaveGroup, leaveGroupOpt...); err != nil {
			fmt.Printf("Failed to subscribe subject: %v\n", err)
			return err
		}

		return nil
	})
}
