<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Users, Loader2, RefreshCcw } from "lucide-vue-next"
import { authUtils } from "@/lib/utils"
import type { Keluarga } from "./columns"

// API Response Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface GetAllKeluargaResponse {
  data: Keluarga[]
  total: number
}

interface InsertKeluargaResponse {
  id: string
}

interface UpdateDeleteKeluargaResponse {
  id: string
  message: string
}

// State management
const keluargaData = ref<Keluarga[]>([])
const isLoading = ref(false)
const isDialogLoading = ref(false)

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedKeluarga = ref<Keluarga | null>(null)

// Statistics
const totalKeluarga = ref(0)
const totalWithCoordinates = ref(0)

// API Functions
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

// Get all keluarga from API
const fetchKeluargaData = async () => {
  isLoading.value = true
  try {
    console.log("Fetching keluarga data from API...")

    const response = await apiRequest<GetAllKeluargaResponse>("/admin/keluarga/get")

    keluargaData.value = response.data.data
    totalKeluarga.value = response.data.total

    // Calculate coordinates statistics
    totalWithCoordinates.value = keluargaData.value.filter((k) => {
      if (!k.koordinat || !Array.isArray(k.koordinat)) return false
      const [lat, lng] = k.koordinat
      return lat !== 0 && lng !== 0
    }).length

    console.log("Keluarga data fetched successfully:", {
      total: totalKeluarga.value,
      withCoordinates: totalWithCoordinates.value,
    })

    toast.success(`Data keluarga berhasil dimuat (${totalKeluarga.value} data)`)
  } catch (error) {
    console.error("Error fetching keluarga data:", error)
    toast.error("Gagal memuat data keluarga. Silakan coba lagi.")

    // Fallback to empty data
    keluargaData.value = []
    totalKeluarga.value = 0
    totalWithCoordinates.value = 0
  } finally {
    isLoading.value = false
  }
}

// Create new keluarga
const createKeluarga = async (keluargaPayload: Partial<Keluarga>) => {
  try {
    console.log("Creating new keluarga:", keluargaPayload)

    const payload = {
      nomor_kk: keluargaPayload.nomor_kk,
      nama_ayah: keluargaPayload.nama_ayah,
      nama_ibu: keluargaPayload.nama_ibu,
      nik_ayah: keluargaPayload.nik_ayah,
      nik_ibu: keluargaPayload.nik_ibu,
      alamat: keluargaPayload.alamat,
      rt: keluargaPayload.rt,
      rw: keluargaPayload.rw,
      id_kelurahan: keluargaPayload.id_kelurahan,
      koordinat: keluargaPayload.koordinat,
    }

    const response = await apiRequest<InsertKeluargaResponse>("/admin/keluarga/insert", {
      method: "POST",
      body: JSON.stringify(payload),
    })

    console.log("Keluarga created successfully:", response)
    toast.success("Data keluarga berhasil ditambahkan")

    // Refresh data
    await fetchKeluargaData()

    return response.data
  } catch (error) {
    console.error("Error creating keluarga:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Data sudah ada. Periksa Nomor KK atau NIK.")
      } else {
        toast.error(`Gagal menambah data keluarga: ${error.message}`)
      }
    } else {
      toast.error("Gagal menambah data keluarga. Silakan coba lagi.")
    }
    throw error
  }
}

// Update existing keluarga
const updateKeluarga = async (keluargaPayload: Keluarga) => {
  try {
    console.log("Updating keluarga:", keluargaPayload)

    const payload = {
      id: keluargaPayload.id,
      nomor_kk: keluargaPayload.nomor_kk,
      nama_ayah: keluargaPayload.nama_ayah,
      nama_ibu: keluargaPayload.nama_ibu,
      nik_ayah: keluargaPayload.nik_ayah,
      nik_ibu: keluargaPayload.nik_ibu,
      alamat: keluargaPayload.alamat,
      rt: keluargaPayload.rt,
      rw: keluargaPayload.rw,
      id_kelurahan: keluargaPayload.id_kelurahan,
      koordinat: keluargaPayload.koordinat,
    }

    const response = await apiRequest<UpdateDeleteKeluargaResponse>("/admin/keluarga/update", {
      method: "PUT",
      body: JSON.stringify(payload),
    })

    console.log("Keluarga updated successfully:", response)
    toast.success("Data keluarga berhasil diperbarui")

    // Refresh data
    await fetchKeluargaData()

    return response.data
  } catch (error) {
    console.error("Error updating keluarga:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Data sudah ada. Periksa Nomor KK atau NIK.")
      } else {
        toast.error(`Gagal memperbarui data keluarga: ${error.message}`)
        showDialog.value = false
      }
    } else {
      toast.error("Gagal memperbarui data keluarga. Silakan coba lagi.")
    }
    throw error
  }
}

// Delete keluarga (soft delete)
const deleteKeluarga = async (keluargaId: string) => {
  try {
    console.log("Deleting keluarga:", keluargaId)

    const response = await apiRequest<UpdateDeleteKeluargaResponse>("/admin/keluarga/delete", {
      method: "DELETE",
      body: JSON.stringify({ id: keluargaId }),
    })

    console.log("Keluarga deleted successfully:", response)
    toast.success("Data keluarga berhasil dihapus")

    // Refresh data
    await fetchKeluargaData()

    return response.data
  } catch (error) {
    console.error("Error deleting keluarga:", error)
    if (error instanceof Error) {
      if (error.message.includes("active balita")) {
        toast.error("Tidak dapat menghapus keluarga yang masih memiliki data balita aktif.")
      } else if (error.message.includes("active laporan")) {
        toast.error("Tidak dapat menghapus keluarga yang masih memiliki laporan aktif.")
      } else {
        toast.error(`Gagal menghapus data keluarga: ${error.message}`)
      }
    } else {
      toast.error("Gagal menghapus data keluarga. Silakan coba lagi.")
    }
    throw error
  }
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedKeluarga.value = null
  showDialog.value = true
}

