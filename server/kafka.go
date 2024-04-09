package main

import (
	"context"
	"event/event"
	"log/slog"

	"go.uber.org/fx"
)

type ReceiveParams struct {
	fx.In

	Receiver event.Receiver
}

func Receiver(params ReceiveParams) {
	err := params.Receiver.Receive(context.Background(), func(_ context.Context, event event.Event) error {
		slog.Info("receive", "key", event.Key(), "value", string(event.Value()))
		return nil
	})

	if err != nil {
		slog.Error("kafka receiver", "receive error", err)
		return
	}
}
