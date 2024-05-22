package chart

import (
	"net/http"
	"sort"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sagini18/message-broker/broker/internal/channelconsumer"
)

func Consumer(c echo.Context, consumerStore *channelconsumer.InMemoryConsumerCache) error {
	allConsumers := consumerStore.GetAll()
	if len(allConsumers) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	var response []time.Time
	for _, consumer := range allConsumers {
		for _, consumer := range consumer {
			response = append(response, consumer.JoinedAt)
		}
	}

	sort.Slice(response, func(i, j int) bool {
		return response[i].Before(response[j])
	})
	return c.JSON(http.StatusOK, response)
}
