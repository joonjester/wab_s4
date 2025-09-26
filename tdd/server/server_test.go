package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"tdd/bib"
	"testing"
)

func TestHTTPServer_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectSuccess  bool
		expectError    string
	}{
		{
			name:           "Admin erstellen erfolgreich",
			requestBody:    `{"username": "testadmin", "role": "admin"}`,
			expectedStatus: http.StatusCreated,
			expectSuccess:  true,
		},
		{
			name:           "Standard User erstellen erfolgreich",
			requestBody:    `{"username": "testuser", "role": "standard"}`,
			expectedStatus: http.StatusCreated,
			expectSuccess:  true,
		},
		{
			name:           "Ungültige Rolle",
			requestBody:    `{"username": "testuser", "role": "invalid"}`,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    "Ungültige Rolle. Verwenden Sie 'admin' oder 'standard'",
		},
		{
			name:           "Ungültiger JSON",
			requestBody:    `{"invalid json"}`,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    "Ungültiger JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			server.createUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError != "" && response.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, response.Error)
			}
		})
	}
}

func TestHTTPServer_CreateMedia(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userRole       bib.UserRole
		requestBody    string
		expectedStatus int
		expectSuccess  bool
		expectError    string
	}{
		{
			name:           "Admin kann Medium erstellen",
			userID:         "1",
			userRole:       bib.Admin,
			requestBody:    `{"title": "Test Buch", "author": "Test Autor"}`,
			expectedStatus: http.StatusCreated,
			expectSuccess:  true,
		},
		{
			name:           "Standard User kann kein Medium erstellen",
			userID:         "2",
			userRole:       bib.StandardUser,
			requestBody:    `{"title": "Test Buch", "author": "Test Autor"}`,
			expectedStatus: http.StatusForbidden,
			expectSuccess:  false,
			expectError:    "keine Berechtigung zum Hinzufügen von Medien",
		},
		{
			name:           "Fehlende User-ID",
			userID:         "",
			userRole:       bib.Admin,
			requestBody:    `{"title": "Test Buch", "author": "Test Autor"}`,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    "X-User-ID header erforderlich",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			// Benutzer erstellen
			if tt.userID != "" {
				server.Ml.AddUser("testuser", tt.userRole)
			}

			req := httptest.NewRequest(http.MethodPost, "/media", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}
			w := httptest.NewRecorder()

			server.createMedia(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError != "" && response.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, response.Error)
			}
		})
	}
}

func TestHTTPServer_SearchMedia(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userRole       bib.UserRole
		query          string
		expectedStatus int
		expectSuccess  bool
		expectedCount  int
	}{
		{
			name:           "Admin kann alle Medien suchen",
			userID:         "1",
			userRole:       bib.Admin,
			query:          "",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
			expectedCount:  2,
		},
		{
			name:           "Standard User kann alle Medien suchen",
			userID:         "2",
			userRole:       bib.StandardUser,
			query:          "",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
			expectedCount:  2,
		},
		{
			name:           "Suche mit Query",
			userID:         "1",
			userRole:       bib.Admin,
			query:          "Go Programming",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
			expectedCount:  1,
		},
		{
			name:           "Suche nach Autor",
			userID:         "1",
			userRole:       bib.Admin,
			query:          "John Doe",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			// Benutzer erstellen
			admin := server.Ml.AddUser("admin", bib.Admin)
			server.Ml.AddUser("user", tt.userRole)

			// Testmedien hinzufügen
			server.Ml.AddMedia(admin.ID, "Go Programming", "John Doe")
			server.Ml.AddMedia(admin.ID, "Python Basics", "Jane Smith")

			// URL mit korrekter Query-Parameter-Kodierung erstellen
			testURL := "/media"
			if tt.query != "" {
				testURL += "?q=" + url.QueryEscape(tt.query)
			}

			req := httptest.NewRequest(http.MethodGet, testURL, nil)
			req.Header.Set("X-User-ID", tt.userID)
			w := httptest.NewRecorder()

			server.searchMedia(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectSuccess && response.Data != nil {
				results, ok := response.Data.([]interface{})
				if !ok {
					t.Error("Expected data to be array")
				} else if len(results) != tt.expectedCount {
					t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
				}
			}
		})
	}
}

func TestHTTPServer_BorrowMedia(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userRole       bib.UserRole
		requestBody    string
		expectedStatus int
		expectSuccess  bool
		expectError    string
	}{
		{
			name:           "Admin kann Medium ausleihen",
			userID:         "1",
			userRole:       bib.Admin,
			requestBody:    `{"media_id": 1}`,
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "Nicht verfügbares Medium",
			userID:         "1",
			userRole:       bib.Admin,
			requestBody:    `{"media_id": 1}`,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    "Medium ist nicht verfügbar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			// Benutzer erstellen
			admin := server.Ml.AddUser("admin", bib.Admin)
			server.Ml.AddUser("user", tt.userRole)

			// Testmedien hinzufügen
			media, _ := server.Ml.AddMedia(admin.ID, "Buch 1", "Autor 1")

			// Erstes Medium als ausgeliehen markieren für Test
			if tt.name == "Nicht verfügbares Medium" {
				media.Status = "borrowed"
			}

			req := httptest.NewRequest(http.MethodPost, "/borrow", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", tt.userID)
			w := httptest.NewRecorder()

			server.borrowMedia(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError != "" && response.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, response.Error)
			}
		})
	}
}

