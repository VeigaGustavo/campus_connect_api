package banco

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AplicarMigracoesEssenciais(contexto context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(contexto, `ALTER TABLE feed_posts ADD COLUMN IF NOT EXISTS content_kind TEXT NOT NULL DEFAULT ''`); err != nil {
		return err
	}
	if _, err := pool.Exec(contexto, `
ALTER TABLE grupos_estudo ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'public';
UPDATE grupos_estudo SET visibility = 'public' WHERE visibility IS NULL OR trim(visibility) = '';
`); err != nil {
		return err
	}
	if _, err := pool.Exec(contexto, `
CREATE TABLE IF NOT EXISTS grupos_membros (
  group_id UUID NOT NULL REFERENCES grupos_estudo(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
  role VARCHAR(20) NOT NULL DEFAULT 'member',
  joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (group_id, user_id)
);
CREATE TABLE IF NOT EXISTS grupos_pedidos_entrada (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_id UUID NOT NULL REFERENCES grupos_estudo(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (group_id, user_id)
);`); err != nil {
		return err
	}
	if _, err := pool.Exec(contexto, `
ALTER TABLE feed_cartoes DROP CONSTRAINT IF EXISTS feed_cartoes_kind_check;
ALTER TABLE feed_cartoes ADD CONSTRAINT feed_cartoes_kind_check CHECK (kind::text = ANY (ARRAY[
  'internship','event','study_group','project','notice','campus_feed','post','reading'
]::text[]));`); err != nil {
		return err
	}
	_, _ = pool.Exec(contexto, `
DO $$
BEGIN
	INSERT INTO perfis_usuario (codigo) SELECT 'empresa'::text
		WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'empresa');
	INSERT INTO perfis_usuario (codigo) SELECT 'comunidade'::text
		WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'comunidade');
	INSERT INTO perfis_usuario (codigo) SELECT 'universidade'::text
		WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'universidade');
EXCEPTION WHEN OTHERS THEN
	RAISE NOTICE 'seed perfis_usuario (empresa/comunidade/universidade) ignorado: %', SQLERRM;
END $$;`)
	return nil
}
