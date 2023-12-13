package repository_test

import (
	"context"
	"order/adapters/repository"
	"order/domain/order"
	"order/pkg/dateTime"
	"sync"
	"testing"

	_ "github.com/lib/pq"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB
var getDbOnce sync.Once

func getDb() *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", "dbname=db user=user host=localhost password=password sslmode=disable")
		if err != nil {
			panic(err)
		}
		err = repository.InitializeDatabaseSchema(db)
		if err != nil {
			panic(err)
		}
	})
	return db
}

func Test_AddOrder(t *testing.T) {
	start, err := dateTime.NewDateTime(10, 00)
	if err != nil {
		panic(err)
	}
	end, err := dateTime.NewDateTime(23, 00)
	if err != nil {
		panic(err)
	}

	db := getDb()
	repo := repository.NewOrderPostgresRepository(db)

	for i := 0; i < 2; i++ {
		createOrder, err := order.NewOrder(uuid.NewString(), "customer", "address", []order.OrderItem{{ProductUUID: uuid.NewString(), Quantity: 1}}, uuid.NewString(), start, end)
		if err != nil {
			panic(err)
		}
		err = repo.AddOrder(context.Background(), createOrder)
		require.NoError(t, err)

		repoOrder, err := repo.ReadOrder(context.Background(), createOrder.Uuid())
		require.NoError(t, err)

		require.Equal(t, createOrder, repoOrder)
	}

}

func Test_ReadOrder(t *testing.T) {
	start, err := dateTime.NewDateTime(10, 00)
	if err != nil {
		panic(err)
	}
	end, err := dateTime.NewDateTime(23, 00)
	if err != nil {
		panic(err)
	}

	db := getDb()
	repo := repository.NewOrderPostgresRepository(db)

	searchOrder, err := order.NewOrder(uuid.NewString(), "customer", "address", []order.OrderItem{{ProductUUID: uuid.NewString(), Quantity: 1}}, uuid.NewString(), start, end)
	if err != nil {
		panic(err)
	}
	err = repo.AddOrder(context.Background(), searchOrder)
	require.NoError(t, err)

	order, err := repo.ReadOrder(context.Background(), searchOrder.Uuid())
	require.NoError(t, err)

	require.Equal(t, searchOrder, order)

}

func Test_UpdateOrder(t *testing.T) {
	start, err := dateTime.NewDateTime(10, 00)
	if err != nil {
		panic(err)
	}
	end, err := dateTime.NewDateTime(23, 00)
	if err != nil {
		panic(err)
	}

	db := getDb()
	repo := repository.NewOrderPostgresRepository(db)

	oldOrder, err := order.NewOrder(uuid.NewString(), "customer", "address", []order.OrderItem{{ProductUUID: uuid.NewString(), Quantity: 1}}, uuid.NewString(), start, end)
	if err != nil {
		panic(err)
	}

	err = repo.AddOrder(context.Background(), oldOrder)
	if err != nil {
		panic(err)
	}

	expectedOrder, _ := order.NewOrder(oldOrder.Uuid(), "new-name", "new-adress", []order.OrderItem{{ProductUUID: uuid.NewString(), Quantity: 2}}, uuid.NewString(), start, end)

	upFn := func(ctx context.Context, oldOrder *order.Order) (*order.Order, error) {
		return order.NewOrder(oldOrder.Uuid(), expectedOrder.CustomerName(), expectedOrder.DeliveryAddress(), []order.OrderItem{{ProductUUID: expectedOrder.OrderItems()[0].ProductUUID, Quantity: 2}}, expectedOrder.PaymentUUID(), start, end)
	}

	_, err = repo.UpdateOrder(context.Background(), oldOrder.Uuid(), upFn)
	require.NoError(t, err)

	assertPersistedTrainingEquals(t, repo, expectedOrder)

}

func assertPersistedTrainingEquals(t *testing.T, repo repository.OrderPostgresRepository, order *order.Order) {
	persistedOrder, err := repo.ReadOrder(
		context.Background(),
		order.Uuid(),
	)
	require.NoError(t, err)

	assertTrainingsEquals(t, order, persistedOrder)
}

func assertTrainingsEquals(t *testing.T, or1, or2 *order.Order) {
	cmpOpts := []cmp.Option{
		cmp.AllowUnexported(
			order.Order{},
			dateTime.DateTime{},
		),
	}
	assert.True(
		t,
		cmp.Equal(or1, or2, cmpOpts...),
		cmp.Diff(or1, or2, cmpOpts...),
	)
}
