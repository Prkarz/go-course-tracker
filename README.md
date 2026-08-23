## COURSE FLOW

A full-stack gamified learning application designed to help users organize their personal education. It allows users to track external YouTube playlist courses, rack up points for their hard work, and leverage AI to summarize their study materials.

* **Live Application:** https://go-course-tracker.onrender.com/

> **Developer's Note:** Parts of the React/Tailwind frontend were built with the assistance of AI tools. This approach was used to bridge skill gaps, learn modern architectural patterns, and accelerate full-stack development while manually architecting the robust Go backend.

### 🚀 Key Features & Security

* **Gamification Engine:** Tracks user progress, calculates streaks, and awards points using transaction-safe PostgreSQL queries.
* **Bulletproof IDOR Protection:** Uses a secure "Context Backpack" to extract user IDs directly from JWTs, completely preventing Insecure Direct Object Reference (IDOR) attacks.
* **CORS Middleware:** Implements a custom Cross-Origin Resource Sharing (CORS) bridge to securely connect the frontend and backend.
* **Smart Input Validation:** Enforces strict YouTube playlist URL and domain parsing to ensure database integrity.

### 🛠️ Tech Stack

* **Backend:** Go, PostgreSQL (pgx driver), bcrypt password hashing, and JWT authentication.
* **Frontend:** React, Vite, and Tailwind CSS.
* **Deployment:** Deployed live using **Render** for hosting and **Neon** for serverless PostgreSQL.

### 💻 How to Run Locally

* **Database Setup:** Ensure you have PostgreSQL installed and running (or use a remote provider like Neon).
* **Environment Variables:** Create a hidden `.env` file in the root directory containing your `DB_URL`, `SECRET_SAUCE` (JWT secret), `PORT`, and `GOOGLE_API_KEY`.
* **Start the Backend:** Open a terminal in the backend folder and execute `go run main.go` to start the server.
* **Start the Frontend:** Open a separate terminal in the UI folder, run `npm install` to download dependencies, and then run `npm run dev` to launch the Vite development server.