package main

import (
	"fmt"
	"net/http"
)

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderData []byte
	var ok bool
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Parameter 'id' not specified", http.StatusBadRequest)
		return
	}
	orderData, ok = cache.Get(id)
	if !ok {
		if orderData, ok = getOrder(id); !ok {
			fmt.Fprintf(w, "Data not found")
			return
		}
	}
	jsonData := addUUIDToJson(id, orderData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createServer() {
	http.HandleFunc("/getOrder", getOrderHandler)

	fmt.Println("The server is running on http://localhost:8080/")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}
