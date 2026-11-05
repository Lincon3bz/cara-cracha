package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"slices"
	"sync"
)

type RepositorioEventos interface {
	Adicionar(e Evento) error
	Listar() []Evento
}

type repositorioMemoria struct {
	trava   sync.Mutex
	eventos []Evento
}

type repositorioArquivo struct {
	trava   sync.Mutex
	caminho string
	eventos []Evento
}

// ler ou criar arquivo
func novoRepositorioArquivo(caminho string) (*repositorioArquivo, error) {
	arq := &repositorioArquivo{caminho: caminho}

	conteudo, err := os.ReadFile(caminho)
	if errors.Is(err, fs.ErrNotExist) {
		return arq, nil
	}

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(conteudo, &arq.eventos)
	if err != nil {
		return nil, err
	}
	return arq, nil
}

func (arq *repositorioArquivo) Adicionar(e Evento) error {
	// travar, append, salvar; se salvar falhar, desfazer o append e devolver o erro
	return nil
}

func (arq *repositorioArquivo) Listar() []Evento {
	arq.trava.Lock()
	defer arq.trava.Unlock()
	return slices.Clone(arq.eventos)
}

func (arq *repositorioArquivo) salvar() error {
	// grava todos os eventos no arquivo (quem chama já está com a trava)
	return nil
}

func (m *repositorioMemoria) Adicionar(e Evento) error {
	m.trava.Lock()
	defer m.trava.Unlock()
	m.eventos = append(m.eventos, e)
	return nil
}

func (m *repositorioMemoria) Listar() []Evento {
	m.trava.Lock()
	defer m.trava.Unlock()
	return slices.Clone(m.eventos)
}

var _ RepositorioEventos = (*repositorioMemoria)(nil)

//var _ RepositorioEventos = (*repositorioArquivo)(nil)
