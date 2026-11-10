package main

import (
	"os"
	"path/filepath"
	"testing"
)

// listar chama repo.Listar e encerra o teste se der erro.
func listar(t *testing.T, repo RepositorioEventos) []Evento {
	t.Helper()                               // se falhar, o erro aponta a linha do teste, não esta
	eventos, err := repo.Listar(t.Context()) // t.Context(): um context que dura enquanto o teste roda
	if err != nil {
		t.Fatal(err)
	}
	return eventos
}

func TestRepositorioMemoria(t *testing.T) {
	repo := &repositorioMemoria{}
	repo.Adicionar(t.Context(), Evento{Nome: "TESTE_EVENTO_TESTE 1", Data: "01/01/2026"})
	repo.Adicionar(t.Context(), Evento{Nome: "TESTE_EVENTO_TESTE 2", Data: "02/02/2026"})

	eventos := listar(t, repo)
	if len(eventos) != 2 {
		t.Errorf("esperava 2 eventos, vieram %d: %v", len(eventos), eventos)
	}
}
func TestRepositorioArquivoSobreviveAoReinicio(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "eventos.json")

	primeiro, err := novoRepositorioArquivo(caminho)
	if err != nil {
		t.Fatal(err)
	}
	err = primeiro.Adicionar(t.Context(), Evento{Nome: "Churrasco dos primos", Data: "31/10/2026"})
	if err != nil {
		t.Fatal(err)
	}
	// "reinicia o servidor": um repositório novo, lendo o mesmo arquivo
	segundo, err := novoRepositorioArquivo(caminho)
	if err != nil {
		t.Fatal(err)
	}
	eventos := listar(t, segundo)
	if len(eventos) != 1 || eventos[0].Nome != "Churrasco dos primos" {
		t.Errorf("depois de reabrir: esperava só o Churrasco dos primos, veio %v", eventos)
	}
}

func TestRepositorioArquivoInexistente(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "eventos.json") // pasta temporária, sem nenhum arquivo dentro

	arq, err := novoRepositorioArquivo(caminho)
	if err != nil {
		t.Fatal("arquivo que não existe não deveria ser erro:", err)
	}
	if len(listar(t, arq)) != 0 {
		t.Errorf("esperava começar vazio, vieram %v", listar(t, arq))
	}
}

func TestRepositorioArquivoLeEventosSalvos(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "eventos.json")
	conteudo := `[{"Nome": "Churrasco dos primos", "Data": "31/10/2026"},
	          {"Nome": "Natal da família", "Data": "25/12/2026"}]`
	os.WriteFile(caminho, []byte(conteudo), 0o644) // simula um arquivo salvo numa execução anterior

	arq, err := novoRepositorioArquivo(caminho)
	if err != nil {
		t.Fatal(err)
	}
	eventos := listar(t, arq)
	if len(eventos) != 2 || eventos[0].Nome != "Churrasco dos primos" {
		t.Errorf("esperava os 2 eventos do arquivo, vieram %v", eventos)
	}
}

func TestRepositorioMemoriaGeraIDs(t *testing.T) {
	repo := &repositorioMemoria{}
	repo.Adicionar(t.Context(), Evento{Nome: "Churrasco dos primos", Data: "31/10/2026"})
	repo.Adicionar(t.Context(), Evento{Nome: "Natal da família", Data: "25/12/2026"})

	eventos := listar(t, repo)
	if eventos[0].ID != 1 || eventos[1].ID != 2 {
		t.Errorf("esperava IDs 1 e 2, vieram %d e %d", eventos[0].ID, eventos[1].ID)
	}
}

func TestRepositorioArquivoIDsContinuamDepoisDoReinicio(t *testing.T) {
	caminho := filepath.Join(t.TempDir(), "eventos.json")

	primeiro, _ := novoRepositorioArquivo(caminho)
	primeiro.Adicionar(t.Context(), Evento{Nome: "Churrasco dos primos", Data: "31/10/2026"}) // ganha o ID 1

	segundo, _ := novoRepositorioArquivo(caminho) // "reinicia o servidor"
	segundo.Adicionar(t.Context(), Evento{Nome: "Natal da família", Data: "25/12/2026"})

	eventos := listar(t, segundo)
	if eventos[1].ID != 2 {
		t.Errorf("depois do reinício, esperava o ID 2, veio %d", eventos[1].ID)
	}
}
