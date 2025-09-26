package bib

import (
	"errors"
	"fmt"
)

type UserRole int

const (
	StandardUser UserRole = iota
	Admin
)

func (ur UserRole) String() string {
	switch ur {
	case StandardUser:
		return "StandardUser"
	case Admin:
		return "Admin"
	default:
		return "Unknown"
	}
}

type Permission int

const (
	SearchMedia Permission = iota
	BorrowMedia
	ReserveMedia
	AddMedia
	EditMedia
	DeleteMedia
)

func (p Permission) String() string {
	switch p {
	case SearchMedia:
		return "SearchMedia"
	case BorrowMedia:
		return "BorrowMedia"
	case ReserveMedia:
		return "ReserveMedia"
	case AddMedia:
		return "AddMedia"
	case EditMedia:
		return "EditMedia"
	case DeleteMedia:
		return "DeleteMedia"
	default:
		return "Unknown"
	}
}

type User struct {
	ID       int
	Username string
	Role     UserRole
}

type Media struct {
	ID         int
	Title      string
	Author     string
	Status     string
	BorrowedBy *int
}

type MediaLibrary struct {
	users       map[int]*User
	media       map[int]*Media
	nextUserID  int
	nextMediaID int
}

func NewMediaLibrary() *MediaLibrary {
	return &MediaLibrary{
		users:       make(map[int]*User),
		media:       make(map[int]*Media),
		nextUserID:  1,
		nextMediaID: 1,
	}
}

func (ml *MediaLibrary) AddUser(username string, role UserRole) *User {
	user := &User{
		ID:       ml.nextUserID,
		Username: username,
		Role:     role,
	}

	ml.users[ml.nextUserID] = user
	ml.nextUserID++

	return user
}

func (ml *MediaLibrary) HasPermission(userID int, permission Permission) bool {
	user, exists := ml.users[userID]
	if !exists {
		return false
	}

	switch user.Role {
	case Admin:
		return true
	case StandardUser:
		switch permission {
		case SearchMedia, BorrowMedia, ReserveMedia:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func (ml *MediaLibrary) AddMedia(userID int, title, author string) (*Media, error) {
	if !ml.HasPermission(userID, AddMedia) {
		return nil, errors.New("keine Berechtigung zum Hinzufügen von Medien")
	}

	media := &Media{
		ID:     ml.nextMediaID,
		Title:  title,
		Author: author,
		Status: "available",
	}
	ml.media[ml.nextMediaID] = media
	ml.nextMediaID++
	return media, nil
}

func (ml *MediaLibrary) EditMedia(userID, mediaID int, title, author string) error {
	if !ml.HasPermission(userID, EditMedia) {
		return errors.New("keine Berechtigung zum Bearbeiten von Medien")
	}

	media, exists := ml.media[mediaID]
	if !exists {
		return errors.New("Medium nicht gefunden")
	}

	if title != "" {
		media.Title = title
	}
	if author != "" {
		media.Author = author
	}
	return nil
}

func (ml *MediaLibrary) DeleteMedia(userID, mediaID int) error {
	if !ml.HasPermission(userID, DeleteMedia) {
		return errors.New("keine Berechtigung zum Löschen von Medien")
	}

	if _, exists := ml.media[mediaID]; !exists {
		return errors.New("Medium nicht gefunden")
	}

	delete(ml.media, mediaID)
	return nil
}

func (ml *MediaLibrary) SearchMedia(userID int, query string) ([]*Media, error) {
	if !ml.HasPermission(userID, SearchMedia) {
		return nil, errors.New("keine Berechtigung zum Suchen von Medien")
	}

	var results []*Media
	for _, media := range ml.media {
		if query == "" ||
			fmt.Sprintf("%s %s", media.Title, media.Author) == query ||
			media.Title == query || media.Author == query {
			results = append(results, media)
		}
	}
	return results, nil
}

func (ml *MediaLibrary) BorrowMedia(userID, mediaID int) error {
	if !ml.HasPermission(userID, BorrowMedia) {
		return errors.New("keine Berechtigung zum Ausleihen von Medien")
	}

	media, exists := ml.media[mediaID]
	if !exists {
		return errors.New("Medium nicht gefunden")
	}

	if media.Status != "available" {
		return errors.New("Medium ist nicht verfügbar")
	}

	media.Status = "borrowed"
	media.BorrowedBy = &userID
	return nil
}

func (ml *MediaLibrary) ReserveMedia(userID, mediaID int) error {
	if !ml.HasPermission(userID, ReserveMedia) {
		return errors.New("keine Berechtigung zum Reservieren von Medien")
	}

	media, exists := ml.media[mediaID]
	if !exists {
		return errors.New("Medium nicht gefunden")
	}

	if media.Status == "borrowed" {
		media.Status = "reserved"
		return nil
	}

	return errors.New("Medium kann nicht reserviert werden")
}
