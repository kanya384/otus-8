package delivery

import (
	"context"

	"stock/internal/delivery/event"
	"stock/internal/delivery/rest"
	"stock/internal/service"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type Delivery struct {
	Rest  *rest.Rest
	Event *message.Router
}

func NewDelivery(
	port int,
	redisClient *redis.Client,
	service *service.Service,
	eventBus *cqrs.EventBus,
) Delivery {
	restApi := rest.New(port, service)

	eventHandler := event.NewHandler(service, eventBus)
	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))
	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	/*id, _ := uuid.Parse("de28f032-3942-48a1-ba50-b9eb26aa2ad3")
	eventBus.Publish(context.Background(), products.ReserveProducts_v1{
		Header:  types.NewEventHeader(),
		OrderId: uuid.New(),
		Products: []products.ReserveProductItem{
			{
				ProductId: id,
				Quantity:  10,
			},
		},
	})*/

	router := event.NewRouter(
		eventHandler,
		eventProcessorConfig,
		watermillLogger,
	)
	return Delivery{
		Rest:  restApi,
		Event: router,
	}
}
