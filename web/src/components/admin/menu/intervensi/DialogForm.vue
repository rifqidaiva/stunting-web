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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { Textarea } from "@/components/ui/textarea"
import {
  CalendarIcon,
  Save,
  X,
  Activity,
  Baby,
  FileText,
  Stethoscope,
  Loader2,
} from "lucide-vue-next"
import { format } from "date-fns"
import { id } from "date-fns/locale"
import { CalendarDate, type DateValue } from "@internationalized/date"
import { authUtils } from "@/lib/utils"
import type { Intervensi } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  intervensi: Intervensi | null
  loading?: boolean
}

interface Emits {
  (e: "close"): void
  (e: "save", intervensi: Intervensi): void
}

// API Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface BalitaOption {
  id: string
  nama: string
  umur: string
  jenis_kelamin: string
  nama_ayah: string
  nama_ibu: string
  nomor_kk: string
  kelurahan: string
  kecamatan: string
}

interface GetAllBalitaResponse {
  data: BalitaOption[]
  total: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data sesuai dengan endpoint
const formData = ref<Partial<Intervensi>>({
  id_balita: "",
  jenis: "",
  tanggal: "",
  deskripsi: "",
  hasil: "",
})

// Date picker state
const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

// Validation errors
const errors = ref<Record<string, string>>({})

// Master data from API
const balitaOptions = ref<BalitaOption[]>([])
const balitaLoading = ref(false)

// Static jenis intervensi options
const jenisIntervensiOptions = [
  {
    value: "gizi",
    label: "Gizi",
    description: "Intervensi terkait nutrisi dan gizi balita",
    icon: "🥗",
    color: "text-green-600",
  },
  {
    value: "kesehatan",
    label: "Kesehatan",
    description: "Pemeriksaan dan pelayanan kesehatan",
    icon: "🏥",
    color: "text-red-600",
  },
  {
    value: "sosial",
    label: "Sosial",
    description: "Dukungan sosial dan konseling keluarga",
    icon: "👥",
    color: "text-blue-600",
  },
]

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
      Authorization: `Bearer ${token}`,
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

// Fetch balita master data
const fetchBalitaOptions = async () => {
  balitaLoading.value = true
  try {
    console.log("Fetching balita master data from API...")

    const response = await apiRequest<GetAllBalitaResponse>("/admin/balita/get")

    balitaOptions.value = response.data.data.map((b) => ({
      id: b.id,
      nama: b.nama,
      umur: b.umur,
      jenis_kelamin: b.jenis_kelamin,
      nama_ayah: b.nama_ayah,
      nama_ibu: b.nama_ibu,
      nomor_kk: b.nomor_kk,
      kelurahan: b.kelurahan,
      kecamatan: b.kecamatan,
    }))

    console.log("Balita options loaded:", balitaOptions.value.length, "items")
  } catch (error) {
    console.error("Error fetching balita options:", error)
    toast.error("Gagal memuat data balita. Menggunakan data default.")

    // Fallback to static data if API fails
    balitaOptions.value = [
      {
        id: "B001",
        nama: "Ahmad Rizki",
        umur: "36 bulan",
        jenis_kelamin: "L",
        nama_ayah: "Budi Santoso",
        nama_ibu: "Siti Nurhaliza",
        nomor_kk: "3201234567890123",
        kelurahan: "Kejaksan",
        kecamatan: "Kejaksan",
      },
      {
        id: "B002",
        nama: "Fatimah Zahra",
        umur: "24 bulan",
        jenis_kelamin: "P",
        nama_ayah: "Ahmad Rahman",
        nama_ibu: "Dewi Sartika",
        nomor_kk: "3201234567890124",
        kelurahan: "Pekalangan",
        kecamatan: "Pekalangan",
      },
    ]
  } finally {
    balitaLoading.value = false
  }
}

// Get selected info
const selectedBalitaInfo = computed(() => {
  if (!formData.value.id_balita) return null
  return balitaOptions.value.find((b) => b.id === formData.value.id_balita) || null
})

const selectedJenisInfo = computed(() => {
  if (!formData.value.jenis) return null
  return jenisIntervensiOptions.find((j) => j.value === formData.value.jenis) || null
})

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Balita validation
  if (!formData.value.id_balita) {
    errors.value.id_balita = "Balita harus dipilih"
  }

  // Jenis validation
  if (!formData.value.jenis) {
    errors.value.jenis = "Jenis intervensi harus dipilih"
  }

  // Tanggal validation
  if (!formData.value.tanggal) {
    errors.value.tanggal = "Tanggal intervensi harus diisi"
  } else {
    const intervensiDate = new Date(formData.value.tanggal)
    const today = new Date()

    if (intervensiDate > today) {
      errors.value.tanggal = "Tanggal intervensi tidak boleh di masa depan"
    }
  }

