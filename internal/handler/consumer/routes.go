package consumer

import (
	"log"

	gonats "github.com/nats-io/nats.go"
	infranats "gochat/internal/infra/nats"
)

func Register() {
	infranats.RegisterConsumer(func() error {
		nc, err := infranats.GetNats()
		if err != nil {
			return err
		}

		// 目前都是純通知, 不需 ack, 也不管舊的
		opts := []gonats.SubOpt{
			gonats.AckNone(),
			gonats.DeliverNew(),
		}

		subs := []struct {
			subject string
			handler gonats.MsgHandler
		}{
			{infranats.SubjectGroupChat, CastSendChat},
			{infranats.SubjectGroupUpdate, CastUpdateGroup},
			{infranats.SubjectGroupLeave, CastLeaveGroup},
			{infranats.SubjectGroupPin, CastPinChat},
			{infranats.SubjectGroupUnpin, CastUnpinChat},
			{infranats.SubjectGroupDel, CastDelChat},
			{infranats.SubjectGroupBan, CastBanUser},
			{infranats.SubjectGroupUnban, CastUnbanUser},
			{infranats.SubjectGroupKick, CastKickUser},
		}

		for _, s := range subs {
			if err := nc.SubscribeSubject(s.subject, s.handler, opts...); err != nil {
				log.Printf("[consumer] subscribe %s failed: %v", s.subject, err)
				return err
			}
		}

		return nil
	})
}
