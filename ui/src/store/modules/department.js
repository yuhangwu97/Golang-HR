import { departmentApi } from '@/services/departmentApi'

const state = {
  departments: [],
  departmentTree: [],
  currentDepartment: null,
  loading: false
}

const getters = {
  departments: state => state.departments,
  departmentTree: state => state.departmentTree,
  currentDepartment: state => state.currentDepartment,
  loading: state => state.loading
}

const mutations = {
  SET_DEPARTMENTS(state, departments) {
    state.departments = departments
  },
  SET_DEPARTMENT_TREE(state, tree) {
    state.departmentTree = tree
  },
  SET_CURRENT_DEPARTMENT(state, department) {
    state.currentDepartment = department
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  }
}

const actions = {
  async fetchDepartments({ commit }) {
    commit('SET_LOADING', true)
    try {
      const departments = await departmentApi.getDepartments()
      commit('SET_DEPARTMENTS', departments)
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async fetchDepartmentTree({ commit }) {
    commit('SET_LOADING', true)
    try {
      const tree = await departmentApi.getDepartmentTree()
      commit('SET_DEPARTMENT_TREE', tree)
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async createDepartment({ dispatch }, departmentData) {
    try {
      const response = await departmentApi.createDepartment(departmentData)
      await dispatch('fetchDepartments')
      await dispatch('fetchDepartmentTree')
      return response
    } catch (error) {
      throw error
    }
  },

  async updateDepartment({ dispatch }, { id, departmentData }) {
    try {
      const response = await departmentApi.updateDepartment(id, departmentData)
      await dispatch('fetchDepartments')
      await dispatch('fetchDepartmentTree')
      return response
    } catch (error) {
      throw error
    }
  },

  async deleteDepartment({ dispatch }, id) {
    try {
      await departmentApi.deleteDepartment(id)
      await dispatch('fetchDepartments')
      await dispatch('fetchDepartmentTree')
    } catch (error) {
      throw error
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