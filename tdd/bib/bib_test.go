package bib

import (
	"testing"
)

func TestMediaLibrary_HasPermission(t *testing.T) {
	tests := []struct {
		name       string
		userRole   UserRole
		permission Permission
		want       bool
	}{
		// Admin Tests
		{"Admin kann Medien suchen", Admin, SearchMedia, true},
		{"Admin kann Medien ausleihen", Admin, BorrowMedia, true},
		{"Admin kann Medien reservieren", Admin, ReserveMedia, true},
		{"Admin kann Medien hinzufügen", Admin, AddMedia, true},
		{"Admin kann Medien bearbeiten", Admin, EditMedia, true},
		{"Admin kann Medien löschen", Admin, DeleteMedia, true},

		// StandardUser Tests
		{"StandardUser kann Medien suchen", StandardUser, SearchMedia, true},
		{"StandardUser kann Medien ausleihen", StandardUser, BorrowMedia, true},
		{"StandardUser kann Medien reservieren", StandardUser, ReserveMedia, true},
		{"StandardUser kann KEINE Medien hinzufügen", StandardUser, AddMedia, false},
		{"StandardUser kann KEINE Medien bearbeiten", StandardUser, EditMedia, false},
		{"StandardUser kann KEINE Medien löschen", StandardUser, DeleteMedia, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			user := ml.AddUser("testuser", tt.userRole)

			if got := ml.HasPermission(user.ID, tt.permission); got != tt.want {
				t.Errorf("HasPermission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMediaLibrary_AddMedia(t *testing.T) {
	tests := []struct {
		name     string
		userRole UserRole
		wantErr  bool
		errMsg   string
	}{
		{"Admin kann Medien hinzufügen", Admin, false, ""},
		{"StandardUser kann KEINE Medien hinzufügen", StandardUser, true, "keine Berechtigung zum Hinzufügen von Medien"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			user := ml.AddUser("testuser", tt.userRole)

			media, err := ml.AddMedia(user.ID, "Test Buch", "Test Autor")

			if tt.wantErr {
				if err == nil {
					t.Errorf("AddMedia() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("AddMedia() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("AddMedia() unexpected error = %v", err)
					return
				}
				if media == nil {
					t.Error("AddMedia() returned nil media")
				}
			}
		})
	}
}

func TestMediaLibrary_EditMedia(t *testing.T) {
	tests := []struct {
		name     string
		userRole UserRole
		wantErr  bool
		errMsg   string
	}{
		{"Admin kann Medien bearbeiten", Admin, false, ""},
		{"StandardUser kann KEINE Medien bearbeiten", StandardUser, true, "keine Berechtigung zum Bearbeiten von Medien"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			admin := ml.AddUser("admin", Admin)
			user := ml.AddUser("testuser", tt.userRole)

			// Erst ein Medium hinzufügen (als Admin)
			media, _ := ml.AddMedia(admin.ID, "Original Titel", "Original Autor")

			// Dann versuchen zu bearbeiten
			err := ml.EditMedia(user.ID, media.ID, "Neuer Titel", "Neuer Autor")

			if tt.wantErr {
				if err == nil {
					t.Errorf("EditMedia() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("EditMedia() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("EditMedia() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestMediaLibrary_DeleteMedia(t *testing.T) {
	tests := []struct {
		name     string
		userRole UserRole
		wantErr  bool
		errMsg   string
	}{
		{"Admin kann Medien löschen", Admin, false, ""},
		{"StandardUser kann KEINE Medien löschen", StandardUser, true, "keine Berechtigung zum Löschen von Medien"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			admin := ml.AddUser("admin", Admin)
			user := ml.AddUser("testuser", tt.userRole)

			// Erst ein Medium hinzufügen (als Admin)
			media, _ := ml.AddMedia(admin.ID, "Test Buch", "Test Autor")

			// Dann versuchen zu löschen
			err := ml.DeleteMedia(user.ID, media.ID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteMedia() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("DeleteMedia() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("DeleteMedia() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestMediaLibrary_BorrowMedia(t *testing.T) {
	tests := []struct {
		name        string
		userRole    UserRole
		mediaStatus string
		wantErr     bool
		errMsg      string
	}{
		{"Admin kann verfügbare Medien ausleihen", Admin, "available", false, ""},
		{"StandardUser kann verfügbare Medien ausleihen", StandardUser, "available", false, ""},
		{"Ausgeliehene Medien können nicht ausgeliehen werden", StandardUser, "borrowed", true, "Medium ist nicht verfügbar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			admin := ml.AddUser("admin", Admin)
			user := ml.AddUser("testuser", tt.userRole)

			// Medium hinzufügen
			media, _ := ml.AddMedia(admin.ID, "Test Buch", "Test Autor")
			media.Status = tt.mediaStatus

			// Ausleihen versuchen
			err := ml.BorrowMedia(user.ID, media.ID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("BorrowMedia() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("BorrowMedia() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("BorrowMedia() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestMediaLibrary_SearchMedia(t *testing.T) {
	tests := []struct {
		name        string
		userRole    UserRole
		query       string
		expectCount int
		wantErr     bool
	}{
		{"Admin kann alle Medien suchen", Admin, "", 2, false},
		{"StandardUser kann alle Medien suchen", StandardUser, "", 2, false},
		{"Suche nach spezifischem Titel", StandardUser, "Go Programming", 1, false},
		{"Suche nach Autor", StandardUser, "John Doe", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ml := NewMediaLibrary()
			admin := ml.AddUser("admin", Admin)
			user := ml.AddUser("testuser", tt.userRole)

			// Testmedien hinzufügen
			ml.AddMedia(admin.ID, "Go Programming", "John Doe")
			ml.AddMedia(admin.ID, "Python Basics", "Jane Smith")

			results, err := ml.SearchMedia(user.ID, tt.query)

			if tt.wantErr {
				if err == nil {
					t.Errorf("SearchMedia() error = nil, wantErr %v", tt.wantErr)
				}
			} else {
				if err != nil {
					t.Errorf("SearchMedia() unexpected error = %v", err)
					return
				}
				if len(results) != tt.expectCount {
					t.Errorf("SearchMedia() returned %d results, want %d", len(results), tt.expectCount)
				}
			}
		})
	}
}
