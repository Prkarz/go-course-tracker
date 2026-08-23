export interface UserStats {
  user_id: number;
  streak_count: number;
  total_points: number;
  last_active_date: string;
}

export interface CourseData {
  id: number;
  owner_id?: number;
  url?: string;
  title?: string;
  summary?: string;
  tags: string[];
  completion_percent: number;
  started_at?: string;
  is_started: boolean;
  thumbnail_url?: string;
  first_video_id?: string;
}

export interface VideoData {
  duration: string;
  title: string;
  video_id: string;
  index: number;
  status: boolean; // is_completed
}

export interface CourseViewerData {
  title: string;
  summary: string;
  playlist: VideoData[];
}

export interface CreateUserRequest {
  username: string;
  email: string;
  password: string;
}

export interface CreateUserResponse {
  message: string;
  reactivated?: boolean;
  user_id?: number;
}

export interface LoginUserRequest {
  email: string;
  password: string;
}

export interface CourseCreationRequest {
  url: string;
  playlist_url?: string;
  title: string;
}

export interface CourseCreationResponse {
  course_id: number;
  created: boolean;
  already_exists: boolean;
  message: string;
}

export interface StartCourseRequest {
  course_id: number;
}

export interface UpdateProgressRequest {
  course_id: number;
  percentage_to_add: number;
}

export interface ProgressResponse {
  message: string;
  points_earned: number;
}

export interface VideoViewedRequest {
  course_id: number;
  video_id: string;
}

export interface VideoViewedResponse {
  message: string;
  points_earned: number;
  new_video: boolean;
}

export interface APIErrorDetail {
  code: string;
  message: string;
}

export interface APIErrorResponse {
  error: APIErrorDetail;
}

export interface UserSession {
  userId: number;
  email?: string;
  username?: string;
  token: string;
  exp?: number;
}

export type ActiveTab = 'dashboard' | 'courses' | 'player' | 'achievements' | 'settings';

export type ToastType = 'success' | 'error' | 'info' | 'reward' | 'warning';

export interface ToastItem {
  id: string;
  type: ToastType;
  title: string;
  message?: string;
  points?: number;
  timestamp: number;
}

export interface BadgeItem {
  id: string;
  title: string;
  description: string;
  icon: string;
  isUnlocked: boolean;
  currentProgress: number;
  targetProgress: number;
  category: 'streak' | 'points' | 'courses' | 'general';
}

export const API_VERSION = '1.0.0';
