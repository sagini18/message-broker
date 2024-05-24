package chart

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Consumer(c echo.Context, consumerStore *channelconsumer.InMemoryConsumerCache) error {
	consumerEvents := consumerStore.GetEventCount()
	if len(consumerEvents) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	return c.JSON(http.StatusOK, consumerEvents)
}
