import React, { useState, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import type { ActiveTab } from '../types';
import { api } from '../services/api';
import {
  Flame,
  Zap,
  LayoutDashboard,
  BookOpen,
  Trophy,
  Video,
  Settings,
  LogOut,
  User,
  PlusCircle,
  Menu,
  X,
  Server,
  Trash2,
} from 'lucide-react';

interface NavbarProps {
  activeTab: ActiveTab;
  setActiveTab: (tab: ActiveTab) => void;
  onOpenCreateCourse: () => void;
  onOpenSettings: () => void;
  onOpenDeleteAccount: () => void;
  hasActiveCourse: boolean;
}

export const Navbar: React.FC<NavbarProps> = ({
  activeTab,
  setActiveTab,
  onOpenCreateCourse,
  onOpenSettings,
  onOpenDeleteAccount,
  hasActiveCourse,
}) => {
  const { user, stats, isAuthenticated, logout, openAuthModal } = useAuth();
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [backendHealth, setBackendHealth] = useState<{ ok: boolean; latency: number }>({
    ok: true,
    latency: 0,
  });

  useEffect(() => {
    let isMounted = true;
    const check = async () => {
      const result = await api.checkBackendHealth();
      if (isMounted) {
        setBackendHealth({ ok: result.ok, latency: result.latency });
      }
    };
    check();
    const interval = setInterval(check, 15000);
    return () => {
      isMounted = false;
      clearInterval(interval);
    };
  }, []);

  const navItems = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'courses', label: 'My Courses', icon: BookOpen },
    {
      id: 'player',
      label: 'Player',
      icon: Video,
      badge: hasActiveCourse ? 'Active' : undefined,
    },
    { id: 'achievements', label: 'Achievements', icon: Trophy },
  ];

  return (
    <header className="sticky top-0 z-40 w-full border-b border-slate-800 bg-slate-950">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <button
              onClick={() => setActiveTab('dashboard')}
              className="flex items-center gap-2.5 focus:outline-none"
            >
              <div className="w-9 h-9 rounded-lg bg-brand-500 flex items-center justify-center">
                <Zap className="w-5 h-5 text-slate-950 fill-current" />
              </div>
              <div className="text-left">
                <span className="font-bold text-lg text-white tracking-tight">
                  Course Flow
                </span>
                <span className="block text-[10px] uppercase font-semibold tracking-wider text-brand-400 -mt-1">
                  AI Learning Platform
                </span>
              </div>
            </button>

            {/* Backend connection pill */}
            <button
              onClick={onOpenSettings}
              className="hidden lg:flex items-center gap-1.5 ml-3 px-2.5 py-1 rounded-md text-xs font-medium bg-slate-900 border border-slate-800 text-slate-300 hover:border-slate-700 transition-colors"
              title="Click to configure backend API URL"
            >
              <span
                className={`w-2 h-2 rounded-full ${
                  backendHealth.ok ? 'bg-emerald-400' : 'bg-rose-500'
                }`}
              />
              <span>Go API: {backendHealth.ok ? 'Online' : 'Offline'}</span>
              {backendHealth.ok && backendHealth.latency > 0 && (
                <span className="text-[10px] text-slate-500">({backendHealth.latency}ms)</span>
              )}
            </button>
          </div>

          {/* Navigation Links */}
          <nav className="hidden md:flex items-center gap-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = activeTab === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => setActiveTab(item.id as ActiveTab)}
                  className={`flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-slate-800 text-white border border-slate-700'
                      : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900'
                  }`}
                >
                  <Icon className={`w-4 h-4 ${isActive ? 'text-brand-400' : 'text-slate-400'}`} />
                  <span>{item.label}</span>
                  {item.badge && (
                    <span className="px-1.5 py-0.2 rounded text-[10px] font-bold bg-brand-500/20 text-brand-400 border border-brand-500/30">
                      {item.badge}
                    </span>
                  )}
                </button>
              );
            })}
          </nav>

          {/* Right Header Actions */}
          <div className="flex items-center gap-3">
            {isAuthenticated ? (
              <>
                {/* Streak Counter */}
                <div
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-gradient-to-r from-orange-500/10 to-amber-500/10 border border-amber-500/20 text-amber-300 cursor-pointer hover:border-amber-500/40 transition-colors"
                  onClick={() => setActiveTab('achievements')}
                  title="Daily Study Streak"
                >
                  <Flame className="w-4 h-4 text-orange-400 fill-orange-400/30 animate-pulse" />
                  <span className="text-sm font-bold tracking-tight">
                    {stats?.streak_count || 0}
                  </span>
                  <span className="text-xs text-amber-400/70 hidden sm:inline">days</span>
                </div>

                {/* Total Points */}
                <div
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-gradient-to-r from-brand-500/10 to-emerald-500/10 border border-brand-500/20 text-brand-300 cursor-pointer hover:border-brand-500/40 transition-colors"
                  onClick={() => setActiveTab('achievements')}
                  title="Total Learning XP / Points"
                >
                  <Zap className="w-4 h-4 text-brand-400 fill-brand-400/30" />
                  <span className="text-sm font-bold tracking-tight">
                    {stats?.total_points || 0}
                  </span>
                  <span className="text-xs text-brand-400/70 hidden sm:inline">XP</span>
                </div>

                {/* Add Course Button */}
                <button
                  onClick={onOpenCreateCourse}
                  className="hidden sm:flex items-center gap-2 px-3.5 py-2 rounded-xl bg-gradient-to-r from-brand-600 to-brand-500 hover:from-brand-500 hover:to-brand-400 text-slate-950 font-bold text-sm shadow-md shadow-brand-500/20 hover:shadow-brand-500/30 transition-all transform active:scale-95"
                >
                  <PlusCircle className="w-4 h-4 text-slate-950" />
                  <span>Add Course</span>
                </button>

                {/* Profile Avatar / Dropdown */}
                <div className="relative">
                  <button
                    onClick={() => setIsProfileMenuOpen(!isProfileMenuOpen)}
                    className="flex items-center gap-2 p-1.5 rounded-xl border border-slate-700 bg-slate-900/80 hover:border-slate-600 transition-colors"
                  >
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-tr from-accent-500 to-brand-400 flex items-center justify-center font-bold text-xs text-slate-950">
                      {user?.username
                        ? user.username.charAt(0).toUpperCase()
                        : user?.email
                        ? user.email.charAt(0).toUpperCase()
                        : 'U'}
                    </div>
                  </button>

                  {/* Dropdown Menu */}
                  {isProfileMenuOpen && (
                    <div
                      className="absolute right-0 mt-2 w-64 rounded-2xl bg-surface-900 border border-slate-800 shadow-2xl py-2 z-50 backdrop-blur-xl animate-in fade-in slide-in-from-top-2 duration-150"
                      onClick={() => setIsProfileMenuOpen(false)}
                    >
                      <div className="px-4 py-3 border-b border-slate-800">
                        <p className="text-xs font-medium text-slate-400">Signed in as</p>
                        <p className="text-sm font-bold text-slate-100 truncate mt-0.5">
                          {user?.username || user?.email || 'Learner'}
                        </p>
                        {user?.username && user?.email && (
                          <p className="text-xs text-slate-400 truncate mt-0.5">
                            {user.email}
                          </p>
                        )}
                        <p className="text-[11px] text-brand-400 font-medium mt-1">
                          User ID #{user?.userId || 0}
                        </p>
                      </div>

                      <div className="py-1">
                        <button
                          onClick={onOpenCreateCourse}
                          className="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-slate-200 hover:bg-slate-800/80 transition-colors text-left"
                        >
                          <PlusCircle className="w-4 h-4 text-brand-400" />
                          <span>Import YouTube Course</span>
                        </button>
                        <button
                          onClick={() => setActiveTab('achievements')}
                          className="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-slate-200 hover:bg-slate-800/80 transition-colors text-left"
                        >
                          <Trophy className="w-4 h-4 text-amber-400" />
                          <span>My Achievements & XP</span>
                        </button>
                        <button
                          onClick={onOpenSettings}
                          className="w-full flex items-center gap-2.5 px-4 py-2.5 text-xs text-slate-200 hover:bg-slate-800/80 transition-colors text-left"
                        >
                          <Server className="w-4 h-4 text-sky-400" />
                          <span>Backend Settings</span>
                        </button>
                      </div>

                      <div className="pt-1 border-t border-slate-800/80">
                        <button
                          onClick={onOpenDeleteAccount}
                          className="w-full flex items-center gap-2.5 px-4 py-2 text-xs text-rose-400 hover:bg-rose-950/40 transition-colors text-left"
                        >
                          <Trash2 className="w-4 h-4" />
                          <span>Delete Account</span>
                        </button>
                        <button
                          onClick={logout}
                          className="w-full flex items-center gap-2.5 px-4 py-2 text-xs text-slate-400 hover:text-slate-200 hover:bg-slate-800/80 transition-colors text-left"
                        >
                          <LogOut className="w-4 h-4" />
                          <span>Sign Out</span>
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              </>
            ) : (
              <div className="flex items-center gap-2">
                <button
                  onClick={() => openAuthModal('login')}
                  className="px-3.5 py-1.5 rounded-xl text-sm font-semibold text-slate-300 hover:text-white hover:bg-slate-800 transition-colors"
                >
                  Sign In
                </button>
                <button
                  onClick={() => openAuthModal('register')}
                  className="px-4 py-1.5 rounded-xl bg-accent-600 hover:bg-accent-500 text-white font-bold text-sm shadow-md shadow-accent-600/30 transition-all transform active:scale-95"
                >
                  Sign Up
                </button>
              </div>
            )}

            {/* Mobile menu button */}
            <button
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
              className="md:hidden p-2 rounded-xl text-slate-400 hover:text-slate-200 hover:bg-slate-800"
            >
              {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
            </button>
          </div>
        </div>

        {/* Mobile Navigation Drawer */}
        {isMobileMenuOpen && (
          <div className="md:hidden py-4 border-t border-slate-800 space-y-2">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = activeTab === item.id;
              return (
                <button
                  key={item.id}
                  onClick={() => {
                    setActiveTab(item.id as ActiveTab);
                    setIsMobileMenuOpen(false);
                  }}
                  className={`w-full flex items-center justify-between px-4 py-3 rounded-xl text-sm font-medium ${
                    isActive
                      ? 'bg-accent-600/20 text-accent-300 border border-accent-500/30'
                      : 'text-slate-300 hover:bg-slate-800/60'
                  }`}
                >
                  <div className="flex items-center gap-3">
                    <Icon className="w-5 h-5" />
                    <span>{item.label}</span>
                  </div>
                  {item.badge && (
                    <span className="px-2 py-0.5 rounded text-xs font-bold bg-brand-500/20 text-brand-400">
                      {item.badge}
                    </span>
                  )}
                </button>
              );
            })}

            {isAuthenticated && (
              <div className="pt-2 border-t border-slate-800">
                <button
                  onClick={() => {
                    onOpenCreateCourse();
                    setIsMobileMenuOpen(false);
                  }}
                  className="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-bold bg-brand-500 text-slate-950"
                >
                  <PlusCircle className="w-5 h-5" />
                  <span>Add New Course</span>
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </header>
  );
};
