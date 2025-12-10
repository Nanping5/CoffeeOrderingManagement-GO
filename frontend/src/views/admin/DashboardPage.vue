<template>
  <div class="dashboard-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-info">
        <h1>仪表板</h1>
        <p>欢迎回来，{{ adminName }}！查看您的咖啡店运营数据</p>
      </div>
      <el-button @click="refreshData" :loading="loading" class="refresh-btn">
        <el-icon><Refresh /></el-icon>
        刷新数据
      </el-button>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card blue">
        <div class="stat-content">
          <div class="stat-info">
            <span class="stat-label">今日订单</span>
            <span class="stat-value">{{ orderStats.today_orders }}</span>
          </div>
          <div class="stat-icon"><el-icon><Document /></el-icon></div>
        </div>
        <div class="stat-footer"><span class="trend up">+12%</span><span class="trend-label">较昨日</span></div>
      </div>
      <div class="stat-card green">
        <div class="stat-content">
          <div class="stat-info">
            <span class="stat-label">今日收入</span>
            <span class="stat-value">¥{{ orderStats.today_revenue?.toFixed(2) || '0.00' }}</span>
          </div>
          <div class="stat-icon"><el-icon><Money /></el-icon></div>
        </div>
        <div class="stat-footer"><span class="trend up">+8%</span><span class="trend-label">较昨日</span></div>
      </div>
      <div class="stat-card orange">
        <div class="stat-content">
          <div class="stat-info">
            <span class="stat-label">总订单数</span>
            <span class="stat-value">{{ orderStats.total_orders }}</span>
          </div>
          <div class="stat-icon"><el-icon><Tickets /></el-icon></div>
        </div>
        <div class="stat-footer"><span class="trend-label">累计数据</span></div>
      </div>
      <div class="stat-card purple">
        <div class="stat-content">
          <div class="stat-info">
            <span class="stat-label">总收入</span>
            <span class="stat-value">¥{{ orderStats.total_revenue?.toFixed(2) || '0.00' }}</span>
          </div>
          <div class="stat-icon"><el-icon><Coin /></el-icon></div>
        </div>
        <div class="stat-footer"><span class="trend-label">累计数据</span></div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="main-content">
      <div class="content-left">
        <!-- 订单状态 -->
        <div class="card">
          <div class="card-header"><h3><el-icon><List /></el-icon> 订单状态</h3></div>
          <div class="card-body">
            <div class="status-grid">
              <div class="status-item pending"><span class="status-count">{{ orderStats.pending_count || 0 }}</span><span class="status-label">待处理</span></div>
              <div class="status-item preparing"><span class="status-count">{{ orderStats.preparing_count || 0 }}</span><span class="status-label">制作中</span></div>
              <div class="status-item ready"><span class="status-count">{{ orderStats.ready_count || 0 }}</span><span class="status-label">待取餐</span></div>
              <div class="status-item completed"><span class="status-count">{{ orderStats.completed_count || 0 }}</span><span class="status-label">已完成</span></div>
              <div class="status-item cancelled"><span class="status-count">{{ orderStats.cancelled_count || 0 }}</span><span class="status-label">已取消</span></div>
            </div>
          </div>
        </div>
        <!-- 热销商品 -->
        <div class="card" v-if="topProducts.length > 0">
          <div class="card-header"><h3><el-icon><TrendCharts /></el-icon> 热销商品 Top 5</h3></div>
          <div class="card-body">
            <div class="product-list">
              <div v-for="(product, index) in topProducts" :key="index" class="product-item">
                <div class="product-rank" :class="'rank-' + (index + 1)">{{ index + 1 }}</div>
                <div class="product-info"><span class="product-name">{{ product.menu_name }}</span><span class="product-sales">销量: {{ product.quantity }}</span></div>
                <div class="product-revenue">¥{{ (product.revenue || 0).toFixed(2) }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="content-right">
        <!-- 会员统计 -->
        <div class="card">
          <div class="card-header"><h3><el-icon><User /></el-icon> 会员统计</h3></div>
          <div class="card-body">
            <div class="member-total"><span class="total-value">{{ orderStats.total_users || 0 }}</span><span class="total-label">注册会员</span></div>
            <div class="member-levels">
              <div class="level-item"><div class="level-badge bronze">铜</div><span class="level-count">{{ memberLevels.bronze || 0 }}</span></div>
              <div class="level-item"><div class="level-badge silver">银</div><span class="level-count">{{ memberLevels.silver || 0 }}</span></div>
              <div class="level-item"><div class="level-badge gold">金</div><span class="level-count">{{ memberLevels.gold || 0 }}</span></div>
              <div class="level-item"><div class="level-badge platinum">白金</div><span class="level-count">{{ memberLevels.platinum || 0 }}</span></div>
            </div>
          </div>
        </div>
        <!-- 菜单统计 -->
        <div class="card">
          <div class="card-header"><h3><el-icon><Menu /></el-icon> 菜单统计</h3></div>
          <div class="card-body">
            <div class="menu-stats">
              <div class="menu-stat-item"><span class="stat-num">{{ menuStats.total }}</span><span class="stat-text">菜单总数</span></div>
              <div class="menu-stat-item available"><span class="stat-num">{{ menuStats.available }}</span><span class="stat-text">已上架</span></div>
              <div class="menu-stat-item unavailable"><span class="stat-num">{{ menuStats.unavailable }}</span><span class="stat-text">已下架</span></div>
            </div>
          </div>
        </div>
        <!-- 快速操作 -->
        <div class="card">
          <div class="card-header"><h3><el-icon><Operation /></el-icon> 快速操作</h3></div>
          <div class="card-body">
            <div class="action-buttons">
              <button class="action-btn primary" @click="goToMenuManage"><el-icon><Menu /></el-icon><span>菜单管理</span></button>
              <button class="action-btn warning" @click="goToOrderManage"><el-icon><List /></el-icon><span>订单管理</span></button>
              <button class="action-btn success" @click="goToAddMenu"><el-icon><Plus /></el-icon><span>添加菜品</span></button>
              <button class="action-btn info" @click="goToFrontend"><el-icon><View /></el-icon><span>查看前台</span></button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Refresh, Menu, Plus, View, Document, Money, Tickets, Coin, List, User, TrendCharts, Operation } from '@element-plus/icons-vue'
