<template>
  <header class="admin-header">
    <!-- 左侧：折叠按钮和面包屑 -->
    <div class="header-left">
      <button class="toggle-btn" @click="$emit('toggle-sidebar')">
        <el-icon :size="20">
          <Fold v-if="!sidebarCollapsed" />
          <Expand v-else />
        </el-icon>
      </button>
      
      <el-breadcrumb separator="/" class="breadcrumb">
        <el-breadcrumb-item :to="{ path: '/admin' }">
          <el-icon><HomeFilled /></el-icon>
          <span>首页</span>
        </el-breadcrumb-item>
        <el-breadcrumb-item v-if="currentPageName">
          {{ currentPageName }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <!-- 右侧：主题切换、通知、用户 -->
    <div class="header-right">
      <!-- 主题切换 -->
      <ThemeToggle />

      <!-- 通知 -->
      <div class="header-action">
        <el-badge :value="notificationCount" :hidden="notificationCount === 0" :max="99">
          <button class="action-btn" @click="showNotifications">
            <el-icon :size="20"><Bell /></el-icon>
          </button>
        </el-badge>
      </div>

      <!-- 用户下拉菜单 -->
      <el-dropdown @command="handleCommand" class="user-dropdown" trigger="click">
        <div class="user-info">
          <el-avatar :size="36" :src="userAvatar" class="user-avatar">
            <el-icon><User /></el-icon>
          </el-avatar>
          <div class="user-details">
            <span class="username">{{ username }}</span>
            <span class="role">管理员</span>
          </div>
          <el-icon class="arrow"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu class="user-menu">
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>
              系统设置
            </el-dropdown-item>
            <el-dropdown-item divided command="backToHome">
              <el-icon><House /></el-icon>
              返回前台
            </el-dropdown-item>
            <el-dropdown-item divided command="logout" class="logout-item">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  HomeFilled,
  Bell,
  User,
  ArrowDown,
  Setting,
  House,
  SwitchButton,
  Fold,
  Expand
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import ThemeToggle from '@/components/ThemeToggle.vue'

defineProps({
  sidebarCollapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle-sidebar'])

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const notificationCount = ref(3)
const username = computed(() => authStore.user?.username || '管理员')
const userAvatar = computed(() => authStore.user?.avatar || '')

const currentPageName = computed(() => {
  const routeNameMap = {
    AdminDashboard: '仪表板',
    AdminMenu: '菜单管理',
    AdminOrders: '订单管理',
    AdminProfile: '个人信息',
    AdminSettings: '系统设置'
  }
  return routeNameMap[route.name] || ''
})

const showNotifications = () => {
  ElMessage.info('暂无新通知')
  notificationCount.value = 0
}

const handleCommand = async (command) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'settings':
      router.push('/admin/settings')
      break
    case 'backToHome':
      router.push('/menu')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '退出确认', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        authStore.logout()
        ElMessage.success('已退出登录')
        router.push('/admin/login')
      } catch (error) {
        // 用户取消
      }
      break
  }
}
</script>

<style lang="scss" scoped>
.admin-header {
  height: 70px;
  background: var(--admin-header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--admin-sidebar-border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  position: sticky;
  top: 0;
  z-index: 99;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;

  .toggle-btn {
    width: 40px;
    height: 40px;
    border: none;
    border-radius: 10px;
    background: var(--admin-btn-bg);
    color: var(--admin-text-secondary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;

    &:hover {
      background: var(--admin-btn-hover);
      color: var(--admin-text-primary);
    }
  }

  .breadcrumb {
    :deep(.el-breadcrumb__item) {
      .el-breadcrumb__inner {
        display: flex;
        align-items: center;
        gap: 6px;
        color: var(--admin-text-secondary);
        font-weight: 500;

        &:hover {
          color: #f59e0b;
        }
      }

      &:last-child .el-breadcrumb__inner {
        color: var(--admin-text-primary);
      }
    }

    :deep(.el-breadcrumb__separator) {
      color: var(--admin-text-muted);
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-action {
  .action-btn {
    width: 40px;
    height: 40px;
    border: none;
    border-radius: 10px;
    background: var(--admin-btn-bg);
    color: var(--admin-text-secondary);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s ease;

    &:hover {
      background: var(--admin-btn-hover);
      color: var(--admin-text-primary);
    }
  }

  :deep(.el-badge__content) {
    background: #ef4444;
    border: none;
  }
}

.user-dropdown {
  cursor: pointer;

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px 16px;
    border-radius: 12px;
    background: var(--admin-btn-bg);
    transition: all 0.2s ease;

    &:hover {
      background: var(--admin-btn-hover);
    }

    .user-avatar {
      background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
      border: 2px solid rgba(245, 158, 11, 0.3);
    }

    .user-details {
      display: flex;
      flex-direction: column;

      .username {
        font-size: 14px;
        font-weight: 600;
        color: var(--admin-text-primary);
      }

      .role {
        font-size: 12px;
        color: var(--admin-text-secondary);
      }
    }

    .arrow {
      font-size: 12px;
      color: var(--admin-text-secondary);
      transition: transform 0.3s ease;
    }
  }

  &:hover .arrow {
    transform: rotate(180deg);
  }
}

@media (max-width: 768px) {
  .admin-header {
    padding: 0 16px;
  }

  .user-details {
    display: none !important;
  }
}
</style>
