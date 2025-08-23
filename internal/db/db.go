package db

import (
	"context"
	"time"
)

// Order — основной заказ
type Order struct {
	UID               string    `db:"order_uid"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	ShardKey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`
}

// Delivery — адрес доставки, отдельная таблица
type Delivery struct {
	ID       int64  `db:"id"`        // PK
	OrderUID string `db:"order_uid"` // FK -> Order.UID
	Name     string `db:"name"`
	Phone    string `db:"phone"`
	Zip      string `db:"zip"`
	City     string `db:"city"`
	Address  string `db:"address"`
	Region   string `db:"region"`
	Email    string `db:"email"`
}

// Payment — информация о платеже
type Payment struct {
	ID           int64  `db:"id"`        // PK
	OrderUID     string `db:"order_uid"` // FK -> Order.UID
	Transaction  string `db:"transaction"`
	RequestID    string `db:"request_id"`
	Currency     string `db:"currency"`
	Provider     string `db:"provider"`
	Amount       int    `db:"amount"`
	PaymentDT    int64  `db:"payment_dt"` // timestamp unix
	Bank         string `db:"bank"`
	DeliveryCost int    `db:"delivery_cost"`
	GoodsTotal   int    `db:"goods_total"`
	CustomFee    int    `db:"custom_fee"`
}

// Item — товары в заказе
type Item struct {
	ID          int64  `db:"id"`        // PK
	OrderUID    string `db:"order_uid"` // FK -> Order.UID
	ChrtID      int64  `db:"chrt_id"`
	TrackNumber string `db:"track_number"`
	Price       int    `db:"price"`
	Rid         string `db:"rid"`
	Name        string `db:"name"`
	Sale        int    `db:"sale"`
	Size        string `db:"size"`
	TotalPrice  int    `db:"total_price"`
	NmID        int64  `db:"nm_id"`
	Brand       string `db:"brand"`
	Status      int    `db:"status"`
}

type DB interface {
	FindByUID(ctx context.Context, uid string) (Order, Delivery, Payment, []Item, error)
	SaveOrder(ctx context.Context, o Order) error
	SaveDelivery(ctx context.Context, d Delivery) error
	SavePayment(ctx context.Context, pay Payment) error
	SaveItem(ctx context.Context, item Item) error

	Close() error
}

type Cachable interface {
	GetUIDs(ctx context.Context, limit int32) ([]string, error)
}

type CachableDB interface {
	Cachable
	DB
}
