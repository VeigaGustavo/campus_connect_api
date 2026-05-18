# Contrato — Tela do grupo (`/groups/:id`)

Alinhado ao Flutter e à API Go. Ver também `docs/contrato_register_grupos.md` e `docs/CONTRATO_API.md`.

## Convenções

| Item | Valor |
|------|--------|
| Base | `http://localhost:8080` |
| JSON | `snake_case` |
| Auth | `Authorization: Bearer <jwt>` ou `?access_token=` (WebSocket) |
| Listas | sempre `[]`, nunca `null` |
| Datas | ISO 8601 UTC |

Perfis com acesso: `padrao`, `comunidade`, `sistema_admin`.

## Checklist API

| Item | Status |
|------|--------|
| `GET /api/groups` auth + `[]` | implementado |
| `visibility` em grupos | implementado |
| `POST /api/groups` com `visibility` | implementado |
| Chat GET/POST | implementado (memória; reinício limpa) |
| WS `?access_token=` | implementado |
| WS broadcast após POST mensagem | implementado |
| `GET /api/feed/posts?group_ids=` | implementado (só posts do grupo) |
| `content_kind` na query do feed | implementado |
| `GET /api/groups/:id/events` | implementado |
| `POST /api/groups/:id/join` | implementado (público) |
| `POST /api/groups/:id/join-requests` | implementado (privado) |
| `GET /api/groups/:id/members` | implementado |
| Aprovar pedidos / promover admin | pendente |

## Endpoints principais

### `GET /api/groups` · `POST /api/groups`

Ver exemplos no contrato Flutter; resposta inclui `visibility` e `author.avatar_url`.

### Chat

- `GET /api/groups/:id/chat/messages` → array ordenado por `created_at` ASC
- `POST /api/groups/:id/chat/messages` → `{ "text": "..." }` + broadcast WS

Mensagem:

```json
{
  "id": "msg-uuid",
  "group_id": "group-uuid",
  "author_id": "user-uuid",
  "author_name": "Maria Silva",
  "text": "Olá",
  "created_at": "2026-05-15T19:41:00Z"
}
```
- `GET /api/groups/:id/chat/ws?access_token=<jwt>`

Evento WS:

```json
{
  "event": "chat_message",
  "message": { "id", "group_id", "author_id", "text", "created_at" }
}
```

### Publicações do grupo

`GET /api/feed/posts?group_ids=<uuid>&page=1&limit=50&content_kind=article` (opcional)

Posts com `publish_scope: group` e `publish_group_id` = id do grupo.  
`author.avatar_url`, `my_reaction: "none"` quando sem reação.

### Entrada

- `POST /api/groups/:id/join` → **200** (grupo público)
- `POST /api/groups/:id/join-requests` → **201** `{ "status": "pending" }` (privado)

### Membros

`GET /api/groups/:id/members` → `[{ "user_id", "name", "avatar_url", "role", "joined_at?" }]`
