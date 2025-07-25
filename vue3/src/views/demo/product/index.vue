<template>
  <div class="app-container">
    
    <el-form :model="queryParams" ref="queryRef" :inline="true" v-show="showSearch" label-width="68px"><el-form-item label="产品编码,对应可监控类型ID" prop="key">
            <el-input v-model="queryParams.key" placeholder="请输入产品编码,对应可监控类型ID" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        <el-form-item label="名字" prop="name">
            <el-input v-model="queryParams.name" placeholder="请输入名字" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        <el-form-item label="云产品ID" prop="cloudProductId">
            <el-input v-model="queryParams.cloudProductId" placeholder="请输入云产品ID" clearable @keyup.enter="handleQuery" />
          </el-form-item>
        <el-form-item label="云实例ID" prop="cloudInstanceId">
            <el-input v-model="queryParams.cloudInstanceId" placeholder="请输入云实例ID" clearable @keyup.enter="handleQuery" />
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
      <el-table-column type="selection" width="55" align="center" /><el-table-column label="主键" align="center" prop="id" />
        <el-table-column label="产品编码,对应可监控类型ID" align="center" prop="key" /><el-table-column label="名字" align="center" prop="name" /><el-table-column label="云产品ID" align="center" prop="cloudProductId" /><el-table-column label="云实例ID" align="center" prop="cloudInstanceId" /><el-table-column label="平台" align="center" prop="platform" /><el-table-column label="协议" align="center" prop="protocol" /><el-table-column label="节点类型" align="center" prop="nodeType" /><el-table-column label="网络类型" align="center" prop="netType" /><el-table-column label="数据类型" align="center" prop="dataFormat" /><el-table-column label="最后一次同步时间" align="center" prop="lastSyncTime" /><el-table-column label="工厂名称" align="center" prop="factory" /><el-table-column label="描述" align="center" prop="description" /><el-table-column label="产品状态" align="center" prop="status" /><el-table-column label="扩展字段" align="center" prop="extra" /><el-table-column label="删除标记" align="center" prop="delFlag" /><el-table-column label="创建日期" align="center" prop="createTime" width="180">
            <template #default="scope">
              <span>{{ parseTime(scope.row.createTime, '{y}-{m}-{d}') }}</span>
            </template>
          </el-table-column>
        <el-table-column label="更新日期" align="center" prop="updateTime" width="180">
            <template #default="scope">
              <span>{{ parseTime(scope.row.updateTime, '{y}-{m}-{d}') }}</span>
            </template>
          </el-table-column>
        <el-table-column label="更新者" align="center" prop="updateBy" /><el-table-column label="创建者" align="center" prop="createBy" /><el-table-column label="生产厂商" align="center" prop="manufacturer" /><el-table-column label="租户id" align="center" prop="tenantId" /><el-table-column label="操作" align="center" class-name="small-padding fixed-width">
        <template #default="scope">
          <el-button link type="primary" icon="Edit" @click="handleUpdate(scope.row)" v-hasPermi="['product:edit']">修改</el-button>
          <el-button link type="primary" icon="Delete" @click="handleDelete(scope.row)" v-hasPermi="['product:remove']">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <pagination v-show="total>0" :total="total" v-model:page="queryParams.pageNum" v-model:limit="queryParams.pageSize" @pagination="getList" />
    
    <el-dialog :title="title" v-model="open" width="500px" append-to-body>
        <el-form ref="productRef" :model="form" :rules="rules" label-width="80px">
        
                <el-form-item label="产品编码,对应可监控类型ID" prop="key">
                  <el-input v-model="form.key" placeholder="请输入产品编码,对应可监控类型ID" />
                </el-form-item>
                <el-form-item label="名字" prop="name">
                  <el-input v-model="form.name" placeholder="请输入名字" />
                </el-form-item>
                <el-form-item label="云产品ID" prop="cloudProductId">
                  <el-input v-model="form.cloudProductId" placeholder="请输入云产品ID" />
                </el-form-item>
                <el-form-item label="云实例ID" prop="cloudInstanceId">
                  <el-input v-model="form.cloudInstanceId" placeholder="请输入云实例ID" />
                </el-form-item>
                <el-form-item label="平台" prop="platform">
                  <el-input v-model="form.platform" placeholder="请输入平台" />
                </el-form-item>
                <el-form-item label="协议" prop="protocol">
                  <el-input v-model="form.protocol" placeholder="请输入协议" />
                </el-form-item>
                <el-form-item label="数据类型" prop="dataFormat">
                  <el-input v-model="form.dataFormat" placeholder="请输入数据类型" />
                </el-form-item>
                <el-form-item label="最后一次同步时间" prop="lastSyncTime">
                  <el-input v-model="form.lastSyncTime" placeholder="请输入最后一次同步时间" />
                </el-form-item>
                <el-form-item label="工厂名称" prop="factory">
                  <el-input v-model="form.factory" placeholder="请输入工厂名称" />
                </el-form-item>
                <el-form-item label="描述" prop="description">
                  <el-input v-model="form.description" type="textarea" placeholder="请输入内容" />
                </el-form-item>
                <el-form-item label="扩展字段" prop="extra">
                  <el-input v-model="form.extra" placeholder="请输入扩展字段" />
                </el-form-item>
                <el-form-item label="删除标记" prop="delFlag">
                  <el-input v-model="form.delFlag" placeholder="请输入删除标记" />
                </el-form-item>
                <el-form-item label="生产厂商" prop="manufacturer">
                  <el-input v-model="form.manufacturer" placeholder="请输入生产厂商" />
                </el-form-item>
                <el-form-item label="租户id" prop="tenantId">
                  <el-input v-model="form.tenantId" placeholder="请输入租户id" />
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
import request from '@/utils/request'

