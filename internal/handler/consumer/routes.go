package consumer

import (
	"gochat/internal/infra"
	infranats "gochat/internal/infra/nats"
)

func Register() {
	infranats.RegisterConsumer(func() error {
		nc, err := infranats.GetNats()
		if err != nil {
			return err
		}
		return nc.SubscribeChatStream(infra.GetEnvConfig().ServerName, GroupChat)
	})
}
