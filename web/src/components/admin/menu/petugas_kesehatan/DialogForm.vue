<script setup lang="ts">
import { ref, watch, computed } from "vue"
import { toast } from "vue-sonner"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import {
  Stethoscope,
  Save,
  X,
  Users,
  Mail,
  Lock,
  Building,
  Eye,
  EyeOff,
  Shield,
  AlertTriangle,
} from "lucide-vue-next"
import type { PetugasKesehatan } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  petugas: PetugasKesehatan | null
}

interface Emits {
  (e: "close"): void
  (e: "save", petugas: PetugasKesehatan): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data sesuai dengan endpoint
const formData = ref<Partial<PetugasKesehatan>>({
  nama: "",
  email: "",
  id_skpd: "",
  password: "",
})

// Password visibility state
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const confirmPassword = ref("")

// Validation errors
const errors = ref<Record<string, string>>({})

// Dummy SKPD options (dalam implementasi nyata, ambil dari API)
const skpdOptions = [
  {
    id: "S001",
    nama: "Puskesmas Kejaksan",
    jenis: "puskesmas",
    alamat: "Jl. Kejaksan No. 123, Cirebon",
    petugas_count: 5,
  },
  {
    id: "S002",
    nama: "Puskesmas Pekalangan",
    jenis: "puskesmas",
    alamat: "Jl. Pekalangan No. 456, Cirebon",
    petugas_count: 3,
  },
  {
    id: "S003",
    nama: "Kelurahan Kejaksan",
    jenis: "kelurahan",
    alamat: "Jl. Kelurahan No. 789, Cirebon",
    petugas_count: 2,
  },
  {
    id: "S004",
    nama: "Kelurahan Pekalangan",
    jenis: "kelurahan",
    alamat: "Jl. Pekalangan Kelurahan No. 321, Cirebon",
    petugas_count: 1,
  },
  {
    id: "S005",
    nama: "Dinas Kesehatan Kota Cirebon",
    jenis: "skpd",
    alamat: "Jl. Dinas Kesehatan No. 654, Cirebon",
    petugas_count: 8,
  },
  {
    id: "S006",
    nama: "Badan Perencanaan Pembangunan Daerah",
    jenis: "skpd",
    alamat: "Jl. Perencanaan No. 987, Cirebon",
    petugas_count: 4,
  },
]

// Get selected SKPD info
const selectedSkpdInfo = computed(() => {
  if (!formData.value.id_skpd) return null
  return skpdOptions.find((s) => s.id === formData.value.id_skpd) || null
})

// Helper function untuk get SKPD color class
const getSkpdColorClass = (jenis: string): string => {
  switch (jenis?.toLowerCase()) {
    case "puskesmas":
      return "bg-blue-100 text-blue-800 border-blue-200"
    case "kelurahan":
      return "bg-green-100 text-green-800 border-green-200"
    case "skpd":
      return "bg-purple-100 text-purple-800 border-purple-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

// Email validation helper
const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

// Password strength checker
const getPasswordStrength = (password: string): { level: number; text: string; color: string } => {
  let score = 0

  if (password.length >= 8) score += 1
  if (/[a-z]/.test(password)) score += 1
  if (/[A-Z]/.test(password)) score += 1
  if (/[0-9]/.test(password)) score += 1
  if (/[^A-Za-z0-9]/.test(password)) score += 1

  if (score <= 2) return { level: 1, text: "Lemah", color: "text-red-600" }
  if (score <= 3) return { level: 2, text: "Sedang", color: "text-yellow-600" }
  if (score <= 4) return { level: 3, text: "Kuat", color: "text-blue-600" }
  return { level: 4, text: "Sangat Kuat", color: "text-green-600" }
}

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Nama validation
  if (!formData.value.nama || formData.value.nama.trim().length < 3) {
    errors.value.nama = "Nama minimal 3 karakter"
  } else if (formData.value.nama.length > 100) {
    errors.value.nama = "Nama maksimal 100 karakter"
  } else if (!/^[a-zA-Z\s.'-]+$/.test(formData.value.nama)) {
    errors.value.nama =
      "Nama hanya boleh mengandung huruf, spasi, titik, apostrof, dan tanda hubung"
  }

  // Email validation
  if (!formData.value.email) {
    errors.value.email = "Email harus diisi"
  } else if (!isValidEmail(formData.value.email)) {
    errors.value.email = "Format email tidak valid"
  } else if (formData.value.email.length > 100) {
    errors.value.email = "Email maksimal 100 karakter"
  }

  // SKPD validation
  if (!formData.value.id_skpd) {
    errors.value.id_skpd = "SKPD harus dipilih"
  }

  // Password validation (hanya untuk create atau jika password diisi pada edit)
  if (props.mode === "create" || formData.value.password) {
    if (!formData.value.password) {
      errors.value.password = "Password harus diisi"
    } else if (formData.value.password.length < 8) {
      errors.value.password = "Password minimal 8 karakter"
    } else if (formData.value.password.length > 50) {
      errors.value.password = "Password maksimal 50 karakter"
    } else {
      const strength = getPasswordStrength(formData.value.password)
      if (strength.level < 2) {
        errors.value.password =
          "Password terlalu lemah. Gunakan kombinasi huruf besar, kecil, angka, dan simbol"
      }
    }

    // Confirm password validation
    if (formData.value.password !== confirmPassword.value) {
      errors.value.confirmPassword = "Konfirmasi password tidak cocok"
    }
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    nama: "",
    email: "",
    id_skpd: "",
    password: "",
  }
  confirmPassword.value = ""
  errors.value = {}
  showPassword.value = false
  showConfirmPassword.value = false
}

// Load form data for edit mode
const loadFormData = (petugas: PetugasKesehatan) => {
  formData.value = {
    id: petugas.id,
    nama: petugas.nama,
    email: petugas.email,
    id_skpd: petugas.id_skpd,
    password: "", // Don't load existing password for security
    // Extended fields for response
    id_pengguna: petugas.id_pengguna,
    skpd: petugas.skpd,
    jenis_skpd: petugas.jenis_skpd,
    intervensi_count: petugas.intervensi_count,
    created_date: petugas.created_date,
    updated_date: petugas.updated_date,
  }
  confirmPassword.value = ""
  errors.value = {}
  showPassword.value = false
  showConfirmPassword.value = false
}

// Handle save
const handleSave = () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi")

    // Scroll to first error
    const firstErrorElement = document.querySelector(".border-red-500")
    if (firstErrorElement) {
      firstErrorElement.scrollIntoView({ behavior: "smooth", block: "center" })
    }

    return
  }

  // Auto-populate data from selected SKPD
  const selectedSkpd = skpdOptions.find((s) => s.id === formData.value.id_skpd)
  if (selectedSkpd) {
    formData.value.skpd = selectedSkpd.nama
    formData.value.jenis_skpd = selectedSkpd.jenis
  }

  // Set default values for new petugas
  if (props.mode === "create") {
    formData.value.intervensi_count = 0
    formData.value.id_pengguna = `USR_${Date.now()}` // Generate user ID
  }

  // Set timestamps
  const now = new Date().toISOString()
  if (props.mode === "create") {
    formData.value.created_date = now
  }
  formData.value.updated_date = now

  // Remove password jika kosong pada edit mode
  if (props.mode === "edit" && !formData.value.password) {
    delete formData.value.password
  }

  emit("save", formData.value as PetugasKesehatan)
  resetForm()
}

// Handle close
const handleClose = () => {
  resetForm()
  emit("close")
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      if (props.petugas && props.mode === "edit") {
        loadFormData(props.petugas)
      } else {
        resetForm()
      }
    }
  }
)