const { proxy } = getCurrentInstance();
// 响应式数据声明
const productList = ref([]);

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
          key: null,
          name: null,
          cloudProductId: null,
          cloudInstanceId: null,
  },
  rules: {
        name: [
          { required: true, message: "名字不能为空", trigger:"blur" }
        ],
        status: [
          { required: true, message: "产品状态不能为空", trigger:"change" }
        ],
  }
});

const { queryParams, form, rules } = toRefs(data);

/** 查询product列表 */
function getList() {
  loading.value = true;request({
      url: '/demo/product/listIotProduct',
      method: 'get',
      params: queryParams.value
    }).then(response => {
      productList.value = response.rows;
      total.value = response.total;
      loading.value = false;
  });
}

// 取消按钮
function cancel() {
  open.value = false;
  reset();
}

function delIotProduct(id) {
   return request({
     url: '/demo/product/' + id,
     method: 'delete'
   })
}

// 表单重置
function reset() {
  form.value = {
            id: undefined,
            key: undefined,
            name: undefined,
            cloudProductId: undefined,
            cloudInstanceId: undefined,
            platform: undefined,
            protocol: undefined,
            nodeType: undefined,
            netType: undefined,
            dataFormat: undefined,
            lastSyncTime: undefined,
            factory: undefined,
            description: undefined,
            status: undefined,
            extra: undefined,
            delFlag: undefined,
            createTime: undefined,
            updateTime: undefined,
            updateBy: undefined,
            createBy: undefined,
            manufacturer: undefined,
            tenantId: undefined,};

  proxy.resetForm("productRef");
}

/** 搜索按钮操作 */
function handleQuery() {
  queryParams.value.pageNum = 1;
  getList();
}

/** 重置按钮操作 */
function resetQuery() {
  proxy.resetForm("queryRef");
  handleQuery();
}

// 多选框选中数据
function handleSelectionChange(selection) {
  ids.value = selection.map(item => item.id);
  single.value = selection.length !== 1;
  multiple.value = !selection.length;
}

/** 新增按钮操作 */
function handleAdd() {
  reset();
  open.value = true;
  title.value = "添加产品";
}

/** 修改按钮操作 */
function handleUpdate(row) {
  reset();
  const id = row.id || ids.value;
  request({
      url: '/demo/product/' + id,
      method: 'get'
  }).then(response => {
    form.value = response.data;open.value = true;
    title.value = "修改产品";
  });
}

/** 提交按钮 */
function submitForm() {
  proxy.$refs["productRef"].validate(valid => {
    if (valid) {let url = '/demo/product';
      if (form.value.id != null) {
        request({ url: url, method: 'put', data: form.value}).then(response => {
                proxy.$message.success("修改成功");
                open.value = false;
                getList();
            });
      } else {
        request({ url: url,method: 'post',data: form.value}).then(response => {
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
  const ids = row.id || ids.value;
  proxy.$confirm('是否确认删除产品编号为"' + ids + '"的数据项？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    return delIotProduct(ids);
  }).then(() => {
    getList();
    proxy.$message.success("删除成功");
  }).catch(() => {});
}

/** 导出按钮操作 */
function handleExport() {
  proxy.download('demo/product/exportIotProduct', {
    ...queryParams.value
  }, 'product_' + Date.now() + '.xlsx');
}
// 初始化加载数据
getList();
</script>
