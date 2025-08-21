<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Activity, Plus, Stethoscope, Heart, Users, Loader2, RefreshCcw } from "lucide-vue-next"
import { authUtils } from "@/lib/utils"
import type { Intervensi } from "./columns"

// API Response Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface GetAllIntervensiResponse {
  data: Intervensi[]
  total: number
}

interface InsertIntervensiResponse {
  id: string
}

interface UpdateDeleteIntervensiResponse {
  id: string
  message: string
}

// State management
const intervensiData = ref<Intervensi[]>([])
const isLoading = ref(false)
const isDialogLoading = ref(false)

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedIntervensi = ref<Intervensi | null>(null)

// Statistics
const totalIntervensi = ref(0)
const totalGizi = ref(0)
const totalKesehatan = ref(0)
const totalSosial = ref(0)
const totalWithPetugas = ref(0)
const totalWithRiwayat = ref(0)

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

// Get all intervensi from API
const fetchIntervensiData = async () => {
  isLoading.value = true
  try {
    console.log("Fetching intervensi data from API...")

    const response = await apiRequest<GetAllIntervensiResponse>("/admin/intervensi/get")

    intervensiData.value = response.data.data
    totalIntervensi.value = response.data.total

    // Calculate statistics
    updateStatistics()

    console.log("Intervensi data fetched successfully:", {
      total: totalIntervensi.value,
      gizi: totalGizi.value,
      kesehatan: totalKesehatan.value,
      sosial: totalSosial.value,
    })

    toast.success(`Data intervensi berhasil dimuat (${totalIntervensi.value} data)`)
  } catch (error) {
    console.error("Error fetching intervensi data:", error)
    toast.error("Gagal memuat data intervensi. Silakan coba lagi.")

    // Fallback to empty data
    intervensiData.value = []
    totalIntervensi.value = 0
    updateStatistics()
  } finally {
    isLoading.value = false
  }
}

// Calculate statistics
const updateStatistics = () => {
  totalIntervensi.value = intervensiData.value.length

  // Jenis-based statistics
  totalGizi.value = intervensiData.value.filter((i) => i.jenis === "gizi").length
  totalKesehatan.value = intervensiData.value.filter((i) => i.jenis === "kesehatan").length
  totalSosial.value = intervensiData.value.filter((i) => i.jenis === "sosial").length

  // Additional statistics
  totalWithPetugas.value = intervensiData.value.filter((i) => i.petugas_count > 0).length
  totalWithRiwayat.value = intervensiData.value.filter((i) => i.riwayat_count > 0).length
}

// Create new intervensi
const createIntervensi = async (intervensiPayload: Partial<Intervensi>) => {
  try {
    console.log("Creating new intervensi:", intervensiPayload)

    const payload = {
      id_balita: intervensiPayload.id_balita,
      jenis: intervensiPayload.jenis,
      tanggal: intervensiPayload.tanggal,
      deskripsi: intervensiPayload.deskripsi,
      hasil: intervensiPayload.hasil,
    }

    const response = await apiRequest<InsertIntervensiResponse>("/admin/intervensi/insert", {
      method: "POST",
      body: JSON.stringify(payload),
    })

    console.log("Intervensi created successfully:", response)
    toast.success("Data intervensi berhasil ditambahkan")

    // Refresh data
    await fetchIntervensiData()

    return response.data
  } catch (error) {
    console.error("Error creating intervensi:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Intervensi dengan balita, jenis, tanggal, dan deskripsi yang sama sudah ada.")
      } else if (error.message.includes("Balita not found")) {
        toast.error("Data balita tidak ditemukan. Pastikan balita masih aktif.")
      } else if (error.message.includes("future")) {
        toast.error("Tanggal intervensi tidak boleh di masa depan.")
      } else if (error.message.includes("10-1000 characters")) {
        toast.error("Deskripsi harus antara 10-1000 karakter.")
      } else if (error.message.includes("5-500 characters")) {
        toast.error("Hasil harus antara 5-500 karakter.")
      } else {
        toast.error(`Gagal menambah data intervensi: ${error.message}`)
      }
    } else {
      toast.error("Gagal menambah data intervensi. Silakan coba lagi.")
    }
    throw error
  }
}

