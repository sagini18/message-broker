package chart

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Consumer(c echo.Context, consumerStore *channelconsumer.InMemoryConsumerCache) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported")
	}

	consumerEvents := consumerStore.GetEventCount()
	data, err := json.Marshal(consumerEvents)
	if err != nil {
		http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
		return err
	}
	fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
	flusher.Flush()

	sseChannel := consumerStore.SSEChannel()

	for {
		select {
		case <-sseChannel:
			consumerEvents := consumerStore.GetEventCount()
			data, err := json.Marshal(consumerEvents)
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
