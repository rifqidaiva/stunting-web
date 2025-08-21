<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, FileText, Users, AlertTriangle, CheckCircle, Clock, Loader2, RefreshCcw } from "lucide-vue-next"
import { authUtils } from "@/lib/utils"
import type { LaporanMasyarakat } from "./columns"

// API Response Types
interface ApiResponse<T> {
  data: T
  message: string
  status_code: number
}

interface GetAllLaporanMasyarakatResponse {
  data: LaporanMasyarakat[]
  total: number
}

interface InsertLaporanMasyarakatResponse {
  id: string
}

interface UpdateDeleteLaporanMasyarakatResponse {
  id: string
  message: string
}

// State management
const laporanData = ref<LaporanMasyarakat[]>([])
const isLoading = ref(false)
const isDialogLoading = ref(false)

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedLaporan = ref<LaporanMasyarakat | null>(null)

// Statistics
const totalLaporan = ref(0)
const totalBelumDiproses = ref(0)
const totalSedangDiproses = ref(0)
const totalSelesai = ref(0)
const totalLaporanAdmin = ref(0)
const totalLaporanMasyarakat = ref(0)

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

// Get all laporan masyarakat from API
const fetchLaporanData = async () => {
  isLoading.value = true
  try {
    console.log("Fetching laporan masyarakat data from API...")
    
    const response = await apiRequest<GetAllLaporanMasyarakatResponse>("/admin/laporan-masyarakat/get")
    
    laporanData.value = response.data.data
    totalLaporan.value = response.data.total
    
    // Calculate statistics
    updateStatistics()

    console.log("Laporan masyarakat data fetched successfully:", {
      total: totalLaporan.value,
      belumDiproses: totalBelumDiproses.value,
      sedangDiproses: totalSedangDiproses.value,
      selesai: totalSelesai.value
    })

    toast.success(`Data laporan masyarakat berhasil dimuat (${totalLaporan.value} data)`)
  } catch (error) {
    console.error("Error fetching laporan masyarakat data:", error)
    toast.error("Gagal memuat data laporan masyarakat. Silakan coba lagi.")
    
    // Fallback to empty data
    laporanData.value = []
    totalLaporan.value = 0
    updateStatistics()
  } finally {
    isLoading.value = false
  }
}

// Calculate statistics
const updateStatistics = () => {
  totalLaporan.value = laporanData.value.length
  
  // Status-based statistics
  totalBelumDiproses.value = laporanData.value.filter((l) => 
    l.status_laporan === "Belum diproses"
  ).length
  
  totalSedangDiproses.value = laporanData.value.filter((l) => 
    ["Sedang diproses", "Belum ditindaklanjuti", "Sudah ditindaklanjuti"].includes(l.status_laporan)
  ).length
  
  totalSelesai.value = laporanData.value.filter((l) => 
    ["Diproses dan data sesuai", "Sudah perbaikan gizi"].includes(l.status_laporan)
  ).length
  
  // Jenis laporan statistics
  totalLaporanAdmin.value = laporanData.value.filter((l) => 
    l.jenis_laporan === "admin"
  ).length
  
  totalLaporanMasyarakat.value = laporanData.value.filter((l) => 
    l.jenis_laporan === "masyarakat"
  ).length
}

