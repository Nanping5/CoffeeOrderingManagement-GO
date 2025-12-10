<template>
  <div class="login-page">
    <el-card class="login-card glass-panel">
      <div class="login-header">
        <h2>用户登录</h2>
        <p>欢迎回到咖啡时光</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="formData"
        :rules="rules"
        @submit.prevent="handleLogin"
        class="login-form"
      >
        <el-form-item prop="email">
          <el-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            :prefix-icon="Message"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
            :prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <div class="login-options">
            <el-checkbox v-model="formData.rememberMe">记住登录状态</el-checkbox>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleLogin"
            class="login-button"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p>还没有账户？
          <router-link to="/auth/register" class="register-link">立即注册</router-link>
        </p>
        <el-divider>或</el-divider>
        <el-button text @click="continueAsGuest" class="guest-btn">
          游客模式继续购物
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Message, Lock } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loginFormRef = ref(null)
const loading = ref(false)

const formData = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    const result = await authStore.login({
      email: formData.email,
      password: formData.password
    })
    
    if (result.success) {
      ElMessage.success('登录成功')
      const redirect = route.query.redirect || '/menu'
      router.push(redirect)
    } else {
      ElMessage.error(result.errors?.[0] || '登录失败')
    }
  } catch (error) {
    if (error !== false) {
      ElMessage.error('登录失败，请重试')
    }
  } finally {
    loading.value = false
  }
}

const continueAsGuest = () => {
  router.push('/menu')
}
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.login-page {
  min-height: 100%;
  @include flex-center;
  padding: $spacing-lg;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: $spacing-xl;
  border-radius: $radius-lg;
  @include glass-effect(0.9, 10px);
}

.login-header {
  text-align: center;
  margin-bottom: $spacing-xl;

  h2 {
    margin: 0 0 $spacing-sm 0;
    font-size: $font-size-2xl;
    font-weight: $font-weight-bold;
    @include gradient-text($primary-gradient);
  }

  p {
    margin: 0;
    color: $text-secondary;
    font-size: $font-size-sm;
  }
}

.login-form {
  .el-form-item {
    margin-bottom: $spacing-lg;
  }
}

.login-options {
  width: 100%;
  @include flex-between;
}

.login-button {
  width: 100%;
  border-radius: $radius-full;
  background: $primary-gradient;
  border: none;
  font-weight: $font-weight-semibold;

  &:hover {
    transform: translateY(-2px);
    box-shadow: $shadow-md;
  }
}

.login-footer {
  text-align: center;
  margin-top: $spacing-lg;

  p {
    margin: 0;
    color: $text-secondary;
    font-size: $font-size-sm;
  }

  .register-link {
    color: $primary-color;
    text-decoration: none;
    font-weight: $font-weight-medium;

    &:hover {
      text-decoration: underline;
    }
  }

  .guest-btn {
    color: $text-secondary;
    font-size: $font-size-sm;
  }
}
</style>
