-- Adobe2API Banana public-name normalization.
-- Clients use nano-banana* names; adobe-firefly-* remains internal only.
-- Run on contabo:
--   docker exec -i newapi-postgres psql -v ON_ERROR_STOP=1 -U root -d new-api \
--     < scripts/migrate_adobe2api_banana_public_names_ssh.sql

BEGIN;

INSERT INTO model_public_aliases (
    internal_name, public_name, created_time, updated_time, deleted_at
)
SELECT
    v.internal_name,
    v.public_name,
    EXTRACT(EPOCH FROM NOW())::BIGINT,
    EXTRACT(EPOCH FROM NOW())::BIGINT,
    NULL
FROM (VALUES
    ('adobe-firefly-nano-banana-pro-1k', 'nano-banana-pro-1k'),
    ('adobe-firefly-nano-banana-pro-2k', 'nano-banana-pro-2k'),
    ('adobe-firefly-nano-banana-pro-4k', 'nano-banana-pro-4k'),
    ('adobe-firefly-nano-banana-1k', 'nano-banana-1k'),
    ('adobe-firefly-nano-banana-2k', 'nano-banana-2k'),
    ('adobe-firefly-nano-banana-4k', 'nano-banana-4k'),
    ('adobe-firefly-nano-banana2-1k', 'nano-banana2-1k'),
    ('adobe-firefly-nano-banana2-2k', 'nano-banana2-2k'),
    ('adobe-firefly-nano-banana2-4k', 'nano-banana2-4k')
) AS v(internal_name, public_name)
ON CONFLICT (internal_name) DO UPDATE SET
    public_name = EXCLUDED.public_name,
    updated_time = EXCLUDED.updated_time,
    deleted_at = NULL;

COMMIT;

SELECT internal_name, public_name
FROM model_public_aliases
WHERE internal_name LIKE 'adobe-firefly-nano-banana%'
ORDER BY internal_name;
