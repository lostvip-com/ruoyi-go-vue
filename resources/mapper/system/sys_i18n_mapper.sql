

-- name: ListSysI18n
select
     t.id
      ,t.locale
      ,t.locale_key
      ,t.locale_name
      ,t.sort
      ,t.remark
      ,t.create_time
      ,t.update_time
      ,t.update_by
      ,t.create_by

from sys_i18n t where 1=1

    {{if (ne .Locale "") }}
        and  t.locale like concat('%', @Locale,'%')
    {{end}}
    {{if (ne .LocaleKey "") }}
        and  t.locale_key like concat('%', @LocaleKey,'%')
    {{end}}
    {{if (ne .LocaleName "") }}
        and  t.locale_name like concat('%', @LocaleName,'%')
    {{end}}
   {{if (ne .Sort 0) }}
        and  t.sort = @Sort
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



