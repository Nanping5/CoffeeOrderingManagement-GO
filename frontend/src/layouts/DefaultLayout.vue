<template>
  <div class="default-layout">
    <!-- 增强版导航栏 -->
    <AppHeaderEnhanced />

    <!-- 主要内容区域 -->
    <main class="main-content">
      <router-view v-slot="{ Component, route }">
        <transition name="page" mode="out-in">
          <div :key="route.path" class="page-wrapper">
            <component :is="Component" />
          </div>
        </transition>
      </router-view>
    </main>

    <!-- 页脚（如果需要） -->
    <AppFooter v-if="showFooter" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import AppHeaderEnhanced from '@/components/AppHeaderEnhanced.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useCartStore } from '@/store/cart'

const cartStore = useCartStore()
const showFooter = ref(false) // 可根据需要显示页脚

onMounted(() => {
  // 初始化购物车
  cartStore.initCart()
})
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.default-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-content {
  flex: 1;
  padding-top: 100px;
  background-color: var(--background-color);
  position: relative;
  
  @include respond-to(md) {
    padding-top: 90px;
  }
}

.page-wrapper {
  width: 100%;
  min-height: calc(100vh - 100px);
  
  @include respond-to(md) {
    min-height: calc(100vh - 90px);
  }
}

// 页面加载状态
:global(.page-loading) {
  overflow: hidden;
}

// 路由过渡动画
.transition-fade-enter-active,
.transition-fade-leave-active {
  transition: opacity 0.3s ease;
}

.transition-fade-enter-from,
.transition-fade-leave-to {
  opacity: 0;
}

// 从左滑入
.transition-slide-left-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-left-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-left-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.transition-slide-left-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

// 从右滑入
.transition-slide-right-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-right-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-right-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.transition-slide-right-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

// 从下滑入
.transition-slide-up-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-up-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-up-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.transition-slide-up-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}

// 从上滑入
.transition-slide-down-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-down-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-slide-down-enter-from {
  opacity: 0;
  transform: translateY(-30px);
}

.transition-slide-down-leave-to {
  opacity: 0;
  transform: translateY(30px);
}

// 缩放过渡
.transition-scale-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-scale-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-scale-enter-from {
  opacity: 0;
  transform: scale(0.9);
}

.transition-scale-leave-to {
  opacity: 0;
  transform: scale(1.1);
}

// 翻转过渡
.transition-flip-enter-active {
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-flip-leave-active {
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.transition-flip-enter-from {
  opacity: 0;
  transform: perspective(800px) rotateY(-90deg);
}

.transition-flip-leave-to {
  opacity: 0;
  transform: perspective(800px) rotateY(90deg);
}
</style>