// Update existing intervensi
const updateIntervensi = async (intervensiPayload: Intervensi) => {
  try {
    console.log("Updating intervensi:", intervensiPayload)

    const payload = {
      id: intervensiPayload.id,
      id_balita: intervensiPayload.id_balita,
      jenis: intervensiPayload.jenis,
      tanggal: intervensiPayload.tanggal,
      deskripsi: intervensiPayload.deskripsi,
      hasil: intervensiPayload.hasil,
    }

    const response = await apiRequest<UpdateDeleteIntervensiResponse>("/admin/intervensi/update", {
      method: "PUT",
      body: JSON.stringify(payload),
    })

    console.log("Intervensi updated successfully:", response)
    toast.success("Data intervensi berhasil diperbarui")

    // Refresh data
    await fetchIntervensiData()

    return response.data
  } catch (error) {
    console.error("Error updating intervensi:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Intervensi dengan balita, jenis, tanggal, dan deskripsi yang sama sudah ada.")
      } else if (error.message.includes("Balita not found")) {
        toast.error("Data balita tidak ditemukan. Pastikan balita masih aktif.")
      } else if (error.message.includes("future")) {
        toast.error("Tanggal intervensi tidak boleh di masa depan.")
      } else if (error.message.includes("10-1000 characters")) {
        toast.error("Deskripsi harus antara 10-1000 karakter.")
      } else if (error.message.includes("5-500 characters")) {
        toast.error("Hasil harus antara 5-500 karakter.")
      } else {
        toast.error(`Gagal memperbarui data intervensi: ${error.message}`)
        showDialog.value = false
      }
    } else {
      toast.error("Gagal memperbarui data intervensi. Silakan coba lagi.")
    }
    throw error
  }
}

// Delete intervensi (soft delete)
const deleteIntervensi = async (intervensiId: string) => {
  try {
    console.log("Deleting intervensi:", intervensiId)

    const response = await apiRequest<UpdateDeleteIntervensiResponse>("/admin/intervensi/delete", {
      method: "DELETE",
      body: JSON.stringify({ id: intervensiId }),
    })

    console.log("Intervensi deleted successfully:", response)
    toast.success("Data intervensi berhasil dihapus")

    // Refresh data
    await fetchIntervensiData()

    return response.data
  } catch (error) {
    console.error("Error deleting intervensi:", error)
    if (error instanceof Error) {
      if (error.message.includes("petugas assigned")) {
        toast.error("Tidak dapat menghapus intervensi yang masih memiliki petugas yang ditugaskan.")
      } else if (error.message.includes("riwayat pemeriksaan")) {
        toast.error(
          "Tidak dapat menghapus intervensi yang masih memiliki riwayat pemeriksaan aktif."
        )
      } else if (error.message.includes("already deleted")) {
        toast.error("Intervensi sudah dihapus sebelumnya.")
      } else {
        toast.error(`Gagal menghapus data intervensi: ${error.message}`)
      }
    } else {
      toast.error("Gagal menghapus data intervensi. Silakan coba lagi.")
    }
    throw error
  }
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedIntervensi.value = null
  showDialog.value = true
}

const handleEdit = (intervensi: Intervensi) => {
  dialogMode.value = "edit"
  selectedIntervensi.value = { ...intervensi }
  showDialog.value = true
}

