package mapper

import (
	"go_server_l0/api"
	"go_server_l0/internal/db"
)

// MapAPIOrderToDB разносит API-модель в базу
func MapAPIOrderToDB(o api.Order) (db.Order, db.Delivery, db.Payment, []db.Item) {
	orderDB := db.Order{
		UID:               o.OrderUID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		ShardKey:          o.ShardKey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated,
		OofShard:          o.OofShard,
	}

	deliveryDB := db.Delivery{
		OrderUID: o.OrderUID,
		Name:     o.Delivery.Name,
		Phone:    o.Delivery.Phone,
		Zip:      o.Delivery.Zip,
		City:     o.Delivery.City,
		Address:  o.Delivery.Address,
		Region:   o.Delivery.Region,
		Email:    o.Delivery.Email,
	}

	paymentDB := db.Payment{
		OrderUID:     o.OrderUID,
		Transaction:  o.Payment.Transaction,
		RequestID:    o.Payment.RequestID,
		Currency:     o.Payment.Currency,
		Provider:     o.Payment.Provider,
		Amount:       o.Payment.Amount,
		PaymentDT:    o.Payment.PaymentDT,
		Bank:         o.Payment.Bank,
		DeliveryCost: o.Payment.DeliveryCost,
		GoodsTotal:   o.Payment.GoodsTotal,
		CustomFee:    o.Payment.CustomFee,
	}

	itemsDB := make([]db.Item, len(o.Items))
	for i, item := range o.Items {
		itemsDB[i] = db.Item{
			OrderUID:    o.OrderUID,
			ChrtID:      item.ChrtID,
			TrackNumber: item.TrackNumber,
			Price:       item.Price,
			Rid:         item.RID,
			Name:        item.Name,
			Sale:        item.Sale,
			Size:        item.Size,
			TotalPrice:  item.TotalPrice,
			NmID:        item.NmID,
			Brand:       item.Brand,
			Status:      item.Status,
		}
	}

	return orderDB, deliveryDB, paymentDB, itemsDB
}

func MapOrderDBToAPI(o db.Order, d db.Delivery, p db.Payment, itemsDB []db.Item) api.Order {
	items := make([]api.Item, len(itemsDB))
	for i, it := range itemsDB {
		items[i] = api.Item{
			ChrtID:      it.ChrtID,
			TrackNumber: it.TrackNumber,
			Price:       it.Price,
			RID:         it.Rid,
			Name:        it.Name,
			Sale:        it.Sale,
			Size:        it.Size,
			TotalPrice:  it.TotalPrice,
			NmID:        it.NmID,
			Brand:       it.Brand,
			Status:      it.Status,
		}
	}

	return api.Order{
		OrderUID:    o.UID,
		TrackNumber: o.TrackNumber,
		Entry:       o.Entry,
		Delivery: api.Delivery{
			Name:    d.Name,
			Phone:   d.Phone,
			Zip:     d.Zip,
			City:    d.City,
			Address: d.Address,
			Region:  d.Region,
			Email:   d.Email,
		},
		Payment: api.Payment{
			Transaction:  p.Transaction,
			RequestID:    p.RequestID,
			Currency:     p.Currency,
			Provider:     p.Provider,
			Amount:       p.Amount,
			PaymentDT:    p.PaymentDT,
			Bank:         p.Bank,
			DeliveryCost: p.DeliveryCost,
			GoodsTotal:   p.GoodsTotal,
			CustomFee:    p.CustomFee,
		},
		Items:             items,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		ShardKey:          o.ShardKey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated,
		OofShard:          o.OofShard,
	}
}
