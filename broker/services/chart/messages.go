package chart

import (
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Messages(c echo.Context, messageQueue *channelconsumer.InMemoryMessageCache) error {
	allMessages := messageQueue.GetAll()

	if len(allMessages) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	var response []time.Time
	for _, messages := range allMessages {
		for _, message := range messages {
			response = append(response, message.ReceivedAt)
		}
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].Before(response[j])
	})
	return c.JSON(http.StatusOK, response)
}
