#!/usr/bin/env python3
"""Manju Gemini Banana 系列：api_doc（统一 Image API sync/async）+ ModelPrice（源站 docker 内执行）。"""

from __future__ import annotations

import argparse
import json
import subprocess
import time

CREATE_ASYNC = {
    "id": "task_img_01HZX8A2...",
    "object": "image.generation",
    "model": "{{model}}",
    "status": "queued",
    "progress": "10%",
    "created_at": 1715923200,
}
QUERY_ASYNC = {
    "id": "task_img_01HZX8A2...",
    "object": "image.generation",
    "status": "completed",
    "progress": "100%",
    "data": [{"url": "https://example.com/image.png"}],
}
CREATE_SYNC = {
    "created": 1715923200,
    "data": [{"url": "https://example.com/image.png"}],
}

REFERENCE_PARAMS = [
    {"name": "image", "description": "图生图：单张参考图 URL 或 data URI。"},
    {"name": "images", "description": "图生图：多张参考图 URL 数组。"},
    {"name": "mask", "description": "局部编辑蒙版 URL 或 data URI（PNG）。"},
]

COMMON_PARAMS_TAIL = [
    {"name": "stream", "description": "须为 false。"},
    {"name": "n", "description": "生成张数：客户端可对同一请求多次调用实现 1–15 张（每次返回 1 张图）。"},
]

ASYNC_ENDPOINTS = [
    {
        "method": "POST",
        "path": "{{base}}/images/generations",
        "description": "异步文生图/图生图（application/json，async 必须为 true，stream 须 false）。",
    },
    {
        "method": "POST",
        "path": "{{base}}/images/edits",
        "description": "图生图 multipart（参考图/蒙版）；与 JSON 传 image/images 等价，均走 Image API。",
    },
    {
        "method": "GET",
        "path": "{{base}}/images/generations/{task_id}",
        "description": "查询异步任务状态与结果 URL。",
    },
    {
        "method": "GET",
        "path": "{{base}}/images/{task_id}/content",
        "description": "下载任务图片（R2 代理地址）。",
    },
]

SYNC_ENDPOINTS = [
    {
        "method": "POST",
        "path": "{{base}}/images/generations",
        "description": "同步文生图（勿传 async，stream 须 false）。",
    },
    {
        "method": "POST",
        "path": "{{base}}/images/edits",
        "description": "图生图 multipart（参考图/蒙版）；与 JSON 传 image/images 等价。",
    },
]

MODELS: list[dict] = [
    {
        "internal": "manju-gemini-banana-pro-4k",
        "public": "gemini-banana-pro-4k",
        "price": 0.18,
        "image_profile_id": "image-tpl-banana-chat",
        "intro_sync": (
            "Gemini Banana Pro 4K 同步（推荐）：POST /v1/images/generations（勿传 async，stream: false），"
            "响应 data[].url 或 b64_json；图生图 JSON 传 image/images，或 multipart POST /v1/images/edits。"
        ),
        "intro_async": (
            "可选异步：POST /v1/images/generations 带 async: true，"
            "GET /v1/images/generations/{task_id} 轮询取 data[].url；4K/多参考图场景推荐。"
        ),
        "basic_size": "4K",
        "full_size": "4K",
        "size_desc": "分辨率档位：1K / 2K / 4K（K 大写；本模型推荐 4K）。",
    },
    {
        "internal": "manju-gemini-banana-flash-lite",
        "public": "gemini-banana-flash-lite",
        "price": 0.075,
        "image_profile_id": "image-tpl-banana-chat-flash-lite",
        "intro_sync": (
            "Gemini Banana Flash Lite 同步（推荐）：POST /v1/images/generations（stream: false），"
            "仅支持 1K；图生图 JSON 传 image/images。"
        ),
        "intro_async": (
            "可选异步：POST /v1/images/generations 带 async: true 后轮询；仅支持 1K。"
        ),
        "basic_size": "1K",
        "full_size": "1K",
        "size_desc": "仅支持 1K（K 大写；省略默认 1K）。勿传 2K/4K。",
    },
    {
        "internal": "manju-gemini-banana-pro-1/2k",
        "public": "gemini-banana-pro-1/2k",
        "price": 0.12,
        "image_profile_id": "image-tpl-banana-chat",
        "intro_sync": (
            "Gemini Banana Pro 1K/2K 同步（推荐）：POST /v1/images/generations（stream: false）；"
            "quality 省略默认 1K；图生图 JSON 传 image/images。"
        ),
        "intro_async": (
            "可选异步：POST /v1/images/generations 带 async: true，GET 轮询取图。"
        ),
        "basic_size": "2K",
        "full_size": "2K",
        "size_desc": "分辨率档位：1K / 2K（K 大写；省略默认 1K）。",
    },
    {
        "internal": "manju-gemini-banana-2.0-1/2k",
        "public": "gemini-banana-2.0-1/2k",
        "price": 0.081,
        "image_profile_id": "image-tpl-banana-chat",
        "intro_sync": (
            "Nano Banana 2.0 1K/2K 同步（推荐）：POST /v1/images/generations（stream: false）；"
            "quality 省略默认 1K；图生图 JSON 传 image/images。"
        ),
        "intro_async": (
            "可选异步：带 async: true 后轮询；多参考图/大图场景可选。"
        ),
        "basic_size": "1K",
        "full_size": "2K",
        "size_desc": "分辨率档位：1K / 2K（K 大写；省略默认 1K）。",
    },
    {
        "internal": "manju-gemini-banana-2.0-4k",
        "public": "gemini-banana-2.0-4k",
        "price": 0.135,
        "image_profile_id": "image-tpl-banana-chat",
        "intro_sync": (
            "Nano Banana 2.0 4K 同步（推荐）：POST /v1/images/generations（stream: false）；"
            "图生图 JSON 传 image/images，或 multipart /v1/images/edits。"
        ),
        "intro_async": (
            "可选异步：带 async: true 后轮询；4K/多参考图场景推荐。"
        ),
        "basic_size": "1K",
        "full_size": "4K",
        "size_desc": "分辨率档位：1K / 2K / 4K（K 大写；本模型推荐 4K）。",
    },
]


