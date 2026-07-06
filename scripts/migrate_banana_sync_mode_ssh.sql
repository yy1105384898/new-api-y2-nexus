-- Banana 改回同步 Image API（去掉 async 队列 + 客户端轮询）
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_banana_sync_mode_ssh.sql

BEGIN;

UPDATE model_ui_param_profiles
SET api_mode = 'images-sync-json',
    poll = '{}',
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE profile_id IN ('image-tpl-banana-chat', 'image-tpl-banana-chat-flash-lite');

COMMIT;

SELECT profile_id, api_mode, poll
FROM model_ui_param_profiles
WHERE profile_id IN ('image-tpl-banana-chat', 'image-tpl-banana-chat-flash-lite');
