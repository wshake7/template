INSERT INTO "public"."sys_casbin_model" VALUES (1, NULL, NULL, 0, 0, 0, 't', '', 0, 'pbac', '[request_definition]
                r = sub, obj, act

                [policy_definition]
                p = sub_rule, obj_rule, act

                [policy_effect]
                e = some(where (p.eft == allow))

                [matchers]
                m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act');

INSERT INTO "public"."sys_data_permission" VALUES (1, NULL, NULL, 'root can access all dict types', 0, 0, 0, 't', 0, 'ROLE', 1, 'sys_dict_type', '["all"]', 'all', 'id', '[]', '{}', 100);
INSERT INTO "public"."sys_data_permission" VALUES (2, NULL, NULL, 'all roles can read all dict types', 0, 0, 0, 't', 0, 'ANY_ROLE', 0, 'sys_dict_type', '["read"]', 'all', 'id', '[]', '{}', 0);
INSERT INTO "public"."sys_data_permission" VALUES (3, NULL, NULL, 'all roles can operate dict types except system:is_enabled', 0, 0, 0, 't', 0, 'ANY_ROLE', 0, 'sys_dict_type', '["write", "delete"]', 'custom', 'id', '[]', '{"id__not": 1}', 0);

INSERT INTO "public"."sys_dict_type" VALUES (1, NULL, NULL, 0, 0, 0, 't', 0, '', 0, 'system:is_enabled', '开关状态');

INSERT INTO "public"."sys_dict_entry" VALUES (1, NULL, NULL, 0, 0, 0, 0, 't', '', 0, '', '启用', '1', '', 1);
INSERT INTO "public"."sys_dict_entry" VALUES (2, NULL, NULL, 0, 0, 0, 0, 't', '', 0, '', '停用', '0', '', 1);

INSERT INTO "public"."sys_role" VALUES (1, NULL, NULL, '', 0, 0, 0, 't', 0, '超级管理员', 'root', NULL, '[]');

INSERT INTO "public"."sys_user" VALUES (1, '2026-05-06 14:13:43.158062+08', '2026-05-06 14:13:43.158062+08', '', 0, 0, 0, 't', 0, 'root', '', '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO', NULL, '', '');
INSERT INTO "public"."sys_user" VALUES (2, '2026-05-06 14:13:43.159636+08', '2026-05-06 14:13:43.159636+08', '', 0, 0, 0, 't', 0, 'admin', '', '$2a$04$ASoVUxXahEpdD9.dxfwsHuUw3PqQ/yAZ0gD2KnqtMAqSGZ4VZCSVO', NULL, '', '');

INSERT INTO "public"."sys_user_role" VALUES (1, NULL, NULL, 0, 0, 0, 1, 1, 0);
