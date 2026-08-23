import React, { useState, useEffect, useCallback, memo } from 'react';
import type { CourseData, CourseViewerData, VideoData } from '../types';
import { api } from '../services/api';
import { useNotification } from '../context/NotificationContext';
import { useAuth } from '../context/AuthContext';
import {
  X,
  CheckCircle2,
  Check,
  ChevronLeft,
  ChevronRight,
  Sparkles,
  Clock,
  ExternalLink,
  BookOpen,
  ListVideo,
  Tag,
  Search,
} from 'lucide-react';

interface CourseDetailModalProps {
  isOpen: boolean;
  course: CourseData | null;
  onClose: () => void;
  onCourseUpdated: () => void;
}

// Memoized Video Player component to completely eliminate video lag and unwanted iframe reloads
const MemoizedVideoPlayer = memo(
  ({ videoId, title }: { videoId: string; title: string }) => {
    return (
      <div className="relative w-full aspect-video rounded-2xl overflow-hidden bg-black shadow-2xl border border-slate-800">
        <iframe
          src={`https://www.youtube-nocookie.com/embed/${videoId}?autoplay=1&rel=0&enablejsapi=1&modestbranding=1`}
          title={title || 'Course Video'}
          className="absolute inset-0 w-full h-full"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
          allowFullScreen
          loading="lazy"
        />
      </div>
    );
  },
  (prev, next) => prev.videoId === next.videoId && prev.title === next.title
);
MemoizedVideoPlayer.displayName = 'MemoizedVideoPlayer';

