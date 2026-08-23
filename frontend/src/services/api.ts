import type {
  CreateUserRequest,
  CreateUserResponse,
  LoginUserRequest,
  CourseCreationRequest,
  CourseCreationResponse,
  CourseData,
  CourseViewerData,
  UserStats,
  ProgressResponse,
  VideoViewedResponse,
  UserSession,
} from '../types';

const STORAGE_KEYS = {
  TOKEN: 'course_tracker_jwt',
  API_BASE_URL: 'course_tracker_api_url',
  USER_INFO: 'course_tracker_user_info',
};

const DEFAULT_API_BASE_URL = '/api';

export function getApiBaseUrl(): string {
  const customUrl = localStorage.getItem(STORAGE_KEYS.API_BASE_URL);
  if (customUrl && customUrl.trim()) {
    return customUrl.trim().replace(/\/+$/, '');
  }
  return DEFAULT_API_BASE_URL;
}

export function setApiBaseUrl(url: string): void {
  if (!url || url.trim() === '') {
    localStorage.removeItem(STORAGE_KEYS.API_BASE_URL);
  } else {
    localStorage.setItem(STORAGE_KEYS.API_BASE_URL, url.trim().replace(/\/+$/, ''));
  }
}

export function getStoredToken(): string | null {
  return localStorage.getItem(STORAGE_KEYS.TOKEN);
}

export function setStoredToken(token: string): void {
  localStorage.setItem(STORAGE_KEYS.TOKEN, token);
}

export function removeStoredToken(): void {
  localStorage.removeItem(STORAGE_KEYS.TOKEN);
  localStorage.removeItem(STORAGE_KEYS.USER_INFO);
}

export function parseJwt(token: string): { user_id?: number; exp?: number } | null {
  try {
    const base64Url = token.split('.')[1];
    if (!base64Url) return null;
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
      window
        .atob(base64)
        .split('')
        .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

export class ApiError extends Error {
  code: string;
  status: number;

  constructor(message: string, code = 'API_ERROR', status = 500) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

async function request<T>(
  endpoint: string,
  options: RequestInit = {},
  isTextResponse = false
): Promise<T> {
  const baseUrl = getApiBaseUrl();
  const token = getStoredToken();

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const url = `${baseUrl}${endpoint}`;

  try {
    const finalUrl = url.startsWith('/api') ? url : `/api${url.startsWith('/') ? '' : '/'}${url}`;
    const response = await fetch(finalUrl, {
      ...options,
      headers,
    });

    if (!response.ok) {
      let errorCode = `HTTP_${response.status}`;
      let errorMessage = `Request failed with status ${response.status}`;

      try {
        const errorData = await response.json();
        if (errorData?.error?.message) {
          errorMessage = errorData.error.message;
          errorCode = errorData.error.code || errorCode;
        } else if (errorData?.message) {
          errorMessage = errorData.message;
        }
      } catch {
        try {
          const text = await response.text();
          if (text) errorMessage = text;
        } catch {
          // ignore
        }
      }

      throw new ApiError(errorMessage, errorCode, response.status);
    }

    if (isTextResponse) {
      const text = await response.text();
      return text as unknown as T;
    }

    const text = await response.text();
    if (!text || text.trim() === '') {
      return {} as T;
    }

    return JSON.parse(text) as T;
  } catch (err: unknown) {
    if (err instanceof ApiError) {
      throw err;
    }
    const message = err instanceof Error ? err.message : 'Network error or server unreachable';
    throw new ApiError(
      `${message}. Please verify the Go backend is running at ${baseUrl}.`,
      'NETWORK_ERROR',
      0
    );
  }
}

export const api = {
  // Auth & User
  async register(data: CreateUserRequest): Promise<CreateUserResponse> {
    return request<CreateUserResponse>('/users/create', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async login(data: LoginUserRequest): Promise<UserSession> {
    const token = await request<string>(
      '/users/login',
      {
        method: 'POST',
        body: JSON.stringify(data),
      },
      true
    );

    const cleanToken = token.trim();
    setStoredToken(cleanToken);

    const claims = parseJwt(cleanToken);
    const session: UserSession = {
      userId: claims?.user_id || 0,
      email: data.email,
      token: cleanToken,
      exp: claims?.exp,
    };

    localStorage.setItem(STORAGE_KEYS.USER_INFO, JSON.stringify(session));
    return session;
  },

  async deleteAccount(): Promise<{ message: string }> {
    const res = await request<{ message: string }>('/users/delete', {
      method: 'POST',
    });
    removeStoredToken();
    return res;
  },

  // Courses
  async getMyCourses(): Promise<CourseData[]> {
    const courses = await request<CourseData[]>('/courses/mycourses', {
      method: 'POST',
      body: JSON.stringify({}),
    });
    return Array.isArray(courses) ? courses : [];
  },

  async createCourse(data: CourseCreationRequest): Promise<CourseCreationResponse> {
    return request<CourseCreationResponse>('/courses/create', {
      method: 'POST',
      body: JSON.stringify({
        title: data.title,
        url: data.url,
        playlist_url: data.url,
      }),
    });
  },

  async getCourseDetail(courseId: number): Promise<CourseViewerData> {
    return request<CourseViewerData>('/courses/detail', {
      method: 'POST',
      body: JSON.stringify({ course_id: courseId }),
    });
  },

  async startCourse(courseId: number): Promise<{ message: string }> {
    return request<{ message: string }>('/courses/start', {
      method: 'POST',
      body: JSON.stringify({ course_id: courseId }),
    });
  },

  async deleteCourse(courseId: number): Promise<{ message: string }> {
    return request<{ message: string }>('/courses/delete', {
      method: 'POST',
      body: JSON.stringify({ course_id: courseId }),
    });
  },

  async updateProgress(courseId: number, percentageToAdd: number): Promise<ProgressResponse> {
    return request<ProgressResponse>('/courses/progress', {
      method: 'POST',
      body: JSON.stringify({
        course_id: courseId,
        percentage_to_add: percentageToAdd,
      }),
    });
  },

  async recordVideoViewed(courseId: number, videoId: string): Promise<VideoViewedResponse> {
    return request<VideoViewedResponse>('/courses/videos/viewed', {
      method: 'POST',
      body: JSON.stringify({
        course_id: courseId,
        video_id: videoId,
      }),
    });
  },

  // Gamification & Stats
  async getMyStats(): Promise<UserStats> {
    return request<UserStats>('/stats/me', {
      method: 'GET',
    });
  },

  // Health check
  async checkBackendHealth(): Promise<{ ok: boolean; latency: number; url: string }> {
    const baseUrl = getApiBaseUrl();
    const start = performance.now();
    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 3500);

      const res = await fetch(`${baseUrl}/stats/me`, {
        method: 'OPTIONS',
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      const latency = Math.round(performance.now() - start);
      return { ok: res.status < 500, latency, url: baseUrl };
    } catch {
      return { ok: false, latency: 0, url: baseUrl };
    }
  },
};