import { menuAPI } from '@/api/menu'
import request from '@/api/index'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const menuStats = reactive({ total: 0, available: 0, unavailable: 0 })
const orderStats = reactive({ total_orders: 0, total_revenue: 0, today_orders: 0, today_revenue: 0, pending_count: 0, preparing_count: 0, ready_count: 0, completed_count: 0, cancelled_count: 0, total_users: 0 })
const memberLevels = reactive({ bronze: 0, silver: 0, gold: 0, platinum: 0 })
const topProducts = ref([])
const adminName = computed(() => authStore.user?.username || '管理员')

const fetchMenuStats = async () => {
  try {
    const response = await menuAPI.getMenuItems({ per_page: 1000 })
    if (response.success) {
      const items = response.data?.items || []
      menuStats.total = items.length
      menuStats.available = items.filter(item => item.is_available).length
      menuStats.unavailable = items.filter(item => !item.is_available).length
    }
  } catch (error) { console.error('获取菜单统计失败:', error) }
}

const fetchOrderStats = async () => {
  try {
    const response = await request.get('/admin/orders/statistics')
    if (response.success && response.data) {
      Object.assign(orderStats, response.data)
      if (response.data.member_levels) Object.assign(memberLevels, response.data.member_levels)
      if (response.data.top_products) topProducts.value = response.data.top_products
    }
  } catch (error) { console.error('获取订单统计失败:', error) }
}

const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([fetchMenuStats(), fetchOrderStats()])
    ElMessage.success('数据刷新成功')
  } catch (error) { ElMessage.error('数据刷新失败') }
  finally { loading.value = false }
}

const goToMenuManage = () => router.push('/admin/menu')
const goToOrderManage = () => router.push('/admin/orders')
const goToAddMenu = () => router.push({ name: 'AdminMenu', query: { action: 'add' } })
const goToFrontend = () => window.open('/menu', '_blank')

onMounted(() => { fetchMenuStats(); fetchOrderStats() })
</script>

