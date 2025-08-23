package handlers

import (
	"encoding/json"
	"go_server_l0/internal/db"
	"go_server_l0/internal/mapper"
	"log"
	"net/http"
)

type OrderHandler struct {
	appDb db.DB
}

// ServeHTTP implements http.Handler.
func (o *OrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	ctx := r.Context()
	query := r.URL.Query()
	uid := query.Get("uid")

	if uid == "" {
		w.WriteHeader(400)
		return
	}

	order, delivery, pay, items, err := o.appDb.FindByUID(ctx, uid)

	if err != nil {
		log.Println("Cannot parse json", err.Error())
		w.WriteHeader(400)
	}

	apiRes := mapper.MapOrderDBToAPI(order, delivery, pay, items)

	json.NewEncoder(w).Encode(apiRes)
}

func NewOrdersHandler(appDb db.DB) http.Handler {
	return &OrderHandler{appDb: appDb}
}
