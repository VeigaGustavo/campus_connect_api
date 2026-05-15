# Contrato de Perfil com o Front

Este documento define as regras de negocio e o contrato de API para evolucao do modulo de perfil.

## Objetivo

O perfil deve consolidar:

- dados basicos do cadastro;
- dados complementares editaveis (sobre mim, cargo, curso, semestre e afins);
- preferencias (gostos, topicos e especialidades);
- historico do usuario (posts, leituras e grupos);
- destaque de comunidade (ex.: atletica, CA) quando aplicavel.

## Regras de negocio

### 1) Dados basicos

- Nome, email, cidade/estado e outros campos de cadastro sao a base do perfil.
- Email e identificador do usuario nao sao alterados por esta API de perfil.

### 2) Dados complementares de perfil

- Usuario pode editar:
  - `about_me`
  - `job_title`
  - `course`
  - `semester`
  - `institution_name` (quando aplicavel)
  - `avatar_image_url` e `cover_image_url` (URLs publicas apos upload no storage; ver secao "Fotos" em `PUT /api/profile`)
- Campos vazios podem ser enviados como string vazia para limpar valor (nas imagens, `""` remove a URL guardada).

### 3) Gostos, topicos e especialidades

- Usuario pode manter listas de:
  - `interests` (gostos gerais)
  - `favorite_topics` (topicos de interesse)
  - `specialties` (areas de dominio)
- Regras sugeridas:
  - sem duplicados (comparacao case-insensitive);
  - limite de 20 itens por lista;
  - tamanho maximo de 50 caracteres por item.

### 4) Historico no perfil

- Historico deve mostrar:
  - posts criados pelo usuario;
  - leituras criadas/publicadas pelo usuario;
  - grupos de estudo criados pelo usuario (criador).
- Historico e somente leitura via API de perfil (origem nos modulos de feed/leitura/grupo).
- Contas com `profile_context` `"organization"` usam o mesmo endpoint `GET /profile/history` (conteudo publicado pela **conta** / `author_id` / `criado_por` desse usuario). A UI pode decidir se exibe ou oculta para telas so de instituicao.

### 5) Destaque de comunidade

- Quando usuario estiver associado a comunidade institucional (atletica/CA), retornar bloco `community_highlight`.
- Caso contrario, retornar `null`.
- Em `profile_context` `"organization"`, `community_highlight` e sempre `null`.

### 6) Privacidade e autorizacao

- `GET /profile` retorna o perfil associado à conta autenticada (pessoa ou instituicao, conforme `profile_context` e `role` do JWT).
- `PUT /profile` permite editar apenas o proprio perfil.
- `GET /profile/history` retorna apenas historico do usuario autenticado.
- Para perfil publico futuro (`GET /profiles/:id`), aplicar mascaramento de campos sensiveis (ex.: email).

### 7) Perfil de pessoa vs perfil de instituicao

O papel da conta vem do JWT (`role`: `padrao`, `sistema_admin`, `comunidade`, `empresa`, `universidade`).

- **`profile_context`: `"user"`** — contas `padrao` e `sistema_admin`: resposta montada a partir de `usuarios` (nome pessoal, avatar, capa, contagens, `community_highlight` quando aplicavel).
- **`profile_context`: `"organization"`** — contas `comunidade`, `empresa`, `universidade`: a tela deve tratar como **perfil da instituicao**, nao da pessoa:
  - `name`, `email`, `city_state`, `about_me`, `institution_name` refletem a organizacao (ver fontes abaixo);
  - `avatar_image_url`, `cover_image_url`: perfil **usuario** vêm de `usuarios`; perfil **organizacao** vêm de `cadastros_usuario.details_json` (apos `PUT` com esses campos);
  - `job_title`, `course`, `semester` costumam vir vazios em modo organizacao;
  - `applications_count`, `groups_count`, `events_count` no objeto raiz permanecem `0` (use os totais em `organization_panel` quando existir);
  - `community_highlight` e sempre `null`;
  - **`organization_panel`**: bloco com listagens e totais do conteudo publicado por essa conta (`criado_por` / `author_id` = usuario autenticado), ate 12 itens por lista:
    - **Empresa:** vagas (`oportunidades`), eventos, posts; `about_me` = sobre a empresa.
    - **Universidade:** `map_url` (URL de mapa / embed), eventos, posts; `about_me` = sobre a instituicao.
    - **Comunidade:** `parent_institution` (instituicao informada no cadastro, campo `institution` no JSON), eventos, grupos de estudo criados por essa conta, posts; `about_me` = sobre a comunidade.

