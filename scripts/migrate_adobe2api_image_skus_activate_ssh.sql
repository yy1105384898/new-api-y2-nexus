-- Adobe2API image SKUs, contract/price gate + activation phase.
-- Run only after NewAPI code rollout and real candidate requests succeed.
-- Current public families: nano-banana2 and gpt-image-2. The legacy
-- nano-banana / nano-banana-pro rows stay staged and disabled.

BEGIN;

DO $$
DECLARE
    prices jsonb;
    family text;
    model_count integer;
    ability_count integer;
BEGIN
    SELECT count(*) INTO model_count
    FROM models
    WHERE model_name LIKE 'adobe-firefly-%'
      AND deleted_at IS NULL;
    IF model_count <> 12 THEN
        RAISE EXCEPTION 'expected 12 Adobe Firefly image models, found %', model_count;
    END IF;

    SELECT count(*) INTO ability_count
    FROM abilities
    WHERE channel_id = 75
      AND model LIKE 'adobe-firefly-%';
    IF ability_count <> 36 THEN
        RAISE EXCEPTION 'expected 36 disabled Adobe abilities, found %', ability_count;
    END IF;

    SELECT value::jsonb INTO prices FROM options WHERE key = 'ModelPrice';
    IF prices IS NULL OR NOT prices ?& ARRAY[
        'adobe-firefly-nano-banana2-1k', 'adobe-firefly-nano-banana2-2k', 'adobe-firefly-nano-banana2-4k',
        'adobe-firefly-gpt-image-2-1k', 'adobe-firefly-gpt-image-2-2k', 'adobe-firefly-gpt-image-2-4k'
    ] THEN
        RAISE EXCEPTION 'configure all 6 public Adobe Firefly ModelPrice keys before activation';
    END IF;

    FOREACH family IN ARRAY ARRAY[
        'adobe-firefly-nano-banana2',
        'adobe-firefly-gpt-image-2'
    ] LOOP
        IF (prices->>(family || '-1k'))::numeric <= 0
           OR (prices->>(family || '-1k'))::numeric >= (prices->>(family || '-2k'))::numeric
           OR (prices->>(family || '-2k'))::numeric >= (prices->>(family || '-4k'))::numeric THEN
            RAISE EXCEPTION 'ModelPrice must satisfy 0 < %-1k < %-2k < %-4k', family, family, family;
        END IF;
    END LOOP;
END $$;

UPDATE models
SET status = 1,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name ~ '^adobe-firefly-(nano-banana2|gpt-image-2)-(1k|2k|4k)$'
  AND deleted_at IS NULL;

UPDATE abilities
SET enabled = TRUE
WHERE channel_id = 75
  AND model ~ '^adobe-firefly-(nano-banana2|gpt-image-2)-(1k|2k|4k)$';

UPDATE models
SET status = 0,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name ~ '^adobe-firefly-(nano-banana-pro|nano-banana)-(1k|2k|4k)$'
  AND deleted_at IS NULL;

UPDATE abilities
SET enabled = FALSE
WHERE channel_id = 75
  AND model ~ '^adobe-firefly-(nano-banana-pro|nano-banana)-(1k|2k|4k)$';

COMMIT;

-- Rolling-restart NewAPI after this transaction so the public-name registry is refreshed.
SELECT model_name, status, image_profile_id
FROM models
WHERE model_name LIKE 'adobe-firefly-%'
ORDER BY model_name;

SELECT "group", count(*) AS enabled_abilities
FROM abilities
WHERE channel_id = 75
  AND model ~ '^adobe-firefly-(nano-banana2|gpt-image-2)-(1k|2k|4k)$'
  AND enabled = TRUE
GROUP BY "group"
ORDER BY "group";
