<template>
  <div class="admin-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <!-- 侧边栏 -->
    <AdminSidebar :collapsed="sidebarCollapsed" @toggle="toggleSidebar" />

    <!-- 主要内容区域 -->
    <div class="admin-main">
      <!-- 顶部导航栏 -->
      <AdminHeader @toggle-sidebar="toggleSidebar" :sidebar-collapsed="sidebarCollapsed" />

      <!-- 页面内容 -->
      <main class="admin-content">
        <router-view v-slot="{ Component, route }">
          <transition name="fade-slide" mode="out-in" appear>
            <div :key="route.path" class="page-wrapper">
              <component :is="Component" />
            </div>
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import AdminHeader from '@/components/admin/AdminHeader.vue'

const sidebarCollapsed = ref(false)

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}
</script>

<style lang="scss" scoped>
.admin-layout {
  display: flex;
  height: 100vh;
  background: var(--admin-bg);
  overflow: hidden;
}

.admin-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: margin-left 0.3s ease;
}

.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: var(--admin-bg-gradient);

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    background: var(--admin-scrollbar-track);
  }

  &::-webkit-scrollbar-thumb {
    background: var(--admin-scrollbar-thumb);
    border-radius: 4px;

    &:hover {
      background: var(--admin-scrollbar-thumb);
      opacity: 0.8;
    }
  }
}

.page-wrapper {
  width: 100%;
  min-height: 100%;
}

// 过渡动画
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

// 响应式设计 - 移动端适配
@media (max-width: 768px) {
  .admin-layout {
    position: relative;
  }

  .admin-main {
    width: 100%;
    margin-left: 0 !important;
  }

  .admin-content {
    padding: 12px;
    padding-bottom: 80px; // 为底部导航留出空间
  }
}

@media (max-width: 480px) {
  .admin-content {
    padding: 8px;
    padding-bottom: 80px;
  }
}
</style>
