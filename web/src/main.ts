import { createApp } from "vue"
import { createRouter, createWebHistory } from "vue-router"

import "./style.css"

// Extend Window interface for development debugging
declare global {
  interface Window {
    authUtils: typeof authUtils
  }
}
import App from "./App.vue"
import Home from "./components/Home.vue"

// Auth
import Auth from "./components/auth/Auth.vue"
import Login from "./components/auth/Login.vue"
import Register from "./components/auth/Register.vue"

// Admin
import Admin from "./components/admin/Admin.vue"
import Keluarga from "./components/admin/menu/keluarga/Keluarga.vue"
import Balita from "./components/admin/menu/balita/Balita.vue"
import LaporanMasyarakat from "./components/admin/menu/laporan_masyarakat/LaporanMasyarakat.vue"
import Intervensi from "./components/admin/menu/intervensi/Intervensi.vue"
import RiwayatPemeriksaan from "./components/admin/menu/riwayat_pemeriksaan/RiwayatPemeriksaan.vue"
import PetugasKesehatan from "./components/admin/menu/petugas_kesehatan/PetugasKesehatan.vue"
import Skpd from "./components/admin/menu/skpd/Skpd.vue"

// Community
import Community from "./components/community/Community.vue"
import CommunityKeluarga from "./components/community/menu/keluarga/Keluarga.vue"
import CommunityBalita from "./components/community/menu/balita/Balita.vue"
import CommunityLaporanMasyarakat from "./components/community/menu/laporan_masyarakat/LaporanMasyarakat.vue"

// Types for API responses
interface ProfileResponse {
  data: {
    data: any
    email: string
    id: string
    role: string
  }
  message: string
  status_code: number
}

// Authentication utilities
const authUtils = {
  // Get token from localStorage
  getToken(): string | null {
    return localStorage.getItem("auth_token")
  },

  // Get user data from localStorage
  getUserData(): any | null {
    const userData = localStorage.getItem("user_data")
    if (userData) {
      try {
        return JSON.parse(userData)
      } catch (error) {
        console.error("Error parsing user data:", error)
        return null
      }
    }
    return null
  },

  // Check if user is authenticated
  isAuthenticated(): boolean {
    const token = this.getToken()
    const userData = this.getUserData()
    return !!(token && userData)
  },

  // Get user role
  getUserRole(): string | null {
    const userData = this.getUserData()
    return userData?.role || null
  },

  // Clear authentication data
  clearAuth(): void {
    localStorage.removeItem("auth_token")
    localStorage.removeItem("user_data")
  },

  // Verify token with API
  async verifyTokenWithAPI(): Promise<ProfileResponse | null> {
    const token = this.getToken()
    if (!token) return null

    try {
      const response = await fetch("/api/auth/profile", {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
          "Content-Type": "application/json",
        },
      })

      const data = await response.json()

      if (!response.ok || data.status_code !== 200) {
        console.error("Token verification failed:", data.message)
        this.clearAuth()
        return null
      }

      // Update stored user data
      const userData = {
        id: data.data.id,
        email: data.data.email,
        role: data.data.role,
        data: data.data.data,
      }
      localStorage.setItem("user_data", JSON.stringify(userData))

      return data
    } catch (error) {
      console.error("Error verifying token:", error)
      this.clearAuth()
      return null
    }
  },
}

