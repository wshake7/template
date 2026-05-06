BEGIN;

ALTER TABLE "public"."sys_data_permission"
    ADD COLUMN IF NOT EXISTS "action_key" varchar(64) NOT NULL DEFAULT 'read';

UPDATE "public"."sys_data_permission"
SET "action" = '["all"]'::jsonb,
    "action_key" = 'all'
WHERE "action" ? 'all';

UPDATE "public"."sys_data_permission"
SET "action" = to_jsonb(
        ARRAY(
            SELECT action_value
            FROM unnest(ARRAY['read', 'write', 'delete']) AS action_value
            WHERE "sys_data_permission"."action" ? action_value
        )
    ),
    "action_key" = array_to_string(
        ARRAY(
            SELECT action_value
            FROM unnest(ARRAY['read', 'write', 'delete']) AS action_value
            WHERE "sys_data_permission"."action" ? action_value
        ),
        ','
    )
WHERE NOT ("action" ? 'all');

UPDATE "public"."sys_data_permission"
SET "action" = '["read"]'::jsonb,
    "action_key" = 'read'
WHERE jsonb_typeof("action") <> 'array'
   OR jsonb_array_length("action") = 0
   OR "action_key" = '';

ALTER TABLE "public"."sys_data_permission"
    ALTER COLUMN "scope_type" SET DEFAULT 'none',
    ALTER COLUMN "scope_field" SET DEFAULT 'id';

ALTER TABLE "public"."sys_user"
    ALTER COLUMN "language_code" SET DEFAULT '';

ALTER TABLE "public"."sys_api_log"
    ALTER COLUMN "referer" TYPE text,
    ALTER COLUMN "request_uri" TYPE text;

DROP INDEX IF EXISTS "public"."idx_sys_user_username_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_casbin_model_name";
DROP INDEX IF EXISTS "public"."idx_sys_data_permission_subject_resource_action_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_dict_type_type_code_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_language_type_code_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_language_entry_code_language_type_id_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_resource_code_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_role_code_deleted_at";
DROP INDEX IF EXISTS "public"."idx_sys_user_role_user_id_role_id_delete_at";

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_user_username_active"
    ON "public"."sys_user" ("username")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_casbin_model_name_active"
    ON "public"."sys_casbin_model" ("name")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_data_permission_subject_resource_action_active"
    ON "public"."sys_data_permission" ("subject_type", "subject_id", "resource_table", "action_key")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_dict_type_type_code_active"
    ON "public"."sys_dict_type" ("type_code")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_dict_entry_type_lang_value_active"
    ON "public"."sys_dict_entry" ("sys_dict_type_id", "language_code", "entry_value")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_language_type_code_active"
    ON "public"."sys_language_type" ("type_code")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_language_type_default_active"
    ON "public"."sys_language_type" ("is_default")
    WHERE "is_default" = TRUE AND "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_language_entry_code_type_active"
    ON "public"."sys_language_entry" ("entry_code", "sys_language_type_id")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_resource_code_active"
    ON "public"."sys_resource" ("code")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_role_code_active"
    ON "public"."sys_role" ("code")
    WHERE "deleted_at" = 0;

CREATE UNIQUE INDEX IF NOT EXISTS "idx_sys_user_role_user_role_active"
    ON "public"."sys_user_role" ("user_id", "role_id")
    WHERE "deleted_at" = 0;

CREATE INDEX IF NOT EXISTS "idx_sys_api_log_created_at"
    ON "public"."sys_api_log" ("created_at");

CREATE INDEX IF NOT EXISTS "idx_sys_api_log_user_created_at"
    ON "public"."sys_api_log" ("sys_user_id", "created_at");

CREATE INDEX IF NOT EXISTS "idx_sys_api_log_success_created_at"
    ON "public"."sys_api_log" ("success", "created_at");

CREATE INDEX IF NOT EXISTS "idx_sys_api_log_status_created_at"
    ON "public"."sys_api_log" ("status_code", "created_at");

ALTER TABLE "public"."sys_data_permission"
    DROP CONSTRAINT IF EXISTS "chk_sys_data_permission_subject_type",
    ADD CONSTRAINT "chk_sys_data_permission_subject_type"
        CHECK ("subject_type" IN ('USER', 'ROLE', 'ANY_USER', 'ANY_ROLE'));

ALTER TABLE "public"."sys_data_permission"
    DROP CONSTRAINT IF EXISTS "chk_sys_data_permission_scope_type",
    ADD CONSTRAINT "chk_sys_data_permission_scope_type"
        CHECK ("scope_type" IN ('all', 'none', 'include', 'exclude', 'custom'));

ALTER TABLE "public"."sys_data_permission"
    DROP CONSTRAINT IF EXISTS "chk_sys_data_permission_action",
    ADD CONSTRAINT "chk_sys_data_permission_action"
        CHECK (
            jsonb_typeof("action") = 'array'
            AND jsonb_array_length("action") > 0
            AND "action" <@ '["all", "read", "write", "delete"]'::jsonb
        );

ALTER TABLE "public"."sys_data_permission"
    DROP CONSTRAINT IF EXISTS "chk_sys_data_permission_scope_values",
    ADD CONSTRAINT "chk_sys_data_permission_scope_values"
        CHECK (
            "scope_type" NOT IN ('include', 'exclude')
            OR (
                jsonb_typeof("scope_values") = 'array'
                AND jsonb_array_length("scope_values") > 0
            )
        );

ALTER TABLE "public"."sys_data_permission"
    DROP CONSTRAINT IF EXISTS "chk_sys_data_permission_any_subject",
    ADD CONSTRAINT "chk_sys_data_permission_any_subject"
        CHECK (
            ("subject_type" IN ('ANY_USER', 'ANY_ROLE') AND "subject_id" = 0)
            OR ("subject_type" IN ('USER', 'ROLE') AND "subject_id" > 0)
        );

COMMIT;
