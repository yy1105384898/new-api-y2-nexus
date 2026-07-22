-- Adobe2API video ability-group alignment.
-- Channel 75 is a combined media product: its customer keys may use IMAGE while calling both image and video APIs.
-- Run on contabo:
--   docker exec -i newapi-postgres psql -v ON_ERROR_STOP=1 -U root -d new-api \
--     < scripts/migrate_adobe2api_video_ability_groups_ssh.sql

BEGIN;

INSERT INTO abilities ("group", model, channel_id, enabled, priority, weight)
SELECT 'IMAGE', v.model, 75, TRUE, 0, 90
FROM (VALUES
    ('adobe-sora2'),
    ('adobe-sora2-pro'),
    ('adobe-veo31'),
    ('adobe-veo31-ref'),
    ('adobe-veo31-fast')
) AS v(model)
WHERE NOT EXISTS (
    SELECT 1
    FROM abilities a
    WHERE a.channel_id = 75
      AND a."group" = 'IMAGE'
      AND a.model = v.model
);

UPDATE abilities
SET enabled = TRUE,
    priority = 0,
    weight = 90
WHERE channel_id = 75
  AND "group" = 'IMAGE'
  AND model IN (
      'adobe-sora2',
      'adobe-sora2-pro',
      'adobe-veo31',
      'adobe-veo31-ref',
      'adobe-veo31-fast'
  );

COMMIT;

SELECT channel_id, "group", model, enabled, priority, weight
FROM abilities
WHERE channel_id = 75
  AND model IN (
      'adobe-sora2',
      'adobe-sora2-pro',
      'adobe-veo31',
      'adobe-veo31-ref',
      'adobe-veo31-fast'
  )
ORDER BY model, "group";