Fontes dos dados em modo organizacao:

- Base: `cadastros_usuario.details_json` (dados do `POST /api/auth/register`).
- **Comunidade:** se existir comunidade criada por esse usuario em `comunidades`, usa-se a **primeira** (por `criado_em`) para `name`, `about_me` (descricao) e `institution_name` (`kind`). Senao, usa-se `community_name` e `community_type` do cadastro.

Se nao existir linha em `cadastros_usuario`, a API devolve um perfil institucional minimo (`id`, `email` da conta) mas ainda preenche **`organization_panel`** com dados agregados do banco (posts, eventos, etc.) quando houver.

---

## API (v2 implementada)

### 1) Obter perfil autenticado

- **GET** `/api/profile`

#### Response 200 (campos comuns)

Sempre incluir:

- `profile_context`: `"user"` | `"organization"`

Exemplo **estudante / admin** (`profile_context`: `"user"`):

```json
{
  "profile_context": "user",
  "id": "uuid",
  "name": "string",
  "email": "string",
  "avatar_image_url": "string",
  "cover_image_url": "string",
  "city_state": "string",
  "about_me": "string",
  "job_title": "string",
  "course": "string",
  "semester": "string",
  "institution_name": "string",
  "interests": ["Backend", "IA"],
  "favorite_topics": ["Sistemas Distribuidos"],
  "specialties": ["Go", "PostgreSQL"],
  "applications_count": 0,
  "groups_count": 0,
  "events_count": 0,
  "community_highlight": {
    "id": "uuid",
    "name": "Atletica de Computacao",
    "kind": "athletic|academic_center|community",
    "role": "member|director|president"
  }
}
```

Exemplo **organizacao** (`profile_context`: `"organization"`): mesma forma de objeto JSON no topo; `community_highlight` = `null`. Inclui **`organization_panel`** (sempre presente em modo organizacao), por exemplo:

```json
"organization_panel": {
  "parent_institution": "so comunidade: texto do cadastro (institution)",
  "map_url": "so universidade: URL salva em cadastro (map_url)",
  "jobs": [{ "id": "uuid", "title": "...", "subtitle": "...", "reference_id": "uuid", "created_at": "RFC3339" }],
  "events": [],
  "groups": [],
  "posts": [{ "id": "uuid", "preview": "trecho do texto", "created_at": "RFC3339" }],
  "jobs_total": 0,
  "events_total": 0,
  "groups_total": 0,
  "posts_total": 0
}
```

- **Empresa:** preenche `jobs`, `events`, `posts` e totais correspondentes; `groups` omitido ou vazio.
- **Universidade:** preenche `map_url`, `events`, `posts`; `jobs` e `groups` vazios.
- **Comunidade:** preenche `parent_institution`, `events`, `groups`, `posts`; `jobs` vazio.

---

### 2) Atualizar perfil

- **PUT** `/api/profile`

#### Fotos — upload rapido (recomendado)

Use **multipart** (nao base64 no JSON). A API redimensiona e comprime no servidor e devolve **so a URL** (resposta leve, sem recarregar `organization_panel`).

| Endpoint | Campo do form | Tamanho final (aprox.) |
|----------|---------------|-------------------------|
| **POST** `/api/profile/avatar` | `file`, `image` ou **`avatar`** (Flutter) | 256×256 JPEG (~30–80 KB) |
| **POST** `/api/profile/cover` | `file`, `image` ou **`cover`** | 1280×420 JPEG (~80–200 KB) |

**Importante:** envie o **ficheiro binario** no multipart. Nao envie URL `blob:...` no JSON — isso e so preview local no browser.

Headers: `Authorization: Bearer <token>`, `Content-Type: multipart/form-data`.

