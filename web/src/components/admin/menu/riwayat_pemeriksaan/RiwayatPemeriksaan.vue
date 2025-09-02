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
import { Badge } from "@/components/ui/badge"
import {
  CalendarIcon,
  Save,
  X,
  Stethoscope,
  Baby,
  Activity,
  FileText,
  Loader2,
} from "lucide-vue-next"
import { format } from "date-fns"
import { id } from "date-fns/locale"
import { CalendarDate, type DateValue } from "@internationalized/date"
import { authUtils } from "@/lib/utils"
import type { RiwayatPemeriksaan } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  riwayat: RiwayatPemeriksaan | null
  loading?: boolean
}

interface Emits {
  (e: "close"): void
  (e: "save", riwayat: RiwayatPemeriksaan): void
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

interface IntervensiOption {
  id: string
  id_balita: string
  nama_balita: string
  jenis: string
  tanggal: string
  deskripsi: string
}

interface LaporanOption {
  id: string
  id_balita: string
  nama_balita: string
  status_laporan: string
  tanggal_laporan: string
  jenis_laporan: string
}

interface GetAllBalitaResponse {
  data: BalitaOption[]
  total: number
}

interface GetAllIntervensiResponse {
  data: IntervensiOption[]
  total: number
}

interface GetAllLaporanResponse {
  data: LaporanOption[]
  total: number
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data sesuai dengan backend API
const formData = ref({
  // Fields untuk create/edit
  id_balita: "",
  id_intervensi: "",
  id_laporan_masyarakat: "", // Wajib sesuai backend
  tanggal: "",
  berat_badan: "",
  tinggi_badan: "",
  status_gizi: "",
  keterangan: "",

  // Fields tambahan untuk edit mode
  id: "",

  // Fields untuk response (tidak dikirim ke backend)
  nama_balita: "",
  umur_balita: "",
  jenis_kelamin: "",
  nama_ayah: "",
  nama_ibu: "",
  nomor_kk: "",
  jenis_intervensi: "",
  tanggal_intervensi: "",
  status_laporan: "",
  tanggal_laporan: "",
  jenis_laporan: "",
  kelurahan: "",
  kecamatan: "",
  created_date: "",
  updated_date: "",
  created_by: "",
  updated_by: "",
})

// Date picker state
const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

// Validation errors
const errors = ref<Record<string, string>>({})

// Master data from API
const balitaOptions = ref<BalitaOption[]>([])
const intervensiOptions = ref<IntervensiOption[]>([])
const laporanOptions = ref<LaporanOption[]>([])

// Loading states
const balitaLoading = ref(false)
const intervensiLoading = ref(false)
const laporanLoading = ref(false)

// Static status gizi options
const statusGiziOptions = [
  {
    value: "normal",
    label: "Normal",
    description: "Status gizi dalam batas normal",
    color: "text-green-600",
  },
  {
    value: "stunting",
    label: "Stunting",
    description: "Kondisi gagal tumbuh pada balita",
    color: "text-orange-600",
  },
  {
    value: "gizi buruk",
    label: "Gizi Buruk",
    description: "Kondisi kekurangan gizi yang parah",
    color: "text-red-600",
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

// Fetch master data
const fetchBalitaOptions = async () => {
  balitaLoading.value = true
  try {
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
    toast.error("Gagal memuat data balita")
  } finally {
    balitaLoading.value = false
  }
}

const fetchIntervensiOptions = async () => {
  intervensiLoading.value = true
  try {
    const response = await apiRequest<GetAllIntervensiResponse>("/admin/intervensi/get")
    intervensiOptions.value = response.data.data.map((i) => ({
      id: i.id,
      id_balita: i.id_balita,
      nama_balita: i.nama_balita,
      jenis: i.jenis,
      tanggal: i.tanggal,
      deskripsi: i.deskripsi,
    }))
    console.log("Intervensi options loaded:", intervensiOptions.value.length, "items")
  } catch (error) {
    console.error("Error fetching intervensi options:", error)
    toast.error("Gagal memuat data intervensi")
  } finally {
    intervensiLoading.value = false
  }
}

const fetchLaporanOptions = async () => {
  laporanLoading.value = true
  try {
    const response = await apiRequest<GetAllLaporanResponse>("/admin/laporan-masyarakat/get")
    laporanOptions.value = response.data.data.map((l) => ({
      id: l.id,
      id_balita: l.id_balita,
      nama_balita: l.nama_balita,
      status_laporan: l.status_laporan,
      tanggal_laporan: l.tanggal_laporan,
      jenis_laporan: l.jenis_laporan,
    }))
    console.log("Laporan options loaded:", laporanOptions.value.length, "items")
  } catch (error) {
    console.error("Error fetching laporan options:", error)
    toast.error("Gagal memuat data laporan masyarakat")
  } finally {
    laporanLoading.value = false
  }
}

// Get selected info
const selectedBalitaInfo = computed(() => {
  if (!formData.value.id_balita) return null
  return balitaOptions.value.find((b) => b.id === formData.value.id_balita) || null
})

const selectedIntervensiInfo = computed(() => {
  if (!formData.value.id_intervensi) return null
  return intervensiOptions.value.find((i) => i.id === formData.value.id_intervensi) || null
})

const selectedLaporanInfo = computed(() => {
  if (!formData.value.id_laporan_masyarakat) return null
  return laporanOptions.value.find((l) => l.id === formData.value.id_laporan_masyarakat) || null
})

const selectedStatusGiziInfo = computed(() => {
  if (!formData.value.status_gizi) return null
  return statusGiziOptions.find((s) => s.value === formData.value.status_gizi) || null
})

// Filtered options based on selected balita
const filteredIntervensiOptions = computed(() => {
  if (!formData.value.id_balita) return []
  return intervensiOptions.value.filter((i) => i.id_balita === formData.value.id_balita)
})

const filteredLaporanOptions = computed(() => {
  if (!formData.value.id_balita) return []
  return laporanOptions.value.filter((l) => l.id_balita === formData.value.id_balita)
})

// Form validation dengan pengecekan lebih ketat
const validateForm = (): boolean => {
  errors.value = {}

  console.log("Validating form data:", formData.value) // Debug log

  // Balita validation
  if (!formData.value.id_balita || formData.value.id_balita.trim() === "") {
    errors.value.id_balita = "Balita harus dipilih"
  }

  // Intervensi validation
  if (!formData.value.id_intervensi || formData.value.id_intervensi.trim() === "") {
    errors.value.id_intervensi = "Intervensi harus dipilih"
  }

  // Laporan validation (WAJIB)
  if (!formData.value.id_laporan_masyarakat || formData.value.id_laporan_masyarakat.trim() === "") {
    errors.value.id_laporan_masyarakat = "Laporan masyarakat harus dipilih"
  }

  // Tanggal validation
  if (!formData.value.tanggal || formData.value.tanggal.trim() === "") {
    errors.value.tanggal = "Tanggal pemeriksaan harus diisi"
  } else {
    // Cek format tanggal YYYY-MM-DD
    const dateRegex = /^\d{4}-\d{2}-\d{2}$/
    if (!dateRegex.test(formData.value.tanggal)) {
      errors.value.tanggal = "Format tanggal harus YYYY-MM-DD"
    } else {
      const pemeriksaanDate = new Date(formData.value.tanggal)
      const today = new Date()
      today.setHours(23, 59, 59, 999) // Set ke akhir hari

      if (pemeriksaanDate > today) {
        errors.value.tanggal = "Tanggal pemeriksaan tidak boleh di masa depan"
      }
    }
  }

  // Berat badan validation
  if (!formData.value.berat_badan || formData.value.berat_badan.toString().trim() === "") {
    errors.value.berat_badan = "Berat badan harus diisi"
  } else {
    const berat = parseFloat(formData.value.berat_badan.toString())
    if (isNaN(berat) || berat < 1 || berat > 50) {
      errors.value.berat_badan = "Berat badan harus antara 1-50 kg"
    }
  }

  // Tinggi badan validation
  if (!formData.value.tinggi_badan || formData.value.tinggi_badan.toString().trim() === "") {
    errors.value.tinggi_badan = "Tinggi badan harus diisi"
  } else {
    const tinggi = parseFloat(formData.value.tinggi_badan.toString())
    if (isNaN(tinggi) || tinggi < 30 || tinggi > 150) {
      errors.value.tinggi_badan = "Tinggi badan harus antara 30-150 cm"
    }
  }

  // Status gizi validation
  if (!formData.value.status_gizi || formData.value.status_gizi.trim() === "") {
    errors.value.status_gizi = "Status gizi harus dipilih"
  } else {
    const allowedStatuses = ["normal", "stunting", "gizi buruk"]
    if (!allowedStatuses.includes(formData.value.status_gizi)) {
      errors.value.status_gizi = "Status gizi tidak valid"
    }
  }

  // Keterangan validation
  if (!formData.value.keterangan || formData.value.keterangan.trim().length < 5) {
    errors.value.keterangan = "Keterangan minimal 5 karakter"
  } else if (formData.value.keterangan.length > 500) {
    errors.value.keterangan = "Keterangan maksimal 500 karakter"
  }

  console.log("Validation errors:", errors.value) // Debug log
  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    id_balita: "",
    id_intervensi: "",
    id_laporan_masyarakat: "",
    tanggal: "",
    berat_badan: "",
    tinggi_badan: "",
    status_gizi: "",
    keterangan: "",
    id: "",
    nama_balita: "",
    umur_balita: "",
    jenis_kelamin: "",
    nama_ayah: "",
    nama_ibu: "",
    nomor_kk: "",
    jenis_intervensi: "",
    tanggal_intervensi: "",
    status_laporan: "",
    tanggal_laporan: "",
    jenis_laporan: "",
    kelurahan: "",
    kecamatan: "",
    created_date: "",
    updated_date: "",
    created_by: "",
    updated_by: "",
  }
  errors.value = {}
  selectedDate.value = undefined
}

// Load form data for edit mode
const loadFormData = (riwayat: RiwayatPemeriksaan) => {
  formData.value = {
    id: riwayat.id,
    id_balita: riwayat.id_balita,
    id_intervensi: riwayat.id_intervensi,
    id_laporan_masyarakat: riwayat.id_laporan_masyarakat,
    tanggal: riwayat.tanggal,
    berat_badan: riwayat.berat_badan,
    tinggi_badan: riwayat.tinggi_badan,
    status_gizi: riwayat.status_gizi,
    keterangan: riwayat.keterangan,

    // Extended fields untuk display
    nama_balita: riwayat.nama_balita || "",
    umur_balita: riwayat.umur_balita || "",
    jenis_kelamin: riwayat.jenis_kelamin || "",
    nama_ayah: riwayat.nama_ayah || "",
    nama_ibu: riwayat.nama_ibu || "",
    nomor_kk: riwayat.nomor_kk || "",
    jenis_intervensi: riwayat.jenis_intervensi || "",
    tanggal_intervensi: riwayat.tanggal_intervensi || "",
    status_laporan: riwayat.status_laporan || "",
    tanggal_laporan: riwayat.tanggal_laporan || "",
    jenis_laporan: riwayat.jenis_laporan || "",
    kelurahan: riwayat.kelurahan || "",
    kecamatan: riwayat.kecamatan || "",
    created_date: riwayat.created_date || "",
    updated_date: riwayat.updated_date || "",
    created_by: riwayat.created_by || "",
    updated_by: riwayat.updated_by || "",
  }
  errors.value = {}

  // Initialize date for calendar
  if (riwayat.tanggal) {
    const jsDate = new Date(riwayat.tanggal)
    selectedDate.value = new CalendarDate(
      jsDate.getFullYear(),
      jsDate.getMonth() + 1,
      jsDate.getDate()
    )
  }
}

// Handle date selection
const handleDateSelect = (date: DateValue | undefined) => {
  if (date) {
    selectedDate.value = date
    const jsDate = new Date(date.year, date.month - 1, date.day)
    formData.value.tanggal = format(jsDate, "yyyy-MM-dd")
    isCalendarOpen.value = false

    // Clear tanggal error if date is valid
    const pemeriksaanDate = jsDate
    const today = new Date()
    today.setHours(23, 59, 59, 999)

    if (pemeriksaanDate <= today) {
      delete errors.value.tanggal
    }
  }
}

// Handle save dengan payload yang tepat
const handleSave = () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi")
    return
  }

