package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"uniproj/db"
	"uniproj/models"

	"github.com/google/uuid"
)

type createOrderRequest struct {
	Items []struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	} `json:"items"`
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		http.Error(w, `{"error":"Server error"}`, http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	orderID := uuid.New()
	total := 0.0
	for _, item := range req.Items {
		var price float64
		var stock int
		err = tx.QueryRow(context.Background(),
			"SELECT price, stock FROM products WHERE id = $1", item.ProductID).
			Scan(&price, &stock)
		if err != nil {
			http.Error(w, `{"error":"Invalid product"}`, http.StatusBadRequest)
			return
		}
		if stock < item.Quantity {
			http.Error(w, `{"error":"Product out of stock"}`, http.StatusBadRequest)
			return
		}
		total += price * float64(item.Quantity)

		_, err = tx.Exec(context.Background(),
			"UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			http.Error(w, `{"error":"Failed to update stock"}`, http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec(context.Background(),
			"INSERT INTO order_items (id, order_id, product_id, quantity, price_at_purchase) VALUES ($1, $2, $3, $4, $5)",
			uuid.New(), orderID, item.ProductID, item.Quantity, price)
		if err != nil {
			http.Error(w, `{"error":"Failed to create order item"}`, http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec(context.Background(),
		"INSERT INTO orders (id, user_id, status, total, created_at) VALUES ($1, $2, $3, $4, $5)",
		orderID, userID, "pending", total, time.Now())
	if err != nil {
		http.Error(w, `{"error":"Failed to create order"}`, http.StatusInternalServerError)
		return
	}

	points := int(total / 10)
	if points > 0 {
		_, err = tx.Exec(context.Background(),
			"INSERT INTO loyalty_points (id, user_id, points, source, created_at) VALUES ($1, $2, $3, $4, $5)",
			uuid.New(), userID, points, "purchase", time.Now())
		if err != nil {
			http.Error(w, `{"error":"Failed to award points"}`, http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, `{"error":"Failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id":      orderID,
		"total":         total,
		"status":        "pending",
		"points_earned": points,
	})
}

func GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	rows, err := db.Conn.Query(context.Background(),
		`SELECT o.id, o.user_id, o.status, o.total, o.created_at,
                oi.id, oi.product_id, oi.quantity, oi.price_at_purchase, p.name
         FROM orders o
         LEFT JOIN order_items oi ON o.id = oi.order_id
         LEFT JOIN products p ON oi.product_id = p.id
         WHERE o.user_id = $1
         ORDER BY o.created_at DESC`, userID)
	if err != nil {
		http.Error(w, `{"error":"Failed to fetch orders"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	ordersMap := make(map[uuid.UUID]*models.Order)
	for rows.Next() {
		var o models.Order
		var oi models.OrderItem
		var oiID, productID uuid.UUID
		var productName string

		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt,
			&oiID, &productID, &oi.Quantity, &oi.PriceAtPurchase, &productName); err != nil {
			http.Error(w, `{"error":"Failed to scan orders"}`, http.StatusInternalServerError)
			return
		}

		oi.ID = oiID
		oi.OrderID = o.ID
		oi.ProductID = productID
		oi.ProductName = productName

		if existing, exists := ordersMap[o.ID]; exists {
			existing.Items = append(existing.Items, oi)
		} else {
			o.Items = []models.OrderItem{oi}
			ordersMap[o.ID] = &o
		}
	}

	var orders []models.Order
	for _, o := range ordersMap {
		orders = append(orders, *o)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
