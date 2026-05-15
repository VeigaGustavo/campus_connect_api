-- Executar uma vez no Postgres (coluna opcional para tipo editorial do post).
ALTER TABLE feed_posts ADD COLUMN IF NOT EXISTS content_kind TEXT NOT NULL DEFAULT '';

-- Cartões de post social no feed agregado (kind = post).
ALTER TABLE feed_cartoes DROP CONSTRAINT IF EXISTS feed_cartoes_kind_check;
ALTER TABLE feed_cartoes ADD CONSTRAINT feed_cartoes_kind_check CHECK (kind::text = ANY (ARRAY[
  'internship','event','study_group','project','notice','campus_feed','post','reading'
]::text[]));
