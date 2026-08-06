# 🎨 SpotSync Frontend — Step-by-Step Implementation Guide (Astro + TypeScript) for Puku CLI

> **Target Platform:** Astro 4+ (TypeScript + Tailwind CSS + Nano Stores + Lucide SVG Icons + Native Fetch API)  
> **Design Theme:** Dark Mode Glassmorphism (Deep Slate background, Neon EV Green `#10b981` accents, Electric Blue `#3b82f6` highlights)  
> **Key Features:** Client-side hydration with Nano Stores, Role-based view guards (Driver & Admin), Dynamic Zone Availability visual progress bars, Concurrency-aware Spot Reservation modal, Driver Reservation history dashboard, Admin Zone creation & System monitor table.

---

## 🤖 Instructions for Puku CLI Execution

When instructed to implement a step (e.g., `"Implement Step 1 of Frontend.MD"`):
1. **Read the target Step completely.** Each step contains full context, file paths, TypeScript interface definitions, Astro component structure, client-side script logic, styling rules, and verification commands.
2. **Ensure complete TypeScript typing.** Do not use `any` types. Define clear interfaces in `src/types/`.
3. **Execute all file creation and modifications** described in that specific step.
4. **Run the provided Verification Commands** to test your implementation before declaring the step complete.

---

## 📁 Frontend Directory Blueprint (Astro + TypeScript)

```
frontend/
├── public/
│   └── favicon.svg
├── src/
│   ├── components/
│   │   ├── Navbar.astro               # Top navigation bar with dynamic client state
│   │   ├── ZoneCard.astro             # Zone card component with dynamic capacity bar
│   │   ├── ZoneFilter.astro           # Type filter tabs & search input
│   │   ├── ReserveModal.astro         # Modal for spot booking with license plate input
│   │   ├── ReservationCard.astro      # Driver reservation card with cancel action
│   │   ├── CreateZoneModal.astro      # Admin modal to create new parking zone
│   │   └── AdminReservationTable.astro# Admin table monitor of all system bookings
│   ├── layouts/
│   │   └── Layout.astro               # Core Astro page container layout & toast script
│   ├── pages/
│   │   ├── index.astro                # Public/Driver landing page & zone catalog
│   │   ├── login.astro                # User login page
│   │   ├── register.astro             # User registration page
│   │   ├── my-reservations.astro      # Driver reservations dashboard
│   │   └── admin.astro                # Admin management & system stats page
│   ├── services/
│   │   ├── api.ts                     # Fetch client wrapper with Bearer token interceptor
│   │   ├── authService.ts             # Auth API endpoints
│   │   ├── zoneService.ts             # Parking Zone API endpoints
│   │   └── reservationService.ts      # Reservation API endpoints
│   ├── store/
│   │   └── authStore.ts               # Nano Stores for reactive client auth state
│   ├── styles/
│   │   └── global.css                 # Tailwind CSS imports & glassmorphism utilities
│   ├── types/
│   │   ├── auth.ts                    # User & Auth interfaces
│   │   ├── zone.ts                    # Zone interfaces
│   │   ├── reservation.ts             # Reservation interfaces
│   │   └── api.ts                     # Standard API response interfaces
│   └── env.d.ts                       # Environment variable typings
├── .env.example
├── .env
├── astro.config.mjs
├── package.json
├── tailwind.config.mjs
└── tsconfig.json
```

---

## 📌 Step 1: Astro Project Initialization & Design System Setup

### 🎯 Objective
Initialize an Astro TypeScript project, install dependencies (`@astrojs/tailwind`, `tailwindcss`, `nanostores`, `@nanostores/persistent`), configure Tailwind CSS with custom glassmorphism tokens, setup environment variables, and create folder structure.

### 📋 Prerequisites
- Node.js 18+ installed.

### 📁 Files to Create / Modify
- `frontend/package.json`
- `frontend/astro.config.mjs`
- `frontend/tailwind.config.mjs`
- `frontend/src/styles/global.css`
- `frontend/.env`
- `frontend/.env.example`
- `frontend/src/env.d.ts`

