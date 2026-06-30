#!/usr/bin/env python3
"""补全 Gemini Banana 系列 models.api_doc 与 image_profile_id（源站执行）。"""

from __future__ import annotations

import json
import subprocess

BANANA_API_DOC = {
    "dispatch_mode": "sync",
    "intro": (
        "Gemini Nano Banana 同步出图：POST /v1/chat/completions（stream 须 false），"
        "图片嵌在 assistant message 的 Markdown data URI 中。"
        "画幅用 extra_body.google.image_config 的 aspect_ratio + image_size 控制，勿传 async，无需轮询。"
        "带参考图/蒙版时 messages[0].content 为 text + image_url 数组（非 /images/edits）。"
    ),
    "endpoints": [
        {
            "method": "POST",
            "path": "{{base}}/chat/completions",
            "description": "同步文生图/图生图（stream 须为 false）。",
        }
    ],
    "basic_request_json": {
        "model": "{{model}}",
        "stream": False,
        "messages": [{"role": "user", "content": "一只橘猫坐在窗台上，午后阳光"}],
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
                    {"type": "text", "text": "将 @图片1 的风格应用到新场景"},
                    {
                        "type": "image_url",
                        "image_url": {"url": "https://example.com/ref.png"},
                    },
                ],
            }
        ],
        "extra_body": {
            "google": {
                "image_config": {
                    "aspect_ratio": "3:2",
                    "image_size": "4K",
                }
            }
        },
    },
    "params": [
        {"name": "model", "description": "必填，固定传模型广场展示名（如 gemini-banana-2.0 / gemini-banana-2.0-pro）。"},
        {
            "name": "messages",
            "description": "必填，OpenAI 消息数组。文生图 user.content 为字符串；图生图/参考图为 content 数组（text + image_url）。",
        },
        {
            "name": "messages[].content[].type=text",
            "description": "提示词文本；可在 text 中用 @图片1、@图片2 引用同条消息中的参考图顺序。",
        },
        {
            "name": "messages[].content[].type=image_url",
            "description": "参考图：公网 URL 或 data:image/...;base64,... 。多张参考图可追加多个 image_url 项。",
        },
        {
            "name": "mask",
            "description": "局部编辑蒙版：作为额外 image_url 传入 content（PNG，透明区域为编辑区），并在 text 中说明蒙版用途。",
        },
        {"name": "stream", "description": "须为 false（Banana 走 JSON 同步响应，非 SSE）。"},
        {
            "name": "extra_body.google.image_config.aspect_ratio",
            "description": "画幅比例：1:1、2:3、3:2、3:4、4:3、4:5、5:4、9:16、16:9、21:9 等；auto 可省略（有参考图时可由上游推断）。",
        },
        {
            "name": "extra_body.google.image_config.image_size",
            "description": "分辨率档位：1K / 2K / 4K（K 大写；省略默认 1K）。",
        },
        {
            "name": "n",
            "description": "生成张数：客户端可对同一请求多次调用实现 1–15 张（每次返回 1 张图）。",
        },
    ],
    "create_response_json": {
        "choices": [
            {
                "message": {
                    "role": "assistant",
                    "content": "![image](data:image/jpeg;base64,...)",
                }
            }
        ]
    },
    "query_response_json": None,
}

IMAGE_PROFILE = "image-tpl-aspect-count-extended"


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
    esc = json.dumps(BANANA_API_DOC, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
    psql(
        f"UPDATE models SET api_doc = '{esc}', "
        f"image_profile_id = '{IMAGE_PROFILE}', "
        f"updated_time = extract(epoch from now())::bigint "
        f"WHERE model_name ILIKE '%banana%' AND deleted_at IS NULL;"
    )
    print("updated banana models api_doc + image_profile_id")

    psql(
        "SELECT model_name, image_profile_id, length(api_doc) AS doc_len, "
        "CASE WHEN api_doc LIKE '%image_url%' THEN 'yes' ELSE 'no' END AS has_ref "
        "FROM models WHERE model_name ILIKE '%banana%' AND deleted_at IS NULL ORDER BY 1;"
    )


if __name__ == "__main__":
    main()
