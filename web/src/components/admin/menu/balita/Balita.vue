<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Baby, Users, Calendar, Loader2, RefreshCcw } from "lucide-vue-next"
import { authUtils } from "@/lib/utils"
import type { Balita } from "./columns"

// API Response Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface GetAllBalitaResponse {
  data: Balita[]
  total: number
}

interface InsertBalitaResponse {
  id: string
}

interface UpdateDeleteBalitaResponse {
  id: string
  message: string
}

// State management
const balitaData = ref<Balita[]>([])
const isLoading = ref(false)
const isDialogLoading = ref(false)

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedBalita = ref<Balita | null>(null)

// Statistics
const totalBalita = ref(0)
const totalLakiLaki = ref(0)
const totalPerempuan = ref(0)
const averageAge = ref(0)

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

// Get all balita from API
const fetchBalitaData = async () => {
  isLoading.value = true
  try {
    console.log("Fetching balita data from API...")

    const response = await apiRequest<GetAllBalitaResponse>("/admin/balita/get")

    balitaData.value = response.data.data
    totalBalita.value = response.data.total

    // Calculate statistics
    updateStatistics()

    console.log("Balita data fetched successfully:", {
      total: totalBalita.value,
      lakiLaki: totalLakiLaki.value,
      perempuan: totalPerempuan.value,
    })

    toast.success(`Data balita berhasil dimuat (${totalBalita.value} data)`)
  } catch (error) {
    console.error("Error fetching balita data:", error)
    toast.error("Gagal memuat data balita. Silakan coba lagi.")

    // Fallback to empty data
    balitaData.value = []
    totalBalita.value = 0
    updateStatistics()
  } finally {
    isLoading.value = false
  }
}

// Calculate statistics
const updateStatistics = () => {
  totalBalita.value = balitaData.value.length
  totalLakiLaki.value = balitaData.value.filter((b) => b.jenis_kelamin === "L").length
  totalPerempuan.value = balitaData.value.filter((b) => b.jenis_kelamin === "P").length

  // Calculate average age in months
  if (balitaData.value.length > 0) {
    const totalMonths = balitaData.value.reduce((sum, balita) => {
      const ageText = balita.umur || "0 bulan"
      const months = parseAgeToMonths(ageText)
      return sum + months
    }, 0)
    averageAge.value = Math.round(totalMonths / balitaData.value.length)
  } else {
    averageAge.value = 0
  }
}

// Helper function to parse age text to months
const parseAgeToMonths = (ageText: string): number => {
  const tahunMatch = ageText.match(/(\d+)\s*tahun/)
  const bulanMatch = ageText.match(/(\d+)\s*bulan/)

  const years = tahunMatch ? parseInt(tahunMatch[1]) : 0
  const months = bulanMatch ? parseInt(bulanMatch[1]) : 0

  return years * 12 + months
}

// Format average age for display
const formatAverageAge = (months: number): string => {
  if (months < 12) {
    return `${months} bulan`
  } else {
    const years = Math.floor(months / 12)
    const remainingMonths = months % 12
    if (remainingMonths === 0) {
      return `${years} tahun`
    }
    return `${years}.${Math.round((remainingMonths / 12) * 10)} tahun`
  }
}

// Create new balita
const createBalita = async (balitaPayload: Partial<Balita>) => {
  try {
    console.log("Creating new balita:", balitaPayload)

    const payload = {
      id_keluarga: balitaPayload.id_keluarga,
      nama: balitaPayload.nama,
      tanggal_lahir: balitaPayload.tanggal_lahir,
      jenis_kelamin: balitaPayload.jenis_kelamin,
      berat_lahir: balitaPayload.berat_lahir,
      tinggi_lahir: balitaPayload.tinggi_lahir,
    }

    const response = await apiRequest<InsertBalitaResponse>("/admin/balita/insert", {
      method: "POST",
      body: JSON.stringify(payload),
    })

    console.log("Balita created successfully:", response)
    toast.success("Data balita berhasil ditambahkan")

    // Refresh data
    await fetchBalitaData()

    return response.data
  } catch (error) {
    console.error("Error creating balita:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Data balita sudah ada. Periksa nama dan tanggal lahir.")
      } else if (error.message.includes("Keluarga not found")) {
        toast.error("Keluarga tidak ditemukan. Pastikan data keluarga masih aktif.")
      } else if (error.message.includes("under 5 years old")) {
        toast.error("Anak harus berusia di bawah 5 tahun (kriteria balita).")
      } else {
        toast.error(`Gagal menambah data balita: ${error.message}`)
      }
    } else {
      toast.error("Gagal menambah data balita. Silakan coba lagi.")
    }
    throw error
  }
}

