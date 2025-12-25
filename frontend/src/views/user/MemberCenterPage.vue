<template>
  <div class="member-center-page">
    <div class="page-header">
      <h1 class="page-title gradient-text">
        <el-icon><Trophy /></el-icon>
        会员中心
      </h1>
      <p class="page-subtitle">管理您的会员权益和积分</p>
    </div>

    <div v-if="authStore.isAuthenticated" class="member-content">
      <!-- 会员等级卡片 -->
      <el-card class="member-level-card glass-panel">
        <div class="level-header">
          <div class="level-info">
            <div class="level-badge" :class="pointsStore.memberLevel">
              <el-icon><Medal /></el-icon>
              <span>{{ pointsStore.memberLevelName }}</span>
            </div>
            <div class="user-info">
              <h3>{{ authStore.userInfo?.username || '会员' }}</h3>
              <p>{{ authStore.userInfo?.email }}</p>
            </div>
          </div>
          <div v-if="pointsStore.nextLevel" class="level-progress">
            <p class="progress-text">
              距离 {{ getLevelName(pointsStore.nextLevel.name) }} 还需 
              <strong>{{ pointsStore.nextLevel.points_needed }}</strong> 积分
            </p>
            <el-progress
              :percentage="levelProgress"
              :stroke-width="10"
              :color="levelProgressColor"
            />
          </div>
        </div>
      </el-card>

      <!-- 积分信息卡片 -->
      <el-row :gutter="20" class="points-cards">
        <el-col :xs="24" :sm="8">
          <el-card class="points-card available glass-panel">
            <div class="points-content">
              <el-icon class="points-icon"><Wallet /></el-icon>
              <div class="points-info">
                <div class="points-value">{{ pointsStore.availablePoints }}</div>
                <div class="points-label">可用积分</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card class="points-card total glass-panel">
            <div class="points-content">
              <el-icon class="points-icon"><Coin /></el-icon>
              <div class="points-info">
                <div class="points-value">{{ pointsStore.totalPoints }}</div>
                <div class="points-label">总积分</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :xs="24" :sm="8">
          <el-card class="points-card lifetime glass-panel">
            <div class="points-content">
              <el-icon class="points-icon"><TrendCharts /></el-icon>
              <div class="points-info">
                <div class="points-value">{{ pointsStore.lifetimePoints }}</div>
                <div class="points-label">历史累计</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 快捷操作 -->
      <el-card class="quick-actions glass-panel">
        <h3>快捷操作</h3>
        <el-row :gutter="16">
          <el-col :xs="12" :sm="6">
            <el-button type="primary" plain @click="goToProfile" class="action-btn">
              <el-icon><User /></el-icon>
              个人资料
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="success" plain @click="goToMenu" class="action-btn">
              <el-icon><Coffee /></el-icon>
              立即购物
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="warning" plain @click="showBenefits = true" class="action-btn">
              <el-icon><Present /></el-icon>
              会员权益
            </el-button>
          </el-col>
          <el-col :xs="12" :sm="6">
            <el-button type="info" plain @click="refreshData" class="action-btn">
              <el-icon><Refresh /></el-icon>
              刷新数据
            </el-button>
          </el-col>
        </el-row>
      </el-card>

      <!-- 积分历史 -->
      <el-card class="points-history glass-panel">
        <div class="card-header">
          <h3>积分记录</h3>
        </div>
        <el-table
          :data="pointsStore.transactions"
          v-loading="pointsStore.loading"
          stripe
          class="points-table"
        >
          <el-table-column prop="transaction_type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag
                :type="pointsStore.getTransactionTypeColor(row.transaction_type)"
                size="small"
              >
                {{ pointsStore.getTransactionTypeText(row.transaction_type) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="说明" min-width="200" />
          <el-table-column prop="points_change" label="积分变动" width="120">
            <template #default="{ row }">
              <span :class="row.points_change > 0 ? 'points-add' : 'points-sub'">
                {{ pointsStore.formatPointsChange(row.points_change) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
        </el-table>

        <el-pagination
          v-if="pointsStore.pagination.total > 0"
          v-model:current-page="currentPage"
          :page-size="pointsStore.pagination.per_page"
          :total="pointsStore.pagination.total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
          class="pagination"
        />
      </el-card>
    </div>

    <!-- 未登录状态 -->
    <div v-else class="not-logged-in">
      <el-icon :size="80" class="empty-icon"><User /></el-icon>
      <h3>请先登录</h3>
      <p>登录后即可查看会员信息和积分</p>
      <el-button type="primary" @click="goToLogin" size="large">
        立即登录
      </el-button>
    </div>

    <!-- 会员权益弹窗 -->
    <el-dialog v-model="showBenefits" title="会员等级权益" width="500px">
      <div class="benefits-content">
        <div v-for="level in memberLevels" :key="level.name" class="benefit-item">
          <div class="benefit-header" :class="level.name">
            <span class="level-name">{{ level.displayName }}</span>
            <span class="level-points">{{ level.minPoints }}+ 积分</span>
          </div>
          <ul class="benefit-list">
            <li>积分获取倍率: {{ level.rate }}x</li>
            <li>最大抵扣比例: {{ level.discount }}%</li>
            <li v-for="b in level.benefits" :key="b">{{ b }}</li>
          </ul>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  Trophy, Medal, Wallet, Coin, TrendCharts, User, Coffee, 
  Present, Refresh 
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { usePointsStore } from '@/store/points'

const router = useRouter()
const authStore = useAuthStore()
const pointsStore = usePointsStore()

const currentPage = ref(1)
const showBenefits = ref(false)

const memberLevels = [
  { name: 'bronze', displayName: '铜牌会员', minPoints: 0, rate: 1.0, discount: 20, benefits: ['基础会员权益'] },
  { name: 'silver', displayName: '银牌会员', minPoints: 1000, rate: 1.2, discount: 30, benefits: ['生日双倍积分'] },
  { name: 'gold', displayName: '金牌会员', minPoints: 5000, rate: 1.5, discount: 40, benefits: ['专属优惠券'] },
  { name: 'platinum', displayName: '白金会员', minPoints: 20000, rate: 2.0, discount: 50, benefits: ['专属客服', '免费配送'] }
]

const levelProgress = computed(() => {
  if (!pointsStore.nextLevel) return 100
  const current = pointsStore.lifetimePoints
  const needed = pointsStore.nextLevel.required_points
  const prevLevel = memberLevels.find(l => l.name === pointsStore.memberLevel)
  const prevPoints = prevLevel?.minPoints || 0
  return Math.min(100, ((current - prevPoints) / (needed - prevPoints)) * 100)
})

const levelProgressColor = computed(() => {
  const colors = { bronze: '#CD7F32', silver: '#C0C0C0', gold: '#FFD700', platinum: '#E5E4E2' }
  return colors[pointsStore.memberLevel] || '#8B7355'
})

const getLevelName = (level) => {
  const names = { bronze: '铜牌', silver: '银牌', gold: '金牌', platinum: '白金' }
  return names[level] || level
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const goToProfile = () => router.push('/user/profile')
const goToMenu = () => router.push('/menu')
const goToLogin = () => router.push('/auth/login')

const handlePageChange = (page) => {
  currentPage.value = page
  pointsStore.fetchTransactions({ page })
}

const refreshData = async () => {
  try {
    await pointsStore.refreshPointsData()
    ElMessage.success('数据已刷新')
  } catch {
    ElMessage.error('刷新失败')
  }
}

onMounted(async () => {
  if (authStore.isAuthenticated) {
    await pointsStore.refreshPointsData()
  }
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.member-center-page {
  min-height: 100vh;
  padding: $spacing-2xl $spacing-lg;
  max-width: 1200px;
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

.member-level-card {
  margin-bottom: $spacing-lg;
  padding: $spacing-xl;
  border-radius: $radius-lg;

  .level-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    gap: $spacing-lg;
  }

  .level-info {
    display: flex;
    align-items: center;
    gap: $spacing-lg;
  }

  .level-badge {
    display: flex;
    align-items: center;
    gap: $spacing-sm;
    padding: $spacing-sm $spacing-lg;
    border-radius: $radius-full;
    font-weight: $font-weight-bold;
    color: white;

    &.bronze { background: linear-gradient(135deg, #CD7F32, #B8860B); }
    &.silver { background: linear-gradient(135deg, #C0C0C0, #808080); }
    &.gold { background: linear-gradient(135deg, #FFD700, #FFA500); }
    &.platinum { background: linear-gradient(135deg, #E5E4E2, #9370DB); }
  }

  .user-info {
    h3 { margin: 0 0 $spacing-xs 0; font-size: $font-size-xl; }
    p { margin: 0; color: $text-secondary; font-size: $font-size-sm; }
  }

  .level-progress {
    flex: 1;
    max-width: 400px;

    .progress-text {
      margin: 0 0 $spacing-sm 0;
      font-size: $font-size-sm;
      color: $text-secondary;
      strong { color: $primary-color; }
    }
  }
}

.points-cards {
  margin-bottom: $spacing-lg;

  .el-col { margin-bottom: $spacing-md; }
}

.points-card {
  padding: $spacing-lg;
  border-radius: $radius-lg;

  .points-content {
    display: flex;
    align-items: center;
    gap: $spacing-md;
  }

  .points-icon {
    font-size: 40px;
    color: $primary-color;
  }

  .points-value {
    font-size: $font-size-3xl;
    font-weight: $font-weight-bold;
    color: $text-primary;
  }

  .points-label {
    font-size: $font-size-sm;
    color: $text-secondary;
  }

  &.available .points-icon { color: $success-color; }
  &.total .points-icon { color: $warning-color; }
  &.lifetime .points-icon { color: $info-color; }
}

.quick-actions {
  margin-bottom: $spacing-lg;
  padding: $spacing-lg;
  border-radius: $radius-lg;

  h3 { margin: 0 0 $spacing-lg 0; }

  .action-btn {
    width: 100%;
    height: 60px;
    flex-direction: column;
    gap: $spacing-xs;
    border-radius: $radius-md;
  }
}

.points-history {
  padding: $spacing-lg;
  border-radius: $radius-lg;

  .card-header {
    margin-bottom: $spacing-lg;
    h3 { margin: 0; }
  }

  .points-add { color: $success-color; font-weight: $font-weight-bold; }
  .points-sub { color: $danger-color; font-weight: $font-weight-bold; }

  .pagination {
    margin-top: $spacing-lg;
    justify-content: center;
  }
}

.not-logged-in {
  @include flex-center;
  flex-direction: column;
  gap: $spacing-lg;
  min-height: 50vh;
  text-align: center;

  .empty-icon { color: $text-light; opacity: 0.5; }
  h3 { margin: 0; color: $text-primary; }
  p { margin: 0; color: $text-secondary; }
}

.benefits-content {
  .benefit-item {
    margin-bottom: $spacing-lg;
    border-radius: $radius-md;
    overflow: hidden;
    border: 1px solid $border-color;

    &:last-child { margin-bottom: 0; }
  }

  .benefit-header {
    padding: $spacing-sm $spacing-md;
    color: white;
    @include flex-between;

    &.bronze { background: linear-gradient(135deg, #CD7F32, #B8860B); }
    &.silver { background: linear-gradient(135deg, #C0C0C0, #808080); }
    &.gold { background: linear-gradient(135deg, #FFD700, #FFA500); }
    &.platinum { background: linear-gradient(135deg, #E5E4E2, #9370DB); }
  }

  .benefit-list {
    padding: $spacing-md;
    margin: 0;
    list-style: none;
    li {
      padding: $spacing-xs 0;
      font-size: $font-size-sm;
      color: $text-secondary;
      &::before { content: '✓ '; color: $success-color; }
    }
  }
}

// 移动端适配
@include respond-to(md) {
  .member-center-page {
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
  
  .member-level-card {
    padding: $spacing-md;
    
    .level-header {
      flex-direction: column;
      align-items: flex-start;
    }
    
    .level-info {
      flex-direction: column;
      align-items: flex-start;
      gap: $spacing-sm;
      width: 100%;
    }
    
    .level-progress {
      width: 100%;
      max-width: none;
      margin-top: $spacing-md;
    }
  }
  
  .points-card {
    padding: $spacing-md;
    
    .points-icon {
      font-size: 32px;
    }
    
    .points-value {
      font-size: $font-size-2xl;
    }
  }
  
  .quick-actions {
    padding: $spacing-md;
    
    h3 {
      font-size: $font-size-lg;
      margin-bottom: $spacing-md;
    }
    
    .action-btn {
      height: 50px;
      font-size: $font-size-sm;
    }
  }
  
  .points-history {
    padding: $spacing-md;
    
    .points-table {
      font-size: $font-size-sm;
      
      :deep(.el-table__header th) {
        font-size: $font-size-xs;
        padding: $spacing-xs;
      }
      
      :deep(.el-table__body td) {
        padding: $spacing-xs;
      }
    }
  }
}

@media (max-width: 480px) {
  .member-center-page {
    padding: $spacing-md $spacing-xs;
  }
  
  .page-header .page-title {
    font-size: $font-size-xl;
  }
  
  .member-level-card {
    .level-badge {
      padding: $spacing-xs $spacing-sm;
      font-size: $font-size-sm;
    }
    
    .user-info {
      h3 { font-size: $font-size-lg; }
      p { font-size: $font-size-xs; }
    }
  }
  
  .points-card {
    .points-content {
      flex-direction: column;
      text-align: center;
    }
    
    .points-value {
      font-size: $font-size-xl;
    }
  }
  
  // 积分历史表格在小屏幕上横向滚动
  .points-history {
    .points-table {
      :deep(.el-table__body-wrapper) {
        overflow-x: auto;
      }
    }
  }
}
</style>
