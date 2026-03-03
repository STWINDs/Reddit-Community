import axios from 'axios'
import { getToken, removeToken } from '@/utils/auth'
import { useAuthStore } from '@/store/modules/auth'
import router from '@/router'

const service = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 5000
})

// 请求拦截器
service.interceptors.request.use(
  config => {
    const token = getToken()
    if (token) {
      config.headers['Authorization'] = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  response => {
    const res = response.data
    // 如果业务状态码不是 1000，则认为是一个错误
    if (res.code !== 1000) {
      let msg = res.msg || '请求失败'
      // 处理校验错误：如果 msg 是对象（校验错误），则提取所有错误信息并合并
      if (typeof msg === 'object' && msg !== null) {
        msg = Object.values(msg).join(', ')
      }
      return Promise.reject(new Error(msg))
    }
    return res
  },
  error => {
    if (error.response?.status === 401) {
      removeToken()
      // We need to access store inside the interceptor or handling it via router/store
      // To avoid circular dependency, we might import store instance if exported, or use useAuthStore inside function
      // Here we assume Pinia is active
      const authStore = useAuthStore()
      authStore.logout()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default service
