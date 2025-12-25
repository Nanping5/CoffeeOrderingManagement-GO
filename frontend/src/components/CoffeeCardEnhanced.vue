<template>
  <div 
    class="coffee-card-enhanced" 
    :class="{ 'out-of-stock': !item.is_available }"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    @mousemove="handleMouseMove"
    :style="cardStyle"
  >
    <!-- 玻璃态卡片容器 -->
    <div class="card-glass-container">
      <!-- 商品图片 -->
      <div class="card-image-wrapper">
        <div class="image-container">
          <img
            v-if="item.image_url"
            :src="item.image_url"
            :alt="item.name"
            @error="handleImageError"
            class="product-image"
            loading="lazy"
          />
          <div v-else class="image-placeholder">
            <el-icon :size="48"><Coffee /></el-icon>
          </div>
          
          <!-- 图片遮罩层 -->
          <div class="image-overlay"></div>
        </div>

        <!-- 折扣标签 -->
        <transition name="badge-bounce">
          <div v-if="hasDiscount" class="discount-badge glass-badge">
            <span class="badge-text">-{{ discountPercentage }}%</span>
          </div>
        </transition>

        <!-- 热门标签 -->
        <transition name="badge-bounce">
          <div v-if="item.is_popular" class="popular-badge glass-badge">
            <el-icon class="badge-icon"><Star /></el-icon>
            <span>热门</span>
          </div>
        </transition>

        <!-- 售罄遮罩 -->
        <transition name="fade">
          <div v-if="!item.is_available" class="out-of-stock-overlay">
            <div class="sold-out-badge">
              <el-icon :size="32"><CircleClose /></el-icon>
              <span>暂时售罄</span>
            </div>
          </div>
        </transition>

        <!-- 快速查看按钮 -->
        <transition name="scale-fade">
          <div v-if="isHovered && item.is_available" class="quick-view-btn">
            <el-button
              type="primary"
              :icon="View"
              circle
              size="large"
              @click.stop="handleViewDetail"
              class="glass-button"
            />
          </div>
        </transition>
      </div>

      <!-- 商品信息 -->
      <div class="card-content">
        <!-- 商品名称和分类 -->
        <div class="item-header">
          <h3 class="item-name gradient-text" :title="item.name">
            {{ item.name }}
          </h3>
          <el-tag size="small" :type="getCategoryType(item.category)" effect="light" class="category-tag">
            {{ categoryLabel }}
          </el-tag>
        </div>

        <!-- 商品描述 -->
        <p class="item-description" :title="item.description">
          {{ item.description || '暂无描述' }}
        </p>

        <!-- 价格信息 -->
        <div class="price-section">
          <div class="price-wrapper">
            <span class="current-price">
              <span class="currency">¥</span>
              <span class="price-value">{{ item.price }}</span>
            </span>
            <transition name="slide-fade">
              <span 
                v-if="item.original_price && item.original_price > item.price" 
                class="original-price"
              >
                ¥{{ item.original_price }}
              </span>
            </transition>
          </div>

          <!-- 销量信息 -->
          <div v-if="item.order_count > 0" class="sales-info">
            <el-icon class="sales-icon"><TrendCharts /></el-icon>
            <span>已售 {{ item.order_count }}</span>
          </div>
        </div>

        <!-- 标签区域 -->
        <transition-group name="tag-stagger" tag="div" class="tags-section" v-if="item.tags && item.tags.length > 0">
          <el-tag
            v-for="(tag, index) in item.tags"
            :key="tag"
            size="small"
            effect="plain"
            class="item-tag"
            :style="{ transitionDelay: `${index * 50}ms` }"
          >
            {{ tag }}
          </el-tag>
        </transition-group>

        <!-- 操作区域 -->
        <div class="action-section">
          <!-- 数量选择器 -->
          <transition name="slide-up">
            <div class="quantity-selector" v-if="item.is_available && showQuantity">
              <el-input-number
                v-model="quantity"
                :min="1"
                :max="99"
                size="small"
                @change="handleQuantityChange"
                class="quantity-input"
              />
            </div>
          </transition>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <!-- 加入购物车按钮 -->
            <el-button
              v-if="item.is_available"
              type="primary"
              size="default"
              :icon="ShoppingCart"
              :loading="adding"
              @click="handleAddToCart"
              class="add-cart-btn ripple-button"
            >
              <span class="button-text">加入购物车</span>
            </el-button>

            <!-- 到货通知按钮 -->
            <el-button
              v-else
              type="warning"
              size="default"
              :icon="Bell"
              @click="handleNotifyWhenAvailable"
              class="notify-btn ripple-button"
            >
              <span class="button-text">到货通知</span>
            </el-button>
          </div>
        </div>
      </div>

      <!-- 悬停光效 -->
      <div class="card-shine" :style="shineStyle"></div>
    </div>

    <!-- 添加成功动画 -->
    <transition name="success-bounce">
      <div v-if="showSuccessAnimation" class="success-animation">
        <el-icon :size="48" class="success-icon"><SuccessFilled /></el-icon>
        <span class="success-text">已添加</span>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Coffee, Star, ShoppingCart, Bell, View, 
  CircleClose, TrendCharts, SuccessFilled 
} from '@element-plus/icons-vue'

