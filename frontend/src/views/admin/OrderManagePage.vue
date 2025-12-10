<template>
  <div class="order-manage-page">
    <div class="page-header">
      <h2>订单管理</h2>
      <div class="header-actions">
        <el-select v-model="statusFilter" placeholder="订单状态" clearable style="width: 120px; margin-right: 12px;">
          <el-option label="全部" value="" />
          <el-option label="待处理" value="pending" />
          <el-option label="制作中" value="preparing" />
          <el-option label="待取餐" value="ready" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>
        <el-button @click="fetchOrders" :icon="Refresh" :loading="loading">刷新</el-button>
      </div>
    </div>

    <el-card>
      <el-table :data="orders" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="order_number" label="订单号" width="140" />
        <el-table-column prop="pickup_code" label="取餐码" width="100">
          <template #default="{ row }">
            <el-tag type="success" size="large" effect="dark">{{ row.pickup_code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="total_price" label="金额" width="100">
          <template #default="{ row }">¥{{ (row.total_price || 0).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="item_count" label="商品数" width="80" align="center" />
        <el-table-column prop="created_at" label="下单时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewOrder(row)">详情</el-button>
            <el-dropdown trigger="click" @command="(cmd) => updateStatus(row.id, cmd)">
              <el-button size="small" type="primary">更新状态</el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="pending">待处理</el-dropdown-item>
                  <el-dropdown-item command="preparing">制作中</el-dropdown-item>
                  <el-dropdown-item command="ready">待取餐</el-dropdown-item>
                  <el-dropdown-item command="completed">已完成</el-dropdown-item>
                  <el-dropdown-item command="cancelled" divided>取消订单</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchOrders"
          @current-change="fetchOrders"
        />
      </div>
    </el-card>

    <!-- 订单详情对话框 -->
    <el-dialog v-model="showDetail" title="订单详情" width="600px">
      <div v-if="currentOrder" class="order-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="订单号">{{ currentOrder.order_number }}</el-descriptions-item>
          <el-descriptions-item label="取餐码">
            <el-tag type="success" size="large" effect="dark">{{ currentOrder.pickup_code }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="订单金额">¥{{ (currentOrder.total_price || 0).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item label="订单状态">
            <el-tag :type="getStatusType(currentOrder.status)">{{ getStatusText(currentOrder.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="下单时间" :span="2">{{ formatTime(currentOrder.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ currentOrder.notes || '无' }}</el-descriptions-item>
        </el-descriptions>
        <h4 style="margin: 16px 0 8px;">订单商品</h4>
        <el-table :data="currentOrder.items || []" size="small">
          <el-table-column prop="menu_name" label="商品" />
          <el-table-column prop="quantity" label="数量" width="80" align="center" />
          <el-table-column prop="unit_price" label="单价" width="100">
            <template #default="{ row }">¥{{ (row.unit_price || 0).toFixed(2) }}</template>
          </el-table-column>
          <el-table-column prop="subtotal" label="小计" width="100">
            <template #default="{ row }">¥{{ (row.subtotal || 0).toFixed(2) }}</template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import request from '@/api/index'

const loading = ref(false)
const orders = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const statusFilter = ref('')
const showDetail = ref(false)
const currentOrder = ref(null)

const statusMap = {
  pending: { text: '待处理', type: 'warning' },
  preparing: { text: '制作中', type: 'primary' },
  ready: { text: '待取餐', type: 'success' },
  completed: { text: '已完成', type: 'info' },
  cancelled: { text: '已取消', type: 'danger' }
}

const getStatusText = (status) => statusMap[status]?.text || status
const getStatusType = (status) => statusMap[status]?.type || 'info'

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const fetchOrders = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: pageSize.value }
    if (statusFilter.value) params.status = statusFilter.value
    const response = await request.get('/admin/orders', params)
    if (response.success && response.data) {
      orders.value = response.data.orders || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    console.error('获取订单失败:', error)
    ElMessage.error('获取订单失败')
  } finally {
    loading.value = false
  }
}

const viewOrder = async (order) => {
  try {
    const response = await request.get(`/orders/${order.id}`)
    if (response.success && response.data) {
      currentOrder.value = response.data
      showDetail.value = true
    }
  } catch (error) {
    // 如果详情API失败，使用列表数据
    currentOrder.value = order
    showDetail.value = true
  }
}

const updateStatus = async (id, status) => {
  try {
    const response = await request.put(`/admin/orders/${id}/status`, { status })
    if (response.success) {
      ElMessage.success('状态更新成功')
      fetchOrders()
    } else {
      ElMessage.error(response.message || '更新失败')
    }
  } catch (error) {
    ElMessage.error('更新状态失败')
  }
}

watch(statusFilter, () => {
  currentPage.value = 1
  fetchOrders()
})

onMounted(() => fetchOrders())
</script>

<style lang="scss" scoped>
.order-manage-page {
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
  h2 { margin: 0; color: #2c3e50; }
}
.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
.order-detail h4 { color: #2c3e50; }
</style>
