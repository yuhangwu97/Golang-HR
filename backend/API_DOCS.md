# Auth API Documentation

## Authentication Endpoints

### 1. Register User
- **POST** `/api/v1/auth/register`
- **Body**: 
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```
- **Response**: 
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "status": "active"
  }
}
```

### 2. Login
- **POST** `/api/v1/auth/login`
- **Body**: 
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```
- **Response**: 
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    }
  }
}
```

### 3. Validate Token
- **POST** `/api/v1/auth/validate`
- **Headers**: 
```
Authorization: Bearer <token>
```
- **Response**: 
```json
{
  "success": true,
  "message": "Token is valid",
  "data": {
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    },
    "valid": true
  }
}
```

### 4. Refresh Token
- **POST** `/api/v1/auth/refresh`
- **Headers**: 
```
Authorization: Bearer <token>
```
- **Response**: 
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "status": "active"
    }
  }
}
```

### 5. Logout
- **POST** `/api/v1/auth/logout`
- **Headers**: 
```
Authorization: Bearer <token>
```
- **Response**: 
```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

## Error Responses

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Invalid or expired token",
  "error": "unauthorized"
}
```

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request data",
  "error": "bad_request"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "internal_server_error"
}
```

## Usage Examples

### Frontend Integration (JavaScript)

```javascript
// Login and store token
const loginResponse = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    email: 'john@example.com',
    password: 'password123'
  })
});

const loginData = await loginResponse.json();
const token = loginData.data.token;

// Store token in localStorage
localStorage.setItem('token', token);

// Validate token before making authenticated requests
const validateResponse = await fetch('/api/v1/auth/validate', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const validateData = await validateResponse.json();
if (validateData.success) {
  // Token is valid, proceed with authenticated requests
  console.log('User:', validateData.data.user);
} else {
  // Token is invalid, redirect to login
  window.location.href = '/login';
}
```

### Vue.js Store Integration

```javascript
// In your Vuex store
export const actions = {
  async validateToken({ commit }) {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('No token found');
    }

    try {
      const response = await this.$axios.post('/api/v1/auth/validate', {}, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.data.success) {
        commit('SET_USER', response.data.data.user);
        return response.data.data.user;
      } else {
        throw new Error('Token validation failed');
      }
    } catch (error) {
      localStorage.removeItem('token');
      throw error;
    }
  }
};
```

## Security Features

1. **JWT Token Expiration**: Tokens expire after 24 hours
2. **Redis Session Management**: Sessions are stored in Redis for fast validation
3. **Password Hashing**: Passwords are hashed using bcrypt
4. **Token Refresh**: Tokens can be refreshed without re-login
5. **Session Invalidation**: Logout removes session from Redis

## Notes

- All endpoints return consistent JSON responses
- Tokens must be included in the `Authorization` header with `Bearer` prefix
- The validate endpoint is primarily used for frontend route guards
- Session data is stored in Redis with automatic expiration
- User passwords are never returned in API responses