  // Deskripsi validation
  if (!formData.value.deskripsi || formData.value.deskripsi.trim().length < 10) {
    errors.value.deskripsi = "Deskripsi minimal 10 karakter"
  } else if (formData.value.deskripsi.length > 1000) {
    errors.value.deskripsi = "Deskripsi maksimal 1000 karakter"
  }

  // Hasil validation
  if (!formData.value.hasil || formData.value.hasil.trim().length < 5) {
    errors.value.hasil = "Hasil intervensi minimal 5 karakter"
  } else if (formData.value.hasil.length > 500) {
    errors.value.hasil = "Hasil intervensi maksimal 500 karakter"
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    id_balita: "",
    jenis: "",
    tanggal: "",
    deskripsi: "",
    hasil: "",
  }
  errors.value = {}
  selectedDate.value = undefined
}

// Load form data for edit mode
const loadFormData = (intervensi: Intervensi) => {
  formData.value = {
    id: intervensi.id,
    id_balita: intervensi.id_balita,
    jenis: intervensi.jenis,
    tanggal: intervensi.tanggal,
    deskripsi: intervensi.deskripsi,
    hasil: intervensi.hasil,
    // Extended fields for response
    nama_balita: intervensi.nama_balita,
    petugas_count: intervensi.petugas_count,
    riwayat_count: intervensi.riwayat_count,
    created_date: intervensi.created_date,
    updated_date: intervensi.updated_date,
    created_by: intervensi.created_by,
    updated_by: intervensi.updated_by,
  }
  errors.value = {}

  // Initialize date for calendar
  if (intervensi.tanggal) {
    const jsDate = new Date(intervensi.tanggal)
    selectedDate.value = new CalendarDate(
      jsDate.getFullYear(),
      jsDate.getMonth() + 1,
      jsDate.getDate()
    )
  }
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

  // Auto-populate data from selected options
  if (selectedBalitaInfo.value) {
    formData.value.nama_balita = selectedBalitaInfo.value.nama
  }

  emit("save", formData.value as Intervensi)
  resetForm()
}

// Handle close
const handleClose = () => {
  resetForm()
  emit("close")
}

// Calendar handlers
const handleDateSelect = (date: DateValue | undefined) => {
  if (date) {
    selectedDate.value = date
    const jsDate = new Date(date.year, date.month - 1, date.day)
    formData.value.tanggal = format(jsDate, "yyyy-MM-dd")
    isCalendarOpen.value = false

    // Clear tanggal error if date is valid
    const intervensiDate = jsDate
    const today = new Date()

    if (intervensiDate <= today) {
      delete errors.value.tanggal
    }
  }
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  async (newVal) => {
    if (newVal) {
      // Load balita options first
      await fetchBalitaOptions()

      if (props.intervensi && props.mode === "edit") {
        loadFormData(props.intervensi)
      } else {
        resetForm()
        // Set default tanggal to today for create mode
        const today = new Date()
        selectedDate.value = new CalendarDate(
          today.getFullYear(),
          today.getMonth() + 1,
          today.getDate()
        )
        formData.value.tanggal = format(today, "yyyy-MM-dd")
      }
    }
  }
)

// Watch for intervensi prop changes
watch(
  () => props.intervensi,
  (newIntervensi) => {
    if (newIntervensi && props.mode === "edit" && props.show) {
      loadFormData(newIntervensi)
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.deskripsi,
  (newVal) => {
    if (newVal && newVal.trim().length >= 10 && newVal.length <= 1000) {
      delete errors.value.deskripsi
    }
  }
)

watch(
  () => formData.value.hasil,
  (newVal) => {
    if (newVal && newVal.trim().length >= 5 && newVal.length <= 500) {
      delete errors.value.hasil
    }
  }
)

// Format tanggal untuk display
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  })
}

