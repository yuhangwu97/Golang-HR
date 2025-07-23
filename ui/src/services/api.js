import axios from "axios";
import { Message } from "element-ui";

class ApiService {
  constructor() {
    this.api = axios.create({
      baseURL: "/api/v1",
      timeout: 10000,
    });

    this.setupInterceptors();
  }

  setupInterceptors() {
    // Request interceptor
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem("token");
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // Response interceptor
    this.api.interceptors.response.use(
      (response) => {
        return response.data;
      },
      (error) => {
        if (error.response) {
          const { status, data } = error.response;

          switch (status) {
            case 401:
              Message.error("未授权，请重新登录");
              localStorage.removeItem("token");
              window.location.href = "#/login";
              break;
            case 403:
              Message.error("权限不足");
              break;
            case 404:
              Message.error("请求的资源不存在");
              break;
            case 500:
              Message.error("服务器内部错误");
              break;
            default:
              Message.error(data?.message || "请求失败");
          }
        } else if (error.request) {
          Message.error("网络错误，请检查网络连接");
        } else {
          Message.error("请求配置错误");
        }

        return Promise.reject(error);
      }
    );
  }

  get(url, config) {
    const token = localStorage.getItem("token");
    if (token) {
      if (!config) {
        config = {};
      }
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return this.api.get(url, config);
  }

  post(url, data, config) {
    const token = localStorage.getItem("token");
    if (token) {
      if (!config) {
        config = {};
      }
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return this.api.post(url, data, config);
  }

  put(url, data, config) {
    const token = localStorage.getItem("token");
    if (token) {
      if (!config) {
        config = {};
      }
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return this.api.put(url, data, config);
  }

  delete(url, config) {
    const token = localStorage.getItem("token");
    if (token) {
      if (!config) {
        config = {};
      }
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return this.api.delete(url, config);
  }

  patch(url, data, config) {
    const token = localStorage.getItem("token");
    if (token) {
      if (!config) {
        config = {};
      }
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return this.api.patch(url, data, config);
  }
}

export const apiService = new ApiService();
export default apiService;
