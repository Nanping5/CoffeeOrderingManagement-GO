<template>
  <aside class="admin-sidebar" :class="{ collapsed }">
    <!-- Logo区域 -->
    <div class="sidebar-logo" @click="$emit('toggle')">
      <div class="logo-icon">
        <el-icon :size="collapsed ? 28 : 32">
          <CoffeeCup />
        </el-icon>
      </div>
      <transition name="fade">
        <h2 v-if="!collapsed">咖啡管理</h2>
      </transition>
    </div>

    <!-- 导航菜单 -->
    <nav class="sidebar-nav">
      <div
        v-for="item in menuItems"
        :key="item.path"
        class="nav-item"
        :class="{ active: isActive(item.path), 'has-children': item.children }"
      >
        <template v-if="!item.children">
          <router-link :to="item.path" class="nav-link">
            <el-icon class="nav-icon"><component :is="item.icon" /></el-icon>
            <transition name="fade">
              <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
            </transition>
          </router-link>
        </template>
        <template v-else>
          <div class="nav-group">
            <div class="nav-link group-header" @click="toggleGroup(item.path)">
              <el-icon class="nav-icon"><component :is="item.icon" /></el-icon>
              <transition name="fade">
                <span v-if="!collapsed" class="nav-label">{{ item.label }}</span>
              </transition>
              <el-icon v-if="!collapsed" class="arrow" :class="{ expanded: expandedGroups.includes(item.path) }">
                <ArrowRight />
              </el-icon>
            </div>
            <transition name="slide">
              <div v-if="!collapsed && expandedGroups.includes(item.path)" class="nav-children">
                <router-link
                  v-for="child in item.children"
                  :key="child.path"
                  :to="child.path"
                  class="nav-link child-link"
                  :class="{ active: isActive(child.path) }"
                >
                  <el-icon class="nav-icon"><component :is="child.icon" /></el-icon>
                  <span class="nav-label">{{ child.label }}</span>
                </router-link>
              </div>
            </transition>
          </div>
        </template>
      </div>
    </nav>

    <!-- 底部信息 -->
    <div class="sidebar-footer">
      <transition name="fade">
        <div v-if="!collapsed" class="footer-info">
          <span class="version">v1.0.0</span>
          <span class="divider">·</span>
          <span class="copyright">咖啡点餐系统</span>
        </div>
      </transition>
    </div>
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  DataAnalysis,
  Menu,
  CoffeeCup,
  List,
  Setting,
  User,
  Tools,
  ArrowRight
} from '@element-plus/icons-vue'

defineProps({
  collapsed: {
    type: Boolean,
    default: false
  }
})

defineEmits(['toggle'])

const route = useRoute()
const expandedGroups = ref(['settings'])

const menuItems = [
  { path: '/admin', icon: DataAnalysis, label: '仪表板' },
  { path: '/admin/menu', icon: Menu, label: '菜单管理' },
  { path: '/admin/orders', icon: List, label: '订单管理' },
  {
    path: 'settings',
    icon: Setting,
    label: '系统设置',
    children: [
      { path: '/admin/profile', icon: User, label: '个人信息' },
      { path: '/admin/settings', icon: Tools, label: '系统配置' }
    ]
  }
]

const isActive = (path) => {
  if (path === '/admin') {
    return route.path === '/admin'
  }
  return route.path.startsWith(path)
}

const toggleGroup = (path) => {
  const index = expandedGroups.value.indexOf(path)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(path)
  }
}
</script>

<style lang="scss" scoped>
.admin-sidebar {
  width: 260px;
  height: 100vh;
  background: var(--admin-sidebar-bg);
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  border-right: 1px solid var(--admin-sidebar-border);
  position: relative;
  z-index: 100;

  &.collapsed {
    width: 80px;
  }
}

.sidebar-logo {
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  border-bottom: 1px solid var(--admin-sidebar-border);
  transition: all 0.3s ease;

  &:hover {
    background: var(--admin-btn-bg);
  }

  .logo-icon {
    width: 48px;
    height: 48px;
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    flex-shrink: 0;
    box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
  }

  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 700;
    color: var(--admin-text-primary);
    white-space: nowrap;
  }
}

.sidebar-nav {
  flex: 1;
  padding: 16px 12px;
  overflow-y: auto;

  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--admin-scrollbar-thumb);
    border-radius: 2px;
  }
}

.nav-item {
  margin-bottom: 4px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 10px;
  color: var(--admin-text-secondary);
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;

  &:hover {
    background: var(--admin-btn-bg);
    color: var(--admin-text-primary);
  }

  &.router-link-exact-active,
  &.active {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
    color: white;
    box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);

    .nav-icon {
      color: white;
    }
  }

  .nav-icon {
    font-size: 20px;
    flex-shrink: 0;
  }

  .nav-label {
    font-size: 14px;
    font-weight: 500;
    white-space: nowrap;
  }

  .arrow {
    margin-left: auto;
    font-size: 12px;
    transition: transform 0.3s ease;

    &.expanded {
      transform: rotate(90deg);
    }
  }
}

.group-header {
  &:hover {
    background: var(--admin-btn-bg);
  }
}

.nav-children {
  margin-top: 4px;
  padding-left: 20px;

  .child-link {
    padding: 12px 16px;
    font-size: 13px;

    .nav-icon {
      font-size: 16px;
    }
  }
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--admin-sidebar-border);

  .footer-info {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 12px;
    color: var(--admin-text-muted);

    .divider {
      opacity: 0.5;
    }
  }
}

// 动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
  overflow: hidden;
}

.slide-enter-to,
.slide-leave-from {
  max-height: 200px;
}

// 响应式
@media (max-width: 768px) {
  .admin-sidebar {
    width: 80px;

    &.collapsed {
      width: 0;
      padding: 0;
      border: none;
    }
  }
}
</style>
