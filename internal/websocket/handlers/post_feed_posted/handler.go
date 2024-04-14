package post_feed_posted

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/wagslane/go-rabbitmq"

	"service-post-notification/internal/schema/components/post"
)

type Handler struct {
	rabbit   *rabbitmq.Conn
	upgrader websocket.Upgrader
}

func New(rabbit *rabbitmq.Conn) *Handler {
	return &Handler{
		rabbit:   rabbit,
		upgrader: websocket.Upgrader{},
	}
}

func (h *Handler) Handle(c echo.Context) error {
	userID := c.Get("userID").(string)
	if userID == "" {
		return websocket.ErrBadHandshake
	}

	conn, err := h.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	go func() {
		consumer, err := rabbitmq.NewConsumer(
			h.rabbit,
			"",
			rabbitmq.WithConsumerOptionsQueueAutoDelete,
			rabbitmq.WithConsumerOptionsQueueExclusive,
			rabbitmq.WithConsumerOptionsRoutingKey(userID),
			rabbitmq.WithConsumerOptionsExchangeName("post-created"),
			rabbitmq.WithConsumerOptionsExchangeKind("topic"),
			rabbitmq.WithConsumerOptionsExchangeDeclare,
		)
		if err != nil {
			c.Logger().Error(err)
		}

		err = consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
			var message post.Post
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				return rabbitmq.NackDiscard
			}

			err = conn.WriteJSON(post.Post{
				PostID:           message.PostID,
				PostText:         message.PostText,
				PostAuthorUserID: message.PostAuthorUserID,
			})
			if err != nil {
				consumer.Close()
			}
			return rabbitmq.Ack
		})
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	return nil
}
