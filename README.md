# Campus Connect API

API REST em **Go** para substituir os mocks do **CampusConnect**: autenticação institucional (e-mail + senha, com espaço para SSO/OIDC), JWT em rotas protegidas e recursos alinhados ao app (feed Descobrir, oportunidades, eventos, grupos de estudo, perfil e candidaturas).

Este repositório está no estágio inicial: contratos JSON (`snake_case` nos campos expostos), tipos de domínio e rotas expostas como stubs — persistência, auth real e paginação vêm nos próximos passos.

Código em português (pacotes e identificadores): `principal.go` (entrada), `internal/modelos`, `internal/respostas`, `internal/manipuladores`.

## Executar localmente

```bash
go run .
```

## URL base

Em desenvolvimento o servidor escuta em **`http://localhost:8080`** (host local, porta **8080**), a menos que você defina variáveis de ambiente (ver abaixo). O app cliente deve usar essa origem como **base URL** e concatenar os paths REST.

Exemplos:

| Método | Path |
|--------|------|
| `GET` | `/health` |
| `GET` | `/discover` — query `filter`: `all`, `internships`, `events`, `groups`, `projects` |
| `GET` | `/opportunities?q=...` |
| `GET` | `/opportunities/{id}` |
| `POST` | `/opportunities/{id}/applications` |
| `GET` | `/events` |
| `GET` | `/events/{id}` |
| `GET` | `/groups` |
| `GET` | `/groups/{id}` |
| `POST` | `/groups/{id}/join` |
| `GET` | `/me` ou `/users/me` |

### Ambiente

- **`LISTEN_ADDR`**: endereço completo de bind (ex. `:8080`, `127.0.0.1:3000`). Se definido, sobrescreve `PORT`.
- **`PORT`**: apenas a porta (ex. `8080`); o servidor usa `:{PORT}` em todas as interfaces.

## Convenções

- Cabeçalhos: `Authorization: Bearer <token>`, `Content-Type: application/json`.
- Erros: `4xx` / `5xx` com corpo `{"code":"...","message":"..."}` (chaves JSON fixas em inglês por compatibilidade com clientes; **valores** de `code` e `message` estão em português neste servidor).

## Roadmap resumido

Auth, `/discover`, `/opportunities`, `/events`, `/groups`, `/me`, `POST /opportunities/{id}/applications`, tipos enumerados estáveis no código (`CategoriaItemDescobrir`, `ModalidadeLocalTrabalho`, `NivelGrupoEstudo`, `TipoAtividadePerfil`; valores JSON iguais ao app) e, opcionalmente, projetos no feed.
