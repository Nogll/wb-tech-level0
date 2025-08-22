package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGSQL struct {
	dbpool *pgxpool.Pool
	url    string
}

// findById implements DB.
func (p PGSQL) FindByUID(ctx context.Context, uid string) (Order, Delivery, Payment, []Item, error) {
	var order Order
	var delivery Delivery
	var payment Payment
	var items []Item

	// Order
	orderRow, err := p.dbpool.Query(ctx, "SELECT * FROM orders WHERE order_uid = $1", uid)
	if err != nil {
		return order, delivery, payment, nil, err
	}
	order, err = pgx.CollectOneRow(orderRow, pgx.RowToStructByName[Order])
	if err != nil {
		return order, delivery, payment, nil, err
	}

	// Delivery
	deliveryRow, err := p.dbpool.Query(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", uid)
	if err != nil {
		return order, delivery, payment, nil, err
	}
	delivery, err = pgx.CollectOneRow(deliveryRow, pgx.RowToStructByName[Delivery])
	if err != nil {
		return order, delivery, payment, nil, err
	}

	// Payment
	paymentRow, err := p.dbpool.Query(ctx, "SELECT * FROM payments WHERE order_uid = $1", uid)
	if err != nil {
		return order, delivery, payment, nil, err
	}
	payment, err = pgx.CollectOneRow(paymentRow, pgx.RowToStructByName[Payment])
	if err != nil {
		return order, delivery, payment, nil, err
	}

	// Items
	itemRows, err := p.dbpool.Query(ctx, "SELECT * FROM items WHERE order_uid = $1", uid)
	if err != nil {
		return order, delivery, payment, nil, err
	}
	items, err = pgx.CollectRows(itemRows, pgx.RowToStructByName[Item])
	if err != nil {
		return order, delivery, payment, nil, err
	}

	return order, delivery, payment, items, nil
}

func (p *PGSQL) SaveOrder(ctx context.Context, o Order) error {
	_, err := p.dbpool.Exec(
		ctx,
		`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT DO NOTHING`,
		o.UID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SmID, o.DateCreated, o.OofShard,
	)
	return err
}

func (p *PGSQL) SaveDelivery(ctx context.Context, d Delivery) error {
	_, err := p.dbpool.Exec(
		ctx,
		`INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT DO NOTHING`,
		d.OrderUID, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email,
	)
	return err
}

func (p *PGSQL) SavePayment(ctx context.Context, pay Payment) error {
	_, err := p.dbpool.Exec(
		ctx,
		`INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		 ON CONFLICT DO NOTHING`,
		pay.OrderUID, pay.Transaction, pay.RequestID, pay.Currency, pay.Provider,
		pay.Amount, pay.PaymentDT, pay.Bank, pay.DeliveryCost, pay.GoodsTotal, pay.CustomFee,
	)
	return err
}

func (p *PGSQL) SaveItem(ctx context.Context, item Item) error {
	_, err := p.dbpool.Exec(
		ctx,
		`INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		 ON CONFLICT DO NOTHING`,
		item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
		item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status, item.OrderUID,
	)
	return err
}

func (p *PGSQL) Close() error {
	p.dbpool.Close()
	return nil
}

func ConnectToPGSQL(dbUrl string) (DB, error) {
	dbpool, err := pgxpool.New(context.Background(), dbUrl)

	if err != nil {
		return nil, err
	}

	return &PGSQL{dbpool, dbUrl}, nil
}
