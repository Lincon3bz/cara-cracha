package main

import (
	"fmt"
	"net/http"
)

func rotas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /eventos", listarEventos)
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
func listarEventos(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "eventos📋")
}

func main() {
	fmt.Println("_-* I N I C I O U   O   S I S T E M A *-_")
	fmt.Println("Servidor ligado! http://localhost:8080")
	err := http.ListenAndServe(":8080", rotas())
	fmt.Println("Server parou:", err)

}
