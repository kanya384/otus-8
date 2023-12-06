package command

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

var marshaler = cqrs.JSONMarshaler{
	GenerateName: cqrs.StructName,
}

func NewBus(publisher message.Publisher, config cqrs.CommandBusConfig) *cqrs.CommandBus {
	commandBus, err := cqrs.NewCommandBusWithConfig(publisher, config)
	if err != nil {
		panic(err)
	}

	return commandBus
}
