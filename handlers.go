package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu       sync.RWMutex
	tasks    map[string]Task
	filePath string
}

func NewStore(filePath string) *Store {
	s := &Store{
		tasks:    make(map[string]Task),
		filePath: filePath,
	}
	s.load()
	return s
}

func (s *Store) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *Store) Get(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *Store) Create(in TaskInput) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	t := Task{
		ID:          newID(),
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[t.ID] = t
	s.persist()
	return t
}

func (s *Store) Update(id string, in TaskInput) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	t.Title = in.Title
	t.Description = in.Description
	t.Status = in.Status
	t.UpdatedAt = time.Now()
	s.tasks[id] = t
	s.persist()
	return t, true
}

func (s *Store) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return false
	}
	delete(s.tasks, id)
	s.persist()
	return true
}

func (s *Store) persist() {
	if s.filePath == "" {
		return
	}
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return
	}
	_ = writeFile(s.filePath, data)
}

func (s *Store) load() {
	if s.filePath == "" {
		return
	}
	data, err := readFile(s.filePath)
	if err != nil {
		return
	}
	var list []Task
	if err := json.Unmarshal(data, &list); err != nil {
		return
	}
	for _, t := range list {
		s.tasks[t.ID] = t
	}
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func tasksHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, store.List())

		case http.MethodPost:
			var in TaskInput
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
				return
			}
			if err := in.Validate(true); err != nil {
				writeError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			task := store.Create(in)
			writeJSON(w, http.StatusCreated, task)

		default:
			writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		}
	}
}

func taskByIDHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")
		id = strings.Trim(id, "/")
		if id == "" {
			writeError(w, http.StatusBadRequest, "id da tarefa não informado")
			return
		}

		switch r.Method {
		case http.MethodPut:
			existing, ok := store.Get(id)
			if !ok {
				writeError(w, http.StatusNotFound, "tarefa não encontrada")
				return
			}
			var in TaskInput
			if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
				writeError(w, http.StatusBadRequest, "corpo da requisição inválido")
				return
			}
			if in.Status == "" {
				in.Status = existing.Status
			}
			if err := in.Validate(false); err != nil {
				writeError(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			updated, _ := store.Update(id, in)
			writeJSON(w, http.StatusOK, updated)

		case http.MethodDelete:
			if ok := store.Delete(id); !ok {
				writeError(w, http.StatusNotFound, "tarefa não encontrada")
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodGet:
			task, ok := store.Get(id)
			if !ok {
				writeError(w, http.StatusNotFound, "tarefa não encontrada")
				return
			}
			writeJSON(w, http.StatusOK, task)

		default:
			writeError(w, http.StatusMethodNotAllowed, "método não permitido")
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
