// src/lib/api/client.ts

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from "axios";
import { toast } from "sonner";
import { TokenManager } from "../auth/token-manager";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

class APIClient {
  private client: AxiosInstance;
  private isRefreshing = false;
  private refreshSubscribers: Array<(token: string) => void> = [];

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        "Content-Type": "application/json",
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = TokenManager.getAccessToken();
        if (token && !TokenManager.isTokenExpired(token)) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor to handle token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error) => {
        const original = error.config;

        if (error.response?.status === 401 && !original._retry) {
          if (this.isRefreshing) {
            // If already refreshing, wait for the new token
            return new Promise((resolve) => {
              this.refreshSubscribers.push((token: string) => {
                original.headers.Authorization = `Bearer ${token}`;
                resolve(this.client(original));
              });
            });
          }

          original._retry = true;
          this.isRefreshing = true;

          try {
            const refreshToken = TokenManager.getRefreshToken();
            if (!refreshToken) {
              throw new Error("No refresh token available");
            }

            const response = await this.client.post("/auth/refresh", {
              refresh_token: refreshToken,
            });

            const {
              access_token,
              refresh_token: newRefreshToken,
              expires_in,
            } = response.data.data;

            TokenManager.setTokens(access_token, newRefreshToken, expires_in);

            // Notify all waiting requests
            this.refreshSubscribers.forEach((callback) =>
              callback(access_token)
            );
            this.refreshSubscribers = [];

            // Retry the original request
            original.headers.Authorization = `Bearer ${access_token}`;
            return this.client(original);
          } catch (refreshError) {
            // Refresh failed, logout user
            TokenManager.clearTokens();
            window.location.href = "/auth/login";
            toast.error("Session expired. Please login again.");
            return Promise.reject(refreshError);
          } finally {
            this.isRefreshing = false;
          }
        }

        return Promise.reject(error);
      }
    );
  }

  private handleError(error: any) {
    if (error.response?.data?.message) {
      return error.response.data;
    }

    if (error.message) {
      return {
        success: false,
        message: error.message,
        error: error.message,
      };
    }

    return {
      success: false,
      message: "An unexpected error occurred",
      error: "Unknown error",
    };
  }

  async get<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.client.get(url, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async post<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.client.post(
        url,
        data,
        config
      );
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async put<T>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.client.put(
        url,
        data,
        config
      );
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  async delete<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response: AxiosResponse<T> = await this.client.delete(url, config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }

  // Method to make requests without auto token attachment (for auth endpoints)
  async publicRequest<T>(
    method: "GET" | "POST" | "PUT" | "DELETE",
    url: string,
    data?: any
  ): Promise<T> {
    try {
      const config: AxiosRequestConfig = {
        method,
        url: `${API_BASE_URL}${url}`,
        data,
        headers: {
          "Content-Type": "application/json",
        },
      };

      const response: AxiosResponse<T> = await axios(config);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  }
}

export const apiClient = new APIClient();
