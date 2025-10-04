# 🚀 Guia Rápido - Todo App

## Iniciar a Aplicação

```bash
# 1. Build
go build -o todo-app ./cmd/server/

# 2. Executar
./todo-app
```

Servidor rodando em: `http://localhost:8080`

## Testar a API

### Criar uma tarefa

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Minha primeira tarefa",
    "description": "Testar a aplicação"
  }'
```

### Listar todas as tarefas

```bash
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks
```

### Atualizar uma tarefa

```bash
# Substitua {task-id} pelo ID retornado na criação
curl -X PUT http://localhost:8080/api/tasks/{task-id} \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Tarefa atualizada",
    "description": "Nova descrição",
    "status": "in_progress"
  }'
```

### Deletar uma tarefa

```bash
curl -X DELETE http://localhost:8080/api/tasks/{task-id} \
  -H "X-User-ID: user-1"
```

## Testar a Interface Web

1. Abra o navegador em: `http://localhost:8080/tasks`
2. Crie tarefas usando o formulário
3. Veja a lista atualizar automaticamente com HTMX
4. Delete tarefas com confirmação

## Executar Testes

```bash
# Todos os testes
go test ./...

# Testes com detalhes
go test -v ./...

# Testes de um pacote específico
go test -v ./internal/domain/application/
go test -v ./internal/domain/service/
go test -v ./internal/usecases/
```

## Verificar o Banco de Dados

```bash
# Instalar sqlite3 (se necessário)
# Ubuntu/Debian: sudo apt install sqlite3
# MacOS: brew install sqlite3

# Conectar ao banco
sqlite3 todo.db

# Queries úteis:
sqlite> .tables
sqlite> SELECT * FROM tasks;
sqlite> SELECT * FROM users;
sqlite> SELECT * FROM task_shares;
sqlite> .quit
```

## Status das Tarefas

- `pending` - Pendente (padrão ao criar)
- `in_progress` - Em Progresso
- `completed` - Concluída

## Exemplos Completos

### Cenário 1: CRUD Completo

```bash
# 1. Criar tarefa
TASK_ID=$(curl -s -X POST http://localhost:8080/api/tasks \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Comprar pão","description":"Padaria da esquina"}' | \
  jq -r '.ID')

echo "Tarefa criada: $TASK_ID"

# 2. Listar tarefas
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks | jq

# 3. Atualizar para "em progresso"
curl -X PUT http://localhost:8080/api/tasks/$TASK_ID \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Comprar pão","description":"Padaria da esquina","status":"in_progress"}'

# 4. Obter tarefa específica
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks/$TASK_ID | jq

# 5. Marcar como concluída
curl -X PUT http://localhost:8080/api/tasks/$TASK_ID \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Comprar pão","description":"Padaria da esquina","status":"completed"}'

# 6. Deletar
curl -X DELETE http://localhost:8080/api/tasks/$TASK_ID \
  -H "X-User-ID: user-1"
```

### Cenário 2: Múltiplos Usuários

```bash
# Usuário 1 cria tarefas
curl -X POST http://localhost:8080/api/tasks \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{"title":"Tarefa do usuário 1","description":"Privada"}'

# Usuário 2 cria tarefas
curl -X POST http://localhost:8080/api/tasks \
  -H "X-User-ID: user-2" \
  -H "Content-Type: application/json" \
  -d '{"title":"Tarefa do usuário 2","description":"Outra privada"}'

# Cada usuário vê apenas suas tarefas
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks | jq
curl -H "X-User-ID: user-2" http://localhost:8080/api/tasks | jq
```

## Troubleshooting

### Erro: "Unauthorized"
- Certifique-se de incluir o header `X-User-ID`

### Erro: "Content-Type must be application/json"
- Adicione o header: `-H "Content-Type: application/json"`

### Erro: "task not found"
- Verifique se o ID da tarefa está correto
- Confirme que você está usando o mesmo `X-User-ID` do criador

### Banco de dados corrompido
```bash
rm todo.db
./todo-app  # Recria automaticamente
```

## Arquitetura

```
┌─────────────────────────────────────────────────┐
│              HTTP Layer (Port 8080)             │
│  Handlers + Middlewares (Auth, CORS, Security) │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│              Use Cases Layer                    │
│  CreateTask, UpdateTask, DeleteTask, etc.       │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│           Domain Services Layer                 │
│    TaskService (business rules)                 │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│            Repository Layer                     │
│   TaskRepo, UserRepo, ShareRepo (interfaces)    │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│          Infrastructure Layer                   │
│     SQLite3 (prepared statements)               │
└─────────────────────────────────────────────────┘
```

## Próximos Passos

Veja o [README.md](./README.md) para funcionalidades futuras e [CLAUDE.md](./CLAUDE.md) para guia completo de desenvolvimento.
