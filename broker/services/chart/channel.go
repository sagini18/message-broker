package chart

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Channel(c echo.Context, channel *channelconsumer.Channel) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported")
	}

	channelEvents := channel.Get()
	data, err := json.Marshal(channelEvents)
	if err != nil {
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
		return err
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
				http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
				return err
			}
			fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
			flusher.Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}
}
