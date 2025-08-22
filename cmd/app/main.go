package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_server_l0/api"
	"go_server_l0/internal/db"
	"go_server_l0/internal/kafkaservice"
	"go_server_l0/internal/mapper"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/segmentio/kafka-go"
)

const apiPrefixV1 = "/api/v1"

func main() {
	// init DB
	var dbUrl = "postgres://user:password@localhost:5432/db"
	var appDb, err = db.ConnectToPGSQL(dbUrl)

	if err != nil {
		panic("DB connection error " + err.Error())
	}

	db.AddCaching(&appDb)
	defer appDb.Close()
	// end init DB

	// init kafka
	kafkaService, err := kafkaservice.ConnectToKafks(
		[]string{"localhost:29092"},
		"orders",
		"go-reader",
		context.Background())

	if err != nil {
		log.Println("Error connecting to kafka", err.Error())
	}

	kafkaService.RegisterListener(func(msg kafka.Message) {
		log.Println("New message")

		var req api.Order
		json.Unmarshal(msg.Value, &req)

		var order, delivery, pay, items = mapper.MapAPIOrderToDB(req)

		uid := order.UID
		delivery.OrderUID = uid
		pay.OrderUID = uid

		var ctx = context.Background()
		appDb.SaveOrder(ctx, order)
		appDb.SaveDelivery(ctx, delivery)
		appDb.SavePayment(ctx, pay)

		for _, item := range items {
			item.OrderUID = uid
			appDb.SaveItem(ctx, item)
		}
	})

	go kafkaService.Listen()
	// end init Kafka

	// templates
	indexTmpl := template.Must(template.ParseFiles("templates/index.html"))

	// static
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	loadOrder := func(query url.Values, ctx context.Context) (api.Order, error) {
		uid := query.Get("uid")

		if uid == "" {
			return api.Order{}, nil
		}

		var order, delivery, pay, items, err = appDb.FindByUID(ctx, uid)

		if err != nil {
			log.Println("Cannot parse json", err.Error())
			return api.Order{}, err
		}

		return mapper.MapOrderDBToAPI(order, delivery, pay, items), nil
	}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.HandleFunc(apiPrefixV1+"/orders", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		switch r.Method {
		case "POST":
			var req api.Order
			var err = json.NewDecoder(r.Body).Decode(&req)

			if err != nil {
				log.Println("Cannot decode", err.Error())
				w.WriteHeader(500)
				return
			}

			var order, delivery, pay, items = mapper.MapAPIOrderToDB(req)

			uid := order.UID
			delivery.OrderUID = uid
			pay.OrderUID = uid

			appDb.SaveOrder(ctx, order)
			appDb.SaveDelivery(ctx, delivery)
			appDb.SavePayment(ctx, pay)

			for _, item := range items {
				item.OrderUID = uid
				appDb.SaveItem(ctx, item)
			}
		case "GET":
			query := r.URL.Query()
			apiRes, err := loadOrder(query, ctx)

			if err != nil {
				log.Println("Cannot load order", err.Error())
				w.WriteHeader(500)
			}

			json.NewEncoder(w).Encode(apiRes)
		default:
			w.WriteHeader(404)
		}
	})

	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		apiRes, err := loadOrder(r.URL.Query(), r.Context())

		if err != nil {
			log.Println("Cannot load oreds,", err.Error())
			w.WriteHeader(500)
		}

		err = indexTmpl.Execute(w, apiRes)

		if err != nil {
			log.Printf("Error executing template: %v", err.Error())
		}
	})

	log.Println("Starting server")
	http.ListenAndServe(":8000", nil)
	log.Println("Finished")
}
