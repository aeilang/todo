package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/aeilang/backend/store"
)

type Server struct {
	query *store.Queries
}

func NewServer(query *store.Queries) *Server {
	return &Server{
		query: query,
	}
}

type TodoDTO struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"create_at"`
}

// GET /todos
func (s *Server) HandleGetTodos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := s.query.GetTodos(ctx)
	if err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	todos := make([]TodoDTO, len(rows))
	for i, row := range rows {
		todo := TodoDTO{
			Id:        int(row.ID),
			Content:   row.Content,
			Completed: row.Completed,
			CreatedAt: row.CreateAt,
		}
		todos[i] = todo
	}

	SendJSON(w, http.StatusOK, todos)
}

// PUT /todo/:id
func (s *Server) HandleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	if err := s.query.UpdateTodoCompleted(ctx, int32(id)); err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// DELETE /todo/:id
func (s *Server) HandleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		SendError(w, http.StatusBadRequest, errors.New("id is required"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	if err := s.query.DeleteTodo(ctx, int32(id)); err != nil {
		SendError(w, http.StatusBadRequest, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type AddTodoPayload struct {
	Content string `json:"content"`
}

// POST /todo
func (s *Server) HandleAddTodo(w http.ResponseWriter, r *http.Request) {
	var payload AddTodoPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	if payload.Content == "" {
		SendError(w, http.StatusBadRequest, errors.New("content is requried"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.query.CreateTodo(ctx, payload.Content); err != nil {
		SendError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created successful!"))
}
