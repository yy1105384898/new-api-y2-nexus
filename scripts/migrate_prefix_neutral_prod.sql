-- 渠道前缀中性化 + tags/description 清洗（生产 PostgreSQL 一次性执行）
-- 目标：internal 注册名改为 cy-* 路由码；public 名与 upstream model_mapping 值不变，下游 API 不受影响。
--
-- 执行前：先部署含 cy-* 前缀判断的 new-api / infinite-canvas 代码。
-- 执行后：滚动重启 new-api（刷新 model_channel_prefixes 内存缓存）。
--
--   docker exec -i newapi-postgres psql -U root -d new-api < migrate_prefix_neutral_prod.sql

BEGIN;

-- ---------------------------------------------------------------------------
-- 0. 文本替换：按前缀长度降序，避免短前缀误伤
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION cy_replace_vendor_prefixes(input text)
RETURNS text
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
  out text := COALESCE(input, '');
BEGIN
  IF out = '' THEN
    RETURN out;
  END IF;

  -- internal 渠道前缀（长串优先）
  out := replace(out, 'oairegbox-', 'cy-sd1-');
  out := replace(out, 'happyhorse-', 'cy-vid1-');
  out := replace(out, 'leonardo-', 'cy-sd4-');
  out := replace(out, '119337-', 'cy-gv1-');
  out := replace(out, 'ctlove-', 'cy-sd3-');
  out := replace(out, 'tengda-', 'cy-veo1-');
  out := replace(out, 'tengd-', 'cy-sd2-');
  out := replace(out, 'yunwu-', 'cy-vid2-');
  out := replace(out, 'geek2-', 'cy-img2-');
  out := replace(out, 'Gulie-', 'cy-img1-');
  out := replace(out, 'gulie-', 'cy-img1-');
  out := replace(out, 'gz-seedance-', 'cy-sd0-seedance-');
  out := replace(out, 'gz-video-', 'cy-sd0-video-');
  out := replace(out, 'gz-', 'cy-sd0-');

  -- 渠道名里的上游域名
  out := replace(out, 'https://newapi-2.oairegbox.cc', '');
  out := replace(out, 'https://newapi.oairegbox.cc', '');
  out := replace(out, 'https://td.geeknow.top', '');
  out := replace(out, 'https://api.119337.xyz', '');
  out := replace(out, 'https://apidoc.geeknow.top', '');

  RETURN out;
END;
$$;

CREATE OR REPLACE FUNCTION cy_strip_vendor_tags(input text)
RETURNS text
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
  out text := COALESCE(input, '');
  tok text;
  deny text[] := ARRAY[
    'oairegbox', 'geeknow', '119337', 'tengda', 'tengd', 'ctlove',
    'happyhorse', 'yunwu', 'leonardo', 'gulie', 'geek2', 'gz',
    'oairegbox', 'special-offer', 'subscription', 'geek2api', 'gulieapi'
  ];
BEGIN
  IF out = '' THEN
    RETURN out;
  END IF;
  FOREACH tok IN ARRAY deny
  LOOP
    out := regexp_replace(out, '(^|,)\s*' || tok || '\s*(,|$)', ',', 'gi');
  END LOOP;
  out := regexp_replace(out, ',+', ',', 'g');
  out := regexp_replace(out, '^,|,$', '', 'g');
  out := trim(both ' ' from out);
  RETURN out;
END;
$$;

CREATE OR REPLACE FUNCTION cy_sanitize_description(input text)
RETURNS text
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
  out text := COALESCE(input, '');
BEGIN
  IF out = '' THEN
    RETURN out;
  END IF;
  out := regexp_replace(out, 'OAIREGBox\s*', '', 'gi');
  out := regexp_replace(out, '119337\s*', '', 'gi');
  out := regexp_replace(out, 'Geeknow\s*', '', 'gi');
  out := regexp_replace(out, 'Geek2(API)?\s*', '', 'gi');
  out := regexp_replace(out, 'Gulie\s*', '', 'gi');
  out := regexp_replace(out, 'Leonardo\s*订阅号\s*', '', 'gi');
  out := regexp_replace(out, 'CTLove\s*', '', 'gi');
  out := regexp_replace(out, 'HappyHorse\s*', '', 'gi');
  out := regexp_replace(out, '云雾\s*', '', 'gi');
  out := regexp_replace(out, '腾达\s*', '', 'gi');
  out := regexp_replace(out, 'GZ\s*/\s*冠臻\s*', '', 'gi');
  out := trim(both ' ' from out);
  RETURN out;
END;
$$;

