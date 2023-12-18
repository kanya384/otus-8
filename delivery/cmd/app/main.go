package main

import (
	"context"
	"fmt"
	"net/http"
	"payments/internal/adapters/repository"
	"payments/internal/app/command"

	"os"
	"os/signal"
	"payments/internal/config/config"
	"payments/internal/ports/message"
	"payments/internal/ports/message/event"
	"payments/internal/ports/rest"
	"payments/pkg/observability"
	"payments/pkg/redis"

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

	repository := repository.NewPaymentPostgresRepository(db)

	redisClient := redis.NewRedisClient(cfg.RedisAddress())
	defer redisClient.Close()

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	redisPublisher := redis.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}
	//defer redisPublisher.Close()

	redisSubscriber := redis.NewRedisSubscriber(redisClient, watermillLogger)
	//defer redisSubscriber.Close()

	traceProvider := observability.ConfigureTraceProvider(cfg.Jaeger.Endpoint)

	eventBus := event.NewBus(redisPublisher)

	eventProcessorConfig := event.NewProcessorConfig(redisClient, watermillLogger)

	messageRouter := message.NewMessageRouter(redisPublisher, redisSubscriber, eventProcessorConfig, watermillLogger)

	errgrp, ctx := errgroup.WithContext(ctx)

	commandHandlers := command.NewHandler(eventBus, repository)
	httpRouter := rest.New(cfg.App.Port, commandHandlers)

	errgrp.Go(func() error {
		return messageRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		// we don't want to start HTTP server before Watermill router (so service won't be healthy before it's ready)
		<-messageRouter.Running()

		err := httpRouter.Run()
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return httpRouter.Stop(context.Background())
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return messageRouter.Close()
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
