package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Evento struct {
	ID   int
	Nome string
	Data string
}

type formEvento struct {
	Nome  string
	Data  string
	Erros map[string]string
}

type paginaEventos struct {
	Eventos []Evento
	Form    formEvento
}

var repo RepositorioEventos = &repositorioMemoria{}

func mostrarEventos(w http.ResponseWriter, r *http.Request, status int, form formEvento) { //entrega o html da página /eventos
	tmpl, err := template.ParseFiles("templates/eventos.html")
	if err != nil {
		http.Error(w, "Erro ao carregar a página", http.StatusInternalServerError)
		return
	}

	eventos, err := repo.Listar(r.Context())
	if err != nil {
		http.Error(w, "Erro ao carregar eventos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	err = tmpl.Execute(w, paginaEventos{Eventos: eventos, Form: form})
	if err != nil {
		fmt.Println("erro ao montar a página:", err)
	}
}

func criarEvento(w http.ResponseWriter, r *http.Request) { //formulário e validação para criação de evento
	form := formEvento{
		Nome:  strings.TrimSpace(r.FormValue("nome")),
		Data:  strings.TrimSpace(r.FormValue("data")),
		Erros: map[string]string{},
	}
	if form.Nome == "" {
		form.Erros["nome"] = "Faltou o nome do evento "
	}
	if form.Data == "" { //TODO: tratar no futuro pra aceitar somente números (colocar '/' sozinho)
		form.Erros["data"] = "Informe a data" //TODO:tratar pra no futuro deixar enviar o formulário (checkbox "ainda não tenho data")

	} else {
		_, err := time.Parse("02/01/2006", form.Data)
		if err != nil {
			form.Erros["data"] = "Informe data válida e no formato DD/MM/AAAA"
		}
	}
	if len(form.Erros) > 0 {
		mostrarEventos(w, r, http.StatusUnprocessableEntity, form)
		return
	}
	err := repo.Adicionar(r.Context(), Evento{Nome: form.Nome, Data: form.Data})
	if err != nil {
		http.Error(w, "Erro ao registrar evento", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/eventos", http.StatusSeeOther)
}

func listarEventos(w http.ResponseWriter, r *http.Request) {
	mostrarEventos(w, r, http.StatusOK, formEvento{})
}

func rotas() http.Handler { //lista as possíveis rotas que o mux entrega
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handleHealth)
	mux.HandleFunc("GET /eventos", listarEventos)
	mux.HandleFunc("POST /eventos", criarEvento)
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}

func main() {
	//arq, err := novoRepositorioArquivo("dados/eventos.json")
	//if err != nil {
	//	log.Fatal("erro ao ler os dados: ", err)
	//}
	//repo = arq

	_ = godotenv.Load()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		log.Fatal("defina DATABASE_URL no arquivo .env")
	}

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatal("DATABASE_URL inválida", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal("não foi possível conectar DB: ", err)
	}
	repo = &repositorioPostgres{pool: pool}

	fmt.Println("Servidor ligado! http://localhost:8080")
	err = http.ListenAndServe(":8080", rotas())
	fmt.Println("Servidor parou:", err)

}
