import React from 'react';
import { useAuth } from '../context/AuthContext';
import type { CourseData, BadgeItem } from '../types';
import {
  Trophy,
  Flame,
  Zap,
  Award,
  BookOpen,
  CheckCircle,
  Sparkles,
  Lock,
  Unlock,
} from 'lucide-react';

interface AchievementsViewProps {
  courses: CourseData[];
  onOpenCreateCourse: () => void;
}

export const AchievementsView: React.FC<AchievementsViewProps> = ({
  courses,
  onOpenCreateCourse,
}) => {
  const { stats, user } = useAuth();

  const totalPoints = stats?.total_points || 0;
  const streak = stats?.streak_count || 0;
  const completedCourses = courses.filter((c) => c.completion_percent >= 100).length;

  const currentLevel = Math.max(1, Math.floor(totalPoints / 100) + 1);
  const currentLevelXP = totalPoints % 100;
  const nextLevelXP = 100;
  const progressToNextLevel = (currentLevelXP / nextLevelXP) * 100;

  const tiers = [
    { level: 1, title: 'Novice Explorer', minXP: 0, icon: '🌱' },
    { level: 2, title: 'Code Apprentice', minXP: 100, icon: '⚡' },
    { level: 3, title: 'Dedicated Scholar', minXP: 200, icon: '📖' },
    { level: 4, title: 'Skill Master', minXP: 300, icon: '🎯' },
    { level: 5, title: 'Grandmaster Learner', minXP: 500, icon: '👑' },
    { level: 6, title: 'Legendary Architect', minXP: 1000, icon: '🌟' },
  ];

  const badges: BadgeItem[] = [
    {
      id: 'first_step',
      title: 'First Step',
      description: 'Create an account and access your learning dashboard.',
      icon: '🌱',
      isUnlocked: true,
      currentProgress: 1,
      targetProgress: 1,
      category: 'general',
    },
    {
      id: 'streak_1',
      title: 'Flame Starter',
      description: 'Reach a 1-day learning streak.',
      icon: '🔥',
      isUnlocked: streak >= 1,
      currentProgress: Math.min(streak, 1),
      targetProgress: 1,
      category: 'streak',
    },
    {
      id: 'streak_3',
      title: 'Consistency Champion',
      description: 'Maintain a 3-day active streak.',
      icon: '⚡',
      isUnlocked: streak >= 3,
      currentProgress: Math.min(streak, 3),
      targetProgress: 3,
      category: 'streak',
    },
    {
      id: 'streak_7',
      title: 'Streak Legend',
      description: 'Achieve a 7-day uninterrupted study streak.',
      icon: '🌟',
      isUnlocked: streak >= 7,
      currentProgress: Math.min(streak, 7),
      targetProgress: 7,
      category: 'streak',
    },
    {
      id: 'points_50',
      title: 'Half Century',
      description: 'Earn 50 Total Points from completed video lessons.',
      icon: '🥉',
      isUnlocked: totalPoints >= 50,
      currentProgress: Math.min(totalPoints, 50),
      targetProgress: 50,
      category: 'points',
    },
    {
      id: 'points_100',
      title: 'Century Club',
      description: 'Earn 100 Total Points & reach Level 2.',
      icon: '🥈',
      isUnlocked: totalPoints >= 100,
      currentProgress: Math.min(totalPoints, 100),
      targetProgress: 100,
      category: 'points',
    },
    {
      id: 'points_500',
      title: 'XP High Roller',
      description: 'Earn 500 Total Points across all courses.',
      icon: '💎',
      isUnlocked: totalPoints >= 500,
      currentProgress: Math.min(totalPoints, 500),
      targetProgress: 500,
      category: 'points',
    },
    {
      id: 'course_1',
      title: 'Playlist Pioneer',
      description: 'Import and track your first YouTube playlist.',
      icon: '📚',
      isUnlocked: courses.length >= 1,
      currentProgress: Math.min(courses.length, 1),
      targetProgress: 1,
      category: 'courses',
    },
    {
      id: 'course_finisher',
      title: 'Course Conqueror',
      description: 'Complete 100% of any course playlist.',
      icon: '🎓',
      isUnlocked: completedCourses >= 1,
      currentProgress: Math.min(completedCourses, 1),
      targetProgress: 1,
      category: 'courses',
    },
  ];

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-slate-900 border border-slate-800 p-6 sm:p-7">
        <div className="grid grid-cols-1 md:grid-cols-12 gap-6 items-center">
          <div className="md:col-span-7 space-y-3">
            <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-md text-xs font-semibold bg-slate-800 text-brand-400 border border-slate-700">
              <Trophy className="w-3.5 h-3.5 text-brand-400" />
              <span>Learner Progression Tier</span>
            </div>

            <h2 className="text-2xl font-bold text-white">
              Level {currentLevel}:{' '}
              <span className="text-brand-400">
                {tiers[Math.min(currentLevel - 1, tiers.length - 1)]?.title || 'Master'}
              </span>
            </h2>

            <p className="text-sm text-slate-400 leading-relaxed">
              Every lesson you complete awards strictly +10 XP. Finish entire courses for a +100 XP
              bonus and unlock higher ranks.
            </p>

            <div className="pt-2">
              <div className="flex items-center justify-between text-xs font-medium text-slate-400 mb-1.5">
                <span>{currentLevelXP} XP Earned</span>
                <span>{nextLevelXP - currentLevelXP} XP needed for Level {currentLevel + 1}</span>
              </div>
              <div className="w-full h-2 rounded-full bg-slate-950 border border-slate-800 overflow-hidden">
                <div
                  className="h-full bg-brand-500 rounded-full"
                  style={{ width: `${Math.min(100, Math.max(0, progressToNextLevel))}%` }}
                />
              </div>
            </div>
          </div>

          <div className="md:col-span-5 grid grid-cols-2 gap-3">
            <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 text-center">
              <Flame className="w-5 h-5 text-amber-400 mx-auto" />
              <div className="mt-1 text-2xl font-bold text-white">{streak}</div>
              <p className="text-xs text-slate-400">Day Streak</p>
            </div>
            <div className="p-4 rounded-xl bg-slate-950 border border-slate-800 text-center">
              <Zap className="w-6 h-6 text-brand-400 mx-auto" />
              <div className="mt-2 text-2xl font-extrabold text-white">{totalPoints}</div>
              <p className="text-[11px] text-slate-400 font-medium">Total XP</p>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-xl font-bold text-white tracking-tight">Badges & Milestones</h3>
            <p className="text-xs text-slate-400">
              Unlock prestigious achievements by maintaining habits and completing playlists.
            </p>
          </div>
          <span className="text-xs font-bold text-slate-400">
            {badges.filter((b) => b.isUnlocked).length} / {badges.length} Unlocked
          </span>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {badges.map((badge) => {
            const isDone = badge.isUnlocked;
            const percent = Math.round((badge.currentProgress / badge.targetProgress) * 100);

            return (
              <div
                key={badge.id}
                className={`p-4 rounded-xl border flex flex-col justify-between ${
                  isDone
                    ? 'bg-slate-900 border-amber-500/30'
                    : 'bg-slate-950 border-slate-800 opacity-60'
                }`}
              >
                <div>
                  <div className="flex items-start justify-between gap-3">
                    <div
                      className={`w-10 h-10 rounded-lg flex items-center justify-center text-xl ${
                        isDone ? 'bg-amber-500/20 border border-amber-500/30' : 'bg-slate-800'
                      }`}
                    >
                      {badge.icon}
                    </div>
                    {isDone ? (
                      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-semibold bg-amber-500/20 text-amber-300 border border-amber-500/30">
                        <Unlock className="w-3 h-3" />
                        <span>Unlocked</span>
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded text-[10px] font-semibold bg-slate-800 text-slate-400 border border-slate-700">
                        <Lock className="w-3 h-3" />
                        <span>Locked</span>
                      </span>
                    )}
                  </div>

                  <h4 className="text-sm font-bold text-white mt-3">{badge.title}</h4>
                  <p className="text-xs text-slate-400 mt-1 leading-relaxed">
                    {badge.description}
                  </p>
                </div>

                <div className="mt-4 pt-3 border-t border-slate-800">
                  <div className="flex items-center justify-between text-[11px] font-medium text-slate-400 mb-1">
                    <span>Progress</span>
                    <span>
                      {badge.currentProgress} / {badge.targetProgress}
                    </span>
                  </div>
                  <div className="w-full h-1.5 rounded-full bg-slate-800 overflow-hidden">
                    <div
                      className={`h-full rounded-full ${
                        isDone ? 'bg-amber-400' : 'bg-slate-600'
                      }`}
                      style={{ width: `${Math.min(100, percent)}%` }}
                    />
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
};
