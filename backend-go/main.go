package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	frontendDir = "../frontend"
	timeout     = 10
)

type Shift struct {
	ShiftID  int       `json:"shift_id"`
	WorkDate time.Time `json:"work_date"`
	Hours    float64   `json:"hours"`
	Rate     float64   `json:"rate"`
}

func getShifts(ctx context.Context, pool *pgxpool.Pool) ([]Shift, error) {
	rows, err := pool.Query(ctx, "SELECT shift_id, work_date, hours, rate FROM shifts")
	if err != nil {
		return nil, fmt.Errorf("Запрос к shifts упал: %w", err)
	}
	defer rows.Close()

	var shifts []Shift
	for rows.Next() {
		var s Shift
		if err := rows.Scan(&s.ShiftID, &s.WorkDate, &s.Hours, &s.Rate); err != nil {
			return nil, fmt.Errorf("не распаковал строку: %w", err)
		}
		shifts = append(shifts, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Хождение по строкам оборвалось: %w", err)
	}
	return shifts, nil
}

func readHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := map[string]string{
		"message": "zdarova ubludok",
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ошибка записи ответа: %v", err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeout)
	defer cancel()
	pool, err := connectDB(ctx)
	if err != nil {
		log.Fatalf("Не удалось достучаться до БД: %v\n", err)
	}
	stat := pool.Stat()
	log.Printf("пул: всего %d, свободно %d, максимум %d",
		stat.TotalConns(), stat.IdleConns(), stat.MaxConns())
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hp", readHealth)
	mux.HandleFunc("GET /shifts", func(w http.ResponseWriter, r *http.Request) {
		shifts, err := getShifts(ctx, pool)
		if err != nil {
			http.Error(w, "не смог получить смены", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shifts)
	})

	mux.Handle("GET /", http.FileServer(http.Dir(frontendDir)))

	log.Println("listen :8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatal(err)
	}

}
