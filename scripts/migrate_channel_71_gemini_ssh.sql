-- 渠道 71 对齐渠道 58：Gemini 类型 + gemini-3-pro-image-preview → generateContent
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_channel_71_gemini_ssh.sql

BEGIN;

UPDATE channels SET
    type = 24,
    status = 1,
    model_mapping = '{
  "manju-gemini-banana-pro-4k": "gemini-3-pro-image-preview"
}'::text
WHERE id = 71;

COMMIT;

SELECT id, name, type, status, base_url, model_mapping FROM channels WHERE id IN (58, 71);
