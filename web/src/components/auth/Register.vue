<script setup lang="ts">
import { ref } from "vue"
import { useRouter } from "vue-router"
import { toast } from "vue-sonner"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Eye, EyeOff, Loader2, UserPlus, Mail, User, MapPin, Lock } from "lucide-vue-next"

const router = useRouter()

// Form data
const formData = ref({
  email: "",
  nama: "", // Changed from 'name' to 'nama' to match API
  alamat: "",
  password: "",
  confirmPassword: "",
})

// Validation errors
const errors = ref<Record<string, string>>({})

// Loading state
const isLoading = ref(false)

// Show password toggles
const showPassword = ref(false)
const showConfirmPassword = ref(false)

// Validate email format
const validateEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// Validate passwords match
const validatePasswords = () => {
  if (formData.value.password !== formData.value.confirmPassword) {
    errors.value.confirmPassword = "Password tidak cocok"
  } else {
    delete errors.value.confirmPassword
  }
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

  // Nama validation
  if (!formData.value.nama) {
    errors.value.nama = "Nama harus diisi"
  } else if (formData.value.nama.length < 2) {
    errors.value.nama = "Nama minimal 2 karakter"
  } else if (formData.value.nama.length > 100) {
    errors.value.nama = "Nama maksimal 100 karakter"
  }

  // Alamat validation
  if (!formData.value.alamat) {
    errors.value.alamat = "Alamat harus diisi"
  } else if (formData.value.alamat.length < 10) {
    errors.value.alamat = "Alamat minimal 10 karakter"
  } else if (formData.value.alamat.length > 500) {
    errors.value.alamat = "Alamat maksimal 500 karakter"
  }

  // Password validation
  if (!formData.value.password) {
    errors.value.password = "Password harus diisi"
  } else if (formData.value.password.length < 6) {
    errors.value.password = "Password minimal 6 karakter"
  } else if (formData.value.password.length > 50) {
    errors.value.password = "Password maksimal 50 karakter"
  }

  // Confirm password validation
  if (!formData.value.confirmPassword) {
    errors.value.confirmPassword = "Konfirmasi password harus diisi"
  }

  // Validate passwords match
  validatePasswords()

  return Object.keys(errors.value).length === 0
}