### 📝 Step Context & Implementation Details

1. **Astro Initialization & Package Setup (`package.json`)**
   Dependencies:
   - `astro` (Astro core framework)
   - `@astrojs/tailwind` (Tailwind CSS integration for Astro)
   - `tailwindcss` & `autoprefixer`
   - `nanostores` & `@nanostores/persistent` (Lightweight client-side reactive state)

2. **Astro Configuration (`astro.config.mjs`)**
   ```js
   import { defineConfig } from 'astro/config';
   import tailwind from '@astrojs/tailwind';

   export default defineConfig({
     integrations: [tailwind()],
     output: 'static', // Build static output with client-side hydration scripts
   });
   ```

3. **Tailwind & Custom CSS (`src/styles/global.css`)**
   Define dark glassmorphism aesthetic:
   - Background: Slate 950 `#020617` / Slate 900 `#0f172a`.
   - Cards: `.glass-card` (`bg-slate-900/80 backdrop-blur-md border border-slate-800/80 rounded-2xl`).
   - Accents: Neon EV Green `#10b981`, Electric Blue `#3b82f6`, Warning Amber `#f59e0b`, Danger Coral `#ef4444`.

4. **Environment Variables (`.env`)**
   ```env
   PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
   ```

### 💻 Execution Commands
```bash
npm create astro@latest frontend -- --template bare --typescript strict
cd frontend
npm install @astrojs/tailwind tailwindcss autoprefixer nanostores @nanostores/persistent
mkdir -p src/components src/layouts src/pages src/services src/store src/styles src/types
```

### ✅ Step 1 Verification
Run Astro dev server:
```bash
npm run dev
```
Expected Output: Astro dev server running at `http://localhost:4321`.

---

## 📌 Step 2: TypeScript Interfaces & Native Fetch API Client Setup

### 🎯 Objective
Define TypeScript domain interfaces matching the Go backend DTOs and create a robust native Fetch API wrapper (`src/services/api.ts`) supporting automatic `Authorization: Bearer <token>` injection and central error handling.

### 📋 Prerequisites
- Step 1 completed.

### 📁 Files to Create / Modify
- `frontend/src/types/api.ts`
- `frontend/src/types/auth.ts`
- `frontend/src/types/zone.ts`
- `frontend/src/types/reservation.ts`
- `frontend/src/services/api.ts`

### 📝 Step Context & Implementation Details

1. **Domain Interfaces**
   - `src/types/api.ts`:
     ```typescript
     export interface ApiResponse<T> {
       success: boolean;
       message: string;
       data: T;
       errors?: string;
     }
     ```
   - `src/types/auth.ts`: `User`, `LoginCredentials`, `RegisterPayload`, `AuthResponse`. Roles: `'driver' | 'admin'`.
   - `src/types/zone.ts`: `ParkingZone` (`id`, `name`, `type`: `'general' | 'ev_charging' | 'covered'`, `total_capacity`, `available_spots`, `price_per_hour`, `created_at`).
   - `src/types/reservation.ts`: `Reservation` (`id`, `user_id`, `zone_id`, `license_plate`, `status`: `'active' | 'completed' | 'cancelled'`, `zone?`, `user?`, `created_at`).

2. **Fetch API Client Wrapper (`src/services/api.ts`)**
   - Base URL from `import.meta.env.PUBLIC_API_BASE_URL`.
   - Export helper function `apiFetch<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>>`.
   - Automatically read token from `localStorage.getItem('spotsync_token')` and attach `Authorization: Bearer ${token}` header if present.
   - Parse response JSON. Handle HTTP status errors (401, 403, 409, 500) with clear throw messages.

### 💻 Execution Commands
```bash
cd frontend
npx astro check
```

### ✅ Step 2 Verification
TypeScript check pass without errors.

---

## 📌 Step 3: Auth State Management with Nano Stores