// Create new laporan masyarakat
const createLaporanMasyarakat = async (laporanPayload: Partial<LaporanMasyarakat>) => {
  try {
    console.log("Creating new laporan masyarakat:", laporanPayload)

    // Handle admin laporan (set id_masyarakat to null if admin)
    const payload = {
      id_masyarakat: laporanPayload.id_masyarakat === "ADMIN" ? null : laporanPayload.id_masyarakat,
      id_balita: laporanPayload.id_balita,
      id_status_laporan: laporanPayload.id_status_laporan,
      tanggal_laporan: laporanPayload.tanggal_laporan,
      hubungan_dengan_balita: laporanPayload.hubungan_dengan_balita,
      nomor_hp_pelapor: laporanPayload.nomor_hp_pelapor,
      nomor_hp_keluarga_balita: laporanPayload.nomor_hp_keluarga_balita,
    }

    const response = await apiRequest<InsertLaporanMasyarakatResponse>("/admin/laporan-masyarakat/insert", {
      method: "POST",
      body: JSON.stringify(payload),
    })

    console.log("Laporan masyarakat created successfully:", response)
    toast.success("Data laporan masyarakat berhasil ditambahkan")

    // Refresh data
    await fetchLaporanData()
    
    return response.data
  } catch (error) {
    console.error("Error creating laporan masyarakat:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Laporan dengan balita dan tanggal yang sama sudah ada dari pelapor ini.")
      } else if (error.message.includes("Balita not found")) {
        toast.error("Data balita tidak ditemukan. Pastikan balita masih aktif.")
      } else if (error.message.includes("Masyarakat not found")) {
        toast.error("Data masyarakat tidak ditemukan. Pastikan pelapor terdaftar.")
      } else if (error.message.includes("Status laporan not found")) {
        toast.error("Status laporan tidak valid. Silakan pilih status yang benar.")
      } else if (error.message.includes("future")) {
        toast.error("Tanggal laporan tidak boleh di masa depan.")
      } else {
        toast.error(`Gagal menambah data laporan masyarakat: ${error.message}`)
      }
    } else {
      toast.error("Gagal menambah data laporan masyarakat. Silakan coba lagi.")
    }
    throw error
  }
}

// Update existing laporan masyarakat
const updateLaporanMasyarakat = async (laporanPayload: LaporanMasyarakat) => {
  try {
    console.log("Updating laporan masyarakat:", laporanPayload)

    // Handle admin laporan (set id_masyarakat to null if admin)
    const payload = {
      id: laporanPayload.id,
      id_masyarakat: laporanPayload.id_masyarakat === "ADMIN" ? null : laporanPayload.id_masyarakat,
      id_balita: laporanPayload.id_balita,
      id_status_laporan: laporanPayload.id_status_laporan,
      tanggal_laporan: laporanPayload.tanggal_laporan,
      hubungan_dengan_balita: laporanPayload.hubungan_dengan_balita,
      nomor_hp_pelapor: laporanPayload.nomor_hp_pelapor,
      nomor_hp_keluarga_balita: laporanPayload.nomor_hp_keluarga_balita,
    }

    const response = await apiRequest<UpdateDeleteLaporanMasyarakatResponse>("/admin/laporan-masyarakat/update", {
      method: "PUT",
      body: JSON.stringify(payload),
    })

    console.log("Laporan masyarakat updated successfully:", response)
    toast.success("Data laporan masyarakat berhasil diperbarui")

    // Refresh data
    await fetchLaporanData()
    
    return response.data
  } catch (error) {
    console.error("Error updating laporan masyarakat:", error)
    if (error instanceof Error) {
      if (error.message.includes("already exists")) {
        toast.error("Laporan dengan balita dan tanggal yang sama sudah ada dari pelapor ini.")
      } else if (error.message.includes("Balita not found")) {
        toast.error("Data balita tidak ditemukan. Pastikan balita masih aktif.")
      } else if (error.message.includes("Masyarakat not found")) {
        toast.error("Data masyarakat tidak ditemukan. Pastikan pelapor terdaftar.")
      } else if (error.message.includes("Status laporan not found")) {
        toast.error("Status laporan tidak valid. Silakan pilih status yang benar.")
      } else if (error.message.includes("future")) {
        toast.error("Tanggal laporan tidak boleh di masa depan.")
      } else {
        toast.error(`Gagal memperbarui data laporan masyarakat: ${error.message}`)
        showDialog.value = false
      }
    } else {
      toast.error("Gagal memperbarui data laporan masyarakat. Silakan coba lagi.")
    }
    throw error
  }
}

// Delete laporan masyarakat (soft delete)
const deleteLaporanMasyarakat = async (laporanId: string) => {
  try {
    console.log("Deleting laporan masyarakat:", laporanId)

    const response = await apiRequest<UpdateDeleteLaporanMasyarakatResponse>("/admin/laporan-masyarakat/delete", {
      method: "DELETE",
      body: JSON.stringify({ id: laporanId }),
    })

    console.log("Laporan masyarakat deleted successfully:", response)
    toast.success("Data laporan masyarakat berhasil dihapus")

    // Refresh data
    await fetchLaporanData()
    
    return response.data
  } catch (error) {
    console.error("Error deleting laporan masyarakat:", error)
    if (error instanceof Error) {
      if (error.message.includes("processed")) {
        toast.error("Tidak dapat menghapus laporan yang sudah diproses.")
      } else if (error.message.includes("already deleted")) {
        toast.error("Laporan sudah dihapus sebelumnya.")
      } else {
        toast.error(`Gagal menghapus data laporan masyarakat: ${error.message}`)
      }
    } else {
      toast.error("Gagal menghapus data laporan masyarakat. Silakan coba lagi.")
    }
    throw error
  }
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedLaporan.value = null
  showDialog.value = true
}