func TestHTTPServer_EditMedia(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userRole       bib.UserRole
		mediaID        string
		requestBody    string
		expectedStatus int
		expectSuccess  bool
		expectError    string
	}{
		{
			name:           "Admin kann Medium bearbeiten",
			userID:         "1",
			userRole:       bib.Admin,
			mediaID:        "1",
			requestBody:    `{"title": "Neuer Titel", "author": "Neuer Autor"}`,
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "Standard User kann kein Medium bearbeiten",
			userID:         "2",
			userRole:       bib.StandardUser,
			mediaID:        "1",
			requestBody:    `{"title": "Neuer Titel"}`,
			expectedStatus: http.StatusForbidden,
			expectSuccess:  false,
			expectError:    "keine Berechtigung zum Bearbeiten von Medien",
		},
		{
			name:           "Ungültige Medium-ID",
			userID:         "1",
			userRole:       bib.Admin,
			mediaID:        "invalid",
			requestBody:    `{"title": "Neuer Titel"}`,
			expectedStatus: http.StatusBadRequest,
			expectSuccess:  false,
			expectError:    "Ungültige Medium-ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			// Benutzer erstellen
			admin := server.Ml.AddUser("admin", bib.Admin)
			server.Ml.AddUser("user", tt.userRole)

			// Testmedium hinzufügen
			server.Ml.AddMedia(admin.ID, "Original Titel", "Original Autor")

			req := httptest.NewRequest(http.MethodPut, "/media/"+tt.mediaID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-User-ID", tt.userID)
			w := httptest.NewRecorder()

			server.editMedia(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError != "" && response.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, response.Error)
			}
		})
	}
}

func TestHTTPServer_DeleteMedia(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		userRole       bib.UserRole
		mediaID        string
		expectedStatus int
		expectSuccess  bool
		expectError    string
	}{
		{
			name:           "Admin kann Medium löschen",
			userID:         "1",
			userRole:       bib.Admin,
			mediaID:        "1",
			expectedStatus: http.StatusOK,
			expectSuccess:  true,
		},
		{
			name:           "Standard User kann kein Medium löschen",
			userID:         "2",
			userRole:       bib.StandardUser,
			mediaID:        "1",
			expectedStatus: http.StatusForbidden,
			expectSuccess:  false,
			expectError:    "keine Berechtigung zum Löschen von Medien",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewMediaLibraryServer()

			// Benutzer erstellen
			admin := server.Ml.AddUser("admin", bib.Admin)
			server.Ml.AddUser("user", tt.userRole)

			// Testmedium hinzufügen
			server.Ml.AddMedia(admin.ID, "Test Buch", "Test Autor")

			req := httptest.NewRequest(http.MethodDelete, "/media/"+tt.mediaID, nil)
			req.Header.Set("X-User-ID", tt.userID)
			w := httptest.NewRecorder()

			server.deleteMedia(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, response.Success)
			}

			if tt.expectError != "" && response.Error != tt.expectError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectError, response.Error)
			}
		})
	}
}

func TestHTTPServer_Integration(t *testing.T) {
	// Vollständiger Integrationstest
	server := NewMediaLibraryServer()

	// 1. Admin erstellen
	req := httptest.NewRequest(http.MethodPost, "/users",
		bytes.NewBufferString(`{"username": "admin", "role": "admin"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.createUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create admin: %d", w.Code)
	}

	// 2. Standard User erstellen
	req = httptest.NewRequest(http.MethodPost, "/users",
		bytes.NewBufferString(`{"username": "user", "role": "standard"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.createUser(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create user: %d", w.Code)
	}

	// 3. Admin fügt Medium hinzu
	req = httptest.NewRequest(http.MethodPost, "/media",
		bytes.NewBufferString(`{"title": "Integration Test Book", "author": "Test Author"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "1")
	w = httptest.NewRecorder()
	server.createMedia(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Admin failed to create media: %d", w.Code)
	}

	// 4. User sucht Medien
	req = httptest.NewRequest(http.MethodGet, "/media", nil)
	req.Header.Set("X-User-ID", "2")
	w = httptest.NewRecorder()
	server.searchMedia(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("User failed to search media: %d", w.Code)
	}

	// 5. User leiht Medium aus
	req = httptest.NewRequest(http.MethodPost, "/borrow",
		bytes.NewBufferString(`{"media_id": 1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "2")
	w = httptest.NewRecorder()
	server.borrowMedia(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("User failed to borrow media: %d", w.Code)
	}

	// 6. User versucht Medium zu löschen (sollte fehlschlagen)
	req = httptest.NewRequest(http.MethodDelete, "/media/1", nil)
	req.Header.Set("X-User-ID", "2")
	w = httptest.NewRecorder()
	server.deleteMedia(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("User should not be able to delete media, got: %d", w.Code)
	}

	t.Log("Integration test passed successfully!")
}
