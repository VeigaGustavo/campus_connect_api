-- Migracao base para modulo feed (posts, comentarios, reacoes e salvos)
-- Execute em ambiente PostgreSQL antes de usar os endpoints /api/feed/posts...

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS feed_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    body_text TEXT NOT NULL,
    attachments JSONB NOT NULL DEFAULT '[]'::jsonb,
    share_link TEXT NOT NULL DEFAULT '',
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_feed_posts_criado_em ON feed_posts (criado_em DESC);
CREATE INDEX IF NOT EXISTS idx_feed_posts_author_id ON feed_posts (author_id);

CREATE TABLE IF NOT EXISTS feed_comentarios (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL REFERENCES feed_posts(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    body_text TEXT NOT NULL,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_feed_comentarios_post_id_criado_em ON feed_comentarios (post_id, criado_em ASC);
CREATE INDEX IF NOT EXISTS idx_feed_comentarios_author_id ON feed_comentarios (author_id);

CREATE TABLE IF NOT EXISTS feed_post_reacoes (
    post_id UUID NOT NULL REFERENCES feed_posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    reaction VARCHAR(8) NOT NULL CHECK (reaction IN ('like', 'dislike')),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_feed_post_reacoes_post_id ON feed_post_reacoes (post_id);
CREATE INDEX IF NOT EXISTS idx_feed_post_reacoes_user_id ON feed_post_reacoes (user_id);

CREATE TABLE IF NOT EXISTS feed_comentario_reacoes (
    comment_id UUID NOT NULL REFERENCES feed_comentarios(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    reaction VARCHAR(8) NOT NULL CHECK (reaction IN ('like', 'dislike')),
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    atualizado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_feed_comentario_reacoes_comment_id ON feed_comentario_reacoes (comment_id);
CREATE INDEX IF NOT EXISTS idx_feed_comentario_reacoes_user_id ON feed_comentario_reacoes (user_id);

CREATE TABLE IF NOT EXISTS feed_posts_salvos (
    post_id UUID NOT NULL REFERENCES feed_posts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    criado_em TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_feed_posts_salvos_user_id_criado_em ON feed_posts_salvos (user_id, criado_em DESC);
CREATE INDEX IF NOT EXISTS idx_feed_posts_salvos_post_id ON feed_posts_salvos (post_id);
