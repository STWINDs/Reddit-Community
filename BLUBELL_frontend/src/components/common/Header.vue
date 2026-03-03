<template>
  <header class="app-header">
    <div class="header-content">
      <div class="logo">
        <router-link to="/home">Bluebell</router-link>
      </div>
      <div class="user-info" v-if="authStore.isAuthenticated">
        <el-dropdown @command="handleCommand">
          <span class="user-dropdown">
            <el-avatar :size="32" :src="authStore.userInfo?.avatar">
              {{ authStore.userInfo?.username?.charAt(0).toUpperCase() }}
            </el-avatar>
            <span class="username">{{ authStore.userInfo?.username }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人资料</el-dropdown-item>
              <el-dropdown-item command="settings">设置</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="auth-actions" v-else>
        <router-link to="/login">登录</router-link>
        <router-link to="/register">注册</router-link>
      </div>
    </div>
  </header>
</template>

<script setup>
import { useAuthStore } from '@/store/modules/auth'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'

const authStore = useAuthStore()
const router = useRouter()

const handleCommand = (command) => {
  switch (command) {
    case 'logout':
      handleLogout()
      break
    case 'profile':
      ElMessage.info('个人资料功能暂未开放')
      break
    case 'settings':
      ElMessage.info('设置功能暂未开放')
      break
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped lang="scss">
.app-header {
  height: 60px;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  display: flex;
  align-items: center;
  padding: 0 20px;
  
  .header-content {
    width: 100%;
    max-width: 1200px;
    margin: 0 auto;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logo {
    font-size: 24px;
    font-weight: bold;
    a {
      text-decoration: none;
      color: #409EFF;
    }
  }

  .user-dropdown {
    display: flex;
    align-items: center;
    cursor: pointer;
    
    .username {
      margin: 0 8px;
      font-size: 14px;
    }
  }
  
  .auth-actions {
    a {
      margin-left: 15px;
      text-decoration: none;
      color: #606266;
      &:hover {
        color: #409EFF;
      }
    }
  }
}
</style>
