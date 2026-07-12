-- Adobe2API Veo 3.1 profile split.
-- Standard/Fast use frame references (max 2); Ref uses asset/image references (max 3).
-- Run on contabo:
--   docker exec -i newapi-postgres psql -v ON_ERROR_STOP=1 -U root -d new-api \
--     < scripts/migrate_adobe2api_veo31_profiles_ssh.sql

BEGIN;

UPDATE model_ui_param_profiles
SET reference_limits = '{"images":2,"videos":0,"audios":0}',
    hints = '[{"text":"Adobe Firefly Veo 3.1 标准版/Fast：reference_mode=frame，最多 2 张首尾帧参考图；支持 4/6/8 秒、16:9/9:16、720p/1080p。"}]',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT,
    deleted_at = NULL
WHERE capability = 'video'
  AND profile_id = 'video-tpl-adobe-veo31-json-async';

INSERT INTO model_ui_param_profiles (
    capability, profile_id, api_mode, requires_reference_media,
    poll, poll_status, reference_limits, params, option_rules, hints,
    created_time, updated_time
)
VALUES (
    'video', 'video-tpl-adobe-veo31-ref-json-async', 'videos-json-async', FALSE,
    '{}', NULL,
    '{"images":3,"videos":0,"audios":0}',
    '{"resolution":{"enabled":true,"options":[{"value":"720p","label":"720p"},{"value":"1080p","label":"1080p"}]},"ratio":{"enabled":true,"options":[{"value":"16:9","label":"横屏"},{"value":"9:16","label":"竖屏"}]},"duration":{"enabled":true,"min":4,"max":8,"numericOptions":[4,6,8]},"generateAudio":{"enabled":true},"watermark":{"enabled":false},"seed":{"enabled":false},"widthHeight":{"enabled":false},"frameInputs":{"enabled":false}}',
    '[]',
    '[{"text":"Adobe Firefly Veo 3.1 Ref：reference_mode=image，最多 3 张主体或素材参考图，不表示首尾帧；支持 4/6/8 秒、16:9/9:16、720p/1080p。"}]',
    EXTRACT(EPOCH FROM NOW())::BIGINT,
    EXTRACT(EPOCH FROM NOW())::BIGINT
)
ON CONFLICT (capability, profile_id) DO UPDATE SET
    api_mode = EXCLUDED.api_mode,
    requires_reference_media = EXCLUDED.requires_reference_media,
    poll = EXCLUDED.poll,
    poll_status = EXCLUDED.poll_status,
    reference_limits = EXCLUDED.reference_limits,
    params = EXCLUDED.params,
    option_rules = EXCLUDED.option_rules,
    hints = EXCLUDED.hints,
    updated_time = EXCLUDED.updated_time,
    deleted_at = NULL;

UPDATE models
SET video_profile_id = CASE model_name
        WHEN 'adobe-veo31-ref' THEN 'video-tpl-adobe-veo31-ref-json-async'
        ELSE 'video-tpl-adobe-veo31-json-async'
    END,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name IN ('adobe-veo31', 'adobe-veo31-ref', 'adobe-veo31-fast')
  AND deleted_at IS NULL;

COMMIT;

SELECT model_name, video_profile_id
FROM models
WHERE model_name IN ('adobe-veo31', 'adobe-veo31-ref', 'adobe-veo31-fast')
  AND deleted_at IS NULL
ORDER BY model_name;

SELECT profile_id, reference_limits, hints
FROM model_ui_param_profiles
WHERE capability = 'video'
  AND profile_id IN (
      'video-tpl-adobe-veo31-json-async',
      'video-tpl-adobe-veo31-ref-json-async'
  )
  AND deleted_at IS NULL
ORDER BY profile_id;
