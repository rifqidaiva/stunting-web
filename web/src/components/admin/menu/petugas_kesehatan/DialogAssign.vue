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
import { Separator } from "@/components/ui/separator"
import {
  UserPlus,
  X,
  Activity,
  Users,
  Building,
  Trash2,
  CheckCircle,
  AlertCircle,
  FileText,
} from "lucide-vue-next"
import type { PetugasKesehatan } from "./columns"

// Interface untuk intervensi yang tersedia
interface IntervensiOption {
  id: string
  nama_balita: string
  jenis: string
  tanggal: string
  deskripsi: string
  petugas_assigned: Array<{
    id: string // ID dari tabel intervensi_petugas
    id_petugas: string
    nama_petugas: string
    skpd: string
  }>
}

// Interface untuk assignment baru
interface NewAssignment {
  id_intervensi: string
  id_petugas_kesehatan: string
}

interface Props {
  show: boolean
  petugas: PetugasKesehatan | null
}

interface Emits {
  (e: "close"): void
  (e: "assign", assignment: NewAssignment): void
  (e: "remove", assignmentId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form state
const selectedIntervensi = ref<string>("")
const searchQuery = ref("")

// Validation errors
const errors = ref<Record<string, string>>({})

// Dummy data intervensi yang tersedia
const intervensiOptions = ref<IntervensiOption[]>([
  {
    id: "I001",
    nama_balita: "Ahmad Rizki (24 bulan)",
    jenis: "gizi",
    tanggal: "2024-02-15",
    deskripsi: "Pemberian Vitamin A dan makanan tambahan",
    petugas_assigned: [
      {
        id: "IP001",
        id_petugas: "P002",
        nama_petugas: "Dr. Sari Wijaya",
        skpd: "Puskesmas Kejaksan",
      },
    ],
  },
  {
    id: "I002",
    nama_balita: "Fatimah Zahra (18 bulan)",
    jenis: "kesehatan",
    tanggal: "2024-02-20",
    deskripsi: "Pemeriksaan kesehatan rutin dan konseling",
    petugas_assigned: [],
  },
  {
    id: "I003",
    nama_balita: "Zainab (30 bulan)",
    jenis: "sosial",
    tanggal: "2024-02-25",
    deskripsi: "Konseling gizi keluarga dan edukasi",
    petugas_assigned: [
      {
        id: "IP002",
        id_petugas: "P001",
        nama_petugas: "Bidan Rina",
        skpd: "Puskesmas Pekalangan",
      },
      {
        id: "IP003",
        id_petugas: "P003",
        nama_petugas: "Ahli Gizi Dewi",
        skpd: "Dinas Kesehatan",
      },
    ],
  },
  {
    id: "I004",
    nama_balita: "Muhammad Hakim (12 bulan)",
    jenis: "gizi",
    tanggal: "2024-03-01",
    deskripsi: "Program makanan tambahan dan suplementasi",
    petugas_assigned: [],
  },
  {
    id: "I005",
    nama_balita: "Aisyah Putri (36 bulan)",
    jenis: "kesehatan",
    tanggal: "2024-03-05",
    deskripsi: "Pemeriksaan berkala dan vaksinasi",
    petugas_assigned: [
      {
        id: "IP004",
        id_petugas: "P004",
        nama_petugas: "Dr. Ahmad Fauzi",
        skpd: "Puskesmas Harjamukti",
      },
    ],
  },
])

// Filtered intervensi berdasarkan pencarian
const filteredIntervensi = computed(() => {
  if (!searchQuery.value) return intervensiOptions.value

  const query = searchQuery.value.toLowerCase()
  return intervensiOptions.value.filter(
    (intervensi) =>
      intervensi.nama_balita.toLowerCase().includes(query) ||
      intervensi.jenis.toLowerCase().includes(query) ||
      intervensi.deskripsi.toLowerCase().includes(query)
  )
})

// Get selected intervensi info
const selectedIntervensiInfo = computed(() => {
  if (!selectedIntervensi.value) return null
  return intervensiOptions.value.find((i) => i.id === selectedIntervensi.value) || null
})

// Check if current petugas is already assigned to selected intervensi
const isAlreadyAssigned = computed(() => {
  if (!selectedIntervensiInfo.value || !props.petugas) return false
  return selectedIntervensiInfo.value.petugas_assigned.some(
    (assignment) => assignment.id_petugas === props.petugas!.id
  )
})

// Get current assignment ID if petugas is already assigned
const currentAssignmentId = computed(() => {
  if (!selectedIntervensiInfo.value || !props.petugas) return null
  const assignment = selectedIntervensiInfo.value.petugas_assigned.find(
    (assignment) => assignment.id_petugas === props.petugas!.id
  )
  return assignment ? assignment.id : null
})

// Helper functions
const getJenisColor = (jenis: string): string => {
  switch (jenis?.toLowerCase()) {
    case "gizi":
      return "bg-blue-100 text-blue-800 border-blue-200"
    case "kesehatan":
      return "bg-green-100 text-green-800 border-green-200"
    case "sosial":
      return "bg-purple-100 text-purple-800 border-purple-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  })
}

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  if (!selectedIntervensi.value) {
    errors.value.intervensi = "Intervensi harus dipilih"
    return false
  }

  if (isAlreadyAssigned.value) {
    errors.value.intervensi = "Petugas sudah ditugaskan pada intervensi ini"
    return false
  }

  return true
}

// Reset form
const resetForm = () => {
  selectedIntervensi.value = ""
  searchQuery.value = ""
  errors.value = {}
}

// Handle assign
const handleAssign = () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi")
    return
  }

