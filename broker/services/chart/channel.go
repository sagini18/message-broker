package chart

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Channel(c echo.Context, channel *channelconsumer.Channel) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusNotImplemented, "Streaming unsupported")
	}

	channelEvents := channel.Get()

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(channelEvents)
	if err != nil {
		return c.JSON(http.StatusNotImplemented, map[string]string{
			"type":    "StreamError",
			"message": "Streaming unsupported",
		})
	}
	fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
	flusher.Flush()

	sseChannel := channel.SSEChannel()

	for {
		select {
		case <-sseChannel:
			channelEvents := channel.Get()
			data, err := json.Marshal(channelEvents)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"type":    "MarshalError",
					"message": "Error in marshalling",
					"cause":   err.Error(),
				})
			}
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
			flusher.Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}
