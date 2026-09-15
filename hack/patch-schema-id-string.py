#!/usr/bin/env python3
"""Patch svalabs forgejo schema: numeric TF id → string for upjet.

Upjet initializes empty external-name as \"\" before Create; terraform-plugin-sdk
cannot parse \"\" as number (EOF / no digits). Treating id as string keeps
IdentifierFromProvider working while Forgejo still returns numeric ids.
Re-run after `make generate` refreshes config/schema.json from TF.
"""
from __future__ import annotations

import json
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = ROOT / "config" / "schema.json"

RESOURCES = [
    "forgejo_organization",
    "forgejo_repository",
    "forgejo_team",
    "forgejo_user",
    "forgejo_gpg_key",
    "forgejo_personal_access_token",
]

# Resources without TF `id` that expose a numeric alt identifier.
ALT_ID_ATTRS = {
    "forgejo_deploy_key": "key_id",
    "forgejo_repository_webhook": "webhook_id",
}


def main() -> int:
    s = json.loads(SCHEMA.read_text())
    if "provider_schemas" in s:
        res = next(iter(s["provider_schemas"].values()))["resource_schemas"]
    else:
        res = s["resource_schemas"]
    changed = []
    for name in RESOURCES:
        attrs = res[name]["block"]["attributes"]
        if attrs.get("id", {}).get("type") == "number":
            attrs["id"]["type"] = "string"
            changed.append(name)
    alt_changed = []
    for name, attr in ALT_ID_ATTRS.items():
        attrs = res[name]["block"]["attributes"]
        if attrs.get(attr, {}).get("type") == "number":
            attrs[attr]["type"] = "string"
            alt_changed.append(f"{name}.{attr}")
    SCHEMA.write_text(json.dumps(s))
    print("patched id number→string:", changed or "(already string)")
    print("patched alt id number→string:", alt_changed or "(already string)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
