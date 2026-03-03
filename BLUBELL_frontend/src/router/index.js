import { createRouter, createWebHistory } from 'vue-router'
import { getToken, isTokenExpired } from '@/utils/auth'
import { useAuthStore } from '@/store/modules/auth'

const routes = [
  {
    path: '/',
    redirect: '/home'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/Register.vue'),
    meta: { requiresGuest: true }
  },
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/home/Index.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const token = getToken()
  const authStore = useAuthStore()

  // 检查token是否过期
  if (token && isTokenExpired(token)) {
    authStore.logout()
    next('/login')
    return
  }

  // 获取用户信息
  if (token && !authStore.user) {
    authStore.getUserInfo().catch(error => {
      console.warn('Failed to get user info:', error)
    })
  }

  const isAuthenticated = !!token || authStore.isAuthenticated

  // 需要认证的页面
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }

  // 游客访问的页面（已登录用户不应访问）
  if (to.meta.requiresGuest && isAuthenticated) {
    next('/home')
    return
  }

  next()
})

export default router
