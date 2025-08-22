package db

import (
	"context"
	"log"
)

type Cache struct {
	Db               DB
	CacheMapOrder    map[string]Order
	CacheMapDelivery map[string]Delivery
	CacheMapPay      map[string]Payment
	CacheMapItems    map[string][]Item
}

// Close implements DB.
func (c *Cache) Close() error {
	return c.Db.Close()
}

// FindByUID implements DB.
func (c *Cache) FindByUID(ctx context.Context, uid string) (Order, Delivery, Payment, []Item, error) {
	order, orderPresent := c.CacheMapOrder[uid]
	delivery, deliveryPresent := c.CacheMapDelivery[uid]
	pay, payPresent := c.CacheMapPay[uid]
	items, itemsPresent := c.CacheMapItems[uid]

	if !orderPresent || !deliveryPresent || !payPresent || !itemsPresent {
		var err error
		order, delivery, pay, items, err = c.Db.FindByUID(ctx, uid)
		if err != nil {
			return Order{}, Delivery{}, Payment{}, []Item{}, err
		}

		c.CacheMapOrder[uid] = order
		c.CacheMapDelivery[uid] = delivery
		c.CacheMapPay[uid] = pay
		c.CacheMapItems[uid] = items

		log.Print("Cache miss")
	}

	return order, delivery, pay, items, nil
}

// SaveDelivery implements DB.
func (c *Cache) SaveDelivery(ctx context.Context, d Delivery) error {
	return c.Db.SaveDelivery(ctx, d)
}

// SaveItem implements DB.
func (c *Cache) SaveItem(ctx context.Context, item Item) error {
	return c.Db.SaveItem(ctx, item)
}

// SaveOrder implements DB.
func (c *Cache) SaveOrder(ctx context.Context, o Order) error {
	return c.Db.SaveOrder(ctx, o)
}

// SavePayment implements DB.
func (c *Cache) SavePayment(ctx context.Context, pay Payment) error {
	return c.Db.SavePayment(ctx, pay)
}

func CacheDb(db DB) (DB, error) {
	var cache = &Cache{Db: db}
	return cache.preLoadCache(), nil
}

func AddCaching(userDb *DB) error {
	cachedDb, err := CacheDb(*userDb)
	if err != nil {
		return err
	}
	*userDb = cachedDb
	return nil
}

func (c *Cache) preLoadCache() *Cache {
	c.CacheMapOrder = map[string]Order{}
	c.CacheMapDelivery = map[string]Delivery{}
	c.CacheMapPay = map[string]Payment{}
	c.CacheMapItems = map[string][]Item{}
	return c
}
