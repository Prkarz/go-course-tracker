import React from 'react';
import { useAuth } from '../context/AuthContext';
import { Flame, Zap, Award, BookOpen, Calendar, ArrowUpRight, PlayCircle } from 'lucide-react';
import type { CourseData } from '../types';

interface StatsHeroProps {
  courses: CourseData[];
  onOpenCreateCourse: () => void;
  onResumeCourse: (course: CourseData) => void;
  onViewAchievements: () => void;
}

export const StatsHero: React.FC<StatsHeroProps> = ({
  courses,
  onOpenCreateCourse,
  onResumeCourse,
  onViewAchievements,
}) => {
  const { user, stats } = useAuth();

  const totalPoints = stats?.total_points || 0;
  const streak = stats?.streak_count || 0;

  // Level calculations: every 100 XP is a level
  const currentLevel = Math.max(1, Math.floor(totalPoints / 100) + 1);
  const currentLevelXP = totalPoints % 100;
  const nextLevelXP = 100;
  const levelProgress = (currentLevelXP / nextLevelXP) * 100;

  const levelTitles: Record<number, string> = {
    1: 'Novice Explorer',
    2: 'Code Apprentice',
    3: 'Dedicated Scholar',
    4: 'Skill Master',
    5: 'Grandmaster Learner',
  };
  const levelTitle = levelTitles[Math.min(currentLevel, 5)] || 'Legendary Scholar';

  const activeCourse = courses.find((c) => c.is_started && c.completion_percent < 100) || courses[0];
  const completedCoursesCount = courses.filter((c) => c.completion_percent >= 100).length;

  const formattedLastActive = stats?.last_active_date
    ? new Date(stats.last_active_date).toLocaleDateString(undefined, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
      })
    : 'Today';

  return (
    <div className="rounded-2xl bg-slate-900 border border-slate-800 p-6 sm:p-7">
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
        {/* Left Column: Greeting & Quick Action */}
        <div className="lg:col-span-6 space-y-4">
          <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-md text-xs font-semibold bg-slate-800 text-brand-400 border border-slate-700">
            <span className="w-1.5 h-1.5 rounded-full bg-brand-400" />
            <span>AI-Driven Learning Platform</span>
          </div>

          <h1 className="text-2xl sm:text-3xl font-bold text-white tracking-tight leading-tight">
            Keep crushing your goals,{' '}
            <span className="text-brand-400">
              {user?.email?.split('@')[0] || 'Learner'}!
            </span>
          </h1>

          <p className="text-sm text-slate-400 max-w-xl leading-relaxed">
            Track YouTube playlists lesson-by-lesson, earn XP, maintain study streaks, and unlock
            Gemini AI summaries for each course.
          </p>

          {/* Quick Buttons */}
          <div className="flex flex-wrap items-center gap-3 pt-1">
            <button
              onClick={onOpenCreateCourse}
              className="px-4 py-2 rounded-lg bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-sm transition-colors flex items-center gap-2"
            >
              <Zap className="w-4 h-4 fill-slate-950" />
              <span>Import Course</span>
            </button>

            {activeCourse && (
              <button
                onClick={() => onResumeCourse(activeCourse)}
                className="px-4 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold text-sm border border-slate-700 transition-colors flex items-center gap-2"
              >
                <PlayCircle className="w-4 h-4 text-brand-400" />
                <span className="truncate max-w-[160px]">Resume: {activeCourse.title || 'Course'}</span>
              </button>
            )}
          </div>
        </div>

        {/* Right Column: Stat Cards Grid */}
        <div className="lg:col-span-6 grid grid-cols-2 gap-3">
          {/* Card 1: Streak */}
          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-400">Daily Streak</span>
              <Flame className="w-4 h-4 text-amber-400 fill-amber-400" />
            </div>
            <div className="mt-2 flex items-baseline gap-1.5">
              <span className="text-2xl font-bold text-white">{streak}</span>
              <span className="text-xs text-slate-400">days</span>
            </div>
            <p className="mt-1 text-[11px] text-slate-400">
              {streak > 0 ? '🔥 On fire! Keep learning today.' : 'Start your streak today!'}
            </p>
          </div>

          {/* Card 2: Total XP & Level */}
          <div
            onClick={onViewAchievements}
            className="p-4 rounded-xl bg-slate-950 border border-slate-800 hover:border-slate-700 transition-colors cursor-pointer"
          >
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-400">Level {currentLevel}</span>
              <Award className="w-4 h-4 text-brand-400" />
            </div>
            <div className="mt-2 flex items-baseline gap-1.5">
              <span className="text-2xl font-bold text-brand-400">{totalPoints}</span>
              <span className="text-xs text-slate-400">XP</span>
            </div>
            <div className="mt-2">
              <div className="w-full h-1.5 bg-slate-800 rounded-full overflow-hidden">
                <div
                  className="h-full bg-brand-400 rounded-full"
                  style={{ width: `${Math.min(100, Math.max(0, levelProgress))}%` }}
                />
              </div>
              <p className="mt-1 text-[10px] text-slate-400 truncate">
                {currentLevelXP}/{nextLevelXP} XP to Level {currentLevel + 1}
              </p>
            </div>
          </div>

          {/* Card 3: Enrolled Courses */}
          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-400">Courses Enrolled</span>
              <BookOpen className="w-4 h-4 text-slate-400" />
            </div>
            <div className="mt-2 flex items-baseline gap-1.5">
              <span className="text-2xl font-bold text-white">
                {courses.length}
              </span>
              <span className="text-xs text-slate-400">courses</span>
            </div>
            <p className="mt-1 text-[11px] text-slate-400">
              {completedCoursesCount} completed (100%)
            </p>
          </div>

          {/* Card 4: Last Active */}
          <div className="p-4 rounded-xl bg-slate-950 border border-slate-800">
            <div className="flex items-center justify-between">
              <span className="text-xs font-semibold text-slate-400">Rank Status</span>
              <Calendar className="w-4 h-4 text-slate-400" />
            </div>
            <div className="mt-2">
              <span className="text-base font-bold text-white block truncate">
                {levelTitle}
              </span>
            </div>
            <p className="mt-1 text-[11px] text-slate-400">
              Active: {formattedLastActive}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
