# Contrato API — Campus Connect (Flutter / Web)

Documento de referência alinhado ao que o app Flutter consome hoje e ao que a API Go implementa em `campus_connect_api`.

## Convenções

| Item | Valor |
|------|--------|
| Base URL | `http://localhost:8080` ou `API_BASE_URL` no build |
| Formato | JSON, `snake_case` |
| Auth | `Authorization: Bearer <access_token>` nas rotas autenticadas |
| Datas | ISO 8601 (`2026-05-15T14:00:00Z`) |
| Erros | HTTP 4xx/5xx; corpo `{ "code", "message" }` |

Variáveis úteis no backend:

- `PUBLIC_BASE_URL` — prefixo para URLs de upload (`/uploads/...` → URL absoluta).
- `UPLOAD_DIR` — pasta de ficheiros (padrão `./data/uploads`).
- `CORS_ORIGIN` — origem permitida (vazio = `*`).

---

## 1. Autenticação

### POST `/api/auth/login`

**Request:** `{ "email", "password" }`

**Response 200:**

```json
{
  "access_token": "eyJ...",
  "token_type": "Bearer",
  "login": "user@edu.br",
  "expires_in": 43200,
  "role": "padrao",
  "user_id": "uuid"
}
```

**Roles:** `padrao` | `comunidade` | `empresa` | `universidade` | `sistema_admin`

### POST `/api/auth/register`

Auth: nenhuma. **Content-Type:** `application/json`. **201 Created.**

**Convenção recomendada:** chaves em `snake_case`. A API também aceita **camelCase** e chaves **case-insensitive** equivalentes (ex.: `companyName` → mesmo que `company_name`; ver nota abaixo).

**Campos base (quase todos os tipos):**

| Campo | Obrigatório | Notas |
|-------|-------------|--------|
| `profile_type` | sim | `estudante` \| `comunidade` \| `empresa` \| `universidade` (trim + minúsculas no servidor) |
| `full_name` | sim | |
| `birth_date` ou `age` | sim* | `birth_date`: `YYYY-MM-DD` ou ISO (`2006-05-15T00:00:00Z`). Idade calculada se `age` omitido |
| `cpf` | sim | Preferir **string** JSON (`"069..."`) para não perder zeros à esquerda; número também aceite |
| `city`, `state` | sim | |
| `email`, `password` | sim | |
| `institution` | condicional | **Obrigatório** para `estudante` e `comunidade`. **Opcional** para `empresa` e `universidade` (pode omitir ou `""`) |

**Por `profile_type`:**

| Tipo | Campos extra obrigatórios |
|------|---------------------------|
| `comunidade` | `community_type`, `community_name`, `group_description`, `group_visibility` (`public` \| `private`); `group_title` opcional |
| `empresa` | `company_name` |
| `universidade` | `institution_name` |
| `estudante` | — (só `institution` comum) |

**Resposta 201 (comunidade):**

```json
{
  "id": "uuid",
  "name": "Maria Silva",
  "email": "maria@email.com",
  "role": "comunidade",
  "profile_type": "comunidade",
  "community_id": "uuid-comunidade",
  "group_id": "uuid-grupo"
}
```

**Resposta 201 (outros tipos):** `id`, `name`, `email`, `role`, `profile_type` (sem `community_id` / `group_id`).

**Erros 400:**

| `code` | `message` (exemplos) |
|--------|----------------------|
| `invalid_json` | corpo não é JSON válido |
| `invalid_registration` | texto explicativo (ex.: `cadastro invalido: company_name obrigatorio para empresa`) — usar para debug/UI |
| `registration_failed` | BD (ex.: e-mail duplicado) |

**Compatibilidade de chaves:** o servidor normaliza nomes de campos (remove `_`/`-`, minúsculas), por isso `profileType` e `profile_type` equivalem a `profiletype`. **Recomendado** manter `snake_case` no contrato Flutter para consistência com o resto da API.

`birth_date` com ISO completo ou só data; `profile_type` com maiúsculas é normalizado no servidor.

### GET `/api/auth/me` (autenticado)

Retorna `{ "id", "email", "role", "name?" }`.

---

## 2. Perfil

### GET `/api/profile` (autenticado)

Inclui:

- `profile_context`: `"user"` | `"organization"`
- `profile_type`: `estudante` | `comunidade` | `empresa` | `universidade` (cadastro)
- `role`: papel da conta (`padrao`, etc.)
- `institution_name`, `course`, `semester`, tags, imagens (URLs absolutas se `PUBLIC_BASE_URL` definido)
- `organization_panel` quando `profile_context` = `"organization"`

### PUT `/api/profile`

Body: `about_me`, `job_title`, `course`, `semester`, `institution_name`, `map_url`, listas de tags.  
**Não** enviar imagens em base64 — usar POST de upload.  
**Response:** mesmo formato do GET.

### Upload de imagens (recomendado)

