# ⚡ Course Flow Frontend

A modern, high-performance, production-ready frontend web application for the **Course Flow Go Backend**. Built with **Vite**, **React**, **TypeScript**, and **Tailwind CSS**.

---

## 🌟 Key Features

1. **AI-Powered Playlist Import & Insights**:
   - Accepts YouTube playlist URLs (`https://www.youtube.com/playlist?list=...`).
   - Automatically communicates with Go backend to fetch YouTube metadata and Gemini AI summaries & tags.
   - Real-time client-side playlist URL validator.

2. **Gamification & Habit Tracking**:
   - **Daily Streaks (🔥)**: Automatically calculates and maintains consecutive learning days.
   - **XP & Levels (💎)**: Awards +25 XP per completed lesson, +50 XP daily login rewards, and +100 XP upon course completion.
   - **Badges & Milestones (🏆)**: Visual progression ladder from *Novice Explorer* to *Grandmaster Learner*.

3. **Interactive Course Studio & YouTube Player**:
   - Zero-lag memoized 16:9 embedded YouTube player.
   - Curriculum playlist sidebar with exact video durations, clean sequential index order, and instant checkmarks.
   - Mark as Completed button with celebration confetti bursts.

4. **Complete Endpoint Coverage**:
   - `POST /users/create` — User registration
   - `POST /users/login` — User authentication & JWT session management
   - `POST /users/delete` — Account deletion
   - `POST /courses/mycourses` — List user courses with completion percentages
   - `POST /courses/create` — Import playlist & trigger Gemini AI analysis
   - `POST /courses/detail` — Load course curriculum & video completion status
   - `POST /courses/start` — Enroll and start course
   - `POST /courses/delete` — Remove course
   - `POST /courses/videos/viewed` — Mark video as completed & reward points
   - `POST /courses/progress` — Update course completion percentage & reward points
   - `GET /stats/me` — Fetch user streak, total points, and last active timestamp

5. **Backend Connection Monitor**:
   - Real-time health check pill with response latency.
   - Built-in Backend Settings modal to customize the Go server API URL (defaults to `http://localhost:8080`).

---

## 🚀 Quick Setup & Run Instructions

### Step 1: Install Dependencies
Open a terminal in the `frontend` folder:
```bash
cd frontend
npm install
```

### Step 2: Start Frontend Dev Server
```bash
npm run dev
```
The application will start at **http://localhost:5173**.

### Step 3: Start the Go Backend Server
In the root directory of the project:
```bash
go run main.go
```
Backend will run at **http://localhost:8080** with CORS enabled.

### Step 4: Build for Production
```bash
npm run build
npm run preview
```
