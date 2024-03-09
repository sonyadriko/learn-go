package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var db *sql.DB

func init() {
	// Open a database connection
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/test_go")
	if err != nil {
		panic(err.Error())
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL")
}

// var items = []Item{
// 	{1, "Item 1", 19.99},
// 	{2, "Item 2", 29.99},
// 	{3, "Item 3", 39.99},
// 	{4, "Item 4", 49.99},
// }

// func getItems(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Handling GET request for items")
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(items)
// }

func getItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling GET request for items")
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name, &item.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

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

	row := db.QueryRow("SELECT * FROM items WHERE id=?", id)

	var item Item
	err = row.Scan(&item.ID, &item.Name, &item.Price)
	if err == sql.ErrNoRows {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling POST request to create a new item")
	w.Header().Set("Content-Type", "application/json")

	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)

	// Insert the new item into the MySQL database
	result, err := db.Exec("INSERT INTO items (name, price) VALUES (?, ?)", newItem.Name, newItem.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted item
	// newItem.ID, _ = result.LastInsertId()

	// Get the ID of the newly inserted item
	newItemID, _ := result.LastInsertId()

	// Assign the obtained ID to the newItem
	newItem.ID = int(newItemID)

	// Return the newly created item in the response
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

	// Update the item in the MySQL database
	_, err = db.Exec("UPDATE items SET name=?, price=? WHERE id=?", updatedItem.Name, updatedItem.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the updated item in the response
	json.NewEncoder(w).Encode(updatedItem)
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

	// Delete the item from the MySQL database
	_, err = db.Exec("DELETE FROM items WHERE id=?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message":"Item deleted successfully"}`))
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