const handleEdit = (laporan: LaporanMasyarakat) => {
  dialogMode.value = "edit"
  selectedLaporan.value = { ...laporan }
  showDialog.value = true
}

const handleDelete = async (laporan: LaporanMasyarakat) => {
  // Determine laporan type for display
  const laporanType = laporan.jenis_laporan === "admin" ? "Administratif" : "Masyarakat"
  const pelaporInfo = laporan.jenis_laporan === "admin" ? "Admin Sistem" : laporan.nama_pelapor

  // Show confirmation dialog
  const isConfirmed = confirm(
    `Apakah Anda yakin ingin menghapus laporan ${laporanType} ini?\n\n` +
    `Detail:\n` +
    `• ID Laporan: ${laporan.id}\n` +
    `• Pelapor: ${pelaporInfo}\n` +
    `• Balita: ${laporan.nama_balita}\n` +
    `• Status: ${laporan.status_laporan}\n` +
    `• Tanggal: ${laporan.tanggal_laporan}\n\n` +
    `Data akan dihapus secara permanen dari sistem.`
  )

  if (!isConfirmed) return

  isLoading.value = true
  try {
    await deleteLaporanMasyarakat(laporan.id)
  } catch (error) {
    // Error handling already done in deleteLaporanMasyarakat function
  } finally {
    isLoading.value = false
  }
}

