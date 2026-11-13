package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"example.com/online-shop/database"
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

func getProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")

		if id == "" {
			products, err := database.GetProducts(db)
			if err != nil {
				log.Printf("Failed to retrieve products: %v", err)
				http.Error(w, "Failed to retrieve products", http.StatusInternalServerError)
				return
			}

			writeJSON(w, http.StatusOK, products)
		} else {
			idNum, err := strconv.Atoi(id)
			if err != nil || idNum < 1 {
				http.Error(w, "Invalid product ID", http.StatusBadRequest)
				return
			}

			product, err := database.GetProductsByID(db, idNum)
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Product Not Found", http.StatusNotFound)
				return
			}

			if err != nil {
				log.Printf("Failed to retrieve product: %v", err)
				http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
				return
			}

			writeJSON(w, http.StatusOK, product)
		}
	}
}

func createProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
			return
		}

		var newProduct product.Product

		err := json.NewDecoder(r.Body).Decode(&newProduct)
		if err != nil {
			http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
			return
		}

		err = validateProduct(newProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		productID, err := database.CreateProduct(db, newProduct)
		if err != nil {
			log.Printf("Could not create product: %v", err)
			http.Error(w, "Could not create product", http.StatusInternalServerError)
			return
		}

		newProduct.ProductID = productID

		writeJSON(w, http.StatusCreated, newProduct)
	}
}

func updateProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var updatedProduct product.Product

		id := r.URL.Query().Get("id")

		if id == "" {
			http.Error(w, "No ID provided", http.StatusBadRequest)
			return
		}

		idNum, err := strconv.Atoi(id)
		if err != nil || idNum < 1 {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&updatedProduct)
		if err != nil {
			http.Error(w, "Invalid JSON body request", http.StatusBadRequest)
			return
		}

		updatedProduct.ProductID = idNum

		err = validateProduct(updatedProduct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = database.UpdateProduct(db, updatedProduct)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No requested product found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("Failed to update product: %v", err)
			http.Error(w, "Could not update product", http.StatusInternalServerError)
			return
		}

		updatedProduct, err = database.GetProductsByID(db, idNum)
		if err != nil {
			log.Printf("Failed to retrieve updated product: %v", err)
			http.Error(w, "Failed to retrieve updated product", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, updatedProduct)
	}
}

func patchProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
		if err != nil || idNum < 1 {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
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

		err = database.PatchProduct(db, update, idNum)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("Failed to patch product: %v", err)
			http.Error(w, "Could not update product", http.StatusInternalServerError)
			return
		}

		updatedProduct, err := database.GetProductsByID(db, idNum)
		if err != nil {
			log.Printf("Failed to retrieve updated product: %v", err)
			http.Error(w, "Failed to retrieve updated product", http.StatusInternalServerError)
			return
		}

		writeJSON(w, http.StatusOK, updatedProduct)
	}
}

func deleteProduct(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")

		if len(strings.TrimSpace(id)) == 0 {
			http.Error(w, "No id given", http.StatusBadRequest)
			return
		}

		idNum, err := strconv.Atoi(id)
		if err != nil || idNum < 1 {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		err = database.DeleteProduct(db, idNum)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		if err != nil {
			log.Printf("Failed to delete product: %v", err)
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Successfully connected to PostgreSQL!")

	http.HandleFunc("/products", getProducts(db))
	http.HandleFunc("/add/product", createProduct(db))
	http.HandleFunc("/change/product", updateProduct(db))
	http.HandleFunc("/update/product", patchProduct(db))
	http.HandleFunc("/delete/product", deleteProduct(db))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
