package main

import (
	"ConcurrentLogProcessor/handler"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/process-logs", handler.ProcessLogsHandler)
	fmt.Println("Server started at :8080")

	err := http.ListenAndServe(":8099", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
