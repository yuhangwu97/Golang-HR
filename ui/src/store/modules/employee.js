import { employeeService } from '@/services/employee'

const state = {
  employees: [],
  currentEmployee: null,
  loading: false,
  pagination: {
    total: 0,
    page: 1,
    pageSize: 10
  }
}

const getters = {
  employees: state => state.employees,
  currentEmployee: state => state.currentEmployee,
  loading: state => state.loading,
  pagination: state => state.pagination
}

const mutations = {
  SET_EMPLOYEES(state, employees) {
    state.employees = employees
  },
  SET_CURRENT_EMPLOYEE(state, employee) {
    state.currentEmployee = employee
  },
  SET_LOADING(state, loading) {
    state.loading = loading
  },
  SET_PAGINATION(state, pagination) {
    state.pagination = { ...state.pagination, ...pagination }
  }
}

const actions = {
  async fetchEmployees({ commit }, params = {}) {
    commit('SET_LOADING', true)
    try {
      const response = await employeeService.getEmployees(params)
      commit('SET_EMPLOYEES', response.data)
      commit('SET_PAGINATION', {
        total: response.total,
        page: response.page,
        pageSize: response.pageSize
      })
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async fetchEmployee({ commit }, id) {
    commit('SET_LOADING', true)
    try {
      const response = await employeeService.getEmployee(id)
      commit('SET_CURRENT_EMPLOYEE', response.data)
      return response
    } catch (error) {
      throw error
    } finally {
      commit('SET_LOADING', false)
    }
  },

  async createEmployee({ dispatch }, employeeData) {
    try {
      const response = await employeeService.createEmployee(employeeData)
      await dispatch('fetchEmployees')
      return response
    } catch (error) {
      throw error
    }
  },

  async updateEmployee({ dispatch }, { id, employeeData }) {
    try {
      const response = await employeeService.updateEmployee(id, employeeData)
      await dispatch('fetchEmployees')
      return response
    } catch (error) {
      throw error
    }
  },

  async deleteEmployee({ dispatch }, id) {
    try {
      await employeeService.deleteEmployee(id)
      await dispatch('fetchEmployees')
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