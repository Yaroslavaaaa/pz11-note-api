package repo

import (
	"errors"
	"sync"
	"time"

	"example.com/notes-api/internal/core"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

type NoteRepoMem struct {
	mu    sync.RWMutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{
		notes: make(map[int64]*core.Note),
		next:  1,
	}
}

func (r *NoteRepoMem) Create(n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n.ID = r.next
	n.CreatedAt = time.Now()
	n.UpdatedAt = nil
	r.notes[n.ID] = &n
	r.next++

	return n.ID, nil
}

func (r *NoteRepoMem) GetByID(id int64) (*core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, exists := r.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}

	noteCopy := *note
	return &noteCopy, nil
}

func (r *NoteRepoMem) GetAll() ([]core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notes := make([]core.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, *note)
	}

	return notes, nil
}

func (r *NoteRepoMem) UpdatePartial(id int64, updates map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	note, exists := r.notes[id]
	if !exists {
		return ErrNoteNotFound
	}

	if title, ok := updates["title"].(string); ok && title != "" {
		note.Title = title
	}

	if content, ok := updates["content"].(string); ok {
		note.Content = content
	}

	now := time.Now()
	note.UpdatedAt = &now

	return nil
}

func (r *NoteRepoMem) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.notes[id]; !exists {
		return ErrNoteNotFound
	}

	delete(r.notes, id)
	return nil
}
