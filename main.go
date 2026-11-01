package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"example.com/online-shop/product"
)

func validateProduct(p product.Product) error {
	if len(strings.TrimSpace(p.Name)) == 0 {
		return errors.New("Name cannot be empty")
	}

	if p.Price <= 0 {
		return errors.New("Price should be positive")
	}

	if p.Stock < 0 {
		return errors.New("Stock cannot be negative")
	}

	return nil
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode products: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_, err = w.Write(encoded)
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func getProducts(w http.ResponseWriter, r *http.Request, miniData []product.Product) {
	id := r.URL.Query().Get("id")

	if len(strings.TrimSpace(id)) == 0 {
		writeJSON(w, http.StatusOK, miniData)
	} else {
		idNum, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid product id", http.StatusBadRequest)
			return
		}

		if idNum >= 1 && idNum <= len(miniData) {
			neededItem := miniData[idNum-1]

			writeJSON(w, http.StatusOK, neededItem)
		} else {
			http.Error(w, "Index out of range", http.StatusNotFound)
		}
	}
}

func createProduct(w http.ResponseWriter, r *http.Request, miniData []product.Product) []product.Product {
	var newProduct product.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return miniData
	}

	err = validateProduct(newProduct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return miniData
	}

	miniData = append(miniData, newProduct)

	writeJSON(w, http.StatusCreated, newProduct)

	return miniData
}

func updateProduct(w http.ResponseWriter, r *http.Request, miniData []product.Product) {
	var newProduct product.Product

	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")

	if len(strings.TrimSpace(id)) == 0 {
		http.Error(w, "No id given", http.StatusBadRequest)
		return
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Wrong id type", http.StatusBadRequest)
		return
	}

	if idNum < 1 || idNum > len(miniData) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	err = validateProduct(newProduct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	miniData[idNum-1] = newProduct

	writeJSON(w, http.StatusOK, newProduct)
}

func patchProduct(w http.ResponseWriter, r *http.Request, miniData []product.Product) {
	var update product.ProductUpdate

	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	id := r.URL.Query().Get("id")

	if len(strings.TrimSpace(id)) == 0 {
		http.Error(w, "No id given", http.StatusBadRequest)
		return
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	if idNum < 1 || idNum > len(miniData) {
		http.Error(w, "Index out of range", http.StatusNotFound)
		return
	}

	if update.Name == nil && update.Price == nil && update.Stock == nil {
		http.Error(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	if update.Name != nil {
		if len(strings.TrimSpace(*update.Name)) == 0 {
			http.Error(w, "Name should be provided", http.StatusBadRequest)
			return
		}
	}

	if update.Price != nil {
		if *update.Price <= 0 {
			http.Error(w, "Price should be positive", http.StatusBadRequest)
			return
		}
	}

	if update.Stock != nil {
		if *update.Stock < 0 {
			http.Error(w, "Stock cannot be negative", http.StatusBadRequest)
			return
		}
	}

	if update.Name != nil {
		miniData[idNum-1].Name = strings.TrimSpace(*update.Name)
	}

	if update.Price != nil {
		miniData[idNum-1].Price = *update.Price
	}

	if update.Stock != nil {
		miniData[idNum-1].Stock = *update.Stock
	}

	writeJSON(w, http.StatusOK, miniData[idNum-1])
}

func deleteProduct(w http.ResponseWriter, r *http.Request, miniData []product.Product) []product.Product {
	id := r.URL.Query().Get("id")

	if len(strings.TrimSpace(id)) == 0 {
		http.Error(w, "No id given", http.StatusBadRequest)
		return miniData
	}

	idNum, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Wrong id type", http.StatusBadRequest)
		return miniData
	}

	if idNum >= 1 && idNum <= len(miniData) {
		miniData = append(miniData[:idNum-1], miniData[idNum:]...)
	} else {
		http.Error(w, "Product not found", http.StatusNotFound)
		return miniData
	}

	w.WriteHeader(http.StatusNoContent)

	return miniData
}

func main() {
	miniData := []product.Product{
		{Name: "Laptop", Price: 800.00, Stock: 10},
		{Name: "Gaming Laptop", Price: 1200.00, Stock: 5},
		{Name: "Mouse", Price: 30.00, Stock: 50},
		{Name: "Keyboard", Price: 75.00, Stock: 25},
		{Name: "Monitor", Price: 250.00, Stock: 15},
	}

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProducts(w, r, miniData)

		case http.MethodPost:
			miniData = createProduct(w, r, miniData)

		case http.MethodDelete:
			miniData = deleteProduct(w, r, miniData)

		case http.MethodPut:
			updateProduct(w, r, miniData)

		case http.MethodPatch:
			patchProduct(w, r, miniData)

		default:
			http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
