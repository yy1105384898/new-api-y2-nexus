-- SD5 Seedance 参考图单张上限 10MB -> 30MB（与其它 Seedance 线路对齐）
-- 源站执行：ssh contabo 'docker exec -i newapi-postgres psql -U root -d new-api < migrate_sd5_reference_image_max_bytes_ssh.sql'

UPDATE model_ui_param_profiles
SET reference_limits = (
    COALESCE(NULLIF(reference_limits, ''), '{}')::jsonb
    || '{"imageMaxBytes":31457280}'::jsonb
)::text,
    updated_time = EXTRACT(EPOCH FROM NOW())::bigint
WHERE profile_id = 'video-tpl-cy-sd5-seedance-933-async';
