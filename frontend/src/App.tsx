import React, { useState, useEffect, useCallback } from 'react';
import { useAuth } from './context/AuthContext';
import { useNotification } from './context/NotificationContext';
import type { CourseData, ActiveTab } from './types';
import { api } from './services/api';
import { Navbar } from './components/Navbar';
import { StatsHero } from './components/StatsHero';
import { CourseList } from './components/CourseList';
import { CreateCourseModal } from './components/CreateCourseModal';
import { CourseDetailModal } from './components/CourseDetailModal';
import { AchievementsView } from './components/AchievementsView';
import { BackendSettingsModal } from './components/BackendSettingsModal';
import { DeleteAccountModal } from './components/DeleteAccountModal';
import { AuthModal } from './components/AuthModal';
import { ToastContainer } from './components/ToastContainer';
import {
  Zap,
  Sparkles,
  BookOpen,
  Trophy,
  Flame,
  PlusCircle,
  PlayCircle,
  Video,
  ArrowRight,
  ShieldCheck,
  Brain,
  CheckCircle2,
} from 'lucide-react';

export function App() {
  const { isAuthenticated, isLoading, openAuthModal, stats, refreshStats } = useAuth();
  const { addToast } = useNotification();

  const [activeTab, setActiveTab] = useState<ActiveTab>('dashboard');
  const [courses, setCourses] = useState<CourseData[]>([]);
  const [isCoursesLoading, setIsCoursesLoading] = useState(false);

  // Active Course for Player & Modals
  const [selectedCourse, setSelectedCourse] = useState<CourseData | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isSettingsModalOpen, setIsSettingsModalOpen] = useState(false);
  const [isDeleteAccountModalOpen, setIsDeleteAccountModalOpen] = useState(false);
  const [startingCourseId, setStartingCourseId] = useState<number | null>(null);

  // Fetch Courses
  const loadCourses = useCallback(async (preferredSelectCourseId?: number) => {
    if (!isAuthenticated) {
      setCourses([]);
      return;
    }
    setIsCoursesLoading(true);
    try {
      const data = await api.getMyCourses();
      setCourses(data);
      if (preferredSelectCourseId) {
        const found = data.find((c) => c.id === preferredSelectCourseId);
        if (found) {
          setSelectedCourse(found);
        }
      } else if (data.length > 0 && !selectedCourse) {
        setSelectedCourse(data[0]);
      }
    } catch (err: unknown) {
      addToast(
        'Failed to load courses',
        err instanceof Error ? err.message : 'Please check your connection',
        'error'
      );
    } finally {
      setIsCoursesLoading(false);
    }
  }, [isAuthenticated, addToast, selectedCourse]);

  useEffect(() => {
    if (isAuthenticated) {
      loadCourses();
      refreshStats();
    }
  }, [isAuthenticated, loadCourses, refreshStats]);

  // Start Course Action
  const handleStartCourse = async (courseId: number) => {
    setStartingCourseId(courseId);
    try {
      const res = await api.startCourse(courseId);
      addToast('Course Started! 🚀', res.message || 'You are now enrolled in this course.', 'success');
      await loadCourses();
      await refreshStats();

      const course = courses.find((c) => c.id === courseId);
      if (course) {
        setSelectedCourse(course);
        setIsDetailModalOpen(true);
      }
    } catch (err: unknown) {
      addToast('Failed to start course', err instanceof Error ? err.message : 'Error', 'error');
    } finally {
      setStartingCourseId(null);
    }
  };

  // Delete Course Action (Deletion Immunity: earned XP is permanent)
  const handleDeleteCourse = async (courseId: number) => {
    if (!window.confirm('Are you sure you want to delete this course from your dashboard? Your earned XP is permanent and will NOT be reduced.')) {
      return;
    }
    try {
      const res = await api.deleteCourse(courseId);
      addToast('Course Removed', res.message || 'Course removed. Your accumulated XP is preserved.', 'info');
      await loadCourses();
      await refreshStats();
      if (selectedCourse?.id === courseId) {
        setSelectedCourse(null);
        setIsDetailModalOpen(false);
      }
    } catch (err: unknown) {
      addToast('Delete failed', err instanceof Error ? err.message : 'Error', 'error');
    }
  };

  // Open Player Modal for Course
  const handleOpenPlayer = (course: CourseData) => {
    setSelectedCourse(course);
    setIsDetailModalOpen(true);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-surface-950 flex flex-col items-center justify-center space-y-4">
        <div className="w-12 h-12 border-4 border-brand-500/20 border-t-brand-500 rounded-full animate-spin" />
        <p className="text-sm font-bold text-slate-300">Initializing Course Flow...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-surface-950 text-slate-100 flex flex-col selection:bg-accent-500 selection:text-white">
      {/* Floating Notifications */}
      <ToastContainer />

      {/* Navigation Bar */}
      <Navbar
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        onOpenCreateCourse={() => {
          if (!isAuthenticated) {
            openAuthModal('login');
          } else {
            setIsCreateModalOpen(true);
          }
        }}
        onOpenSettings={() => setIsSettingsModalOpen(true)}
        onOpenDeleteAccount={() => setIsDeleteAccountModalOpen(true)}
        hasActiveCourse={!!selectedCourse}
      />

      {/* Main Content Area */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8 space-y-8">
        {!isAuthenticated ? (
          /* Guest / Landing View */
          <div className="py-12 space-y-12">
            {/* Hero Section */}
            <div className="rounded-2xl bg-slate-900 border border-slate-800 p-8 sm:p-12 text-center max-w-4xl mx-auto overflow-hidden">
              <div className="space-y-5">
                <div className="inline-flex items-center gap-2 px-3 py-1 rounded-md text-xs font-semibold bg-slate-800 text-brand-400 border border-slate-700">
                  <Sparkles className="w-4 h-4 text-brand-400" />
                  <span>Go Backend + Gemini AI Learning Platform</span>
                </div>

                <h1 className="text-3xl sm:text-5xl font-bold text-white tracking-tight leading-tight">
                  Supercharge your YouTube learning with{' '}
                  <span className="text-brand-400">
                    AI Insights & Progress Tracking
                  </span>
                </h1>

                <p className="text-base text-slate-400 max-w-2xl mx-auto leading-relaxed">
                  Turn any YouTube playlist into a structured course. Track video lessons, build
                  daily study streaks, gain XP, and receive automatic AI summaries.
                </p>

                <div className="flex flex-wrap items-center justify-center gap-3 pt-2">
                  <button
                    onClick={() => openAuthModal('register')}
                    className="px-6 py-3 rounded-lg bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-sm transition-colors flex items-center gap-2"
                  >
                    <span>Get Started Free</span>
                    <ArrowRight className="w-4 h-4" />
                  </button>

                  <button
                    onClick={() => openAuthModal('login')}
                    className="px-6 py-3 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold text-sm border border-slate-700 transition-colors"
                  >
                    <span>Sign In to Account</span>
                  </button>
                </div>
              </div>
            </div>

            {/* Features Highlight Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
              <div className="p-6 rounded-2xl bg-surface-900 border border-slate-800/80 space-y-3">
                <div className="w-12 h-12 rounded-xl bg-brand-500/10 border border-brand-500/30 flex items-center justify-center text-brand-400">
                  <Brain className="w-6 h-6" />
                </div>
                <h3 className="text-base font-bold text-white">Google Gemini AI Insights</h3>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Automatically analyzes YouTube playlists to generate concise course summaries and
                  smart curriculum topic tags.
                </p>
              </div>

              <div className="p-6 rounded-2xl bg-surface-900 border border-slate-800/80 space-y-3">
                <div className="w-12 h-12 rounded-xl bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400">
                  <Flame className="w-6 h-6" />
                </div>
                <h3 className="text-base font-bold text-white">Streaks & XP Gamification</h3>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Earn +10 XP per video and +100 XP per completed course. Level up from Novice to
                  Grandmaster and keep your streak alive.
                </p>
              </div>

              <div className="p-6 rounded-2xl bg-surface-900 border border-slate-800/80 space-y-3">
                <div className="w-12 h-12 rounded-xl bg-accent-500/10 border border-accent-500/30 flex items-center justify-center text-accent-400">
                  <Video className="w-6 h-6" />
                </div>
                <h3 className="text-base font-bold text-white">Distraction-Free Video Studio</h3>
                <p className="text-xs text-slate-400 leading-relaxed">
                  Watch course playlists inside an organized environment with duration trackers,
                  completion checkmarks, and video search.
                </p>
              </div>
            </div>
          </div>
        ) : (
          /* Authenticated Dashboard Tabs */
          <>
            {/* TAB: DASHBOARD */}
            {activeTab === 'dashboard' && (
              <div className="space-y-8 animate-in fade-in duration-200">
                {/* Stats Hero Section */}
                <StatsHero
                  courses={courses}
                  onOpenCreateCourse={() => setIsCreateModalOpen(true)}
                  onResumeCourse={handleOpenPlayer}
                  onViewAchievements={() => setActiveTab('achievements')}
                />

                {/* Courses Section */}
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h2 className="text-xl font-bold text-white tracking-tight">Your Courses</h2>
                      <p className="text-xs text-slate-400">
                        Manage your active learning playlists and track your progress.
                      </p>
                    </div>

                    <button
                      onClick={() => setIsCreateModalOpen(true)}
                      className="px-4 py-2 rounded-xl bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-xs shadow-md shadow-brand-500/20 transition-all flex items-center gap-1.5"
                    >
                      <PlusCircle className="w-4 h-4" />
                      <span>Add Course</span>
                    </button>
                  </div>

                  <CourseList
                    courses={courses}
                    onStartCourse={handleStartCourse}
                    onOpenPlayer={handleOpenPlayer}
                    onDeleteCourse={handleDeleteCourse}
                    onOpenCreateCourse={() => setIsCreateModalOpen(true)}
                    startingCourseId={startingCourseId}
                  />
                </div>
              </div>
            )}

            {/* TAB: COURSES */}
            {activeTab === 'courses' && (
              <div className="space-y-6 animate-in fade-in duration-200">
                <div className="flex items-center justify-between">
                  <div>
                    <h2 className="text-2xl font-extrabold text-white tracking-tight">
                      My Learning Playlists
                    </h2>
                    <p className="text-xs text-slate-400 mt-0.5">
                      Explore, filter, and track all your enrolled YouTube courses.
                    </p>
                  </div>

                  <button
                    onClick={() => setIsCreateModalOpen(true)}
                    className="px-4 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-sm shadow-md shadow-brand-500/20 transition-all flex items-center gap-2"
                  >
                    <PlusCircle className="w-4 h-4" />
                    <span>Import YouTube Course</span>
                  </button>
                </div>

                <CourseList
                  courses={courses}
                  onStartCourse={handleStartCourse}
                  onOpenPlayer={handleOpenPlayer}
                  onDeleteCourse={handleDeleteCourse}
                  onOpenCreateCourse={() => setIsCreateModalOpen(true)}
                  startingCourseId={startingCourseId}
                />
              </div>
            )}

            {/* TAB: PLAYER */}
            {activeTab === 'player' && (
              <div className="space-y-6 animate-in fade-in duration-200">
                {selectedCourse ? (
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <span className="text-xs font-semibold text-brand-400 uppercase tracking-wider">
                          Active Course Player
                        </span>
                        <h2 className="text-xl sm:text-2xl font-extrabold text-white">
                          {selectedCourse.title}
                        </h2>
                      </div>
                      <button
                        onClick={() => setIsDetailModalOpen(true)}
                        className="px-4 py-2 rounded-xl bg-accent-600 hover:bg-accent-500 text-white font-bold text-xs shadow-md"
                      >
                        Open Fullscreen Studio
                      </button>
                    </div>

                    <div className="p-8 rounded-3xl bg-surface-900 border border-slate-800 text-center space-y-4">
                      <div className="w-16 h-16 rounded-2xl bg-accent-600/20 text-accent-400 flex items-center justify-center mx-auto">
                        <Video className="w-8 h-8" />
                      </div>
                      <h3 className="text-lg font-bold text-white">Course Player Ready</h3>
                      <p className="text-xs text-slate-400 max-w-md mx-auto">
                        Open the interactive YouTube video studio with playlist curriculum, AI
                        insights drawer, and video completion tracker.
                      </p>
                      <button
                        onClick={() => setIsDetailModalOpen(true)}
                        className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-gradient-to-r from-brand-500 to-brand-400 text-slate-950 font-bold text-sm shadow-lg shadow-brand-500/20"
                      >
                        <PlayCircle className="w-5 h-5" />
                        <span>Launch Course Video Player</span>
                      </button>
                    </div>
                  </div>
                ) : (
                  <div className="p-12 rounded-3xl bg-surface-900 border border-slate-800 text-center max-w-md mx-auto space-y-4">
                    <BookOpen className="w-12 h-12 text-slate-500 mx-auto" />
                    <h3 className="text-lg font-bold text-white">No course selected</h3>
                    <p className="text-xs text-slate-400">
                      Pick a course from your dashboard or import a new YouTube playlist to start watching.
                    </p>
                    <button
                      onClick={() => setActiveTab('courses')}
                      className="px-5 py-2.5 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold"
                    >
                      Browse Courses
                    </button>
                  </div>
                )}
              </div>
            )}

            {/* TAB: ACHIEVEMENTS */}
            {activeTab === 'achievements' && (
              <AchievementsView
                courses={courses}
                onOpenCreateCourse={() => setIsCreateModalOpen(true)}
              />
            )}
          </>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-800/80 bg-surface-950 py-6 mt-12 text-center text-xs text-slate-500">
        <div className="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-3">
          <div className="flex items-center gap-2 font-medium">
            <span className="w-2 h-2 rounded-full bg-brand-500" />
            <span className="text-slate-400">Course Flow — AI-Driven Learning Platform</span>
          </div>
          <div className="flex items-center gap-4 text-slate-400">
            <button onClick={() => setIsSettingsModalOpen(true)} className="hover:text-slate-200 transition-colors">
              Backend Settings
            </button>
            <span>•</span>
            <a
              href="https://github.com/Prkarz/course-tracker"
              target="_blank"
              rel="noreferrer"
              className="hover:text-slate-200 transition-colors"
            >
              GitHub Repository
            </a>
          </div>
        </div>
      </footer>

      {/* Modals Container */}
      <AuthModal />

      <CreateCourseModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCourseCreated={async (newCourseId?: number) => {
          await loadCourses(newCourseId);
          await refreshStats();
          if (newCourseId) {
            setActiveTab('courses');
          }
        }}
      />

      <CourseDetailModal
        isOpen={isDetailModalOpen}
        course={selectedCourse}
        onClose={() => setIsDetailModalOpen(false)}
        onCourseUpdated={() => {
          loadCourses();
          refreshStats();
        }}
      />

      <BackendSettingsModal
        isOpen={isSettingsModalOpen}
        onClose={() => setIsSettingsModalOpen(false)}
      />

      <DeleteAccountModal
        isOpen={isDeleteAccountModalOpen}
        onClose={() => setIsDeleteAccountModalOpen(false)}
      />
    </div>
  );
}

export default App;
