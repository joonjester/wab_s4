package main

import (
	"fmt"
	"log"
	"net/http"
	"tdd/bib"
	"tdd/server"
)

func main() {
	fmt.Println("TDD application is getting started")

	newServer := server.NewMediaLibraryServer()

	// Beispieldaten erstellen
	admin := newServer.Ml.AddUser("admin", bib.Admin)
	user := newServer.Ml.AddUser("user", bib.StandardUser)

	// Beispielmedien hinzufügen
	newServer.Ml.AddMedia(admin.ID, "Go Programming", "John Doe")
	newServer.Ml.AddMedia(admin.ID, "Python Basics", "Jane Smith")

	fmt.Printf("Server gestartet auf :8080\n")
	fmt.Printf("Beispielbenutzer:\n")
	fmt.Printf("- Admin (ID: %d)\n", admin.ID)
	fmt.Printf("- Standard User (ID: %d)\n", user.ID)
	fmt.Println("\nAPI Endpunkte:")
	fmt.Println("POST /users - Benutzer erstellen")
	fmt.Println("GET /media?q=query - Medien suchen")
	fmt.Println("POST /media - Medium hinzufügen (nur Admins)")
	fmt.Println("PUT /media/{id} - Medium bearbeiten (nur Admins)")
	fmt.Println("DELETE /media/{id} - Medium löschen (nur Admins)")
	fmt.Println("POST /borrow - Medium ausleihen")
	fmt.Println("POST /reserve - Medium reservieren")
	fmt.Println("\nVerwenden Sie X-User-ID Header für Authentifizierung")

	mux := newServer.SetupRoutes()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

