-- Banana 画布参数与 hints：画幅改比例值、启用画质、关闭自定义像素、hints 改客户向说明
-- contabo: docker exec -i newapi-postgres psql -U root -d new-api < migrate_banana_profile_params_ssh.sql
-- 或重跑: go run ./scripts/seed_model_ui_params/main.go -force

BEGIN;

UPDATE model_ui_param_profiles
SET params = '{
  "quality": {"enabled": true, "options": [{"value": "auto", "label": "自动"}, {"value": "low", "label": "1K"}, {"value": "medium", "label": "2K"}, {"value": "high", "label": "4K"}]},
  "aspectRatio": {"enabled": true, "options": [
    {"value": "1:1", "label": "1:1", "size": "1:1", "width": 1, "height": 1, "icon": "square"},
    {"value": "16:9", "label": "16:9", "size": "16:9", "width": 16, "height": 9, "icon": "landscape"},
    {"value": "9:16", "label": "9:16", "size": "9:16", "width": 9, "height": 16, "icon": "portrait"},
    {"value": "3:2", "label": "3:2", "size": "3:2", "width": 3, "height": 2, "icon": "landscape"},
    {"value": "2:3", "label": "2:3", "size": "2:3", "width": 2, "height": 3, "icon": "portrait"},
    {"value": "4:3", "label": "4:3", "size": "4:3", "width": 4, "height": 3, "icon": "landscape"},
    {"value": "3:4", "label": "3:4", "size": "3:4", "width": 3, "height": 4, "icon": "portrait"},
    {"value": "4:5", "label": "4:5", "size": "4:5", "width": 4, "height": 5, "icon": "portrait"},
    {"value": "5:4", "label": "5:4", "size": "5:4", "width": 5, "height": 4, "icon": "landscape"},
    {"value": "21:9", "label": "21:9", "size": "21:9", "width": 21, "height": 9, "icon": "landscape"},
    {"value": "auto", "label": "自动", "width": 0, "height": 0, "icon": "auto"}
  ]},
  "customDimensions": {"enabled": false},
  "count": {"enabled": true, "min": 1, "max": 4, "quickCount": 4},
  "background": {"enabled": false},
  "outputFormat": {"enabled": false},
  "outputCompression": {"enabled": false},
  "moderation": {"enabled": false}
}'::jsonb::text,
    hints = '[
  {"text": "连接参考图后，在提示词中用 @图片1 说明要改动的素材。"},
  {"text": "画幅选「自动」时，有参考图会按参考图比例出图。"},
  {"when": {"modelIncludes": "4k"}, "text": "本模型支持最高 4K 画质，可在参数面板选择。"},
  {"when": {"modelExcludes": "4k"}, "text": "本模型最高支持 2K 画质；名称带 4k 的型号才支持 4K。"}
]'::jsonb::text,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE profile_id = 'image-tpl-banana-chat';

UPDATE model_ui_param_profiles
SET hints = '[
  {"text": "轻量快速出图，仅支持 1K 画质（约 1024px）。"},
  {"text": "画幅请选 1:1、16:9 等比例，或选「自动」按参考图推断。"},
  {"text": "连接参考图后，在提示词中用 @图片1 说明要改动的素材。"}
]'::jsonb::text,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE profile_id = 'image-tpl-banana-chat-flash-lite';

UPDATE model_ui_param_profiles
SET hints = '[
  {"text": "Flash Lite 图像模型仅支持 1K 出图（约 1024px），不支持 2K/4K。"},
  {"text": "请使用 1:1、16:9 等比例；画质选 1K 或自动即可。"},
  {"text": "连接参考图后，在提示词中用 @图片1 说明要改动的素材。"}
]'::jsonb::text,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE profile_id = 'image-tpl-aspect-count-flash-lite';

COMMIT;

SELECT profile_id, left(params, 80) AS params_preview, hints
FROM model_ui_param_profiles
WHERE profile_id IN ('image-tpl-banana-chat', 'image-tpl-banana-chat-flash-lite', 'image-tpl-aspect-count-flash-lite');
