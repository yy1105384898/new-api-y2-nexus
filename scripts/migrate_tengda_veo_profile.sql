-- Tengda Geeknow Veo 3.1: profile + model bindings + metadata.
-- Safe to re-run. Internal names use hyphens: tengda-veo-3-1 / tengda-veo-3-1-fast.

INSERT INTO model_ui_param_profiles (
    capability, profile_id, api_mode, requires_reference_media,
    poll, poll_status, reference_limits, params, option_rules, hints,
    created_time, updated_time
)
VALUES (
    'video',
    'video-tpl-async-ratio-duration-ref1-veo',
    'videos-json-async',
    FALSE,
    '{}',
    NULL,
    '{"images":5,"videos":0,"audios":0}',
    '{"resolution":{"enabled":false},"ratio":{"enabled":true,"options":[{"value":"16:9","label":"横屏"},{"value":"9:16","label":"竖屏"}]},"duration":{"enabled":true,"min":1,"max":30,"hint":"单位为秒；文档常见值为 8、10，具体以当前模型/账户为准。"},"generateAudio":{"enabled":false},"watermark":{"enabled":false},"seed":{"enabled":false},"widthHeight":{"enabled":false},"frameInputs":{"enabled":false}}',
    '[]',
    '[{"text":"Geeknow Veo 3.1：文生视频或参考图；画幅 1280×720 / 720×1280；duration 为整数秒（文档示例 8/10）；input_reference 可传 URL/base64，多张时效果取决于模型。"}]',
    EXTRACT(EPOCH FROM NOW())::BIGINT,
    EXTRACT(EPOCH FROM NOW())::BIGINT
)
ON CONFLICT (capability, profile_id) DO UPDATE SET
    api_mode = EXCLUDED.api_mode,
    reference_limits = EXCLUDED.reference_limits,
    params = EXCLUDED.params,
    hints = EXCLUDED.hints,
    updated_time = EXCLUDED.updated_time;

-- Rename legacy internal names (if still present).
UPDATE models SET model_name = 'tengda-veo-3-1', updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name = 'tengda-veo_3_1' AND deleted_at IS NULL;
UPDATE models SET model_name = 'tengda-veo-3-1-fast', updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name = 'tengda-veo_3_1-fast' AND deleted_at IS NULL;
UPDATE abilities SET model = 'tengda-veo-3-1' WHERE model = 'tengda-veo_3_1';
UPDATE abilities SET model = 'tengda-veo-3-1-fast' WHERE model = 'tengda-veo_3_1-fast';

UPDATE models SET
    video_profile_id = 'video-tpl-async-ratio-duration-ref1-veo',
    description = '腾达 Geeknow Veo 3.1 标准质量。文生/图生，1280×720 或 720×1280，duration 整数秒。',
    tags = 'video,veo,tengda,geeknow',
    vendor_id = 6,
    endpoints = '{"openai-video":{"path":"/v1/videos","method":"POST"}}',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name = 'tengda-veo-3-1' AND deleted_at IS NULL;

UPDATE models SET
    video_profile_id = 'video-tpl-async-ratio-duration-ref1-veo',
    description = '腾达 Geeknow Veo 3.1 Fast。更快生成，参数同标准版。',
    tags = 'video,veo,tengda,geeknow,fast',
    vendor_id = 6,
    endpoints = '{"openai-video":{"path":"/v1/videos","method":"POST"}}',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE model_name = 'tengda-veo-3-1-fast' AND deleted_at IS NULL;

DELETE FROM model_public_aliases WHERE internal_name IN ('tengda-veo_3_1', 'tengda-veo_3_1-fast');
