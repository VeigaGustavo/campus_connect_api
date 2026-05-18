# Contrato — Register (comunidade) e grupos

Alinhado com o front Flutter e a API Go em `feature/perfil`.

## POST `/api/auth/register`

- **Auth:** nenhuma (público)
- **Content-Type:** `application/json`
- **Sucesso:** `201 Created`

### Campos comuns (todos os tipos)

| Campo | Tipo | Obrigatório | Descrição |
|-------|------|-------------|-----------|
| `profile_type` | string | sim | `estudante` \| `comunidade` \| `empresa` \| `universidade` |
| `full_name` | string | sim | Nome completo |
| `birth_date` | string | sim* | `YYYY-MM-DD` (idade calculada no back se `age` não vier) |
| `age` | number | sim* | Alternativa a `birth_date` |
| `cpf` | string | sim | Preferir string JSON; número aceite (zeros iniciais podem perder-se) |
| `institution` | string | condicional | **Obrigatório:** `estudante`, `comunidade`. **Opcional:** `empresa`, `universidade` |
| `city` | string | sim | |
| `state` | string | sim | |
| `email` | string | sim | |
| `password` | string | sim | |

\* Pelo menos um de `birth_date` ou `age` válido.

### Quando `profile_type` = `comunidade` (Atlética / CA)

| Campo | Tipo | Obrigatório | Valores |
|-------|------|-------------|---------|
| `community_type` | string | sim | `atletica` \| `ca` |
| `community_name` | string | sim | Nome da atlética / CA |
| `group_description` | string | sim | Descrição do grupo |
| `group_visibility` | string | sim | `public` \| `private` |
| `group_title` | string | não | Se vazio, o back usa `community_name` |

**Exemplo empresa (sem instituição):**

```json
{
  "profile_type": "empresa",
  "full_name": "Gustavo Antunes",
  "birth_date": "2006-05-15",
  "cpf": "06927249150",
  "city": "Palmas",
  "state": "TO",
  "email": "contato@empresa.com",
  "password": "SenhaForte123!",
  "company_name": "Veiga.dev"
}
```

**Exemplo de request (comunidade):**

```json
{
  "profile_type": "comunidade",
  "full_name": "Maria Silva",
  "birth_date": "2000-05-15",
  "cpf": "12345678901",
  "institution": "Universidade Federal",
  "city": "São Paulo",
  "state": "SP",
  "email": "maria@email.com",
  "password": "SenhaForte123!",
  "community_type": "atletica",
  "community_name": "Atlética Engenharia",
  "group_title": "Grupo Geral da Atlética",
  "group_description": "Grupo oficial para estudos e avisos da atlética.",
  "group_visibility": "public"
}
```

**Resposta 201 (comunidade):**

```json
{
  "id": "uuid-do-usuario",
  "name": "Maria Silva",
  "email": "maria@email.com",
  "role": "comunidade",
  "profile_type": "comunidade",
  "community_id": "uuid-da-comunidade",
  "group_id": "uuid-do-grupo"
}
```

`community_id` e `group_id` **só** aparecem quando `profile_type` = `comunidade`.

**Erros:**

| Status | code | Quando |
|--------|------|--------|
| 400 | `invalid_json` | JSON inválido |
| 400 | `invalid_registration` | Validação; campo **`message`** traz explicação (ex.: `cadastro invalido: company_name...`). Chaves JSON podem ser snake_case ou camelCase equivalente (normalização no servidor) |
| 400 | `registration_failed` | Erro de BD (ex.: migração `visibility` não aplicada) |

### O que o back faz na mesma transação (comunidade)

1. Cria `usuarios` com `role` = `comunidade`
2. Grava `cadastros_usuario.details_json` (inclui `group_*`, `birth_date`, etc.)
3. Insere em `comunidades` → devolve `community_id`
4. Insere em `grupos_estudo`:
   - `titulo` ← `group_title` ou `community_name`
   - `field_of_study` ← `"Atlética"` ou `"Centro Acadêmico"`
   - `description` ← `group_description`
   - `level` = `beginner`, `member_count` = 1
   - `visibility` ← `group_visibility` (migração no arranque da API)
   - `criado_por` = id do novo usuário
5. Cria cartão no feed (`kind`: `study_group`)

---

## Depois do cadastro

### POST `/api/auth/login`

Resposta com `access_token` — necessário para rotas abaixo.

### GET `/api/groups`

- **Auth:** `Bearer` token
- **Sucesso:** `200` — array de grupos

```json
[
  {
    "id": "uuid-grupo",
    "author_id": "uuid-usuario",
    "author": {
      "id": "uuid-usuario",
      "name": "Maria Silva",
      "avatar_url": "",
      "role": "comunidade"
    },
    "title": "Grupo Geral da Atlética",
    "field_of_study": "Atlética",
    "description": "Grupo oficial...",
    "level": "beginner",
    "member_count": 1,
    "schedule_label": "",
    "visibility": "public"
  }
]
```

---

## Resumo por `profile_type`

| profile_type | Campos específicos |
|--------------|-------------------|
| `estudante` | — |
| `comunidade` | `community_type`, `community_name`, `group_description`, `group_visibility`, `group_title?` |
| `empresa` | `company_name` (+ `company_cnpj?`, `company_description?`) |
| `universidade` | `institution_name` (+ sigla, tipo, descrição opcionais) |

---

## Pendente (não implementado)

| Ação | Endpoint sugerido |
|------|---------------------|
| Entrar (público) | `POST /api/groups/:id/join` |
| Pedir entrada (privado) | `POST /api/groups/:id/join-requests` |
| Aprovar/rejeitar | `POST /api/groups/:id/join-requests/:requestId/accept` |
| Promover admin (máx. 3) | `POST /api/groups/:id/admins` |
| Listar membros | `GET /api/groups/:id/members` |

O botão “Entrar” no Flutter continua mock até estes endpoints existirem.
