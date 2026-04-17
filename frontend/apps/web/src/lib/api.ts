const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const DEFAULT_STALE_TIME = 5 * 60 * 1000;
export const LONG_STALE_TIME = 10 * 60 * 1000;

interface ApiResponse<T> {
  data: T;
}

interface ApiOptions {
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE';
  body?: unknown;
  token?: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getStoredToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('token');
  }

  async request<T>(path: string, options: ApiOptions = {}): Promise<T> {
    const { method = 'GET', body, token } = options;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };

    const authToken = token || this.getStoredToken();
    if (authToken) {
      headers['Authorization'] = `Bearer ${authToken}`;
    }

    const response = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Request failed' }));
      throw new ApiError(response.status, error.error || 'Request failed');
    }

    return response.json();
  }

  get<T>(path: string, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'GET', token });
  }

  post<T>(path: string, body: unknown, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'POST', body, token });
  }

  put<T>(path: string, body: unknown, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'PUT', body, token });
  }

  delete<T>(path: string, token?: string): Promise<T> {
    return this.request<T>(path, { method: 'DELETE', token });
  }

  async getAndUnwrap<T>(path: string, token?: string): Promise<T> {
    const response = await this.get<ApiResponse<T>>(path, token);
    return response.data;
  }
}

export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

export const api = new ApiClient(API_BASE);
