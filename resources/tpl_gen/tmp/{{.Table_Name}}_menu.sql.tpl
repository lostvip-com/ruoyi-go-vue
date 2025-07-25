/*
==========================================================================
LV自动生成菜单SQL,只生成一次,按需修改.
生成日期：{{FmtTime .CreateTime}}
生成路径: tmp/sql/{{.BusinessName}}/{{.Table_Name}}_menu.sql
生成人：{{.FunctionAuthor}}
==========================================================================
*/

-- name: create_menu
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}', '{{.ParentMenuId}}', '1', '{{.BusinessName}}', '{{.ModuleName}}/{{.BusinessName}}/index', 1, 0, 'C', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:list', '#', 'admin', sysdate(), '', null, '{{.FunctionName}}菜单');

SELECT @parentId := LAST_INSERT_ID();

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}查询', @parentId, '1',  '#', '', 1, 0, 'F', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:query',        '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}新增', @parentId, '2',  '#', '', 1, 0, 'F', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:new',          '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}修改', @parentId, '3',  '#', '', 1, 0, 'F', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:edit',         '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}删除', @parentId, '4',  '#', '', 1, 0, 'F', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:del',       '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('{{.FunctionName}}导出', @parentId, '5',  '#', '', 1, 0, 'F', '0', '0', '{{.ModuleName}}:{{.BusinessName}}:export',       '#', 'admin', sysdate(), '', null, '');