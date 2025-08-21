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
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { CalendarIcon, Save, X, Baby, Users, Loader2 } from "lucide-vue-next"
import { format } from "date-fns"
import { id } from "date-fns/locale"
import { CalendarDate, type DateValue } from "@internationalized/date"
import { authUtils } from "@/lib/utils"
import type { Balita } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  balita: Balita | null
  loading?: boolean
}

interface Emits {
  (e: "close"): void
  (e: "save", balita: Balita): void
}

// API Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface KeluargaOption {
  id: string
  nomor_kk: string
  nama_ayah: string
  nama_ibu: string
  kelurahan: string
  kecamatan: string
}

interface GetAllKeluargaResponse {
  data: KeluargaOption[]
  total: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data
const formData = ref<Partial<Balita>>({
  id_keluarga: "",
  nama: "",
  tanggal_lahir: "",
  jenis_kelamin: "L",
  berat_lahir: "",
  tinggi_lahir: "",
})

// Date picker state
const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

// Validation errors
const errors = ref<Record<string, string>>({})

// Keluarga options from API
const keluargaOptions = ref<KeluargaOption[]>([])
const keluargaLoading = ref(false)

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

// Fetch keluarga master data
const fetchKeluargaOptions = async () => {
  keluargaLoading.value = true
  try {
    console.log("Fetching keluarga master data from API...")

    const response = await apiRequest<GetAllKeluargaResponse>("/admin/keluarga/get")

    keluargaOptions.value = response.data.data.map((k) => ({
      id: k.id,
      nomor_kk: k.nomor_kk,
      nama_ayah: k.nama_ayah,
      nama_ibu: k.nama_ibu,
      kelurahan: k.kelurahan,
      kecamatan: k.kecamatan,
    }))

    console.log("Keluarga options loaded:", keluargaOptions.value.length, "items")
  } catch (error) {
    console.error("Error fetching keluarga options:", error)
    toast.error("Gagal memuat data keluarga. Menggunakan data default.")

    // Fallback to empty array if API fails
    keluargaOptions.value = []
  } finally {
    keluargaLoading.value = false
  }
}

// Get formatted keluarga label for display
const getKeluargaLabel = (keluarga: KeluargaOption): string => {
  return `${keluarga.nomor_kk} - ${keluarga.nama_ayah} & ${keluarga.nama_ibu}`
}

// Get selected keluarga object
const getSelectedKeluarga = (id: string): KeluargaOption | undefined => {
  return keluargaOptions.value.find((k) => k.id === id)
}

// Helper functions untuk input restriction
const restrictToNumbers = (value: string) => value.replace(/\D/g, "")
const restrictToLetters = (value: string) => value.replace(/[^a-zA-Z\s]/g, "")

// Helper function untuk calculate umur
const calculateAge = (birthDate: string): string => {
  if (!birthDate) return ""

  const birth = new Date(birthDate)
  const now = new Date()

  let years = now.getFullYear() - birth.getFullYear()
  let months = now.getMonth() - birth.getMonth()

  if (now.getDate() < birth.getDate()) {
    months--
  }

  const totalMonths = years * 12 + months

  if (totalMonths < 0) return "0 bulan"

  if (totalMonths < 12) {
    return `${totalMonths} bulan`
  } else {
    const remainingMonths = totalMonths % 12
    if (remainingMonths === 0) {
      return `${years} tahun`
    }
    return `${years} tahun ${remainingMonths} bulan`
  }
}

// Computed untuk umur
const currentAge = computed(() => {
  return formData.value.tanggal_lahir ? calculateAge(formData.value.tanggal_lahir) : ""
})

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Keluarga validation
  if (!formData.value.id_keluarga) {
    errors.value.id_keluarga = "Keluarga harus dipilih"
  }

  // Nama validation
  if (!formData.value.nama) {
    errors.value.nama = "Nama balita harus diisi"
  } else if (formData.value.nama.length < 2) {
    errors.value.nama = "Nama balita minimal 2 karakter"
  } else if (formData.value.nama.length > 50) {
    errors.value.nama = "Nama balita maksimal 50 karakter"
  } else if (!/^[a-zA-Z\s]+$/.test(formData.value.nama)) {
    errors.value.nama = "Nama balita hanya boleh huruf dan spasi"
  }