export const CourseDetailModal: React.FC<CourseDetailModalProps> = ({
  isOpen,
  course,
  onClose,
  onCourseUpdated,
}) => {
  const [detailData, setDetailData] = useState<CourseViewerData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeVideoIndex, setActiveVideoIndex] = useState(0);
  const [isRecordingProgress, setIsRecordingProgress] = useState(false);
  const [playlistSearch, setPlaylistSearch] = useState('');
  const [isAiInsightsExpanded, setIsAiInsightsExpanded] = useState(true);

  const { addToast } = useNotification();
  const { refreshStats } = useAuth();

  useEffect(() => {
    let isMounted = true;
    if (isOpen && course) {
      setIsLoading(true);
      api
        .getCourseDetail(course.id)
        .then((data) => {
          if (isMounted) {
            setDetailData(data);
            if (data.playlist && data.playlist.length > 0) {
              const firstUncompletedIndex = data.playlist.findIndex((v) => !v.status);
              setActiveVideoIndex(firstUncompletedIndex >= 0 ? firstUncompletedIndex : 0);
            }
            setIsLoading(false);
          }
        })
        .catch((err) => {
          if (isMounted) {
            addToast('Failed to load course playlist', err.message, 'error');
            setIsLoading(false);
          }
        });
    }
    return () => {
      isMounted = false;
    };
  }, [isOpen, course, addToast]);

  const playlist = detailData?.playlist || [];
  const currentVideo: VideoData | undefined = playlist[activeVideoIndex];
  const completedVideosCount = playlist.filter((v) => v.status).length;
  const totalVideos = playlist.length;
  const calculatedPercent =
    totalVideos > 0
      ? Math.round((completedVideosCount / totalVideos) * 100)
      : Math.round(course?.completion_percent || 0);

  const handleMarkVideoCompleted = useCallback(
    async (video: VideoData, index: number) => {
      if (!video.video_id || !course) return;
      setIsRecordingProgress(true);

      try {
        const res = await api.recordVideoViewed(course.id, video.video_id);

        setDetailData((prev) => {
          if (!prev) return prev;
          const newPlaylist = [...prev.playlist];
          newPlaylist[index] = { ...newPlaylist[index], status: true };
          return { ...prev, playlist: newPlaylist };
        });

        addToast(
          res.new_video ? 'Lesson Completed! ⚡' : 'Lesson Already Completed',
          `+${res.points_earned || 10} XP awarded for "${video.title || 'Lesson'}"`,
          'reward',
          res.points_earned || 10
        );

        await refreshStats();
        onCourseUpdated();

        if (index < playlist.length - 1) {
          setActiveVideoIndex(index + 1);
        }
      } catch (err: unknown) {
        addToast('Progress update failed', err instanceof Error ? err.message : 'Error', 'error');
      } finally {
        setIsRecordingProgress(false);
      }
    },
    [course, playlist.length, addToast, refreshStats, onCourseUpdated]
  );

  if (!isOpen || !course) return null;

  const filteredPlaylist = playlist.filter((v) =>
    (v.title || '').toLowerCase().includes(playlistSearch.toLowerCase())
  );

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-2 sm:p-4 md:p-6 bg-slate-950/80 overflow-y-auto">
      <div className="relative w-full max-w-6xl rounded-2xl bg-slate-900 border border-slate-800 overflow-hidden flex flex-col max-h-[94vh]">
        {/* Modal Header */}
        <div className="p-4 border-b border-slate-800 bg-slate-950 flex items-center justify-between gap-4 shrink-0">
          <div className="flex items-center gap-3 min-w-0">
            <div className="w-9 h-9 rounded-lg bg-slate-800 border border-slate-700 flex items-center justify-center shrink-0">
              <BookOpen className="w-5 h-5 text-brand-400" />
            </div>
            <div className="min-w-0">
              <h2 className="text-base sm:text-lg font-bold text-white truncate">
                {detailData?.title || course.title || 'Course Player'}
              </h2>
              <div className="flex items-center gap-3 text-xs text-slate-400 mt-0.5">
                <span>
                  {completedVideosCount} / {totalVideos} lessons completed
                </span>
                <span>•</span>
                <span className="font-semibold text-brand-400">{calculatedPercent}% overall progress</span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-2 shrink-0">
            {course.url && (
              <a
                href={course.url}
                target="_blank"
                rel="noreferrer"
                className="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-xs font-semibold text-slate-200 transition-colors"
              >
                <ExternalLink className="w-3.5 h-3.5" />
                <span>YouTube Playlist</span>
              </a>
            )}
            <button
              onClick={onClose}
              className="p-2 rounded-xl text-slate-400 hover:text-slate-200 hover:bg-slate-800 transition-colors"
            >
              <X className="w-5 h-5" />
            </button>
          </div>
        </div>

        {/* Dynamic Progress Bar */}
        <div className="w-full h-1.5 bg-slate-800 shrink-0 overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-brand-500 via-emerald-400 to-accent-400 transition-all duration-500"
            style={{ width: `${Math.min(100, Math.max(0, calculatedPercent))}%` }}
          />
        </div>

        {isLoading ? (
          <div className="flex-1 flex flex-col items-center justify-center p-12 space-y-4">
            <div className="w-10 h-10 border-4 border-brand-500/20 border-t-brand-500 rounded-full animate-spin" />
            <p className="text-sm font-semibold text-slate-300">Loading course curriculum & video stream...</p>
          </div>
        ) : (
          <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 min-h-0 overflow-y-auto lg:overflow-hidden">
            {/* Main Player Column */}
            <div className="lg:col-span-8 p-4 sm:p-6 overflow-y-auto space-y-4 border-b lg:border-b-0 lg:border-r border-slate-800/80">
              {currentVideo ? (
                <>
                  {/* Memoized Video Player (Zero-Lag Playback) */}
                  <MemoizedVideoPlayer
                    videoId={currentVideo.video_id}
                    title={currentVideo.title || 'Course Video'}
                  />

                  {/* Video Control Bar & Mark Completed Action */}
                  <div className="rounded-2xl bg-slate-950/70 border border-slate-800/80 p-4 space-y-3">
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                      <div>
                        <div className="flex items-center gap-2 flex-wrap">
                          <span className="px-2 py-0.5 rounded text-[11px] font-bold bg-accent-500/20 text-accent-300 border border-accent-500/30">
                            Lesson #{activeVideoIndex + 1} of {totalVideos}
                          </span>
                          {currentVideo.duration && (
                            <span className="flex items-center gap-1 text-xs text-slate-400 font-mono">
                              <Clock className="w-3 h-3" />
                              <span>{currentVideo.duration}</span>
                            </span>
                          )}
                          {currentVideo.status && (
                            <span className="flex items-center gap-1 text-xs font-bold text-emerald-400">
                              <CheckCircle2 className="w-3.5 h-3.5" />
                              <span>Completed</span>
                            </span>
                          )}
                        </div>
                        <h3 className="text-base sm:text-lg font-bold text-white mt-1.5 leading-snug">
                          {currentVideo.title}
                        </h3>
                      </div>

                      <div className="shrink-0">
                        <button
                          onClick={() => handleMarkVideoCompleted(currentVideo, activeVideoIndex)}
                          disabled={isRecordingProgress}
                          className={`px-4 py-2 rounded-lg font-bold text-xs transition-colors flex items-center gap-1.5 ${
                            currentVideo.status
                              ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 hover:bg-emerald-500/30'
                              : 'bg-brand-500 hover:bg-brand-400 text-slate-950'
                          }`}
                        >
                          <Check className="w-4 h-4" />
                          <span>
                            {currentVideo.status ? 'Completed (+10 XP)' : 'Mark as Completed (+10 XP)'}
                          </span>
                        </button>
                      </div>
                    </div>

                    {/* Navigation Buttons */}
                    <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between gap-2">
                      <button
                        onClick={() => setActiveVideoIndex((prev) => Math.max(0, prev - 1))}
                        disabled={activeVideoIndex === 0}
                        className="px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-800 hover:bg-slate-800 text-slate-200 text-xs font-semibold disabled:opacity-40 flex items-center gap-1.5 transition-colors"
                      >
                        <ChevronLeft className="w-4 h-4" />
                        <span>Previous Lesson</span>
                      </button>
                      <button
                        onClick={() => setActiveVideoIndex((prev) => Math.min(totalVideos - 1, prev + 1))}
                        disabled={activeVideoIndex >= totalVideos - 1}
                        className="px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-800 hover:bg-slate-800 text-slate-200 text-xs font-semibold disabled:opacity-40 flex items-center gap-1.5 transition-colors"
                      >
                        <span>Next Lesson</span>
                        <ChevronRight className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  {/* AI Summary and Tags */}
                  {(detailData?.summary || course.summary) && (
                    <div className="rounded-2xl bg-surface-950/60 border border-slate-800/80 p-4 space-y-2">
                      <div
                        className="flex items-center justify-between cursor-pointer select-none"
                        onClick={() => setIsAiInsightsExpanded(!isAiInsightsExpanded)}
                      >
                        <div className="flex items-center gap-2 text-brand-300 font-bold text-sm">
                          <Sparkles className="w-4 h-4 text-brand-400" />
                          <span>Course Overview & AI Insights</span>
                        </div>
                        <span className="text-xs text-slate-500">
                          {isAiInsightsExpanded ? 'Hide' : 'Show'}
                        </span>
                      </div>

                      {isAiInsightsExpanded && (
                        <div className="pt-2 space-y-3">
                          <p className="text-xs text-slate-300 leading-relaxed">
                            {detailData?.summary || course.summary}
                          </p>

                          {course.tags && course.tags.length > 0 && (
                            <div className="flex flex-wrap gap-1.5 pt-1">
                              {course.tags.map((tag, idx) => (
                                <span
                                  key={idx}
                                  className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-medium bg-brand-500/15 text-brand-300 border border-brand-500/30"
                                >
                                  <Tag className="w-2.5 h-2.5" />
                                  <span>{tag}</span>
                                </span>
                              ))}
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  )}
                </>
              ) : (
                <div className="text-center py-12 text-slate-400">
                  <p>No video selected.</p>
                </div>
              )}
            </div>

            {/* Sidebar Playlist Column with Clean Sequential Indexing */}
            <div className="lg:col-span-4 flex flex-col bg-surface-950/40 min-h-0">
              <div className="p-4 border-b border-slate-800/80 space-y-3 shrink-0">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <ListVideo className="w-4 h-4 text-slate-400" />
                    <h4 className="text-sm font-bold text-slate-200">Playlist Curriculum</h4>
                  </div>
                  <span className="text-xs font-semibold text-slate-400">
                    {playlist.length} lessons
                  </span>
                </div>

                <div className="relative">
                  <Search className="w-3.5 h-3.5 text-slate-500 absolute left-3 top-2.5" />
                  <input
                    type="text"
                    placeholder="Search lessons..."
                    value={playlistSearch}
                    onChange={(e) => setPlaylistSearch(e.target.value)}
                    className="w-full pl-9 pr-3 py-1.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 placeholder-slate-500 text-xs focus:outline-none focus:border-brand-500 transition-all"
                  />
                </div>
              </div>

              {/* Scrollable Playlist Items */}
              <div className="flex-1 overflow-y-auto p-2 space-y-1.5">
                {filteredPlaylist.map((video, idx) => {
                  const actualIndex = playlist.findIndex((v) => v.video_id === video.video_id);
                  const isCurrent = actualIndex === activeVideoIndex;
                  const itemIndex = video.index || (actualIndex >= 0 ? actualIndex + 1 : idx + 1);

                  return (
                    <div
                      key={video.video_id || idx}
                      onClick={() => setActiveVideoIndex(actualIndex >= 0 ? actualIndex : idx)}
                      className={`group flex items-start gap-3 p-3 rounded-xl cursor-pointer border transition-all ${
                        isCurrent
                          ? 'bg-brand-500/15 border-brand-500/40 text-white shadow-md'
                          : 'bg-slate-900/40 hover:bg-slate-900 border-transparent hover:border-slate-800 text-slate-300'
                      }`}
                    >
                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          handleMarkVideoCompleted(video, actualIndex >= 0 ? actualIndex : idx);
                        }}
                        className={`mt-0.5 shrink-0 w-5 h-5 rounded-md flex items-center justify-center transition-colors ${
                          video.status
                            ? 'bg-emerald-500 text-slate-950'
                            : 'border border-slate-700 hover:border-slate-500 text-transparent hover:text-slate-400'
                        }`}
                        title={video.status ? 'Completed' : 'Click to complete (+10 XP)'}
                      >
                        <Check className="w-3.5 h-3.5 stroke-[3]" />
                      </button>

                      <div className="flex-1 min-w-0">
                        <div className="flex items-center justify-between gap-1 text-[11px] text-slate-500 mb-0.5 font-mono">
                          <span>#{itemIndex}</span>
                          {video.duration && (
                            <span className="text-slate-400">{video.duration}</span>
                          )}
                        </div>
                        <p
                          className={`text-xs font-semibold leading-snug line-clamp-2 ${
                            isCurrent ? 'text-white font-bold' : 'text-slate-200 group-hover:text-white'
                          } ${video.status ? 'line-through text-slate-500' : ''}`}
                        >
                          {video.title}
                        </p>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
