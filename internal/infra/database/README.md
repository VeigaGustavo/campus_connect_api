# Convencao da camada `infra/database`

- Um arquivo por agregado/modulo (`postgres_evento.go`, `postgres_grupo.go`, etc.).
- Helpers compartilhados ficam em arquivo proprio (`postgres_feed.go`, `postgres_autorizacao.go`).
- Evitar arquivos "demais", "misc" ou multi-dominio.
- `service` contem regra de negocio; `repository/infra` contem persistencia.
- Metodos de admin e dono continuam no mesmo agregado, nao em arquivo global.

Essa organizacao facilita extracao gradual de modulos para microservicos.
