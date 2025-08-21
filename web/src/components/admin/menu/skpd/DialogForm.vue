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
import { Building, Save, X, Shield, MapPin } from "lucide-vue-next"
import type { Skpd } from "./columns"

interface Props {
  show: boolean
  mode: "create" | "edit"
  skpd: Skpd | null
}

interface Emits {
  (e: "close"): void
  (e: "save", skpd: Skpd): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data sesuai dengan endpoint
const formData = ref<Partial<Skpd>>({
  skpd: "",
  jenis: "",
})

// Validation errors
const errors = ref<Record<string, string>>({})

// Jenis SKPD options dengan detail informasi
const jenisSkpdOptions = [
  {
    value: "puskesmas",
    label: "Puskesmas",
    description: "Pusat Kesehatan Masyarakat - Unit pelaksana teknis kesehatan",
    icon: "🏥",
    color: "blue"
  },
  {
    value: "kelurahan",
    label: "Kelurahan",
    description: "Kelurahan - Unit administrasi pemerintahan setingkat desa",
    icon: "🏛️",
    color: "green"
  },
  {
    value: "skpd",
    label: "SKPD",
    description: "Satuan Kerja Perangkat Daerah - Unit kerja pemerintah daerah",
    icon: "🏢",
    color: "purple"
  }
]

// Get selected jenis info
const selectedJenisInfo = computed(() => {
  if (!formData.value.jenis) return null
  return jenisSkpdOptions.find(j => j.value === formData.value.jenis) || null
})

// Helper function untuk get color class berdasarkan jenis
const getJenisColorClass = (jenis: string): string => {
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

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Nama SKPD validation
  if (!formData.value.skpd || formData.value.skpd.trim().length < 3) {
    errors.value.skpd = "Nama SKPD minimal 3 karakter"
  } else if (formData.value.skpd.length > 100) {
    errors.value.skpd = "Nama SKPD maksimal 100 karakter"
  } else {
    // Validasi format nama SKPD
    const validPrefixes = ["Dinas", "Kantor", "Badan", "Sekretariat", "Puskesmas", "Kelurahan"]
    const hasValidPrefix = validPrefixes.some(prefix => 
      formData.value.skpd!.toLowerCase().startsWith(prefix.toLowerCase())
    )
    
    if (!hasValidPrefix && formData.value.jenis !== "skpd") {
      errors.value.skpd = "Nama SKPD harus dimulai dengan: " + validPrefixes.join(", ")
    }
  }

  // Jenis validation
  if (!formData.value.jenis) {
    errors.value.jenis = "Jenis SKPD harus dipilih"
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    skpd: "",
    jenis: "",
  }
  errors.value = {}
}

// Load form data for edit mode
const loadFormData = (skpd: Skpd) => {
  formData.value = {
    id: skpd.id,
    skpd: skpd.skpd,
    jenis: skpd.jenis,
    // Extended fields for response
    petugas_count: skpd.petugas_count,
    created_date: skpd.created_date,
    updated_date: skpd.updated_date,
  }
  errors.value = {}
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

  // Set default petugas_count untuk SKPD baru
  if (props.mode === "create") {
    formData.value.petugas_count = 0
  }

  // Set timestamps
  const now = new Date().toISOString()
  if (props.mode === "create") {
    formData.value.created_date = now
  }
  formData.value.updated_date = now

  emit("save", formData.value as Skpd)
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
      if (props.skpd && props.mode === "edit") {
        loadFormData(props.skpd)
      } else {
        resetForm()
      }
    }
  }
)

