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
- Campos vazios podem ser enviados como string vazia para limpar valor.

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

### 5) Destaque de comunidade

- Quando usuario estiver associado a comunidade institucional (atletica/CA), retornar bloco `community_highlight`.
- Caso contrario, retornar `null`.

### 6) Privacidade e autorizacao

- `GET /profile` retorna perfil do usuario autenticado.
- `PUT /profile` permite editar apenas o proprio perfil.
- `GET /profile/history` retorna apenas historico do usuario autenticado.
- Para perfil publico futuro (`GET /profiles/:id`), aplicar mascaramento de campos sensiveis (ex.: email).

---

## API (v2 implementada)

### 1) Obter perfil autenticado

- **GET** `/api/profile`

#### Response 200

```json
{
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

---

### 2) Atualizar perfil

- **PUT** `/api/profile`

#### Request

```json
{
  "about_me": "Sou estudante de SI focado em backend e dados.",
  "job_title": "Estagiario de Engenharia de Software",
  "course": "Sistemas de Informacao",
  "semester": "6",
  "institution_name": "Universidade X",
  "interests": ["Backend", "Cloud"],
  "favorite_topics": ["Arquitetura", "APIs"],
  "specialties": ["Go", "SQL"]
}
```

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

