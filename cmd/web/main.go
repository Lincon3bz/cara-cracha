package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Evento struct {
	Nome string
	Data string
}

var eventos = []Evento{
	{Nome: "Churrasco dos primos", Data: "31/02/2206"},
	{Nome: "Aniversário do tatu", Data: "10/09/1990"},
}

func criarEvento(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("nome")
	data := r.FormValue("data")
	eventos = append(eventos, Evento{Nome: nome, Data: data})
	http.Redirect(w, r, "/eventos", http.StatusSeeOther)
}

func rotas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /eventos", listarEventos)
	mux.HandleFunc("POST /eventos", criarEvento)
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func listarEventos(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/eventos.html")
	if err != nil {
		http.Error(w, "Erro ao carregar a página", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, eventos)
	if err != nil {
		fmt.Println("Erro ao montar a página:", err)
	}
}

func main() {
	fmt.Println("Servidor ligado! http://localhost:8080")
	err := http.ListenAndServe(":8080", rotas())
	fmt.Println("Server parou:", err)

}
