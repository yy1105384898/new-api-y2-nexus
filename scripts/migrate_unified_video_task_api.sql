-- 视频对外协议收口：所有画布 profile 只调用 POST/GET /v1/videos。
-- 上游 chat、Grok generations、multipart 等差异由 oaivideo vendor 层适配。

BEGIN;

UPDATE model_ui_param_profiles
SET payload_builder = CASE api_mode
        WHEN 'videos-form' THEN 'openai-form'
        WHEN 'chat-completions' THEN 'chat-video'
        WHEN 'video-generations' THEN 'grok-generations'
        ELSE payload_builder
    END,
    api_mode = 'videos-json-async',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE capability = 'video'
  AND deleted_at IS NULL;

UPDATE model_ui_param_registries
SET poll_defaults = '{"videos-json-async":{"delayMs":5000,"maxAttempts":120}}',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE capability = 'video'
  AND deleted_at IS NULL;

COMMIT;
