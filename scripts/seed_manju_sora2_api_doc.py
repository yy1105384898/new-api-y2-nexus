#!/usr/bin/env python3
"""Manju OpenAI Sora2（manju-openai-sora2）：api_doc + ModelPrice（源站 docker 内执行）。

对齐 Manju Apifox：https://ssnsuyettr.apifox.cn/
- 创建：POST /v1/chat/completions（sora2_duration / sora2_ratio / messages / stream:false）
- 查询：GET /v1/videos/{task_id}
"""

from __future__ import annotations

import json
import subprocess
import time

MODEL = "manju-openai-sora2"
UPSTREAM = "sora2"
PROFILE = "video-tpl-manju-sora-async"
PRICE_PER_SECOND = 0.40

ENDPOINTS = [
    {
        "method": "POST",
        "path": "{{base}}/chat/completions",
        "description": "创建 Sora2 视频任务（stream 须 false；返回 task 对象含 poll_url）。",
    },
    {
        "method": "GET",
        "path": "{{base}}/videos/{task_id}",
        "description": "查询任务状态；完成时返回 raw_data.video_url / video.url。",
    },
    {
        "method": "GET",
        "path": "{{base}}/tasks/{task_id}",
        "description": "查询任务结果别名接口。",
    },
]

PARAMS = [
    {"name": "model", "description": f"必填，固定传 {MODEL}（映射上游 {UPSTREAM}）。"},
    {"name": "stream", "description": "须为 false。"},
    {"name": "messages", "description": "必填，messages[0].content 为视频提示词。"},
    {"name": "prompt", "description": "（兼容）与 messages 二选一，平台会转为 messages。"},
    {"name": "sora2_duration", "description": "时长（秒），如 8、12；与 seconds/duration 等价。"},
    {"name": "sora2_ratio", "description": "画幅比例：16:9、9:16、1:1；与 size/aspect_ratio 等价。"},
    {"name": "input_reference", "description": "图生视频：单张参考图 URL。"},
    {"name": "seconds", "description": "（兼容 /v1/videos）时长秒数。"},
    {"name": "size", "description": "（兼容 /v1/videos）1280x720、720x1280 等。"},
]

CREATE_RESP = {
    "id": "sora2-fb22482c0bde",
    "task_id": "sora2-fb22482c0bde",
    "platform": "sora2",
    "status": "running",
    "progress": 0,
    "poll_url": "https://manjuapi.com/v1/videos/sora2-fb22482c0bde",
    "properties": {"duration": "8", "aspect_ratio": "16:9", "output_resolution": "720p"},
    "object": "task",
}

QUERY_RESP = {
    "id": "sora2-fb22482c0bde",
    "platform": "sora2",
    "status": "succeeded",
    "progress": 100,
    "raw_data": {"video_url": "https://example.com/output.mp4"},
    "video": {"url": "https://example.com/output.mp4"},
}


def build_api_doc() -> dict:
    return {
        "dispatch_mode": "async",
        "intro": (
            "Manju Sora2：创建走 POST /v1/chat/completions（stream:false），"
            "响应为 task 对象；轮询 GET /v1/videos/{task_id} 取片。"
        ),
        "endpoints": ENDPOINTS,
        "basic_request_json": {
            "model": MODEL,
            "stream": False,
            "sora2_duration": "8",
            "sora2_ratio": "16:9",
            "messages": [{"role": "user", "content": "雨夜城市街道，电影感镜头缓慢推进"}],
        },
        "request_json": {
            "model": MODEL,
            "stream": False,
            "sora2_duration": "12",
            "sora2_ratio": "16:9",
            "input_reference": "https://example.com/ref.png",
            "messages": [{"role": "user", "content": "主体缓慢向镜头靠近，电影感运镜"}],
        },
        "params": PARAMS,
        "create_response_json": CREATE_RESP,
        "query_response_json": QUERY_RESP,
    }


def psql(sql: str) -> str:
    result = subprocess.run(
        [
            "docker",
            "exec",
            "newapi-postgres",
            "psql",
            "-U",
            "root",
            "-d",
            "new-api",
            "-t",
            "-A",
            "-c",
            sql,
        ],
        check=True,
        capture_output=True,
        text=True,
    )
    return result.stdout.strip()


def psql_exec(sql: str) -> None:
    subprocess.run(
        [
            "docker",
            "exec",
            "newapi-postgres",
            "psql",
            "-U",
            "root",
            "-d",
            "new-api",
            "-v",
            "ON_ERROR_STOP=1",
            "-c",
            sql,
        ],
        check=True,
    )


def merge_model_price(updates: dict[str, float]) -> None:
    raw = psql("SELECT value::text FROM options WHERE key='ModelPrice'")
    prices = json.loads(raw)
    prices.update(updates)
    payload = json.dumps(prices, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    psql_exec(f"UPDATE options SET value='{payload}' WHERE key='ModelPrice'")


def main() -> None:
    now = int(time.time())
    doc = build_api_doc()
    esc = json.dumps(doc, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    endpoints = '{"openai-chat-video": {"path": "/v1/chat/completions", "method": "POST"}, "openai-video": {"path": "/v1/videos", "method": "GET"}}'
    psql_exec(
        f"UPDATE models SET api_doc = '{esc}', "
        f"video_profile_id = '{PROFILE}', "
        f"endpoints = '{endpoints}', "
        f"updated_time = {now} "
        f"WHERE model_name = '{MODEL}' AND deleted_at IS NULL;"
    )
    merge_model_price({MODEL: PRICE_PER_SECOND})
    print(f"api_doc + ModelPrice updated: {MODEL} @ {PRICE_PER_SECOND} USD/s")
    psql_exec(
        "SELECT model_name, video_profile_id, length(api_doc) AS doc_len "
        f"FROM models WHERE model_name = '{MODEL}' AND deleted_at IS NULL;"
    )


if __name__ == "__main__":
    main()
