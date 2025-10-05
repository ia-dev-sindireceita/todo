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
- ✅ **Rate Limiting**: Proteção contra ataques DoS e brute-force com limites configuráveis

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

### 3. Configuração (Opcional)

Variáveis de ambiente disponíveis:

```bash
# Rate limiting (padrões configurados para segurança)
export RATE_LIMIT_GENERAL=100    # Requisições por minuto para rotas normais
export RATE_LIMIT_AUTH=5          # Requisições por minuto para rotas de autenticação
export RATE_LIMIT_WINDOW=60       # Janela de tempo em segundos

# Trusted Proxies (Segurança contra IP Spoofing)
# Lista de IPs de proxies/load balancers confiáveis separados por vírgula
# Se não configurado, apenas RemoteAddr é usado (mais seguro)
# Exemplo: export TRUSTED_PROXIES="127.0.0.1,10.0.0.1"
export TRUSTED_PROXIES=""

# JWT Secret (OBRIGATÓRIO em produção)
export JWT_SECRET="your-secret-key-here"

# Executar
./todo-app
```

#### ⚠️ Importante: Configuração de Proxies Confiáveis

O rate limiting usa o endereço IP do cliente para limitar requisições. Por padrão, **apenas o IP real da conexão (`RemoteAddr`) é usado**, ignorando headers HTTP como `X-Forwarded-For` e `X-Real-IP`.

**Quando usar `TRUSTED_PROXIES`:**
- ✅ Aplicação está atrás de proxy reverso (nginx, Apache)
- ✅ Aplicação está atrás de load balancer (AWS ELB, GCP Load Balancer)
- ✅ Aplicação está atrás de CDN (Cloudflare, CloudFront)

**Como configurar:**
1. Identifique os IPs dos seus proxies/load balancers
2. Configure `TRUSTED_PROXIES` com esses IPs separados por vírgula
3. Apenas requisições vindas desses IPs poderão definir o IP do cliente via headers

**Exemplo de configuração:**
```bash
# Nginx/Apache local
export TRUSTED_PROXIES="127.0.0.1"

# Load balancer interno
export TRUSTED_PROXIES="10.0.1.10,10.0.1.11"

# Múltiplos proxies
export TRUSTED_PROXIES="127.0.0.1,10.0.1.10,172.16.0.5"
```

**⚠️ Segurança:**
- **NÃO** configure `TRUSTED_PROXIES` se não estiver usando proxy
- **NÃO** adicione IPs que você não controla
- Headers de proxy podem ser facilmente forjados se não validados corretamente

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

### Rate Limiting

Todas as rotas possuem rate limiting:

- **Rotas normais**: 100 requisições/minuto por IP
- **Rotas de autenticação** (`/api/auth/*`, `/web/auth/*`): 5 requisições/minuto por IP

Headers de resposta:
- `X-RateLimit-Limit`: Limite total de requisições
- `X-RateLimit-Remaining`: Requisições restantes
- `X-RateLimit-Reset`: Timestamp Unix quando o limite será resetado
- `Retry-After`: Segundos até poder tentar novamente (apenas em 429)

Quando o limite é excedido, retorna HTTP 429 (Too Many Requests).

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
