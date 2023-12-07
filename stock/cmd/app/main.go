package main

import (
	"context"
	stdHTTP "net/http"
	"os"
	"os/signal"
	"stock/internal/delivery"
	"stock/internal/delivery/event"
	"stock/internal/repository"
	"stock/internal/service"
	"stock/pkg/observability"
	"stock/pkg/redis"
	"strconv"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	traceProvider := observability.ConfigureTraceProvider()
	traceDB, err := otelsql.Open("postgres", os.Getenv("POSTGRES_URL"),
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

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher watermillMessage.Publisher
	redisPublisher = redis.NewRedisPublisher(redisClient, watermillLogger)

	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}
	redisPublisher = observability.TracingPublisherDecorator{Publisher: redisPublisher}

	//redisSubscriber := redis.NewRedisSubscriber(redisClient, watermillLogger)
	eventBus := event.NewBus(redisPublisher)

	err = repository.InitializeDatabaseSchema(db)
	if err != nil {
		panic(err)
	}

	repository := repository.NewProductRepository(db)
	service := service.New(repository)

	port, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		panic(err)
	}
	delivery := delivery.NewDelivery(port, redisClient, &service, eventBus)

	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return delivery.Event.Run(ctx)
	})

	errgrp.Go(func() error {
		<-delivery.Event.Running()

		err := delivery.Rest.Run()
		if err != nil && err != stdHTTP.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return traceProvider.Shutdown(context.Background())
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return delivery.Rest.Stop(context.Background())
	})

	if err := errgrp.Wait(); err != nil {
		panic(err)
	}
}