// Update existing balita
const updateBalita = async (balitaPayload: Balita) => {
  try {
    console.log("Updating balita:", balitaPayload)

    const payload = {
      id: balitaPayload.id,
      id_keluarga: balitaPayload.id_keluarga,
      nama: balitaPayload.nama,
      tanggal_lahir: balitaPayload.tanggal_lahir,
      jenis_kelamin: balitaPayload.jenis_kelamin,
      berat_lahir: balitaPayload.berat_lahir,
      tinggi_lahir: balitaPayload.tinggi_lahir,
    }

    const response = await apiRequest<UpdateDeleteBalitaResponse>("/admin/balita/update", {
      method: "PUT",
      body: JSON.stringify(payload),
    })

    console.log("Balita updated successfully:", response)
    toast.success("Data balita berhasil diperbarui")

    // Refresh data
    await fetchBalitaData()

    return response.data
  } catch (error) {
    console.error("Error updating balita:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Data balita sudah ada. Periksa nama dan tanggal lahir.")
      } else if (error.message.includes("Keluarga not found")) {
        toast.error("Keluarga tidak ditemukan. Pastikan data keluarga masih aktif.")
      } else if (error.message.includes("under 5 years old")) {
        toast.error("Anak harus berusia di bawah 5 tahun (kriteria balita).")
      } else {
        toast.error(`Gagal memperbarui data balita: ${error.message}`)
        showDialog.value = false
      }
    } else {
      toast.error("Gagal memperbarui data balita. Silakan coba lagi.")
    }
    throw error
  }
}

// Delete balita (soft delete)
const deleteBalita = async (balitaId: string) => {
  try {
    console.log("Deleting balita:", balitaId)

    const response = await apiRequest<UpdateDeleteBalitaResponse>("/admin/balita/delete", {
      method: "DELETE",
      body: JSON.stringify({ id: balitaId }),
    })

    console.log("Balita deleted successfully:", response)
    toast.success("Data balita berhasil dihapus")

    // Refresh data
    await fetchBalitaData()

    return response.data
  } catch (error) {
    console.error("Error deleting balita:", error)
    if (error instanceof Error) {
      if (error.message.includes("active laporan")) {
        toast.error("Tidak dapat menghapus balita yang masih memiliki laporan aktif.")
      } else if (error.message.includes("active riwayat pemeriksaan")) {
        toast.error("Tidak dapat menghapus balita yang masih memiliki riwayat pemeriksaan aktif.")
      } else {
        toast.error(`Gagal menghapus data balita: ${error.message}`)
      }
    } else {
      toast.error("Gagal menghapus data balita. Silakan coba lagi.")
    }
    throw error
  }
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedBalita.value = null
  showDialog.value = true
}

const handleEdit = (balita: Balita) => {
  dialogMode.value = "edit"
  selectedBalita.value = { ...balita }
  showDialog.value = true
}

