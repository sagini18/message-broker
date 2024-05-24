package chart

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Channel(c echo.Context, channel *channelconsumer.Channel) error {
	channelEvents := channel.Get()
	if len(channelEvents) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	return c.JSON(http.StatusOK, channelEvents)
}
