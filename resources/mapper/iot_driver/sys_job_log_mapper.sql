

-- name: ListSysJobLog
select
     t.create_time
      ,t.exception_info
      ,t.status
      ,t.job_message
      ,t.invoke_target
      ,t.job_group
      ,t.job_name
      ,t.job_log_id

from sys_job_log t where 1=1 and t.del_flag=0

    {{if (ne .ExceptionInfo "") }}
        and  t.exception_info like concat('%', @ExceptionInfo,'%')
    {{end}}
    {{if (ne .Status "") }}
        and  t.status like concat('%', @Status,'%')
    {{end}}
    {{if (ne .JobMessage "") }}
        and  t.job_message like concat('%', @JobMessage,'%')
    {{end}}
    {{if (ne .InvokeTarget "") }}
        and  t.invoke_target like concat('%', @InvokeTarget,'%')
    {{end}}

    order by
        {{if (eq .SortName "jobLogId") }}
           job_log_id
        {{else}}
           job_log_id
        {{end}}

        {{if (eq .SortOrder "asc") }}
             asc
        {{else}}
             desc
        {{end}}



