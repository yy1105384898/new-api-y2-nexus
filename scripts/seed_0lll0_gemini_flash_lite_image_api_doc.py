#!/usr/bin/env python3
"""写入 0lll0-gemini-3.1-flash-lite-image 的 profile、api_doc（源站执行）。"""

from __future__ import annotations

import json
import subprocess
import time

ENDPOINTS = [
    {
        "method": "POST",
        "path": "{{base}}/chat/completions",
        "description": "同步文生图/图生图（stream 须为 false）。",
    }
]

PARAMS = [
    {"name": "model", "description": "必填，固定传 0lll0-gemini-3.1-flash-lite-image。"},
    {
        "name": "messages",
        "description": "必填，OpenAI 消息数组。文生图 user.content 为字符串；图生图/参考图为 content 数组（text + image_url）。",
    },
    {
        "name": "messages[].content[].type=text",
        "description": "提示词文本；可在 text 中用 @图1、@图2 引用同条消息中的参考图顺序。",
    },
    {
        "name": "messages[].content[].type=image_url",
        "description": "参考图：公网 URL 或 data:image/...;base64,...。多张参考图可追加多个 image_url 项。",
    },
    {
        "name": "mask",
        "description": "局部编辑蒙版：作为额外 image_url 传入 content（PNG，透明区域为编辑区），并在 text 中说明蒙版用途。",
    },
    {"name": "stream", "description": "须为 false（同步 JSON 响应，非 SSE）。"},
    {
        "name": "extra_body.google.image_config.aspect_ratio",
        "description": "画幅比例：1:1、2:3、3:2、3:4、4:3、4:5、5:4、9:16、16:9、21:9 等；auto 可省略（有参考图时可由上游推断）。",
    },
    {
        "name": "extra_body.google.image_config.image_size",
        "description": "仅支持 1K（K 大写）；省略或 auto 时上游默认 1K。勿传 2K/4K。",
    },
    {
        "name": "n",
        "description": "生成张数：客户端可对同一请求多次调用实现 1–4 张（每次返回 1 张图）。",
    },
]

CREATE_RESP = {
    "choices": [
        {
            "message": {
                "role": "assistant",
                "content": "![image](data:image/jpeg;base64,...)",
            }
        }
    ]
}

DOC = {
    "dispatch_mode": "sync",
    "intro": (
        "Gemini 3.1 Flash Lite 图像生成（0lll0 渠道）。按次 ¥0.20/张，仅支持 1K 出图。"
        "POST /v1/chat/completions（stream 须 false），图片嵌在 assistant message 的 Markdown data URI 中。"
        "画幅用 extra_body.google.image_config 的 aspect_ratio 控制；image_size 仅可 1K 或省略，勿传 2K/4K。"
        "带参考图/蒙版时 messages[0].content 为 text + image_url 数组（非 /images/generations）。"
    ),
    "endpoints": ENDPOINTS,
    "params": PARAMS,
    "basic_request_json": {
        "model": "{{model}}",
        "stream": False,
        "messages": [{"role": "user", "content": "一只橘猫坐在窗台上，水彩画风格，午后阳光"}],
        "extra_body": {
            "google": {
                "image_config": {
                    "aspect_ratio": "1:1",
                    "image_size": "1K",
                }
            }
        },
    },
    "request_json": {
        "model": "{{model}}",
        "stream": False,
        "messages": [
            {
                "role": "user",
                "content": [
                    {"type": "text", "text": "将 @图1 的风格应用到新场景"},
                    {"type": "image_url", "image_url": {"url": "https://example.com/ref.png"}},
                ],
            }
        ],
        "extra_body": {
            "google": {
                "image_config": {
                    "aspect_ratio": "3:2",
                    "image_size": "1K",
                }
            }
        },
    },
    "create_response_json": CREATE_RESP,
    "query_response_json": None,
}

MODEL_NAME = "0lll0-gemini-3.1-flash-lite-image"
IMAGE_PROFILE = "image-tpl-aspect-count-flash-lite"

PROFILE_PARAMS = {
    "quality": {
        "enabled": True,
        "options": [{"value": "auto", "label": "自动"}, {"value": "low", "label": "1K"}],
    },
    "aspectRatio": {
        "enabled": True,
        "options": [
            {"value": "1:1", "label": "1:1", "size": "1:1", "width": 1, "height": 1, "icon": "square"},
            {"value": "16:9", "label": "16:9", "size": "16:9", "width": 16, "height": 9, "icon": "landscape"},
            {"value": "9:16", "label": "9:16", "size": "9:16", "width": 9, "height": 16, "icon": "portrait"},
            {"value": "3:2", "label": "3:2", "size": "3:2", "width": 3, "height": 2, "icon": "landscape"},
            {"value": "2:3", "label": "2:3", "size": "2:3", "width": 2, "height": 3, "icon": "portrait"},
            {"value": "4:3", "label": "4:3", "size": "4:3", "width": 4, "height": 3, "icon": "landscape"},
            {"value": "3:4", "label": "3:4", "size": "3:4", "width": 3, "height": 4, "icon": "portrait"},
            {"value": "21:9", "label": "21:9", "size": "21:9", "width": 21, "height": 9, "icon": "landscape"},
            {"value": "auto", "label": "自动", "width": 0, "height": 0, "icon": "auto"},
        ],
    },
    "customDimensions": {"enabled": False},
    "count": {"enabled": True, "min": 1, "max": 4, "quickCount": 4},
    "background": {"enabled": False},
    "outputFormat": {"enabled": False},
    "outputCompression": {"enabled": False},
    "moderation": {"enabled": False},
}

PROFILE_HINTS = [
    {"text": "Flash Lite 图像模型仅支持 1K 出图（约 1024px），不支持 2K/4K。"},
    {"text": "请使用 1:1、16:9 等比例；quality 选 1K 或自动即可。"},
]


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


def upsert_profile() -> None:
    now = int(time.time())
    params_esc = json.dumps(PROFILE_PARAMS, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    hints_esc = json.dumps(PROFILE_HINTS, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    psql(
        f"""
        INSERT INTO model_ui_param_profiles (
            capability, profile_id, match, sort_order, api_mode,
            requires_reference_media, poll, reference_limits, params, option_rules, hints,
            created_time, updated_time
        ) VALUES (
            'image', '{IMAGE_PROFILE}', '["flash-lite-image"]', 0, 'chat-completions',
            false, '{{}}', '{{}}', '{params_esc}', '[]', '{hints_esc}',
            {now}, {now}
        )
        ON CONFLICT (capability, profile_id) DO UPDATE SET
            match = EXCLUDED.match,
            api_mode = EXCLUDED.api_mode,
            params = EXCLUDED.params,
            hints = EXCLUDED.hints,
            updated_time = EXCLUDED.updated_time;
        """
    )
    print(f"upserted profile {IMAGE_PROFILE}")


def main() -> None:
    upsert_profile()
    payload = DOC.copy()
    esc = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    psql(
        f"UPDATE models SET api_doc = '{esc}', image_profile_id = '{IMAGE_PROFILE}', "
        f"updated_time = extract(epoch from now())::bigint "
        f"WHERE model_name = '{MODEL_NAME}' AND deleted_at IS NULL;"
    )
    print(f"updated {MODEL_NAME}")
    psql(
        f"SELECT model_name, image_profile_id, length(api_doc) AS doc_len "
        f"FROM models WHERE model_name = '{MODEL_NAME}' AND deleted_at IS NULL;"
    )


if __name__ == "__main__":
    main()
