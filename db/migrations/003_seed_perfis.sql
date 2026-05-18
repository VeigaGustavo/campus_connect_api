INSERT INTO perfis_usuario (codigo)
SELECT 'padrao'::text
WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'padrao');

INSERT INTO perfis_usuario (codigo)
SELECT 'sistema_admin'::text
WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'sistema_admin');

INSERT INTO perfis_usuario (codigo)
SELECT 'empresa'::text
WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'empresa');

INSERT INTO perfis_usuario (codigo)
SELECT 'comunidade'::text
WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'comunidade');

INSERT INTO perfis_usuario (codigo)
SELECT 'universidade'::text
WHERE NOT EXISTS (SELECT 1 FROM perfis_usuario p WHERE p.codigo = 'universidade');
