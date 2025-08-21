<script setup lang="ts">
import { ref, watch, computed, onMounted } from "vue"
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
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { 
  FileText, 
  Save, 
  X, 
  Users, 
  Phone, 
  Baby, 
  AlertTriangle,
  Calendar as CalendarIcon,
  UserCheck,
  Building,
  Shield,
  Loader2
} from "lucide-vue-next"
import { format } from "date-fns"
import { id } from "date-fns/locale"
import { CalendarDate, type DateValue } from "@internationalized/date"
import { authUtils } from "@/lib/utils"
import type { LaporanMasyarakat } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  laporan: LaporanMasyarakat | null
  loading?: boolean
}

interface Emits {
  (e: "close"): void
  (e: "save", laporan: LaporanMasyarakat): void
}

// API Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface MasyarakatOption {
  id: string
  nama: string
  email: string
  alamat: string
  nomor_hp: string
  kelurahan: string
  kecamatan: string
}

interface GetAllMasyarakatResponse {
  data: MasyarakatOption[]
  total: number
}

interface BalitaOption {
  id: string
  nama: string
  umur: string
  jenis_kelamin: string
  nama_ayah: string
  nama_ibu: string
  nomor_kk: string
  alamat: string
  kelurahan: string
  kecamatan: string
  nomor_hp_keluarga: string
}

interface GetAllBalitaResponse {
  data: BalitaOption[]
  total: number
}

interface StatusLaporanOption {
  id: string
  status: string
  description: string
}