// Watch for skpd prop changes
watch(
  () => props.skpd,
  (newSkpd) => {
    if (newSkpd && props.mode === "edit" && props.show) {
      loadFormData(newSkpd)
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.skpd,
  (newVal) => {
    if (newVal && newVal.trim().length >= 3) {
      delete errors.value.skpd
    }
  }
)

watch(
  () => formData.value.jenis,
  (newVal) => {
    if (newVal) {
      delete errors.value.jenis
    }
  }
)

// Helper untuk mendapatkan contoh nama berdasarkan jenis
const getContohNama = (jenis: string): string[] => {
  switch (jenis?.toLowerCase()) {
    case "puskesmas":
      return [
        "Puskesmas Kejaksan",
        "Puskesmas Pekalangan",
        "Puskesmas Harjamukti"
      ]
    case "kelurahan":
      return [
        "Kelurahan Kejaksan",
        "Kelurahan Pekalangan", 
        "Kelurahan Harjamukti"
      ]
    case "skpd":
      return [
        "Dinas Kesehatan Kota Cirebon",
        "Badan Perencanaan Pembangunan Daerah",
        "Kantor Kependudukan dan Pencatatan Sipil"
      ]
    default:
      return []
  }
}
</script>

<template>
  <Dialog
    :open="show"
    @update:open="handleClose">
    <DialogContent class="w-[95vw] max-w-3xl max-h-[95vh] overflow-hidden p-0">
      <div class="flex flex-col h-full max-h-[95vh]">
        <!-- Header -->
        <div class="p-4 sm:p-6 border-b bg-background">
          <DialogHeader class="space-y-2">
            <DialogTitle class="text-lg sm:text-xl flex items-center gap-2">
              <Building class="h-5 w-5 text-blue-600" />
              {{ mode === "create" ? "Tambah SKPD Baru" : "Edit Data SKPD" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data SKPD baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi SKPD. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">
            <!-- Jenis SKPD -->
            <Card :class="errors.jenis ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Shield class="h-4 w-4" />
                  Jenis SKPD *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Select v-model="formData.jenis" class="">
                  <SelectTrigger class="w-full" :class="errors.jenis ? 'border-red-500' : ''">
                    <SelectValue>
                      <template v-if="selectedJenisInfo">
                        {{ selectedJenisInfo.label }}
                      </template>
                      <template v-else>
                        Pilih jenis SKPD
                      </template>
                    </SelectValue>
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem
                      v-for="jenis in jenisSkpdOptions"
                      :key="jenis.value"
                      :value="jenis.value">
                      <div class="flex items-center gap-3">
                        <span class="text-lg">{{ jenis.icon }}</span>
                        <div class="flex flex-col">
                          <span class="font-medium">{{ jenis.label }}</span>
                          <span class="text-xs text-muted-foreground">{{ jenis.description }}</span>
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
                  :class="getJenisColorClass(selectedJenisInfo.value)">
                  <div class="flex items-center gap-2 mb-2">
                    <span class="text-lg">{{ selectedJenisInfo.icon }}</span>
                    <div>
                      <div class="font-medium text-sm">{{ selectedJenisInfo.label }}</div>
                      <div class="text-xs opacity-90">{{ selectedJenisInfo.description }}</div>
                    </div>
                  </div>
                  <div class="text-xs opacity-80 mt-2">
                    <strong>Contoh nama:</strong>
                    <ul class="list-disc list-inside mt-1 space-y-0.5">
                      <li v-for="contoh in getContohNama(selectedJenisInfo.value)" :key="contoh">
                        {{ contoh }}
                      </li>
                    </ul>
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Nama SKPD -->
            <Card :class="errors.skpd ? 'border-red-500' : ''">
              <CardHeader class="pb-3">
                <CardTitle class="flex items-center gap-2 text-base">
                  <Building class="h-4 w-4" />
                  Nama SKPD *
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div class="space-y-2">
                  <Label
                    for="skpd"
                    class="text-sm font-medium">
                    Nama Lengkap SKPD
                  </Label>
                  <Input
                    id="skpd"
                    v-model="formData.skpd"
                    placeholder="Masukkan nama lengkap SKPD..."
                    :class="errors.skpd ? 'border-red-500' : ''" />
                  <p
                    v-if="errors.skpd"
                    class="text-sm text-red-500">
                    {{ errors.skpd }}
                  </p>
                  <div class="flex justify-between items-center text-xs text-muted-foreground">
                    <span>Minimum 3 karakter</span>
                    <span>{{ formData.skpd?.length || 0 }}/100 karakter</span>
                  </div>
                </div>

                <!-- Preview nama berdasarkan jenis -->
                <div
                  v-if="formData.jenis && !formData.skpd"
                  class="mt-3 p-3 bg-gray-50 rounded-md border border-gray-200">
                  <div class="text-sm text-gray-700 mb-2">
                    <strong>💡 Saran nama untuk {{ selectedJenisInfo?.label }}:</strong>
                  </div>
                  <div class="space-y-1">
                    <button
                      v-for="contoh in getContohNama(formData.jenis)"
                      :key="contoh"
                      @click="formData.skpd = contoh"
                      class="block w-full text-left text-xs p-2 hover:bg-gray-100 rounded border border-dashed border-gray-300 transition-colors">
                      {{ contoh }}
                    </button>
                  </div>
                  <div class="text-xs text-gray-500 mt-2">
                    Klik salah satu untuk menggunakan sebagai template
                  </div>
                </div>
              </CardContent>
            </Card>

            <!-- Info Box -->
            <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <MapPin class="h-5 w-5 text-blue-600 mt-0.5 flex-shrink-0" />
                <div class="text-sm text-blue-800">
                  <p class="font-medium mb-2">📋 Panduan Pengisian SKPD:</p>
                  <ul class="space-y-1 text-xs">
                    <li>• <strong>Jenis:</strong> Pilih kategori yang sesuai dengan unit kerja</li>
                    <li>• <strong>Nama:</strong> Gunakan nama resmi dan lengkap sesuai SK pembentukan</li>
                    <li>• <strong>Puskesmas:</strong> Harus mengandung kata "Puskesmas" dalam nama</li>
                    <li>• <strong>Kelurahan:</strong> Harus mengandung kata "Kelurahan" dalam nama</li>
                    <li>• <strong>SKPD:</strong> Dapat berupa Dinas, Badan, Kantor, atau Sekretariat</li>
                  </ul>
                </div>
              </div>
            </div>

            <!-- Live Preview -->
            <div
              v-if="formData.skpd && formData.jenis"
              class="bg-green-50 border border-green-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <div class="text-2xl">✅</div>
                <div>
                  <div class="text-sm font-medium text-green-900 mb-1">Preview SKPD:</div>
                  <div class="flex items-center gap-2">
                    <span class="text-lg">{{ selectedJenisInfo?.icon }}</span>
                    <div>
                      <div class="font-medium text-green-800">{{ formData.skpd }}</div>
                      <div class="text-xs text-green-600">{{ selectedJenisInfo?.label }} • Kota Cirebon</div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Debug Info (untuk development) -->
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 text-xs">
              <strong>🐛 Debug Info:</strong><br />
              Mode: {{ mode }}<br />
              Has Data: {{ !!props.skpd }}<br />
              Form ID: {{ formData.id || "New" }}<br />
              Nama SKPD: {{ formData.skpd || "Empty" }}<br />
              Jenis: {{ formData.jenis || "Not selected" }}<br />
              Petugas Count: {{ formData.petugas_count || 0 }}
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