const handleDelete = async (balita: Balita) => {
  // Show confirmation dialog
  const isConfirmed = confirm(
    `Apakah Anda yakin ingin menghapus data balita ${balita.nama}?\n\n` +
      `Detail:\n` +
      `• Nama: ${balita.nama}\n` +
      `• Jenis Kelamin: ${balita.jenis_kelamin === "L" ? "Laki-laki" : "Perempuan"}\n` +
      `• Umur: ${balita.umur}\n` +
      `• Orang Tua: ${balita.nama_ayah} & ${balita.nama_ibu}\n\n` +
      `Data akan dihapus secara permanen dari sistem.`
  )

  if (!isConfirmed) return

  isLoading.value = true
  try {
    await deleteBalita(balita.id)
  } catch (error) {
    // Error handling already done in deleteBalita function
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (balita: Balita | Partial<Balita>) => {
  isDialogLoading.value = true

  try {
    if (dialogMode.value === "create") {
      await createBalita(balita)
    } else {
      await updateBalita(balita as Balita)
    }

    showDialog.value = false
  } catch (error) {
    // Error handling already done in create/update functions
    // Keep dialog open so user can fix issues
  } finally {
    isDialogLoading.value = false
  }
}

// Custom event handlers (for DataTable component events)
const handleCustomEvents = async (event: Event) => {
  const customEvent = event as CustomEvent

  if (event.type === "edit-balita") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-balita") {
    await handleDelete(customEvent.detail)
  }
}

// Lifecycle hooks
onMounted(async () => {
  // Add event listeners for custom events from DataTable
  document.addEventListener("edit-balita", handleCustomEvents)
  document.addEventListener("delete-balita", handleCustomEvents)

  // Fetch initial data
  await fetchBalitaData()
})

onUnmounted(() => {
  document.removeEventListener("edit-balita", handleCustomEvents)
  document.removeEventListener("delete-balita", handleCustomEvents)
})

// Refresh data function (can be called manually)
const refreshData = async () => {
  await fetchBalitaData()
}

// Error retry function
const retryFetchData = async () => {
  toast.info("Mencoba mengambil data ulang...")
  await fetchBalitaData()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Baby class="h-8 w-8 text-blue-600" />
          Data Balita
        </h1>
        <p class="text-gray-600">Kelola data balita (anak usia di bawah 5 tahun) di Kota Cirebon</p>
      </div>
      <div class="flex gap-2">
        <!-- Refresh Button -->
        <Button
          @click="refreshData"
          :disabled="isLoading"
          variant="outline"
          size="sm">
          <Loader2
            v-if="isLoading"
            class="h-4 w-4 mr-2 animate-spin" />
          <RefreshCcw
            v-else
            class="h-4 w-4 mr-2" />
          Refresh
        </Button>

        <!-- Add Button -->
        <Button
          @click="handleCreate"
          :disabled="isLoading"
          class="gap-2">
          <Plus class="h-4 w-4" />
          Tambah Balita
        </Button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Balita</CardTitle>
          <Baby class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : totalBalita }}
          </div>
          <p class="text-xs text-muted-foreground">Balita terdaftar</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Laki-laki</CardTitle>
          <Users class="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600">
            {{ isLoading ? "..." : totalLakiLaki }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ totalBalita > 0 ? Math.round((totalLakiLaki / totalBalita) * 100) : 0 }}% dari total
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Perempuan</CardTitle>
          <Users class="h-4 w-4 text-pink-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-pink-600">
            {{ isLoading ? "..." : totalPerempuan }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ totalBalita > 0 ? Math.round((totalPerempuan / totalBalita) * 100) : 0 }}% dari total
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Usia Rata-rata</CardTitle>
          <Calendar class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : formatAverageAge(averageAge) }}
          </div>
          <p class="text-xs text-muted-foreground">Rata-rata balita</p>
        </CardContent>
      </Card>
    </div>

    <!-- Loading State -->
    <div
      v-if="isLoading && balitaData.length === 0"
      class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">Memuat Data Balita</h3>
        <p class="text-gray-600">Sedang mengambil data dari server...</p>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="!isLoading && balitaData.length === 0"
      class="text-center py-12">
      <div class="mx-auto mb-4 w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
        <Baby class="h-8 w-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Tidak Ada Data Balita</h3>
      <p class="text-gray-600 mb-4">Belum ada data balita yang terdaftar dalam sistem.</p>
      <div class="flex gap-2 justify-center">
        <Button
          @click="retryFetchData"
          variant="outline">
          Coba Lagi
        </Button>
        <Button @click="handleCreate">
          <Plus class="h-4 w-4 mr-2" />
          Tambah Balita Pertama
        </Button>
      </div>
    </div>

    <!-- Data Table -->
    <Card v-else>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="flex items-center gap-2">
              <Baby class="h-5 w-5" />
              Daftar Balita
            </CardTitle>
            <CardDescription>
              Data balita yang terdaftar dalam sistem pemantauan stunting. Termasuk informasi orang
              tua, data lahir, dan lokasi.
            </CardDescription>
          </div>
          <div
            v-if="isLoading"
            class="flex items-center gap-2 text-sm text-muted-foreground">
            <Loader2 class="h-4 w-4 animate-spin" />
            Memperbarui...
          </div>
        </div>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTable :data="balitaData" />
      </CardContent>
    </Card>

    <!-- Dialog Form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :balita="selectedBalita"
      :loading="isDialogLoading"
      @close="showDialog = false"
      @save="handleSave" />

    <!-- Development Info -->
    <Card class="border-yellow-200 bg-yellow-50">
      <CardHeader>
        <CardTitle class="text-sm text-yellow-800">🔧 Development Info</CardTitle>
      </CardHeader>
      <CardContent class="text-xs text-yellow-700 space-y-1">
        <div><strong>API Endpoints:</strong></div>
        <div>• GET /api/admin/balita/get - Fetch all data</div>
        <div>• POST /api/admin/balita/insert - Create new</div>
        <div>• PUT /api/admin/balita/update - Update existing</div>
        <div>• DELETE /api/admin/balita/delete - Soft delete</div>
        <div class="pt-2"><strong>Current Status:</strong></div>
        <div>• Loading: {{ isLoading }}</div>
        <div>• Total Data: {{ balitaData.length }}</div>
        <div>• Auth Token: {{ authUtils.getToken() ? "✅ Valid" : "❌ Missing" }}</div>
        <div>• Average Age: {{ formatAverageAge(averageAge) }}</div>
      </CardContent>
    </Card>
  </div>
</template>