### 🎯 Objective
Create a lightweight, reactive client-side store (`src/store/authStore.ts`) using Nano Stores persistent state to manage authentication state, user roles, JWT tokens, and login/logout methods across Astro pages.

### 📋 Prerequisites
- Step 2 completed.

### 📁 Files to Create / Modify
- `frontend/src/services/authService.ts`
- `frontend/src/store/authStore.ts`

### 📝 Step Context & Implementation Details

1. **Auth Service (`src/services/authService.ts`)**
   Implement methods calling backend API:
   - `register(payload: RegisterPayload)` → `POST /auth/register`
   - `login(credentials: LoginCredentials)` → `POST /auth/login`

2. **Nano Stores Auth Manager (`src/store/authStore.ts`)**
   Use `persistentAtom` from `@nanostores/persistent`:
   ```typescript
   import { persistentAtom } from '@nanostores/persistent';
   import type { User } from '../types/auth';

   export const $token = persistentAtom<string | null>('spotsync_token', null);
   export const $user = persistentAtom<User | null>('spotsync_user', null, {
     encode: JSON.stringify,
     decode: JSON.parse,
   });

   export function setSession(token: string, user: User) {
     $token.set(token);
     $user.set(user);
   }

   export function clearSession() {
     $token.set(null);
     $user.set(null);
   }

   export function isAuthenticated(): boolean {
     return !!$token.get();
   }

   export function isAdmin(): boolean {
     const u = $user.get();
     return u?.role === 'admin';
   }
   ```

### 💻 Execution Commands
```bash
cd frontend
npx astro check
```

### ✅ Step 3 Verification
TypeScript verification pass.

---

## 📌 Step 4: Master Layout, Navbar & Dynamic Client Hydration

### 🎯 Objective
Create the core Astro layout component (`Layout.astro`), responsive top navigation bar (`Navbar.astro`), and client-side page authentication guards.

### 📋 Prerequisites
- Step 3 completed.

### 📁 Files to Create / Modify
- `frontend/src/layouts/Layout.astro`
- `frontend/src/components/Navbar.astro`

### 📝 Step Context & Implementation Details

1. **Master Layout (`src/layouts/Layout.astro`)**
   - HTML boilerplate, Google Fonts (Inter / Outfit), global CSS import (`../styles/global.css`).
   - Slot for page content `<slot />`.
   - Embedded Toast notification container (simple custom toast notification script or lightweight toast system).

2. **Navbar Component (`src/components/Navbar.astro`)**
   - Brand Logo with Zap icon + "SpotSync".
   - Links: "Parking Zones" (`/`), "My Reservations" (`/my-reservations`), "Admin Dashboard" (`/admin`).
   - Interactive `<script>` inside `Navbar.astro`:
     - Listens to Nano Stores `$user` updates.
     - Dynamically renders logged-in user profile, role badge (`admin` emerald / `driver` blue), or "Sign In" / "Register" buttons.
     - Handles Logout button click -> calls `clearSession()` -> redirects to `/login`.

### 💻 Execution Commands
```bash
cd frontend
npm run dev
```

### ✅ Step 4 Verification
Open `http://localhost:4321`. Navbar renders with brand styling and dynamic user state container.

---

## 📌 Step 5: Authentication Pages (`login.astro` & `register.astro`)

### 🎯 Objective
Build interactive Login and Register Astro pages featuring modern dark glassmorphism forms, role selection toggle, quick demo login buttons, and instant client redirects.

### 📋 Prerequisites
- Step 4 completed.

### 📁 Files to Create / Modify
- `frontend/src/pages/login.astro`
- `frontend/src/pages/register.astro`

### 📝 Step Context & Implementation Details

1. **Login Page (`src/pages/login.astro`)**
   - Glassmorphism form card with Email and Password input fields.
   - Quick Demo Account Buttons:
     - `Fill Driver` (`john@spotsync.com` / `password123`)
     - `Fill Admin` (`admin@spotsync.com` / `adminpassword`)
   - `<script>` tag: Intercept form submit → call `authService.login()` → invoke `setSession(data.token, data.user)` → redirect to `/` or `/admin`.