const handleDelete = async (intervensi: Intervensi) => {
  // Show confirmation dialog
  const isConfirmed = confirm(
    `Apakah Anda yakin ingin menghapus intervensi ini?\n\n` +
      `Detail:\n` +
      `• Balita: ${intervensi.nama_balita}\n` +
      `• Jenis: ${intervensi.jenis}\n` +
      `• Tanggal: ${intervensi.tanggal}\n` +
      `• Deskripsi: ${intervensi.deskripsi?.substring(0, 50)}${
        intervensi.deskripsi?.length > 50 ? "..." : ""
      }\n` +
      `• Petugas Assigned: ${intervensi.petugas_count}\n` +
      `• Riwayat Pemeriksaan: ${intervensi.riwayat_count}\n\n` +
      `Data akan dihapus secara permanen dari sistem.`
  )

  if (!isConfirmed) return

  isLoading.value = true
  try {
    await deleteIntervensi(intervensi.id)
  } catch (error) {
    // Error handling already done in deleteIntervensi function
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (intervensi: Intervensi | Partial<Intervensi>) => {
  isDialogLoading.value = true

  try {
    if (dialogMode.value === "create") {
      await createIntervensi(intervensi)
    } else {
      await updateIntervensi(intervensi as Intervensi)
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

  if (event.type === "edit-intervensi") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-intervensi") {
    await handleDelete(customEvent.detail)
  }
}

// Lifecycle hooks
onMounted(async () => {
  // Add event listeners for custom events from DataTable
  document.addEventListener("edit-intervensi", handleCustomEvents)
  document.addEventListener("delete-intervensi", handleCustomEvents)

  // Fetch initial data
  await fetchIntervensiData()
})

onUnmounted(() => {
  document.removeEventListener("edit-intervensi", handleCustomEvents)
  document.removeEventListener("delete-intervensi", handleCustomEvents)
})

// Refresh data function (can be called manually)
const refreshData = async () => {
  await fetchIntervensiData()
}

// Error retry function
const retryFetchData = async () => {
  toast.info("Mencoba mengambil data ulang...")
  await fetchIntervensiData()
}

// Helper function to get jenis percentage
const getJenisPercentage = (count: number): number => {
  return totalIntervensi.value > 0 ? Math.round((count / totalIntervensi.value) * 100) : 0
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Activity class="h-8 w-8 text-blue-600" />
          Data Intervensi
        </h1>
        <p class="text-gray-600">Kelola data intervensi stunting untuk balita di Kota Cirebon</p>
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
          Tambah Intervensi
        </Button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-6 gap-4">
      <!-- Total Intervensi -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Intervensi</CardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : totalIntervensi }}
          </div>
          <p class="text-xs text-muted-foreground">Intervensi dilakukan</p>
        </CardContent>
      </Card>

      <!-- Intervensi Gizi -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Gizi</CardTitle>
          <div class="text-2xl">🥗</div>
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">
            {{ isLoading ? "..." : totalGizi }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisPercentage(totalGizi) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Intervensi Kesehatan -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Kesehatan</CardTitle>
          <Stethoscope class="h-4 w-4 text-red-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-red-600">
            {{ isLoading ? "..." : totalKesehatan }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisPercentage(totalKesehatan) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Intervensi Sosial -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Sosial</CardTitle>
          <Users class="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600">
            {{ isLoading ? "..." : totalSosial }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisPercentage(totalSosial) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- With Petugas -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Ada Petugas</CardTitle>
          <Users class="h-4 w-4 text-purple-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-purple-600">
            {{ isLoading ? "..." : totalWithPetugas }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisPercentage(totalWithPetugas) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- With Riwayat -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Ada Riwayat</CardTitle>
          <Heart class="h-4 w-4 text-pink-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-pink-600">
            {{ isLoading ? "..." : totalWithRiwayat }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisPercentage(totalWithRiwayat) }}% dari total
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Loading State -->
    <div
      v-if="isLoading && intervensiData.length === 0"
      class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">Memuat Data Intervensi</h3>
        <p class="text-gray-600">Sedang mengambil data dari server...</p>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="!isLoading && intervensiData.length === 0"
      class="text-center py-12">
      <div class="mx-auto mb-4 w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
        <Activity class="h-8 w-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Tidak Ada Data Intervensi</h3>
      <p class="text-gray-600 mb-4">Belum ada data intervensi yang terdaftar dalam sistem.</p>
      <div class="flex gap-2 justify-center">
        <Button
          @click="retryFetchData"
          variant="outline">
          Coba Lagi
        </Button>
        <Button @click="handleCreate">
          <Plus class="h-4 w-4 mr-2" />
          Tambah Intervensi Pertama
        </Button>
      </div>
    </div>

    <!-- Data Table -->
    <Card v-else>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="flex items-center gap-2">
              <Activity class="h-5 w-5" />
              Daftar Intervensi
            </CardTitle>
            <CardDescription>
              Data intervensi yang dilakukan terhadap balita stunting, termasuk informasi balita,
              jenis intervensi, tanggal pelaksanaan, deskripsi, dan hasil yang dicapai.
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
        <DataTable :data="intervensiData" />
      </CardContent>
    </Card>

    <!-- Dialog Form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :intervensi="selectedIntervensi"
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
        <div>• GET /api/admin/intervensi/get - Fetch all data</div>
        <div>• POST /api/admin/intervensi/insert - Create new</div>
        <div>• PUT /api/admin/intervensi/update - Update existing</div>
        <div>• DELETE /api/admin/intervensi/delete - Soft delete</div>
        <div class="pt-2"><strong>Master Data Endpoints:</strong></div>
        <div>• GET /api/admin/balita/get - Balita dropdown</div>
        <div class="pt-2"><strong>Current Status:</strong></div>
        <div>• Loading: {{ isLoading }}</div>
        <div>• Total Data: {{ intervensiData.length }}</div>
        <div>• Auth Token: {{ authUtils.getToken() ? "✅ Valid" : "❌ Missing" }}</div>
        <div>
          • Gizi: {{ totalGizi }}, Kesehatan: {{ totalKesehatan }}, Sosial: {{ totalSosial }}
        </div>
        <div>• With Petugas: {{ totalWithPetugas }}, With Riwayat: {{ totalWithRiwayat }}</div>
      </CardContent>
    </Card>
  </div>
</template>
