package main

import (
	"context"
	"event/event"
	"go.uber.org/fx"
	"log/slog"
)

type ReceiveParams struct {
	fx.In

	Receiver event.Receiver
}

func Receiver(params ReceiveParams) {
	err := params.Receiver.Receive(context.Background(), func(ctx context.Context, event event.Event) error {
		slog.Info("receive", "key", event.Key(), "value", string(event.Value()))
		return nil
	})

	if err != nil {
		slog.Error("kafka receiver", "receive error", err)
		return
	}
}
