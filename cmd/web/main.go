package main

import (
	"fmt"
	"net/http"
)

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "_ - F U N C I O N A N D O - _")
}

func main() {
	fmt.Println("_-* I N I C I O U   O   S I S T E M A *-_")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	fmt.Println("Servidor ligado! http://localhost:8080")
	err := http.ListenAndServe(":8080", mux)
	fmt.Println("Server parou:", err)

}
