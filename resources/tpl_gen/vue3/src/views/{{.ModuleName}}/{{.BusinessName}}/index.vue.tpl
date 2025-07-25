{{- $tagS :="{{" -}}
{{- $tagE :="}}" -}}
<template>
  <div class="app-container">
    <!-- 查询表单 -->
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
      {{- range $column := .Columns -}}
      {{- if eq "1" $column.IsQuery -}}
        {{- $dictType := $column.DictType -}}
        {{- $comment := $column.ColumnComment -}}
        {{- if eq $column.HtmlType "input" -}}
          <el-form-item label="{{$comment}}" prop="{{$column.GoField}}">
            <el-input v-model="queryParams.{{$column.GoField}}" placeholder="请输入{{$comment}}" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        <!-- 下拉/单选框（带字典） -->
        {{- else if and (or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")) (ne $dictType "") -}}
          <el-form-item label="{{$comment}}" prop="{{$column.GoField}}">
            <el-select v-model="queryParams.{{$column.GoField}}" placeholder="请选择{{$comment}}" clearable>
              <el-option v-for="dict in {{$dictType}}" :key="dict.value" :label="dict.label" :value="dict.value"/>
            </el-select>
          </el-form-item>
        <!-- 下拉/单选框（无字典） -->
        {{- else if and (or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")) $dictType -}}
          <el-form-item label="{{$comment}}" prop="{{$column.GoField}}">
            <el-select v-model="queryParams.{{$column.GoField}}" placeholder="请选择{{$comment}}" clearable>
              <el-option label="请选择字典生成" value="" />
            </el-select>
          </el-form-item>
        <!-- 日期选择（非范围） -->
        {{- else if and (eq $column.HtmlType "datetime") (ne $column.QueryType "BETWEEN") -}}
          <el-form-item label="{{$comment}}" prop="{{$column.GoField}}">
            <el-date-picker clearable v-model="queryParams.{{$column.GoField}}" type="date" value-format="YYYY-MM-DD" placeholder="请选择{{$comment}}" />
          </el-form-item>
        <!-- 日期范围选择 -->
        {{- else if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN") -}}
          <el-form-item label="{{$comment}}" style="width: 308px">
            <el-date-picker v-model="daterange{{upperFirst $column.GoField}}" value-format="YYYY-MM-DD" type="daterange" range-separator="-" start-placeholder="开始日期" end-placeholder="结束日期"></el-date-picker>
          </el-form-item>
        {{- end -}}
      {{- end -}}
      {{- end -}}
      <!-- 搜索/重置按钮 -->
      <el-form-item>
        <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 操作按钮行 -->
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" @click="handleAdd" v-hasPermi="['{{.BusinessName}}:add']">新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="success" plain icon="Edit" :disabled="single" @click="handleUpdate" v-hasPermi="['{{.BusinessName}}:edit']">修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="handleDelete" v-hasPermi="['{{.BusinessName}}:remove']">删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain icon="Download" @click="handleExport" v-hasPermi="['{{.BusinessName}}:export']">导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <!-- 数据表格 -->
    <el-table v-loading="loading" :data="{{.BusinessName}}List" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      {{- range $column := .Columns -}}
        {{- $GoField := $column.GoField -}}
        {{- $comment := $column.ColumnComment -}}
        <!-- 主键列 -->
        {{- if eq "1" $column.IsPk -}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$GoField}}" />
        <!-- 日期列 -->
        {{- else if (eq $column.HtmlType "datetime") -}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$GoField}}" width="180">
            <template #default="scope">
              <span>{{$tagS}} parseTime(scope.row.{{$GoField}}, '{y}-{m}-{d}') {{$tagE}}</span>
            </template>
          </el-table-column>
        <!-- 图片列 -->
        {{- else if (eq $column.HtmlType "imageUpload") -}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$GoField}}" width="100">
            <template #default="scope">
              <image-preview :src="scope.row.{{$GoField}}" :width="50" :height="50"/>
            </template>
          </el-table-column>
        {{- else if (ne $column.DictType "") -}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$GoField}}">
            <template #default="scope">
              {{- if eq $column.HtmlType "checkbox" -}}
                <dict-tag :options="{{$column.DictType}}" :value="scope.row.{{$GoField}} ? scope.row.{{$GoField}}.split(',') : []"/>
              {{- else -}}
                <dict-tag :options="{{$column.DictType}}" :value="scope.row.{{$GoField}}"/>
              {{- end -}}
            </template>
          </el-table-column>
        {{- else if and  (ne $GoField "") -}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$GoField}}" />
        {{- end -}}
      {{- end  -}}
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['{{.BusinessName}}:edit']">修改</el-button>
          <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['{{.BusinessName}}:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <!-- 分页组件 -->
    <pagination v-show="total>0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />
    <!-- 添加/修改对话框 -->
    <el-dialog :title="title" v-model="open" width="500px" append-to-body>
        <el-form ref="{{.BusinessName}}Ref" :model="form" :rules="rules" label-width="80px">
        {{ range $column := .Columns }}
          {{- $field := $column.GoField  -}}
          {{ if and (eq "1" $column.IsInsert) (ne "1" $column.IsPk) }}
              {{- $comment := $column.ColumnComment -}}
              {{- $dictType := $column.DictType}}
              {{- if eq $column.HtmlType "input" }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-input v-model="form.{{$field}}" placeholder="请输入{{$comment}}" />
                </el-form-item>
          {{- else if eq $column.HtmlType "imageUpload" }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <image-upload v-model="form.{{$field}}"/>
                </el-form-item>
          {{- else if eq $column.HtmlType "fileUpload" -}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <file-upload v-model="form.{{$field}}"/>
                </el-form-item>
          {{- else if eq $column.HtmlType "editor" }}
                <el-form-item label="{{$comment}}">
                  <editor v-model="form.{{$field}}" :min-height="192"/>
                </el-form-item>
          {{- else if and (eq $column.HtmlType "select") (ne $dictType "") }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-select v-model="form.{{$field}}" placeholder="请选择{{$comment}}">
                    <el-option v-for="dict in {{$dictType}}" :key="dict.value" :label="dict.label"
                      {{- if or (eq $column.GoType "Integer") (eq $column.GoType "Long")}} :value="parseInt(dict.value)" {{else}} :value="dict.value" {{ end -}}></el-option>
                  </el-select>
                </el-form-item>
          <!-- 多选框（带字典） -->
          {{- else if and (eq $column.HtmlType "checkbox") (ne $dictType "") }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-checkbox-group v-model="form.{{$field}}">
                    <el-checkbox v-for="dict in {{$dictType}}" :key="dict.value" :label="dict.value">{{$tagS}} dict.label {{$tagE}} </el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              <!-- 单选框（带字典） -->
          {{- else if and (eq $column.HtmlType "radio") (ne $dictType "") }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-radio-group v-model="form.{{$field}}">
                    <el-radio v-for="dict in {{$dictType}}" :key="dict.value"
                      {{- if or (eq $column.GoType "Integer") (eq $column.GoType "Long") }}
                       :label="parseInt(dict.value)"
                      {{- else -}}
                        :label="dict.value"
                      {{- end -}}>
                      {{$tagS}} dict.label {{$tagE}}
                      </el-radio>
                  </el-radio-group>
                </el-form-item>
              <!-- 日期选择 -->
          {{- else if eq $column.HtmlType "datetime" }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-date-picker clearable v-model="form.{{$field}}" type="date" value-format="YYYY-MM-DD" placeholder="请选择{{$comment}}"></el-date-picker>
                </el-form-item>
          {{- else if eq $column.HtmlType "textarea" }}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-input v-model="form.{{$field}}" type="textarea" placeholder="请输入内容" />
                </el-form-item>
          {{- end }}
          {{- end -}}
        {{- end -}}
        <!-- 子表部分（示例） -->
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="cancel">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

{{raw "<"}}script setup name="{{.BusinessName}}"{{raw ">"}}
import request from '@/utils/request'

const { proxy } = getCurrentInstance();
// 响应式数据声明
const {{.BusinessName}}List = ref([]);

const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
{{- range $column := .Columns -}}
  {{- if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
    const daterange{{upperFirst $column.GoField}} = ref([]);
  {{- end -}}
{{- end }}

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    {{- range $column := .Columns -}}
        {{- if eq "1" $column.IsQuery}}
          {{$column.GoField}}: null,
        {{- end -}}
    {{- end }}
  },
  rules: {
    {{- range $column := .Columns -}}
    {{- if eq "1" $column.IsRequired}}
        {{$column.GoField}}: [
          { required: true, message: "{{$column.ColumnComment}}不能为空", trigger: {{- if or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")}}"change"{{else}}"blur"{{end}} }
        ],
    {{- end -}}
    {{- end }}
  }
});

const { queryParams, form, rules } = toRefs(data);

/** 查询{{.BusinessName}}列表 */
function getList() {
  loading.value = true;
  {{- range $column := .Columns -}}
      {{- if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN") -}}
      queryParams.value.params = {};
      {{break}}
      {{- end -}}
  {{- end -}}
  {{- range $column := .Columns -}}
      {{- if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN") }}
          const attrName = "{{upperFirst $column.GoField}}";
          if (daterange${attrName}.value.length > 0) {
            queryParams.value.params["begin${attrName}"] = daterange${attrName}.value[0];
            queryParams.value.params["end${attrName}"] = daterange${attrName}.value[1];
          }
      {{- end -}}
  {{- end -}}
  request({
      url: '/{{.ModuleName}}/{{.BusinessName}}/list{{.ClassName}}',
      method: 'get',
      params: queryParams.value
    }).then(response => {
      {{.BusinessName}}List.value = response.rows;
      total.value = response.total;
      loading.value = false;
  });
}

// 取消按钮
function cancel() {
  open.value = false;
  reset();
}

function del{{.ClassName}}(id) {
   return request({
     url: '/{{.ModuleName}}/{{.BusinessName}}/' + id,
     method: 'delete'
   })
}

// 表单重置
function reset() {
  form.value = {
    {{- range $column := .Columns -}}
        {{- if eq $column.HtmlType "checkbox" }}
            {{$column.GoField}}: [],
        {{- else }}
            {{$column.GoField}}: undefined,
        {{- end }}
    {{- end -}}
  };

  proxy.resetForm("{{.BusinessName}}Ref");
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

/** 重置按钮操作 */
function resetQuery() {
  {{- range $column := .Columns -}}
      {{- if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
      daterange{{upperFirst $column.GoField}}.value = [];
      {{- end}}
  {{- end }}
  proxy.resetForm("queryRef");
  handleQuery();
}

// 多选框选中数据
function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.{{.PkColumn.GoField}});
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

/** 新增按钮操作 */
function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加{{.FunctionName}}";
}

/** 修改按钮操作 */
function handleUpdate(row) {
  reset();
  const id = row.{{.PkColumn.GoField}} || ids.value;
  request({
      url: '/{{.ModuleName}}/{{.BusinessName}}/' + id,
      method: 'get'
  }).then(response => {
    form.value = response.data;
    {{- range $column := .Columns -}}
        {{- if eq $column.HtmlType "checkbox"}}
            form.value.{{$column.GoField}} = form.value.{{$column.GoField}}?.split(",") || [];
        {{end}}
    {{- end -}}

    open.value = true;
    title.value = "修改{{.FunctionName}}";
  });
}

/** 提交按钮 */
function submitForm() {
  proxy.$refs["{{.BusinessName}}Ref"].validate(valid => {
    if (valid) {
      {{- range $column := .Columns -}}
          {{- if eq $column.HtmlType "checkbox"}}
          form.value.{{$column.GoField}} = form.value.{{$column.GoField}}?.join(",") || "";
          {{- end}}
      {{- end -}}

      let url = '/{{.ModuleName}}/{{.BusinessName}}';
      if (form.value.{{.PkColumn.GoField}} != null) {
        request({ url: url, method: 'put', data: form.value}).then(response => {
                proxy.$message.success("修改成功");
            });
      } else {
        request({ url: url,method: 'post',data: form.value}).then(response => {
                proxy.$message.success("新增成功");
            });
      }
      open.value = false;
      getList();
    }
  });
}

/** 删除按钮操作 */
function handleDelete(row) {
  const ids = row.{{.PkColumn.GoField}} || ids.value;
  proxy.$confirm('是否确认删除{{.FunctionName}}编号为"' + ids + '"的数据项？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return del{{.ClassName}}(ids);
  }).then(() => {
    getList();
    proxy.$message.success("删除成功");
  }).catch(() => {});
}

/** 导出按钮操作 */
function handleExport() {
  proxy.download('{{.ModuleName}}/{{.BusinessName}}/export{{.ClassName}}', {
    ...queryParams.value
  }, '{{.BusinessName}}_' + Date.now() + '.xlsx');
}
// 初始化加载数据
getList();
{{raw "<"}}/script{{raw ">"}}
