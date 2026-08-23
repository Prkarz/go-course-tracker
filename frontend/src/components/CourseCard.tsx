import React, { useState } from 'react';
import type { CourseData } from '../types';
import {
  BookOpen,
  PlayCircle,
  CheckCircle,
  ExternalLink,
  Trash2,
  Tag,
  ChevronDown,
  ChevronUp,
  Sparkles,
  Play,
  Video,
} from 'lucide-react';

interface CourseCardProps {
  course: CourseData;
  onStartCourse: (courseId: number) => void;
  onOpenPlayer: (course: CourseData) => void;
  onDeleteCourse: (courseId: number) => void;
  isStarting?: boolean;
}

export const CourseCard: React.FC<CourseCardProps> = ({
  course,
  onStartCourse,
  onOpenPlayer,
  onDeleteCourse,
  isStarting = false,
}) => {
  const [isSummaryExpanded, setIsSummaryExpanded] = useState(false);
  const [imgError, setImgError] = useState(false);
  const isCompleted = course.completion_percent >= 100;
  const isStarted = course.is_started;

  const playlistUrl = course.url || '';
  const playlistMatch = playlistUrl.match(/[?&]list=([^&]+)/);
  const playlistId = playlistMatch ? playlistMatch[1] : '';

  // Determine official YouTube thumbnail URL
  const getThumbnailUrl = (): string | null => {
    if (course.thumbnail_url && course.thumbnail_url.trim()) {
      return course.thumbnail_url.trim();
    }
    if (course.first_video_id && course.first_video_id.trim()) {
      return `https://img.youtube.com/vi/${course.first_video_id.trim()}/hqdefault.jpg`;
    }
    // Check if course URL contains direct video ID (e.g. v=VIDEO_ID)
    const videoMatch = playlistUrl.match(/[?&]v=([^&]+)/);
    if (videoMatch && videoMatch[1]) {
      return `https://img.youtube.com/vi/${videoMatch[1]}/hqdefault.jpg`;
    }
    return null;
  };

  const thumbnailUrl = getThumbnailUrl();

  return (
    <div className="rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 transition-colors flex flex-col justify-between overflow-hidden">
      <div>
        {/* Course Card Header with Official YouTube Thumbnail */}
        <div className="relative h-40 bg-slate-950 overflow-hidden border-b border-slate-800">
          {thumbnailUrl && !imgError ? (
            <>
              <img
                src={thumbnailUrl}
                alt={course.title || 'Course Thumbnail'}
                onError={() => setImgError(true)}
                loading="lazy"
                className="w-full h-full object-cover object-center"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-slate-950 via-slate-950/40 to-transparent" />
            </>
          ) : (
            <div className="w-full h-full bg-slate-950 p-4 flex flex-col justify-between relative overflow-hidden">
              <div className="absolute bottom-3 right-3 text-slate-800">
                <Video className="w-12 h-12" />
              </div>
            </div>
          )}

          {/* Top Bar Badges and Actions */}
          <div className="absolute top-2.5 left-2.5 right-2.5 z-10 flex items-center justify-between">
            <div>
              {isCompleted ? (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-semibold bg-emerald-500 text-slate-950">
                  <CheckCircle className="w-3 h-3" />
                  <span>Completed</span>
                </span>
              ) : isStarted ? (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-semibold bg-brand-500 text-slate-950">
                  <PlayCircle className="w-3 h-3" />
                  <span>In Progress</span>
                </span>
              ) : (
                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-semibold bg-slate-900 text-slate-300 border border-slate-700">
                  <span>Not Started</span>
                </span>
              )}
            </div>

            <div className="flex items-center gap-1">
              {playlistUrl && (
                <a
                  href={playlistUrl}
                  target="_blank"
                  rel="noreferrer"
                  className="p-1 rounded bg-slate-950/90 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800 transition-colors"
                  title="Open YouTube Playlist"
                >
                  <ExternalLink className="w-3.5 h-3.5" />
                </a>
              )}
              <button
                onClick={() => onDeleteCourse(course.id)}
                className="p-1 rounded bg-slate-950/90 hover:bg-rose-950 text-slate-300 hover:text-rose-400 border border-slate-800 transition-colors"
                title="Delete Course (Earned XP is preserved)"
              >
                <Trash2 className="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          {/* Bottom Title & Playlist ID on Header */}
          <div className="absolute bottom-2.5 left-2.5 right-2.5 z-10">
            <h3 className="text-sm font-bold text-white line-clamp-1 leading-snug">
              {course.title || 'Untitled Course'}
            </h3>
            {playlistId && (
              <p className="text-[10px] text-slate-400 font-mono truncate">
                Playlist: {playlistId}
              </p>
            )}
          </div>
        </div>

        {/* Card Body */}
        <div className="p-4 space-y-3">
          {/* Progress Bar */}
          <div>
            <div className="flex items-center justify-between text-xs font-medium mb-1">
              <span className="text-slate-400">Progress</span>
              <span className={isCompleted ? 'text-emerald-400 font-bold' : 'text-brand-400 font-bold'}>
                {Math.round(course.completion_percent)}%
              </span>
            </div>
            <div className="w-full h-1.5 rounded-full bg-slate-800 overflow-hidden">
              <div
                className={`h-full rounded-full ${
                  isCompleted ? 'bg-emerald-500' : 'bg-brand-500'
                }`}
                style={{ width: `${Math.min(100, Math.max(0, course.completion_percent))}%` }}
              />
            </div>
          </div>

          {/* AI Summary Expander */}
          {course.summary && (
            <div className="rounded-lg bg-slate-950 border border-slate-800 p-2.5 text-xs">
              <div
                className="flex items-center justify-between cursor-pointer select-none text-slate-300 hover:text-white font-medium"
                onClick={() => setIsSummaryExpanded(!isSummaryExpanded)}
              >
                <div className="flex items-center gap-1 text-brand-400">
                  <Sparkles className="w-3 h-3" />
                  <span>AI Overview</span>
                </div>
                {isSummaryExpanded ? (
                  <ChevronUp className="w-3.5 h-3.5 text-slate-400" />
                ) : (
                  <ChevronDown className="w-3.5 h-3.5 text-slate-400" />
                )}
              </div>
              <p
                className={`mt-1 text-slate-400 leading-relaxed ${
                  isSummaryExpanded ? 'line-clamp-none' : 'line-clamp-2'
                }`}
              >
                {course.summary}
              </p>
            </div>
          )}

          {/* Topic Tags */}
          {course.tags && course.tags.length > 0 && (
            <div className="flex flex-wrap gap-1">
              {course.tags.slice(0, 3).map((tag, idx) => (
                <span
                  key={idx}
                  className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-medium bg-slate-800 text-slate-300 border border-slate-700"
                >
                  <Tag className="w-2.5 h-2.5 text-slate-400" />
                  <span>{tag}</span>
                </span>
              ))}
              {course.tags.length > 3 && (
                <span className="px-1 py-0.5 text-[10px] text-slate-500">
                  +{course.tags.length - 3}
                </span>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Card Action Button */}
      <div className="p-4 pt-0">
        {!isStarted ? (
          <button
            onClick={() => onStartCourse(course.id)}
            disabled={isStarting}
            className="w-full py-2 rounded-lg bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-xs transition-colors flex items-center justify-center gap-1.5 disabled:opacity-60"
          >
            <Play className="w-3.5 h-3.5 fill-slate-950" />
            <span>{isStarting ? 'Starting...' : 'Start Course'}</span>
          </button>
        ) : (
          <button
            onClick={() => onOpenPlayer(course)}
            className="w-full py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-white font-semibold text-xs border border-slate-700 transition-colors flex items-center justify-center gap-1.5"
          >
            <PlayCircle className="w-3.5 h-3.5 text-brand-400" />
            <span>{isCompleted ? 'Review Lessons' : 'Open Course Player'}</span>
          </button>
        )}
      </div>
    </div>
  );
};
