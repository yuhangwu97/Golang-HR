import { apiService } from './api'

class SystemApiService {
  // 用户管理
  async getUsers(params = {}) {
    return apiService.get('/system/users', { params })
  }

  async getUser(id) {
    return apiService.get(`/system/users/${id}`)
  }

  async createUser(userData) {
    return apiService.post('/system/users', userData)
  }

  async updateUser(id, userData) {
    return apiService.put(`/system/users/${id}`, userData)
  }

  async deleteUser(id) {
    return apiService.delete(`/system/users/${id}`)
  }

  async resetUserPassword(id, newPassword) {
    return apiService.put(`/system/users/${id}/reset-password`, { newPassword })
  }

  async toggleUserStatus(id) {
    return apiService.put(`/system/users/${id}/toggle-status`)
  }

  // 角色管理
  async getRoles(params = {}) {
    return apiService.get('/system/roles', { params })
  }

  async getRole(id) {
    return apiService.get(`/system/roles/${id}`)
  }

  async createRole(roleData) {
    return apiService.post('/system/roles', roleData)
  }

  async updateRole(id, roleData) {
    return apiService.put(`/system/roles/${id}`, roleData)
  }

  async deleteRole(id) {
    return apiService.delete(`/system/roles/${id}`)
  }

  async assignRolesToUser(userId, roleIds) {
    return apiService.post(`/system/users/${userId}/roles`, { roleIds })
  }

  async getUserRoles(userId) {
    return apiService.get(`/system/users/${userId}/roles`)
  }

  async getRoleUsers(roleId) {
    return apiService.get(`/system/roles/${roleId}/users`)
  }

  // 权限管理
  async getPermissions(params = {}) {
    return apiService.get('/system/permissions', { params })
  }

  async getPermission(id) {
    return apiService.get(`/system/permissions/${id}`)
  }

  async createPermission(permissionData) {
    return apiService.post('/system/permissions', permissionData)
  }

  async updatePermission(id, permissionData) {
    return apiService.put(`/system/permissions/${id}`, permissionData)
  }

  async deletePermission(id) {
    return apiService.delete(`/system/permissions/${id}`)
  }

  async assignPermissionsToRole(roleId, permissionIds) {
    return apiService.post(`/system/roles/${roleId}/permissions`, { permissionIds })
  }

  async getRolePermissions(roleId) {
    return apiService.get(`/system/roles/${roleId}/permissions`)
  }

  async getPermissionTree() {
    return apiService.get('/system/permissions/tree')
  }

  // 菜单管理
  async getMenus() {
    return apiService.get('/system/menus')
  }

  async getUserMenus() {
    return apiService.get('/system/menus/user')
  }
}

export const systemApiService = new SystemApiService()
export const roleApi = {
  getRoles: (params) => systemApiService.getRoles(params),
  getRole: (id) => systemApiService.getRole(id),
  createRole: (roleData) => systemApiService.createRole(roleData),
  updateRole: (id, roleData) => systemApiService.updateRole(id, roleData),
  deleteRole: (id) => systemApiService.deleteRole(id),
  assignPermissions: (roleId, permissionIds) => systemApiService.assignPermissionsToRole(roleId, permissionIds),
  getRolePermissions: (roleId) => systemApiService.getRolePermissions(roleId),
  getRoleUsers: (roleId) => systemApiService.getRoleUsers(roleId)
}

export const permissionApi = {
  getPermissions: (params) => systemApiService.getPermissions(params),
  getPermission: (id) => systemApiService.getPermission(id),
  createPermission: (permissionData) => systemApiService.createPermission(permissionData),
  updatePermission: (id, permissionData) => systemApiService.updatePermission(id, permissionData),
  deletePermission: (id) => systemApiService.deletePermission(id),
  getPermissionTree: () => systemApiService.getPermissionTree()
}

export default systemApiService