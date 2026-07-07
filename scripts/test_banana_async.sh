#!/usr/bin/env bash
# Banana 异步出图 E2E 测试
# 用法:
#   export NEW_API_TOKEN='sk-你的key'
#   export NEW_API_BASE='https://vip-api.cangyuansuanli.cn'   # 或 http://127.0.0.1:3000
#   ./scripts/test_banana_async.sh
#
# contabo 上（从 DB 取 user 58 的 key）:
#   ssh contabo 'bash -s' < scripts/test_banana_async.sh
set -euo pipefail

BASE="${NEW_API_BASE:-http://127.0.0.1:3000}"
MODEL="${BANANA_MODEL:-gemini-banana-pro-4k}"

if [ -z "${NEW_API_TOKEN:-}" ]; then
  if command -v docker >/dev/null && docker ps --format '{{.Names}}' | grep -q newapi-postgres; then
    NEW_API_TOKEN=$(docker exec newapi-postgres psql -U root -d new-api -t -A -c \
      "SELECT key FROM tokens WHERE user_id=58 AND status=1 ORDER BY id DESC LIMIT 1;" | tr -d '[:space:]')
  fi
fi
if [ -z "${NEW_API_TOKEN:-}" ]; then
  echo "请设置 NEW_API_TOKEN=sk-..." >&2
  exit 1
fi

echo "BASE=$BASE MODEL=$MODEL TOKEN=${NEW_API_TOKEN:0:12}..."

BODY=$(cat <<EOF
{"model":"$MODEL","prompt":"a cute red apple on white background","size":"1024x1024","quality":"low","async":true,"stream":false,"n":1}
EOF
)

echo "=== POST async ==="
CREATE=$(curl -sS -w "\n__HTTP__:%{http_code} __TIME__:%{time_total}s" -X POST "$BASE/v1/images/generations" \
  -H "Authorization: Bearer $NEW_API_TOKEN" \
  -H "Content-Type: application/json" \
  -d "$BODY")
echo "$CREATE" | sed 's/__HTTP__/HTTP:/;s/__TIME__/TIME:/' | head -c 1500

TASK_ID=$(echo "$CREATE" | sed 's/__HTTP__.*//' | python3 -c "import sys,json; print(json.load(sys.stdin).get('id',''))" 2>/dev/null || true)
if [ -z "$TASK_ID" ]; then
  echo "NO task id — 创建失败" >&2
  exit 1
fi
echo "TASK_ID=$TASK_ID"

echo "=== POLL (每 5s，最多 3 分钟) ==="
for i in $(seq 1 36); do
  sleep 5
  POLL=$(curl -sS -w "\n__HTTP__:%{http_code}" "$BASE/v1/images/generations/$TASK_ID" \
    -H "Authorization: Bearer $NEW_API_TOKEN")
  BODY_ONLY=$(echo "$POLL" | sed 's/__HTTP__.*//')
  STATUS=$(echo "$BODY_ONLY" | python3 -c "import sys,json; print(json.load(sys.stdin).get('status',''))" 2>/dev/null || echo "?")
  echo "attempt=$i status=$STATUS"
  echo "$BODY_ONLY" | head -c 500
  echo
  case "$STATUS" in
    completed)
      echo "=== 出图成功 ==="
      echo "$BODY_ONLY" | python3 -c "
import sys,json
d=json.load(sys.stdin)
data=d.get('data') or []
urls=[item.get('url') for item in data if isinstance(item,dict) and item.get('url')]
print('model:', d.get('model'))
print('urls:', urls)
print('error:', d.get('error'))
"
      exit 0
      ;;
    failed)
      echo "=== 失败 ==="
      echo "$BODY_ONLY"
      exit 2
      ;;
  esac
done
echo "=== 3 分钟仍未完成 ==="
exit 3
