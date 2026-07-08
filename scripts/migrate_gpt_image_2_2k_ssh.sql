-- gpt-image-2-2k（cy-img2-gpt-image-2-2k / Gulie 渠道 72）：绑定 Gulie 2K profile + 描述
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_gpt_image_2_2k_ssh.sql

UPDATE models
SET image_profile_id = 'image-tpl-gulie-2k',
    description = 'GPT-Image-2 2K 经济档（Gulie 线路），size 传画幅比例；参考图 multipart /images/edits。',
    updated_time = extract(epoch from now())::bigint
WHERE model_name = 'cy-img2-gpt-image-2-2k'
  AND deleted_at IS NULL;

SELECT model_name, image_profile_id, left(description, 80) AS desc_preview
FROM models
WHERE model_name = 'cy-img2-gpt-image-2-2k'
  AND deleted_at IS NULL;
