
# cara-crachá

Servidor web em Go. Por enquanto, ele mostra uma página com a lista de eventos e um formulário para cadastrar novos.

## Rotas

| Método | Caminho    | O que faz                                                                 |
|--------|------------|---------------------------------------------------------------------------|
| GET    | `/health`  | Responde `ok`                                                             |
| GET    | `/eventos` | Página HTML com a lista de eventos e o formulário de novo evento          |
| POST   | `/eventos` | Valida o formulário; se válido, cria o evento e redireciona para `/eventos` |

## Como rodar

Requer Go 1.26.1 (versão declarada no `go.mod`). Rode **a partir da raiz do repositório**, porque o template é carregado por caminho relativo:

```bash
go run ./cmd/web
```

O servidor sobe em `http://localhost:8080`.

## Testes

```bash
go test -race ./...
go test ./... -count=1
```

- `TestHealth` — confere status 200 e corpo `ok` em `/health`.
- `TestCriarEventosAoMesmoTempo` — dispara 1000 `POST /eventos` em goroutines simultâneas e confere se exatamente 1000 eventos novos foram gravados.

## Estrutura

```
cmd/web/main.go        # tipos, estado, handlers, rotas e main()
cmd/web/main_test.go   # testes HTTP com httptest
templates/eventos.html # página de eventos (html/template)
```

## Decisões técnicas e trade-offs

### 1. Só a biblioteca padrão
O `go.mod` não tem nenhum `require`: roteamento, templates, validação e concorrência usam apenas `net/http`, `html/template`, `time`, `sync` e `slices`.
- **Ganho:** nenhuma dependência para instalar ou atualizar; dá para ver de perto o que a stdlib já resolve sozinha.
- **Custo:** o que um framework entregaria pronto (validação declarativa, middlewares encadeados) teria que ser escrito à mão.

### 2. `http.ServeMux` com método no padrão da rota
As rotas são registradas como `"GET /eventos"` e `"POST /eventos"` (sintaxe disponível desde o Go 1.22).
- **Ganho:** a mesma URL lista e cria, sem `if r.Method == ...` dentro do handler. Métodos não registrados recebem `405 Method Not Allowed` com o cabeçalho `Allow`, sem código extra.
- **Custo:** menos recursos que roteadores de terceiros, como grupos de rotas.

### 3. `rotas()` separada de `main()`
`rotas()` monta o mux e devolve um `http.Handler`; `main()` só chama `ListenAndServe`.
- **Ganho:** os testes chamam `rotas().ServeHTTP(rec, req)` com `httptest`, sem abrir porta nem subir servidor.
- **Custo:** a porta (`:8080`) está fixa em `main()`; trocar exige editar o código.

### 4. `html/template` em vez de `text/template`
- **Ganho:** escape automático do que vai para o HTML. O nome do evento é digitado pelo usuário, então isso evita injeção de HTML/JavaScript (XSS).
- **Custo:** praticamente nenhum para este uso.

### 5. Validação no servidor, com erro por campo
`criarEvento` aplica `strings.TrimSpace` nos campos, exige nome e data e valida a data com `time.Parse("02/01/2006", ...)`. Os erros vão num `map[string]string` indexado pelo nome do campo, e o template mostra cada mensagem ao lado do input correspondente. Os valores digitados voltam preenchidos no formulário.
- **Ganho:** o usuário não perde o que digitou. `TrimSpace` impede nome feito só de espaços. `time.Parse` rejeita datas impossíveis, como `31/02/2026`, não só o formato errado.
- **Custo:** sem validação no navegador, cada erro custa uma ida e volta ao servidor.

### 6. Post/Redirect/Get com `303 See Other`
Depois de criar o evento, o handler redireciona para `GET /eventos` em vez de renderizar a página direto.
- **Ganho:** apertar F5 depois de cadastrar não reenvia o formulário, então não cria evento duplicado.
- **Custo:** duas requisições por cadastro (o POST e o GET do redirecionamento).

### 7. Render direto no `ResponseWriter`
`w.WriteHeader(status)` é chamado antes de `tmpl.Execute(w, ...)`, sem buffer intermediário.
- **Ganho:** simples e sem alocação extra.
- **Custo:** se `Execute` falhar no meio, o status (200 ou 422) já foi enviado e a página pode chegar incompleta. O erro só é impresso no terminal com `fmt.Println`.

### 8. Teste de concorrência com 1000 goroutines
`TestCriarEventosAoMesmoTempo` manda 1000 POSTs ao mesmo tempo e confere a contagem.
- **Ganho:** prova que o mutex protege o `append`; sem ele, a contagem pode sair errada e o `-race` acusa corrida de dados.
- **Custo:** o teste depende do estado global do pacote.
