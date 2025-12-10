<template>
  <div class="admin-login-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <div class="login-container animate__animated animate__fadeIn">
      <!-- 左侧品牌区 -->
      <div class="brand-section">
        <div class="brand-content">
          <div class="logo-wrapper">
            <el-icon :size="60"><CoffeeCup /></el-icon>
          </div>
          <h1>咖啡点餐系统</h1>
          <p>Coffee Ordering Management</p>
          <div class="features">
            <div class="feature-item">
              <el-icon><DataAnalysis /></el-icon>
              <span>数据统计分析</span>
            </div>
            <div class="feature-item">
              <el-icon><Menu /></el-icon>
              <span>菜单管理</span>
            </div>
            <div class="feature-item">
              <el-icon><List /></el-icon>
              <span>订单处理</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录表单 -->
      <div class="form-section">
        <div class="form-header">
          <h2>管理员登录</h2>
          <p>请输入您的管理员账号</p>
        </div>

        <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" class="login-form">
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="管理员账号"
              size="large"
              :prefix-icon="User"
              clearable
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="密码"
              size="large"
              :prefix-icon="Lock"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              :loading="loading"
              @click="handleLogin"
              class="login-button"
            >
              <el-icon v-if="!loading"><Right /></el-icon>
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>

        <div class="form-footer">
          <el-link type="info" @click="goToHome">
            <el-icon><ArrowLeft /></el-icon>
            返回首页
          </el-link>
        </div>

        <!-- 演示账户 -->
        <div class="demo-account">
          <el-divider>演示账户</el-divider>
          <div class="account-info" @click="fillDemo">
            <span>账号: <code>admin</code></span>
            <span>密码: <code>admin123</code></span>
            <el-tag size="small" type="info">点击填充</el-tag>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock, ArrowLeft, CoffeeCup, DataAnalysis, Menu, List, Right } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({ username: '', password: '' })

const loginRules = {
  username: [{ required: true, message: '请输入管理员账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  const valid = await loginFormRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const result = await authStore.adminLogin(loginForm)
    if (result.success) {
      ElMessage.success('登录成功！')
      await new Promise(resolve => setTimeout(resolve, 500))
      router.push(route.query.redirect || '/admin')
    } else {
      ElMessage.error(result.message || '账号或密码错误')
    }
  } catch (error) {
    ElMessage.error('登录失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const goToHome = () => router.push('/menu')

const fillDemo = () => {
  loginForm.username = 'admin'
  loginForm.password = 'admin123'
}
</script>

<style lang="scss" scoped>
.admin-login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 20px;
  position: relative;
  overflow: hidden;
}

.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;

  .circle {
    position: absolute;
    border-radius: 50%;
    background: linear-gradient(135deg, rgba(139, 69, 19, 0.3), rgba(160, 82, 45, 0.1));
    animation: float 20s infinite ease-in-out;
  }

  .circle-1 { width: 400px; height: 400px; top: -100px; right: -100px; animation-delay: 0s; }
  .circle-2 { width: 300px; height: 300px; bottom: -50px; left: -50px; animation-delay: -5s; }
  .circle-3 { width: 200px; height: 200px; top: 50%; left: 50%; animation-delay: -10s; }
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(20px, -20px) scale(1.05); }
  50% { transform: translate(-10px, 10px) scale(0.95); }
  75% { transform: translate(-20px, -10px) scale(1.02); }
}

.login-container {
  display: flex;
  width: 100%;
  max-width: 900px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
}

.brand-section {
  flex: 1;
  background: linear-gradient(135deg, #8b4513 0%, #a0522d 50%, #cd853f 100%);
  padding: 50px 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;

  .brand-content { text-align: center; }

  .logo-wrapper {
    width: 100px;
    height: 100px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 24px;
    backdrop-filter: blur(10px);
  }

  h1 { font-size: 28px; font-weight: 700; margin-bottom: 8px; }
  p { opacity: 0.8; font-size: 14px; margin-bottom: 40px; }

  .features {
    text-align: left;
    .feature-item {
      display: flex;
      align-items: center;
      gap: 12px;
      padding: 12px 20px;
      background: rgba(255, 255, 255, 0.1);
      border-radius: 12px;
      margin-bottom: 12px;
      backdrop-filter: blur(5px);
      font-size: 14px;
      .el-icon { font-size: 20px; }
    }
  }
}

.form-section {
  flex: 1;
  padding: 50px 40px;
  display: flex;
  flex-direction: column;
  justify-content: center;

  .form-header {
    margin-bottom: 32px;
    h2 { font-size: 28px; font-weight: 700; color: #2c3e50; margin-bottom: 8px; }
    p { color: #7f8c8d; font-size: 14px; }
  }

  .login-form {
    .el-form-item { margin-bottom: 24px; }
    :deep(.el-input__wrapper) {
      border-radius: 12px;
      padding: 4px 16px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      &:hover, &.is-focus { box-shadow: 0 4px 16px rgba(139, 69, 19, 0.2); }
    }
  }

  .login-button {
    width: 100%;
    height: 50px;
    font-size: 16px;
    font-weight: 600;
    border-radius: 12px;
    background: linear-gradient(135deg, #8b4513 0%, #a0522d 100%);
    border: none;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    &:hover {
      background: linear-gradient(135deg, #a0522d 0%, #8b4513 100%);
      transform: translateY(-2px);
      box-shadow: 0 8px 24px rgba(139, 69, 19, 0.4);
    }
  }

  .form-footer {
    text-align: center;
    margin-top: 24px;
    .el-link { font-size: 14px; display: inline-flex; align-items: center; gap: 4px; }
  }

  .demo-account {
    margin-top: 32px;
    .el-divider { margin: 16px 0; :deep(.el-divider__text) { color: #909399; font-size: 12px; } }
    .account-info {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 16px;
      padding: 12px;
      background: #f8f9fa;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s ease;
      font-size: 13px;
      color: #666;
      code { background: #e9ecef; padding: 2px 8px; border-radius: 4px; color: #8b4513; font-weight: 600; }
      &:hover { background: #e9ecef; }
    }
  }
}

@media (max-width: 768px) {
  .login-container { flex-direction: column; max-width: 400px; }
  .brand-section { padding: 30px 20px; .features { display: none; } h1 { font-size: 22px; } p { margin-bottom: 0; } }
  .form-section { padding: 30px 20px; }
}
</style>
