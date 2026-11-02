package main

import "testing"

func TestRepositorioMemoria(t *testing.T) {
	repo := &repositorioMemoria{}
	repo.Adicionar(Evento{Nome: "TESTE_EVENTO_TESTE 1", Data: "01/01/2026"})
	repo.Adicionar(Evento{Nome: "TESTE_EVENTO_TESTE 2", Data: "02/02/2026"})

	eventos := repo.Listar()
	if len(eventos) != 2 {
		t.Errorf("esperava 2 eventos, vieram %d: %v", len(eventos), eventos)
	}
}
