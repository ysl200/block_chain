package main

import (
	"blockchain/controller"
	"blockchain/storage"
	"blockchain/web"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.NewSource(time.Now().UnixNano())

	nodeController := controller.NewNodeController()

	log.Println("服务启动: http://localhost:8080")

	http.HandleFunc("/add", nodeController.HandleAddNode)
	http.HandleFunc("/store", storage.HandleStoreData)
	http.HandleFunc("/list", nodeController.HandleListNodes)

	web.StartWebServer()

	http.ListenAndServe(":8080", nil)
}
