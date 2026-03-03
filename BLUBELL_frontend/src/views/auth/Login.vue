<template>
  <div class="login-container">
    <div class="login-box">
      <h2>用户登录</h2>
      <LoginForm @submit="handleLogin" :loading="loading" />
      <div class="form-footer">
        <span>还没有账号？</span>
        <router-link to="/register">立即注册</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/modules/auth'
import LoginForm from '@/components/auth/LoginForm.vue'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const handleLogin = async (credentials) => {
  try {
    loading.value = true
    await authStore.login(credentials)
    ElMessage.success('登录成功')
    router.replace({ name: 'Home' })
  } catch (error) {
    ElMessage.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background-color: #f5f7fa;

  .login-box {
    width: 400px;
    padding: 40px;
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    
    h2 {
      text-align: center;
      margin-bottom: 30px;
      color: #303133;
    }

    .form-footer {
      margin-top: 20px;
      text-align: center;
      font-size: 14px;
      color: #606266;
      
      a {
        color: #409EFF;
        text-decoration: none;
        margin-left: 5px;
        
        &:hover {
          text-decoration: underline;
        }
      }
    }
  }
}
</style>
