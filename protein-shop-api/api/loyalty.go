package api

import (
    "context"
    "encoding/json"
    "net/http"

    "uniproj/db"
    "uniproj/models"
)

func LoyaltyPointsHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := r.Context().Value("user_id").(string)
    if !ok {
        http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
        return
    }

    rows, err := db.Conn.Query(context.Background(),
        "SELECT id, user_id, points, source, created_at FROM loyalty_points WHERE user_id = $1 ORDER BY created_at DESC", userID)
    if err != nil {
        http.Error(w, `{"error":"Failed to fetch points"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var history []models.LoyaltyPoint
    for rows.Next() {
        var lp models.LoyaltyPoint
        if err := rows.Scan(&lp.ID, &lp.UserID, &lp.Point ubs, &lp.Source, &lp.CreatedAt); err != nil {
            http.Error(w, `{"error":"Failed to scan points"}`, http.StatusInternalServerError)
            return
        }
        history = append(history, lp)
    }

    var balance int
    err = db.Conn.QueryRow(context.Background(),
        "SELECT COALESCE(SUM(points), 0) FROM loyalty_points WHERE user_id = $1", userID).
        Scan(&balance)
    if err != nil {
        http.Error(w, `{"error":"Failed to calculate balance"}`, http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "balance": balance,
        "history": history,
    })
}