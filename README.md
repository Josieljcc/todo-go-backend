# Todo Go Backend API

API RESTful para gerenciamento de tarefas desenvolvida em Go, com autenticação JWT e sistema de atribuição de tarefas entre usuários.

## Funcionalidades

- ✅ Autenticação JWT (registro e login)
- ✅ CRUD completo de tarefas
- ✅ Tarefas por usuário (cada usuário vê apenas suas tarefas)
- ✅ Tipos de tarefas: casa, trabalho, lazer, saúde
- ✅ Tempo de realização (due_date) para cada tarefa
- ✅ Usuários podem criar tarefas para outros usuários
- ✅ Filtros por tipo e status de conclusão

## Tecnologias

- **Go 1.21+**
- **Gin** - Framework web
- **GORM** - ORM para banco de dados
- **SQLite** - Banco de dados (pode ser facilmente trocado por PostgreSQL/MySQL)
- **JWT-Go** - Autenticação JWT
- **Bcrypt** - Hash de senhas

## Estrutura do Projeto

O projeto segue as melhores práticas de arquitetura em Go, utilizando separação de responsabilidades:

```
todo-go-backend/
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── config/                  # Configurações da aplicação
│   ├── database/                # Conexão e setup do banco de dados
│   ├── errors/                  # Erros customizados da aplicação
│   ├── handlers/                # Handlers HTTP (camada de apresentação)
│   ├── middleware/              # Middlewares (autenticação, etc)
│   ├── models/                  # Modelos de dados (entidades)
│   ├── repositories/            # Camada de acesso a dados (Repository Pattern)
│   └── services/                # Camada de lógica de negócio (Service Layer)
├── pkg/
│   └── utils/                   # Utilitários (JWT, password hashing)
└── go.mod
```

### Arquitetura

O projeto utiliza uma arquitetura em camadas:

1. **Handlers**: Recebem requisições HTTP, validam entrada e chamam services
2. **Services**: Contêm a lógica de negócio da aplicação
3. **Repositories**: Abstraem o acesso aos dados, permitindo fácil troca de banco de dados
4. **Models**: Definem as entidades do domínio
5. **Errors**: Erros customizados com códigos HTTP apropriados

## Instalação

1. Clone o repositório:
```bash
git clone <repository-url>
cd todo-go-backend
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure as variáveis de ambiente (opcional):
```bash
# Crie um arquivo .env ou exporte as variáveis
export PORT=8080
export JWT_SECRET=your-secret-key-change-in-production
export DATABASE_PATH=todo.db
```

4. Execute a aplicação:
```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

## Endpoints

### Autenticação

#### Registrar usuário
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "usuario",
  "email": "usuario@example.com",
  "password": "senha123"
}
```

**Resposta:**
```json
{
  "message": "User created successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "usuario",
    "email": "usuario@example.com"
  }
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "usuario",
  "password": "senha123"
}
```

**Resposta:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "usuario",
    "email": "usuario@example.com"
  }
}
```

### Tarefas (Requer autenticação)

Todas as rotas de tarefas requerem o header:
```
Authorization: Bearer <token>
```

#### Criar tarefa
```http
POST /api/v1/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Limpar a casa",
  "description": "Limpar todos os cômodos",
  "type": "casa",
  "due_date": "2024-12-31T23:59:59Z",
  "user_id": 2  // Opcional: criar tarefa para outro usuário
}
```

**Tipos válidos:** `casa`, `trabalho`, `lazer`, `saude`

#### Listar tarefas
```http
GET /api/v1/tasks?type=casa&completed=false
Authorization: Bearer <token>
```

**Query parameters opcionais:**
- `type`: Filtrar por tipo (casa, trabalho, lazer, saude)
- `completed`: Filtrar por status (true/false)

#### Obter tarefa específica
```http
GET /api/v1/tasks/:id
Authorization: Bearer <token>
```

