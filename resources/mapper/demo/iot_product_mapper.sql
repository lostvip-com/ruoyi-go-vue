

-- name: ListIotProduct
select
     t.id
      ,t.key
      ,t.name
      ,t.cloud_product_id
      ,t.cloud_instance_id
      ,t.platform
      ,t.protocol
      ,t.node_type
      ,t.net_type
      ,t.data_format
      ,t.last_sync_time
      ,t.factory
      ,t.description
      ,t.status
      ,t.extra
      ,t.del_flag
      ,t.create_time
      ,t.update_time
      ,t.update_by
      ,t.create_by
      ,t.manufacturer
      ,t.tenant_id

from iot_product t where 1=1 and t.del_flag=0

    {{if (ne .Key "") }}
        and  t.key like concat('%', @Key,'%')
    {{end}}
    {{if (ne .Name "") }}
        and  t.name like concat('%', @Name,'%')
    {{end}}
    {{if (ne .CloudProductId "") }}
        and  t.cloud_product_id like concat('%', @CloudProductId,'%')
    {{end}}
    {{if (ne .CloudInstanceId "") }}
        and  t.cloud_instance_id like concat('%', @CloudInstanceId,'%')
    {{end}}

    order by
        {{if (eq .SortName "id") }}
           id
        {{else}}
           id
        {{end}}

        {{if (eq .SortOrder "asc") }}
             asc
        {{else}}
             desc
        {{end}}



