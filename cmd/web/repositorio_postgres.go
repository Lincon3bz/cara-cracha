package main

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repositorioPostgres struct {
	pool *pgxpool.Pool
}

var _ RepositorioEventos = (*repositorioPostgres)(nil)

func (p *repositorioPostgres) Listar(ctx context.Context) ([]Evento, error) {
	linhas, err := p.pool.Query(ctx, "SELECT id, nome , data FROM eventos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var eventos []Evento
	for linhas.Next() {
		var e Evento
		err := linhas.Scan(&e.ID, &e.Nome, &e.Data)
		if err != nil {
			return nil, err
		}
		eventos = append(eventos, e)
	}
	return eventos, linhas.Err()
}

func (p *repositorioPostgres) Adicionar(_ context.Context, _ Evento) error {
	return errors.New("adicionar no postgres ainda não foi implementado")
}
