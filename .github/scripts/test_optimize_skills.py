#!/usr/bin/env python3
from __future__ import annotations

import shutil
import unittest
from contextlib import redirect_stdout
from io import StringIO
from pathlib import Path

import optimize_skills


class ParseJsonResponseTest(unittest.TestCase):
    def test_parses_plain_json(self) -> None:
        payload = optimize_skills.parse_json_response('{"skill_rewrites": [], "summary": "ok"}')

        self.assertEqual(payload["skill_rewrites"], [])
        self.assertEqual(payload["summary"], "ok")

    def test_parses_fenced_json(self) -> None:
        payload = optimize_skills.parse_json_response(
            '```json\n{"skill_rewrites": [], "summary": "ok"}\n```'
        )

        self.assertEqual(payload["skill_rewrites"], [])

    def test_prefers_payload_with_skill_rewrites(self) -> None:
        payload = optimize_skills.parse_json_response(
            'Example: {"ignored": true}\nFinal:\n{"skill_rewrites": [], "summary": "ok"}'
        )

        self.assertEqual(payload["summary"], "ok")

    def test_allows_trailing_commas(self) -> None:
        payload = optimize_skills.parse_json_response('{"skill_rewrites": [], "summary": "ok",}')

        self.assertEqual(payload["summary"], "ok")

    def test_rejects_non_json_response(self) -> None:
        with self.assertRaises(optimize_skills.SkillOptimizerError):
            optimize_skills.parse_json_response("not json at all")

    def test_apply_updates_rewrites_and_archives_skill(self) -> None:
        tmp = Path(__file__).resolve().parent / "_tmp_skill_test"
        if tmp.exists():
            shutil.rmtree(tmp)
        tmp.mkdir()
        try:
            skill_path = tmp / "SKILL.md"
            skill_path.write_text("# Skill: Old\n\nold\n", encoding="utf-8")
            original_target = optimize_skills.SKILL_TARGETS["backend-admin-service"]
            optimize_skills.SKILL_TARGETS["backend-admin-service"] = optimize_skills.SkillTarget(
                skill_path,
                "test skill",
            )
            try:
                with redirect_stdout(StringIO()):
                    applied = optimize_skills.apply_updates(
                        {
                            "skill_rewrites": [
                                {
                                    "skill_id": "backend-admin-service",
                                    "content": "# Skill: New\n\nnew\n",
                                }
                            ]
                        },
                        {"backend-admin-service"},
                        "1234567890abcdef",
                        dry_run=False,
                    )
            finally:
                optimize_skills.SKILL_TARGETS["backend-admin-service"] = original_target

            self.assertEqual(applied, 1)
            self.assertEqual(skill_path.read_text(encoding="utf-8"), "# Skill: New\n\nnew\n")
            self.assertEqual(
                (skill_path.parent / "archive" / "SKILL.1234567890ab.md").read_text(encoding="utf-8"),
                "# Skill: Old\n\nold\n",
            )
        finally:
            if tmp.exists():
                shutil.rmtree(tmp)


if __name__ == "__main__":
    unittest.main()
