package handlers

import (
	"go_server_l0/api"
	"go_server_l0/internal/db"
	"go_server_l0/internal/mapper"
	"html/template"
	"log"
	"net/http"
)

type OrderController struct {
	appDb     db.DB
	indexTmpl *template.Template
}

// ServeHTTP implements http.Handler.
func (o *OrderController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}

	ctx := r.Context()
	query := r.URL.Query()
	uid := query.Get("uid")

	var apiRes api.Order
	if uid != "" {
		var order, delivery, pay, items, err = o.appDb.FindByUID(ctx, uid)
		if err != nil {
			log.Println("Cannot parse json", err.Error())
			apiRes = api.Order{}
		} else {
			apiRes = mapper.MapOrderDBToAPI(order, delivery, pay, items)
		}
	} else {
		apiRes = api.Order{}
	}

	err := o.indexTmpl.Execute(w, apiRes)

	if err != nil {
		log.Printf("Error executing template: %v", err.Error())
	}
}

func NewOrderControllerHandler(appDb db.DB, indexTmpl *template.Template) http.Handler {
	return &OrderController{
		appDb:     appDb,
		indexTmpl: indexTmpl,
	}
}
