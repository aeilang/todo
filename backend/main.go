package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/aeilang/backend/store"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("POSTGRE_URL")
	if dbUrl == "" {
		panic("dbURL is required")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}

	query := store.New(db)
	srv := NewServer(query)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /todos", srv.HandleGetTodos)
	mux.HandleFunc("PUT /todo/{id}", srv.HandleUpdateTodo)
	mux.HandleFunc("DELETE /todo/{id}", srv.HandleDeleteTodo)
	mux.HandleFunc("POST /todo", srv.HandleAddTodo)

	serv := http.Server{
		Addr:    ":8080",
		Handler: CORS(mux),
	}

	fmt.Println("listen on", serv.Addr)

	if err := serv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func CORS(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
