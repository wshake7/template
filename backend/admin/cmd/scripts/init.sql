INSERT INTO
  "public"."sys_casbin_model" (
    "id",
    "created_at",
    "updated_at",
    "deleted_at",
    "created_by",
    "updated_by",
    "deleted_by",
    "is_enabled",
    "remark",
    "name",
    "content"
  )
VALUES
  (
    1,
    NULL,
    NULL,
    NULL,
    0,
    0,
    0,
    true,
    '',
    'pbac',
    '[request_definition]
            r = sub, obj, act

            [policy_definition]
            p = sub_rule, obj_rule, act

            [policy_effect]
            e = some(where (p.eft == allow))

            [matchers]
            m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act'
  );

INSERT INTO
  "public"."sys_role" (
    "id",
    "created_at",
    "updated_at",
    "remark",
    "created_by",
    "updated_by",
    "is_enabled",
    "deleted_at",
    "name",
    "code",
    "parent_id",
    "child_ids"
  )
VALUES
  (
    1,
    NULL,
    NULL,
    '',
    0,
    0,
    true,
    NULL,
    '超级管理员',
    'root',
    NULL,
    '[]'
  );

INSERT INTO
  "public"."sys_dict_type" (
    "id",
    "created_at",
    "updated_at",
    "created_by",
    "updated_by",
    "is_enabled",
    "sort_order",
    "description",
    "deleted_at",
    "type_code",
    "type_name",
    "remark"
  )
VALUES
  (
    1,
    null,
    null,
    0,
    0,
    true,
    0,
    '',
    0,
    'system:is_enabled',
    '开关状态',
    ''
  );

INSERT INTO
  "public"."sys_dict_entry" (
    "id",
    "created_at",
    "updated_at",
    "deleted_at",
    "created_by",
    "updated_by",
    "sort_order",
    "is_enabled",
    "remark",
    "entry_label",
    "entry_value",
    "language_code",
    "sys_dict_type_id"
  )
VALUES
  (
    1,
    null,
    null,
    NULL,
    0,
    0,
    0,
    true,
    '',
    '启用',
    '1',
    '',
    1
  );

INSERT INTO
  "public"."sys_dict_entry" (
    "id",
    "created_at",
    "updated_at",
    "deleted_at",
    "created_by",
    "updated_by",
    "sort_order",
    "is_enabled",
    "remark",
    "entry_label",
    "entry_value",
    "language_code",
    "sys_dict_type_id"
  )
VALUES
  (
    2,
    null,
    null,
    NULL,
    0,
    0,
    0,
    true,
    '',
    '停用',
    '0',
    '',
    1
  );

INSERT INTO
  "public"."sys_data_permission" (
    "id",
    "created_at",
    "updated_at",
    "remark",
    "created_by",
    "updated_by",
    "is_enabled",
    "deleted_at",
    "subject_type",
    "subject_id",
    "resource_table",
    "action",
    "scope_type",
    "scope_values",
    "conditions",
    "priority"
  )
VALUES
  (
    1,
    NULL,
    NULL,
    '',
    0,
    0,
    true,
    NULL,
    'ROLE',
    1,
    'sys_dict_type',
    '["read"]',
    'all',
    '[]',
    '{"type_code": "system:is_enabled"}',
    0
  );
