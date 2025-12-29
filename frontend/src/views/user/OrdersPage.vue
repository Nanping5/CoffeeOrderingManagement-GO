<template>
  <div class="orders-page">
    <div class="page-header">
      <h1 class="page-title gradient-text">
        <el-icon><List /></el-icon>
        我的订单
      </h1>
      <p class="page-subtitle">查看您的历史订单</p>
    </div>

    <div v-if="authStore.isAuthenticated" class="orders-content">
      <!-- 筛选栏 -->
      <div class="filter-bar glass-panel">
        <el-radio-group v-model="statusFilter" @change="fetchOrders">
          <el-radio-button label="">全部</el-radio-button>
          <el-radio-button label="pending">待处理</el-radio-button>
          <el-radio-button label="preparing">制作中</el-radio-button>
          <el-radio-button label="ready">待取餐</el-radio-button>
          <el-radio-button label="completed">已完成</el-radio-button>
        </el-radio-group>
      </div>

      <!-- 订单列表 -->
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else-if="orders.length > 0" class="orders-list">
        <el-card v-for="order in orders" :key="order.id" class="order-card glass-panel">
          <div class="order-header">
            <div class="order-info">
              <span class="order-number">订单号: {{ order.order_number }}</span>
              <span class="order-time">{{ formatTime(order.created_at) }}</span>
            </div>
            <el-tag :type="getStatusType(order.status)" size="small">
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>

          <div class="pickup-code-display" v-if="order.status !== 'completed' && order.status !== 'cancelled'">
            <span class="label">取餐码</span>
            <span class="code">{{ order.pickup_code }}</span>
          </div>

          <div class="order-items">
            <div v-for="item in order.items" :key="item.id" class="order-item">
              <span class="item-name">{{ item.menu_name }}</span>
              <span class="item-quantity">x{{ item.quantity }}</span>
              <span class="item-price">¥{{ item.subtotal?.toFixed(2) }}</span>
            </div>
          </div>

          <div class="order-footer">
            <div class="order-summary">
              <span v-if="order.points_used > 0" class="points-used">
                积分抵扣: -{{ order.points_used }}
              </span>
              <span v-if="order.points_earned > 0" class="points-earned">
                获得积分: +{{ order.points_earned }}
              </span>
            </div>
            <div class="order-total">
              共 {{ order.item_count }} 件，实付 <span class="price">¥{{ order.final_payment_amount?.toFixed(2) || order.total_price?.toFixed(2) }}</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 空状态 -->
      <div v-else class="empty-state">
        <el-icon :size="80" class="empty-icon"><List /></el-icon>
        <p>暂无订单记录</p>
        <el-button type="primary" @click="goToMenu">去点餐</el-button>
      </div>

      <!-- 分页 -->
      <div v-if="pagination.pages > 1" class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          :page-size="pagination.per_page"
          :total="pagination.total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- 未登录状态 -->
    <div v-else class="not-logged-in">
      <el-icon :size="80" class="empty-icon"><List /></el-icon>
      <h3>请先登录</h3>
      <p>登录后即可查看您的订单</p>
      <el-button type="primary" @click="goToLogin" size="large">立即登录</el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { List } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import orderAPI from '@/api/order'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const orders = ref([])
const statusFilter = ref('')
const pagination = reactive({
  page: 1,
  per_page: 10,
  total: 0,
  pages: 0
})