// Props
const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  categoryLabel: {
    type: String,
    default: '咖啡'
  }
})

// Emits
const emit = defineEmits(['add-to-cart', 'view-detail', 'notify-when-available'])

// 响应式数据
const quantity = ref(1)
const adding = ref(false)
const isHovered = ref(false)
const showQuantity = ref(false)
const showSuccessAnimation = ref(false)
const mouseX = ref(0)
const mouseY = ref(0)
const cardRotateX = ref(0)
const cardRotateY = ref(0)

// 计算属性
const hasDiscount = computed(() => {
  return props.item.original_price && props.item.original_price > props.item.price
})

const discountPercentage = computed(() => {
  if (!hasDiscount.value) return 0
  return Math.round(
    ((props.item.original_price - props.item.price) / props.item.original_price) * 100
  )
})

const cardStyle = computed(() => {
  if (!isHovered.value) return {}
  return {
    transform: `perspective(1000px) rotateX(${cardRotateX.value}deg) rotateY(${cardRotateY.value}deg) translateZ(10px)`
  }
})

const shineStyle = computed(() => {
  return {
    background: `radial-gradient(circle at ${mouseX.value}% ${mouseY.value}%, rgba(255, 255, 255, 0.3) 0%, transparent 50%)`,
    opacity: isHovered.value ? 1 : 0
  }
})

// 方法
const handleMouseEnter = () => {
  isHovered.value = true
  showQuantity.value = true
}

const handleMouseLeave = () => {
  isHovered.value = false
  cardRotateX.value = 0
  cardRotateY.value = 0
  setTimeout(() => {
    if (!isHovered.value) {
      showQuantity.value = false
    }
  }, 300)
}

const handleMouseMove = (event) => {
  if (!isHovered.value) return
  
  const card = event.currentTarget
  const rect = card.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  
  mouseX.value = (x / rect.width) * 100
  mouseY.value = (y / rect.height) * 100
  
  // 3D倾斜效果
  const centerX = rect.width / 2
  const centerY = rect.height / 2
  const rotateY = ((x - centerX) / centerX) * 5
  const rotateX = ((centerY - y) / centerY) * 5
  
  cardRotateX.value = rotateX
  cardRotateY.value = rotateY
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
}

const handleQuantityChange = (value) => {
  if (value < 1) quantity.value = 1
  if (value > 99) quantity.value = 99
}

const handleAddToCart = async () => {
  if (adding.value) return

  adding.value = true

  try {
    await emit('add-to-cart', props.item, quantity.value)
    
    // 显示成功动画
    showSuccessAnimation.value = true
    setTimeout(() => {
      showSuccessAnimation.value = false
    }, 1500)
    
    // 重置数量
    quantity.value = 1
  } catch (error) {
    console.error('添加到购物车失败:', error)
  } finally {
    adding.value = false
  }
}

const handleViewDetail = () => {
  emit('view-detail', props.item)
}

const handleNotifyWhenAvailable = () => {
  emit('notify-when-available', props.item)
  ElMessage.success('我们会在商品到货时通知您')
}

const getCategoryType = (category) => {
  const typeMap = {
    coffee: 'warning',
    tea: 'success',
    dessert: 'danger',
    other: 'info'
  }
  return typeMap[category] || 'info'
}
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.coffee-card-enhanced {
  height: 100%;
  position: relative;
  transition: transform $transition-base $ease-out;
  transform-style: preserve-3d;
  
  &:hover {
    z-index: 10;
  }
  
  &.out-of-stock {
    .card-glass-container {
      opacity: 0.7;
    }
  }
}

.card-glass-container {
  height: 100%;
  @include glass-effect(0.15, 12px);
  border-radius: $radius-lg;
  overflow: hidden;
  box-shadow: $shadow-lg;
  transition: all $transition-base $ease-out;
  position: relative;
  
  .coffee-card-enhanced:hover & {
    box-shadow: $shadow-float;
    border-color: rgba($primary-color, 0.3);
  }
}

