-- 下线 GZ 直连 Seedance（videos-json-gz / video-tpl-json-*）
-- 将仍绑定 json profile 的模型改指向 async Seedance profile；禁用 gz-seedance 渠道能力。

-- 1. 删除 json-gz profile 模板
DELETE FROM model_ui_param_profiles
WHERE profile_id LIKE 'video-tpl-json-%'
   OR api_mode = 'videos-json-gz';

-- 2. 原 GZ 档位 → 满血 async profile（Pro/Fast 720p、15s）
UPDATE models SET video_profile_id = 'video-tpl-seedance-async'
WHERE video_profile_id IN (
    'video-tpl-json-pro-audio-seed-ref4v1',
    'video-tpl-json-fast-audio-seed-ref4v1',
    'video-tpl-json-15s-pro-audio-ref4v1a3',
    'video-tpl-json-15s-fast-audio-ref4v1a1'
);

-- 3. 原 GZ 10s / plain / sz → 统一 Seedance async profile
UPDATE models SET video_profile_id = 'video-tpl-seedance-async'
WHERE video_profile_id IN (
    'video-tpl-json-10s-seed-ref9',
    'video-tpl-json-art-plain-ref4',
    'video-tpl-json-720-duration-ref4'
);

-- 4. 禁用 gz-seedance 渠道能力（保留 oairegbox/ctlove Seedance）
UPDATE abilities SET enabled = false
WHERE model LIKE 'gz-seedance-%' OR model LIKE 'gz-video-%';

-- 5. 清理 poll_defaults 中的 videos-json-gz（若存在）
UPDATE model_ui_param_registries
SET poll_defaults = (poll_defaults::jsonb - 'videos-json-gz')::text
WHERE poll_defaults::jsonb ? 'videos-json-gz';
