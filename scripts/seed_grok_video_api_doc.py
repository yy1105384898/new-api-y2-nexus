#!/usr/bin/env python3
"""写入 Grok 视频两模型的 models.api_doc（源站执行，对齐统一 /v1/videos 协议）。"""

from __future__ import annotations

import json
import subprocess

ENDPOINTS = [
    {"method": "POST", "path": "{{base}}/videos", "description": "通过统一视频任务接口创建任务（application/json 或 multipart/form-data）。"},
    {"method": "GET", "path": "{{base}}/videos/{task_id}", "description": "通过统一视频任务接口查询状态和结果。"},
    {"method": "GET", "path": "{{base}}/videos/{task_id}/content", "description": "下载已完成任务的成片。"},
]

PROMPT_PARAM = {
    "name": "prompt",
    "description": "必填，视频描述提示词；上限 4096 字符。",
}

CREATE_RESP = {
    "id": "task_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "model": "{{model}}",
    "object": "video",
    "status": "queued",
    "task_id": "task_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    "progress": 0,
    "created_at": 1780000000,
}

QUERY_RESP = {
    "id": "task_xxx",
    "object": "video",
    "model": "{{model}}",
    "status": "completed",
    "progress": 100,
    "created_at": 1780000000,
    "seconds": "6",
    "metadata": {"video_url": "{{base}}/videos/task_xxx/content"},
}

DOCS: dict[str, dict] = {
    "cy-gv1-grok-video": {
        "dispatch_mode": "async",
        "intro": (
            "Grok 异步视频：POST /v1/videos 提交，GET /v1/videos/{task_id} 轮询至 SUCCESS。"
            "支持文生、单图/多参考图生视频，以及 video_url 视频编辑（可不传 image_urls）。"
            "文生/单图最长 15 秒；多参考图最多 7 张且 seconds>10 时自动按 10 秒处理。"
        ),
        "params": [
            {"name": "model", "description": "必填，固定传模型广场展示名（如 cy-gv1-grok-video）。"},
            PROMPT_PARAM,
            {"name": "seconds", "description": "时长（秒），可选 4、6、8、10、12、15；默认建议 4。"},
            {"name": "duration", "description": "与 seconds 等价。"},
            {
                "name": "aspect_ratio",
                "description": "画幅比例：1:1、16:9、9:16、4:3、3:4、3:2、2:3；默认 16:9。",
            },
            {"name": "resolution", "description": "清晰度：720p 或 480p。"},
            {
                "name": "image_urls",
                "description": "参考图 HTTPS 直链或 data:image/...;base64,... 数组；最多 7 张。多图时 seconds 上限 10。",
            },
            {
                "name": "video_url",
                "description": "可选，源视频 HTTPS 直链，用于视频编辑；可不传 image_urls。",
            },
        ],
        "endpoints": ENDPOINTS,
        "basic_request_json": {
            "model": "{{model}}",
            "prompt": "A cinematic shot of a red sports car driving through rainy neon streets at night",
            "seconds": 6,
            "resolution": "720p",
            "aspect_ratio": "16:9",
        },
        "request_json": {
            "model": "{{model}}",
            "prompt": "Animate the product with a slow rotating camera",
            "seconds": 6,
            "image_urls": ["https://example.com/product.png"],
            "resolution": "720p",
            "aspect_ratio": "9:16",
        },
        "examples": [
            {
                "title": "视频编辑",
                "request_json": {
                    "model": "{{model}}",
                    "prompt": "Add subtle rainbow in the sky",
                    "seconds": 4,
                    "resolution": "720p",
                    "aspect_ratio": "16:9",
                    "video_url": "https://example.com/source.mp4",
                },
            }
        ],
        "create_response_json": CREATE_RESP,
        "query_response_json": QUERY_RESP,
    },
    "cy-gv1-grok-video-1.5": {
        "dispatch_mode": "async",
        "intro": (
            "Grok 1.5 单图生视频：POST /v1/videos 提交，GET /v1/videos/{task_id} 轮询至 SUCCESS。"
            "必须且只能 1 张图片参考（image_urls / image），不支持纯文生、不支持视频参考；"
            "画幅仅 16:9 / 9:16；清晰度 480p/720p。"
        ),
        "params": [
            {"name": "model", "description": "必填，固定传 cy-gv1-grok-video-1.5。"},
            PROMPT_PARAM,
            {"name": "seconds", "description": "时长（秒），可选 4、6、8、10、12、15。"},
            {"name": "duration", "description": "与 seconds 等价。"},
            {"name": "aspect_ratio", "description": "仅支持 16:9 或 9:16。"},
            {"name": "resolution", "description": "清晰度：720p 或 480p。"},
            {
                "name": "image_urls",
                "description": "必填且只能 1 张参考图；HTTPS 直链或 data URL。勿传视频 URL。",
            },
            {
                "name": "image",
                "description": "与 image_urls 单图等价，格式：{\"url\": \"https://...\"}。",
            },
        ],
        "endpoints": ENDPOINTS,
        "basic_request_json": {
            "model": "{{model}}",
            "prompt": "Gentle camera push-in, water flowing",
            "seconds": 4,
            "image_urls": ["https://example.com/product.png"],
            "resolution": "720p",
            "aspect_ratio": "16:9",
        },
        "request_json": {
            "model": "{{model}}",
            "prompt": "Gentle camera push-in, water flowing",
            "duration": 4,
            "image": {"url": "https://example.com/product.png"},
            "resolution": "720p",
            "aspect_ratio": "9:16",
        },
        "create_response_json": CREATE_RESP,
        "query_response_json": QUERY_RESP,
    },
}


def psql(sql: str) -> None:
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


def main() -> None:
    for model_name, doc in DOCS.items():
        payload = json.dumps(doc, ensure_ascii=False).replace("'", "''")
        psql(f"UPDATE models SET api_doc = '{payload}' WHERE model_name = '{model_name}';")
        print(f"updated {model_name}")


if __name__ == "__main__":
    main()
