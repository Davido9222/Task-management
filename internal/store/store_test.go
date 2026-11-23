package store

import (
	"context"
	"testing"

	"github.com/example/go-task-api/internal/model"
)

func TestMemoryStore_CRUD(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	// Create
	task, err := s.Create(ctx, "Buy milk", "2 liters")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if task.ID == "" {
		t.Error("expected non-empty ID")
	}
	if task.Title != "Buy milk" {
		t.Errorf("expected title 'Buy milk', got %q", task.Title)
	}
	if task.Completed {
		t.Error("new task should not be completed")
	}

	// Get
	got, err := s.Get(ctx, task.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if got.Title != task.Title {
		t.Errorf("Get returned wrong title: %q", got.Title)
	}

	// List
	list, err := s.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 task, got %d", len(list))
	}

	// Update
	completed := true
	newTitle := "Buy oat milk"
	updated, err := s.Update(ctx, task.ID, model.UpdateTaskRequest{
		Title:     &newTitle,
		Completed: &completed,
	})
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Title != "Buy oat milk" {
		t.Errorf("expected updated title, got %q", updated.Title)
	}
	if !updated.Completed {
		t.Error("expected completed=true")
	}

	// Filter completed
	trueVal := true
	completedList, err := s.List(ctx, &trueVal)
	if err != nil {
		t.Fatalf("List completed failed: %v", err)
	}
	if len(completedList) != 1 {
		t.Errorf("expected 1 completed task, got %d", len(completedList))
	}

	// Delete
	if err := s.Delete(ctx, task.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Get after delete
	_, err = s.Get(ctx, task.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_NotFound(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStore()

	_, err := s.Get(ctx, "non-existent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}

	err = s.Delete(ctx, "non-existent")
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}
