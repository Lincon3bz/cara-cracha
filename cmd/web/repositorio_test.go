package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositorioMemoria(t *testing.T) {
	repo := &repositorioMemoria{}
	repo.Adicionar(Evento{Nome: "TESTE_EVENTO_TESTE 1", Data: "01/01/2026"})
	repo.Adicionar(Evento{Nome: "TESTE_EVENTO_TESTE 2", Data: "02/02/2026"})

	eventos := repo.Listar()
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
	err = primeiro.Adicionar(Evento{Nome: "Churrasco dos primos", Data: "31/10/2026"})
	if err != nil {
		t.Fatal(err)
	}

	// "reinicia o servidor": um repositório novo, lendo o mesmo arquivo
	segundo, err := novoRepositorioArquivo(caminho)
	if err != nil {
		t.Fatal(err)
	}
	eventos := segundo.Listar()
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
	if len(arq.Listar()) != 0 {
		t.Errorf("esperava começar vazio, vieram %v", arq.Listar())
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
	eventos := arq.Listar()
	if len(eventos) != 2 || eventos[0].Nome != "Churrasco dos primos" {
		t.Errorf("esperava os 2 eventos do arquivo, vieram %v", eventos)
	}
}
