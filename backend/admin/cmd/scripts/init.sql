insert into sys_casbin_model (created_at, updated_at, deleted_at, created_by, updated_by, deleted_by, is_enabled,
                              remark, name, content)
values (null, null, null, 0, 0, 0, true, '', 'pbac', '[request_definition]
r = sub, obj, act

[policy_definition]
p = sub_rule, obj_rule, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act');

