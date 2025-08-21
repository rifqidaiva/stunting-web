<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { toast } from "vue-sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Eye, EyeOff, Loader2, LogIn } from "lucide-vue-next"

const router = useRouter()

// Form data
const formData = ref({
  email: "",
  password: "",
})

// Validation errors
const errors = ref<Record<string, string>>({})

// Loading state
const isLoading = ref(false)

// Show password toggle
const showPassword = ref(false)

// API Response interfaces
interface LoginResponse {
  data: {
    token: string
  }
  message: string
  status_code: number
}

interface ProfileResponse {
  data: {
    data: any // Additional user data
    email: string
    id: string
    role: string
  }
  message: string
  status_code: number
}

// Validate email format
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Email validation
  if (!formData.value.email) {
    errors.value.email = "Email harus diisi"
  } else if (!validateEmail(formData.value.email)) {
    errors.value.email = "Format email tidak valid"
  }

  // Password validation
  if (!formData.value.password) {
    errors.value.password = "Password harus diisi"
  } else if (formData.value.password.length < 6) {
    errors.value.password = "Password minimal 6 karakter"
  }

  return Object.keys(errors.value).length === 0
}

// Login API request
const loginRequest = async (loginData: {
  email: string
  password: string
}): Promise<LoginResponse> => {
  try {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(loginData),
    })

    const data = await response.json()

    if (!response.ok || data.status_code !== 200) {
      throw new Error(data.message || "Login gagal")
    }

    return data
  } catch (error) {
    console.error("Login error:", error)
    throw error
  }
}

// Get user profile using token
const getProfileRequest = async (token: string): Promise<ProfileResponse> => {
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
      throw new Error(data.message || "Gagal mengambil data profil")
    }

    return data
  } catch (error) {
    console.error("Profile fetch error:", error)
    throw error
  }
}

// Handle form submission
const handleSubmit = async () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi")
    return
  }

  isLoading.value = true

  try {
    // Step 1: Login to get token
    const loginData = {
      email: formData.value.email.trim().toLowerCase(),
      password: formData.value.password,
    }

    console.log("Login attempt:", { email: loginData.email }) // Don't log password

    const loginResponse = await loginRequest(loginData)
    console.log("Login successful:", {
      message: loginResponse.message,
      status_code: loginResponse.status_code,
    })

    // Store token
    const token = loginResponse.data.token
    localStorage.setItem("auth_token", token)

    // Step 2: Get user profile using token
    console.log("Fetching user profile...")
    const profileResponse = await getProfileRequest(token)
    console.log("Profile fetched:", {
      id: profileResponse.data.id,
      email: profileResponse.data.email,
      role: profileResponse.data.role,
    })

    // Store user data
    const userData = {
      id: profileResponse.data.id,
      email: profileResponse.data.email,
      role: profileResponse.data.role,
      data: profileResponse.data.data, // Additional user data
    }
    localStorage.setItem("user_data", JSON.stringify(userData))

    // Show success message
    toast.success(`${loginResponse.message} Selamat datang!`)

    // Redirect based on user role
    const userRole = profileResponse.data.role
    console.log("Redirecting user with role:", userRole)

    if (userRole === "admin") {
      router.push("/admin/")
    } else if (userRole === "petugas_kesehatan") {
      router.push("/health-worker")
    } else if (userRole === "masyarakat") {
      router.push("/community")
    } else {
      // Default redirect untuk role yang tidak dikenal
      router.push("/dashboard")
    }
  } catch (error) {
    console.error("Login process failed:", error)

    // Clear stored token if any error occurs
    localStorage.removeItem("auth_token")
    localStorage.removeItem("user_data")

    if (error instanceof Error) {
      // Handle specific error messages
      if (error.message.includes("profil")) {
        toast.error("Login berhasil, tetapi gagal mengambil data profil. Silakan coba lagi.")
      } else {
        toast.error(error.message)
      }
    } else {
      toast.error("Terjadi kesalahan saat login. Silakan coba lagi.")
    }
  } finally {
    isLoading.value = false
  }
}

// Clear specific error when user starts typing
const clearError = (field: string) => {
  if (errors.value[field]) {
    delete errors.value[field]
  }
}

