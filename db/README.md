# Migrações PostgreSQL

Scripts em `db/migrations/` (ordem lexicográfica). Aplicados automaticamente no arranque da API ou manualmente:

```bash
go run ./cmd/migrate
```

Requer `DATABASE_URL` (ou `DATABASE_URL_FILE`) no `.env`.

## Banco novo (Docker Swarm)

Na primeira criação do volume Postgres, pode montar `db/migrations` em `/docker-entrypoint-initdb.d/` (ver comentário em `infra/swarm/db-stack.yml`).

## Ficheiros

| Ficheiro | Conteúdo |
|----------|----------|
| `001_extensions.sql` | `pgcrypto` |
| `002_schema_base.sql` | Tabelas core (usuários, feed_cartoes, conteúdos) |
| `003_seed_perfis.sql` | Perfis `padrao`, `sistema_admin`, `empresa`, `comunidade`, `universidade` |
| `004_feed_social.sql` | Posts, comentários, reações, salvos |
| `005_perfil_comunidades.sql` | Colunas de perfil e `usuario_comunidades` |
| `006_grupos_e_feed_cartoes.sql` | Membros de grupos, pedidos de entrada, CHECK em `feed_cartoes` |

Todos os scripts são idempotentes (`IF NOT EXISTS`, `ADD COLUMN IF NOT EXISTS`).
