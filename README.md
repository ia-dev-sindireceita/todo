# Todo App - Gerenciador de Tarefas

Aplicação de gerenciamento de tarefas com compartilhamento, seguindo arquitetura hexagonal, TDD e práticas de segurança.

## 🏗️ Arquitetura

- **Hexagonal Architecture** (Ports and Adapters)
- **Test-Driven Development** (TDD)
- **Domain-Driven Design** (DDD)
- **Frontend**: HTMX + Tailwind CSS (design minimalista)
- **Backend**: Go 1.24.5
- **Database**: SQLite3

## 📁 Estrutura do Projeto

```
internal/
├── domain/
│   ├── application/    # Entities e Value Objects com validações
│   ├── repository/     # Interfaces de repositórios (ports)
│   └── service/        # Regras de negócio gerais
├── usecases/          # Casos de uso específicos
└── infrastructure/
    ├── database/      # Implementações SQLite com prepared statements
    ├── http/          # Handlers e middlewares HTTP
    └── templates/     # Templates HTML com HTMX
```

## 🔒 Práticas de Segurança Implementadas

- ✅ **Prepared Statements**: Todas as queries SQL usam prepared statements (proteção contra SQL injection)
- ✅ **Validação em Entities**: Todas as validações acontecem na camada de domínio
- ✅ **Security Headers**: X-Content-Type-Options, X-Frame-Options, CSP, HSTS
- ✅ **Input Sanitization**: Validação de tipos, tamanhos e formatos
- ✅ **Error Handling**: Erros genéricos para o cliente, detalhes apenas em logs

## 🚀 Como Executar

### 1. Build

```bash
go build -o todo-app ./cmd/server/
```

### 2. Executar

```bash
./todo-app
```

O servidor iniciará em `http://localhost:8080`

## 🧪 Testes

Executar todos os testes:

```bash
go test ./...
```

Executar testes com verbose:

```bash
go test -v ./...
```

## 📡 API REST

### Autenticação

Para testar a API, inclua o header `X-User-ID`:

```bash
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks
```

### Endpoints

#### Criar Tarefa
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Comprar mantimentos",
    "description": "Leite, pão, ovos"
  }'
```

#### Listar Tarefas
```bash
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks
```

#### Listar Tarefas Compartilhadas
```bash
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks/shared
```

#### Obter Tarefa
```bash
curl -H "X-User-ID: user-1" http://localhost:8080/api/tasks/{id}
```

#### Atualizar Tarefa
```bash
curl -X PUT http://localhost:8080/api/tasks/{id} \
  -H "X-User-ID: user-1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Título atualizado",
    "description": "Nova descrição",
    "status": "in_progress"
  }'
```

#### Deletar Tarefa
```bash
curl -X DELETE http://localhost:8080/api/tasks/{id} \
  -H "X-User-ID: user-1"
```

## 🎨 Frontend (HTMX + Tailwind)

Acesse `http://localhost:8080/tasks` no navegador para usar a interface web.

Recursos:
- Criar tarefas sem JavaScript
- Listar tarefas em tempo real
- Deletar tarefas com confirmação
- Design minimalista com Tailwind CSS
- Progressive enhancement (funciona sem JS)

## 🗄️ Banco de Dados

O arquivo SQLite `todo.db` é criado automaticamente na primeira execução.

### Schema

```sql
-- Usuários
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME NOT NULL
);

-- Tarefas
CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    owner_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

-- Compartilhamentos
CREATE TABLE task_shares (
    task_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    PRIMARY KEY (task_id, user_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

## 📝 Status das Tasks

- `pending` - Pendente
- `in_progress` - Em Progresso
- `completed` - Concluída

## 🔮 Próximas Implementações

- [ ] Sistema completo de autenticação (JWT/Sessions)
- [ ] Compartilhamento de tarefas via interface web
- [ ] Filtros e busca de tarefas
- [ ] Edição inline com HTMX
- [ ] Drag & drop para alterar status
- [ ] Notificações em tempo real
- [ ] Export para CSV/JSON
- [ ] Dark mode

## 📚 Referências

- [CLAUDE.md](./CLAUDE.md) - Guia completo de desenvolvimento
- Go 1.24.5
- HTMX 1.9.10
- Tailwind CSS 3.x
- SQLite3

## 🛡️ Segurança

Este projeto implementa as práticas definidas no [CLAUDE.md](./CLAUDE.md):
- Defense in Depth
- Fail Securely
- Least Privilege
- Zero Trust
- Security by Default

**Nota**: O sistema de autenticação atual (X-User-ID header) é apenas para demonstração.
Em produção, use JWT, OAuth ou sessões seguras com bcrypt/argon2.
