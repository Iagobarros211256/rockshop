package main

import (
	"log"
	"os"

	"github.com/Iagobarross211256/rockshop/internal/handlers"
	"github.com/Iagobarros211256/rockshop/internal/store"
)

// this was done by chat gpt. read it carefully
// this will be replaced on future comnmits beacase 
// is a manual db and server handling

func main() {
	dataFile := os.Getenv("DB_FILE")
	if datafile == "" {
		datafile = "data/store.json"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	st, err := store.NewJSONStore(datafile)
	if err != nil {
		log.Fatalf("faailed to init store: %v", err)
	}
}

r := gin.Default()

//handlers
ph := handlers.NewProductHandler(st)
oh := handlers.NewOrderHandler(st)
hh := handlers.NewHHealthCheckHandler()

api := r.Group("/api/v1")
{
	api.GET("/health", hh.Health)
	api.GET("/products", ph.ListProducts)
	api.GET("/products/:id", ph.GetProduct)
	api.POST("/products", ph.CreateProduct)
	api.PUT("/products/:id", ph.UpdateProduct)
	api.DELETE("/products/:id", ph.DeleteProduct)
	api.POST("/orders", oh.CreateOrder)
	api.GET("/orders/:id", oh.GetOrder)
	api.GET("/orders", oh.ListOrders)
}
log.Printf("RockShop starting on :%s (DB: %s)", port, dataFile)