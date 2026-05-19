CREATE TABLE IF NOT EXISTS perfis_usuario (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS usuarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome TEXT NOT NULL DEFAULT '',
    email TEXT NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    perfil_id UUID NOT NULL REFERENCES perfis_usuario (id),
    ativo BOOLEAN NOT NULL DEFAULT true,
    city_state TEXT NOT NULL DEFAULT '',
    avatar_image_url TEXT NOT NULL DEFAULT '',
    cover_image_url TEXT NOT NULL DEFAULT '',
    about_me TEXT NOT NULL DEFAULT '',
    job_title TEXT NOT NULL DEFAULT '',
    course TEXT NOT NULL DEFAULT '',
    semester TEXT NOT NULL DEFAULT '',
    institution_name TEXT NOT NULL DEFAULT '',
    course_and_semester TEXT NOT NULL DEFAULT '',
    interests JSONB NOT NULL DEFAULT '[]'::jsonb,
    favorite_topics JSONB NOT NULL DEFAULT '[]'::jsonb,
    specialties JSONB NOT NULL DEFAULT '[]'::jsonb,
    applications_count INT NOT NULL DEFAULT 0,
    groups_count INT NOT NULL DEFAULT 0,
    events_count INT NOT NULL DEFAULT 0,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_usuarios_email ON usuarios (lower(trim(email)));

CREATE TABLE IF NOT EXISTS cadastros_usuario (
    usuario_id UUID PRIMARY KEY REFERENCES usuarios (id) ON DELETE CASCADE,
    profile_type TEXT NOT NULL,
    details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS comunidades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome TEXT NOT NULL,
    kind TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_comunidades_criado_por ON comunidades (criado_por);

CREATE TABLE IF NOT EXISTS grupos_estudo (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo TEXT NOT NULL,
    field_of_study TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    level VARCHAR(32) NOT NULL DEFAULT 'beginner',
    member_count INT NOT NULL DEFAULT 0,
    schedule_label TEXT NOT NULL DEFAULT '',
    visibility TEXT NOT NULL DEFAULT 'public',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_grupos_estudo_criado_por ON grupos_estudo (criado_por);
CREATE INDEX IF NOT EXISTS idx_grupos_estudo_criado_em ON grupos_estudo (criado_em DESC);

CREATE TABLE IF NOT EXISTS feed_cartoes (
    id TEXT PRIMARY KEY,
    kind TEXT NOT NULL,
    titulo TEXT NOT NULL DEFAULT '',
    subtitle TEXT NOT NULL DEFAULT '',
    excerpt TEXT NOT NULL DEFAULT '',
    meta_primary TEXT NOT NULL DEFAULT '',
    meta_secondary TEXT NOT NULL DEFAULT '',
    reference_id TEXT NOT NULL DEFAULT '',
    visibility_scope TEXT NOT NULL DEFAULT 'all',
    visibility_group_id UUID REFERENCES grupos_estudo (id) ON DELETE SET NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_feed_cartoes_kind ON feed_cartoes (kind);
CREATE INDEX IF NOT EXISTS idx_feed_cartoes_criado_em ON feed_cartoes (criado_em DESC);
CREATE INDEX IF NOT EXISTS idx_feed_cartoes_reference ON feed_cartoes (kind, reference_id);

CREATE TABLE IF NOT EXISTS eventos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    start_at TIMESTAMPTZ NOT NULL,
    location TEXT NOT NULL DEFAULT '',
    organizer TEXT NOT NULL DEFAULT '',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_eventos_criado_em ON eventos (criado_em DESC);

CREATE TABLE IF NOT EXISTS projetos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_projetos_criado_em ON projetos (criado_em DESC);

CREATE TABLE IF NOT EXISTS oportunidades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo TEXT NOT NULL,
    company_name TEXT NOT NULL DEFAULT '',
    short_description TEXT NOT NULL DEFAULT '',
    full_description TEXT NOT NULL DEFAULT '',
    apply_deadline TIMESTAMPTZ NOT NULL,
    work_location VARCHAR(16) NOT NULL DEFAULT 'remote',
    type_label TEXT NOT NULL DEFAULT '',
    requirements JSONB NOT NULL DEFAULT '[]'::jsonb,
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_oportunidades_criado_em ON oportunidades (criado_em DESC);

CREATE TABLE IF NOT EXISTS leitura_semanal (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind VARCHAR(32) NOT NULL DEFAULT 'article',
    titulo TEXT NOT NULL,
    source TEXT NOT NULL DEFAULT '',
    excerpt TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    meta_label TEXT NOT NULL DEFAULT '',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_leitura_semanal_criado_em ON leitura_semanal (criado_em DESC);

CREATE TABLE IF NOT EXISTS avisos_universidade (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    titulo TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    criado_por UUID NOT NULL REFERENCES usuarios (id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_avisos_universidade_criado_em ON avisos_universidade (criado_em DESC);
