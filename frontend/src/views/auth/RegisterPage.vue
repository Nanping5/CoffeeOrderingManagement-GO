<template>
  <div class="register-page">
    <el-card class="register-card glass-panel">
      <div class="register-header">
        <h2>用户注册</h2>
        <p>加入咖啡时光，享受专属优惠</p>
      </div>

      <el-form
        ref="registerFormRef"
        :model="formData"
        :rules="rules"
        @submit.prevent="handleRegister"
        class="register-form"
        label-position="top"
      >
        <el-form-item prop="username" label="用户名">
          <el-input
            v-model="formData.username"
            placeholder="请输入用户名"
            :prefix-icon="User"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item prop="email" label="邮箱">
          <el-input
            v-model="formData.email"
            placeholder="请输入邮箱"
            :prefix-icon="Message"
            size="large"
            clearable
          />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="12">
            <el-form-item prop="first_name" label="名字">
              <el-input
                v-model="formData.first_name"
                placeholder="名字"
                size="large"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item prop="last_name" label="姓氏">
              <el-input
                v-model="formData.last_name"
                placeholder="姓氏"
                size="large"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item prop="phone" label="手机号">
          <el-input
            v-model="formData.phone"
            placeholder="请输入手机号"
            :prefix-icon="Phone"
            size="large"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password" label="密码">
          <el-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码（至少8位）"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>

        <el-form-item prop="confirmPassword" label="确认密码">
          <el-input
            v-model="formData.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            :prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            @click="handleRegister"
            class="register-button"
          >
            {{ loading ? '注册中...' : '立即注册' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="register-footer">
        <p>已有账户？
          <router-link to="/auth/login" class="login-link">立即登录</router-link>
        </p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Message, Lock, Phone } from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

const registerFormRef = ref(null)
const loading = ref(false)

const formData = reactive({
  username: '',
  email: '',
  first_name: '',
  last_name: '',
  phone: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== formData.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

// 手机号验证器（允许为空）
const validatePhone = (rule, value, callback) => {
  if (!value || value === '') {
    callback() // 允许为空
  } else if (!/^1[3-9]\d{9}$/.test(value)) {
    callback(new Error('请输入正确的11位手机号'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度3-20位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { validator: validatePhone, trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 8, message: '密码至少8位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return
  
  try {
    await registerFormRef.value.validate()
    loading.value = true
    
    const result = await authStore.register({
      username: formData.username,
      email: formData.email,
      password: formData.password,
      first_name: formData.first_name,
      last_name: formData.last_name,
      phone: formData.phone
    })
    
    if (result.success) {
      ElMessage.success('注册成功！已赠送50积分')
      router.push('/menu')
    } else {
      ElMessage.error(result.errors?.[0] || '注册失败')
    }
  } catch (error) {
    if (error !== false) {
      ElMessage.error('注册失败，请重试')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
@import '@/styles/variables.scss';
@import '@/styles/mixins.scss';

.register-page {
  min-height: 100%;
  @include flex-center;
  padding: $spacing-lg;
}

.register-card {
  width: 100%;
  max-width: 480px;
  padding: $spacing-xl;
  border-radius: $radius-lg;
  @include glass-effect(0.9, 10px);
}

.register-header {
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

.register-form {
  .el-form-item {
    margin-bottom: $spacing-md;
  }
}

.register-button {
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

.register-footer {
  text-align: center;
  margin-top: $spacing-lg;

  p {
    margin: 0;
    color: $text-secondary;
    font-size: $font-size-sm;
  }

  .login-link {
    color: $primary-color;
    text-decoration: none;
    font-weight: $font-weight-medium;

    &:hover {
      text-decoration: underline;
    }
  }
}

// 移动端适配
@include respond-to(md) {
  .register-page {
    padding: $spacing-md;
    align-items: flex-start;
    padding-top: $spacing-lg;
  }
  
  .register-card {
    padding: $spacing-lg;
    max-width: 100%;
    
    :deep(.el-card__body) {
      padding: 0;
    }
  }
  
  .register-header {
    margin-bottom: $spacing-lg;
    
    h2 {
      font-size: $font-size-xl;
    }
  }
  
  .register-form {
    .el-form-item {
      margin-bottom: $spacing-sm;
    }
    
    // 姓名字段在移动端堆叠显示
    .el-row {
      .el-col {
        width: 100%;
        max-width: 100%;
        flex: 0 0 100%;
      }
    }
  }
}

@media (max-width: 480px) {
  .register-page {
    padding: $spacing-sm;
    padding-top: $spacing-md;
  }
  
  .register-card {
    padding: $spacing-md;
  }
  
  .register-header {
    margin-bottom: $spacing-md;
    
    h2 {
      font-size: $font-size-lg;
    }
    
    p {
      font-size: $font-size-xs;
    }
  }
}
</style>
