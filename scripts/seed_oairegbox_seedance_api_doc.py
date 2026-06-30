#!/usr/bin/env python3
"""写入 OAIREGBox Seedance 满血三模型的 api_doc 与 video_profile_id（源站执行）。"""

from __future__ import annotations

import json
import subprocess

SHARED_PARAMS = [
    {"name": "model", "description": "必填，固定传当前模型别名。"},
    {"name": "prompt", "description": "必填。多素材时在 prompt 用 @Image1/@Video1/@Audio1 引用。"},
    {"name": "aspect_ratio", "description": "画幅：16:9（默认）、9:16、1:1、21:9、3:4、4:3。"},
    {"name": "duration", "description": "时长秒数，4–15 任意整数。"},
    {"name": "image_url", "description": "主参考图（公网 URL 或 data:image Base64）。"},
    {"name": "extra_images", "description": "额外参考图数组，与 image_url 合计 ≤9。"},
    {"name": "extra_videos", "description": "参考视频 ≤3（mp4/mov，单条 2–15s，24–60fps，≤50MB）。"},
    {"name": "extra_audios", "description": "参考音频 ≤3。"},
    {"name": "first_image_url", "description": "首尾帧：开始画面（须与 last_image_url 成对）。"},
    {"name": "last_image_url", "description": "首尾帧：结束画面。"},
]

CREATE_RESP = {
    "id": "task_01HZX8A2...",
    "status": "queued",
    "created_at": "2026-05-17T08:00:00Z",
}

QUERY_RESP = {
    "id": "task_01HZX8A2...",
    "status": "completed",
    "video_url": "https://example.com/output.mp4",
}

DOCS: dict[str, dict] = {
    "oairegbox-seedance-pro-720p": {
        "intro": "Seedance 2.0 满血 Pro 720p。按秒计费（¥0.65/s × duration）。支持 9 图 / 3 视频 / 3 音频全参考，输出 720p H.264 + AAC，无水印。POST /v1/videos 提交后轮询取片。",
        "basic_request_json": {
            "model": "oairegbox-seedance-pro-720p",
            "prompt": "雨夜霓虹街道，镜头缓慢推进，电影感光影",
            "aspect_ratio": "16:9",
            "duration": 8,
        },
        "request_json": {
            "model": "oairegbox-seedance-pro-720p",
            "prompt": "以 @Image1 的人物、@Video1 的运镜，配合 @Audio1 的节奏生成广告",
            "aspect_ratio": "16:9",
            "duration": 10,
            "image_url": "https://cdn.example.com/main.jpg",
            "extra_images": ["https://cdn.example.com/ref.jpg"],
            "extra_videos": ["https://cdn.example.com/ref.mp4"],
            "extra_audios": ["https://cdn.example.com/ref.mp3"],
        },
    },
    "oairegbox-seedance-fast-720p": {
        "intro": "Seedance 2.0 满血 Fast 720p。按秒计费（¥0.60/s × duration）。快速出片，参数与 Pro 720p 一致，支持全参考素材。",
        "basic_request_json": {
            "model": "oairegbox-seedance-fast-720p",
            "prompt": "雪山日出航拍，镜头平稳推进",
            "aspect_ratio": "16:9",
            "duration": 6,
        },
        "request_json": {
            "model": "oairegbox-seedance-fast-720p",
            "prompt": "让 @Image1 中的人物在 @Video1 的场景里行走",
            "aspect_ratio": "16:9",
            "duration": 8,
            "image_url": "https://cdn.example.com/person.jpg",
            "extra_videos": ["https://cdn.example.com/scene.mp4"],
        },
    },
    "oairegbox-seedance-pro-1080p": {
        "intro": "Seedance 2.0 满血 Pro 1080p。按秒计费（¥1.50/s × duration）。超清 1080p 输出，支持全参考素材，适合大屏/商用成片。",
        "basic_request_json": {
            "model": "oairegbox-seedance-pro-1080p",
            "prompt": "产品特写，柔光棚拍，缓慢环绕",
            "aspect_ratio": "16:9",
            "duration": 8,
        },
        "request_json": {
            "model": "oairegbox-seedance-pro-1080p",
            "prompt": "以 @Image1 的产品为主体，参考 @Video1 的运镜",
            "aspect_ratio": "16:9",
            "duration": 10,
            "image_url": "https://cdn.example.com/product.jpg",
            "extra_videos": ["https://cdn.example.com/motion-ref.mp4"],
        },
    },
}

PROFILE = "video-tpl-seedance-async"


def psql(sql: str) -> None:
    subprocess.run(
        ["docker", "exec", "newapi-postgres", "psql", "-U", "root", "-d", "new-api", "-v", "ON_ERROR_STOP=1", "-c", sql],
        check=True,
    )


def main() -> None:
    for model_name, slice_doc in DOCS.items():
        payload = {
            "dispatch_mode": "async",
            "intro": slice_doc["intro"],
            "params": SHARED_PARAMS,
            "basic_request_json": slice_doc["basic_request_json"],
            "request_json": slice_doc["request_json"],
            "create_response_json": CREATE_RESP,
            "query_response_json": QUERY_RESP,
        }
        esc = json.dumps(payload, ensure_ascii=False, separators=(",", ":")).replace("'", "''")
        psql(
            f"UPDATE models SET api_doc = '{esc}', video_profile_id = '{PROFILE}', "
            f"updated_time = extract(epoch from now())::bigint "
            f"WHERE model_name = '{model_name}' AND deleted_at IS NULL;"
        )
        print(f"updated {model_name}")

    psql(
        "SELECT model_name, video_profile_id, length(api_doc) AS doc_len "
        "FROM models WHERE model_name LIKE 'oairegbox-seedance-%' AND status=1 ORDER BY 1;"
    )


if __name__ == "__main__":
    main()