.card-image-wrapper {
  position: relative;
  width: 100%;
  height: 220px;
  overflow: hidden;
  
  .image-container {
    width: 100%;
    height: 100%;
    position: relative;
    
    .product-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
      transition: transform $transition-slow $ease-out;
      
      .coffee-card-enhanced:hover & {
        transform: scale(1.1);
      }
    }
    
    .image-placeholder {
      width: 100%;
      height: 100%;
      @include flex-center;
      background: $warm-gradient;
      color: rgba(255, 255, 255, 0.8);
      
      .el-icon {
        transition: transform $transition-base $ease-bounce;
      }
      
      .coffee-card-enhanced:hover & .el-icon {
        transform: scale(1.2) rotate(10deg);
      }
    }
    
    .image-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: linear-gradient(180deg, transparent 0%, rgba(0, 0, 0, 0.3) 100%);
      opacity: 0;
      transition: opacity $transition-base $ease-out;
      
      .coffee-card-enhanced:hover & {
        opacity: 1;
      }
    }
  }
}

.glass-badge {
  @include glass-effect(0.2, 8px);
  padding: $spacing-xs $spacing-sm;
  border-radius: $radius-md;
  font-size: $font-size-sm;
  font-weight: $font-weight-bold;
  box-shadow: $shadow-sm;
  backdrop-filter: blur(8px);
}