  // Prepare payload sesuai backend API
  const payload = {
    ...(props.mode === "edit" && { id: formData.value.id }),
    id_balita: formData.value.id_balita,
    id_intervensi: formData.value.id_intervensi,
    id_laporan_masyarakat: formData.value.id_laporan_masyarakat,
    tanggal: formData.value.tanggal,
    berat_badan: formData.value.berat_badan.toString(),
    tinggi_badan: formData.value.tinggi_badan.toString(),
    status_gizi: formData.value.status_gizi,
    keterangan: formData.value.keterangan.trim(),

    // Extended fields untuk komponen parent
    nama_balita: selectedBalitaInfo.value?.nama || "",
    umur_balita: selectedBalitaInfo.value?.umur || "",
    jenis_kelamin: selectedBalitaInfo.value?.jenis_kelamin || "",
    nama_ayah: selectedBalitaInfo.value?.nama_ayah || "",
    nama_ibu: selectedBalitaInfo.value?.nama_ibu || "",
    nomor_kk: selectedBalitaInfo.value?.nomor_kk || "",
    jenis_intervensi: selectedIntervensiInfo.value?.jenis || "",
    tanggal_intervensi: selectedIntervensiInfo.value?.tanggal || "",
    status_laporan: selectedLaporanInfo.value?.status_laporan || "",
    tanggal_laporan: selectedLaporanInfo.value?.tanggal_laporan || "",
    jenis_laporan: selectedLaporanInfo.value?.jenis_laporan || "",
    kelurahan: selectedBalitaInfo.value?.kelurahan || "",
    kecamatan: selectedBalitaInfo.value?.kecamatan || "",
  }

