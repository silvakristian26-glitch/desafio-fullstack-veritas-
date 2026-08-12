# Mini Kanban de Tarefas — Desafio Fullstack Veritas

Kanban simples com três colunas fixas (**A Fazer**, **Em Progresso**, **Concluídas**), construído com **React** (frontend) e **Go** (backend, API REST).

![User Flow](docs/user-flow.png)

## Stack

- **Backend:** Go 1.21+, apenas biblioteca padrão (`net/http`), sem frameworks externos.
- **Frontend:** React 18 + Vite (JavaScript puro, sem TypeScript).
- **Persistência:** arquivo JSON local (`backend/tasks.json`), gerado automaticamente. Sem banco de dados — atende ao escopo mínimo com bônus de persistência.

## Como rodar

### Pré-requisitos
- Go 1.21+
- Node.js 18+ e npm

### Backend

```bash
cd backend
go run .
```

A API sobe em `http://localhost:8080`. As tarefas são persistidas em `backend/tasks.json` (path configurável via variável de ambiente `DATA_FILE`; porta via `PORT`).

### Frontend

Em outro terminal:

```bash
cd frontend
npm install
npm run dev
```

A aplicação sobe em `http://localhost:5173`. Por padrão ela aponta para `http://localhost:8080` (ver `frontend/.env.example` — copie para `.env` caso queira apontar para outra URL de API via `VITE_API_URL`).

## Endpoints da API

| Método | Rota          | Descrição                          |
|--------|---------------|-------------------------------------|
| GET    | `/tasks`      | Lista todas as tarefas              |
| GET    | `/tasks/{id}` | Retorna uma tarefa específica       |
| POST   | `/tasks`      | Cria uma tarefa                     |
| PUT    | `/tasks/{id}` | Atualiza título/descrição/status    |
| DELETE | `/tasks/{id}` | Remove uma tarefa                   |
| GET    | `/health`     | Healthcheck simples                 |

Status válidos: `a_fazer`, `em_progresso`, `concluida`.

Exemplo de criação:

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Estudar Go","description":"Revisar goroutines","status":"a_fazer"}'
```

## Funcionalidades implementadas

- CRUD completo de tarefas (criar, editar, mover, excluir).
- Três colunas fixas com contador de tarefas por coluna.
- Mover tarefas entre colunas por **drag-and-drop** ou pelo seletor no card (dois caminhos redundantes, para acessibilidade e para o caso de o navegador do avaliador bloquear o drag nativo).
- Feedback visual de carregamento e banner de erro com opção de fechar.
- Validação de título obrigatório e status válido, tanto no frontend (evita chamada desnecessária) quanto no backend (fonte da verdade).
- Persistência em arquivo JSON — os dados sobrevivem a um restart do backend.
- CORS liberado para o frontend consumir a API.

## Decisões técnicas

- **Sem framework HTTP no backend:** o escopo é pequeno o suficiente para `net/http` puro (`ServeMux`) atender bem, sem trazer dependências externas que exigiriam acesso à internet no build/avaliação.
- **IDs gerados com `crypto/rand`** (hex de 16 bytes) em vez de uma lib de UUID, pelo mesmo motivo acima — projeto builda offline, sem `go.sum` externo.
- **Armazenamento em memória (`map[string]Task`) + `sync.RWMutex`**, com persistência em disco a cada mutação (create/update/delete). É simples, thread-safe e evita a complexidade de um banco para o escopo do desafio.
- **Duas formas de mover tarefa** (drag-and-drop e `<select>` no card): drag-and-drop é o requisito bônus, mas manter o seletor garante que a funcionalidade principal (mover entre colunas) não dependa de um recurso que pode falhar dependendo do navegador/dispositivo do avaliador.
- **Atualização otimista no `handleMove`:** o card muda de coluna na UI imediatamente e é revertido se a API falhar, para a experiência parecer fluida mesmo com uma chamada de rede no meio.
- **Sem TypeScript:** optei por JS puro no tempo disponível para o desafio, priorizando cobrir os requisitos funcionais com folga.

## Limitações conhecidas

- Sem autenticação/usuários — qualquer cliente que acesse a API pode ler/escrever tarefas.
- Persistência em arquivo único (JSON) não é adequada para concorrência pesada ou múltiplas instâncias do backend — serve para o escopo do desafio, não para produção.
- Sem testes automatizados (endpoints foram validados manualmente via `curl` durante o desenvolvimento).
- Sem Docker (mencionado como bônus, não implementado por prioridade de tempo).
- Ordenação das tarefas dentro de cada coluna segue a ordem retornada pelo mapa em memória, não a ordem de criação/arraste (não há campo de posição/`order`).

## Melhorias futuras

- Adicionar testes (Go: `net/http/httptest`; React: Vitest + Testing Library).
- Persistir em um banco leve (SQLite) em vez de JSON, mantendo a mesma interface do `Store`.
- Campo de posição (`order`) por tarefa para preservar a ordem de arraste dentro da coluna.
- Autenticação simples (token) e CORS restrito a uma origem específica.
- Dockerfile + docker-compose para subir backend e frontend com um comando só.

## Estrutura do repositório

```
/backend
  main.go       — servidor HTTP, rotas, helpers de arquivo/ID
  handlers.go   — handlers REST e Store (CRUD + persistência)
  models.go     — Task, TaskInput e validação
/frontend
  package.json
  src/
    App.jsx             — estado global, chamadas à API, loading/erro
    api.js               — client HTTP da API
    constants.js         — definição das 3 colunas
    components/
      Column.jsx
      TaskCard.jsx
      TaskModal.jsx
/docs
  user-flow.png   — diagrama do fluxo de uso
  user-flow.dot    — fonte Graphviz do diagrama (para editar)
README.md
```