.discount-badge {
  position: absolute;
  top: $spacing-md;
  left: $spacing-md;
  background: linear-gradient(135deg, #e74c3c, #c0392b);
  color: white;
  z-index: 2;
  
  .badge-text {
    display: block;
  }
}

.popular-badge {
  position: absolute;
  top: $spacing-md;
  right: $spacing-md;
  background: linear-gradient(135deg, #f39c12, #e67e22);
  color: white;
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  z-index: 2;
  
  .badge-icon {
    font-size: $font-size-sm;
  }
}

.out-of-stock-overlay {
  @include absolute-center;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.7);
  @include flex-center;
  z-index: 3;
  
  .sold-out-badge {
    @include flex-center;
    flex-direction: column;
    gap: $spacing-sm;
    color: white;
    font-size: $font-size-lg;
    font-weight: $font-weight-bold;
    background: rgba($danger-color, 0.9);
    padding: $spacing-md $spacing-lg;
    border-radius: $radius-lg;
    @include glass-effect(0.1, 10px);
  }
}

.quick-view-btn {
  @include absolute-center;
  z-index: 4;
  
  .glass-button {
    width: 56px;
    height: 56px;
    @include glass-effect(0.2, 10px);
    border: 2px solid rgba(255, 255, 255, 0.3);
    box-shadow: $shadow-lg;
    
    &:hover {
      transform: scale(1.1);
      box-shadow: $shadow-xl;
    }
  }
}

.card-content {
  padding: $spacing-lg;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: $spacing-sm;
  
  .item-name {
    flex: 1;
    font-size: $font-size-xl;
    font-weight: $font-weight-bold;
    line-height: $line-height-tight;
    @include text-ellipsis(2);
    @include gradient-text($primary-gradient);
  }
  
  .category-tag {
    flex-shrink: 0;
  }
}

.item-description {
  color: $text-secondary;
  font-size: $font-size-sm;
  line-height: $line-height-relaxed;
  @include text-ellipsis(2);
}

.price-section {
  .price-wrapper {
    display: flex;
    align-items: baseline;
    gap: $spacing-sm;
    margin-bottom: $spacing-xs;
    
    .current-price {
      display: flex;
      align-items: baseline;
      color: $danger-color;
      font-weight: $font-weight-bold;
      
      .currency {
        font-size: $font-size-lg;
      }
      
      .price-value {
        font-size: $font-size-3xl;
      }
    }
    
    .original-price {
      font-size: $font-size-base;
      color: $text-light;
      text-decoration: line-through;
    }
  }
  
  .sales-info {
    display: flex;
    align-items: center;
    gap: $spacing-xs;
    font-size: $font-size-xs;
    color: $text-light;
    
    .sales-icon {
      font-size: $font-size-sm;
    }
  }
}

.tags-section {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-xs;
  
  .item-tag {
    font-size: $font-size-xs;
    transition: all $transition-base $ease-out;
    
    &:hover {
      transform: translateY(-2px);
    }
  }
}

.action-section {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
  
  .quantity-selector {
    .quantity-input {
      width: 120px;
    }
  }
  
  .action-buttons {
    display: flex;
    gap: $spacing-sm;
    
    .add-cart-btn,
    .notify-btn {
      flex: 1;
      border-radius: $radius-full;
      font-weight: $font-weight-semibold;
      transition: all $transition-base $ease-out;
      @include ripple-effect;
      
      &:hover {
        transform: translateY(-2px);
        box-shadow: $shadow-md;
      }
      
      &:active {
        transform: translateY(0);
      }
    }
    
    .add-cart-btn {
      background: $primary-gradient;
      border: none;
    }
    
    .notify-btn {
      background: linear-gradient(135deg, $warning-color, darken($warning-color, 10%));
      border: none;
    }
  }
}

.card-shine {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  transition: opacity $transition-base $ease-out;
  border-radius: $radius-lg;
}

.success-animation {
  @include absolute-center;
  @include flex-center;
  flex-direction: column;
  gap: $spacing-sm;
  background: rgba(255, 255, 255, 0.95);
  padding: $spacing-xl;
  border-radius: $radius-lg;
  box-shadow: $shadow-xl;
  z-index: 100;
  
  .success-icon {
    color: $success-color;
    animation: bounce 0.6s $ease-bounce;
  }
  
  .success-text {
    font-size: $font-size-lg;
    font-weight: $font-weight-bold;
    color: $success-color;
  }
}

// 动画
.badge-bounce-enter-active {
  animation: bounceIn 0.6s $ease-bounce;
}

.scale-fade-enter-active,
.scale-fade-leave-active {
  transition: all $transition-base $ease-out;
}

.scale-fade-enter-from,
.scale-fade-leave-to {
  opacity: 0;
  transform: scale(0.8);
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all $transition-base $ease-out;
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}

.slide-up-enter-active {
  animation: slideInUp 0.4s $ease-out;
}

.success-bounce-enter-active {
  animation: bounceIn 0.5s $ease-bounce;
}

.success-bounce-leave-active {
  animation: fadeOut 0.3s $ease-out;
}

.tag-stagger-enter-active {
  transition: all 0.4s $ease-out;
}

.tag-stagger-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

// 响应式 - 移动端适配
@include respond-to(md) {
  .card-image-wrapper {
    height: 160px;
  }
  
  .card-content {
    padding: $spacing-sm $spacing-md;
    gap: $spacing-sm;
  }
  
  .item-header {
    .item-name {
      font-size: $font-size-base !important;
      -webkit-line-clamp: 1;
    }
    
    .category-tag {
      font-size: 10px;
      padding: 2px 6px;
    }
  }
  
  .item-description {
    font-size: $font-size-xs;
    -webkit-line-clamp: 1;
  }
  
  .price-section {
    .price-wrapper {
      .current-price {
        .currency {
          font-size: $font-size-base;
        }
        
        .price-value {
          font-size: $font-size-xl !important;
        }
      }
      
      .original-price {
        font-size: $font-size-xs;
      }
    }
    
    .sales-info {
      display: none;
    }
  }
  
  .tags-section {
    display: none;
  }
  
  .action-section {
    .quantity-selector {
      display: none;
    }
    
    .action-buttons {
      .add-cart-btn,
      .notify-btn {
        padding: $spacing-xs $spacing-sm;
        font-size: $font-size-xs;
        
        .button-text {
          display: none;
        }
        
        .el-icon {
          margin: 0;
        }
      }
    }
  }
  
  .quick-view-btn {
    display: none;
  }
  
  .discount-badge,
  .popular-badge {
    padding: 2px 6px;
    font-size: 10px;
    
    .badge-icon {
      font-size: 10px;
    }
  }
}

// 超小屏幕 - 单列布局时的优化
@media (max-width: 480px) {
  .card-glass-container {
    display: flex;
    flex-direction: row;
    height: auto;
  }
  
  .card-image-wrapper {
    width: 120px;
    height: 120px;
    flex-shrink: 0;
    border-radius: $radius-lg 0 0 $radius-lg;
    
    .image-container {
      border-radius: $radius-lg 0 0 $radius-lg;
    }
  }
  
  .card-content {
    flex: 1;
    padding: $spacing-sm;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }
  
  .item-header {
    flex-direction: column;
    align-items: flex-start;
    gap: $spacing-xs;
    
    .item-name {
      font-size: $font-size-base !important;
    }
  }
  
  .item-description {
    display: none;
  }
  
  .price-section {
    .price-wrapper {
      .current-price .price-value {
        font-size: $font-size-lg !important;
      }
    }
  }
  
  .action-section {
    margin-top: auto;
    
    .action-buttons {
      .add-cart-btn,
      .notify-btn {
        width: 100%;
        justify-content: center;
        
        .button-text {
          display: inline;
          font-size: $font-size-xs;
        }
      }
    }
  }
  
  .discount-badge {
    top: $spacing-xs;
    left: $spacing-xs;
  }
  
  .popular-badge {
    top: $spacing-xs;
    right: auto;
    left: $spacing-xs;
    top: auto;
    bottom: $spacing-xs;
  }
}
</style>
