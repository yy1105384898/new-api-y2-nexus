-- Adobe2API API-doc audit corrections: profile limits and staged-model abilities.
-- This migration changes no prices and does not restart services.

BEGIN;

UPDATE model_ui_param_profiles
SET reference_limits = '{"images":1,"videos":0,"audios":0}',
    hints = '[{"text":"Adobe Firefly Sora2/Sora2 Pro：支持 4/8/12 秒、16:9/9:16、单张帧参考图、negative_prompt 与声音开关。"}]',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT,
    deleted_at = NULL
WHERE capability = 'video'
  AND profile_id = 'video-tpl-adobe-sora2-json-async';

-- The legacy nano-banana family is staged (models.status=0) and must not keep live abilities.
UPDATE abilities
SET enabled = FALSE
WHERE channel_id = 75
  AND model ~ '^adobe-firefly-nano-banana-(1k|2k|4k)$';

-- Every active Adobe media product is available to combined IMAGE keys as well as VIDEO/general keys.
UPDATE abilities
SET enabled = TRUE
WHERE channel_id = 75
  AND model ~ '^adobe-(sora2(-pro)?|veo31(-ref|-fast)?|firefly-(nano-banana-pro|nano-banana2|gpt-image-2)-(1k|2k|4k))$';

COMMIT;

SELECT profile_id, reference_limits, hints
FROM model_ui_param_profiles
WHERE capability = 'video'
  AND profile_id = 'video-tpl-adobe-sora2-json-async';

SELECT model, bool_or(enabled) AS any_enabled
FROM abilities
WHERE channel_id = 75
  AND model LIKE 'adobe-firefly-nano-banana-%'
GROUP BY model
ORDER BY model;