| Método | Path | Campo multipart |
|--------|------|-----------------|
| POST | `/api/profile/avatar` | `avatar`, `file` ou `image` |
| POST | `/api/profile/cover` | `cover`, `file` ou `image` |

**Response 200** (exemplo avatar):

```json
{
  "avatar_image_url": "http://localhost:8080/uploads/avatars/avatar_<id>.jpg",
  "avatar_url": "http://localhost:8080/uploads/avatars/avatar_<id>.jpg",
  "url": "http://localhost:8080/uploads/avatars/avatar_<id>.jpg"
}
```

### GET `/api/profile/history?limit=20`

```json
{ "items": [{ "id", "kind", "title", "subtitle", "reference_id", "created_at" }] }
```

`kind`: `post` | `reading` | `group`

---

## 3. Feed da home

### GET `/api/feed?filter=all&group_ids=id1,id2`

**filter:** `all` | `internships` | `events` | `groups` | `projects` | `readings` | `notices`

```json
{
  "items": [{
    "id", "kind", "title", "subtitle", "excerpt",
    "meta_primary", "meta_secondary", "reference_id",
    "publish_scope", "publish_group_id"
  }]
}
```

---

## 4. Oportunidades

Contrato completo (campos, erros, diferenças Flutter): **[contrato_oportunidades.md](./contrato_oportunidades.md)**.

| Método | Rota | Auth |
|--------|------|------|
| GET | `/api/opportunities` | Não — array na raiz |
| GET | `/api/opportunities/{id}` | Não |
| POST | `/api/opportunities` | `empresa` \| `sistema_admin` |
| PUT | `/api/opportunities/{id}` | Dono ou admin |
| DELETE | `/api/opportunities/{id}` | Dono ou admin |
| GET | `/api/opportunities/{id}/applicants` | empresa/admin (mock hoje) |

**Pendente:** `POST /api/opportunities/{id}/apply`, `GET /api/opportunities/types`, candidatos reais em BD.

---

## 5. Eventos

- `GET /api/events` — array
- CRUD: `POST/PUT/DELETE /api/events`, `/api/events/{id}`

---

## 6. Grupos

### GET `/api/groups` (autenticado)

Array de grupos; inclui `visibility` (`public` | `private`):

```json
[{
  "id": "uuid",
  "author_id": "uuid",
  "author": { "id", "name", "avatar_url", "role" },
  "title": "Grupo Geral da Atlética",
  "field_of_study": "Atlética",
  "description": "...",
  "level": "beginner",
  "member_count": 1,
  "schedule_label": "",
  "visibility": "public"
}]
```

- CRUD + chat, ficheiros, reuniões, associações

**Pendente:** `POST /api/groups/{id}/join`, pedidos de entrada em grupos privados, admins.

---

## 7. Leituras

### GET `/api/reading/weekly`

Query opcional: `kind` = `campus_news` | `magazine` | `article` (aliases `revista` → `magazine`, `artigo` → `article` na validação do POST).

```json
{ "items": [{ "id", "kind", "title", "source", "excerpt", "image_url", "meta_label", "author_id", "author" }] }
```

**Revista:** no POST use `"kind": "magazine"` (o back também aceita `"revista"` e normaliza).

### POST `/api/reading/weekly`

Auth: `universidade` | `comunidade` | `empresa` | `sistema_admin`. Corpo: `kind`, `title`, `source`, `excerpt` obrigatórios; `image_url`, `meta_label` opcionais.

---

## 8. Posts sociais

### GET `/api/feed/posts?page=1&limit=20&group_ids=&author_id=&include_comments=false`

Lista paginada com **post completo** (mesmo formato do detalhe). Auth obrigatória.

| Query | Default | Descrição |
|-------|---------|-----------|
| `page` | 1 | Página (1-based) |
| `limit` | 20 (máx. 50) | Itens por página |
| `group_ids` | — | CSV de grupos; inclui posts `publish_scope=group` visíveis |
| `author_id` | — | Filtrar posts de um autor |
| `include_comments` | false | `true` para trazer comentários em cada item |

```json
{
  "items": [{ "...": "igual GET /api/feed/posts/{id}" }],
  "total": 42,
  "page": 1,
  "limit": 20,
  "has_more": true
}
```

- `GET /api/feed/posts/{id}`
- `GET/POST /api/feed/posts/{id}/comments`
- `PUT /api/feed/posts/{id}/reaction` — `{ "reaction": "like" | "dislike" }`
- `PUT /api/feed/posts/{id}/save` — `{ "saved": true }`

---

## Checklist perfil (Flutter)

- [x] GET devolve `institution_name`, `course`, `semester`
- [x] PUT persiste e devolve perfil atualizado
- [x] `interests`, `favorite_topics`, `specialties` como `[]string`
- [x] `profile_type` + `role` no GET
- [x] URLs de mídia absolutas com `PUBLIC_BASE_URL`
- [x] POST multipart avatar/cover com resposta compatível

Detalhe adicional do módulo perfil: `docs/contrato_perfil_front.md`.
