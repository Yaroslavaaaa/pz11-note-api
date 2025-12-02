package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"example.com/notes-api/internal/core"
	"example.com/notes-api/internal/repo"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Repo *repo.NoteRepoMem
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UpdateNoteRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// CreateNote создает новую заметку
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n core.Note

	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if strings.TrimSpace(n.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Title is required")
		return
	}

	id, err := h.Repo.Create(n)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	createdNote, err := h.Repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve created note")
		return
	}

	respondWithJSON(w, http.StatusCreated, createdNote)
}

// GetNote возвращает заметку по ID
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := h.Repo.GetByID(id)
	if err != nil {
		if err == repo.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Note not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to get note")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, note)
}

// GetAllNotes возвращает все заметки
func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.Repo.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get notes")
		return
	}

	if notes == nil {
		notes = []core.Note{}
	}

	respondWithJSON(w, http.StatusOK, notes)
}

// PatchNote - частичное обновление (PATCH)
func (h *Handler) PatchNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var update UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if update.Title == nil && update.Content == nil {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	if update.Title != nil && strings.TrimSpace(*update.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	updates := make(map[string]interface{})
	if update.Title != nil {
		updates["title"] = *update.Title
	}
	if update.Content != nil {
		updates["content"] = *update.Content
	}

	err = h.Repo.UpdatePartial(id, updates)
	if err != nil {
		if err == repo.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Note not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to update note")
		}
		return
	}

	updatedNote, err := h.Repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve updated note")
		return
	}

	respondWithJSON(w, http.StatusOK, updatedNote)
}

// DeleteNote удаляет заметку
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	err = h.Repo.Delete(id)
	if err != nil {
		if err == repo.ErrNoteNotFound {
			respondWithError(w, http.StatusNotFound, "Note not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Failed to delete note")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, SuccessResponse{
		Message: "Note deleted successfully",
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(payload)
}