Exemplo (fetch):

```js
const fd = new FormData();
fd.append("file", file); // File do input, ja comprimido no client se possivel
await fetch("/api/profile/avatar", { method: "POST", headers: { Authorization: `Bearer ${token}` }, body: fd });
// 200: { "avatar_image_url": "http://localhost:8080/uploads/avatars/avatar_<uuid>.jpg" }
```

Variaveis de ambiente (opcional):

- `UPLOAD_DIR` — pasta de gravacao (padrao `./data/uploads`).
- `PUBLIC_BASE_URL` — prefixo da URL publica (ex. `http://localhost:8080`); sem isso a API devolve caminho relativo `/uploads/...`.

**No front:** apos o POST, atualize so o estado local com a URL devolvida; **nao** chame `PUT /api/profile` nem `GET /api/profile` so por causa da foto.

Limite de entrada: **5 MB** por ficheiro. Formatos: JPEG, PNG, GIF.

#### Fotos via URL (`PUT /api/profile`) — legado

Campos opcionais `avatar_image_url` e `cover_image_url` (URL publica ja hospedada).

- **Omissao do campo**: mantém o valor atual.
- **`""`**: remove a imagem.
- **`data:image/...` (base64)**: rejeitado com erro — use os endpoints POST acima.

- **Usuario**: grava em `usuarios`.
- **Organizacao**: grava em `cadastros_usuario.details_json`.

#### Request — usuario (`padrao` / `sistema_admin`)

```json
{
  "about_me": "Sou estudante de SI focado em backend e dados.",
  "job_title": "Estagiario de Engenharia de Software",
  "course": "Sistemas de Informacao",
  "semester": "6",
  "institution_name": "Universidade X",
  "avatar_image_url": "https://cdn.exemplo.com/avatars/u1.jpg",
  "cover_image_url": "https://cdn.exemplo.com/covers/u1.jpg",
  "interests": ["Backend", "Cloud"],
  "favorite_topics": ["Arquitetura", "APIs"],
  "specialties": ["Go", "SQL"]
}
```

#### Request — mapeamento para instituicao (mesmo JSON; interpretacao por `role` do JWT)

| Campo no body (`PUT`) | Comunidade (`comunidade`) | Empresa (`empresa`) | Universidade (`universidade`) |
|----------------------|---------------------------|----------------------|----------------------------------|
| `institution_name`   | Nome da comunidade (e `community_name` no cadastro; sincroniza `comunidades` se existir) | Razao social (`company_name`) | Nome da instituicao (`institution_name`) |
| `about_me`           | Descricao (`description` em `comunidades` se existir) | Descricao (`company_description`) | Descricao (`institution_description`) |
| `course`             | Tipo / categoria (`community_type`) | (nao usado) | Sigla (`institution_acronym`) |
| `job_title`          | (nao usado) | CNPJ (`company_cnpj`) | (nao usado) |
| `semester`           | (nao usado) | (nao usado) | Tipo de instituicao (`institution_type`) |
| `map_url`            | (nao usado) | (nao usado) | URL do mapa (`map_url` no cadastro; ex. link Google Maps ou iframe permitido pelo front) |
| `avatar_image_url`   | Logotipo / foto (`details_json`) | Idem | Idem |
| `cover_image_url`    | Capa (`details_json`) | Idem | Idem |

Listas `interests`, `favorite_topics`, `specialties` em modo organizacao sao persistidas em `cadastros_usuario.details_json` (nao nas colunas de `usuarios`).

#### Response 200

- Retorna o perfil completo atualizado (mesmo formato do `GET /api/profile`).

---

### 3) Obter historico consolidado

- **GET** `/api/profile/history?limit=20`

#### Response 200

```json
{
  "items": [
    {
      "id": "uuid",
      "kind": "post|reading|group",
      "title": "string",
      "subtitle": "string",
      "reference_id": "uuid",
      "created_at": "2026-04-26T13:00:00Z"
    }
  ]
}
```

---

## Erros padrao

- `400` `invalid_json` / `validation_error`
- `401` `unauthorized`
- `403` `forbidden`
- `404` `not_found`
- `500` `server_error`

