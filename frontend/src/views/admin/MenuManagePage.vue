<template>
  <div class="menu-manage-page">
    <div class="page-header">
      <h2>菜单管理</h2>
      <el-button
        type="primary"
        @click="showAddDialog"
        :icon="Plus"
      >
        添加商品
      </el-button>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-bar">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchForm.keyword"
            placeholder="搜索商品名称"
            :prefix-icon="Search"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="searchForm.category"
            placeholder="选择分类"
            clearable
            @change="handleSearch"
          >
            <el-option label="全部分类" value="" />
            <el-option label="咖啡" value="coffee" />
            <el-option label="茶饮" value="tea" />
            <el-option label="甜点" value="pastry" />
            <el-option label="三明治" value="sandwich" />
            <el-option label="果汁" value="juice" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select
            v-model="searchForm.available"
            placeholder="商品状态"
            clearable
            @change="handleSearch"
          >
            <el-option label="全部状态" value="" />
            <el-option label="上架" value="true" />
            <el-option label="下架" value="false" />
          </el-select>
        </el-col>
        <el-col :span="3">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 菜单列表 -->
    <el-table
      :data="menuItems"
      v-loading="loading"
      style="width: 100%"
    >
      <el-table-column prop="image_url" label="图片" width="80">
        <template #default="scope">
          <div v-if="scope.row.image_url">
            <el-image
              :src="scope.row.image_url"
              style="width: 60px; height: 60px; border-radius: 4px;"
              fit="cover"
              @error="handleImageError"
            />
          </div>
          <div v-else class="admin-image-placeholder">
            <el-icon><Picture /></el-icon>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="name" label="商品名称" width="150" />
      <el-table-column prop="category" label="分类" width="100">
        <template #default="scope">
          <el-tag :type="getCategoryType(scope.row.category)">
            {{ getCategoryLabel(scope.row.category) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="100">
        <template #default="scope">
          ¥{{ scope.row.price }}
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column prop="is_available" label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.is_available ? 'success' : 'danger'">
            {{ scope.row.is_available ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button
            size="small"
            type="primary"
            @click="handleEdit(scope.row)"
          >
            编辑
          </el-button>
          <el-button
            size="small"
            :type="scope.row.is_available ? 'warning' : 'success'"
            @click="handleToggleAvailable(scope.row)"
          >
            {{ scope.row.is_available ? '下架' : '上架' }}
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑商品' : '添加商品'"
      width="600px"
      @close="resetForm"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类">
            <el-option label="咖啡" value="coffee" />
            <el-option label="茶饮" value="tea" />
            <el-option label="甜点" value="pastry" />
            <el-option label="三明治" value="sandwich" />
            <el-option label="果汁" value="juice" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格" prop="price">
          <el-input-number
            v-model="form.price"
            :min="0.01"
            :max="999999.99"
            :precision="2"
            placeholder="请输入价格"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="请输入商品描述"
          />
        </el-form-item>
        <el-form-item label="商品图片" prop="image_url">
          <el-input v-model="form.image_url" placeholder="请输入图片URL" />
          <div class="image-preview" v-if="form.image_url">
            <el-image
              :src="form.image_url"
              style="width: 100px; height: 100px; border-radius: 4px; margin-top: 10px;"
              fit="cover"
            />
          </div>
        </el-form-item>
        <el-form-item label="状态" prop="is_available">
          <el-switch
            v-model="form.is_available"
            active-text="上架"
            inactive-text="下架"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '更新' : '添加' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Picture } from '@element-plus/icons-vue'
import { menuAPI } from '@/api/menu'
import { formatDate } from '@/utils/date'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref()

// 搜索表单
const searchForm = reactive({
  keyword: '',
  category: '',
  available: ''
})

// 表单数据
const form = reactive({
  name: '',
  category: '',
  price: 0,
  description: '',
  image_url: '',
  is_available: true
})

// 分页数据
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 菜单列表
const menuItems = ref([])

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' },
    { max: 100, message: '商品名称不能超过100个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '价格必须大于0', trigger: 'blur' }
  ]
}

// 获取菜单列表
const fetchMenuItems = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize.value,
      available_only: searchForm.available !== '' ? searchForm.available === 'true' : false
    }

    if (searchForm.keyword) {
      params.keyword = searchForm.keyword
    }

    if (searchForm.category) {
      params.category = searchForm.category
    }

    const response = await menuAPI.getMenuItems(params)
    if (response.success) {
      menuItems.value = response.data.items
      total.value = response.data.total
    } else {
      ElMessage.error('获取菜单列表失败')
    }
  } catch (error) {
    console.error('获取菜单列表失败:', error)
    ElMessage.error('获取菜单列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
  fetchMenuItems()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1
  fetchMenuItems()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  fetchMenuItems()
}

// 显示添加对话框
const showAddDialog = () => {
  isEdit.value = false
  dialogVisible.value = true
  resetForm()
}

// 重置表单
const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  Object.assign(form, {
    name: '',
    category: '',
    price: 0,
    description: '',
    image_url: '',
    is_available: true
  })
  isEdit.value = false
}

