package main

import (
	"context"
	"fmt"
	"order/internal/adapters/repository"
	"order/internal/app/command"
	"order/internal/app/command/saga"

	"order/internal/config/config"
	"order/internal/ports/message"
	"order/internal/ports/message/event"
	"order/internal/ports/rest"
	"order/pkg/observability"
	"order/pkg/redis"
	"os"
	"os/signal"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.InitConfig("")
	if err != nil {
		panic(err)
	}

	traceDB, err := otelsql.Open("postgres", cfg.PostgresDsn(),
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL),
		otelsql.WithDBName("db"))
	if err != nil {
		panic(err)
	}

	db := sqlx.NewDb(traceDB, "postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := repository.InitializeDatabaseSchema(db); err != nil {
		panic(fmt.Errorf("failed to initialize database schema: %w", err))
	}

	repository := repository.NewOrderPostgresRepository(db)

	redisClient := redis.NewRedisClient(cfg.RedisAddress())
	defer redisClient.Close()

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	redisPublisher := redis.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}
	defer redisPublisher.Close()

	redisSubscriber := redis.NewRedisSubscriber(redisClient, watermillLogger)
	defer redisSubscriber.Close()

	traceProvider := observability.ConfigureTraceProvider(cfg.Jaeger.Endpoint)

	eventBus := event.NewBus(redisPublisher)
	saga := saga.NewOrderSaga(eventBus, repository)

	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	messageRouter := message.NewMessageRouter(redisPublisher, redisSubscriber, saga, eventProcessorConfig, watermillLogger)

	errgrp, ctx := errgroup.WithContext(ctx)

	commandHandlers := command.NewHandler(eventBus, repository)
	httpRouter := rest.New(cfg.App.Port, commandHandlers)

	errgrp.Go(func() error {
		return messageRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-messageRouter.Running()
		fmt.Printf("message router started\n")
		//TODO: start http router start here
		return httpRouter.Run()
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return httpRouter.Stop(context.Background())
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return traceProvider.Shutdown(context.Background())
	})

	err = errgrp.Wait()
	if err != nil {
		panic(err)
	}
}
