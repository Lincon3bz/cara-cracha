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
	eventos = append(eventos, Evento{Nome: form.Nome, Data: form.Data})

}

func (m *repositorioMemoria) Listar() []Evento {
	m.trava.Lock()
	defer m.trava.Unlock()
	foto := slices.Clone(eventos)
}

var _ RepositorioEventos = (*repositorioMemoria)(nil)
