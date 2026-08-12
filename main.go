package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
)

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte("fallback-id-0000"))
	}
	return hex.EncodeToString(b)
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

func main() {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "tasks.json"
	}

	store := NewStore(dataFile)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", tasksHandler(store))
	mux.HandleFunc("/tasks/", taskByIDHandler(store))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	handler := corsMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Kanban API rodando em http://localhost:%s (dados em %s)", port, dataFile)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
