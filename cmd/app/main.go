package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go_server_l0/api"
	"go_server_l0/internal/config"
	"go_server_l0/internal/db"
	"go_server_l0/internal/handlers"
	"go_server_l0/internal/kafkaservice"
	"go_server_l0/internal/mapper"
	"html/template"
	"log"
	"net/http"

	"github.com/segmentio/kafka-go"
)

const apiPrefixV1 = "/api/v1"

func main() {
	// load config
	appConfig := config.LoadConfig()

	// init DB
	var appDb, err = db.ConnectToPGSQL(appConfig.DbUrl)

	if err != nil {
		panic("DB connection error " + err.Error())
	}

	db.AddCaching(&appDb, context.Background())
	defer appDb.Close()
	// end init DB

	// init kafka
	kafkaService, err := kafkaservice.ConnectToKafks(
		appConfig.Brokers,
		appConfig.Topic,
		appConfig.GroupId,
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
	indexTmpl := template.Must(template.ParseFiles(appConfig.IndexTemplate))

	// static
	fs := http.FileServer(http.Dir(appConfig.StaticDir))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "pong")
	})

	http.Handle(apiPrefixV1+"/orders", handlers.NewOrdersHandler(appDb))
	http.Handle("/index.html", handlers.NewOrderControllerHandler(appDb, indexTmpl))

	log.Println("Starting server")
	http.ListenAndServe(appConfig.Addr, nil)
	log.Println("Finished")
}
