package bot

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/wagslane/go-rabbitmq"

	"github.com/carfloresf/financial-chat/internal/constants"
	"github.com/carfloresf/financial-chat/internal/queue"
	"github.com/carfloresf/financial-chat/internal/stooq"
)

type Executor struct {
	stooq *stooq.Client
	queue queue.Queuer
}

func NewExecutor(stooq *stooq.Client, queue queue.Queuer) *Executor {
	return &Executor{
		stooq: stooq,
		queue: queue,
	}
}

func (c *Executor) Execute(delivery rabbitmq.Delivery) rabbitmq.Action {
	requestMessage := cast.ToString(delivery.Body)

	log.Printf("commandmanager received: %s", cast.ToString(delivery.Body))

	responseMessage := []byte("")

	switch {
	case strings.HasPrefix(requestMessage, constants.CommandStock):
		stockCode := strings.TrimPrefix(requestMessage, constants.CommandStock)
		responseMessage = c.getStockMessage(stockCode)

	case strings.HasPrefix(requestMessage, constants.CommandHelp):
		responseMessage = []byte("available commands: " + constants.CommandStock + "[stock_code], " + constants.CommandHelp)

	case strings.HasPrefix(requestMessage, "/"):
		responseMessage = []byte("unknown command received")
	}

	log.Printf("response requestMessage to command %s: %s", delivery.Body, string(responseMessage))

	err := c.queue.Publish(responseMessage, delivery.Delivery.CorrelationId, "response")
	if err != nil {
		log.Errorf("error publishing response: %s", err)

		return rabbitmq.Ack
	}

	return rabbitmq.Ack
}

func (c *Executor) getStockMessage(stockCode string) []byte {
	var responseMessage []byte

	stock, err := c.stooq.GetStockData(context.Background(), stockCode)
	if err != nil {
		log.Errorf("error getting stock data: %s", err)

		responseMessage = []byte("error getting stock data")
	} else {
		if stock.Close != "N/D" {
			responseMessage = []byte(fmt.Sprintf("%s quote is $%s per share",
				strings.ToUpper(stockCode), stock.Close))
		} else {
			responseMessage = []byte(fmt.Sprintf("%s not found", strings.ToUpper(stockCode)))
		}
	}

	return responseMessage
}
