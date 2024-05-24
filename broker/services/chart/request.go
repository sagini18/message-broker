package chart

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Request(c echo.Context, requestCounter *channelconsumer.RequestCounter) error {
	requestEvents := requestCounter.GetEventCount()
	if len(requestEvents) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	return c.JSON(http.StatusOK, requestEvents)
}
