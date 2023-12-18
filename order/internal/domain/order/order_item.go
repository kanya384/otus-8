package order

type OrderItem struct {
	ProductUUID string `db:"product_uuid"`
	Price       int    `db:"price"`
	Quantity    int    `db:"quantity"`
}
