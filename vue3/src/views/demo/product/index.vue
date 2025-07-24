

<template>
  <div class="app-container">
    
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px">
        
        
        
          <el-form-item label="产品编码,对应可监控类型ID" prop="Key">
            <el-input v-model="queryParams.Key" placeholder="请输入产品编码,对应可监控类型ID" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        
        
        
        
          <el-form-item label="名字" prop="Name">
            <el-input v-model="queryParams.Name" placeholder="请输入名字" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        
        
        
        
          <el-form-item label="云产品ID" prop="CloudProductId">
            <el-input v-model="queryParams.CloudProductId" placeholder="请输入云产品ID" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        
        
        
        
          <el-form-item label="云实例ID" prop="CloudInstanceId">
            <el-input v-model="queryParams.CloudInstanceId" placeholder="请输入云实例ID" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        
      
      <el-form-item>
        <el-button type="primary" icon="Search" @click="handleQuery">搜索</el-button>
        <el-button icon="Refresh" @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    
    <el-row :gutter="10" class="mb8">
      <el-col :span="1.5">
        <el-button type="primary" plain icon="Plus" @click="handleAdd" v-hasPermi="['product:add']">新增</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="success" plain icon="Edit" :disabled="single" @click="handleUpdate" v-hasPermi="['product:edit']">修改</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="danger" plain icon="Delete" :disabled="multiple" @click="handleDelete" v-hasPermi="['product:remove']">删除</el-button>
      </el-col>
      <el-col :span="1.5">
        <el-button type="warning" plain icon="Download" @click="handleExport" v-hasPermi="['product:export']">导出</el-button>
      </el-col>
      <right-toolbar v-model:showSearch="showSearch" @queryTable="getList"></right-toolbar>
    </el-row>

    
    <el-table v-loading="loading" :data="productList" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" align="center" />
        
        <el-table-column label="主键" align="center" prop="Id" />
        
        
        <el-table-column label="产品编码,对应可监控类型ID" align="center" prop="Key" />
        
        <el-table-column label="名字" align="center" prop="Name" />
        
        <el-table-column label="云产品ID" align="center" prop="CloudProductId" />
        
        <el-table-column label="云实例ID" align="center" prop="CloudInstanceId" />
        
        <el-table-column label="平台" align="center" prop="Platform" />
        
        <el-table-column label="协议" align="center" prop="Protocol" />
        
        <el-table-column label="节点类型" align="center" prop="NodeType" />
        
        <el-table-column label="网络类型" align="center" prop="NetType" />
        
        <el-table-column label="数据类型" align="center" prop="DataFormat" />
        
        <el-table-column label="最后一次同步时间" align="center" prop="LastSyncTime" />
        
        <el-table-column label="工厂名称" align="center" prop="Factory" />
        
        <el-table-column label="描述" align="center" prop="Description" />
        
        <el-table-column label="产品状态" align="center" prop="Status" />
        
        <el-table-column label="扩展字段" align="center" prop="Extra" />
        
        <el-table-column label="删除标记" align="center" prop="DelFlag" />
        
        <el-table-column label="创建日期" align="center" prop="CreateTime" width="180">
            <template #default="scope">
              <span>{{ parseTime(scope.row.CreateTime, '{y}-{m}-{d}') }}</span>
            </template>
          </el-table-column>
        
        
        <el-table-column label="更新日期" align="center" prop="UpdateTime" width="180">
            <template #default="scope">
              <span>{{ parseTime(scope.row.UpdateTime, '{y}-{m}-{d}') }}</span>
            </template>
          </el-table-column>
        
        
        <el-table-column label="更新者" align="center" prop="UpdateBy" />
        
        <el-table-column label="创建者" align="center" prop="CreateBy" />
        
        <el-table-column label="生产厂商" align="center" prop="Manufacturer" />
        
        <el-table-column label="租户id" align="center" prop="TenantId" />
      <el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['product:edit']">修改</el-button>
          <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['product:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    
    <pagination v-show="total>0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />

    
    <el-dialog :title="title" v-model="open" width="500px" append-to-body>
      <el-form ref="productRef" :model="form" :rules="rules" label-width="80px">
              
              <el-form-item label="产品编码,对应可监控类型ID" prop="Key">
                  <el-input v-model="form.Key" placeholder="请输入产品编码,对应可监控类型ID" />
                </el-form-item>
              
              <el-form-item label="名字" prop="Name">
                  <el-input v-model="form.Name" placeholder="请输入名字" />
                </el-form-item>
              
              <el-form-item label="云产品ID" prop="CloudProductId">
                  <el-input v-model="form.CloudProductId" placeholder="请输入云产品ID" />
                </el-form-item>
              
              <el-form-item label="云实例ID" prop="CloudInstanceId">
                  <el-input v-model="form.CloudInstanceId" placeholder="请输入云实例ID" />
                </el-form-item>
              
              <el-form-item label="平台" prop="Platform">
                  <el-input v-model="form.Platform" placeholder="请输入平台" />
                </el-form-item>
              
              <el-form-item label="协议" prop="Protocol">
                  <el-input v-model="form.Protocol" placeholder="请输入协议" />
                </el-form-item>
              
              
              
              
              
              <el-form-item label="数据类型" prop="DataFormat">
                  <el-input v-model="form.DataFormat" placeholder="请输入数据类型" />
                </el-form-item>
              
              <el-form-item label="最后一次同步时间" prop="LastSyncTime">
                  <el-input v-model="form.LastSyncTime" placeholder="请输入最后一次同步时间" />
                </el-form-item>
              
              <el-form-item label="工厂名称" prop="Factory">
                  <el-input v-model="form.Factory" placeholder="请输入工厂名称" />
                </el-form-item>
              
              <el-form-item label="描述" prop="Description">
                  <el-input v-model="form.Description" type="textarea" placeholder="请输入内容" />
                </el-form-item>
              
              
              
              <el-form-item label="扩展字段" prop="Extra">
                  <el-input v-model="form.Extra" placeholder="请输入扩展字段" />
                </el-form-item>
              
              <el-form-item label="删除标记" prop="DelFlag">
                  <el-input v-model="form.DelFlag" placeholder="请输入删除标记" />
                </el-form-item>
              
              <el-form-item label="生产厂商" prop="Manufacturer">
                  <el-input v-model="form.Manufacturer" placeholder="请输入生产厂商" />
                </el-form-item>
              
              <el-form-item label="租户id" prop="TenantId">
                  <el-input v-model="form.TenantId" placeholder="请输入租户id" />
                </el-form-item>
        
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

<script setup name="product">
import { list"product", get"product", del"product", add"product", update"product" } from "@/api/demo/product";

const { proxy } = getCurrentInstance();

const "product"List = ref([]);

const open = ref(false);
const loading = ref(true);
const showSearch = ref(true);
const ids = ref([]);
const single = ref(true);
const multiple = ref(true);
const total = ref(0);
const title = ref("");

const data = reactive({
  form: {},
  queryParams: {
    pageNum: 1,
    pageSize: 10,
          "Key": null,
          "Name": null,
          "CloudProductId": null,
          "CloudInstanceId": null,
  },
  rules: {
        "Name": [
          { required: true, message: "名字不能为空", trigger:"blur" }
        ]
        "Status": [
          { required: true, message: "产品状态不能为空", trigger:"change" }
        ]
  }
});

const { queryParams, form, rules } = toRefs(data);

  
function getList() {
  loading.value = true;
  list"product"(queryParams.value).then(response => {
    "product"List.value = response.rows;
    total.value = response.total;
    loading.value = false;
  });
}


function cancel() {
  open.value = false;
  reset();
}


function reset() {
  form.value = {
            "Id": null
            "Key": null
            "Name": null
            "CloudProductId": null
            "CloudInstanceId": null
            "Platform": null
            "Protocol": null
            "NodeType": null
            "NetType": null
            "DataFormat": null
            "LastSyncTime": null
            "Factory": null
            "Description": null
            "Status": null
            "Extra": null
            "DelFlag": null
            "CreateTime": null
            "UpdateTime": null
            "UpdateBy": null
            "CreateBy": null
            "Manufacturer": null
            "TenantId": null};

  proxy.resetForm("productRef");
}

 
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

 
function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}


function handleSelectionChange(selection) {
  ids.value = selection.map(item => item."Id");
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

 
function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加product";
}

 
function handleUpdate(row) {
  reset();
  const _"Id" = row?."Id" || ids.value;
  get"product"(_"Id").then(response => {
    form.value = response.data;open.value = true;
    title.value = "修改product";
  });
}

 
function submitForm() {
  proxy.$refs["productRef"].validate(valid => {
    if (valid) {if (form.value."Id" != null) {
        update"product"(form.value).then(response => {
          proxy.$message.success("修改成功");
          open.value = false;
          getList();
        });
      } else {
        add"product"(form.value).then(response => {
          proxy.$message.success("新增成功");
          open.value = false;
          getList();
        });
      }
    }
  });
}

 
function handleDelete(row) {
  const _"Id"s = row?."Id" || ids.value;
  proxy.$confirm('是否确认删除product编号为"' + _"Id"s + '"的数据项？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return del"product"(_"Id"s);
  }).then(() => {
    getList();
    proxy.$message.success("删除成功");
  }).catch(() => {});
}


 
function handleExport() {
  proxy.download('demo/product/export', {
    ...queryParams.value
  }, 'product_' + Date.now() + '.xlsx');
}


getList();
</script>