  if (!props.petugas) {
    toast.error("Data petugas tidak valid")
    return
  }

  const assignment: NewAssignment = {
    id_intervensi: selectedIntervensi.value,
    id_petugas_kesehatan: props.petugas.id,
  }

  emit("assign", assignment)
  resetForm()
}

// Handle remove assignment
const handleRemove = () => {
  if (!currentAssignmentId.value) {
    toast.error("Assignment ID tidak ditemukan")
    return
  }

  if (confirm("Apakah Anda yakin ingin menghapus penugasan ini?")) {
    emit("remove", currentAssignmentId.value)
    resetForm()
  }
}

// Handle close
const handleClose = () => {
  resetForm()
  emit("close")
}

// Watch for dialog visibility
watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      resetForm()
    }
  }
)

// Real-time validation
watch(
  () => selectedIntervensi.value,
  () => {
    if (selectedIntervensi.value) {
      delete errors.value.intervensi
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
              <UserPlus class="h-5 w-5 text-blue-600" />
              Assign Petugas ke Intervensi
            </DialogTitle>
            <DialogDescription class="text-sm">
              Tugaskan petugas kesehatan ke intervensi tertentu atau kelola penugasan yang sudah
              ada.
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Info Petugas -->
            <Card
              v-if="petugas"
              class="bg-blue-50 border-blue-200">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base text-blue-900">
                  <Users class="h-4 w-4" />
                  Petugas yang Akan Ditugaskan
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="flex items-start gap-3">
                  <div class="bg-blue-100 rounded-full p-2">
                    <Users class="h-5 w-5 text-blue-600" />
                  </div>
                  <div class="flex-1">
                    <div class="font-medium text-blue-900">{{ petugas.nama }}</div>
                    <div class="text-sm text-blue-700">{{ petugas.email }}</div>
                    <div class="flex items-center gap-2 mt-2">
                      <Building class="h-4 w-4 text-blue-600" />
                      <span class="text-sm text-blue-800">{{ petugas.skpd }}</span>
                      <Badge
                        variant="outline"
                        :class="
                          petugas.jenis_skpd === 'puskesmas'
                            ? 'border-blue-300 text-blue-700'
                            : 'border-green-300 text-green-700'
                        ">
                        {{ petugas.jenis_skpd }}
                      </Badge>
                    </div>
                    <div class="text-xs text-blue-600 mt-1">
                      {{ petugas.intervensi_count }} intervensi aktif
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Search Intervensi -->
            <Card>
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Activity class="h-4 w-4" />
                  Cari Intervensi
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-3">
                  <div>
                    <Label
                      for="search"
                      class="text-sm font-medium">
                      Cari berdasarkan nama balita, jenis, atau deskripsi
                    </Label>
                    <Input
                      id="search"
                      v-model="searchQuery"
                      placeholder="Ketik untuk mencari intervensi..."
                      class="mt-1" />
                  </div>
                  <div class="text-xs text-muted-foreground">
                    Ditemukan {{ filteredIntervensi.length }} intervensi
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Pilih Intervensi -->
            <Card :class="errors.intervensi ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Activity class="h-4 w-4" />
                  Pilih Intervensi *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select v-model="selectedIntervensi">
                  <SelectTrigger :class="errors.intervensi ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedIntervensiInfo">
                        {{ selectedIntervensiInfo.nama_balita }} -
                        {{ selectedIntervensiInfo.jenis }}
                      </template>
                      <template v-else> Pilih intervensi untuk ditugaskan </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="intervensi in filteredIntervensi"
                      :key="intervensi.id"
                      :value="intervensi.id">
                      <div class="flex flex-col gap-1 py-1">
                        <div class="flex items-center gap-2">
                          <span class="font-medium">{{ intervensi.nama_balita }}</span>
                          <Badge
                            variant="outline"
                            :class="getJenisColor(intervensi.jenis)"
                            class="text-xs">
                            {{ intervensi.jenis }}
                          </Badge>
                        </div>
                        <div class="text-xs text-muted-foreground">
                          📅 {{ formatDate(intervensi.tanggal) }}
                        </div>
                        <div class="text-xs text-muted-foreground">
                          {{ intervensi.deskripsi }}
                        </div>
                        <div class="text-xs text-blue-600">
                          👥 {{ intervensi.petugas_assigned.length }} petugas ditugaskan
                        </div>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
                <p
                  v-if="errors.intervensi"
                  class="text-sm text-red-500 mt-2">
                  {{ errors.intervensi }}
                </p>
              </CardContent>
            </Card>

            <!-- Detail Intervensi Terpilih -->
            <Card
              v-if="selectedIntervensiInfo"
              class="bg-green-50 border-green-200">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base text-green-900">
                  <FileText class="h-4 w-4" />
                  Detail Intervensi
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-3">
                  <!-- Info Intervensi -->
                  <div class="flex items-start gap-3">
                    <div class="bg-green-100 rounded-full p-2">
                      <Activity class="h-5 w-5 text-green-600" />
                    </div>
                    <div class="flex-1">
                      <div class="font-medium text-green-900">
                        {{ selectedIntervensiInfo.nama_balita }}
                      </div>
                      <div class="flex items-center gap-2 mt-1">
                        <Badge
                          variant="outline"
                          :class="getJenisColor(selectedIntervensiInfo.jenis)"
                          class="text-xs">
                          {{ selectedIntervensiInfo.jenis }}
                        </Badge>
                        <span class="text-xs text-green-700">
                          📅 {{ formatDate(selectedIntervensiInfo.tanggal) }}
                        </span>
                      </div>
                      <div class="text-sm text-green-800 mt-2">
                        {{ selectedIntervensiInfo.deskripsi }}
                      </div>
                    </div>
                  </div>

                  <Separator class="my-3" />

                  <!-- Petugas yang Sudah Ditugaskan -->
                  <div>
                    <div class="font-medium text-green-900 mb-2">
                      Petugas yang Sudah Ditugaskan ({{
                        selectedIntervensiInfo.petugas_assigned.length
                      }})
                    </div>
                    <div
                      v-if="selectedIntervensiInfo.petugas_assigned.length === 0"
                      class="text-sm text-green-600 italic">
                      Belum ada petugas yang ditugaskan
                    </div>
                    <div
                      v-else
                      class="space-y-2">
                      <div
                        v-for="assignment in selectedIntervensiInfo.petugas_assigned"
                        :key="assignment.id"
                        class="flex items-center justify-between p-2 bg-white rounded border border-green-200">
                        <div class="flex items-center gap-2">
                          <Users class="h-4 w-4 text-green-600" />
                          <div>
                            <div class="font-medium text-sm text-green-900">
                              {{ assignment.nama_petugas }}
                            </div>
                            <div class="text-xs text-green-700">{{ assignment.skpd }}</div>
                          </div>
                        </div>
                        <div class="flex items-center gap-2">
                          <CheckCircle class="h-4 w-4 text-green-600" />
                          <span class="text-xs text-green-600">Ditugaskan</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- Status Assignment untuk Petugas Saat Ini -->
                  <div
                    v-if="petugas"
                    class="mt-4 p-3 rounded border">
                    <div class="flex items-center gap-2 mb-2">
                      <AlertCircle
                        class="h-4 w-4"
                        :class="isAlreadyAssigned ? 'text-orange-600' : 'text-blue-600'" />
                      <span
                        class="font-medium text-sm"
                        :class="isAlreadyAssigned ? 'text-orange-900' : 'text-blue-900'">
                        Status Penugasan {{ petugas.nama }}
                      </span>
                    </div>
                    <div
                      v-if="isAlreadyAssigned"
                      class="text-sm text-orange-800">
                      ✅ Petugas sudah ditugaskan pada intervensi ini
                    </div>
                    <div
                      v-else
                      class="text-sm text-blue-800">
                      ➕ Petugas belum ditugaskan pada intervensi ini
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Info Box -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <UserPlus class="h-5 w-5 text-blue-600 mt-0.5 flex-shrink-0" />
                <div class="text-sm text-blue-800">
                  <p class="font-medium mb-2">📋 Panduan Assignment Petugas:</p>
                  <ul class="space-y-1 text-xs">
                    <li>
                      • <strong>Assign:</strong> Tugaskan petugas ke intervensi yang belum memiliki
                      petugas atau butuh tambahan
                    </li>
                    <li>
                      • <strong>Remove:</strong> Hapus penugasan jika petugas tidak lagi bertanggung
                      jawab
                    </li>
                    <li>
                      • <strong>Multiple Assignment:</strong> Satu intervensi dapat ditangani oleh
                      beberapa petugas
                    </li>
                    <li>
                      • <strong>Search:</strong> Gunakan fitur pencarian untuk menemukan intervensi
                      dengan cepat
                    </li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Debug Info (untuk development) -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Petugas: {{ petugas?.nama || "None" }}<br />
              Petugas ID: {{ petugas?.id || "None" }}<br />
              Selected Intervensi: {{ selectedIntervensiInfo?.nama_balita || "None" }}<br />
              Already Assigned: {{ isAlreadyAssigned }}<br />
              Assignment ID: {{ currentAssignmentId || "None" }}<br />
              Total Intervensi: {{ filteredIntervensi.length }}
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
              Tutup
            </Button>

            <!-- Conditional Action Buttons -->
            <div class="flex gap-2">
              <Button
                v-if="isAlreadyAssigned && currentAssignmentId"
                variant="destructive"
                @click="handleRemove"
                class="w-full sm:w-auto">
                <Trash2 class="h-4 w-4 mr-2" />
                Remove Assignment
              </Button>

              <Button
                v-if="!isAlreadyAssigned && selectedIntervensi"
                @click="handleAssign"
                class="w-full sm:w-auto">
                <UserPlus class="h-4 w-4 mr-2" />
                Assign Petugas
              </Button>
            </div>
          </DialogFooter>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>
