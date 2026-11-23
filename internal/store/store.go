package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/example/go-task-api/internal/model"
)

var (
	ErrNotFound = errors.New("task not found")
)

// Store defines the interface for task persistence.
// Using an interface makes the handlers easy to test and allows swapping storage later.
type Store interface {
	Create(ctx context.Context, title, description string) (*model.Task, error)
	Get(ctx context.Context, id string) (*model.Task, error)
	List(ctx context.Context, completed *bool) ([]*model.Task, error)
	Update(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id string) error
}

// MemoryStore is a concurrent-safe in-memory implementation of Store.
type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

// NewMemoryStore creates a new in-memory store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*model.Task),
	}
}

func (s *MemoryStore) Create(ctx context.Context, title, description string) (*model.Task, error) {
	// Respect context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	now := time.Now().UTC()
	task := &model.Task{
		ID:          newID(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()

	return task, nil
}

func (s *MemoryStore) Get(ctx context.Context, id string) (*model.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	task, ok := s.tasks[id]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrNotFound
	}

	// Return a copy to avoid external mutation
	copy := *task
	return &copy, nil
}

func (s *MemoryStore) List(ctx context.Context, completed *bool) ([]*model.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*model.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if completed != nil && t.Completed != *completed {
			continue
		}
		copy := *t
		result = append(result, &copy)
	}

	return result, nil
}

func (s *MemoryStore) Update(ctx context.Context, id string, req model.UpdateTaskRequest) (*model.Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Completed != nil {
		task.Completed = *req.Completed
	}
	task.UpdatedAt = time.Now().UTC()

	copy := *task
	return &copy, nil
}

func (s *MemoryStore) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}

	delete(s.tasks, id)
	return nil
}

// newID generates a random 16-byte hex ID (32 characters).
func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
