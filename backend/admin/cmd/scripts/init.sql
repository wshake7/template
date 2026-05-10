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
(1, NULL, NULL, 0, 0, 0, 0, TRUE, '', 0, '', '启用', '1', '', 1),
(2, NULL, NULL, 0, 0, 0, 0, TRUE, '', 0, '', '停用', '0', '', 1);

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

INSERT INTO "public"."sys_resource_menu"
("id", "created_at", "updated_at", "created_by", "updated_by", "deleted_by", "remark", "sort_order", "metadata", "is_enabled", "deleted_at", "menu_type", "path", "redirect", "alias", "name", "component", "parent_id", "tree_path")
VALUES
(1, NULL, NULL, 0, 0, 0, '', 0, '{"order":0}', TRUE, 0, 'MENU', '/dashboard', '', '', '控制台', '', NULL, '/1/'),
(2, NULL, NULL, 0, 0, 0, '', 10, '{"order":10}', TRUE, 0, 'CATALOG', '/system', '', '', '系统管理', '', NULL, '/2/'),
(3, NULL, NULL, 0, 0, 0, '', 10, '{"order":10}', TRUE, 0, 'MENU', '/system/language', '', '', '语言管理', '', 2, '/2/3/'),
(4, NULL, NULL, 0, 0, 0, '', 20, '{"order":20}', TRUE, 0, 'MENU', '/system/dict', '', '', '数据字典', '', 2, '/2/4/'),
(5, NULL, NULL, 0, 0, 0, '', 30, '{"order":30}', TRUE, 0, 'MENU', '/system/api/log', '', '', 'API日志', '', 2, '/2/5/'),
(6, NULL, NULL, 0, 0, 0, '', 40, '{"order":40}', TRUE, 0, 'MENU', '/system/resource/menu', '', '', '菜单资源', '', 2, '/2/6/'),
(9, NULL, NULL, 0, 0, 0, '', 50, '{"order":50}', TRUE, 0, 'MENU', '/system/resource/api', '', '', 'API资源', '', 2, '/2/9/'),
(7, NULL, NULL, 0, 0, 0, '', 20, '{"order":20}', TRUE, 0, 'CATALOG', '/account', '', '', '账号管理', '', NULL, '/7/'),
(8, NULL, NULL, 0, 0, 0, '', 10, '{"order":10}', TRUE, 0, 'MENU', '/account/role', '', '', '角色管理', '', 7, '/7/8/');