const fetchOrders = async () => {
  if (!authStore.isAuthenticated) return
  
  try {
    loading.value = true
    const response = await orderAPI.getMyOrders({
      page: pagination.page,
      per_page: pagination.per_page,
      status: statusFilter.value || undefined
    })
    
    if (response.success && response.data) {
      orders.value = response.data.orders || []
      pagination.total = response.data.total || 0
      pagination.pages = response.data.pages || 0
    }
  } catch (error) {
    console.error('获取订单失败:', error)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  pagination.page = page
  fetchOrders()
}

const getStatusType = (status) => {
  const types = {
    pending: 'warning',
    preparing: 'primary',
    ready: 'success',
    completed: 'info',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待处理',
    preparing: '制作中',
    ready: '待取餐',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const formatTime = (time) => {
  if (!time) return ''
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const goToMenu = () => router.push('/menu')
const goToLogin = () => router.push('/auth/login')

onMounted(() => {
  if (authStore.isAuthenticated) {
    fetchOrders()
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.orders-page {
  min-height: 100vh;
  padding: $spacing-2xl $spacing-lg;
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: $spacing-2xl;

  .page-title {
    font-size: $font-size-4xl;
    font-weight: $font-weight-bold;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: $spacing-md;
    margin-bottom: $spacing-sm;
  }

  .page-subtitle {
    font-size: $font-size-lg;
    color: $text-secondary;
  }
}

.filter-bar {
  padding: $spacing-md;
  margin-bottom: $spacing-lg;
  border-radius: $radius-lg;
  text-align: center;
}

.orders-list {
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.order-card {
  border-radius: $radius-lg;
  
  .order-header {
    @include flex-between;
    margin-bottom: $spacing-md;
    padding-bottom: $spacing-md;
    border-bottom: 1px solid rgba(0, 0, 0, 0.1);

    .order-info {
      display: flex;
      flex-direction: column;
      gap: $spacing-xs;

      .order-number {
        font-weight: $font-weight-semibold;
        color: $text-primary;
      }

      .order-time {
        font-size: $font-size-sm;
        color: $text-secondary;
      }
    }
  }

  .pickup-code-display {
    background: linear-gradient(135deg, $primary-color, darken($primary-color, 10%));
    border-radius: $radius-md;
    padding: $spacing-md;
    text-align: center;
    margin-bottom: $spacing-md;
    color: white;

    .label {
      display: block;
      font-size: $font-size-sm;
      opacity: 0.9;
      margin-bottom: $spacing-xs;
    }

    .code {
      font-size: $font-size-2xl;
      font-weight: $font-weight-bold;
      letter-spacing: 4px;
    }
  }

  .order-items {
    margin-bottom: $spacing-md;

    .order-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: $spacing-sm 0;
      border-bottom: 1px dashed rgba(0, 0, 0, 0.05);

      &:last-child {
        border-bottom: none;
      }

      .item-name {
        flex: 1;
        color: $text-primary;
      }

      .item-quantity {
        color: $text-secondary;
        margin: 0 $spacing-md;
      }

      .item-price {
        color: $danger-color;
        font-weight: $font-weight-medium;
      }
    }
  }

  .order-footer {
    @include flex-between;
    padding-top: $spacing-md;
    border-top: 1px solid rgba(0, 0, 0, 0.1);

    .order-summary {
      display: flex;
      gap: $spacing-md;
      font-size: $font-size-sm;

      .points-used {
        color: $warning-color;
      }

      .points-earned {
        color: $success-color;
      }
    }

    .order-total {
      color: $text-secondary;

      .price {
        color: $danger-color;
        font-size: $font-size-lg;
        font-weight: $font-weight-bold;
      }
    }
  }
}

.empty-state, .not-logged-in {
  @include flex-center;
  flex-direction: column;
  gap: $spacing-lg;
  min-height: 50vh;
  text-align: center;

  .empty-icon { color: $text-light; opacity: 0.5; }
  h3 { margin: 0; color: $text-primary; }
  p { margin: 0; color: $text-secondary; }
}

.loading-state {
  padding: $spacing-xl;
}

.pagination-wrapper {
  @include flex-center;
  margin-top: $spacing-xl;
}

// 移动端适配
@include respond-to(md) {
  .orders-page {
    padding: $spacing-lg $spacing-sm;
  }
  
  .page-header {
    margin-bottom: $spacing-lg;
    
    .page-title {
      font-size: $font-size-2xl;
      gap: $spacing-sm;
    }
    
    .page-subtitle {
      font-size: $font-size-base;
    }
  }
  
  .filter-bar {
    padding: $spacing-sm;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    
    :deep(.el-radio-group) {
      display: flex;
      flex-wrap: nowrap;
      
      .el-radio-button {
        flex-shrink: 0;
        
        .el-radio-button__inner {
          padding: $spacing-xs $spacing-sm;
          font-size: $font-size-xs;
        }
      }
    }
  }
  
  .order-card {
    :deep(.el-card__body) {
      padding: $spacing-md;
    }
    
    .order-header {
      flex-direction: column;
      align-items: flex-start;
      gap: $spacing-sm;
      
      .order-info {
        width: 100%;
        
        .order-number {
          font-size: $font-size-sm;
        }
        
        .order-time {
          font-size: $font-size-xs;
        }
      }
    }
    
    .pickup-code-display {
      padding: $spacing-md;
      
      .code {
        font-size: $font-size-xl;
        letter-spacing: 2px;
      }
    }
    
    .order-items .order-item {
      flex-wrap: wrap;
      gap: $spacing-xs;
      
      .item-name {
        flex: 1 1 100%;
        font-size: $font-size-sm;
      }
      
      .item-quantity, .item-price {
        font-size: $font-size-sm;
      }
    }
    
    .order-footer {
      flex-direction: column;
      align-items: flex-start;
      gap: $spacing-sm;
      
      .order-summary {
        flex-wrap: wrap;
        font-size: $font-size-xs;
      }
      
      .order-total {
        width: 100%;
        text-align: right;
        font-size: $font-size-sm;
        
        .price {
          font-size: $font-size-base;
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .orders-page {
    padding: $spacing-md $spacing-xs;
  }
  
  .page-header .page-title {
    font-size: $font-size-xl;
  }
  
  .filter-bar {
    :deep(.el-radio-button .el-radio-button__inner) {
      padding: 6px 10px;
      font-size: 11px;
    }
  }
}
</style>
