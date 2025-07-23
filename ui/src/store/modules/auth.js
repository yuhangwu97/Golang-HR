import { authService } from '@/services/auth'

const state = {
  user: null,
  token: localStorage.getItem('token') || null,
  loading: false
}

const getters = {
  isAuthenticated: state => !!state.token && !!state.user,
  user: state => state.user,
  token: state => state.token,
  loading: state => state.loading
}

const mutations = {
  SET_USER(state, user) {
    state.user = user
  },
  SET_TOKEN(state, token) {
    state.token = token
    if (token) {
      localStorage.setItem('token', token)
    } else {
      localStorage.removeItem('token')
    }
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  }
}

const actions = {
  async login({ commit }, credentials) {
    commit('SET_LOADING', true)
    try {
      const response = await authService.login(credentials)
      // 处理API响应的数据结构
      const responseData = response.data || response
      const token = responseData.token || responseData.data?.token
      const user = responseData.user || responseData.data?.user
      
      if (token) {
        commit('SET_TOKEN', token)
      }
      if (user) {
        commit('SET_USER', user)
      }
      return responseData
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async register({ commit }, userData) {
    commit('SET_LOADING', true)
    try {
      const response = await authService.register(userData)
      commit('SET_TOKEN', response.token)
      commit('SET_USER', response.user)
      return response
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async validateToken({ commit, state }) {
    if (!state.token) {
      throw new Error('No token')
    }
    try {
      const response = await authService.validateToken()
      commit('SET_USER', response.user)
      return response
    } catch (error) {
      commit('SET_TOKEN', null)
      commit('SET_USER', null)
      throw error
    }
  },

  async logout({ commit }) {
    try {
      await authService.logout()
    } catch (error) {
      // 忽略登出错误
    } finally {
      commit('SET_TOKEN', null)
      commit('SET_USER', null)
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}