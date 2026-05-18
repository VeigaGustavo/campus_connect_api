ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS about_me TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS job_title TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS course TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS semester TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS institution_name TEXT NOT NULL DEFAULT '';
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS favorite_topics JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS specialties JSONB NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now();

CREATE TABLE IF NOT EXISTS usuario_comunidades (
    usuario_id UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    comunidade_id UUID NOT NULL REFERENCES comunidades (id) ON DELETE CASCADE,
    role TEXT NOT NULL DEFAULT 'member',
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (usuario_id, comunidade_id)
);

CREATE INDEX IF NOT EXISTS idx_usuario_comunidades_usuario ON usuario_comunidades (usuario_id, criado_em DESC);