CREATE OR REPLACE FUNCTION cy_rename_json_object_keys(j jsonb)
RETURNS jsonb
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE
  k text;
  v jsonb;
  nk text;
  out jsonb := '{}'::jsonb;
BEGIN
  IF j IS NULL OR j = 'null'::jsonb THEN
    RETURN j;
  END IF;
  FOR k, v IN SELECT * FROM jsonb_each(j)
  LOOP
    nk := cy_replace_vendor_prefixes(k);
    out := out || jsonb_build_object(nk, v);
  END LOOP;
  RETURN out;
END;
$$;

CREATE OR REPLACE FUNCTION cy_rename_json_object_keys_text(input text)
RETURNS text
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
  IF input IS NULL OR btrim(input) = '' THEN
    RETURN input;
  END IF;
  RETURN cy_rename_json_object_keys(input::jsonb)::text;
END;
$$;

-- ---------------------------------------------------------------------------
-- 1. model_channel_prefixes：改 prefix + 中性 note
-- ---------------------------------------------------------------------------
UPDATE model_channel_prefixes SET
  prefix = 'cy-gv1-',
  note = 'Grok 视频 · video-generations',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = '119337-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-sd1-',
  note = 'Seedance/Omni 线路 A',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'oairegbox-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-veo1-',
  note = 'Veo 3.1 异步 JSON',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'tengda-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-sd3-',
  note = 'Seedance 2.0 线路 C',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'ctlove-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-vid1-',
  note = '通用视频线路',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'happyhorse-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-vid2-',
  note = '聚合视频线路',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'yunwu-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-sd4-',
  note = 'Seedance 2.0 订阅号',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'leonardo-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-sd0-',
  note = 'Seedance 历史线路（已下线）',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'gz-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-sd2-',
  note = 'Seedance 2.0 特惠',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'tengd-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-img1-',
  note = '生图线路 A',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'Gulie-' AND deleted_at IS NULL;

UPDATE model_channel_prefixes SET
  prefix = 'cy-img2-',
  note = '生图线路 B（4K）',
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE prefix = 'geek2-' AND deleted_at IS NULL;

INSERT INTO model_channel_prefixes (prefix, note, enabled, sort_order, created_time, updated_time)
VALUES
  ('cy-sd2-', 'Seedance 2.0 特惠', TRUE, 106, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
  ('cy-img1-', '生图线路 A', TRUE, 130, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT),
  ('cy-img2-', '生图线路 B（4K）', TRUE, 131, EXTRACT(EPOCH FROM NOW())::BIGINT, EXTRACT(EPOCH FROM NOW())::BIGINT)
ON CONFLICT (prefix) DO UPDATE SET
  note = EXCLUDED.note,
  enabled = EXCLUDED.enabled,
  updated_time = EXCLUDED.updated_time;

-- ---------------------------------------------------------------------------
-- 2. models 表：internal 注册名
-- ---------------------------------------------------------------------------
UPDATE models SET
  model_name = cy_replace_vendor_prefixes(model_name),
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE deleted_at IS NULL
  AND model_name ~* '(^|.*)(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)';

-- ---------------------------------------------------------------------------
-- 2b. models 表：用户可见元信息（/api/pricing → 模型广场 / 查看文档）
--     model_name 对外已是 public；description / tags / api_doc 必须去商家词。
--     vendor_id 保留（Google/OpenAI/xAI 等能力厂商，非转售渠道名）。
-- ---------------------------------------------------------------------------
UPDATE models SET
  description = cy_sanitize_description(cy_replace_vendor_prefixes(COALESCE(description, ''))),
  tags = cy_strip_vendor_tags(cy_replace_vendor_prefixes(COALESCE(tags, ''))),
  api_doc = NULLIF(
    cy_sanitize_description(cy_replace_vendor_prefixes(COALESCE(api_doc, ''))),
    ''
  ),
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE deleted_at IS NULL
  AND (
    COALESCE(description, '') ~* '(OAIREGBox|119337|Geeknow|Geek2|Gulie|Leonardo|CTLove|HappyHorse|云雾|腾达|oairegbox|geeknow|tengda|tengd|yunwu|ctlove|happyhorse|gulie|geek2)'
    OR COALESCE(tags, '') ~* '(oairegbox|geeknow|119337|tengda|tengd|ctlove|happyhorse|yunwu|leonardo|gulie|geek2|subscription|special-offer)'
    OR COALESCE(api_doc, '') ~* '(oairegbox|119337|tengda|tengd-|gulie|geek2|leonardo|yunwu|ctlove|happyhorse|OAIREGBox|Geeknow|119337)'
  );

-- ---------------------------------------------------------------------------
-- 3. abilities / channels
-- ---------------------------------------------------------------------------
UPDATE abilities SET
  model = cy_replace_vendor_prefixes(model)
WHERE model ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)';