def build_params(public: str, size_desc: str, *, async_mode: bool) -> list[dict]:
    params = [
        {
            "name": "model",
            "description": f"必填，固定传模型广场展示名（{public}）。",
        },
        {
            "name": "prompt",
            "description": "必填，文生图提示词；图生图时在 prompt 中说明 @图片1 等引用顺序。",
        },
        {
            "name": "size",
            "description": "兼容旧 OpenAI 尺寸/比例；仅用于推断 aspect_ratio/output_resolution，不推荐与 aspect_ratio 混用。",
        },
        {
            "name": "aspect_ratio",
            "description": "推荐传画幅比例：1:1、16:9、9:16、4:3、3:4；auto 可省略并由参考图推断。",
        },
        {
            "name": "output_resolution",
            "description": f"{size_desc} Adobe2API / Manju 均会读取该字段。",
        },
        {
            "name": "image_size",
            "description": "output_resolution 的兼容别名；若同时传入，须与 output_resolution 保持一致。",
        },
        {
            "name": "quality",
            "description": "OpenAI 风格别名：low=1K、medium=2K、high=4K；推荐直接传 output_resolution。",
        },
    ] + REFERENCE_PARAMS + COMMON_PARAMS_TAIL
    if async_mode:
        params.append({"name": "async", "description": "异步模式必填 true。"})
    return params


def build_basic_request(public: str, basic_size: str, *, async_mode: bool) -> dict:
    body = {
        "model": public,
        "stream": False,
        "prompt": "一只橘猫坐在窗台上，午后阳光",
        "aspect_ratio": "1:1",
        "output_resolution": basic_size,
        "image_size": basic_size,
    }
    if async_mode:
        body["async"] = True
    return body


def build_full_request(public: str, full_size: str, *, async_mode: bool) -> dict:
    body = {
        "model": public,
        "stream": False,
        "prompt": "将 @图片1 的风格应用到新场景",
        "aspect_ratio": "16:9",
        "output_resolution": full_size,
        "image_size": full_size,
        "image": "https://example.com/ref.png",
    }
    if async_mode:
        body["async"] = True
    return body


def build_api_doc(spec: dict) -> dict:
    public = spec["public"]
    return {
        "modes": {
            "async": {
                "dispatch_mode": "async",
                "intro": spec["intro_async"],
                "endpoints": ASYNC_ENDPOINTS,
                "basic_request_json": build_basic_request(public, spec["basic_size"], async_mode=True),
                "request_json": build_full_request(public, spec["full_size"], async_mode=True),
                "params": build_params(public, spec["size_desc"], async_mode=True),
                "create_response_json": CREATE_ASYNC,
                "query_response_json": QUERY_ASYNC,
            },
            "sync": {
                "dispatch_mode": "sync",
                "intro": spec["intro_sync"],
                "endpoints": SYNC_ENDPOINTS,
                "basic_request_json": build_basic_request(public, spec["basic_size"], async_mode=False),
                "request_json": build_full_request(public, spec["full_size"], async_mode=False),
                "params": build_params(public, spec["size_desc"], async_mode=False),
                "create_response_json": CREATE_SYNC,
            },
        }
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
    parser = argparse.ArgumentParser(description="更新 Manju Gemini Banana 模型文档与价格")
    parser.add_argument(
        "--docs-only",
        action="store_true",
        help="只更新 api_doc/profile，不修改 ModelPrice",
    )
    args = parser.parse_args()

    now = int(time.time())
    price_updates: dict[str, float] = {}

    for spec in MODELS:
        doc = build_api_doc(spec)
        esc = json.dumps(doc, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
        internal = spec["internal"].replace("'", "''")
        profile = spec["image_profile_id"]
        psql_exec(
            f"UPDATE models SET api_doc = '{esc}', "
            f"image_profile_id = '{profile}', "
            f"updated_time = {now} "
            f"WHERE model_name = '{internal}' AND deleted_at IS NULL;"
        )
        price_updates[spec["internal"]] = spec["price"]
        print(f"api_doc updated: {spec['internal']} -> public {spec['public']}")

    if args.docs_only:
        print("ModelPrice unchanged (--docs-only)")
    else:
        merge_model_price(price_updates)
        print("ModelPrice updated:")
        for k, v in price_updates.items():
            print(f"  {k}: {v}")

    psql_exec(
        "SELECT model_name, image_profile_id, length(api_doc) AS doc_len "
        "FROM models WHERE model_name LIKE 'manju-gemini-banana%' AND deleted_at IS NULL ORDER BY 1;"
    )


if __name__ == "__main__":
    main()