2. **Register Page (`src/pages/register.astro`)**
   - Inputs: Full Name, Email, Password.
   - Account Role Selector Cards:
     - 🚗 **Driver Account** (Default - reserve parking & EV spots)
     - ⚡ **Admin Account** (Manage parking zones & system overview)
   - `<script>` tag: Form submit → call `authService.register()` → redirect to `/login` with success alert.

### 💻 Execution Commands
```bash
cd frontend
npm run dev
```

### ✅ Step 5 Verification
Navigate to `http://localhost:4321/register` and `http://localhost:4321/login`. Test input forms, quick demo fill buttons, and role card toggle.

---

## 📌 Step 6: Zone Catalog, Dynamic Spot Calculations & Reserve Modal

### 🎯 Objective
Build the main Parking Zone catalog page (`index.astro`), Zone Card component with visual capacity bars (`ZoneCard.astro`), Filter bar (`ZoneFilter.astro`), and the Concurrency-Aware Spot Reservation Modal (`ReserveModal.astro`).

### 📋 Prerequisites
- Step 5 completed.

### 📁 Files to Create / Modify
- `frontend/src/services/zoneService.ts`
- `frontend/src/services/reservationService.ts`
- `frontend/src/components/ZoneCard.astro`
- `frontend/src/components/ZoneFilter.astro`
- `frontend/src/components/ReserveModal.astro`
- `frontend/src/pages/index.astro`

### 📝 Step Context & Implementation Details

1. **Services**
   - `zoneService.ts`: `getZones()`, `getZoneById(id)`.
   - `reservationService.ts`: `createReservation({ zone_id, license_plate })`.

2. **Zone Card Component (`src/components/ZoneCard.astro`)**
   Renders:
   - Zone Name (e.g., "Terminal 1 EV Charging").
   - Type Badge (`ev_charging` green lightning bolt, `general` blue car, `covered` purple shield).
   - Price per hour tag (`$5.50/hr`).
   - **Dynamic Capacity Progress Bar:**
     - Width percentage: `(available_spots / total_capacity) * 100%`.
     - Color logic: Green (>50%), Yellow (20-50%), Red (<20%).
     - Counter label: `"14 of 20 spots available"`.
   - "Reserve Spot" action button (opens `ReserveModal` with selected `zone_id`).

3. **Reserve Spot Modal (`src/components/ReserveModal.astro`)**
   - Native HTML `<dialog>` or overlay modal.
   - Input: Vehicle License Plate (`license_plate`).
   - `<script>`: Submit form → send request to `reservationService.createReservation()` → handle `409 Conflict` (Zone is full) gracefully with error notice → close modal & reload zones on success.

4. **Home Page (`src/pages/index.astro`)**
   - Hero banner: "Smart Parking & EV Charging Reservation".
   - Stats Counter bar (Total Zones, Available EV Spots).
   - Filter Tabs: `All Zones`, `⚡ EV Charging`, `🚗 General`, `🛡️ Covered`.
   - Grid layout rendering zone cards.

### 💻 Execution Commands
```bash
cd frontend
npm run dev
```

### ✅ Step 6 Verification
Fetch zones from Go backend. Click "Reserve Spot" on an EV Charging zone, enter vehicle plate `EV-101`, submit, and verify spot count updates dynamically.

---

## 📌 Step 7: Driver "My Reservations" Dashboard

### 🎯 Objective
Build the driver's reservation management dashboard page (`my-reservations.astro`), displaying booking history, license plate badges, and instant cancellation functionality.

### 📋 Prerequisites
- Step 6 completed.

### 📁 Files to Create / Modify
- `frontend/src/components/ReservationCard.astro`
- `frontend/src/pages/my-reservations.astro`

### 📝 Step Context & Implementation Details