#### Atualizar tarefa
```http
PUT /api/v1/tasks/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Título atualizado",
  "completed": true,
  "due_date": "2024-12-31T23:59:59Z"
}
```

#### Deletar tarefa
```http
DELETE /api/v1/tasks/:id
Authorization: Bearer <token>
```

## Exemplos de Uso

### Criar tarefa para outro usuário

```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Authorization: Bearer <seu-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Revisar código",
    "description": "Revisar PR #123",
    "type": "trabalho",
    "due_date": "2024-12-25T18:00:00Z",
    "user_id": 2
  }'
```

### Listar tarefas pendentes de casa

```bash
curl -X GET "http://localhost:8080/api/v1/tasks?type=casa&completed=false" \
  -H "Authorization: Bearer <seu-token>"
```

## Testes

Execute os testes com:

```bash
go test ./...
```

Para ver cobertura:

```bash
go test -cover ./...
```

**Nota**: Os testes que utilizam SQLite requerem CGO habilitado. Se você encontrar erros relacionados ao CGO, certifique-se de que o CGO está habilitado no seu ambiente Go.

Para executar testes específicos:

```bash
go test ./internal/services/... -v
go test ./internal/handlers/... -v
```

## Variáveis de Ambiente

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `PORT` | Porta do servidor | `8080` |
| `JWT_SECRET` | Chave secreta para JWT | `your-secret-key-change-in-production` |
| `DATABASE_PATH` | Caminho do arquivo SQLite | `todo.db` |
| `DATABASE_HOST` | Host do MySQL (se usando MySQL) | - |
| `DATABASE_PORT` | Porta do MySQL | `3306` |
| `DATABASE_USER` | Usuário do MySQL | - |
| `DATABASE_PASSWORD` | Senha do MySQL | - |
| `DATABASE_NAME` | Nome do banco de dados MySQL | - |

## Swagger Documentation

The API is fully documented with Swagger/OpenAPI. After starting the server, you can access the interactive documentation at:

```
http://localhost:8080/swagger/index.html
```

To regenerate the Swagger documentation after making changes:

```bash
swag init -g cmd/api/main.go
```

## Docker

### Using Docker Compose

The easiest way to run the application is using Docker Compose, which will start both the API and MySQL database:

```bash
# Build and start containers
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop containers
docker-compose down

# Stop and remove volumes (clean database)
docker-compose down -v
```

The API will be available at `http://localhost:8080` and MySQL at `localhost:3306`.

### Environment Variables for Docker

You can create a `.env` file in the root directory with the following variables:

```env
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
DATABASE_HOST=mysql
DATABASE_PORT=3306
DATABASE_USER=todo_user
DATABASE_PASSWORD=todo_password
DATABASE_NAME=todo_db
MYSQL_ROOT_PASSWORD=root_password
```

### Building Docker Image Manually

```bash
# Build the image
docker build -t todo-api .

# Run the container (requires MySQL to be running)
docker run -p 8080:8080 \
  -e DATABASE_HOST=mysql \
  -e DATABASE_USER=todo_user \
  -e DATABASE_PASSWORD=todo_password \
  -e DATABASE_NAME=todo_db \
  -e JWT_SECRET=your-secret-key \
  todo-api
```

## Roadmap

Para ver o plano completo de melhorias futuras, consulte o arquivo [ROADMAP.md](./ROADMAP.md).

### Próximas Melhorias Prioritárias

- [ ] Paginação nas listagens
- [ ] Busca por texto nas tarefas
- [ ] Notificações de tarefas próximas do vencimento
- [x] Suporte a múltiplos bancos de dados (SQLite, MySQL)
- [x] Documentação Swagger/OpenAPI
- [x] Docker e Docker Compose
- [ ] Rate limiting
- [ ] Logs estruturados
- [ ] Cache com Redis
- [ ] Refresh tokens
- [ ] CI/CD Pipeline

## Licença

MIT

