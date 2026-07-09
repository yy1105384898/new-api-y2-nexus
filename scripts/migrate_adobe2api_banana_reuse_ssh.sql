-- Adobe2API channel 75: reuse current Gemini Banana public models through Adobe2API.
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_adobe2api_banana_reuse_ssh.sql

BEGIN;

-- Keep the public/internal model names already used by clients, but route channel 75
-- to Adobe2API's upstream model names.
WITH existing_models AS (
    SELECT c.id, item.model, item.ord
    FROM channels c
    CROSS JOIN LATERAL unnest(string_to_array(COALESCE(c.models, ''), ',')) WITH ORDINALITY AS item(model, ord)
    WHERE c.id = 75 AND btrim(item.model) <> ''
),
wanted_models AS (
    SELECT * FROM (VALUES
        ('manju-gemini-banana-pro-4k', 10001),
        ('manju-gemini-banana-2.0-4k', 10002)
    ) AS v(model, ord)
),
merged_models AS (
    SELECT model, min(ord) AS ord
    FROM (
        SELECT btrim(model) AS model, ord FROM existing_models
        UNION ALL
        SELECT model, ord FROM wanted_models
    ) s
    GROUP BY model
),
existing_groups AS (
    SELECT c.id, item.grp, item.ord
    FROM channels c
    CROSS JOIN LATERAL unnest(string_to_array(COALESCE(c."group", ''), ',')) WITH ORDINALITY AS item(grp, ord)
    WHERE c.id = 75 AND btrim(item.grp) <> ''
),
wanted_groups AS (
    SELECT * FROM (VALUES
        ('IMAGE', 10001),
        ('全模型-无claude/gpt', 10002)
    ) AS v(grp, ord)
),
merged_groups AS (
    SELECT grp, min(ord) AS ord
    FROM (
        SELECT btrim(grp) AS grp, ord FROM existing_groups
        UNION ALL
        SELECT grp, ord FROM wanted_groups
    ) s
    GROUP BY grp
)
UPDATE channels
SET
    models = (SELECT string_agg(model, ',' ORDER BY ord) FROM merged_models),
    model_mapping = (
        COALESCE(NULLIF(model_mapping, '')::jsonb, '{}'::jsonb)
        || '{
          "manju-gemini-banana-pro-4k": "nano-banana-pro",
          "manju-gemini-banana-2.0-4k": "nano-banana2"
        }'::jsonb
    )::text,
    "group" = (SELECT string_agg(grp, ',' ORDER BY ord) FROM merged_groups),
    status = 1
WHERE id = 75;

DELETE FROM abilities
WHERE channel_id = 75
  AND model IN ('manju-gemini-banana-pro-4k', 'manju-gemini-banana-2.0-4k');

INSERT INTO abilities ("group", model, channel_id, enabled, priority, weight)
SELECT g.grp, m.model, 75, true, 0, 90
FROM (VALUES
    ('manju-gemini-banana-pro-4k'),
    ('manju-gemini-banana-2.0-4k')
) AS m(model)
CROSS JOIN (VALUES ('IMAGE'), ('全模型-无claude/gpt')) AS g(grp);

UPDATE models AS m
SET
    description = v.description,
    tags = v.tags,
    endpoints = '["openai"]',
    image_profile_id = 'image-tpl-banana-chat',
    status = 1,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
FROM (VALUES
    ('manju-gemini-banana-pro-4k', 'Gemini Banana Pro 4K。同步/异步出图，支持 4K。', 'image,gemini,banana,pro,4k'),
    ('manju-gemini-banana-2.0-4k', 'Nano Banana 2.0 4K。同步/异步出图，支持 4K。', 'image,gemini,banana,2.0,4k')
) AS v(model_name, description, tags)
WHERE m.model_name = v.model_name
  AND m.deleted_at IS NULL;

COMMIT;

SELECT id, name, "group", models, model_mapping, status
FROM channels
WHERE id = 75;

SELECT channel_id, "group", model, enabled, priority, weight
FROM abilities
WHERE channel_id = 75
  AND model IN ('manju-gemini-banana-pro-4k', 'manju-gemini-banana-2.0-4k')
ORDER BY model, "group";
