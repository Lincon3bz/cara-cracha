-- +goose Up
CREATE TABLE eventos (
    id integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nome text NOT NULL,
    data text NOT NULL
);

-- +goose Down
DROP TABLE eventos;