// Watch for petugas prop changes
watch(
  () => props.petugas,
  (newPetugas) => {
    if (newPetugas && props.mode === "edit" && props.show) {
      loadFormData(newPetugas)
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.nama,
  (newVal) => {
    if (newVal && newVal.trim().length >= 3) {
      delete errors.value.nama
    }
  }
)

watch(
  () => formData.value.email,
  (newVal) => {
    if (newVal && isValidEmail(newVal)) {
      delete errors.value.email
    }
  }
)

watch(
  () => formData.value.password,
  (newVal) => {
    if (newVal && newVal.length >= 8 && getPasswordStrength(newVal).level >= 2) {
      delete errors.value.password
    }
  }
)

watch(
  () => confirmPassword.value,
  (newVal) => {
    if (newVal && newVal === formData.value.password) {
      delete errors.value.confirmPassword
    }
  }
)
</script>

<template>
  <Dialog
    :open="show"
    @update:open="handleClose">
    <DialogContent class="w-[95vw] max-w-4xl max-h-[95vh] overflow-hidden p-0">
      <div class="flex flex-col h-full max-h-[95vh]">
        <!-- Header -->
        <div class="p-4 sm:p-6 border-b bg-background">
          <DialogHeader class="space-y-2">
            <DialogTitle class="text-lg sm:text-xl flex items-center gap-2">
              <Stethoscope class="h-5 w-5 text-blue-600" />
              {{
                mode === "create" ? "Tambah Petugas Kesehatan Baru" : "Edit Data Petugas Kesehatan"
              }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah petugas kesehatan baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi petugas kesehatan. Kosongkan password jika tidak ingin mengubahnya."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Data Pribadi -->
            <Card :class="errors.nama ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Users class="h-4 w-4" />
                  Data Pribadi
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-4">
                  <!-- Nama Lengkap -->
                  <div class="space-y-2">
                    <Label
                      for="nama"
                      class="text-sm font-medium">
                      Nama Lengkap *
                    </Label>
                    <Input
                      id="nama"
                      v-model="formData.nama"
                      placeholder="Masukkan nama lengkap petugas..."
                      :class="errors.nama ? 'border-red-500' : ''" />
                    <p
                      v-if="errors.nama"
                      class="text-sm text-red-500">
                      {{ errors.nama }}
                    </p>
                    <div class="flex justify-between items-center text-xs text-muted-foreground">
                      <span>Minimum 3 karakter, hanya huruf dan spasi</span>
                      <span>{{ formData.nama?.length || 0 }}/100 karakter</span>
                    </div>
                  </div>

                  <!-- Email -->
                  <div class="space-y-2">
                    <Label
                      for="email"
                      class="text-sm font-medium">
                      Email *
                    </Label>
                    <div class="relative">
                      <Mail
                        class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                      <Input
                        id="email"
                        v-model="formData.email"
                        type="email"
                        placeholder="nama@example.com"
                        class="pl-10"
                        :class="errors.email ? 'border-red-500' : ''" />
                    </div>
                    <p
                      v-if="errors.email"
                      class="text-sm text-red-500">
                      {{ errors.email }}
                    </p>
                    <div class="flex justify-between items-center text-xs text-muted-foreground">
                      <span>Email akan digunakan untuk login</span>
                      <span>{{ formData.email?.length || 0 }}/100 karakter</span>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- SKPD Assignment -->
            <Card :class="errors.id_skpd ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Building class="h-4 w-4" />
                  Penugasan SKPD *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select v-model="formData.id_skpd">
                  <SelectTrigger class="w-full" :class="errors.id_skpd ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedSkpdInfo">
                        {{ selectedSkpdInfo.nama }}
                      </template>
                      <template v-else>
                        <span class="text-muted-foreground">Pilih SKPD tempat bertugas</span>
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="skpd in skpdOptions"
                      :key="skpd.id"
                      :value="skpd.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <Badge
                            variant="outline"
                            :class="getSkpdColorClass(skpd.jenis)"
                            class="text-xs">
                            {{ skpd.jenis.charAt(0).toUpperCase() + skpd.jenis.slice(1) }}
                          </Badge>
                          <span class="font-medium">{{ skpd.nama }}</span>
                        </div>
                        <div class="text-xs text-muted-foreground">{{ skpd.alamat }}</div>
                        <div class="text-xs text-blue-600">
                          {{ skpd.petugas_count }} petugas aktif
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_skpd"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_skpd }}
                </p>
                <div
                  v-if="selectedSkpdInfo"
                  class="mt-3 p-3 rounded-md border"
                  :class="getSkpdColorClass(selectedSkpdInfo.jenis)">
                  <div class="flex items-start gap-3">
                    <Building class="h-5 w-5 mt-0.5 opacity-70" />
                    <div class="flex-1">
                      <div class="font-medium text-sm">{{ selectedSkpdInfo.nama }}</div>
                      <div class="text-xs opacity-90 mt-1">{{ selectedSkpdInfo.alamat }}</div>
                      <div class="flex items-center gap-4 mt-2 text-xs">
                        <span>{{
                          selectedSkpdInfo.jenis.charAt(0).toUpperCase() +
                          selectedSkpdInfo.jenis.slice(1)
                        }}</span>
                        <span>{{ selectedSkpdInfo.petugas_count }} petugas aktif</span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Password & Security -->
            <Card :class="errors.password || errors.confirmPassword ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Lock class="h-4 w-4" />
                  Keamanan Akun
                  <div
                    v-if="mode === 'edit'"
                    class="ml-auto">
                    <Badge
                      variant="outline"
                      class="text-xs">
                      Opsional untuk edit
                    </Badge>
                  </div>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-4">
                  <!-- Password -->
                  <div class="space-y-2">
                    <Label
                      for="password"
                      class="text-sm font-medium">
                      Password
                      {{ mode === "create" ? "*" : "(kosongkan jika tidak ingin mengubah)" }}
                    </Label>
                    <div class="relative">
                      <Lock
                        class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                      <Input
                        id="password"
                        v-model="formData.password"
                        :type="showPassword ? 'text' : 'password'"
                        placeholder="Masukkan password..."
                        class="pl-10 pr-10"
                        :class="errors.password ? 'border-red-500' : ''" />
                      <button
                        type="button"
                        @click="showPassword = !showPassword"
                        class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600">
                        <Eye
                          v-if="showPassword"
                          class="h-4 w-4" />
                        <EyeOff
                          v-else
                          class="h-4 w-4" />
                      </button>
                    </div>
                    <p
                      v-if="errors.password"
                      class="text-sm text-red-500">
                      {{ errors.password }}
                    </p>
                    <!-- Password Strength Indicator -->
                    <div
                      v-if="formData.password"
                      class="space-y-2">
                      <div class="flex items-center justify-between text-xs">
                        <span>Kekuatan Password:</span>
                        <span :class="getPasswordStrength(formData.password).color">
                          {{ getPasswordStrength(formData.password).text }}
                        </span>
                      </div>
                      <div class="w-full bg-gray-200 rounded-full h-2">
                        <div
                          class="h-2 rounded-full transition-all duration-300"
                          :class="{
                            'bg-red-500': getPasswordStrength(formData.password).level === 1,
                            'bg-yellow-500': getPasswordStrength(formData.password).level === 2,
                            'bg-blue-500': getPasswordStrength(formData.password).level === 3,
                            'bg-green-500': getPasswordStrength(formData.password).level === 4,
                          }"
                          :style="{
                            width: `${(getPasswordStrength(formData.password).level / 4) * 100}%`,
                          }"></div>
                      </div>
                    </div>
                  </div>

                  <!-- Confirm Password -->
                  <div
                    v-if="formData.password"
                    class="space-y-2">
                    <Label
                      for="confirmPassword"
                      class="text-sm font-medium">
                      Konfirmasi Password *
                    </Label>
                    <div class="relative">
                      <Lock
                        class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
                      <Input
                        id="confirmPassword"
                        v-model="confirmPassword"
                        :type="showConfirmPassword ? 'text' : 'password'"
                        placeholder="Ulangi password..."
                        class="pl-10 pr-10"
                        :class="errors.confirmPassword ? 'border-red-500' : ''" />
                      <button
                        type="button"
                        @click="showConfirmPassword = !showConfirmPassword"
                        class="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600">
                        <Eye
                          v-if="showConfirmPassword"
                          class="h-4 w-4" />
                        <EyeOff
                          v-else
                          class="h-4 w-4" />
                      </button>
                    </div>
                    <p
                      v-if="errors.confirmPassword"
                      class="text-sm text-red-500">
                      {{ errors.confirmPassword }}
                    </p>
                    <!-- Password Match Indicator -->
                    <div
                      v-if="confirmPassword"
                      class="flex items-center gap-2 text-xs">
                      <div
                        class="w-3 h-3 rounded-full"
                        :class="
                          confirmPassword === formData.password ? 'bg-green-500' : 'bg-red-500'
                        "></div>
                      <span
                        :class="
                          confirmPassword === formData.password ? 'text-green-600' : 'text-red-600'
                        ">
                        {{
                          confirmPassword === formData.password
                            ? "Password cocok"
                            : "Password tidak cocok"
                        }}
                      </span>
                    </div>
                  </div>

                  <!-- Password Requirements -->
                  <div class="bg-blue-50 border border-blue-200 rounded-md p-3">
                    <div class="flex items-start gap-2">
                      <Shield class="h-4 w-4 text-blue-600 mt-0.5 flex-shrink-0" />
                      <div class="text-xs text-blue-800">
                        <p class="font-medium mb-1">Persyaratan Password:</p>
                        <ul class="space-y-0.5 list-disc list-inside">
                          <li>Minimal 8 karakter</li>
                          <li>Kombinasi huruf besar dan kecil</li>
                          <li>Minimal 1 angka</li>
                          <li>Minimal 1 karakter khusus (@, #, $, dll)</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Info Box -->
            <div class="bg-orange-50 border border-orange-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <AlertTriangle class="h-5 w-5 text-orange-600 mt-0.5 flex-shrink-0" />
                <div class="text-sm text-orange-800">
                  <p class="font-medium mb-2">📋 Informasi Penting:</p>
                  <ul class="space-y-1 text-xs">
                    <li>• <strong>Email:</strong> Akan digunakan sebagai username untuk login</li>
                    <li>
                      • <strong>SKPD:</strong> Menentukan wilayah kerja dan akses data petugas
                    </li>
                    <li>• <strong>Password:</strong> Harus kuat dan unik untuk keamanan sistem</li>
                    <li v-if="mode === 'edit'">
                      • <strong>Edit Password:</strong> Kosongkan jika tidak ingin mengubah password
                    </li>
                    <li>• <strong>Akun:</strong> Akan aktif otomatis setelah data disimpan</li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Live Preview -->
            <div
              v-if="formData.nama && formData.email && selectedSkpdInfo"
              class="bg-green-50 border border-green-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <div class="text-2xl">✅</div>
                <div class="flex-1">
                  <div class="text-sm font-medium text-green-900 mb-2">Preview Petugas:</div>
                  <div class="space-y-2">
                    <div class="flex items-center gap-2">
                      <Users class="h-4 w-4 text-green-600" />
                      <span class="font-medium text-green-800">{{ formData.nama }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <Mail class="h-4 w-4 text-green-600" />
                      <span class="text-green-700">{{ formData.email }}</span>
                    </div>
                    <div class="flex items-center gap-2">
                      <Building class="h-4 w-4 text-green-600" />
                      <Badge
                        variant="outline"
                        :class="getSkpdColorClass(selectedSkpdInfo.jenis)"
                        class="text-xs">
                        {{
                          selectedSkpdInfo.jenis.charAt(0).toUpperCase() +
                          selectedSkpdInfo.jenis.slice(1)
                        }}
                      </Badge>
                      <span class="text-green-700">{{ selectedSkpdInfo.nama }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Debug Info (untuk development) -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Mode: {{ mode }}<br />
              Has Data: {{ !!props.petugas }}<br />
              Form ID: {{ formData.id || "New" }}<br />
              Nama: {{ formData.nama || "Empty" }}<br />
              Email: {{ formData.email || "Empty" }}<br />
              SKPD: {{ selectedSkpdInfo?.nama || "Not selected" }}<br />
              Password: {{ formData.password ? "Set" : "Empty" }}<br />
              Confirm Password: {{ confirmPassword ? "Set" : "Empty" }}
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="p-4 sm:p-6 border-t bg-background">
          <DialogFooter class="flex flex-col-reverse sm:flex-row gap-3">
            <Button
              variant="outline"
              @click="handleClose"
              class="w-full sm:w-auto">
              <X class="h-4 w-4 mr-2" />
              Batal
            </Button>
            <Button
              @click="handleSave"
              class="w-full sm:w-auto">
              <Save class="h-4 w-4 mr-2" />
              {{ mode === "create" ? "Simpan" : "Perbarui" }}
            </Button>
          </DialogFooter>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