  // Tanggal Lahir validation
  if (!formData.value.tanggal_lahir) {
    errors.value.tanggal_lahir = "Tanggal lahir harus diisi"
  } else {
    const birthDate = new Date(formData.value.tanggal_lahir)
    const today = new Date()
    const fiveYearsAgo = new Date()
    fiveYearsAgo.setFullYear(today.getFullYear() - 5)

    if (birthDate > today) {
      errors.value.tanggal_lahir = "Tanggal lahir tidak boleh di masa depan"
    } else if (birthDate < fiveYearsAgo) {
      errors.value.tanggal_lahir = "Anak harus berusia di bawah 5 tahun (kriteria balita)"
    }
  }

  // Jenis Kelamin validation
  if (!formData.value.jenis_kelamin) {
    errors.value.jenis_kelamin = "Jenis kelamin harus dipilih"
  }

  // Berat Lahir validation
  if (!formData.value.berat_lahir) {
    errors.value.berat_lahir = "Berat lahir harus diisi"
  } else {
    const berat = parseInt(formData.value.berat_lahir)
    if (isNaN(berat) || berat < 500 || berat > 6000) {
      errors.value.berat_lahir = "Berat lahir harus antara 500-6000 gram"
    }
  }

  // Tinggi Lahir validation
  if (!formData.value.tinggi_lahir) {
    errors.value.tinggi_lahir = "Tinggi lahir harus diisi"
  } else {
    const tinggi = parseInt(formData.value.tinggi_lahir)
    if (isNaN(tinggi) || tinggi < 25 || tinggi > 65) {
      errors.value.tinggi_lahir = "Tinggi lahir harus antara 25-65 cm"
    }
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    id_keluarga: "",
    nama: "",
    tanggal_lahir: "",
    jenis_kelamin: "L",
    berat_lahir: "",
    tinggi_lahir: "",
  }
  errors.value = {}
  selectedDate.value = undefined
}

// Load form data for edit mode
const loadFormData = (balita: Balita) => {
  formData.value = {
    id: balita.id,
    id_keluarga: balita.id_keluarga,
    nama: balita.nama,
    tanggal_lahir: balita.tanggal_lahir,
    jenis_kelamin: balita.jenis_kelamin,
    berat_lahir: balita.berat_lahir,
    tinggi_lahir: balita.tinggi_lahir,
    nomor_kk: balita.nomor_kk,
    nama_ayah: balita.nama_ayah,
    nama_ibu: balita.nama_ibu,
    umur: balita.umur,
    kelurahan: balita.kelurahan,
    kecamatan: balita.kecamatan,
    created_date: balita.created_date,
    updated_date: balita.updated_date,
  }
  errors.value = {}

  // Initialize date for calendar
  if (balita.tanggal_lahir) {
    const jsDate = new Date(balita.tanggal_lahir)
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

  // Auto-populate data keluarga berdasarkan id_keluarga dari API data
  const selectedKeluarga = getSelectedKeluarga(formData.value.id_keluarga!)
  if (selectedKeluarga) {
    formData.value.nomor_kk = selectedKeluarga.nomor_kk
    formData.value.nama_ayah = selectedKeluarga.nama_ayah
    formData.value.nama_ibu = selectedKeluarga.nama_ibu
    formData.value.kelurahan = selectedKeluarga.kelurahan
    formData.value.kecamatan = selectedKeluarga.kecamatan
  }

  // Calculate umur
  formData.value.umur = currentAge.value

  emit("save", formData.value as Balita)
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
    // Convert DateValue to string format for backend
    const jsDate = new Date(date.year, date.month - 1, date.day)
    formData.value.tanggal_lahir = format(jsDate, "yyyy-MM-dd")
    isCalendarOpen.value = false

    // Clear tanggal_lahir error if date is valid
    const birthDate = jsDate
    const today = new Date()
    const fiveYearsAgo = new Date()
    fiveYearsAgo.setFullYear(today.getFullYear() - 5)

    if (birthDate <= today && birthDate >= fiveYearsAgo) {
      delete errors.value.tanggal_lahir
    }
  }
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  async (newVal) => {
    console.log("👀 Dialog visibility changed:", newVal)

    if (newVal) {
      // Load keluarga options first
      await fetchKeluargaOptions()

      if (props.balita && props.mode === "edit") {
        loadFormData(props.balita)
      } else {
        resetForm()
      }
    }
  }
)