UPDATE channels SET
  models = cy_replace_vendor_prefixes(models),
  name = cy_sanitize_description(cy_replace_vendor_prefixes(name)),
  model_mapping = cy_rename_json_object_keys_text(model_mapping)
WHERE models ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)'
   OR model_mapping ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)'
   OR name ~* '(oairegbox|geeknow|119337|td\.geeknow)';

-- ---------------------------------------------------------------------------
-- 4. public 别名（internal 侧改名；public 名不变）
-- ---------------------------------------------------------------------------
UPDATE model_public_aliases SET
  internal_name = cy_replace_vendor_prefixes(internal_name),
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE internal_name ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)';

-- ---------------------------------------------------------------------------
-- 5. options JSON（ModelPrice / Ratio / billing 等键名）
-- ---------------------------------------------------------------------------
UPDATE options SET value = cy_rename_json_object_keys_text(value)
WHERE key IN (
  'ModelPrice', 'ModelRatio', 'CompletionRatio', 'CacheRatio', 'CreateCacheRatio',
  'ImageRatio', 'AudioRatio', 'AudioCompletionRatio', 'TopUpGroupRatio'
)
AND value ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)';

UPDATE options SET value = cy_rename_json_object_keys_text(value)
WHERE key IN ('billing_setting.billing_mode', 'billing_setting.billing_expr', 'billing_setting.request_unit')
AND value ~* '(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)';

-- ---------------------------------------------------------------------------
-- 6. UI profile hints / rules（用户可见文案）
-- ---------------------------------------------------------------------------
UPDATE model_ui_param_profiles SET
  option_rules = cy_replace_vendor_prefixes(option_rules),
  hints = cy_replace_vendor_prefixes(
    replace(
      replace(hints, '（oairegbox-seedance-2.0-*）', ''),
      'oairegbox-seedance-2.0', 'seedance-2.0'
    )
  ),
  updated_time = EXTRACT(EPOCH FROM NOW())::BIGINT
WHERE option_rules ~* '(oairegbox|119337|tengda|tengd-|leonardo|gulie|geek2|ctlove|yunwu|happyhorse)'
   OR hints ~* '(oairegbox|119337|tengda|tengd-|leonardo|gulie|geek2|ctlove|yunwu|happyhorse)';

COMMIT;

-- ---------------------------------------------------------------------------
-- 7. 验收（应无 internal 旧前缀；public 名抽样）
-- ---------------------------------------------------------------------------
SELECT 'model_channel_prefixes' AS section, prefix, note
FROM model_channel_prefixes
WHERE deleted_at IS NULL
ORDER BY sort_order, prefix;

SELECT 'models_still_old_prefix' AS check, model_name
FROM models
WHERE deleted_at IS NULL
  AND model_name ~* '^(oairegbox|119337|tengda|tengd-|ctlove|happyhorse|yunwu|leonardo|Gulie|gulie|geek2|gz-)'
LIMIT 20;

SELECT 'tags_still_vendor' AS check, model_name, tags
FROM models
WHERE deleted_at IS NULL
  AND tags ~* '(oairegbox|geeknow|119337|tengda|tengd|ctlove|happyhorse|yunwu|leonardo|gulie|geek2)'
LIMIT 20;

SELECT 'description_still_vendor' AS check, model_name, left(description, 80) AS description
FROM models
WHERE deleted_at IS NULL
  AND description ~* '(OAIREGBox|119337|Geeknow|Geek2|Gulie|Leonardo订阅号|CTLove|HappyHorse|云雾|腾达)'
LIMIT 20;

SELECT 'api_doc_still_vendor' AS check, model_name
FROM models
WHERE deleted_at IS NULL
  AND api_doc ~* '(oairegbox|119337|OAIREGBox|Geeknow|tengda|tengd-|gulie|geek2|leonardo|yunwu|ctlove|happyhorse)'
LIMIT 20;

SELECT 'public_name_sample' AS check, m.model_name AS internal, a.public_name
FROM model_public_aliases a
JOIN models m ON m.model_name = a.internal_name AND m.deleted_at IS NULL
WHERE a.deleted_at IS NULL
ORDER BY a.public_name
LIMIT 15;

-- 清理临时函数（可选，保留便于二次 patch）
-- DROP FUNCTION IF EXISTS cy_rename_json_object_keys(jsonb);
-- DROP FUNCTION IF EXISTS cy_sanitize_description(text);
-- DROP FUNCTION IF EXISTS cy_strip_vendor_tags(text);
-- DROP FUNCTION IF EXISTS cy_replace_vendor_prefixes(text);
