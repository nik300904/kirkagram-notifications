package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"kirkagram-notification/internal/transport/ws"
	"log/slog"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	log      *slog.Logger
	group    string
}

func NewConsumer(groupName string, log *slog.Logger) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group, err := sarama.NewConsumerGroup([]string{"localhost:29092"}, groupName, config)
	if err != nil {
		panic(err)
	}

	return &Consumer{
		consumer: group,
		log:      log,
		group:    groupName,
	}
}

func (c *Consumer) ConsumeMessages(ctx context.Context, wsManager *ws.WebSocketManager, topic string) error {
	handler := consumerGroupHandler{
		log:       c.log,
		wsManager: wsManager,
	}

	for {
		err := c.consumer.Consume(ctx, []string{topic}, handler)
		if err != nil {
			return fmt.Errorf("error from consumer: %w", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type consumerGroupHandler struct {
	log       *slog.Logger
	wsManager *ws.WebSocketManager
}

func (h consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.log.Info("Starting to consume messages")
	for msg := range claim.Messages() {
		h.log.Info("Received message",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"value", string(msg.Value))
		fmt.Println("msg", string(msg.Value))

		switch msg.Topic {
		case "like":
			if err := h.wsManager.SendMessageLike(msg.Value); err != nil {
				return err
			}
		case "post":
			if err := h.wsManager.SendMessagePost(msg.Value); err != nil {
				return err
			}
		case "follow":
			if err := h.wsManager.SendMessageFollow(msg.Value); err != nil {
				return err
			}
		case "unfollow":
			if err := h.wsManager.SendMessageUnFollow(msg.Value); err != nil {
				return err
			}
		}

		sess.MarkMessage(msg, "")
	}

	h.log.Info("Finished consuming messages")
	return nil
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
