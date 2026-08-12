package main

import (
	"errors"
	"strings"
	"time"
)

type Status string

const (
	StatusTodo       Status = "a_fazer"
	StatusInProgress Status = "em_progresso"
	StatusDone       Status = "concluida"
)

var validStatuses = map[Status]bool{
	StatusTodo:       true,
	StatusInProgress: true,
	StatusDone:       true,
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

func (t *TaskInput) Validate(requireStatus bool) error {
	t.Title = strings.TrimSpace(t.Title)
	if t.Title == "" {
		return errors.New("o campo 'title' é obrigatório")
	}
	if t.Status == "" {
		if requireStatus {
			t.Status = StatusTodo
		}
		return nil
	}
	if !validStatuses[t.Status] {
		return errors.New("status inválido: use 'a_fazer', 'em_progresso' ou 'concluida'")
	}
	return nil
}
