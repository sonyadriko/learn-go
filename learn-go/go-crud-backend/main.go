package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *gorm.DB

type Item struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Price     float64    `json:"price"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index;type:timestamp"`
}

func init() {
	var err error
	dsn := "root:@tcp(localhost: 3306)/test_go"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// Test the connection
	err = db.Exec("SELECT 1").Error
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to MySQL")
}

func getItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling GET request for items")
	w.Header().Set("Content-Type", "application/json")

	var items []Item
	db.Find(&items)

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

	var item Item
	result := db.First(&item, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling POST request to create a new item")
	w.Header().Set("Content-Type", "application/json")

	var newItem Item
	_ = json.NewDecoder(r.Body).Decode(&newItem)
	db.Create(&newItem)

	json.NewEncoder(w).Encode(newItem)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling PUT request to update an item")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		http.Error(w, "Item ID is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	var existingItem Item
	result := db.First(&existingItem, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Item not found", http.StatusNotFound)
		} else {
			fmt.Println("Error retrieving existing item:", result.Error)
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	var updatedItem Item
	if err := json.NewDecoder(r.Body).Decode(&updatedItem); err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("Updating item with ID %d\n", id)
	fmt.Printf("Updated item: %+v\n", updatedItem)

	// Update only the specified fields
	result = db.Model(&existingItem).Updates(updatedItem)
	if result.Error != nil {
		fmt.Println("Error updating item:", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated item from the database
	result = db.First(&existingItem, id)
	if result.Error != nil {
		fmt.Println("Error retrieving updated item:", result.Error)
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(existingItem)
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling DELETE request to delete an item")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		http.Error(w, "Item ID is missing", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Delete the item from the MySQL database
	result := db.Delete(&Item{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(`{"message":"Item deleted successfully"}`))
}

func deletea(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling DELETE request to soft delete an item")
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	idStr, ok := params["id"]
	if !ok {
		http.Error(w, "Item ID is missing", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}
	current := time.Now()
	// Soft delete by updating the deleted_at field
	db1 := db.Model(Item{}).Where("id = ?", id).Updates(map[string]interface{}{
		"deleted_at": &current,
	})
	if db1.Error != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}
	// item := Item{ID: int32(id)}
	// db.Save(&item) // update the item with a SoftDeletedAt timestamp
}

// Di sisi server (Go)
func searchItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling GET request for search items")
	w.Header().Set("Content-Type", "application/json")

	// Ambil nilai-nilai pencarian dari query parameter "search"
	searchParam := r.URL.Query().Get("search")
	if searchParam == "" {
		http.Error(w, "Parameter 'search' is required", http.StatusBadRequest)
		return
	}

	// Buat variabel untuk menyimpan hasil pencarian
	var items []Item

	// Lakukan pencarian berdasarkan nilai pencarian
	result := db.Where("name LIKE ?", "%"+searchParam+"%").Find(&items)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirim hasil pencarian sebagai respons
	json.NewEncoder(w).Encode(items)
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
	router.HandleFunc("/api/itemss/{id}", deletea).Methods("PUT")
	router.HandleFunc("/api/items/{id}", deleteItem).Methods("DELETE")
	router.HandleFunc("/api/search/items", searchItems).Methods("GET")

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
