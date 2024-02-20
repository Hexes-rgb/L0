package main

import (
	"fmt"
	"net/http"
)

func ordersHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "orders.html")
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderData []byte
	var ok bool
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	if orderData, ok = cache.Get(id); !ok {
		if orderData, ok = getOrder(id); !ok {
			http.Error(w, "Data not found", http.StatusNotFound)
			return
		}
	}
	jsonData := addUUIDToJson(id, orderData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func createServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/orders", ordersHandler)
	http.HandleFunc("/getOrder", getOrderHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orders" && r.URL.Path != "/getOrder" {
			http.Redirect(w, r, "/orders", http.StatusSeeOther)
			return
		}
		http.NotFound(w, r)
	})

	fmt.Println("The server is running on http://localhost:" + config.ServicePort + "/orders")
	if err := http.ListenAndServe(":"+config.ServicePort, nil); err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}