// Register API request
const registerRequest = async (registerData: { email: string; nama: string; alamat: string; password: string }) => {
  try {
    const response = await fetch("/api/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(registerData),
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.message || "Registrasi gagal")
    }

    const data = await response.json()
    return data
  } catch (error) {
    console.error("Registration error:", error)
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
    // Prepare data for API (exclude confirmPassword)
    const registerData = {
      email: formData.value.email.trim().toLowerCase(),
      nama: formData.value.nama.trim(),
      alamat: formData.value.alamat.trim(),
      password: formData.value.password,
    }

    console.log("Registration attempt:", { ...registerData, password: "******" }) // Don't log password

    // Call register API
    const response = await registerRequest(registerData)

    // Handle successful registration
    console.log("Registration successful:", response)
    toast.success("Registrasi berhasil! Silakan login dengan akun Anda.")

    // Redirect to login page
    router.push("/auth/login")

  } catch (error) {
    console.error("Registration failed:", error)
    if (error instanceof Error) {
      // Handle specific error messages
      if (error.message.includes("email")) {
        errors.value.email = "Email sudah terdaftar"
      } else {
        toast.error(error.message)
      }
    } else {
      toast.error("Terjadi kesalahan saat registrasi. Silakan coba lagi.")
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

const toggleConfirmPasswordVisibility = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-green-50 to-emerald-100">
    <Card class="w-full max-w-md shadow-xl">
      <CardHeader class="text-center">
        <div class="mx-auto mb-4 w-12 h-12 bg-green-600 rounded-full flex items-center justify-center">
          <UserPlus class="h-6 w-6 text-white" />
        </div>
        <CardTitle class="text-2xl font-bold text-gray-900">Daftar</CardTitle>
        <CardDescription class="text-gray-600">
          Buat akun baru untuk mengakses sistem monitoring stunting
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
              class="text-sm font-medium text-gray-700 flex items-center gap-1">
              <Mail class="h-3 w-3" />
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

          <!-- Nama -->
          <div class="space-y-2">
            <Label
              for="nama"
              class="text-sm font-medium text-gray-700 flex items-center gap-1">
              <User class="h-3 w-3" />
              Nama Lengkap *
            </Label>
            <Input
              id="nama"
              v-model="formData.nama"
              @input="clearError('nama')"
              type="text"
              placeholder="Masukkan nama lengkap Anda"
              :class="errors.nama ? 'border-red-500 focus:border-red-500' : ''"
              :disabled="isLoading"
              autocomplete="name"
              required />
            <p
              v-if="errors.nama"
              class="text-sm text-red-500 flex items-center gap-1">
              <span class="text-red-500">⚠️</span>
              {{ errors.nama }}
            </p>
            <div class="text-xs text-gray-500">
              {{ formData.nama.length }}/100 karakter
            </div>
          </div>

          <!-- Alamat -->
          <div class="space-y-2">
            <Label
              for="alamat"
              class="text-sm font-medium text-gray-700 flex items-center gap-1">
              <MapPin class="h-3 w-3" />
              Alamat Lengkap *
            </Label>
            <Textarea
              id="alamat"
              v-model="formData.alamat"
              @input="clearError('alamat')"
              placeholder="Masukkan alamat lengkap Anda (Jalan, RT/RW, Kelurahan, Kecamatan)"
              :class="errors.alamat ? 'border-red-500 focus:border-red-500' : ''"
              :disabled="isLoading"
              class="min-h-[80px] resize-none"
              autocomplete="address-line1"
              required />
            <p
              v-if="errors.alamat"
              class="text-sm text-red-500 flex items-center gap-1">
              <span class="text-red-500">⚠️</span>
              {{ errors.alamat }}
            </p>
            <div class="text-xs text-gray-500">
              {{ formData.alamat.length }}/500 karakter
            </div>
          </div>

          <!-- Password -->
          <div class="space-y-2">
            <Label
              for="password"
              class="text-sm font-medium text-gray-700 flex items-center gap-1">
              <Lock class="h-3 w-3" />
              Password *
            </Label>
            <div class="relative">
              <Input
                id="password"
                v-model="formData.password"
                @input="clearError('password')"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Minimal 6 karakter"
                :class="errors.password ? 'border-red-500 focus:border-red-500 pr-10' : 'pr-10'"
                :disabled="isLoading"
                autocomplete="new-password"
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

          <!-- Konfirmasi Password -->
          <div class="space-y-2">
            <Label
              for="confirmPassword"
              class="text-sm font-medium text-gray-700 flex items-center gap-1">
              <Lock class="h-3 w-3" />
              Konfirmasi Password *
            </Label>
            <div class="relative">
              <Input
                id="confirmPassword"
                v-model="formData.confirmPassword"
                @input="clearError('confirmPassword')"
                @blur="validatePasswords"
                :type="showConfirmPassword ? 'text' : 'password'"
                placeholder="Ulangi password"
                :class="errors.confirmPassword ? 'border-red-500 focus:border-red-500 pr-10' : 'pr-10'"
                :disabled="isLoading"
                autocomplete="new-password"
                required />
              <button
                type="button"
                @click="toggleConfirmPasswordVisibility"
                class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
                :disabled="isLoading">
                <Eye
                  v-if="!showConfirmPassword"
                  class="h-4 w-4" />
                <EyeOff
                  v-else
                  class="h-4 w-4" />
              </button>
            </div>
            <p
              v-if="errors.confirmPassword"
              class="text-sm text-red-500 flex items-center gap-1">
              <span class="text-red-500">⚠️</span>
              {{ errors.confirmPassword }}
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
            <UserPlus
              v-else
              class="mr-2 h-4 w-4" />
            {{ isLoading ? "Memproses..." : "Daftar" }}
          </Button>
        </form>

        <!-- Login Link -->
        <div class="mt-6 text-center text-sm text-gray-600">
          Sudah punya akun?
          <router-link
            to="/auth/login"
            class="font-medium text-green-600 hover:text-green-500 underline">
            Masuk sekarang
          </router-link>
        </div>

        <!-- Development Info -->
        <div class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md text-xs text-yellow-800">
          <p class="font-medium">🔧 Info Development:</p>
          <p>API Endpoint: <code>/api/auth/register</code></p>
          <p>Method: <code>POST</code></p>
          <p>Data: <code>{{ JSON.stringify({ 
            email: formData.email || 'email@example.com', 
            nama: formData.nama || 'Nama Lengkap',
            alamat: formData.alamat || 'Alamat Lengkap',
            password: '******' 
          }) }}</code></p>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
