package chart

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Messages(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache) error {
	messageEvents := messageQueue.GetEventCount()
	if len(messageEvents) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	return c.JSON(http.StatusOK, messageEvents)
}
