<template>
  <div class="theme-toggle">
    <el-tooltip
      :content="isDark ? '切换到浅色模式' : '切换到深色模式'"
      placement="bottom"
    >
      <button
        class="theme-toggle-btn glass-button"
        @click="handleToggle"
        :aria-label="isDark ? '切换到浅色模式' : '切换到深色模式'"
      >
        <transition name="theme-icon" mode="out-in">
          <el-icon :size="20" :key="isDark ? 'moon' : 'sun'">
            <Moon v-if="isDark" />
            <Sunny v-else />
          </el-icon>
        </transition>
      </button>
    </el-tooltip>
  </div>
</template>

<script setup>
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useTheme } from '@/composables/useTheme'

const { isDark, toggleTheme } = useTheme()

const handleToggle = () => {
  toggleTheme()
  
  // 添加切换动画效果
  const body = document.body
  body.classList.add('theme-transitioning')
  
  setTimeout(() => {
    body.classList.remove('theme-transitioning')
  }, 300)
}
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.theme-toggle {
  .theme-toggle-btn {
    @include glass-effect(0.15, 10px);
    width: 40px;
    height: 40px;
    border-radius: $radius-full;
    border: 1px solid rgba(255, 255, 255, 0.2);
    background: transparent;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all $transition-base $ease-out;
    color: $primary-color;
    
    &:hover {
      transform: scale(1.1) rotate(10deg);
      box-shadow: $shadow-md;
      background: rgba(255, 255, 255, 0.2);
    }
    
    &:active {
      transform: scale(0.95);
    }
    
    @include dark-mode {
      color: $accent-color;
      border-color: rgba(255, 255, 255, 0.1);
      
      &:hover {
        background: rgba(255, 255, 255, 0.1);
      }
    }
  }
}

// 主题图标切换动画
.theme-icon-enter-active,
.theme-icon-leave-active {
  transition: all $transition-base $ease-out;
}

.theme-icon-enter-from {
  opacity: 0;
  transform: rotate(-180deg) scale(0);
}

.theme-icon-leave-to {
  opacity: 0;
  transform: rotate(180deg) scale(0);
}

// 全局主题切换动画
:global(body.theme-transitioning) {
  * {
    transition: background-color $transition-base $ease-out,
                color $transition-base $ease-out,
                border-color $transition-base $ease-out !important;
  }
}
</style>