1. **Reservation Card (`src/components/ReservationCard.astro`)**
   Renders:
   - Reservation ID `#105`.
   - Zone details & Zone Type icon.
   - Vehicle License Plate badge (`ABC-1234`).
   - Timestamp.
   - Status Badge: `active` (Green), `cancelled` (Gray).
   - Action Button: "Cancel Spot" (visible when status is `active`). Triggers confirmation dialog -> calls `DELETE /reservations/:id`.

2. **My Reservations Page (`src/pages/my-reservations.astro`)**
   - Auth guard check in `<script>`: If not authenticated, redirect to `/login`.
   - Fetches driver reservations via `reservationService.getMyReservations()`.
   - Displays list of `ReservationCard` items or empty state vector.

### 💻 Execution Commands
```bash
cd frontend
npm run dev
```

### ✅ Step 7 Verification
Navigate to `http://localhost:4321/my-reservations`. Active driver reservations render. Clicking "Cancel Spot" changes status to `cancelled` and frees up zone capacity on backend.

---

## 📌 Step 8: Admin Dashboard (Zone Creation & Global System Monitor)

### 🎯 Objective
Build the Admin Dashboard (`admin.astro`) containing analytical overview cards, a modal for creating new parking zones (`CreateZoneModal.astro`), and a full system reservations monitor table (`AdminReservationTable.astro`).

### 📋 Prerequisites
- Step 7 completed.

### 📁 Files to Create / Modify
- `frontend/src/components/CreateZoneModal.astro`
- `frontend/src/components/AdminReservationTable.astro`
- `frontend/src/pages/admin.astro`

### 📝 Step Context & Implementation Details

1. **Create Zone Modal (`src/components/CreateZoneModal.astro`)**
   Form inputs:
   - Zone Name (text)
   - Zone Type (dropdown: `ev_charging`, `general`, `covered`)
   - Total Capacity (integer > 0)
   - Price Per Hour (decimal > 0)
   - `<script>` submit handler → calls `zoneService.createZone()` with Admin Bearer token → closes modal & refetches zone list.

2. **Admin Reservations Table (`src/components/AdminReservationTable.astro`)**
   Renders comprehensive system bookings table:
   - Booking ID | User Email | Zone Name | License Plate | Status | Reserved At
   - Filter dropdown by status (`All`, `Active`, `Cancelled`).

3. **Admin Page (`src/pages/admin.astro`)**
   - Auth Guard script: Verifies `isAdmin()` is true. Redirects non-admins to `/`.
   - Analytics Cards: Total Zones, System Active Bookings, Total EV Capacity.
   - Section 1: Parking Zone List + "Create Zone" trigger button.
   - Section 2: Global System Reservations Monitor Table.

### 💻 Execution Commands
```bash
cd frontend
npm run dev
```

### ✅ Step 8 Verification
Log in as Admin (`admin@spotsync.com`). Navigate to `/admin`.
- Click "Add New Zone" → Create zone "Terminal 3 EV Hub" → Appears in catalog.
- System Reservations table lists all driver bookings.

---

## 📌 Step 9: Toast Notifications & Final Production Build

### 🎯 Objective
Polish client UI interactions, ensure error feedback toasts on all API operations, verify mobile/desktop responsiveness, and confirm a clean Astro production build.

### 📋 Prerequisites
- Step 8 completed.

### 📁 Files to Modify
- `frontend/src/layouts/Layout.astro`
- `frontend/src/styles/global.css`

### 📝 Step Context & Implementation Details
1. Ensure client script toast helper handles success/error alerts cleanly.
2. Verify responsive layout across Mobile (375px), Tablet (768px), and Desktop (1440px).
3. Run complete Astro build test.

### 💻 Execution Commands
```bash
cd frontend
npm run build
npm run preview
```

### ✅ Step 9 Verification Checklist
- `npm run build` generates production static build directory `dist/` without TypeScript or syntax errors.
- Preview server runs cleanly at `http://localhost:4321`.
- Complete flow functional: Register/Login -> Browse Zones -> Reserve Spot with License Plate -> View My Reservations -> Cancel Spot -> Admin Zone Creation & Monitoring.
