#!/usr/bin/env python3
"""Audit enabled abilities models for public-name prefix stripping collisions."""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from collections import defaultdict


def strip_channel_prefix(name: str, prefixes: tuple[str, ...]) -> str:
    trimmed = name.strip()
    for prefix in prefixes:
        if trimmed.startswith(prefix):
            return trimmed[len(prefix):]
    return trimmed


def run_sql(container: str, sql: str) -> list[str]:
    cmd = [
        "docker",
        "exec",
        container,
        "sqlite3",
        "-batch",
        "/data/one-api.db",
        sql,
    ]
    result = subprocess.run(cmd, capture_output=True, text=True, check=False)
    if result.returncode != 0:
        raise RuntimeError(result.stderr.strip() or result.stdout.strip())
    lines = [line.strip() for line in result.stdout.splitlines() if line.strip()]
    return lines


def load_prefixes(container: str | None, prefixes_file: str | None) -> tuple[str, ...]:
    if prefixes_file:
        with open(prefixes_file, encoding="utf-8") as fh:
            lines = [line.strip() for line in fh if line.strip()]
    elif container:
        sql = (
            "SELECT prefix FROM model_channel_prefixes "
            "WHERE enabled=1 AND deleted_at IS NULL "
            "ORDER BY sort_order, prefix;"
        )
        lines = run_sql(container, sql)
    else:
        raise RuntimeError("provide --container or --prefixes-file")

    normalized: list[str] = []
    for line in lines:
        prefix = line.strip()
        if prefix and not prefix.endswith("-"):
            prefix += "-"
        if prefix:
            normalized.append(prefix)
    return tuple(normalized)


def audit_models(models: list[str], prefixes: tuple[str, ...]) -> dict:
    public_to_internal: dict[str, list[str]] = defaultdict(list)
    for model in models:
        public = strip_channel_prefix(model, prefixes)
        public_to_internal[public].append(model)

    collisions = {
        public: sorted(internals)
        for public, internals in public_to_internal.items()
        if len(internals) > 1
    }
    return {
        "total_models": len(models),
        "unique_public": len(public_to_internal),
        "prefix_count": len(prefixes),
        "collisions": collisions,
        "mapping": {
            internal: strip_channel_prefix(internal, prefixes)
            for internal in sorted(models)
        },
    }


def main() -> int:
    parser = argparse.ArgumentParser(description="Audit model public-name collisions")
    parser.add_argument(
        "--container",
        default="new-api",
        help="Docker container name hosting new-api sqlite DB",
    )
    parser.add_argument(
        "--models-file",
        help="Plain text file with one model name per line (skip docker for models)",
    )
    parser.add_argument(
        "--prefixes-file",
        help="Plain text file with one prefix per line (skip docker for prefixes)",
    )
    parser.add_argument("--json", action="store_true", help="Print JSON report")
    args = parser.parse_args()

    if args.models_file:
        with open(args.models_file, encoding="utf-8") as fh:
            models = [line.strip() for line in fh if line.strip()]
    else:
        sql = "SELECT DISTINCT model FROM abilities WHERE enabled=1 ORDER BY model;"
        models = run_sql(args.container, sql)

    prefixes = load_prefixes(
        None if args.prefixes_file else args.container,
        args.prefixes_file,
    )
    if not prefixes:
        print("warning: no enabled channel prefixes configured", file=sys.stderr)

    report = audit_models(models, prefixes)
    if args.json:
        print(json.dumps(report, indent=2, ensure_ascii=False))
    else:
        print(f"enabled models: {report['total_models']}")
        print(f"channel prefixes: {report['prefix_count']}")
        print(f"unique public names: {report['unique_public']}")
        print(f"collisions: {len(report['collisions'])}")
        for public, internals in sorted(report["collisions"].items()):
            print(f"  {public}: {', '.join(internals)}")

    return 1 if report["collisions"] else 0


if __name__ == "__main__":
    sys.exit(main())
