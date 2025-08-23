package db

import (
	"context"
	"go_server_l0/internal/config"
	"log"
)

type Cache struct {
	Db               CachableDB
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

// GetUIDs implements Cachable.
func (c *Cache) GetUIDs(ctx context.Context, limit int32) ([]string, error) {
	return c.Db.GetUIDs(ctx, limit)
}

func CacheDb(db CachableDB, ctx context.Context) (CachableDB, error) {
	var cache = &Cache{Db: db}
	return cache.preLoadCache(ctx), nil
}

func AddCaching(userDb *CachableDB, ctx context.Context) error {
	cachedDb, err := CacheDb(*userDb, ctx)
	if err != nil {
		return err
	}
	*userDb = cachedDb
	return nil
}

func (c *Cache) preLoadCache(ctx context.Context) *Cache {
	c.CacheMapOrder = map[string]Order{}
	c.CacheMapDelivery = map[string]Delivery{}
	c.CacheMapPay = map[string]Payment{}
	c.CacheMapItems = map[string][]Item{}

	appConfig := config.LoadConfig()

	uids, err := c.Db.GetUIDs(ctx, appConfig.PreloadLimit)

	if err != nil {
		log.Println("Cannot preload cache")
		return c
	}

	log.Printf("Preloaded %v orders", len(uids))

	for _, uid := range uids {
		var order, delivery, pay, items, err = c.Db.FindByUID(ctx, uid)
		if err != nil {
			log.Printf("Error preloading %v", uid)
		}
		c.CacheMapOrder[uid] = order
		c.CacheMapDelivery[uid] = delivery
		c.CacheMapPay[uid] = pay
		c.CacheMapItems[uid] = items
	}

	return c
}
