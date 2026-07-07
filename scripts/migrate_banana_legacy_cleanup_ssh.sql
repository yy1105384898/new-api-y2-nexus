-- 下线 legacy Banana 模型（全面切 Manju 新 public 名）+ 清理 ModelPrice
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_banana_legacy_cleanup_ssh.sql

BEGIN;

-- 1. 软删 legacy internal 模型（Manju 系列保留）
UPDATE models SET
    deleted_at = NOW(),
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE deleted_at IS NULL
  AND model_name IN (
    '0lll0-gemini-banana-2.0-pro',
    'byte-gemini-banana-2.0',
    'cy-img2-banana-2.0-pro',
    'cy-img2-banana-flash-image-preview',
    'niming-gemini-banana-2.0-pro'
  );

-- 2. ModelPrice 移除 legacy key（保留 manju-* 与无关项）
UPDATE options SET value = (
    SELECT COALESCE(jsonb_object_agg(k, v), '{}'::jsonb)::text
    FROM jsonb_each(value::jsonb) AS t(k, v)
    WHERE k NOT IN (
        '0lll0-gemini-banana-2.0-pro',
        'byte-gemini-banana-2.0',
        'cy-img2-banana-2.0-pro',
        'cy-img2-banana-flash-image-preview',
        'niming-gemini-banana-2.0-pro'
    )
)
WHERE key = 'ModelPrice';

COMMIT;

SELECT 'models' AS section, model_name, deleted_at IS NOT NULL AS deleted
FROM models WHERE model_name ILIKE '%banana%' ORDER BY 1;

SELECT 'model_price_banana' AS section, k AS model_key
FROM jsonb_each((SELECT value::jsonb FROM options WHERE key = 'ModelPrice')) AS t(k, v)
WHERE k ILIKE '%banana%'
ORDER BY 1;
