-- OAIREGBox Seedance 实际成片时长补结算（2026-07-11）。
--
-- 执行方式（仅在包含 duration 解析/实际 MP4 时长结算修复的版本部署后）：
--   docker exec -i newapi-postgres psql -v ON_ERROR_STOP=1 -U root -d new-api \
--     < scripts/backfill_oairegbox_duration_billing_20260711.sql
--
-- 每条任务均已用 ffprobe 核验实际成片时长。脚本通过
-- private_data.billing_adjustments.oairegbox_duration_v1 保证幂等。

BEGIN;

CREATE TEMP TABLE verified_duration_adjustments (
    task_id bigint PRIMARY KEY,
    prior_quota bigint NOT NULL,
    actual_seconds integer NOT NULL
) ON COMMIT DROP;

INSERT INTO verified_duration_adjustments (task_id, prior_quota, actual_seconds) VALUES
    (8004,  500000,  4),
    (8753,  225000,  4),
    (8754,  225000,  4),
    (8777,  325000,  4),
    (8790,  225000,  4),
    (9251, 1300000,  6),
    (9462, 3000000,  6),
    (9515, 6000000, 15);

DO $$
DECLARE
    found_count integer;
    invalid_count integer;
BEGIN
    SELECT count(*) INTO found_count
    FROM tasks t JOIN verified_duration_adjustments v ON v.task_id = t.id;
    IF found_count <> 8 THEN
        RAISE EXCEPTION 'expected 8 verified tasks, found %', found_count;
    END IF;

    SELECT count(*) INTO invalid_count
    FROM tasks t
    JOIN verified_duration_adjustments v ON v.task_id = t.id
    WHERE t.status <> 'SUCCESS'
       OR t.private_data::jsonb #>> '{billing_source}' <> 'wallet'
       OR (
            NOT coalesce(t.private_data::jsonb #> '{billing_adjustments}' ? 'oairegbox_duration_v1', false)
            AND t.quota <> v.prior_quota
       );
    IF invalid_count <> 0 THEN
        RAISE EXCEPTION '% verified tasks no longer match the audited billing snapshot', invalid_count;
    END IF;
END $$;

CREATE TEMP TABLE pending_duration_adjustments ON COMMIT DROP AS
SELECT
    t.id AS task_id,
    t.user_id,
    t.channel_id,
    t.group AS user_group,
    t.created_at,
    t.properties->>'origin_model_name' AS model_name,
    (t.private_data::jsonb #>> '{token_id}')::bigint AS token_id,
    (t.private_data::jsonb #>> '{billing_context,other_ratios,seconds}')::numeric AS billed_seconds,
    v.actual_seconds,
    t.quota AS prior_quota,
    round(
        (t.private_data::jsonb #>> '{billing_context,model_price}')::numeric
        * 500000
        * coalesce(nullif((t.private_data::jsonb #>> '{billing_context,group_ratio}')::numeric, 0), 1)
        * coalesce(nullif((t.private_data::jsonb #>> '{billing_context,other_ratios,size}')::numeric, 0), 1)
        * v.actual_seconds
    )::bigint AS actual_quota
FROM tasks t
JOIN verified_duration_adjustments v ON v.task_id = t.id
WHERE NOT coalesce(t.private_data::jsonb #> '{billing_adjustments}' ? 'oairegbox_duration_v1', false);

ALTER TABLE pending_duration_adjustments ADD COLUMN quota_delta bigint;
UPDATE pending_duration_adjustments SET quota_delta = actual_quota - prior_quota;

UPDATE users u
SET quota = u.quota - x.net_delta,
    used_quota = u.used_quota + x.positive_delta
FROM (
    SELECT user_id,
           sum(quota_delta) AS net_delta,
           sum(greatest(quota_delta, 0)) AS positive_delta
    FROM pending_duration_adjustments
    GROUP BY user_id
) x
WHERE u.id = x.user_id;

-- 这些核验任务均使用 unlimited token；保留通用分支，避免未来复用时漏扣有限令牌。
UPDATE tokens tk
SET remain_quota = tk.remain_quota - x.net_delta
FROM (
    SELECT token_id, sum(quota_delta) AS net_delta
    FROM pending_duration_adjustments
    WHERE token_id IS NOT NULL
    GROUP BY token_id
) x
WHERE tk.id = x.token_id AND NOT tk.unlimited_quota;

UPDATE channels c
SET used_quota = c.used_quota + x.positive_delta
FROM (
    SELECT channel_id, sum(greatest(quota_delta, 0)) AS positive_delta
    FROM pending_duration_adjustments
    GROUP BY channel_id
) x
WHERE c.id = x.channel_id;

INSERT INTO logs (
    user_id, created_at, type, content, username, token_name, model_name,
    quota, prompt_tokens, completion_tokens, use_time, is_stream,
    channel_id, token_id, "group", other
)
SELECT
    p.user_id,
    extract(epoch FROM now())::bigint,
    CASE WHEN p.quota_delta > 0 THEN 2 ELSE 6 END,
    'OAIREGBox Seedance 实际成片时长补结算',
    u.username,
    coalesce(tk.name, ''),
    p.model_name,
    abs(p.quota_delta),
    0, 0, 0, false,
    p.channel_id,
    coalesce(p.token_id, 0),
    p.user_group,
    jsonb_build_object(
        'is_task', true,
        'task_id', p.task_id,
        'billing_adjustment', 'oairegbox_duration_v1',
        'billed_seconds', p.billed_seconds,
        'actual_seconds', p.actual_seconds,
        'prior_quota', p.prior_quota,
        'actual_quota', p.actual_quota,
        'quota_delta', p.quota_delta
    )::text
FROM pending_duration_adjustments p
JOIN users u ON u.id = p.user_id
LEFT JOIN tokens tk ON tk.id = p.token_id
WHERE p.quota_delta <> 0;

-- 同步修正历史小时看板额度，不增加调用次数。
UPDATE quota_data q
SET quota = q.quota + p.quota_delta
FROM pending_duration_adjustments p
JOIN users u ON u.id = p.user_id
WHERE q.user_id = p.user_id
  AND q.username = u.username
  AND q.model_name = p.model_name
  AND q.created_at = p.created_at - (p.created_at % 3600)
  AND p.quota_delta <> 0;

INSERT INTO quota_data (user_id, username, model_name, created_at, token_used, count, quota)
SELECT p.user_id, u.username, p.model_name,
       p.created_at - (p.created_at % 3600), 0, 0, p.quota_delta
FROM pending_duration_adjustments p
JOIN users u ON u.id = p.user_id
WHERE p.quota_delta <> 0
  AND NOT EXISTS (
      SELECT 1 FROM quota_data q
      WHERE q.user_id = p.user_id
        AND q.username = u.username
        AND q.model_name = p.model_name
        AND q.created_at = p.created_at - (p.created_at % 3600)
  );

UPDATE tasks t
SET quota = p.actual_quota,
    private_data = jsonb_set(
        jsonb_set(
            t.private_data::jsonb,
            '{billing_adjustments}',
            coalesce(t.private_data::jsonb->'billing_adjustments', '{}'::jsonb),
            true
        ),
        '{billing_adjustments,oairegbox_duration_v1}',
        jsonb_build_object(
            'applied_at', extract(epoch FROM now())::bigint,
            'billed_seconds', p.billed_seconds,
            'actual_seconds', p.actual_seconds,
            'prior_quota', p.prior_quota,
            'actual_quota', p.actual_quota,
            'quota_delta', p.quota_delta
        ),
        true
    )::json
FROM pending_duration_adjustments p
WHERE t.id = p.task_id;

SELECT
    count(*) AS adjusted_tasks,
    sum(greatest(quota_delta, 0)) AS charged_quota,
    sum(greatest(-quota_delta, 0)) AS refunded_quota,
    sum(quota_delta) AS net_quota
FROM pending_duration_adjustments;

COMMIT;