const routes = [
  {
    path: "/",
    component: Home,
    meta: { requiresAuth: false, allowedRoles: [] },
  },
  {
    path: "/auth",
    component: Auth,
    meta: { requiresAuth: false, allowedRoles: [] },
    children: [
      {
        path: "",
        redirect: "/auth/login",
      },
      {
        path: "login",
        component: Login,
        meta: { title: "Login", requiresAuth: false, allowedRoles: [] },
      },
      {
        path: "register",
        component: Register,
        meta: { title: "Register", requiresAuth: false, allowedRoles: [] },
      },
    ],
  },
  {
    path: "/admin",
    component: Admin,
    meta: { requiresAuth: true, allowedRoles: ["admin"] },
    children: [
      {
        path: "",
        redirect: "/admin/keluarga",
      },
      {
        path: "keluarga",
        component: Keluarga,
        meta: { title: "Admin - Keluarga", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "balita",
        component: Balita,
        meta: { title: "Admin - Balita", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "laporan-masyarakat",
        component: LaporanMasyarakat,
        meta: { title: "Admin - Laporan Masyarakat", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "intervensi",
        component: Intervensi,
        meta: { title: "Admin - Intervensi", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "riwayat-pemeriksaan",
        component: RiwayatPemeriksaan,
        meta: { title: "Admin - Riwayat Pemeriksaan", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "petugas-kesehatan",
        component: PetugasKesehatan,
        meta: { title: "Admin - Petugas Kesehatan", requiresAuth: true, allowedRoles: ["admin"] },
      },
      {
        path: "skpd",
        component: Skpd,
        meta: { title: "Admin - SKPD", requiresAuth: true, allowedRoles: ["admin"] },
      },
    ],
  },
  {
    path: "/community",
    component: Community,
    meta: { requiresAuth: true, allowedRoles: ["masyarakat"] },
    children: [
      {
        path: "",
        redirect: "/community/keluarga",
      },
      {
        path: "keluarga",
        component: CommunityKeluarga,
        meta: { title: "Community - Keluarga", requiresAuth: true, allowedRoles: ["masyarakat"] },
      },
      {
        path: "balita",
        component: CommunityBalita,
        meta: { title: "Community - Balita", requiresAuth: true, allowedRoles: ["masyarakat"] },
      },
      {
        path: "laporan-masyarakat",
        component: CommunityLaporanMasyarakat,
        meta: {
          title: "Community - Laporan Masyarakat",
          requiresAuth: true,
          allowedRoles: ["masyarakat"],
        },
      },
    ],
  },
  // Catch-all route for 404 pages
  {
    path: "/:pathMatch(.*)*",
    redirect: "/",
  },
]

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({ left: 0, top: 0 })
      }, 150)
    })
  },
  routes,
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  console.log(`Navigating from ${from.path} to ${to.path}`)

  // Set document title
  document.title = (to.meta.title as string) || "Stunting Web"

  // Check if route requires authentication
  const requiresAuth = to.meta.requiresAuth as boolean
  const allowedRoles = to.meta.allowedRoles as string[]

  // If route doesn't require authentication, allow access
  if (!requiresAuth) {
    // If user is already authenticated and trying to access auth pages, redirect to appropriate dashboard
    if (to.path.startsWith("/auth") && authUtils.isAuthenticated()) {
      const userRole = authUtils.getUserRole()
      console.log("User already authenticated, redirecting based on role:", userRole)

      if (userRole === "admin") {
        next("/admin")
      } else if (userRole === "petugas_kesehatan") {
        next("/petugas")
      } else if (userRole === "masyarakat") {
        next("/community")
      } else {
        next("/")
      }
    } else {
      next()
    }
    return
  }

  // Route requires authentication
  console.log("Route requires authentication, checking...")

  // Check if user has token
  const token = authUtils.getToken()
  if (!token) {
    console.log("No token found, redirecting to login")
    next(`/auth/login?redirect=${encodeURIComponent(to.fullPath)}`)
    return
  }

  // Verify token with API
  console.log("Token found, verifying with API...")
  try {
    const profileData = await authUtils.verifyTokenWithAPI()

    if (!profileData) {
      console.log("Token verification failed, redirecting to login")
      next(`/auth/login?redirect=${encodeURIComponent(to.fullPath)}`)
      return
    }

    const userRole = profileData.data.role
    console.log("Token verified, user role:", userRole)

    // Check if user role is allowed for this route
    if (allowedRoles.length > 0 && !allowedRoles.includes(userRole)) {
      console.log("User role not allowed for this route")

      // Redirect to appropriate dashboard based on role
      if (userRole === "admin") {
        next("/admin")
      } else if (userRole === "petugas_kesehatan") {
        next("/petugas") // Add petugas routes if needed
      } else if (userRole === "masyarakat") {
        next("/community")
      } else {
        next("/")
      }
      return
    }

    // User is authenticated and has the right role
    console.log("Access granted")
    next()
  } catch (error) {
    console.error("Error during authentication check:", error)
    authUtils.clearAuth()
    next(`/auth/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }
})

// Handle navigation errors
router.onError((error) => {
  console.error("Router error:", error)
})

const app = createApp(App)
app.use(router)
app.mount("#app")

// Add global auth utilities to window for debugging (development only)
if (import.meta.env.DEV) {
  window.authUtils = authUtils
}