// Watch for balita prop changes
watch(
  () => props.balita,
  (newBalita) => {
    if (newBalita && props.mode === "edit" && props.show) {
      loadFormData(newBalita)
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.nama,
  (newVal) => {
    if (newVal && newVal.length >= 2 && /^[a-zA-Z\s]+$/.test(newVal)) {
      delete errors.value.nama
    }
  }
)

watch(
  () => formData.value.berat_lahir,
  (newVal) => {
    if (newVal) {
      const berat = parseInt(newVal)
      if (!isNaN(berat) && berat >= 500 && berat <= 6000) {
        delete errors.value.berat_lahir
      }
    }
  }
)

watch(
  () => formData.value.tinggi_lahir,
  (newVal) => {
    if (newVal) {
      const tinggi = parseInt(newVal)
      if (!isNaN(tinggi) && tinggi >= 25 && tinggi <= 65) {
        delete errors.value.tinggi_lahir
      }
    }
  }
)

// Load keluarga options on component mount
onMounted(async () => {
  await fetchKeluargaOptions()
})
</script>

<template>
  <Dialog
    :open="show"
    @update:open="handleClose">
    <DialogContent class="w-[95vw] max-w-2xl max-h-[95vh] overflow-hidden p-0">
      <div class="flex flex-col h-full max-h-[95vh]">
        <!-- Header -->
        <div class="p-4 sm:p-6 border-b bg-background">
          <DialogHeader class="space-y-2">
            <DialogTitle class="text-lg sm:text-xl flex items-center gap-2">
              <Baby class="h-5 w-5" />
              {{ mode === "create" ? "Tambah Balita Baru" : "Edit Data Balita" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data balita baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi balita. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Pilih Keluarga -->
            <Card :class="errors.id_keluarga ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Users class="h-4 w-4" />
                  Pilih Keluarga *
                  <Loader2
                    v-if="keluargaLoading"
                    class="h-3 w-3 animate-spin" />
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select
                  v-model="formData.id_keluarga"
                  :disabled="keluargaLoading">
                  <SelectTrigger class="w-full" :class="errors.id_keluarga ? 'border-red-500' : ''">
                    <SelectValue>
                      {{
                        formData.id_keluarga
                          ? getSelectedKeluarga(formData.id_keluarga)?.nomor_kk
                          : "Pilih Keluarga"
                      }}
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="keluarga in keluargaOptions"
                      :key="keluarga.id"
                      :value="keluarga.id">
                      {{ getKeluargaLabel(keluarga) }}
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_keluarga"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_keluarga }}
                </p>

                <!-- Show selected keluarga info -->
                <div
                  v-if="formData.id_keluarga && !keluargaLoading"
                  class="mt-2 text-xs text-muted-foreground bg-muted/50 p-2 rounded">
                  <strong>Dipilih:</strong>
                  {{ getSelectedKeluarga(formData.id_keluarga)?.nama_ayah }} &
                  {{ getSelectedKeluarga(formData.id_keluarga)?.nama_ibu }} -
                  {{ getSelectedKeluarga(formData.id_keluarga)?.kelurahan }},
                  {{ getSelectedKeluarga(formData.id_keluarga)?.kecamatan }}
                </div>
              </CardContent>
            </Card>

            <!-- Data Balita -->
            <div class="space-y-4">
              <!-- Nama Balita -->
              <div class="space-y-2">
                <Label
                  for="nama"
                  class="text-sm font-medium flex items-center gap-2">
                  Nama Balita *
                  <span class="text-xs text-muted-foreground">
                    (2-50 karakter, huruf dan spasi)
                  </span>
                </Label>
                <Input
                  id="nama"
                  v-model="formData.nama"
                  @input="(e: Event) => {
                    const target = e.target as HTMLInputElement
                    formData.nama = restrictToLetters(target.value).substring(0, 50)
                  }"
                  placeholder="Contoh: Andi Pratama"
                  maxlength="50"
                  :class="errors.nama ? 'border-red-500' : ''" />
                <p
                  v-if="errors.nama"
                  class="text-sm text-red-500">
                  {{ errors.nama }}
                </p>
              </div>

              <!-- Tanggal Lahir & Jenis Kelamin -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label
                    for="tanggal_lahir"
                    class="text-sm font-medium flex items-center gap-2">
                    Tanggal Lahir *
                    <span class="text-xs text-muted-foreground">(balita < 5 tahun)</span>
                  </Label>
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
                          errors.tanggal_lahir && 'border-red-500',
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
                            : "Pilih tanggal lahir"
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
                    v-if="errors.tanggal_lahir"
                    class="text-sm text-red-500">
                    {{ errors.tanggal_lahir }}
                  </p>
                  <p
                    v-if="currentAge"
                    class="text-sm text-blue-600 font-medium">
                    Umur: {{ currentAge }}
                  </p>
                </div>

                <div class="space-y-2">
                  <Label
                    for="jenis_kelamin"
                    class="text-sm font-medium">
                    Jenis Kelamin *
                  </Label>
                  <Select v-model="formData.jenis_kelamin">
                    <SelectTrigger :class="errors.jenis_kelamin ? 'border-red-500' : ''">
                      <SelectValue placeholder="Pilih jenis kelamin" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="L">Laki-laki</SelectItem>
                      <SelectItem value="P">Perempuan</SelectItem>
                    </SelectContent>
                  </Select>
                  <p
                    v-if="errors.jenis_kelamin"
                    class="text-sm text-red-500">
                    {{ errors.jenis_kelamin }}
                  </p>
                </div>
              </div>

              <!-- Berat & Tinggi Lahir -->
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div class="space-y-2">
                  <Label
                    for="berat_lahir"
                    class="text-sm font-medium flex items-center gap-2">
                    Berat Lahir *
                    <span class="text-xs text-muted-foreground">(500-6000 gram)</span>
                  </Label>
                  <div class="relative">
                    <Input
                      id="berat_lahir"
                      v-model="formData.berat_lahir"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.berat_lahir = restrictToNumbers(target.value).substring(0, 4)
                      }"
                      placeholder="3500"
                      maxlength="4"
                      inputmode="numeric"
                      :class="errors.berat_lahir ? 'border-red-500' : ''"
                      class="pr-12" />
                    <span
                      class="absolute right-3 top-1/2 transform -translate-y-1/2 text-sm text-muted-foreground">
                      gram
                    </span>
                  </div>
                  <p
                    v-if="errors.berat_lahir"
                    class="text-sm text-red-500">
                    {{ errors.berat_lahir }}
                  </p>
                </div>

                <div class="space-y-2">
                  <Label
                    for="tinggi_lahir"
                    class="text-sm font-medium flex items-center gap-2">
                    Tinggi Lahir *
                    <span class="text-xs text-muted-foreground">(25-65 cm)</span>
                  </Label>
                  <div class="relative">
                    <Input
                      id="tinggi_lahir"
                      v-model="formData.tinggi_lahir"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.tinggi_lahir = restrictToNumbers(target.value).substring(0, 2)
                      }"
                      placeholder="50"
                      maxlength="2"
                      inputmode="numeric"
                      :class="errors.tinggi_lahir ? 'border-red-500' : ''"
                      class="pr-8" />
                    <span
                      class="absolute right-3 top-1/2 transform -translate-y-1/2 text-sm text-muted-foreground">
                      cm
                    </span>
                  </div>
                  <p
                    v-if="errors.tinggi_lahir"
                    class="text-sm text-red-500">
                    {{ errors.tinggi_lahir }}
                  </p>
                </div>
              </div>
            </div>

            <!-- Info Box -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
              <p class="text-sm text-blue-800 flex items-start gap-2">
                <Baby class="h-4 w-4 mt-0.5 flex-shrink-0" />
                <span>
                  <strong>Kriteria Balita:</strong><br />
                  • Anak berusia di bawah 5 tahun (balita)<br />
                  • Berat lahir normal: 2500-4000 gram<br />
                  • Tinggi lahir normal: 45-55 cm<br />
                  • Data ini akan digunakan untuk pemantauan pertumbuhan
                </span>
              </p>
            </div>

            <!-- Development Info -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Mode: {{ mode }}<br />
              Has Data: {{ !!props.balita }}<br />
              Form ID: {{ formData.id || "New" }}<br />
              Umur: {{ currentAge }}<br />
              Keluarga Options: {{ keluargaOptions.length }} loaded<br />
              Loading: {{ keluargaLoading }}
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
              :disabled="loading || keluargaLoading"
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
