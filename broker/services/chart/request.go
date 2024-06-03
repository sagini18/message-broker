package chart

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Request(c echo.Context, requestCounter *channelconsumer.RequestCounter) error {
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return c.JSON(http.StatusNotImplemented, map[string]string{
			"type":    "StreamError",
			"message": "Streaming unsupported",
		})
	}

	requestEvents := requestCounter.GetEventCount()
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	data, err := json.Marshal(requestEvents)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"type":    "MarshalError",
			"message": "Error in marshalling",
			"cause":   err.Error(),
		})
	}
	fmt.Fprintf(c.Response().Writer, "data: %s\n\n", data)
	flusher.Flush()

	sseChannel := requestCounter.SSEChannel()

	for {
		select {
		case <-sseChannel:
			requestEvents := requestCounter.GetEventCount()
			data, err := json.Marshal(requestEvents)
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