  console.log("Sending payload:", payload) // Debug log

  emit("save", payload as RiwayatPemeriksaan)
  resetForm()
}

// Handle close
const handleClose = () => {
  resetForm()
  emit("close")
}

// Watch for dialog visibility
watch(
  () => props.show,
  async (newVal) => {
    if (newVal) {
      // Load all master data
      await Promise.all([fetchBalitaOptions(), fetchIntervensiOptions(), fetchLaporanOptions()])

      if (props.riwayat && props.mode === "edit") {
        loadFormData(props.riwayat)
      } else {
        resetForm()
        // Set default date to today for create mode
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

// Watch for riwayat prop changes
watch(
  () => props.riwayat,
  (newRiwayat) => {
    if (newRiwayat && props.mode === "edit" && props.show) {
      loadFormData(newRiwayat)
    }
  },
  { deep: true }
)

// Clear dependent selections when balita changes
watch(
  () => formData.value.id_balita,
  () => {
    formData.value.id_intervensi = ""
    formData.value.id_laporan_masyarakat = ""
  }
)

// Real-time validation
watch(
  () => formData.value.keterangan,
  (newVal) => {
    if (newVal && newVal.trim().length >= 5 && newVal.length <= 500) {
      delete errors.value.keterangan
    }
  }
)

// Load master data on component mount
onMounted(async () => {
  await Promise.all([fetchBalitaOptions(), fetchIntervensiOptions(), fetchLaporanOptions()])
})
</script>

<template>
  <!-- Template tetap sama seperti sebelumnya -->
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
              {{ mode === "create" ? "Tambah Riwayat Pemeriksaan" : "Edit Riwayat Pemeriksaan" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form untuk menambah riwayat pemeriksaan baru."
                  : "Perbarui informasi riwayat pemeriksaan."
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
                        {{ balitaLoading ? "Memuat balita..." : "Pilih balita" }}
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
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_balita"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_balita }}
                </p>
              </CardContent>
            </Card>

            <!-- Pilih Intervensi -->
            <Card :class="errors.id_intervensi ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Activity class="h-4 w-4" />
                  Pilih Intervensi *
                  <Loader2
                    v-if="intervensiLoading"
                    class="h-3 w-3 animate-spin" />
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select
                  v-model="formData.id_intervensi"
                  :disabled="intervensiLoading || !formData.id_balita">
                  <SelectTrigger
                    class="w-full"
                    :class="errors.id_intervensi ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedIntervensiInfo">
                        {{ selectedIntervensiInfo.jenis }} - {{ selectedIntervensiInfo.tanggal }}
                      </template>
                      <template v-else>
                        {{
                          !formData.id_balita
                            ? "Pilih balita terlebih dahulu"
                            : intervensiLoading
                            ? "Memuat intervensi..."
                            : "Pilih intervensi"
                        }}
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="intervensi in filteredIntervensiOptions"
                      :key="intervensi.id"
                      :value="intervensi.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <span class="font-medium">{{ intervensi.jenis }}</span>
                          <span class="text-xs text-gray-600">{{ intervensi.tanggal }}</span>
                        </div>
                        <div class="text-xs text-muted-foreground">
                          {{ intervensi.deskripsi?.substring(0, 50)
                          }}{{ intervensi.deskripsi?.length > 50 ? "..." : "" }}
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_intervensi"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_intervensi }}
                </p>
              </CardContent>
            </Card>

            <!-- Pilih Laporan Masyarakat -->
            <Card :class="errors.id_laporan_masyarakat ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <FileText class="h-4 w-4" />
                  Pilih Laporan Masyarakat *
                  <Loader2
                    v-if="laporanLoading"
                    class="h-3 w-3 animate-spin" />
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select
                  v-model="formData.id_laporan_masyarakat"
                  :disabled="laporanLoading || !formData.id_balita">
                  <SelectTrigger
                    class="w-full"
                    :class="errors.id_laporan_masyarakat ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedLaporanInfo">
                        {{ selectedLaporanInfo.status_laporan }} -
                        {{ selectedLaporanInfo.tanggal_laporan }}
                      </template>
                      <template v-else>
                        {{
                          !formData.id_balita
                            ? "Pilih balita terlebih dahulu"
                            : laporanLoading
                            ? "Memuat laporan..."
                            : "Pilih laporan masyarakat"
                        }}
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="laporan in filteredLaporanOptions"
                      :key="laporan.id"
                      :value="laporan.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <span class="font-medium">{{ laporan.status_laporan }}</span>
                          <Badge
                            variant="outline"
                            class="text-xs"
                            >{{ laporan.jenis_laporan }}</Badge
                          >
                        </div>
                        <div class="text-xs text-muted-foreground">
                          📅 {{ laporan.tanggal_laporan }}
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.id_laporan_masyarakat"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.id_laporan_masyarakat }}
                </p>
              </CardContent>
            </Card>

            <!-- Tanggal & Pengukuran -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <!-- Tanggal Pemeriksaan -->
              <Card :class="errors.tanggal ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="flex items-center gap-2 text-base">
                    <CalendarIcon class="h-4 w-4" />
                    Tanggal *
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
                            : "Pilih tanggal"
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
                </CardContent>
              </Card>

              <!-- Berat Badan -->
              <Card :class="errors.berat_badan ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="text-base">Berat Badan *</CardTitle>
                </CardHeader>
                <CardContent>
                  <div class="relative">
                    <Input
                      v-model="formData.berat_badan"
                      type="number"
                      step="0.1"
                      min="1"
                      max="50"
                      placeholder="10.5"
                      :class="errors.berat_badan ? 'border-red-500' : ''" />
                    <span
                      class="absolute right-3 top-1/2 transform -translate-y-1/2 text-sm text-gray-500"
                      >kg</span
                    >
                  </div>
                  <p
                    v-if="errors.berat_badan"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.berat_badan }}
                  </p>
                </CardContent>
              </Card>

              <!-- Tinggi Badan -->
              <Card :class="errors.tinggi_badan ? 'border-red-500' : ''">
                <CardHeader class="pb-3">
                  <CardTitle class="text-base">Tinggi Badan *</CardTitle>
                </CardHeader>
                <CardContent>
                  <div class="relative">
                    <Input
                      v-model="formData.tinggi_badan"
                      type="number"
                      step="0.1"
                      min="30"
                      max="150"
                      placeholder="80.0"
                      :class="errors.tinggi_badan ? 'border-red-500' : ''" />
                    <span
                      class="absolute right-3 top-1/2 transform -translate-y-1/2 text-sm text-gray-500"
                      >cm</span
                    >
                  </div>
                  <p
                    v-if="errors.tinggi_badan"
                    class="text-sm text-red-500 mt-2">
                    {{ errors.tinggi_badan }}
                  </p>
                </CardContent>
              </Card>
            </div>

            <!-- Status Gizi -->
            <Card :class="errors.status_gizi ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base"> Status Gizi * </CardTitle>
              </CardHeader>
              <CardContent>
                <Select v-model="formData.status_gizi">
                  <SelectTrigger
                    class="w-full"
                    :class="errors.status_gizi ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedStatusGiziInfo">
                        <div class="flex items-center gap-2">
                          <span :class="selectedStatusGiziInfo.color">{{
                            selectedStatusGiziInfo.label
                          }}</span>
                        </div>
                      </template>
                      <template v-else> Pilih status gizi </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="status in statusGiziOptions"
                      :key="status.value"
                      :value="status.value">
                      <div class="flex flex-col py-1">
                        <span
                          class="font-medium"
                          :class="status.color"
                          >{{ status.label }}</span
                        >
                        <span class="text-xs text-muted-foreground">{{ status.description }}</span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.status_gizi"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.status_gizi }}
                </p>
              </CardContent>
            </Card>

            <!-- Keterangan -->
            <Card :class="errors.keterangan ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base"> Keterangan * </CardTitle>
              </CardHeader>
              <CardContent>
                <Textarea
                  v-model="formData.keterangan"
                  placeholder="Catatan hasil pemeriksaan, rekomendasi, dan tindak lanjut..."
                  rows="3"
                  :class="errors.keterangan ? 'border-red-500' : ''" />
                <p
                  v-if="errors.keterangan"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.keterangan }}
                </p>
                <div class="flex justify-between items-center mt-2">
                  <p class="text-xs text-muted-foreground">Minimum 5 karakter</p>
                  <p class="text-xs text-muted-foreground">
                    {{ formData.keterangan?.length || 0 }}/500 karakter
                  </p>
                </div>
              </CardContent>
            </Card>

            <!-- Debug Info -->
            <Card class="border-yellow-200 bg-yellow-50">
              <CardHeader>
                <CardTitle class="text-sm text-yellow-800">🐛 Form Debug Info</CardTitle>
              </CardHeader>
              <CardContent class="text-xs text-yellow-700 space-y-1">
                <div><strong>Form Data:</strong></div>
                <div>• Balita ID: {{ formData.id_balita || "Empty" }}</div>
                <div>• Intervensi ID: {{ formData.id_intervensi || "Empty" }}</div>
                <div>• Laporan ID: {{ formData.id_laporan_masyarakat || "Empty" }}</div>
                <div>• Tanggal: {{ formData.tanggal || "Empty" }}</div>
                <div>• BB: {{ formData.berat_badan || "Empty" }}</div>
                <div>• TB: {{ formData.tinggi_badan || "Empty" }}</div>
                <div>• Status: {{ formData.status_gizi || "Empty" }}</div>
                <div>• Keterangan: {{ formData.keterangan?.length || 0 }} chars</div>
                <div class="pt-2"><strong>Validation:</strong></div>
                <div>• Errors: {{ Object.keys(errors).length }}</div>
                <div v-if="Object.keys(errors).length > 0">
                  • {{ Object.keys(errors).join(", ") }}
                </div>
                <div class="pt-2"><strong>Master Data:</strong></div>
                <div>• Balita: {{ balitaOptions.length }} loaded</div>
                <div>• Intervensi: {{ intervensiOptions.length }} loaded</div>
                <div>• Laporan: {{ laporanOptions.length }} loaded</div>
                <div>• Filtered Intervensi: {{ filteredIntervensiOptions.length }}</div>
                <div>• Filtered Laporan: {{ filteredLaporanOptions.length }}</div>
              </CardContent>
            </Card>
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
              :disabled="loading || balitaLoading || intervensiLoading || laporanLoading"
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