// Load balita options on component mount
onMounted(async () => {
  await fetchBalitaOptions()
})
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
              <Activity class="h-5 w-5 text-blue-600" />
              {{ mode === "create" ? "Tambah Intervensi Baru" : "Edit Data Intervensi" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data intervensi baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi intervensi. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Pilih Balita -->
            <Card :class="errors.id_balita ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Baby class="h-4 w-4" />
                  Pilih Balita *
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
                        {{
                          balitaLoading ? "Memuat balita..." : "Pilih balita yang akan diintervensi"
                        }}
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
                          <span class="text-xs text-gray-600"
                            >({{ balita.umur }}, {{ balita.jenis_kelamin }})</span
                          >
                        </div>
                        <div class="text-xs text-muted-foreground">
                          👨‍👩‍👧‍👦 {{ balita.nama_ayah }} & {{ balita.nama_ibu }}
                        </div>
                        <div class="text-xs text-muted-foreground">
                          🏠 {{ balita.kelurahan }}, {{ balita.kecamatan }}
                        </div>
                        <div class="text-xs text-blue-600">📄 KK: {{ balita.nomor_kk }}</div>
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
                  class="mt-3 p-3 bg-blue-50 rounded-md border border-blue-200">
                  <div class="text-sm">
                    <div class="font-medium text-blue-900 mb-1">
                      {{ selectedBalitaInfo.nama }} ({{ selectedBalitaInfo.umur }},
                      {{ selectedBalitaInfo.jenis_kelamin }})
                    </div>
                    <div class="text-blue-700 text-xs space-y-1">
                      <div>
                        👨‍👩‍👧‍👦 Orang Tua: {{ selectedBalitaInfo.nama_ayah }} &
                        {{ selectedBalitaInfo.nama_ibu }}
                      </div>
                      <div>
                        🏠 Alamat: {{ selectedBalitaInfo.kelurahan }},
                        {{ selectedBalitaInfo.kecamatan }}
                      </div>
                      <div>📄 No. KK: {{ selectedBalitaInfo.nomor_kk }}</div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Jenis Intervensi & Tanggal -->
            <div class="grid grid-cols-1 gap-6">
              <!-- Jenis Intervensi -->
              <Card :class="errors.jenis ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="flex items-center gap-2 text-base">
                    <Stethoscope class="h-4 w-4" />
                    Jenis Intervensi *
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <Select v-model="formData.jenis">
                    <SelectTrigger
                      class="w-full"
                      :class="errors.jenis ? 'border-red-500' : ''">
                      <SelectValue>
                        <template v-if="selectedJenisInfo">
                          <div class="flex items-center gap-2">
                            <span class="text-lg">{{ selectedJenisInfo.icon }}</span>
                            <span>{{ selectedJenisInfo.label }}</span>
                          </div>
                        </template>
                        <template v-else> Pilih jenis intervensi </template>
                      </SelectValue>
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="jenis in jenisIntervensiOptions"
                        :key="jenis.value"
                        :value="jenis.value">
                        <div class="flex items-center gap-3 py-1">
                          <span class="text-lg">{{ jenis.icon }}</span>
                          <div class="flex flex-col">
                            <span
                              class="font-medium"
                              :class="jenis.color"
                              >{{ jenis.label }}</span
                            >
                            <span class="text-xs text-muted-foreground">{{
                              jenis.description
                            }}</span>
                          </div>
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p
                    v-if="errors.jenis"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.jenis }}
                  </p>
                  <div
                    v-if="selectedJenisInfo"
                    class="mt-3 p-3 rounded-md border"
                    :class="
                      selectedJenisInfo.value === 'gizi'
                        ? 'bg-green-50 border-green-200'
                        : selectedJenisInfo.value === 'kesehatan'
                        ? 'bg-red-50 border-red-200'
                        : 'bg-blue-50 border-blue-200'
                    ">
                    <div class="flex items-center gap-2">
                      <span class="text-lg">{{ selectedJenisInfo.icon }}</span>
                      <div>
                        <div
                          class="text-sm font-medium"
                          :class="selectedJenisInfo.color">
                          {{ selectedJenisInfo.label }}
                        </div>
                        <div class="text-xs opacity-80">{{ selectedJenisInfo.description }}</div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <!-- Tanggal Intervensi -->
              <Card :class="errors.tanggal ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="flex items-center gap-2 text-base">
                    <CalendarIcon class="h-4 w-4" />
                    Tanggal Intervensi *
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
                        :class="[
                          'w-full justify-start text-left font-normal',
                          !selectedDate && 'text-muted-foreground',
                          errors.tanggal && 'border-red-500',
                        ]">
                        <CalendarIcon class="mr-2 h-4 w-4" />
                        {{
                          selectedDate
                            ? format(
                                new Date(
                                  selectedDate.year,
                                  selectedDate.month - 1,
                                  selectedDate.day
                                ),
                                "PPP",
                                { locale: id }
                              )
                            : "Pilih tanggal intervensi"
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
                    v-if="errors.tanggal"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.tanggal }}
                  </p>
                  <div
                    v-if="formData.tanggal"
                    class="mt-2 p-2 bg-blue-50 rounded-md">
                    <p class="text-xs text-blue-800">
                      <strong>Tanggal yang dipilih:</strong> {{ formatDate(formData.tanggal) }}
                    </p>
                  </div>
                </CardContent>
              </Card>
            </div>

            <!-- Deskripsi Intervensi -->
            <Card :class="errors.deskripsi ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <FileText class="h-4 w-4" />
                  Deskripsi Intervensi *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Textarea
                  v-model="formData.deskripsi"
                  placeholder="Jelaskan secara detail intervensi yang akan/telah dilakukan, termasuk metode, alat yang digunakan, prosedur, dan target yang ingin dicapai..."
                  rows="4"
                  :class="errors.deskripsi ? 'border-red-500' : ''" />
                <p
                  v-if="errors.deskripsi"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.deskripsi }}
                </p>
                <div class="flex justify-between items-center mt-2">
                  <p class="text-xs text-muted-foreground">Minimum 10 karakter</p>
                  <p class="text-xs text-muted-foreground">
                    {{ formData.deskripsi?.length || 0 }}/1000 karakter
                  </p>
                </div>
              </CardContent>
            </Card>

            <!-- Hasil Intervensi -->
            <Card :class="errors.hasil ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <FileText class="h-4 w-4" />
                  Hasil Intervensi *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Textarea
                  v-model="formData.hasil"
                  placeholder="Jelaskan hasil dari intervensi yang telah dilakukan, dampak yang terlihat, respons balita, dan evaluasi keberhasilan intervensi..."
                  rows="3"
                  :class="errors.hasil ? 'border-red-500' : ''" />
                <p
                  v-if="errors.hasil"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.hasil }}
                </p>
                <div class="flex justify-between items-center mt-2">
                  <p class="text-xs text-muted-foreground">Minimum 5 karakter</p>
                  <p class="text-xs text-muted-foreground">
                    {{ formData.hasil?.length || 0 }}/500 karakter
                  </p>
                </div>
              </CardContent>
            </Card>

            <!-- Info Box -->
            <div class="bg-orange-50 border border-orange-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <Activity class="h-5 w-5 text-orange-600 mt-0.5 flex-shrink-0" />
                <div class="text-sm text-orange-800">
                  <p class="font-medium mb-2">📋 Panduan Pengisian Intervensi:</p>
                  <ul class="space-y-1 text-xs">
                    <li>
                      • <strong>Balita:</strong> Pastikan memilih balita yang tepat berdasarkan
                      status gizi
                    </li>
                    <li>
                      • <strong>Jenis:</strong> Pilih jenis intervensi sesuai dengan kebutuhan
                      balita
                    </li>
                    <li>
                      • <strong>Tanggal:</strong> Tanggal pelaksanaan intervensi (tidak boleh di
                      masa depan)
                    </li>
                    <li>
                      • <strong>Deskripsi:</strong> Jelaskan detail metode dan prosedur intervensi
                      (10-1000 karakter)
                    </li>
                    <li>
                      • <strong>Hasil:</strong> Evaluasi dampak dan keberhasilan intervensi (5-500
                      karakter)
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Live Preview -->
            <div
              v-if="selectedBalitaInfo && selectedJenisInfo && formData.tanggal"
              class="bg-green-50 border border-green-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <div class="text-2xl">✅</div>
                <div class="flex-1">
                  <div class="text-sm font-medium text-green-900 mb-2">Preview Intervensi:</div>
                  <div class="space-y-2 text-xs text-green-800">
                    <div>
                      <strong>Balita:</strong> {{ selectedBalitaInfo.nama }} ({{
                        selectedBalitaInfo.umur
                      }})
                    </div>
                    <div>
                      <strong>Jenis:</strong> {{ selectedJenisInfo.icon }}
                      {{ selectedJenisInfo.label }}
                    </div>
                    <div><strong>Tanggal:</strong> {{ formatDate(formData.tanggal) }}</div>
                    <div v-if="formData.deskripsi">
                      <strong>Deskripsi:</strong> {{ formData.deskripsi.substring(0, 100)
                      }}{{ formData.deskripsi.length > 100 ? "..." : "" }}
                    </div>
                    <div v-if="formData.hasil">
                      <strong>Hasil:</strong> {{ formData.hasil.substring(0, 100)
                      }}{{ formData.hasil.length > 100 ? "..." : "" }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Debug Info (untuk development) -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Mode: {{ mode }}<br />
              Has Data: {{ !!props.intervensi }}<br />
              Form ID: {{ formData.id || "New" }}<br />
              Selected Balita: {{ selectedBalitaInfo?.nama || "None" }}<br />
              Selected Jenis: {{ selectedJenisInfo?.label || "None" }}<br />
              Tanggal: {{ formData.tanggal || "Not set" }}<br />
              Balita Options: {{ balitaOptions.length }} loaded<br />
              Loading: {{ balitaLoading }}
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
              :disabled="loading || balitaLoading"
              class="w-full sm:w-auto">
              <Loader2
                v-if="loading"
                class="h-4 w-4 mr-2 animate-spin" />
              <Save
                v-else
                class="h-4 w-4 mr-2" />
              {{ loading ? "Menyimpan..." : mode === "create" ? "Simpan" : "Perbarui" }}
            </Button>
          </DialogFooter>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
