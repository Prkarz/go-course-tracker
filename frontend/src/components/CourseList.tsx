import React, { useState, useMemo } from 'react';
import type { CourseData } from '../types';
import { CourseCard } from './CourseCard';
import {
  Search,
  Filter,
  PlusCircle,
  BookOpen,
  Sparkles,
  SlidersHorizontal,
  X,
  ArrowUpDown,
} from 'lucide-react';

interface CourseListProps {
  courses: CourseData[];
  onStartCourse: (courseId: number) => void;
  onOpenPlayer: (course: CourseData) => void;
  onDeleteCourse: (courseId: number) => void;
  onOpenCreateCourse: () => void;
  startingCourseId: number | null;
}

type StatusFilter = 'all' | 'in_progress' | 'completed' | 'not_started';
type SortOption = 'recent' | 'progress_desc' | 'progress_asc' | 'title_asc';

export const CourseList: React.FC<CourseListProps> = ({
  courses,
  onStartCourse,
  onOpenPlayer,
  onDeleteCourse,
  onOpenCreateCourse,
  startingCourseId,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('all');
  const [selectedTag, setSelectedTag] = useState<string>('all');
  const [sortBy, setSortBy] = useState<SortOption>('recent');

  const allTags = useMemo(() => {
    const set = new Set<string>();
    courses.forEach((c) => {
      if (Array.isArray(c.tags)) {
        c.tags.forEach((t) => set.add(t));
      }
    });
    return Array.from(set);
  }, [courses]);

  const filteredCourses = useMemo(() => {
    return courses
      .filter((course) => {
        if (searchQuery.trim()) {
          const q = searchQuery.toLowerCase();
          const matchesTitle = (course.title || '').toLowerCase().includes(q);
          const matchesSummary = (course.summary || '').toLowerCase().includes(q);
          const matchesTag = Array.isArray(course.tags) && course.tags.some((t) => t.toLowerCase().includes(q));
          if (!matchesTitle && !matchesSummary && !matchesTag) return false;
        }

        if (statusFilter === 'in_progress') {
          if (!course.is_started || course.completion_percent >= 100) return false;
        } else if (statusFilter === 'completed') {
          if (course.completion_percent < 100) return false;
        } else if (statusFilter === 'not_started') {
          if (course.is_started) return false;
        }

        if (selectedTag !== 'all') {
          if (!Array.isArray(course.tags) || !course.tags.includes(selectedTag)) return false;
        }

        return true;
      })
      .sort((a, b) => {
        if (sortBy === 'progress_desc') {
          return b.completion_percent - a.completion_percent;
        }
        if (sortBy === 'progress_asc') {
          return a.completion_percent - b.completion_percent;
        }
        if (sortBy === 'title_asc') {
          return (a.title || '').localeCompare(b.title || '');
        }
        return b.id - a.id;
      });
  }, [courses, searchQuery, statusFilter, selectedTag, sortBy]);

  return (
    <div className="space-y-6">
      <div className="rounded-2xl bg-surface-900 border border-slate-800/80 p-4 sm:p-5 shadow-lg space-y-4">
        <div className="flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="relative w-full md:w-80">
            <Search className="w-4 h-4 text-slate-400 absolute left-3.5 top-3.5" />
            <input
              type="text"
              placeholder="Search by title, topics, or tags..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-9 py-2.5 rounded-xl bg-slate-950 border border-slate-800 text-slate-100 placeholder-slate-500 text-sm focus:outline-none focus:border-accent-500 transition-all"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 top-3 text-slate-400 hover:text-slate-200"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>

          <div className="flex items-center gap-1 p-1 rounded-xl bg-slate-950 border border-slate-800 text-xs font-semibold overflow-x-auto max-w-full">
            <button
              onClick={() => setStatusFilter('all')}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                statusFilter === 'all'
                  ? 'bg-slate-800 text-white shadow-sm'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              All ({courses.length})
            </button>
            <button
              onClick={() => setStatusFilter('in_progress')}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                statusFilter === 'in_progress'
                  ? 'bg-accent-600/30 text-accent-300 border border-accent-500/30'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              In Progress ({courses.filter((c) => c.is_started && c.completion_percent < 100).length})
            </button>
            <button
              onClick={() => setStatusFilter('completed')}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                statusFilter === 'completed'
                  ? 'bg-emerald-500/20 text-emerald-300 border border-emerald-500/30'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              Completed ({courses.filter((c) => c.completion_percent >= 100).length})
            </button>
            <button
              onClick={() => setStatusFilter('not_started')}
              className={`px-3 py-1.5 rounded-lg transition-all ${
                statusFilter === 'not_started'
                  ? 'bg-slate-800 text-slate-200'
                  : 'text-slate-400 hover:text-slate-200'
              }`}
            >
              Not Started ({courses.filter((c) => !c.is_started).length})
            </button>
          </div>

          <div className="flex items-center gap-2 w-full md:w-auto justify-end">
            <div className="flex items-center gap-1.5 text-xs text-slate-400">
              <ArrowUpDown className="w-3.5 h-3.5" />
              <span>Sort:</span>
            </div>
            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value as SortOption)}
              className="px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-slate-200 text-xs focus:outline-none focus:border-accent-500"
            >
              <option value="recent">Recently Added</option>
              <option value="progress_desc">Progress: High to Low</option>
              <option value="progress_asc">Progress: Low to High</option>
              <option value="title_asc">Title: A to Z</option>
            </select>
          </div>
        </div>

        {allTags.length > 0 && (
          <div className="pt-3 border-t border-slate-800/80 flex items-center gap-2 overflow-x-auto pb-1 text-xs">
            <span className="text-slate-400 shrink-0 font-medium">Topic Tags:</span>
            <button
              onClick={() => setSelectedTag('all')}
              className={`px-2.5 py-1 rounded-full text-xs font-semibold transition-all ${
                selectedTag === 'all'
                  ? 'bg-brand-500/20 text-brand-400 border border-brand-500/40'
                  : 'bg-slate-950 text-slate-400 border border-slate-800 hover:text-slate-200'
              }`}
            >
              All Topics
            </button>
            {allTags.map((tag) => (
              <button
                key={tag}
                onClick={() => setSelectedTag(tag === selectedTag ? 'all' : tag)}
                className={`px-2.5 py-1 rounded-full text-xs font-semibold shrink-0 transition-all ${
                  selectedTag === tag
                    ? 'bg-accent-600/30 text-accent-300 border border-accent-500/40'
                    : 'bg-slate-950 text-slate-400 border border-slate-800 hover:text-slate-200'
                }`}
              >
                #{tag}
              </button>
            ))}
          </div>
        )}
      </div>

      {filteredCourses.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredCourses.map((course) => (
            <CourseCard
              key={course.id}
              course={course}
              onStartCourse={onStartCourse}
              onOpenPlayer={onOpenPlayer}
              onDeleteCourse={onDeleteCourse}
              isStarting={startingCourseId === course.id}
            />
          ))}
        </div>
      ) : (
        <div className="rounded-3xl bg-surface-900 border border-slate-800 p-12 text-center max-w-lg mx-auto space-y-4">
          <div className="w-16 h-16 rounded-2xl bg-slate-800/80 flex items-center justify-center mx-auto text-slate-400">
            <BookOpen className="w-8 h-8" />
          </div>
          {courses.length === 0 ? (
            <>
              <h3 className="text-lg font-bold text-white">No courses added yet</h3>
              <p className="text-xs text-slate-400 max-w-sm mx-auto leading-relaxed">
                Add any YouTube learning playlist to unlock AI summaries, track completion, and earn streak XP.
              </p>
              <button
                onClick={onOpenCreateCourse}
                className="mt-2 inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-brand-500 hover:bg-brand-400 text-slate-950 font-bold text-sm shadow-lg shadow-brand-500/20 transition-all"
              >
                <PlusCircle className="w-4 h-4" />
                <span>Import First YouTube Playlist</span>
              </button>
            </>
          ) : (
            <>
              <h3 className="text-lg font-bold text-white">No matching courses found</h3>
              <p className="text-xs text-slate-400">
                Try clearing your search query or filter tags to see all courses.
              </p>
              <button
                onClick={() => {
                  setSearchQuery('');
                  setStatusFilter('all');
                  setSelectedTag('all');
                }}
                className="mt-2 px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold transition-colors"
              >
                Clear Filters
              </button>
            </>
          )}
        </div>
      )}
    </div>
  );
};