interface GetAllStatusLaporanResponse {
  data: StatusLaporanOption[]
  total: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data sesuai dengan endpoint
const formData = ref<Partial<LaporanMasyarakat>>({
  id_masyarakat: "",
  id_balita: "",
  id_status_laporan: "",
  hubungan_dengan_balita: "",
  nomor_hp_pelapor: "",
  nomor_hp_keluarga_balita: "",
  tanggal_laporan: "",
})

// Date picker state
const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

// Validation errors
const errors = ref<Record<string, string>>({})

// Master data from API
const masyarakatOptions = ref<MasyarakatOption[]>([])
const balitaOptions = ref<BalitaOption[]>([])
const statusLaporanOptions = ref<StatusLaporanOption[]>([])

// Loading states
const masyarakatLoading = ref(false)
const balitaLoading = ref(false)
const statusLoading = ref(false)

// Enhanced pelapor options dengan admin
const pelaporOptions = computed(() => {
  // Simulasi data admin saat ini (dalam implementasi nyata, ambil dari auth store/context)
  // Get current admin data from API/auth
  const adminData = authUtils.getUserData() // Assuming this returns admin data
  const currentAdmin = {
    id: "ADMIN",
    type: "admin",
    nama: adminData?.nama || "Admin Sistem",
    email: adminData?.email || "admin@stuntingweb.com",
    alamat: "Dinas Komunikasi Informatika dan Statistik",
    nomor_hp: "-",
    kelurahan: "-",
    kecamatan: "-",
    isAdmin: true
  }

  // Combine admin option with masyarakat options
  return [
    currentAdmin,
    ...masyarakatOptions.value.map(m => ({
      id: m.id,
      type: "masyarakat",
      nama: m.nama,
      email: m.email,
      alamat: m.alamat,
      nomor_hp: m.nomor_hp,
      kelurahan: m.kelurahan,
      kecamatan: m.kecamatan,
      isAdmin: false
    }))
  ]
})

// API request function
const apiRequest = async <T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<ApiResponse<T>> => {
  const token = authUtils.getToken()
  if (!token) {
    throw new Error("No authentication token found")
  }

  const defaultOptions: RequestInit = {
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${token}`,
    },
  }

  const config = { ...defaultOptions, ...options }
  config.headers = { ...defaultOptions.headers, ...options.headers }

  try {
    const response = await fetch(`/api${endpoint}`, config)
    const data = await response.json()

    if (!response.ok || data.status_code !== 200) {
      throw new Error(data.message || `HTTP error! status: ${response.status}`)
    }

    return data
  } catch (error) {
    console.error(`API request failed for ${endpoint}:`, error)
    throw error
  }
}

// Fetch masyarakat master data
const fetchMasyarakatOptions = async () => {
  masyarakatLoading.value = true
  try {
    console.log("Fetching masyarakat master data from API...")
    
    const response = await apiRequest<GetAllMasyarakatResponse>("/admin/master-masyarakat")
    
    masyarakatOptions.value = response.data.data
    console.log("Masyarakat options loaded:", masyarakatOptions.value.length, "items")
    
  } catch (error) {
    console.error("Error fetching masyarakat options:", error)
    toast.error("Gagal memuat data masyarakat. Menggunakan data default.")
    
    // Fallback to static data if API fails
    masyarakatOptions.value = [
      {
        id: "M001",
        nama: "Siti Nurhaliza",
        email: "siti@example.com",
        alamat: "Jl. Merdeka No. 123, Kejaksan",
        nomor_hp: "081234567890",
        kelurahan: "Kejaksan",
        kecamatan: "Kejaksan"
      },
      {
        id: "M002",
        nama: "Ali Pelapor",
        email: "ali@example.com",
        alamat: "Jl. Sudirman No. 456, Pekalangan",
        nomor_hp: "081234567892",
        kelurahan: "Pekalangan",
        kecamatan: "Pekalangan"
      }
    ]
  } finally {
    masyarakatLoading.value = false
  }
}

// Fetch balita master data
const fetchBalitaOptions = async () => {
  balitaLoading.value = true
  try {
    console.log("Fetching balita master data from API...")
    
    const response = await apiRequest<GetAllBalitaResponse>("/admin/balita/get")
    
    balitaOptions.value = response.data.data.map(b => ({
      id: b.id,
      nama: b.nama,
      umur: b.umur,
      jenis_kelamin: b.jenis_kelamin,
      nama_ayah: b.nama_ayah,
      nama_ibu: b.nama_ibu,
      nomor_kk: b.nomor_kk,
      alamat: b.alamat,
      kelurahan: b.kelurahan,
      kecamatan: b.kecamatan,
      nomor_hp_keluarga: "081234567891" // This might need to be mapped from the actual API response
    }))
    
    console.log("Balita options loaded:", balitaOptions.value.length, "items")
    
  } catch (error) {
    console.error("Error fetching balita options:", error)
    toast.error("Gagal memuat data balita. Menggunakan data default.")
    
    // Fallback to static data if API fails
    balitaOptions.value = [
      {
        id: "B001",
        nama: "Ahmad Fauzi",
        umur: "24 bulan",
        jenis_kelamin: "L",
        nama_ayah: "Budi Santoso",
        nama_ibu: "Siti Nurhaliza",
        nomor_kk: "1234567890123456",
        alamat: "Jl. Merdeka No. 123",
        kelurahan: "Kejaksan",
        kecamatan: "Kejaksan",
        nomor_hp_keluarga: "081234567891"
      }
    ]
  } finally {
    balitaLoading.value = false
  }
}

// Fetch status laporan master data
const fetchStatusLaporanOptions = async () => {
  statusLoading.value = true
  try {
    console.log("Fetching status laporan master data from API...")
    
    const response = await apiRequest<GetAllStatusLaporanResponse>("/admin/master-status-laporan")
    
    statusLaporanOptions.value = response.data.data
    console.log("Status laporan options loaded:", statusLaporanOptions.value.length, "items")
    
  } catch (error) {
    console.error("Error fetching status laporan options:", error)
    toast.error("Gagal memuat data status laporan. Menggunakan data default.")
    
    // Fallback to static data if API fails
    statusLaporanOptions.value = [
      {
        id: "1",
        status: "Belum Diproses",
        description: "Laporan baru masuk, menunggu verifikasi"
      },
      {
        id: "2",
        status: "Sedang Diproses",
        description: "Laporan sedang dalam tahap verifikasi dan penanganan"
      },
      {
        id: "3",
        status: "Diproses dan Data Sesuai",
        description: "Laporan telah diverifikasi dan data sesuai"
      },
      {
        id: "4",
        status: "Ditolak",
        description: "Laporan ditolak karena data tidak sesuai atau duplikasi"
      }
    ]
  } finally {
    statusLoading.value = false
  }
}

// Get selected info
const selectedPelaporInfo = computed(() => {
  if (!formData.value.id_masyarakat) return null
  return pelaporOptions.value.find(p => p.id === formData.value.id_masyarakat) || null
})

const selectedBalitaInfo = computed(() => {
  if (!formData.value.id_balita) return null
  return balitaOptions.value.find(b => b.id === formData.value.id_balita) || null
})

const selectedStatusInfo = computed(() => {
  if (!formData.value.id_status_laporan) return null
  return statusLaporanOptions.value.find(s => s.id === formData.value.id_status_laporan) || null
})

// Check if selected pelapor is admin
const isAdminPelapor = computed(() => {
  return selectedPelaporInfo.value?.id === "ADMIN"
})

const getStatusColor = (status: string): string => {
  switch (status) {
    case "Belum Diproses":
      return "text-gray-600"
    case "Sedang Diproses":
      return "text-blue-600"
    case "Diproses dan Data Sesuai":
      return "text-green-600"
    case "Ditolak":
      return "text-red-600"
    default:
      return "text-gray-600"
  }
}

const validatePhoneNumber = (phone: string): boolean => {
  const phoneRegex = /^(\+62|62|0)8[1-9][0-9]{6,9}$/
  return phoneRegex.test(phone)
}

// Handle date selection
const handleDateSelect = (date: DateValue | undefined) => {
  selectedDate.value = date
  if (date) {
    formData.value.tanggal_laporan = format(
      new Date(date.year, date.month - 1, date.day),
      "yyyy-MM-dd"
    )
    delete errors.value.tanggal_laporan
  }
  isCalendarOpen.value = false
}

// Auto-fill phone numbers when balita is selected
watch(
  () => formData.value.id_balita,
  (newBalitaId) => {
    if (newBalitaId && selectedBalitaInfo.value) {
      // Auto-fill nomor HP keluarga balita
      if (!formData.value.nomor_hp_keluarga_balita) {
        formData.value.nomor_hp_keluarga_balita = selectedBalitaInfo.value.nomor_hp_keluarga
      }
    }
  }
)

// Auto-fill phone numbers when pelapor is selected
watch(
  () => formData.value.id_masyarakat,
  (newPelaporId) => {
    if (newPelaporId && selectedPelaporInfo.value) {
      // Auto-fill nomor HP pelapor
      if (!formData.value.nomor_hp_pelapor) {
        formData.value.nomor_hp_pelapor = selectedPelaporInfo.value.nomor_hp
      }
      
      // Auto-fill hubungan jika admin
      if (selectedPelaporInfo.value.isAdmin && !formData.value.hubungan_dengan_balita) {
        formData.value.hubungan_dengan_balita = "Petugas Administrasi"
      }
    }
  }
)

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Pelapor validation
  if (!formData.value.id_masyarakat) {
    errors.value.id_masyarakat = "Pelapor harus dipilih"
  }

  // Balita validation
  if (!formData.value.id_balita) {
    errors.value.id_balita = "Balita harus dipilih"
  }

  // Status laporan validation
  if (!formData.value.id_status_laporan) {
    errors.value.id_status_laporan = "Status laporan harus dipilih"
  }

  // Hubungan validation
  if (!formData.value.hubungan_dengan_balita) {
    errors.value.hubungan_dengan_balita = "Hubungan dengan balita harus diisi"
  } else if (formData.value.hubungan_dengan_balita.trim().length < 2) {
    errors.value.hubungan_dengan_balita = "Hubungan dengan balita minimal 2 karakter"
  } else if (formData.value.hubungan_dengan_balita.length > 50) {
    errors.value.hubungan_dengan_balita = "Hubungan dengan balita maksimal 50 karakter"
  }

  // Tanggal laporan validation
  if (!formData.value.tanggal_laporan) {
    errors.value.tanggal_laporan = "Tanggal laporan harus diisi"
  } else {
    const today = new Date()
    const laporanDate = new Date(formData.value.tanggal_laporan)
    if (laporanDate > today) {
      errors.value.tanggal_laporan = "Tanggal laporan tidak boleh di masa depan"
    }
  }

  // Phone number validation
  if (!formData.value.nomor_hp_pelapor) {
    errors.value.nomor_hp_pelapor = "Nomor HP pelapor harus diisi"
  } else if (!validatePhoneNumber(formData.value.nomor_hp_pelapor)) {
    errors.value.nomor_hp_pelapor = "Format nomor HP tidak valid (contoh: 081234567890)"
  }

  if (!formData.value.nomor_hp_keluarga_balita) {
    errors.value.nomor_hp_keluarga_balita = "Nomor HP keluarga balita harus diisi"
  } else if (!validatePhoneNumber(formData.value.nomor_hp_keluarga_balita)) {
    errors.value.nomor_hp_keluarga_balita = "Format nomor HP tidak valid (contoh: 081234567890)"
  }

  return Object.keys(errors.value).length === 0
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

  // Prepare data for submission
  const submissionData = { ...formData.value }

  // Handle admin pelapor - set id_masyarakat to "ADMIN" for processing
  if (isAdminPelapor.value) {
    submissionData.id_masyarakat = "ADMIN" // This will be converted to null in the parent component
  }

  // Auto-populate data from selected options
  if (selectedPelaporInfo.value && !isAdminPelapor.value) {
    submissionData.nama_pelapor = selectedPelaporInfo.value.nama
    submissionData.email_pelapor = selectedPelaporInfo.value.email
  } else if (isAdminPelapor.value) {
    submissionData.nama_pelapor = "Admin Sistem"
    submissionData.jenis_laporan = "admin"
  } else {
    submissionData.jenis_laporan = "masyarakat"
  }

  if (selectedBalitaInfo.value) {
    submissionData.nama_balita = selectedBalitaInfo.value.nama
    submissionData.nama_ayah = selectedBalitaInfo.value.nama_ayah
    submissionData.nama_ibu = selectedBalitaInfo.value.nama_ibu
    submissionData.nomor_kk = selectedBalitaInfo.value.nomor_kk
    submissionData.alamat = selectedBalitaInfo.value.alamat
    submissionData.kelurahan = selectedBalitaInfo.value.kelurahan
    submissionData.kecamatan = selectedBalitaInfo.value.kecamatan
  }

  if (selectedStatusInfo.value) {
    submissionData.status_laporan = selectedStatusInfo.value.status
  }

  emit("save", submissionData as LaporanMasyarakat)
  resetForm()
}

// Reset form
const resetForm = () => {
  formData.value = {
    id_masyarakat: "",
    id_balita: "",
    id_status_laporan: "",
    hubungan_dengan_balita: "",
    nomor_hp_pelapor: "",
    nomor_hp_keluarga_balita: "",
    tanggal_laporan: "",
  }
  selectedDate.value = undefined
  errors.value = {}
}

// Load form data for edit mode
const loadFormData = (laporan: LaporanMasyarakat) => {
  // Handle admin laporan case
  const pelaporId = laporan.jenis_laporan === "admin" ? "ADMIN" : laporan.id_masyarakat

  formData.value = {
    id: laporan.id,
    id_masyarakat: pelaporId,
    id_balita: laporan.id_balita,
    id_status_laporan: laporan.id_status_laporan,
    hubungan_dengan_balita: laporan.hubungan_dengan_balita,
    nomor_hp_pelapor: laporan.nomor_hp_pelapor,
    nomor_hp_keluarga_balita: laporan.nomor_hp_keluarga_balita,
    tanggal_laporan: laporan.tanggal_laporan,
    // Extended fields for response
    nama_pelapor: laporan.nama_pelapor,
    email_pelapor: laporan.email_pelapor,
    nama_balita: laporan.nama_balita,
    nama_ayah: laporan.nama_ayah,
    nama_ibu: laporan.nama_ibu,
    nomor_kk: laporan.nomor_kk,
    alamat: laporan.alamat,
    kelurahan: laporan.kelurahan,
    kecamatan: laporan.kecamatan,
    status_laporan: laporan.status_laporan,
    jenis_laporan: laporan.jenis_laporan,
    created_date: laporan.created_date,
    updated_date: laporan.updated_date,
  }

  // Set date picker
  if (laporan.tanggal_laporan) {
    const date = new Date(laporan.tanggal_laporan)
    selectedDate.value = new CalendarDate(
      date.getFullYear(),
      date.getMonth() + 1,
      date.getDate()
    )
  }

  errors.value = {}
}

// Handle close
const handleClose = () => {
  resetForm()
  emit("close")
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  async (newVal) => {
    if (newVal) {
      // Load all master data
      await Promise.all([
        fetchMasyarakatOptions(),
        fetchBalitaOptions(),
        fetchStatusLaporanOptions()
      ])
      
      if (props.laporan && props.mode === "edit") {
        loadFormData(props.laporan)
      } else {
        resetForm()
        // Set default date to today for create mode
        const today = new Date()
        selectedDate.value = new CalendarDate(
          today.getFullYear(),
          today.getMonth() + 1,
          today.getDate()
        )
        formData.value.tanggal_laporan = format(today, "yyyy-MM-dd")
      }
    }
  }
)

// Watch for laporan prop changes
watch(
  () => props.laporan,
  (newLaporan) => {
    if (newLaporan && props.mode === "edit" && props.show) {
      loadFormData(newLaporan)
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.nomor_hp_pelapor,
  (newVal) => {
    if (newVal && validatePhoneNumber(newVal)) {
      delete errors.value.nomor_hp_pelapor
    }
  }
)

watch(
  () => formData.value.nomor_hp_keluarga_balita,
  (newVal) => {
    if (newVal && validatePhoneNumber(newVal)) {
      delete errors.value.nomor_hp_keluarga_balita
    }
  }
)

watch(
  () => formData.value.hubungan_dengan_balita,
  (newVal) => {
    if (newVal && newVal.trim().length >= 2 && newVal.length <= 50) {
      delete errors.value.hubungan_dengan_balita
    }
  }
)

// Load master data on component mount
onMounted(async () => {
  await Promise.all([
    fetchMasyarakatOptions(),
    fetchBalitaOptions(),
    fetchStatusLaporanOptions()
  ])
})
</script>

<template>
  <Dialog
    :open="show"
    @update:open="handleClose">
    <DialogContent class="w-[95vw] max-w-5xl max-h-[95vh] overflow-hidden p-0">
      <div class="flex flex-col h-full max-h-[95vh]">
        <!-- Header -->
        <div class="p-4 sm:p-6 border-b bg-background">
          <DialogHeader class="space-y-2">
            <DialogTitle class="text-lg sm:text-xl flex items-center gap-2">
              <FileText class="h-5 w-5 text-blue-600" />
              {{ mode === "create" ? "Tambah Laporan Masyarakat" : "Edit Laporan Masyarakat" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah laporan masyarakat baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi laporan masyarakat. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Data Pelapor -->
            <Card :class="errors.id_masyarakat ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Users class="h-4 w-4" />
                  Data Pelapor *
                  <Loader2 
                    v-if="masyarakatLoading" 
                    class="h-3 w-3 animate-spin" />
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select 
                  v-model="formData.id_masyarakat"
                  :disabled="masyarakatLoading">
                  <SelectTrigger 
                    class="w-full" 
                    :class="errors.id_masyarakat ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedPelaporInfo">
                        <div class="flex items-center gap-2">
                          <Shield v-if="selectedPelaporInfo.isAdmin" class="h-4 w-4 text-blue-600" />
                          <Users v-else class="h-4 w-4 text-gray-600" />
                          <span>{{ selectedPelaporInfo.nama }}</span>
                          <Badge v-if="selectedPelaporInfo.isAdmin" variant="outline" class="bg-blue-50 text-blue-700 border-blue-200">
                            Admin
                          </Badge>
                        </div>
                      </template>
                      <template v-else>
                        {{ masyarakatLoading ? "Memuat pelapor..." : "Pilih pelapor" }}
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="pelapor in pelaporOptions"
                      :key="pelapor.id"
                      :value="pelapor.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <Shield v-if="pelapor.isAdmin" class="h-4 w-4 text-blue-600" />
                          <Users v-else class="h-4 w-4 text-gray-600" />
                          <span class="font-medium">{{ pelapor.nama }}</span>
                          <Badge v-if="pelapor.isAdmin" variant="outline" class="bg-blue-50 text-blue-700 border-blue-200 text-xs">
                            Admin
                          </Badge>
                        </div>
                        <div class="text-xs text-muted-foreground">
                          📧 {{ pelapor.email }}
                        </div>
                        <div class="text-xs" :class="pelapor.isAdmin ? 'text-blue-600' : 'text-blue-600'">
                          {{ pelapor.isAdmin ? '🏛️' : '📍' }} {{ pelapor.alamat }}
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_masyarakat"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_masyarakat }}
                </p>
                
                <!-- Preview pelapor -->
                <div
                  v-if="selectedPelaporInfo"
                  class="mt-3 p-3 rounded-md border"
                  :class="selectedPelaporInfo.isAdmin ? 'bg-blue-50 border-blue-200' : 'bg-green-50 border-green-200'">
                  <div class="text-sm">
                    <div class="flex items-center gap-2 mb-1">
                      <Shield v-if="selectedPelaporInfo.isAdmin" class="h-4 w-4 text-blue-600" />
                      <Users v-else class="h-4 w-4 text-green-600" />
                      <span class="font-medium" :class="selectedPelaporInfo.isAdmin ? 'text-blue-900' : 'text-green-900'">
                        {{ selectedPelaporInfo.nama }}
                      </span>
                      <Badge v-if="selectedPelaporInfo.isAdmin" variant="outline" class="bg-blue-100 text-blue-800 border-blue-300 text-xs">
                        Admin Sistem
                      </Badge>
                    </div>
                    <div class="text-xs space-y-1" :class="selectedPelaporInfo.isAdmin ? 'text-blue-700' : 'text-green-700'">
                      <div>📧 {{ selectedPelaporInfo.email }}</div>
                      <div>{{ selectedPelaporInfo.isAdmin ? '🏛️' : '📍' }} {{ selectedPelaporInfo.alamat }}</div>
                    </div>
                    <div v-if="selectedPelaporInfo.isAdmin" class="mt-2 text-xs text-blue-600 bg-blue-100 rounded px-2 py-1">
                      ℹ️ Laporan ini akan tercatat sebagai laporan administratif
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Data Balita -->
            <Card :class="errors.id_balita ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Baby class="h-4 w-4" />
                  Data Balita yang Dilaporkan *
                  <Loader2 
                    v-if="balitaLoading" 
                    class="h-3 w-3 animate-spin" />
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select 
                  v-model="formData.id_balita"
                  :disabled="balitaLoading">
                  <SelectTrigger 
                    class="w-full" 
                    :class="errors.id_balita ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedBalitaInfo">
                        {{ selectedBalitaInfo.nama }} ({{ selectedBalitaInfo.umur }})
                      </template>
                      <template v-else>
                        {{ balitaLoading ? "Memuat balita..." : "Pilih balita yang dilaporkan" }}
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="balita in balitaOptions"
                      :key="balita.id"
                      :value="balita.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <span class="font-medium">{{ balita.nama }}</span>
                          <span class="text-xs text-gray-600">({{ balita.umur }}, {{ balita.jenis_kelamin }})</span>
                        </div>
                        <div class="text-xs text-muted-foreground">
                          👨‍👩‍👧‍👦 {{ balita.nama_ayah }} & {{ balita.nama_ibu }}
                        </div>
                        <div class="text-xs text-muted-foreground">
                          🏠 {{ balita.alamat }}
                        </div>
                        <div class="text-xs text-blue-600">
                          📄 KK: {{ balita.nomor_kk }}
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_balita"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_balita }}
                </p>
                <div
                  v-if="selectedBalitaInfo"
                  class="mt-3 p-3 bg-green-50 rounded-md border border-green-200">
                  <div class="text-sm">
                    <div class="font-medium text-green-900 mb-1">
                      {{ selectedBalitaInfo.nama }} ({{ selectedBalitaInfo.umur }}, {{ selectedBalitaInfo.jenis_kelamin }})
                    </div>
                    <div class="text-green-700 text-xs space-y-1">
                      <div>👨‍👩‍👧‍👦 Orang Tua: {{ selectedBalitaInfo.nama_ayah }} & {{ selectedBalitaInfo.nama_ibu }}</div>
                      <div>📞 Kontak Keluarga: {{ selectedBalitaInfo.nomor_hp_keluarga }}</div>
                      <div>🏠 Alamat: {{ selectedBalitaInfo.alamat }}</div>
                      <div>🏘️ {{ selectedBalitaInfo.kelurahan }}, {{ selectedBalitaInfo.kecamatan }}</div>
                      <div>📄 No. KK: {{ selectedBalitaInfo.nomor_kk }}</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Hubungan dengan Balita -->
            <Card :class="errors.hubungan_dengan_balita ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <UserCheck class="h-4 w-4" />
                  Hubungan dengan Balita *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-2">
                  <Label
                    for="hubungan_dengan_balita"
                    class="text-sm font-medium">
                    Hubungan dengan Balita *
                  </Label>
                  <Input
                    id="hubungan_dengan_balita"
                    v-model="formData.hubungan_dengan_balita"
                    :placeholder="isAdminPelapor ? 'Contoh: Petugas Administrasi, Koordinator Program, dll' : 'Contoh: Ibu, Ayah, Kakek, Nenek, Paman, Bibi, Tetangga, dll'"
                    :class="errors.hubungan_dengan_balita ? 'border-red-500' : ''" />
                  <p
                    v-if="errors.hubungan_dengan_balita"
                    class="text-sm text-red-500">
                    {{ errors.hubungan_dengan_balita }}
                  </p>
                  <div class="flex justify-between items-center text-xs text-muted-foreground">
                    <span>{{ isAdminPelapor ? 'Jelaskan kapasitas/peran admin dalam melaporkan' : 'Jelaskan hubungan pelapor dengan balita yang dilaporkan' }}</span>
                    <span>{{ formData.hubungan_dengan_balita?.length || 0 }}/50 karakter</span>
                  </div>
                </div>
                
                <!-- Preview hubungan -->
                <div
                  v-if="formData.hubungan_dengan_balita && formData.hubungan_dengan_balita.trim().length >= 2"
                  class="mt-3 p-3 rounded-md border"
                  :class="isAdminPelapor ? 'bg-blue-50 border-blue-200' : 'bg-purple-50 border-purple-200'">
                  <div class="flex items-center gap-2 text-sm">
                    <Shield v-if="isAdminPelapor" class="h-4 w-4 text-blue-600" />
                    <UserCheck v-else class="h-4 w-4 text-purple-600" />
                    <div>
                      <div class="font-medium" :class="isAdminPelapor ? 'text-blue-900' : 'text-purple-900'">
                        {{ formData.hubungan_dengan_balita }}
                      </div>
                      <div class="text-xs" :class="isAdminPelapor ? 'text-blue-700' : 'text-purple-700'">
                        {{ isAdminPelapor ? 'Kapasitas admin' : 'Hubungan pelapor' }} dengan {{ selectedBalitaInfo?.nama || 'balita' }}
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Contoh hubungan untuk panduan -->
                <div class="mt-3 p-3 rounded-md border" :class="isAdminPelapor ? 'bg-blue-50 border-blue-200' : 'bg-green-50 border-green-200'">
                  <div class="text-xs" :class="isAdminPelapor ? 'text-blue-800' : 'text-green-800'">
                    <div class="font-medium mb-1">💡 {{ isAdminPelapor ? 'Contoh kapasitas admin:' : 'Contoh hubungan yang umum:' }}</div>
                    <div class="grid grid-cols-2 gap-1">
                      <template v-if="isAdminPelapor">
                        <span>• Petugas Administrasi</span>
                        <span>• Koordinator Program</span>
                        <span>• Supervisor Lapangan</span>
                        <span>• Verifikator Data</span>
                        <span>• Petugas Monitoring</span>
                        <span>• Admin Sistem</span>
                      </template>
                      <template v-else>
                        <span>• Ibu kandung</span>
                        <span>• Ayah kandung</span>
                        <span>• Kakek/Nenek</span>
                        <span>• Paman/Bibi</span>
                        <span>• Tetangga</span>
                        <span>• Kader Posyandu</span>
                        <span>• Keluarga lainnya</span>
                        <span>• Petugas Kesehatan</span>
                      </template>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Kontak Informasi -->
            <Card>
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Phone class="h-4 w-4" />
                  Informasi Kontak
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <!-- Nomor HP Pelapor -->
                  <div class="space-y-2">
                    <Label
                      for="nomor_hp_pelapor"
                      class="text-sm font-medium">
                      Nomor HP Pelapor *
                    </Label>
                    <Input
                      id="nomor_hp_pelapor"
                      v-model="formData.nomor_hp_pelapor"
                      placeholder="081234567890"
                      type="tel"
                      :class="errors.nomor_hp_pelapor ? 'border-red-500' : ''" />
                    <p
                      v-if="errors.nomor_hp_pelapor"
                      class="text-sm text-red-500">
                      {{ errors.nomor_hp_pelapor }}
                    </p>
                    <div class="text-xs text-muted-foreground">
                      Format: 08xxxxxxxxxx (10-13 digit)
                    </div>
                  </div>

                  <!-- Nomor HP Keluarga Balita -->
                  <div class="space-y-2">
                    <Label
                      for="nomor_hp_keluarga_balita"
                      class="text-sm font-medium">
                      Nomor HP Keluarga Balita *
                    </Label>
                    <Input
                      id="nomor_hp_keluarga_balita"
                      v-model="formData.nomor_hp_keluarga_balita"
                      placeholder="081234567890"
                      type="tel"
                      :class="errors.nomor_hp_keluarga_balita ? 'border-red-500' : ''" />
                    <p
                      v-if="errors.nomor_hp_keluarga_balita"
                      class="text-sm text-red-500">
                      {{ errors.nomor_hp_keluarga_balita }}
                    </p>
                    <div class="text-xs text-muted-foreground">
                      Format: 08xxxxxxxxxx (10-13 digit)
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Tanggal dan Status -->
            <div class="grid grid-cols-1 gap-6">
              <!-- Tanggal Laporan -->
              <Card :class="errors.tanggal_laporan ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="flex items-center gap-2 text-base">
                    <CalendarIcon class="h-4 w-4" />
                    Tanggal Laporan *
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Popover
                    :open="isCalendarOpen"
                    @update:open="isCalendarOpen = $event">
                    <PopoverTrigger as-child>
                      <Button
                        variant="outline"
                        role="combobox"
                        :class=" [
                          'w-full justify-start text-left font-normal',
                          !selectedDate && 'text-muted-foreground',
                          errors.tanggal_laporan && 'border-red-500',
                        ]">
                        <CalendarIcon class="mr-2 h-4 w-4" />
                        {{
                          selectedDate
                            ? format(
                                new Date(selectedDate.year, selectedDate.month - 1, selectedDate.day),
                                "PPP",
                                { locale: id }
                              )
                            : "Pilih tanggal laporan"
                        }}
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent
                      class="w-auto p-0"
                      align="start">
                      <Calendar
                        :model-value="selectedDate"
                        @update:model-value="handleDateSelect"
                        initial-focus />
                    </PopoverContent>
                  </Popover>
                  <p
                    v-if="errors.tanggal_laporan"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.tanggal_laporan }}
                  </p>
                </CardContent>
              </Card>

              <!-- Status Laporan -->
              <Card :class="errors.id_status_laporan ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="flex items-center gap-2 text-base">
                    <Building class="h-4 w-4" />
                    Status Laporan *
                    <Loader2 
                      v-if="statusLoading" 
                      class="h-3 w-3 animate-spin" />
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Select 
                    v-model="formData.id_status_laporan"
                    :disabled="statusLoading">
                    <SelectTrigger 
                      class="w-full" 
                      :class="errors.id_status_laporan ? 'border-red-500' : ''">
                      <SelectValue>
                        <template v-if="formData.id_status_laporan">
                          {{ statusLaporanOptions.find(status => status.id === formData.id_status_laporan)?.status }}
                        </template>
                        <template v-else>
                          {{ statusLoading ? "Memuat status..." : "Pilih status laporan" }}
                        </template>
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="status in statusLaporanOptions"
                        :key="status.id"
                        :value="status.id">
                        <div class="flex flex-col gap-1 py-1">
                          <div class="font-medium" :class="getStatusColor(status.status)">{{ status.status }}</div>
                          <div class="text-xs text-muted-foreground">{{ status.description }}</div>
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p
                    v-if="errors.id_status_laporan"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.id_status_laporan }}
                  </p>
                  <div
                    v-if="selectedStatusInfo"
                    class="mt-3 p-3 rounded-md border bg-gray-50">
                    <div class="text-sm">
                      <div class="font-medium" :class="getStatusColor(selectedStatusInfo.status)">{{ selectedStatusInfo.status }}</div>
                      <div class="text-xs mt-1 opacity-80">{{ selectedStatusInfo.description }}</div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>

            <!-- Info Box -->
            <div class="bg-orange-50 border border-orange-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <AlertTriangle class="h-5 w-5 text-orange-600 mt-0.5 flex-shrink-0" />
                <div class="text-sm text-orange-800">
                  <p class="font-medium mb-2">📋 Panduan Pelaporan:</p>
                  <ul class="space-y-1 text-xs">
                    <li>• <strong>Pelapor:</strong> Pilih admin saat ini atau masyarakat yang melaporkan kondisi balita</li>
                    <li>• <strong>Admin:</strong> Laporan admin akan tercatat sebagai laporan administratif (id_masyarakat = NULL)</li>
                    <li>• <strong>Balita:</strong> Pilih balita yang dilaporkan memiliki masalah gizi</li>
                    <li>• <strong>Hubungan:</strong> Jelaskan hubungan pelapor dengan balita atau kapasitas admin</li>
                    <li>• <strong>Kontak:</strong> Pastikan nomor HP aktif untuk follow-up</li>
                    <li>• <strong>Status:</strong> Pilih status sesuai tahap penanganan laporan</li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Live Preview -->
            <div
              v-if="selectedPelaporInfo && selectedBalitaInfo && formData.hubungan_dengan_balita"
              class="bg-green-50 border border-green-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <div class="text-2xl">✅</div>
                <div class="flex-1">
                  <div class="text-sm font-medium text-green-900 mb-2">Preview Laporan:</div>
                  <div class="space-y-2 text-xs text-green-800">
                    <div>
                      <strong>Pelapor:</strong> {{ selectedPelaporInfo.nama }}
                      <Badge v-if="selectedPelaporInfo.isAdmin" variant="outline" class="ml-2 bg-blue-50 text-blue-700 border-blue-200 text-xs">
                        Admin
                      </Badge>
                    </div>
                    <div>
                      <strong>Kapasitas:</strong> {{ formData.hubungan_dengan_balita }} {{ selectedPelaporInfo.isAdmin ? '(Administratif)' : `dari ${selectedBalitaInfo.nama}` }}
                    </div>
                    <div>
                      <strong>Balita:</strong> {{ selectedBalitaInfo.nama }} ({{ selectedBalitaInfo.umur }})
                    </div>
                    <div>
                      <strong>Orang Tua:</strong> {{ selectedBalitaInfo.nama_ayah }} & {{ selectedBalitaInfo.nama_ibu }}
                    </div>
                    <div>
                      <strong>Status:</strong> {{ selectedStatusInfo?.status || "Belum dipilih" }}
                    </div>
                    <div v-if="formData.tanggal_laporan">
                      <strong>Tanggal:</strong> {{ format(new Date(formData.tanggal_laporan), "PPP", { locale: id }) }}
                    </div>
                    <div v-if="selectedPelaporInfo.isAdmin" class="text-blue-600 bg-blue-100 rounded px-2 py-1 mt-2">
                      📄 Laporan Administratif: id_masyarakat akan diset NULL di database
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Debug Info (untuk development) -->
            <div 
              class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Mode: {{ mode }}<br />
              Has Data: {{ !!props.laporan }}<br />
              Form ID: {{ formData.id || "New" }}<br />
              Pelapor: {{ selectedPelaporInfo?.nama || "Not selected" }}<br />
              Pelapor Type: {{ selectedPelaporInfo?.type || "None" }}<br />
              Is Admin: {{ isAdminPelapor }}<br />
              Balita: {{ selectedBalitaInfo?.nama || "Not selected" }}<br />
              Hubungan: {{ formData.hubungan_dengan_balita || "Not entered" }}<br />
              Status: {{ selectedStatusInfo?.status || "Not selected" }}<br />
              Tanggal: {{ formData.tanggal_laporan || "Not selected" }}<br />
              Will set id_masyarakat to NULL: {{ isAdminPelapor }}<br />
              API Loading: Masyarakat={{ masyarakatLoading }}, Balita={{ balitaLoading }}, Status={{ statusLoading }}
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
              :disabled="loading || masyarakatLoading || balitaLoading || statusLoading"
              class="w-full sm:w-auto">
              <Loader2 
                v-if="loading" 
                class="h-4 w-4 mr-2 animate-spin" />
              <Save 
                v-else 
                class="h-4 w-4 mr-2" />
              {{ loading ? "Menyimpan..." : (mode === "create" ? "Simpan" : "Perbarui") }}
            </Button>
          </DialogFooter>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>