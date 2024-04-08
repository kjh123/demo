package event

import (
	"context"
	"event/event"
	"event/kafka"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"event",
	fx.Provide(
		NewSend,
		NewReceive,

		fx.Annotate(NewSend, fx.As(new(event.Sender))),
		fx.Annotate(NewReceive, fx.As(new(event.Receiver))),
	),

	fx.Invoke(func(sender *Send, lifecycle fx.Lifecycle) {
		lifecycle.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return sender.Close()
			}},
		)
	}),

	fx.Invoke(func(receive *Receive, lifecycle fx.Lifecycle) {
		lifecycle.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return receive.Close()
			}},
		)
	}),
)

type Params struct {
	fx.In

	KafkaAddr string `name:"kafka_addr"`
	Topic     string `name:"kafka_topic"`
}

type Send struct {
	Sender event.Sender
}

func (s *Send) Send(ctx context.Context, msg event.Event) error {
	return s.Sender.Send(ctx, msg)
}

func (s *Send) Close() error {
	return s.Sender.Close()
}

func NewSend(params Params) *Send {
	sender, err := kafka.NewKafkaSender([]string{params.KafkaAddr}, params.Topic)
	if err != nil {
		panic(err)
	}

	return &Send{sender}
}

type Receive struct {
	Receiver event.Receiver
}

func (r *Receive) Receive(ctx context.Context, handler event.Handler) error {
	return r.Receiver.Receive(ctx, handler)
}

func (r *Receive) Close() error {
	return r.Receiver.Close()
}

func NewReceive(params Params) *Receive {
	receiver, err := kafka.NewKafkaReceiver([]string{params.KafkaAddr}, params.Topic)
	if err != nil {
		panic(err)
	}

	return &Receive{receiver}
}
