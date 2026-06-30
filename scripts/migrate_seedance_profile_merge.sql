-- 合并 Seedance async profile：ref9v3 + ref9v3a3 → video-tpl-seedance-async
-- 经济档 / 满血档差异由 infinite-canvas seedance-capability 按模型名判定。

-- 1. 将仍绑定旧 profile 的模型改指向统一 profile
UPDATE models SET video_profile_id = 'video-tpl-seedance-async'
WHERE video_profile_id IN (
    'video-tpl-async-ratio-duration-frame-ref9v3',
    'video-tpl-async-ratio-duration-frame-ref9v3a3'
);

-- 2. 删除已合并的旧 profile 模板
DELETE FROM model_ui_param_profiles
WHERE profile_id IN (
    'video-tpl-async-ratio-duration-frame-ref9v3',
    'video-tpl-async-ratio-duration-frame-ref9v3a3'
);