<style lang="scss" scoped>
.dashboard-page { min-height: 100%; }
.page-header {
  display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 24px;
  .header-info {
    h1 { margin: 0 0 8px 0; font-size: 28px; font-weight: 700; color: var(--admin-text-primary); }
    p { margin: 0; font-size: 14px; color: var(--admin-text-secondary); }
  }
  .refresh-btn { background: var(--admin-btn-bg); border: 1px solid var(--admin-card-border); color: var(--admin-text-primary); border-radius: 10px; padding: 12px 20px;
    &:hover { background: var(--admin-btn-hover); }
  }
}
.stats-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 24px;
  @media (max-width: 1200px) { grid-template-columns: repeat(2, 1fr); }
  @media (max-width: 600px) { grid-template-columns: 1fr; }
}
.stat-card {
  background: var(--admin-card-bg); border-radius: 16px; padding: 24px; border: 1px solid var(--admin-card-border); transition: all 0.3s ease;
  &:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1); }
  &.blue .stat-icon { background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%); }
  &.green .stat-icon { background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); }
  &.orange .stat-icon { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); }
  &.purple .stat-icon { background: linear-gradient(135deg, #a855f7 0%, #7c3aed 100%); }
  .stat-content { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px; }
  .stat-info {
    .stat-label { display: block; font-size: 14px; color: var(--admin-text-secondary); margin-bottom: 8px; }
    .stat-value { display: block; font-size: 28px; font-weight: 700; color: var(--admin-text-primary); }
  }
  .stat-icon { width: 48px; height: 48px; border-radius: 12px; display: flex; align-items: center; justify-content: center; color: white; font-size: 24px; }
  .stat-footer { display: flex; align-items: center; gap: 8px;
    .trend { font-size: 13px; font-weight: 600; &.up { color: #22c55e; } &.down { color: #ef4444; } }
    .trend-label { font-size: 12px; color: var(--admin-text-muted); }
  }
}
.main-content { display: grid; grid-template-columns: 1fr 380px; gap: 24px;
  @media (max-width: 1200px) { grid-template-columns: 1fr; }
}
.content-left, .content-right { display: flex; flex-direction: column; gap: 24px; }
.card {
  background: var(--admin-card-bg); border-radius: 16px; border: 1px solid var(--admin-card-border); overflow: hidden;
  .card-header { padding: 20px 24px; border-bottom: 1px solid var(--admin-card-border); background: var(--admin-card-header-bg);
    h3 { margin: 0; font-size: 16px; font-weight: 600; color: var(--admin-text-primary); display: flex; align-items: center; gap: 10px;
      .el-icon { color: #f59e0b; }
    }
  }
  .card-body { padding: 24px; }
}

.status-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 16px;
  @media (max-width: 768px) { grid-template-columns: repeat(3, 1fr); }
  .status-item { text-align: center; padding: 16px; border-radius: 12px; background: var(--admin-item-bg);
    .status-count { display: block; font-size: 28px; font-weight: 700; margin-bottom: 4px; }
    .status-label { font-size: 12px; color: var(--admin-text-secondary); }
    &.pending .status-count { color: #f59e0b; }
    &.preparing .status-count { color: #3b82f6; }
    &.ready .status-count { color: #22c55e; }
    &.completed .status-count { color: #6b7280; }
    &.cancelled .status-count { color: #ef4444; }
  }
}
.product-list {
  .product-item { display: flex; align-items: center; gap: 16px; padding: 16px 0; border-bottom: 1px solid var(--admin-card-border);
    &:last-child { border-bottom: none; }
    .product-rank { width: 32px; height: 32px; border-radius: 8px; display: flex; align-items: center; justify-content: center; font-weight: 700; font-size: 14px; background: var(--admin-item-bg); color: var(--admin-text-secondary);
      &.rank-1 { background: linear-gradient(135deg, #ffd700 0%, #ffb700 100%); color: #000; }
      &.rank-2 { background: linear-gradient(135deg, #c0c0c0 0%, #a0a0a0 100%); color: #000; }
      &.rank-3 { background: linear-gradient(135deg, #cd7f32 0%, #b87333 100%); color: #fff; }
    }
    .product-info { flex: 1;
      .product-name { display: block; font-size: 14px; font-weight: 500; color: var(--admin-text-primary); margin-bottom: 4px; }
      .product-sales { font-size: 12px; color: var(--admin-text-muted); }
    }
    .product-revenue { font-size: 16px; font-weight: 600; color: #22c55e; }
  }
}
.member-total { text-align: center; padding: 24px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px; margin-bottom: 20px;
  .total-value { display: block; font-size: 40px; font-weight: 700; color: white; }
  .total-label { font-size: 14px; color: rgba(255, 255, 255, 0.8); }
}
.member-levels { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px;
  .level-item { text-align: center;
    .level-badge { width: 40px; height: 40px; border-radius: 50%; display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 600; margin: 0 auto 8px;
      &.bronze { background: #cd7f32; color: white; }
      &.silver { background: #c0c0c0; color: #333; }
      &.gold { background: #ffd700; color: #333; }
      &.platinum { background: linear-gradient(135deg, #e5e4e2 0%, #a0a0a0 100%); color: #333; }
    }
    .level-count { font-size: 18px; font-weight: 600; color: var(--admin-text-primary); }
  }
}
.menu-stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px;
  .menu-stat-item { text-align: center; padding: 16px; border-radius: 12px; background: var(--admin-item-bg);
    .stat-num { display: block; font-size: 28px; font-weight: 700; color: var(--admin-text-primary); margin-bottom: 4px; }
    .stat-text { font-size: 12px; color: var(--admin-text-secondary); }
    &.available .stat-num { color: #22c55e; }
    &.unavailable .stat-num { color: #ef4444; }
  }
}
.action-buttons { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px;
  .action-btn { display: flex; align-items: center; justify-content: center; gap: 8px; padding: 14px 16px; border: none; border-radius: 10px; font-size: 14px; font-weight: 500; cursor: pointer; transition: all 0.2s ease;
    &.primary { background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%); color: white; }
    &.warning { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); color: white; }
    &.success { background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); color: white; }
    &.info { background: var(--admin-btn-bg); color: var(--admin-text-primary); border: 1px solid var(--admin-card-border); }
    &:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15); }
  }
}
</style>
