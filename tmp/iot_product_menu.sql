/*
==========================================================================
LV自动生成菜单SQL,只生成一次,按需修改.
生成日期：2025-07-24 07:57:14
生成路径: tmp/sql/product/iot_product_menu.sql
生成人：lv
==========================================================================
*/

-- name: create_menu
insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product', '1062', '1', 'product', 'demo/product/index', 1, 0, 'C', '0', '0', 'demo:product:list', '#', 'admin', sysdate(), '', null, 'product菜单');

SELECT @parentId := LAST_INSERT_ID();

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product查询', @parentId, '1',  '#', '', 1, 0, 'F', '0', '0', 'demo:product:query',        '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product新增', @parentId, '2',  '#', '', 1, 0, 'F', '0', '0', 'demo:product:add',          '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product修改', @parentId, '3',  '#', '', 1, 0, 'F', '0', '0', 'demo:product:edit',         '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product删除', @parentId, '4',  '#', '', 1, 0, 'F', '0', '0', 'demo:product:remove',       '#', 'admin', sysdate(), '', null, '');

insert into sys_menu (menu_name, parent_id, order_num, path, component, is_frame, is_cache, menu_type, visible, status, perms, icon, create_by, create_time, update_by, update_time, remark)
values('product导出', @parentId, '5',  '#', '', 1, 0, 'F', '0', '0', 'demo:product:export',       '#', 'admin', sysdate(), '', null, '');