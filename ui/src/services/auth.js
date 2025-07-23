import { apiService } from './api'

class AuthService {
  async login(credentials) {
    return apiService.post('/auth/login', credentials)
  }

  async register(userData) {
    return apiService.post('/auth/register', userData)
  }

  async validateToken() {
    return apiService.post('/auth/validate')
  }

  async logout() {
    return apiService.post('/auth/logout')
  }

  async updateProfile(profileData) {
    return apiService.put('/auth/profile', profileData)
  }

  async changePassword(oldPassword, newPassword) {
    return apiService.post('/auth/change-password', {
      oldPassword,
      newPassword
    })
  }

  async resetPassword(email) {
    return apiService.post('/auth/reset-password', { email })
  }

  async confirmResetPassword(token, newPassword) {
    return apiService.post('/auth/confirm-reset-password', {
      token,
      newPassword
    })
  }
}

export const authService = new AuthService()
export default authService