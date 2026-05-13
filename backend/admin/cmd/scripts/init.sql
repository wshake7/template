INSERT INTO "public"."sys_casbin_model"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "is_enabled", "remark", "deleted_at", "name", "content")
VALUES
    (1, NULL, NULL, 0, 0, 0, TRUE, '', 0, 'pbac', '[request_definition]
                r = sub, obj, act

                [policy_definition]
                p = sub_rule, obj_rule, act

                [policy_effect]
                e = some(where (p.eft == allow))

                [matchers]
                m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act');

INSERT INTO "public"."sys_data_permission"
("id", "created_at", "updated_at", "remark", "created_by", "updated_by", "deleted_by", "is_enabled", "deleted_at", "subject_type", "subject_id", "resource_table", "action", "action_key", "scope_type", "scope_field", "scope_values", "conditions", "priority")
VALUES
    (1, NULL, NULL, 'root can access all dict types', 0, 0, 0, TRUE, 0, 'ROLE', 1, 'sys_dict_type', '["all"]', 'all', 'all', 'id', '[]', '{}', 100),
    (2, NULL, NULL, 'all roles can read all dict types', 0, 0, 0, TRUE, 0, 'ANY_ROLE', 0, 'sys_dict_type', '["read"]', 'read', 'all', 'id', '[]', '{}', 0),
    (3, NULL, NULL, 'all roles can operate dict types except system:is_enabled', 0, 0, 0, TRUE, 0, 'ANY_ROLE', 0, 'sys_dict_type', '["write","delete"]', 'write,delete', 'custom', 'id', '[]', '{"id__not": 1}', 0);

INSERT INTO "public"."sys_dict_type"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "is_enabled", "sort_order", "remark", "deleted_at", "type_code", "type_name")
VALUES
    (1, NULL, NULL, 0, 0, 0, TRUE, 0, '', 0, 'system:is_enabled', '开关状态');

INSERT INTO "public"."sys_dict_entry"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "sort_order", "is_enabled", "remark", "deleted_at", "label_component", "entry_label", "entry_value", "language_code", "sys_dict_type_id")
VALUES
    (1, NULL, NULL, 0, 0, 0, 0, TRUE, '', 0, '<Tag color="error">${EntryLabel}</Tag>', '启用', '1', '', 1),
    (2, NULL, NULL, 0, 0, 0, 0, TRUE, '', 0, '<Tag color="success">${EntryLabel}</Tag>', '停用', '0', '', 1);

INSERT INTO "public"."sys_role"
("id", "created_at", "updated_at", "remark", "created_by", "updated_by", "deleted_by", "is_enabled", "deleted_at", "name", "code", "parent_id")
VALUES
    (1, NULL, NULL, '', 0, 0, 0, TRUE, 0, '超级管理员', 'root', NULL);

INSERT INTO "public"."sys_user"
("id", "created_at", "updated_at", "remark", "created_by", "updated_by", "deleted_by", "is_enabled", "deleted_at", "username", "nickname", "password", "language_code")
VALUES
    (1, '2026-05-06 14:13:43.158062+08', '2026-05-06 14:13:43.158062+08', '', 0, 0, 0, TRUE, 0, 'root', '', '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO', ''),
    (2, '2026-05-06 14:13:43.159636+08', '2026-05-06 14:13:43.159636+08', '', 0, 0, 0, TRUE, 0, 'admin', '', '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO', '');

INSERT INTO "public"."sys_user_role"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "user_id", "role_id", "deleted_at")
VALUES
    (1, NULL, NULL, 0, 0, 0, 1, 1, 0);

INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (1, '2026-05-09 15:44:29.247468+08', '2026-05-09 17:31:30.430037+08', 1, 1, 0, '', 0, '{"icon": "", "hidden": false}', 't', 0, 'CATALOG', '/dashboard', '', '', '控制台1', '', NULL, '/1/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (2, '2026-05-09 16:01:45.414731+08', '2026-05-09 17:54:40.486485+08', 1, 1, 0, '', 0, '{"icon": "SettingOutlined", "hidden": false}', 't', 0, 'CATALOG', '/system', '', '', '系统管理', '', NULL, '/2/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (3, '2026-05-09 16:13:02.748554+08', '2026-05-10 20:06:34.713589+08', 1, 1, 0, '', -2, '{"icon": "", "order": -2, "hidden": false}', 't', 0, 'MENU', '/system/resource/menu', '', '', '菜单管理', '/system/resource.menu.tsx', 2, '/2/3/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (4, '2026-05-10 20:10:00.656317+08', '2026-05-10 20:10:45.771044+08', 1, 1, 0, '', 0, '{"icon": "ReconciliationOutlined", "order": 0, "hidden": false}', 't', 0, 'CATALOG', '/logger', '', '', '日志信息', '', NULL, '/10/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (5, '2026-05-09 16:42:05.166618+08', '2026-05-10 19:54:52.382881+08', 1, 1, 0, '', 0, '{"icon": "", "hidden": false}', 't', 0, 'MENU', '/system/language', '', '', '语言管理', '/system/language.tsx', 2, '/2/5/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (6, '2026-05-09 16:42:36.997454+08', '2026-05-10 19:54:57.037816+08', 1, 1, 0, '', 0, '{"icon": "", "hidden": false}', 't', 0, 'MENU', '/system/dict', '', '', '字典管理', '/system/dict.tsx', 2, '/2/6/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (7, '2026-05-09 16:49:51.215257+08', '2026-05-10 20:28:47.004123+08', 1, 1, 0, '', 0, '{"icon": "UserOutlined", "order": 0, "hidden": false}', 't', 0, 'CATALOG', '/account', '', '', '账号管理', '', NULL, '/7/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (8, '2026-05-09 16:50:35.892181+08', '2026-05-09 17:55:29.888606+08', 1, 1, 0, '', 0, '{"icon": "", "hidden": false}', 't', 0, 'MENU', '/account/role', '', '', '角色管理', '/account/role.tsx', 7, '/7/8/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (9, '2026-05-10 17:38:27.019193+08', '2026-05-10 20:06:43.984546+08', 1, 1, 0, '', -1, '{"icon": "", "order": -1, "hidden": false}', 't', 0, 'MENU', '/system/resource/api', '', '', 'API管理', '/system/resource.api.tsx', 2, '/2/9/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (10, '2026-05-09 16:17:20.233627+08', '2026-05-10 20:10:15.368034+08', 1, 1, 0, '', 0, '{"icon": "", "order": 0, "hidden": false}', 't', 0, 'MENU', '/logger/api/log', '', '', 'API日志', '/logger/api.log.tsx', 4, '/10/4/');
INSERT INTO "public"."sys_resource_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path") VALUES (11, '2026-05-10 20:28:36.585473+08', '2026-05-10 20:28:36.587573+08', 1, 1, 0, '', 0, '{"icon": "", "order": 0, "hidden": false}', 't', 0, 'MENU', '/account/user', '', '', '用户管理', '/account/user.tsx', 7, '/7/11/');

INSERT INTO "public"."sys_resource_api"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "is_enabled", "deleted_at", "module", "path", "method")
VALUES
    (1, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/list', 'POST'),
    (2, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/tree', 'GET'),
    (3, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/:id/permissions', 'GET'),
    (4, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/create', 'POST'),
    (5, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/update', 'POST'),
    (6, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/del', 'POST'),
    (7, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'role', '/api/sys/role/permissions', 'POST'),
    (8, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'user', '/api/sys/user/list', 'POST'),
    (9, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'user', '/api/sys/user/create', 'POST'),
    (10, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'user', '/api/sys/user/update', 'POST'),
    (11, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'user', '/api/sys/user/del', 'POST'),
    (12, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/type/list', 'POST'),
    (13, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/type/create', 'POST'),
    (14, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/type/update', 'POST'),
    (15, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/type/del', 'POST'),
    (16, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/list', 'POST'),
    (17, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/match', 'POST'),
    (18, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/create', 'POST'),
    (19, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/update', 'POST'),
    (20, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/del', 'POST'),
    (21, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'dict', '/api/sys/dict/entry/batch/copy', 'POST'),
    (22, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/type/list', 'POST'),
    (23, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/type/create', 'POST'),
    (24, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/type/update', 'POST'),
    (25, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/type/del', 'POST'),
    (26, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/entry/list', 'POST'),
    (27, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/entry/create', 'POST'),
    (28, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/entry/update', 'POST'),
    (29, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/entry/del', 'POST'),
    (30, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'language', '/api/sys/language/entry/batch/create', 'POST'),
    (31, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'api_log', '/api/sys/api/log/list', 'POST'),
    (32, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'api_log', '/api/sys/api/log/detail', 'POST'),
    (33, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_menu', '/api/sys/resource/menu/list', 'POST'),
    (34, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_menu', '/api/sys/resource/menu/tree', 'GET'),
    (35, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_menu', '/api/sys/resource/menu/create', 'POST'),
    (36, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_menu', '/api/sys/resource/menu/update', 'POST'),
    (37, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_menu', '/api/sys/resource/menu/del', 'POST'),
    (38, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_api', '/api/sys/resource/api/list', 'POST'),
    (39, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_api', '/api/sys/resource/api/create', 'POST'),
    (40, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_api', '/api/sys/resource/api/update', 'POST'),
    (41, NULL, NULL, 0, 0, 0, '', 0, TRUE, 0, 'resource_api', '/api/sys/resource/api/del', 'POST');

INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (1, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 2, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (2, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 3, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (3, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 4, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (4, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 5, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (5, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 6, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (6, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 7, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (7, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 8, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (8, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 9, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (9, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 10, 0);
INSERT INTO "public"."sys_role_menu" ("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "menu_id", "deleted_at") VALUES (10, '2026-05-12 16:54:34.859724+08', '2026-05-12 16:54:34.859724+08', 1, 1, 0, 1, 11, 0);

INSERT INTO "public"."sys_role_api"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "role_id", "api_id", "deleted_at")
VALUES
    (1, NULL, NULL, 0, 0, 0, 1, 1, 0),
    (2, NULL, NULL, 0, 0, 0, 1, 2, 0),
    (3, NULL, NULL, 0, 0, 0, 1, 3, 0),
    (4, NULL, NULL, 0, 0, 0, 1, 4, 0),
    (5, NULL, NULL, 0, 0, 0, 1, 5, 0),
    (6, NULL, NULL, 0, 0, 0, 1, 6, 0),
    (7, NULL, NULL, 0, 0, 0, 1, 7, 0),
    (8, NULL, NULL, 0, 0, 0, 1, 8, 0),
    (9, NULL, NULL, 0, 0, 0, 1, 9, 0),
    (10, NULL, NULL, 0, 0, 0, 1, 10, 0),
    (11, NULL, NULL, 0, 0, 0, 1, 11, 0),
    (12, NULL, NULL, 0, 0, 0, 1, 12, 0),
    (13, NULL, NULL, 0, 0, 0, 1, 13, 0),
    (14, NULL, NULL, 0, 0, 0, 1, 14, 0),
    (15, NULL, NULL, 0, 0, 0, 1, 15, 0),
    (16, NULL, NULL, 0, 0, 0, 1, 16, 0),
    (17, NULL, NULL, 0, 0, 0, 1, 17, 0),
    (18, NULL, NULL, 0, 0, 0, 1, 18, 0),
    (19, NULL, NULL, 0, 0, 0, 1, 19, 0),
    (20, NULL, NULL, 0, 0, 0, 1, 20, 0),
    (21, NULL, NULL, 0, 0, 0, 1, 21, 0),
    (22, NULL, NULL, 0, 0, 0, 1, 22, 0),
    (23, NULL, NULL, 0, 0, 0, 1, 23, 0),
    (24, NULL, NULL, 0, 0, 0, 1, 24, 0),
    (25, NULL, NULL, 0, 0, 0, 1, 25, 0),
    (26, NULL, NULL, 0, 0, 0, 1, 26, 0),
    (27, NULL, NULL, 0, 0, 0, 1, 27, 0),
    (28, NULL, NULL, 0, 0, 0, 1, 28, 0),
    (29, NULL, NULL, 0, 0, 0, 1, 29, 0),
    (30, NULL, NULL, 0, 0, 0, 1, 30, 0),
    (31, NULL, NULL, 0, 0, 0, 1, 31, 0),
    (32, NULL, NULL, 0, 0, 0, 1, 32, 0),
    (33, NULL, NULL, 0, 0, 0, 1, 33, 0),
    (34, NULL, NULL, 0, 0, 0, 1, 34, 0),
    (35, NULL, NULL, 0, 0, 0, 1, 35, 0),
    (36, NULL, NULL, 0, 0, 0, 1, 36, 0),
    (37, NULL, NULL, 0, 0, 0, 1, 37, 0),
    (38, NULL, NULL, 0, 0, 0, 1, 38, 0),
    (39, NULL, NULL, 0, 0, 0, 1, 39, 0),
    (40, NULL, NULL, 0, 0, 0, 1, 40, 0),
    (41, NULL, NULL, 0, 0, 0, 1, 41, 0);

