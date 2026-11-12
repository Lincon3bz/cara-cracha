package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
)

type RepositorioEventos interface {
	Adicionar(ctx context.Context, e Evento) error
	Listar(ctx context.Context) ([]Evento, error)
}

func proximoID(eventos []Evento) int {
	maior := 0
	for _, e := range eventos {
		if e.ID > maior {
			maior = e.ID
		}
	}
	return maior + 1
}

type repositorioMemoria struct {
	trava   sync.Mutex
	eventos []Evento
}

var _ RepositorioEventos = (*repositorioMemoria)(nil) //compilador confere se struct satisfaz interface

func (m *repositorioMemoria) Adicionar(_ context.Context, e Evento) error {
	m.trava.Lock()
	defer m.trava.Unlock()
	e.ID = proximoID(m.eventos)
	m.eventos = append(m.eventos, e)
	return nil
}

func (m *repositorioMemoria) Listar(context.Context) ([]Evento, error) {
	m.trava.Lock()
	defer m.trava.Unlock()
	return slices.Clone(m.eventos), nil
}

type repositorioArquivo struct {
	trava   sync.Mutex
	caminho string
	eventos []Evento
}

var _ RepositorioEventos = (*repositorioArquivo)(nil) //compilador confere se struct satisfaz interface

// abre o repositório lendo o arquivo, se existir
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

func (arq *repositorioArquivo) Listar(_ context.Context) ([]Evento, error) {
	arq.trava.Lock()
	defer arq.trava.Unlock()
	return slices.Clone(arq.eventos), nil
}

func (arq *repositorioArquivo) salvar() error {
	conteudo, err := json.MarshalIndent(arq.eventos, "", "    ")
	if err != nil {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(arq.caminho), 0o755)
	if err != nil {
		return nil
	}

	return os.WriteFile(arq.caminho, conteudo, 0o644)
}

func (arq *repositorioArquivo) Adicionar(_ context.Context, e Evento) error {
	arq.trava.Lock()
	defer arq.trava.Unlock()

	e.ID = proximoID(arq.eventos)
	arq.eventos = append(arq.eventos, e)
	err := arq.salvar()
	if err != nil {
		arq.eventos = arq.eventos[:len(arq.eventos)-1]
		return err
	}
	return nil
}