const handleEdit = (keluarga: Keluarga) => {
  dialogMode.value = "edit"
  selectedKeluarga.value = { ...keluarga }
  showDialog.value = true
}

const handleDelete = async (keluarga: Keluarga) => {
  // Show confirmation dialog
  const isConfirmed = confirm(
    `Apakah Anda yakin ingin menghapus data keluarga ${keluarga.nama_ayah}?\n\n` +
      `Detail:\n` +
      `• Nomor KK: ${keluarga.nomor_kk}\n` +
      `• Nama Ayah: ${keluarga.nama_ayah}\n` +
      `• Nama Ibu: ${keluarga.nama_ibu}\n\n` +
      `Data akan dihapus secara permanen dari sistem.`
  )

  if (!isConfirmed) return

  isLoading.value = true
  try {
    await deleteKeluarga(keluarga.id)
  } catch (error) {
    // Error handling already done in deleteKeluarga function
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (keluarga: Keluarga | Partial<Keluarga>) => {
  isDialogLoading.value = true

  try {
    if (dialogMode.value === "create") {
      await createKeluarga(keluarga)
    } else {
      await updateKeluarga(keluarga as Keluarga)
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

  if (event.type === "edit-keluarga") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-keluarga") {
    await handleDelete(customEvent.detail)
  }
}

// Lifecycle hooks
onMounted(async () => {
  // Add event listeners for custom events from DataTable
  document.addEventListener("edit-keluarga", handleCustomEvents)
  document.addEventListener("delete-keluarga", handleCustomEvents)

  // Fetch initial data
  await fetchKeluargaData()
})

onUnmounted(() => {
  document.removeEventListener("edit-keluarga", handleCustomEvents)
  document.removeEventListener("delete-keluarga", handleCustomEvents)
})

// Refresh data function (can be called manually)
const refreshData = async () => {
  await fetchKeluargaData()
}

// Error retry function
const retryFetchData = async () => {
  toast.info("Mencoba mengambil data ulang...")
  await fetchKeluargaData()
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Data Keluarga</h1>
        <p class="text-gray-600">Kelola data keluarga balita di Kota Cirebon</p>
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
          <RefreshCcw class="h-4 w-4" />
          Refresh
        </Button>

        <!-- Add Button -->
        <Button
          @click="handleCreate"
          :disabled="isLoading"
          class="gap-2">
          <Plus class="h-4 w-4" />
          Tambah Keluarga
        </Button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Keluarga</CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : totalKeluarga }}
          </div>
          <p class="text-xs text-muted-foreground">Keluarga terdaftar</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Dengan Koordinat</CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : totalWithCoordinates }}
          </div>
          <p class="text-xs text-muted-foreground">Lokasi terverifikasi</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Persentase Koordinat</CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{
              isLoading
                ? "..."
                : totalKeluarga > 0
                ? Math.round((totalWithCoordinates / totalKeluarga) * 100) + "%"
                : "0%"
            }}
          </div>
          <p class="text-xs text-muted-foreground">Data koordinat lengkap</p>
        </CardContent>
      </Card>
    </div>

    <!-- Loading State -->
    <div
      v-if="isLoading && keluargaData.length === 0"
      class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">Memuat Data Keluarga</h3>
        <p class="text-gray-600">Sedang mengambil data dari server...</p>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="!isLoading && keluargaData.length === 0"
      class="text-center py-12">
      <div class="mx-auto mb-4 w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
        <Users class="h-8 w-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Tidak Ada Data Keluarga</h3>
      <p class="text-gray-600 mb-4">Belum ada data keluarga yang terdaftar dalam sistem.</p>
      <div class="flex gap-2 justify-center">
        <Button
          @click="retryFetchData"
          variant="outline">
          Coba Lagi
        </Button>
        <Button @click="handleCreate">
          <Plus class="h-4 w-4 mr-2" />
          Tambah Keluarga Pertama
        </Button>
      </div>
    </div>

    <!-- Data Table -->
    <Card v-else>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle>Daftar Keluarga</CardTitle>
            <CardDescription>
              Data keluarga yang terdaftar dalam sistem pemantauan stunting
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
        <DataTable :data="keluargaData" />
      </CardContent>
    </Card>

    <!-- Dialog Form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :keluarga="selectedKeluarga"
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
        <div>• GET /api/admin/keluarga/get - Fetch all data</div>
        <div>• POST /api/admin/keluarga/insert - Create new</div>
        <div>• PUT /api/admin/keluarga/update - Update existing</div>
        <div>• DELETE /api/admin/keluarga/delete - Soft delete</div>
        <div class="pt-2"><strong>Current Status:</strong></div>
        <div>• Loading: {{ isLoading }}</div>
        <div>• Total Data: {{ keluargaData.length }}</div>
        <div>• Auth Token: {{ authUtils.getToken() ? "✅ Valid" : "❌ Missing" }}</div>
      </CardContent>
    </Card>
  </div>
</template>