const handleSave = async (laporan: LaporanMasyarakat | Partial<LaporanMasyarakat>) => {
  isDialogLoading.value = true
  
  try {
    if (dialogMode.value === "create") {
      await createLaporanMasyarakat(laporan)
    } else {
      await updateLaporanMasyarakat(laporan as LaporanMasyarakat)
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

  if (event.type === "edit-laporan") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-laporan") {
    await handleDelete(customEvent.detail)
  }
}

// Lifecycle hooks
onMounted(async () => {
  // Add event listeners for custom events from DataTable
  document.addEventListener("edit-laporan", handleCustomEvents)
  document.addEventListener("delete-laporan", handleCustomEvents)

  // Fetch initial data
  await fetchLaporanData()
})

onUnmounted(() => {
  document.removeEventListener("edit-laporan", handleCustomEvents)
  document.removeEventListener("delete-laporan", handleCustomEvents)
})

// Refresh data function (can be called manually)
const refreshData = async () => {
  await fetchLaporanData()
}

// Error retry function
const retryFetchData = async () => {
  toast.info("Mencoba mengambil data ulang...")
  await fetchLaporanData()
}

// Helper function to get jenis laporan percentage
const getJenisLaporanPercentage = (count: number): number => {
  return totalLaporan.value > 0 ? Math.round((count / totalLaporan.value) * 100) : 0
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <FileText class="h-8 w-8 text-slate-600" />
          Data Laporan Masyarakat
        </h1>
        <p class="text-gray-600">
          Kelola laporan balita dari masyarakat dan administrator di Kota Cirebon
        </p>
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
          Tambah Laporan
        </Button>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-6 gap-4">
      <!-- Total Laporan -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Laporan</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ isLoading ? "..." : totalLaporan }}
          </div>
          <p class="text-xs text-muted-foreground">Laporan terdaftar</p>
        </CardContent>
      </Card>

      <!-- Belum Diproses -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Belum Diproses</CardTitle>
          <Clock class="h-4 w-4 text-gray-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-gray-600">
            {{ isLoading ? "..." : totalBelumDiproses }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisLaporanPercentage(totalBelumDiproses) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Sedang Diproses -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Sedang Diproses</CardTitle>
          <AlertTriangle class="h-4 w-4 text-orange-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-orange-600">
            {{ isLoading ? "..." : totalSedangDiproses }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisLaporanPercentage(totalSedangDiproses) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Selesai -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Selesai</CardTitle>
          <CheckCircle class="h-4 w-4 text-green-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-green-600">
            {{ isLoading ? "..." : totalSelesai }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisLaporanPercentage(totalSelesai) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Laporan Admin -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Laporan Admin</CardTitle>
          <Users class="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600">
            {{ isLoading ? "..." : totalLaporanAdmin }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisLaporanPercentage(totalLaporanAdmin) }}% dari total
          </p>
        </CardContent>
      </Card>

      <!-- Laporan Masyarakat -->
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Laporan Masyarakat</CardTitle>
          <Users class="h-4 w-4 text-purple-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-purple-600">
            {{ isLoading ? "..." : totalLaporanMasyarakat }}
          </div>
          <p class="text-xs text-muted-foreground">
            {{ getJenisLaporanPercentage(totalLaporanMasyarakat) }}% dari total
          </p>
        </CardContent>
      </Card>
    </div>

    <!-- Loading State -->
    <div
      v-if="isLoading && laporanData.length === 0"
      class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="h-8 w-8 animate-spin mx-auto mb-4 text-primary" />
        <h3 class="text-lg font-medium text-gray-900 mb-2">Memuat Data Laporan Masyarakat</h3>
        <p class="text-gray-600">Sedang mengambil data dari server...</p>
      </div>
    </div>

    <!-- Error State -->
    <div
      v-else-if="!isLoading && laporanData.length === 0"
      class="text-center py-12">
      <div class="mx-auto mb-4 w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center">
        <FileText class="h-8 w-8 text-gray-400" />
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">Tidak Ada Data Laporan Masyarakat</h3>
      <p class="text-gray-600 mb-4">
        Belum ada laporan masyarakat yang terdaftar dalam sistem.
      </p>
      <div class="flex gap-2 justify-center">
        <Button
          @click="retryFetchData"
          variant="outline">
          Coba Lagi
        </Button>
        <Button @click="handleCreate">
          <Plus class="h-4 w-4 mr-2" />
          Tambah Laporan Pertama
        </Button>
      </div>
    </div>

    <!-- Data Table -->
    <Card v-else>
      <CardHeader>
        <div class="flex items-center justify-between">
          <div>
            <CardTitle class="flex items-center gap-2">
              <FileText class="h-5 w-5" />
              Daftar Laporan Masyarakat
            </CardTitle>
            <CardDescription>
              Data laporan balita dari masyarakat dan administrator yang terdaftar dalam sistem pemantauan stunting.
              Termasuk informasi pelapor, balita, status penanganan, dan kontak.
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
        <DataTable :data="laporanData" />
      </CardContent>
    </Card>

    <!-- Dialog Form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :laporan="selectedLaporan"
      :loading="isDialogLoading"
      @close="showDialog = false"
      @save="handleSave" />

    <!-- Development Info -->
    <Card
      class="border-yellow-200 bg-yellow-50">
      <CardHeader>
        <CardTitle class="text-sm text-yellow-800">🔧 Development Info</CardTitle>
      </CardHeader>
      <CardContent class="text-xs text-yellow-700 space-y-1">
        <div><strong>API Endpoints:</strong></div>
        <div>• GET /api/admin/laporan-masyarakat/get - Fetch all data</div>
        <div>• POST /api/admin/laporan-masyarakat/insert - Create new</div>
        <div>• PUT /api/admin/laporan-masyarakat/update - Update existing</div>
        <div>• DELETE /api/admin/laporan-masyarakat/delete - Soft delete</div>
        <div class="pt-2"><strong>Master Data Endpoints:</strong></div>
        <div>• GET /api/admin/master-masyarakat - Masyarakat dropdown</div>
        <div>• GET /api/admin/balita/get - Balita dropdown</div>
        <div>• GET /api/admin/master-status-laporan - Status laporan dropdown</div>
        <div class="pt-2"><strong>Current Status:</strong></div>
        <div>• Loading: {{ isLoading }}</div>
        <div>• Total Data: {{ laporanData.length }}</div>
        <div>• Auth Token: {{ authUtils.getToken() ? "✅ Valid" : "❌ Missing" }}</div>
        <div>• Admin Reports: {{ totalLaporanAdmin }}</div>
        <div>• Community Reports: {{ totalLaporanMasyarakat }}</div>
      </CardContent>
    </Card>
  </div>
</template>
