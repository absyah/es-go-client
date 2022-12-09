package main

import (
	"log"
	"net/http"

	"es-go-client/delivery"
	ESClient "es-go-client/repository"
)

func main() {
	esClient := ESClient.NewClient("http://localhost:9200")
	err := esClient.Ping()
	if err != nil {
		log.Fatal("failed to connect to Elasticsearch: ", err)
	}

	server := delivery.Server{ESClient: esClient}
	http.HandleFunc("/insert", server.InsertIndexHandler)
	http.HandleFunc("/update", server.UpdateIndexHandler)
	http.HandleFunc("/delete", server.DeleteIndexHandler)
	http.HandleFunc("/search", server.SearchHandler)
	http.HandleFunc("/ping", server.PingHandler)

	log.Println("listening server on port 8080")
	http.ListenAndServe(":8080", nil)
}
