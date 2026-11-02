package main

import (
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

func (m *repositorioMemoria) Adicionar(e Evento) error {
	m.trava.Lock()
	defer m.trava.Unlock()
	eventos = append(m.eventos, Evento{Nome: e.Nome, Data: e.Data})
	return nil
}

func (m *repositorioMemoria) Listar() []Evento {
	m.trava.Lock()
	defer m.trava.Unlock()
	return slices.Clone(eventos)

}

var _ RepositorioEventos = (*repositorioMemoria)(nil)
