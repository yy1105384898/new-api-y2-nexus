-- Banana / Adobe2API image params: align profile hints and safe aspect ratios
-- with Adobe2API's real OpenAI-compatible input fields.
--
-- contabo:
-- docker exec -i newapi-postgres psql -U root -d new-api < migrate_banana_adobe_params_docs_ssh.sql
-- Then refresh model api_doc with:
-- python3 scripts/seed_manju_gemini_banana_api_doc.py

BEGIN;

UPDATE model_ui_param_profiles
SET params = '{
  "quality": {"enabled": true, "options": [{"value": "auto", "label": "自动"}, {"value": "low", "label": "1K"}, {"value": "medium", "label": "2K"}, {"value": "high", "label": "4K"}]},
  "aspectRatio": {"enabled": true, "options": [
    {"value": "1:1", "label": "1:1", "size": "1:1", "width": 1, "height": 1, "icon": "square"},
    {"value": "16:9", "label": "16:9", "size": "16:9", "width": 16, "height": 9, "icon": "landscape"},
    {"value": "9:16", "label": "9:16", "size": "9:16", "width": 9, "height": 16, "icon": "portrait"},
    {"value": "4:3", "label": "4:3", "size": "4:3", "width": 4, "height": 3, "icon": "landscape"},
    {"value": "3:4", "label": "3:4", "size": "3:4", "width": 3, "height": 4, "icon": "portrait"},
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
  {"text": "API 推荐显式传 aspect_ratio 与 output_resolution（1K/2K/4K）；image_size 是兼容别名，若同时传需与 output_resolution 一致。"},
  {"text": "画质别名：low=1K、medium=2K、high=4K；画布会把 4K 转成 output_resolution:\"4K\"。"},
  {"when": {"modelIncludes": "4k"}, "text": "本模型支持最高 4K 画质，可在参数面板选择。"},
  {"when": {"modelExcludes": "4k"}, "text": "本模型最高支持 2K 画质；名称带 4k 的型号才支持 4K。"}
]'::jsonb::text,
    updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE profile_id = 'image-tpl-banana-chat';

COMMIT;

SELECT profile_id, params, hints
FROM model_ui_param_profiles
WHERE profile_id = 'image-tpl-banana-chat';
