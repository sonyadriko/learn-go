package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var items = []Item{
	{1, "Item 1", 19.99},
	{2, "Item 2", 29.99},
	{3, "Item 3", 39.99},
	{4, "Item 4", 49.99},
}

func getItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling GET request for items")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling GET request for a single item")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	for _, item := range items {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling POST request to create a new item")
	w.Header().Set("Content-Type", "application/json")
	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)
	newItem.ID = len(items) + 1
	items = append(items, newItem)
	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling PUT request to update an item")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var updatedItem Item
	_ = json.NewDecoder(r.Body).Decode(&updatedItem)

	for i, item := range items {
		if item.ID == id {
			items[i] = updatedItem
			json.NewEncoder(w).Encode(updatedItem)
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling DELETE request to delete an item")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			w.Write([]byte(`{"message":"Item deleted successfully"}`))
			return
		}
	}

	http.Error(w, "Item not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()

	// Use rs/cors middleware
	// handler := cors.Default().Handler(router)
	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(router)

	router.HandleFunc("/api/items", getItems).Methods("GET")
	router.HandleFunc("/api/items/{id}", getItem).Methods("GET")
	router.HandleFunc("/api/items", createItem).Methods("POST")
	router.HandleFunc("/api/items/{id}", updateItem).Methods("PUT")
	router.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