// 编辑处理
const handleEdit = (row) => {
  isEdit.value = true
  dialogVisible.value = true
  Object.assign(form, {
    id: row.id,
    name: row.name,
    category: row.category,
    price: row.price,
    description: row.description,
    image_url: row.image_url,
    is_available: row.is_available
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    let response
    if (isEdit.value) {
      response = await menuAPI.updateMenuItem(form.id, form)
    } else {
      response = await menuAPI.createMenuItem(form)
    }

    if (response.success) {
      ElMessage.success(isEdit.value ? '商品更新成功' : '商品添加成功')
      dialogVisible.value = false
      fetchMenuItems()
    } else {
      ElMessage.error(response.errors?.join(', ') || '操作失败')
    }
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

// 切换可用状态
const handleToggleAvailable = async (row) => {
  try {
    const newStatus = !row.is_available
    const action = newStatus ? '上架' : '下架'

    await ElMessageBox.confirm(
      `确定要${action}商品"${row.name}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await menuAPI.updateMenuItem(row.id, {
      is_available: newStatus
    })

    if (response.success) {
      ElMessage.success(`商品${action}成功`)
      fetchMenuItems()
    } else {
      ElMessage.error(`${action}失败`)
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新状态失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 删除处理
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除商品"${row.name}"吗？此操作不可撤销！`,
      '确认删除',
      {
        confirmButtonText: '确定删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    const response = await menuAPI.deleteMenuItem(row.id)
    if (response.success) {
      ElMessage.success('商品删除成功')
      fetchMenuItems()
    } else {
      ElMessage.error('删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 获取分类标签
const getCategoryLabel = (category) => {
  const labels = {
    coffee: '咖啡',
    tea: '茶饮',
    pastry: '甜点',
    sandwich: '三明治',
    juice: '果汁'
  }
  return labels[category] || category
}

// 获取分类标签类型
const getCategoryType = (category) => {
  const types = {
    coffee: '',
    tea: 'success',
    pastry: 'warning',
    sandwich: 'info',
    juice: 'primary'
  }
  return types[category] || ''
}

// 处理图片加载错误
const handleImageError = (event) => {
  event.target.style.display = 'none'
  const parent = event.target.parentElement
  if (parent) {
    const placeholder = document.createElement('div')
    placeholder.className = 'admin-image-placeholder'
    placeholder.innerHTML = '<i class="el-icon-picture"></i>'
    parent.appendChild(placeholder)
  }
}

// 组件挂载
onMounted(() => {
  fetchMenuItems()
})
</script>

<style lang="scss" scoped>
.menu-manage-page {
  padding: 20px;
  background: #f8f9fa;
  min-height: 100vh;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  h2 {
    margin: 0;
    color: #2c3e50;
    font-weight: 600;
  }
}

.search-bar {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.el-table {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.image-preview {
  margin-top: 10px;
}

.admin-image-placeholder {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #f0f0f0 0%, #d0d0d0 100%);
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 24px;
}

:deep(.el-dialog__body) {
  padding: 20px;
}
</style>