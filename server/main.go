package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// One row per date for your stacked chart
type SessionsSeries struct {
	Date      string `json:"date"`
	Classical int64  `json:"classical"`
	Quantum   int64  `json:"quantum"`
}

var db *pgxpool.Pool

func main() {
	// Build DSN from env vars
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getenv("PGUSER", "user"),
		getenv("PGPASSWORD", "password"),
		getenv("PGHOST", "localhost"),
		getenv("PGPORT", "5432"),
		getenv("PGDATABASE", "rag_db"),
		getenv("PGSSLMODE", "disable"),
	)

	// Connect pool
	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer db.Close()

	// Routes
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/powerball/sessions", sessionsHandler)

	log.Println("🚀 Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Hello from Go API 👋")
}

// Returns JSON like:
// [
//
//	{"date":"2024-12-18","classical":648,"quantum":3},
//	{"date":"2024-12-19","classical":79,"quantum":1},
//	{"date":"2025-04-05","classical":241,"quantum":2}
//
// ]
func sessionsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := db.Query(ctx, `
		SELECT
		  generated_at::date AS d,
		  SUM(CASE WHEN source = 'classical' THEN 1 ELSE 0 END) AS classical,
		  SUM(CASE WHEN source = 'quantum'   THEN 1 ELSE 0 END) AS quantum
		FROM powerball_sessions
		GROUP BY 1
		ORDER BY 1;
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []SessionsSeries
	for rows.Next() {
		var d time.Time
		var classical, quantum int64
		if err := rows.Scan(&d, &classical, &quantum); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		out = append(out, SessionsSeries{
			Date:      d.Format("2006-01-02"),
			Classical: classical,
			Quantum:   quantum,
		})
	}
	if rows.Err() != nil {
		http.Error(w, rows.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}
