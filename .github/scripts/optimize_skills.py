#!/usr/bin/env python3
"""Optimize project skills from backend/front code changes.

The script intentionally keeps the provider boundary small: workflow plumbing,
diff collection, and skill-file updates are provider-agnostic; only the chat
completion call lives in the provider implementation.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import subprocess
import sys
import textwrap
import urllib.error
import urllib.request
from dataclasses import dataclass
from pathlib import Path
from typing import Any


REPO_ROOT = Path(__file__).resolve().parents[2]
SKILL_ROOT = REPO_ROOT / ".agents" / "skills"
MAX_DIFF_CHARS = int(os.getenv("AI_SKILL_MAX_DIFF_CHARS", "60000"))
OPTIMIZER_SECTION = "## 自动优化记录"


@dataclass(frozen=True)
class SkillTarget:
    path: Path
    description: str


SKILL_TARGETS = {
    "frontend-admin-react": SkillTarget(
        SKILL_ROOT / "frontend" / "admin-react" / "SKILL.md",
        "React 管理后台、前端路由、API、状态、主题、Mock 和测试",
    ),
    "backend-admin-service": SkillTarget(
        SKILL_ROOT / "backend" / "admin-service" / "SKILL.md",
        "GoFiber 管理后台服务、路由、业务逻辑、中间件、配置和 Swagger",
    ),
    "backend-go-common": SkillTarget(
        SKILL_ROOT / "backend" / "go-common" / "SKILL.md",
        "通用 Go 工具库：日志、配置、ID、加密、集合、转换、IP/Geo 等",
    ),
    "backend-orm-crud": SkillTarget(
        SKILL_ROOT / "backend" / "orm-crud" / "SKILL.md",
        "GORM CRUD、分页过滤、排序、proto/OpenAPI 辅助与 ORM 生成链路",
    ),
}


class SkillOptimizerError(RuntimeError):
    pass


class ChatProvider:
    name = "base"

    def complete(self, system_prompt: str, user_prompt: str) -> str:
        raise NotImplementedError


class DeepSeekProvider(ChatProvider):
    name = "deepseek"

    def __init__(self) -> None:
        self.api_key = os.getenv("DEEPSEEK_API_KEY", "")
        self.base_url = os.getenv("DEEPSEEK_BASE_URL", "https://api.deepseek.com").rstrip("/")
        self.model = os.getenv("DEEPSEEK_MODEL", os.getenv("AI_MODEL", "deepseek-v4-pro"))

    def complete(self, system_prompt: str, user_prompt: str) -> str:
        if not self.api_key:
            raise SkillOptimizerError(
                "DEEPSEEK_API_KEY is not set; configure the repository secret to enable AI skill optimization."
            )

        payload: dict[str, Any] = {
            "model": self.model,
            "messages": [
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_prompt},
            ],
            "temperature": 0.2,
            "stream": False,
        }
        if os.getenv("DEEPSEEK_THINKING", "").lower() in {"1", "true", "yes"}:
            payload["reasoning_effort"] = os.getenv("DEEPSEEK_REASONING_EFFORT", "high")
            payload["extra_body"] = {"thinking": {"type": "enabled"}}

        request = urllib.request.Request(
            f"{self.base_url}/chat/completions",
            data=json.dumps(payload).encode("utf-8"),
            headers={
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json",
            },
            method="POST",
        )
        try:
            with urllib.request.urlopen(request, timeout=90) as response:
                data = json.loads(response.read().decode("utf-8"))
        except urllib.error.HTTPError as exc:
            body = exc.read().decode("utf-8", errors="replace")
            raise SkillOptimizerError(f"DeepSeek API request failed: HTTP {exc.code}: {body}") from exc
        except urllib.error.URLError as exc:
            raise SkillOptimizerError(f"DeepSeek API request failed: {exc}") from exc

        try:
            return data["choices"][0]["message"]["content"]
        except (KeyError, IndexError, TypeError) as exc:
            raise SkillOptimizerError(f"DeepSeek API returned an unexpected payload: {data}") from exc


PROVIDERS = {
    DeepSeekProvider.name: DeepSeekProvider,
}


def run_git(args: list[str]) -> str:
    result = subprocess.run(
        ["git", *args],
        cwd=REPO_ROOT,
        check=True,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    return result.stdout


def normalize_base(base: str, head: str) -> str:
    if base and not re.fullmatch(r"0+", base):
        return base
    try:
        return run_git(["rev-parse", f"{head}^"]).strip()
    except subprocess.CalledProcessError:
        return run_git(["rev-list", "--max-parents=0", head]).splitlines()[0]


def changed_files(base: str, head: str) -> list[str]:
    output = run_git(["diff", "--name-only", "--diff-filter=ACMRT", base, head, "--", "backend", "front"])
    return [line.strip() for line in output.splitlines() if line.strip()]


def collect_diff(base: str, head: str, files: list[str]) -> str:
    if not files:
        return ""
    output = run_git(["diff", "--no-ext-diff", "--unified=80", base, head, "--", *files])
    if len(output) <= MAX_DIFF_CHARS:
        return output
    return output[:MAX_DIFF_CHARS] + "\n\n[diff truncated by AI_SKILL_MAX_DIFF_CHARS]\n"


def infer_skill_ids(files: list[str]) -> list[str]:
    ids: set[str] = set()
    for file in files:
        if file.startswith("front/"):
            ids.add("frontend-admin-react")
        if file.startswith("backend/admin/"):
            ids.add("backend-admin-service")
        elif file.startswith("backend/go-common/"):
            ids.add("backend-go-common")
        elif file.startswith("backend/orm-crud/"):
            ids.add("backend-orm-crud")
    return sorted(ids)


def read_skill_context(skill_ids: list[str]) -> str:
    parts: list[str] = []
    for skill_id in skill_ids:
        target = SKILL_TARGETS[skill_id]
        if target.path.exists():
            parts.append(
                f"### {skill_id}\nPath: {target.path.relative_to(REPO_ROOT)}\n"
                f"Scope: {target.description}\n\n{target.path.read_text(encoding='utf-8')[:8000]}"
            )
    return "\n\n".join(parts)


def build_prompts(files: list[str], diff: str, skill_ids: list[str]) -> tuple[str, str]:
    allowed = {
        skill_id: str(SKILL_TARGETS[skill_id].path.relative_to(REPO_ROOT))
        for skill_id in skill_ids
    }
    system_prompt = textwrap.dedent(
        """
        你是 template 仓库的技能维护助手。你的任务是根据 backend/front 的代码变更，
        提炼可长期复用的项目技能改进，而不是复述一次性提交内容。

        只返回严格 JSON，不要 Markdown 包裹，不要额外解释。JSON 结构：
        {
          "updates": [
            {
              "skill_id": "allowed skill id",
              "title": "不超过 24 个中文字符的标题",
              "content": "- 第一条\\n- 第二条"
            }
          ],
          "summary": "一句话说明本次是否学到了可沉淀内容"
        }

        规则：
        - 只使用用户提供的 diff 和现有技能内容，不要编造未出现的事实。
        - 只记录以后做类似任务会有帮助的流程、约定、命令、路径或风险。
        - 不要记录具体 commit SHA、作者、时间、临时 bug、一次性需求。
        - content 使用 1 到 5 条 Markdown bullet，中文为主，必要时保留路径和命令。
        - 如果没有可沉淀内容，返回 {"updates": [], "summary": "..."}。
        """
    ).strip()
    user_prompt = textwrap.dedent(
        f"""
        可更新的技能：
        {json.dumps(allowed, ensure_ascii=False, indent=2)}

        变更文件：
        {json.dumps(files, ensure_ascii=False, indent=2)}

        现有技能摘录：
        {read_skill_context(skill_ids)}

        Git diff：
        {diff}
        """
    ).strip()
    return system_prompt, user_prompt


def parse_json_response(content: str) -> dict[str, Any]:
    try:
        return json.loads(content)
    except json.JSONDecodeError:
        match = re.search(r"\{.*\}", content, re.DOTALL)
        if not match:
            raise SkillOptimizerError(f"AI response was not JSON: {content}")
        return json.loads(match.group(0))


def sanitize_title(value: Any) -> str:
    title = str(value or "代码变更经验").strip()
    return re.sub(r"[\r\n#]+", " ", title)[:40]


def sanitize_content(value: Any) -> str:
    content = str(value or "").strip()
    content = re.sub(r"<!--.*?-->", "", content, flags=re.DOTALL).strip()
    lines = [line.rstrip() for line in content.splitlines() if line.strip()]
    if not lines:
        return ""
    normalized: list[str] = []
    for line in lines[:8]:
        normalized.append(line if line.lstrip().startswith("- ") else f"- {line.lstrip('- ').strip()}")
    return "\n".join(normalized)


def append_update(skill_path: Path, marker: str, title: str, content: str) -> bool:
    text = skill_path.read_text(encoding="utf-8")
    if marker in text:
        return False
    entry = f"\n\n<!-- {marker} -->\n### {title}\n\n{content}\n"
    if OPTIMIZER_SECTION in text:
        text = text.rstrip() + entry
    else:
        text = text.rstrip() + f"\n\n{OPTIMIZER_SECTION}" + entry
    skill_path.write_text(text + "\n", encoding="utf-8")
    return True


def apply_updates(payload: dict[str, Any], allowed_ids: set[str], head: str, dry_run: bool) -> int:
    updates = payload.get("updates", [])
    if not isinstance(updates, list):
        raise SkillOptimizerError("AI response field 'updates' must be a list.")

    applied = 0
    for index, update in enumerate(updates, start=1):
        if not isinstance(update, dict):
            continue
        skill_id = str(update.get("skill_id", "")).strip()
        if skill_id not in allowed_ids:
            print(f"Skipping update with unsupported skill_id: {skill_id}", file=sys.stderr)
            continue
        content = sanitize_content(update.get("content"))
        if not content:
            continue
        title = sanitize_title(update.get("title"))
        marker = f"ai-skill-optimizer:{head[:12]}:{index}"
        skill_path = SKILL_TARGETS[skill_id].path
        if dry_run:
            print(f"\n--- {skill_path.relative_to(REPO_ROOT)} :: {title}\n{content}")
            applied += 1
            continue
        if append_update(skill_path, marker, title, content):
            applied += 1
    return applied


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--base", default=os.getenv("BASE_SHA", ""))
    parser.add_argument("--head", default=os.getenv("HEAD_SHA", "HEAD"))
    parser.add_argument("--provider", default=os.getenv("AI_PROVIDER", "deepseek"))
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    if args.provider not in PROVIDERS:
        names = ", ".join(sorted(PROVIDERS))
        raise SkillOptimizerError(f"Unsupported AI provider '{args.provider}'. Available providers: {names}")

    head = run_git(["rev-parse", args.head]).strip()
    base = normalize_base(args.base, head)
    files = changed_files(base, head)
    if not files:
        print("No backend/front changes detected; nothing to optimize.")
        return 0

    skill_ids = infer_skill_ids(files)
    if not skill_ids:
        print("Changes did not match a known skill target; nothing to optimize.")
        return 0

    diff = collect_diff(base, head, files)
    system_prompt, user_prompt = build_prompts(files, diff, skill_ids)

    try:
        provider = PROVIDERS[args.provider]()
        content = provider.complete(system_prompt, user_prompt)
    except SkillOptimizerError as exc:
        print(f"::warning::{exc}")
        return 0

    payload = parse_json_response(content)
    applied = apply_updates(payload, set(skill_ids), head, args.dry_run)
    summary = str(payload.get("summary", "")).strip()
    print(f"AI skill optimizer summary: {summary or 'no summary'}")
    print(f"Applied skill updates: {applied}")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except subprocess.CalledProcessError as exc:
        print(exc.stderr or str(exc), file=sys.stderr)
        raise SystemExit(exc.returncode)
    except SkillOptimizerError as exc:
        print(f"::error::{exc}", file=sys.stderr)
        raise SystemExit(1)