// Toggle password visibility
const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-blue-50 to-indigo-100">
    <Card class="w-full max-w-sm shadow-xl">
      <CardHeader class="text-center">
        <div
          class="mx-auto mb-4 w-12 h-12 bg-blue-600 rounded-full flex items-center justify-center">
          <LogIn class="h-6 w-6 text-white" />
        </div>
        <CardTitle class="text-2xl font-bold text-gray-900">Masuk</CardTitle>
        <CardDescription class="text-gray-600">
          Masukkan email dan password Anda untuk masuk ke akun
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          @submit.prevent="handleSubmit"
          class="space-y-4">
          <!-- Email -->
          <div class="space-y-2">
            <Label
              for="email"
              class="text-sm font-medium text-gray-700">
              Email *
            </Label>
            <Input
              id="email"
              v-model="formData.email"
              @input="clearError('email')"
              type="email"
              placeholder="contoh@email.com"
              :class="errors.email ? 'border-red-500 focus:border-red-500' : ''"
              :disabled="isLoading"
              autocomplete="email"
              required />
            <p
              v-if="errors.email"
              class="text-sm text-red-500 flex items-center gap-1">
              <span class="text-red-500">⚠️</span>
              {{ errors.email }}
            </p>
          </div>

          <!-- Password -->
          <div class="space-y-2">
            <Label
              for="password"
              class="text-sm font-medium text-gray-700">
              Password *
            </Label>
            <div class="relative">
              <Input
                id="password"
                v-model="formData.password"
                @input="clearError('password')"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Masukkan password"
                :class="errors.password ? 'border-red-500 focus:border-red-500 pr-10' : 'pr-10'"
                :disabled="isLoading"
                autocomplete="current-password"
                required />
              <button
                type="button"
                @click="togglePasswordVisibility"
                class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                :disabled="isLoading">
                <Eye
                  v-if="!showPassword"
                  class="h-4 w-4" />
                <EyeOff
                  v-else
                  class="h-4 w-4" />
              </button>
            </div>
            <p
              v-if="errors.password"
              class="text-sm text-red-500 flex items-center gap-1">
              <span class="text-red-500">⚠️</span>
              {{ errors.password }}
            </p>
          </div>

          <!-- Submit Button -->
          <Button
            type="submit"
            class="w-full"
            :disabled="isLoading">
            <Loader2
              v-if="isLoading"
              class="mr-2 h-4 w-4 animate-spin" />
            <LogIn
              v-else
              class="mr-2 h-4 w-4" />
            {{ isLoading ? "Memproses..." : "Masuk" }}
          </Button>
        </form>

        <!-- Loading Steps Indicator -->
        <div
          v-if="isLoading"
          class="mt-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
          <div class="text-xs text-blue-800">
            <div class="flex items-center gap-2 mb-1">
              <div class="w-2 h-2 bg-blue-600 rounded-full animate-pulse"></div>
              <span>Memverifikasi kredensial...</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-2 h-2 bg-gray-300 rounded-full"></div>
              <span class="text-gray-600">Mengambil data profil...</span>
            </div>
          </div>
        </div>

        <!-- Register Link -->
        <div class="mt-6 text-center text-sm text-gray-600">
          Belum punya akun?
          <router-link
            to="/auth/register"
            class="font-medium text-blue-600 hover:text-blue-500 underline">
            Daftar sekarang
          </router-link>
        </div>

        <!-- Development Info -->
        <div
          class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md text-xs text-yellow-800">
          <p class="font-medium">🔧 Info Development:</p>
          <div class="space-y-1 mt-1">
            <p><strong>Step 1:</strong> POST <code>/api/auth/login</code></p>
            <p class="ml-2">
              Data:
              <code>{{
                JSON.stringify({ email: formData.email || "email@example.com", password: "******" })
              }}</code>
            </p>
            <p class="ml-2">Response: <code>{ data: { token }, message, status_code }</code></p>

            <p><strong>Step 2:</strong> GET <code>/api/auth/profile</code></p>
            <p class="ml-2">Headers: <code>Authorization: Bearer {token}</code></p>
            <p class="ml-2">
              Response: <code>{ data: { id, email, role, data }, message, status_code }</code>
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
