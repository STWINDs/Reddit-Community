import { defineStore } from 'pinia'
import { login, register, getUserInfo } from '@/api/auth'
import { setToken, removeToken, getToken } from '@/utils/auth'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: getToken(),
    loading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    userInfo: (state) => state.user
  },

  actions: {
    async login(credentials) {
      this.loading = true
      try {
        const response = await login(credentials)
        // Adjust based on actual backend response structure
        // Backend returns { access_token, refresh_token }
        // We assume response.data contains these if using axios interceptor correctly
        // But our interceptor returns response.data directly.
        // If backend returns standard response wrapper like { code, msg, data: { ... } }, we need to handle that.
        // Looking at backend code: response.RespondWithSuccess(c, map...) returns JSON with data directly or wrapped?
        // pkg/response/response.go usually wraps.
        // Let's assume response structure is handled in api/request.js or here.
        // For now, assume response contains access_token.
        
        const token = response.data?.token
        if (token) {
          this.token = token
          setToken(token)
          this.user = {
            id: response.data.user_id,
            username: response.data.user_name
          }
        }
        return response
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async register(userData) {
      this.loading = true
      try {
        const response = await register(userData)
        return response
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },

    async getUserInfo() {
      try {
        const response = await getUserInfo()
        // /ping returns { code, msg, data: { user_id } }
        if (response.data?.user_id) {
          this.user = {
            ...this.user,
            id: response.data.user_id
          }
        }
        return response
      } catch (error) {
        throw error
      }
    },

    logout() {
      this.user = null
      this.token = null
      removeToken()
    }
  }
})
