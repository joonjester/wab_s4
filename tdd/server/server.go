package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tdd/bib"
)

// HTTP Request/Response Structs
type CreateUserRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" oder "standard"
}

type CreateMediaRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type EditMediaRequest struct {
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}

type BorrowRequest struct {
	MediaID int `json:"media_id"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// MediaLibraryServer wraps MediaLibrary mit HTTP-Funktionalität
type MediaLibraryServer struct {
	Ml *bib.MediaLibrary
}

func NewMediaLibraryServer() *MediaLibraryServer {
	return &MediaLibraryServer{
		Ml: bib.NewMediaLibrary(),
	}
}

// Helper function to get user ID from header
func (s *MediaLibraryServer) getUserIDFromHeader(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, errors.New("X-User-ID header erforderlich")
	}
	return strconv.Atoi(userIDStr)
}

// Helper function to send JSON response
func (s *MediaLibraryServer) sendJSON(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// POST /users - Benutzer erstellen
func (s *MediaLibraryServer) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur POST erlaubt",
		})
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültiger JSON",
		})
		return
	}

	var role bib.UserRole
	switch strings.ToLower(req.Role) {
	case "admin":
		role = bib.Admin
	case "standard":
		role = bib.StandardUser
	default:
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültige Rolle. Verwenden Sie 'admin' oder 'standard'",
		})
		return
	}

	user := s.Ml.AddUser(req.Username, role)
	s.sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    user,
	})
}

// POST /media - Medium hinzufügen (nur Admins)
func (s *MediaLibraryServer) createMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur POST erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req CreateMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültiger JSON",
		})
		return
	}

	media, err := s.Ml.AddMedia(userID, req.Title, req.Author)
	if err != nil {
		s.sendJSON(w, http.StatusForbidden, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    media,
	})
}

// PUT /media/{id} - Medium bearbeiten (nur Admins)
func (s *MediaLibraryServer) editMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur PUT erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Extract media ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/media/")
	mediaID, err := strconv.Atoi(path)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültige Medium-ID",
		})
		return
	}

	var req EditMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültiger JSON",
		})
		return
	}

	err = s.Ml.EditMedia(userID, mediaID, req.Title, req.Author)
	if err != nil {
		s.sendJSON(w, http.StatusForbidden, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    "Medium erfolgreich bearbeitet",
	})
}

// DELETE /media/{id} - Medium löschen (nur Admins)
func (s *MediaLibraryServer) deleteMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur DELETE erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Extract media ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/media/")
	mediaID, err := strconv.Atoi(path)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültige Medium-ID",
		})
		return
	}

	err = s.Ml.DeleteMedia(userID, mediaID)
	if err != nil {
		s.sendJSON(w, http.StatusForbidden, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    "Medium erfolgreich gelöscht",
	})
}

// GET /media - Medien suchen (alle Benutzer)
func (s *MediaLibraryServer) searchMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur GET erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	query := r.URL.Query().Get("q")
	results, err := s.Ml.SearchMedia(userID, query)
	if err != nil {
		s.sendJSON(w, http.StatusForbidden, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    results,
	})
}

// POST /borrow - Medium ausleihen (alle Benutzer)
func (s *MediaLibraryServer) borrowMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur POST erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req BorrowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültiger JSON",
		})
		return
	}

	err = s.Ml.BorrowMedia(userID, req.MediaID)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    "Medium erfolgreich ausgeliehen",
	})
}

// POST /reserve - Medium reservieren (alle Benutzer)
func (s *MediaLibraryServer) reserveMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   "Nur POST erlaubt",
		})
		return
	}

	userID, err := s.getUserIDFromHeader(r)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	var req BorrowRequest // Reuse same struct
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "Ungültiger JSON",
		})
		return
	}

	err = s.Ml.ReserveMedia(userID, req.MediaID)
	if err != nil {
		s.sendJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    "Medium erfolgreich reserviert",
	})
}

// Router setup
func (s *MediaLibraryServer) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/users", s.createUser)
	mux.HandleFunc("/media", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			s.searchMedia(w, r)
		} else if r.Method == http.MethodPost {
			s.createMedia(w, r)
		} else {
			s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
				Success: false,
				Error:   "Methode nicht erlaubt",
			})
		}
	})
	mux.HandleFunc("/media/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			s.editMedia(w, r)
		} else if r.Method == http.MethodDelete {
			s.deleteMedia(w, r)
		} else {
			s.sendJSON(w, http.StatusMethodNotAllowed, APIResponse{
				Success: false,
				Error:   "Methode nicht erlaubt",
			})
		}
	})
	mux.HandleFunc("/borrow", s.borrowMedia)
	mux.HandleFunc("/reserve", s.reserveMedia)

	return mux
}
