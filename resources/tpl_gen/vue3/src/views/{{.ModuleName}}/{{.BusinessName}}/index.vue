<template>
  <div class="app-container">
    <!-- 查询表单 -->
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
      {{range $column := .Columns}}
      {{if $column.Query}}
        {{$dictType := $column.DictType}}
        {{$parentheseIndex := strings.Index $column.ColumnComment "（"}}
        {{if ne $parentheseIndex -1}}
          {{$comment := substr $column.ColumnComment 0 $parentheseIndex}}
        {{else}}
          {{$comment := $column.ColumnComment}}
        {{end}}
        <!-- 输入框 -->
        {{if eq $column.HtmlType "input"}}
          <el-form-item label="{{$comment}}" prop="{{$column.JavaField}}">
            <el-input v-model="queryParams.{{$column.JavaField}}" placeholder="请输入{{$comment}}" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        <!-- 下拉/单选框（带字典） -->
        {{else if and (or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")) (ne $dictType "")}}
          <el-form-item label="{{$comment}}" prop="{{$column.JavaField}}">
            <el-select v-model="queryParams.{{$column.JavaField}}" placeholder="请选择{{$comment}}" clearable>
              <el-option v-for="dict in {{$dictType}}" :key="dict.value" :label="dict.label" :value="dict.value"/>
            </el-select>
          </el-form-item>
        <!-- 下拉/单选框（无字典） -->
        {{else if and (or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")) $dictType}}
          <el-form-item label="{{$comment}}" prop="{{$column.JavaField}}">
            <el-select v-model="queryParams.{{$column.JavaField}}" placeholder="请选择{{$comment}}" clearable>
              <el-option label="请选择字典生成" value="" />
            </el-select>
          </el-form-item>
        <!-- 日期选择（非范围） -->
        {{else if and (eq $column.HtmlType "datetime") (ne $column.QueryType "BETWEEN")}}
          <el-form-item label="{{$comment}}" prop="{{$column.JavaField}}">
            <el-date-picker clearable v-model="queryParams.{{$column.JavaField}}" type="date" value-format="YYYY-MM-DD" placeholder="请选择{{$comment}}" />
          </el-form-item>
        <!-- 日期范围选择 -->
        {{else if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
          <el-form-item label="{{$comment}}" style="width: 308px">
            <el-date-picker
              v-model="daterange{{upperFirst $column.JavaField}}"
              value-format="YYYY-MM-DD"
              type="daterange"
              range-separator="-"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
            ></el-date-picker>
          </el-form-item>
        {{end}}
      {{end}}
      {{end}}
      <!-- 搜索/重置按钮 -->
      <el-form-item>
        <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <!-- 操作按钮行 -->
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button
          type="primary"
          plain
          icon="Plus"
          @click="handleAdd"
          v-hasPermi="['{{.PermissionPrefix}}:add']"
        >新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="success"
          plain
          icon="Edit"
          :disabled="single"
          @click="handleUpdate"
          v-hasPermi="['{{.PermissionPrefix}}:edit']"
        >修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="danger"
          plain
          icon="Delete"
          :disabled="multiple"
          @click="handleDelete"
          v-hasPermi="['{{.PermissionPrefix}}:remove']"
        >删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button
          type="warning"
          plain
          icon="Download"
          @click="handleExport"
          v-hasPermi="['{{.PermissionPrefix}}:export']"
        >导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    <!-- 数据表格 -->
    <el-table v-loading="loading" :data="{{.BusinessName}}List" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
      {{range $column := .Columns}}
        {{$javaField := $column.JavaField}}
        {{$parentheseIndex := strings.Index $column.ColumnComment "（"}}
        {{if ne $parentheseIndex -1}}
          {{$comment := substr $column.ColumnComment 0 $parentheseIndex}}
        {{else}}
          {{$comment := $column.ColumnComment}}
        {{end}}
        <!-- 主键列 -->
        {{if $column.Pk}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$javaField}}" />
        <!-- 日期列 -->
        {{else if and $column.List (eq $column.HtmlType "datetime")}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$javaField}}" width="180">
            <template #default="scope">
              <span>{{ parseTime(scope.row.{{$javaField}}, '{y}-{m}-{d}') }}</span>
            </template>
          </el-table-column>
        <!-- 图片列 -->
        {{else if and $column.List (eq $column.HtmlType "imageUpload")}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$javaField}}" width="100">
            <template #default="scope">
              <image-preview :src="scope.row.{{$javaField}}" :width="50" :height="50"/>
            </template>
          </el-table-column>
        <!-- 字典列 -->
        {{else if and $column.List (ne $column.DictType "")}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$javaField}}">
            <template #default="scope">
              {{if eq $column.HtmlType "checkbox"}}
                <dict-tag :options="{{$column.DictType}}" :value="scope.row.{{$javaField}} ? scope.row.{{$javaField}}.split(',') : []"/>
              {{else}}
                <dict-tag :options="{{$column.DictType}}" :value="scope.row.{{$javaField}}"/>
              {{end}}
            </template>
          </el-table-column>
        <!-- 普通列 -->
        {{else if and $column.List (ne $javaField "")}}
          <el-table-column label="{{$comment}}" align="center" prop="{{$javaField}}" />
        {{end}}
      {{end}}
      <!-- 操作列 -->
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['{{.PermissionPrefix}}:edit']">修改</el-button>
          <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['{{.PermissionPrefix}}:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页组件 -->
    <pagination
      v-show="total>0"
      :total="total"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      @pagination="getList"
    />

    <!-- 添加/修改对话框 -->
    <el-dialog :title="title" v-model="open" width="500px" append-to-body>
      <el-form ref="{{.BusinessName}}Ref" :model="form" :rules="rules" label-width="80px">
        {{range $column := .Columns}}
          {{$field := $column.JavaField}}
          {{if and $column.Insert (not $column.Pk)}}
            {{if or $column.UsableColumn (not $column.SuperColumn)}}
              {{$parentheseIndex := strings.Index $column.ColumnComment "（"}}
              {{if ne $parentheseIndex -1}}
                {{$comment := substr $column.ColumnComment 0 $parentheseIndex}}
              {{else}}
                {{$comment := $column.ColumnComment}}
              {{end}}
              {{$dictType := $column.DictType}}
              <!-- 输入框 -->
              {{if eq $column.HtmlType "input"}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-input v-model="form.{{$field}}" placeholder="请输入{{$comment}}" />
                </el-form-item>
              <!-- 图片上传 -->
              {{else if eq $column.HtmlType "imageUpload"}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <image-upload v-model="form.{{$field}}"/>
                </el-form-item>
              <!-- 文件上传 -->
              {{else if eq $column.HtmlType "fileUpload"}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <file-upload v-model="form.{{$field}}"/>
                </el-form-item>
              <!-- 富文本 -->
              {{else if eq $column.HtmlType "editor"}}
                <el-form-item label="{{$comment}}">
                  <editor v-model="form.{{$field}}" :min-height="192"/>
                </el-form-item>
              <!-- 下拉框（带字典） -->
              {{else if and (eq $column.HtmlType "select") (ne $dictType "")}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-select v-model="form.{{$field}}" placeholder="请选择{{$comment}}">
                    <el-option
                      v-for="dict in {{$dictType}}"
                      :key="dict.value"
                      :label="dict.label"
                      {{if or (eq $column.JavaType "Integer") (eq $column.JavaType "Long")}}
                        :value="parseInt(dict.value)"
                      {{else}}
                        :value="dict.value"
                      {{end}}
                    ></el-option>
                  </el-select>
                </el-form-item>
              <!-- 多选框（带字典） -->
              {{else if and (eq $column.HtmlType "checkbox") (ne $dictType "")}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-checkbox-group v-model="form.{{$field}}">
                    <el-checkbox
                      v-for="dict in {{$dictType}}"
                      :key="dict.value"
                      :label="dict.value"
                    >{{dict.label}}</el-checkbox>
                  </el-checkbox-group>
                </el-form-item>
              <!-- 单选框（带字典） -->
              {{else if and (eq $column.HtmlType "radio") (ne $dictType "")}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-radio-group v-model="form.{{$field}}">
                    <el-radio
                      v-for="dict in {{$dictType}}"
                      :key="dict.value"
                      {{if or (eq $column.JavaType "Integer") (eq $column.JavaType "Long")}}
                        :label="parseInt(dict.value)"
                      {{else}}
                        :label="dict.value"
                      {{end}}
                    >{{dict.label}}</el-radio>
                  </el-radio-group>
                </el-form-item>
              <!-- 日期选择 -->
              {{else if eq $column.HtmlType "datetime"}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-date-picker clearable
                    v-model="form.{{$field}}"
                    type="date"
                    value-format="YYYY-MM-DD"
                    placeholder="请选择{{$comment}}"
                  ></el-date-picker>
                </el-form-item>
              <!-- 文本域 -->
              {{else if eq $column.HtmlType "textarea"}}
                <el-form-item label="{{$comment}}" prop="{{$field}}">
                  <el-input v-model="form.{{$field}}" type="textarea" placeholder="请输入内容" />
                </el-form-item>
              {{end}}
            {{end}}
          {{end}}
        {{end}}
        <!-- 子表部分（示例） -->
        {{if .Table.Sub}}
          <el-divider content-position="center">{{.SubTable.FunctionName}}信息</el-divider>
          <el-row :gutter="10" class="mb8">
            <el-col :span="1.5">
              <el-button type="primary" icon="Plus" @click="handleAdd{{.SubClassName}}">添加</el-button>
            </el-col>
            <el-col :span="1.5">
              <el-button type="danger" icon="Delete" @click="handleDelete{{.SubClassName}}">删除</el-button>
            </el-col>
          </el-row>
          <el-table :data="{{.SubClassName}}List" :row-class-name="row{{.SubClassName}}Index" @selection-change="handle{{.SubClassName}}SelectionChange" ref="{{.SubClassName}}">
            <el-table-column type="selection" width="50" align="center" />
            <el-table-column label="序号" align="center" prop="index" width="50"/>
            {{range $column := .SubTable.Columns}}
              {{$javaField := $column.JavaField}}
              {{$parentheseIndex := strings.Index $column.ColumnComment "（"}}
              {{if ne $parentheseIndex -1}}
                {{$comment := substr $column.ColumnComment 0 $parentheseIndex}}
              {{else}}
                {{$comment := $column.ColumnComment}}
              {{end}}
              {{if or $column.Pk (eq $javaField .SubTableFkclassName)}}
              {{else if and $column.List (eq $column.HtmlType "input")}}
                <el-table-column label="{{$comment}}" prop="{{$javaField}}" width="150">
                  <template #default="scope">
                    <el-input v-model="scope.row.{{$javaField}}" placeholder="请输入{{$comment}}" />
                  </template>
                </el-table-column>
              {{else if and $column.List (eq $column.HtmlType "datetime")}}
                <el-table-column label="{{$comment}}" prop="{{$javaField}}" width="240">
                  <template #default="scope">
                    <el-date-picker clearable
                      v-model="scope.row.{{$javaField}}"
                      type="date"
                      value-format="YYYY-MM-DD"
                      placeholder="请选择{{$comment}}"
                    ></el-date-picker>
                  </template>
                </el-table-column>
              {{end}}
            {{end}}
          </el-table>
        {{end}}
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

<script setup name="{{.BusinessName}}">
import { list{{.BusinessName}}, get{{.BusinessName}}, del{{.BusinessName}}, add{{.BusinessName}}, update{{.BusinessName}} } from "@/api/{{.ModuleName}}/{{.BusinessName}}";

const { proxy } = getCurrentInstance();
{{if ne .Dicts ""}}
  {{$dictsNoSymbol := strings.Replace .Dicts "'" "" -1}}
  const { {{$dictsNoSymbol}} } = proxy.useDict({{.Dicts}});
{{end}}

// 响应式数据声明
const {{.BusinessName}}List = ref([]);
{{if .Table.Sub}}
const {{.SubClassName}}List = ref([]);
{{end}}
const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
{{if .Table.Sub}}
const checked{{.SubClassName}} = ref([]);
{{end}}
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");
{{range $column := .Columns}}
{{if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
const daterange{{upperFirst $column.JavaField}} = ref([]);
{{end}}
{{end}}

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
    {{range $column := .Columns}}
    {{if $column.Query}}
    {{$column.JavaField}}: null{{if ne $loop.Last}},{{end}}
    {{end}}
    {{end}}
  },
  rules: {
    {{range $column := .Columns}}
    {{if $column.Required}}
    {{$column.JavaField}}: [
      { required: true, message: "{{substr $column.ColumnComment 0 (strings.Index $column.ColumnComment "（")}}不能为空", trigger: {{if or (eq $column.HtmlType "select") (eq $column.HtmlType "radio")}}"change"{{else}}"blur"{{end}} }
    ]{{if ne $loop.Last}},{{end}}
    {{end}}
    {{end}}
  }
});

const { queryParams, form, rules } = toRefs(data);

/** 查询{{.FunctionName}}列表 */
function getList() {
  loading.value = true;
  {{range $column := .Columns}}
  {{if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
  queryParams.value.params = {};
  {{break}}
  {{end}}
  {{end}}
  {{range $column := .Columns}}
  {{if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
  const attrName = "{{upperFirst $column.JavaField}}";
  if (daterange${attrName}.value.length > 0) {
    queryParams.value.params["begin${attrName}"] = daterange${attrName}.value[0];
    queryParams.value.params["end${attrName}"] = daterange${attrName}.value[1];
  }
  {{end}}
  {{end}}
  list{{.BusinessName}}(queryParams.value).then(response => {
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

// 表单重置
function reset() {
  form.value = {
    {{range $column := .Columns}}
    {{if eq $column.HtmlType "checkbox"}}
    {{$column.JavaField}}: []{{if ne $loop.Last}},{{end}}
    {{else}}
    {{$column.JavaField}}: null{{if ne $loop.Last}},{{end}}
    {{end}}
    {{end}}
  };
  {{if .Table.Sub}}
  {{.SubClassName}}List.value = [];
  {{end}}
  proxy.resetForm("{{.BusinessName}}Ref");
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

/** 重置按钮操作 */
function resetQuery() {
  {{range $column := .Columns}}
  {{if and (eq $column.HtmlType "datetime") (eq $column.QueryType "BETWEEN")}}
  daterange{{upperFirst $column.JavaField}}.value = [];
  {{end}}
  {{end}}
  proxy.resetForm("queryRef");
  handleQuery();
}

// 多选框选中数据
function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.{{.PkColumn.JavaField}});
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
  const _{{.PkColumn.JavaField}} = row?.{{.PkColumn.JavaField}} || ids.value;
  get{{.BusinessName}}(_{{.PkColumn.JavaField}}).then(response => {
    form.value = response.data;
    {{range $column := .Columns}}
    {{if eq $column.HtmlType "checkbox"}}
    form.value.{{$column.JavaField}} = form.value.{{$column.JavaField}}?.split(",") || [];
    {{end}}
    {{end}}
    {{if .Table.Sub}}
    {{.SubClassName}}List.value = response.data.{{.SubClassName}}List || [];
    {{end}}
    open.value = true;
    title.value = "修改{{.FunctionName}}";
  });
}

/** 提交按钮 */
function submitForm() {
  proxy.$refs["{{.BusinessName}}Ref"].validate(valid => {
    if (valid) {
      {{range $column := .Columns}}
      {{if eq $column.HtmlType "checkbox"}}
      form.value.{{$column.JavaField}} = form.value.{{$column.JavaField}}?.join(",") || "";
      {{end}}
      {{end}}
      {{if .Table.Sub}}
      form.value.{{.SubClassName}}List = {{.SubClassName}}List.value;
      {{end}}
      if (form.value.{{.PkColumn.JavaField}} != null) {
        update{{.BusinessName}}(form.value).then(response => {
          proxy.$message.success("修改成功");
          open.value = false;
          getList();
        });
      } else {
        add{{.BusinessName}}(form.value).then(response => {
          proxy.$message.success("新增成功");
          open.value = false;
          getList();
        });
      }
    }
  });
}

/** 删除按钮操作 */
function handleDelete(row) {
  const _{{.PkColumn.JavaField}}s = row?.{{.PkColumn.JavaField}} || ids.value;
  proxy.$confirm('是否确认删除{{.FunctionName}}编号为"' + _{{.PkColumn.JavaField}}s + '"的数据项？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return del{{.BusinessName}}(_{{.PkColumn.JavaField}}s);
  }).then(() => {
    getList();
    proxy.$message.success("删除成功");
  }).catch(() => {});
}

{{if .Table.Sub}}
/** {{.SubTable.FunctionName}}序号 */
function row{{.SubClassName}}Index({ row, rowIndex }) {
  row.index = rowIndex + 1;
}

/** {{.SubTable.FunctionName}}添加按钮操作 */
function handleAdd{{.SubClassName}}() {
  let obj = {};
  {{range $column := .SubTable.Columns}}
  {{if and (not $column.Pk) (ne $column.JavaField .SubTableFkclassName)}}
  obj.{{$column.JavaField}} = "";
  {{end}}
  {{end}}
  {{.SubClassName}}List.value.push(obj);
}

/** {{.SubTable.FunctionName}}删除按钮操作 */
function handleDelete{{.SubClassName}}() {
  if (checked{{.SubClassName}}.value.length === 0) {
    proxy.$message.error("请先选择要删除的{{.SubTable.FunctionName}}数据");
  } else {
    {{.SubClassName}}List.value = {{.SubClassName}}List.value.filter(item => 
      !checked{{.SubClassName}}.value.includes(item.index)
    );
  }
}

/** 子表复选框选中数据 */
function handle{{.SubClassName}}SelectionChange(selection) {
  checked{{.SubClassName}}.value = selection.map(item => item.index);
}
{{end}}

/** 导出按钮操作 */
function handleExport() {
  proxy.download('{{.ModuleName}}/{{.BusinessName}}/export', {
    ...queryParams.value
  }, '{{.BusinessName}}_' + Date.now() + '.xlsx');
}

// 初始化加载数据
getList();
</script>
