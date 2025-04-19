package api

import (
	"context"
	"encoding/json"
	"net/http"

	"uniproj/db"
	"uniproj/models"
)

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Conn.Query(context.Background(), "SELECT id, name, description, price, stock, created_at FROM products")
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch products"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt); err != nil {
			http.Error(w, `{"error":"Failed to scan products"}`, http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
