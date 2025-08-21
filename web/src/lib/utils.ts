import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import type { Updater } from "@tanstack/vue-table"
import { type Ref } from "vue"
import { toast } from "vue-sonner"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function valueUpdater<T extends Updater<any>>(updaterOrValue: T, ref: Ref) {
  ref.value = typeof updaterOrValue === "function" ? updaterOrValue(ref.value) : updaterOrValue
}

export interface UserData {
  id: string
  email: string
  role: string
  nama?: string
  data?: any
}

export const authUtils = {
  // Get token from localStorage
  getToken(): string | null {
    try {
      return localStorage.getItem("auth_token")
    } catch (error) {
      console.error("Error getting token:", error)
      return null
    }
  },

  // Get user data from localStorage
  getUserData(): UserData | null {
    try {
      const userData = localStorage.getItem("user_data")
      if (userData) {
        return JSON.parse(userData)
      }
    } catch (error) {
      console.error("Error parsing user data:", error)
      // Clear corrupted data
      localStorage.removeItem("user_data")
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

  // Get user name with fallback
  getUserName(): string {
    const userData = this.getUserData()
    return userData?.data?.nama || userData?.nama || "User"
  },

  // Get user email with fallback
  getUserEmail(): string {
    const userData = this.getUserData()
    return userData?.email || "user@example.com"
  },

  // Clear all authentication data
  clearAuth(): void {
    try {
      // Clear localStorage
      localStorage.removeItem("auth_token")
      localStorage.removeItem("user_data")
      localStorage.removeItem("refresh_token")
      localStorage.removeItem("user_preferences")

      // Clear sessionStorage
      sessionStorage.clear()

      // Clear any cookies if used
      document.cookie.split(";").forEach((c) => {
        const eqPos = c.indexOf("=")
        const name = eqPos > -1 ? c.substr(0, eqPos) : c
        document.cookie = `${name}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
      })

      console.log("Authentication data cleared successfully")
    } catch (error) {
      console.error("Error clearing auth data:", error)
    }
  },

  // Logout with API call
  async logout(redirectTo: string = "/auth/login"): Promise<void> {
    const token = this.getToken()

    try {
      // Call logout API if token exists
      if (token) {
        console.log("Calling logout API...")
        const response = await fetch("/api/auth/logout", {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
            "Content-Type": "application/json",
          },
        })

        if (!response.ok) {
          console.warn("Logout API failed, but continuing with local cleanup")
        } else {
          const data = await response.json()
          console.log("Logout API successful:", data)
        }
      }
    } catch (error) {
      console.error("Logout API error:", error)
    } finally {
      // Always clear local data regardless of API response
      this.clearAuth()

      // Show success message
      toast.success("Logout berhasil. Sampai jumpa!")

      // Redirect using window.location for clean navigation
      if (typeof window !== "undefined") {
        setTimeout(() => {
          window.location.href = redirectTo
        }, 500) // Small delay to show toast
      }
    }
  },

  // Force logout (without API call)
  forceLogout(message: string = "Sesi berakhir, silakan login kembali"): void {
    this.clearAuth()
    toast.warning(message)

    if (typeof window !== "undefined") {
      setTimeout(() => {
        window.location.href = "/auth/login"
      }, 1000) // Show toast before redirect
    }
  },

  // Update user data in localStorage
  updateUserData(userData: Partial<UserData>): void {
    try {
      const currentData = this.getUserData()
      if (currentData) {
        const updatedData = { ...currentData, ...userData }
        localStorage.setItem("user_data", JSON.stringify(updatedData))
        console.log("User data updated:", updatedData)
      }
    } catch (error) {
      console.error("Error updating user data:", error)
    }
  },

  // Refresh token and user data from API
  async refreshUserData(): Promise<UserData | null> {
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

      if (!response.ok) {
        console.error("Failed to refresh user data")
        return null
      }

      const data = await response.json()
      if (data.status_code === 200) {
        const userData = {
          id: data.data.id,
          email: data.data.email,
          role: data.data.role,
          data: data.data.data,
        }

        // Update stored user data
        localStorage.setItem("user_data", JSON.stringify(userData))
        console.log("User data refreshed from API:", userData)

        return userData
      }

      return null
    } catch (error) {
      console.error("Error refreshing user data:", error)
      return null
    }
  },
}
