<script setup lang="ts">
import { ref, watch } from "vue"
import type { LaporanMasyarakat } from "./columns"
import type { DateValue } from "reka-ui"
import { id } from "date-fns/locale"
import { toast } from "vue-sonner"
import { format } from "date-fns"
import { CalendarDate } from "@internationalized/date"
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
import { CalendarIcon, Save, X, Baby, Users } from "lucide-vue-next"

interface Props {
  show: boolean
  mode: "create" | "edit"
  laporan: LaporanMasyarakat | null
}

interface Emits {
  (e: "close"): void
  (e: "save", laporan: LaporanMasyarakat): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formData = ref<Partial<LaporanMasyarakat>>({
  id: "",
  id_balita: "",
  id_masyarakat: "",
  id_status_laporan: "",
  hubungan_dengan_balita: "",
  nomor_hp_keluarga_balita: "",
  nomor_hp_pelapor: "",
  tanggal_laporan: "",
})

// Date picker state
const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

const balitaOptions = [
  { id: "1", name: "Balita 1" },
  { id: "2", name: "Balita 2" },
  { id: "3", name: "Balita 3" },
]

const masyarakatOptions = [
  { id: "1", name: "Masyarakat 1" },
  { id: "2", name: "Masyarakat 2" },
  { id: "3", name: "Masyarakat 3" },
]

const statusLaporanOptions = [
  { id: "1", name: "Status 1" },
  { id: "2", name: "Status 2" },
  { id: "3", name: "Status 3" },
]

// Validation errors
const errors = ref<Record<string, string>>({})

const validateForm = () => {
  errors.value = {}
  if (!formData.value.id_balita) {
    errors.value.id_balita = "Balita is required."
  }
  if (!formData.value.id_masyarakat) {
    errors.value.id_masyarakat = "Masyarakat is required."
  }
  if (!formData.value.id_status_laporan) {
    errors.value.id_status_laporan = "Status laporan is required."
  }
  if (!formData.value.hubungan_dengan_balita) {
    errors.value.hubungan_dengan_balita = "Hubungan dengan balita is required."
  }

  if (!formData.value.nomor_hp_keluarga_balita) {
    errors.value.nomor_hp_keluarga_balita = "Nomor HP keluarga balita is required."
  } else if (!/^\d{10,15}$/.test(formData.value.nomor_hp_keluarga_balita)) {
    errors.value.nomor_hp_keluarga_balita = "Nomor HP keluarga balita harus antara 10 dan 15 digit."
  }

  if (!formData.value.nomor_hp_pelapor) {
    errors.value.nomor_hp_pelapor = "Nomor HP pelapor is required."
  } else if (!/^\d{10,15}$/.test(formData.value.nomor_hp_pelapor)) {
    errors.value.nomor_hp_pelapor = "Nomor HP pelapor harus antara 10 dan 15 digit."
  }

  if (!formData.value.tanggal_laporan) {
    errors.value.tanggal_laporan = "Tanggal laporan is required."
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    id: "",
    id_balita: "",
    id_masyarakat: "",
    id_status_laporan: "",
    hubungan_dengan_balita: "",
    nomor_hp_keluarga_balita: "",
    nomor_hp_pelapor: "",
    tanggal_laporan: "",
  }
  selectedDate.value = undefined
  errors.value = {}
}

// Load data from edit mode
const loadFormData = (laporan: LaporanMasyarakat) => {
  formData.value = {
    ...laporan,
  }
}

// Handle save
const handleSave = () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi.")

    // Scroll to first error
    const firstErrorElement = document.querySelector(".border-red-500")
    if (firstErrorElement) {
      firstErrorElement.scrollIntoView({ behavior: "smooth", block: "center" })
    }

    return
  }

  emit("save", formData.value as LaporanMasyarakat)
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
    formData.value.tanggal_laporan = format(jsDate, "yyyy-MM-dd")
    isCalendarOpen.value = false
  }
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      if (props.laporan && props.mode === "edit") {
        loadFormData(props.laporan)
      } else {
        resetForm()
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
      // Initialize date for calendar
      if (newLaporan.tanggal_laporan) {
        const jsDate = new Date(newLaporan.tanggal_laporan)
        selectedDate.value = new CalendarDate(
          jsDate.getFullYear(),
          jsDate.getMonth() + 1,
          jsDate.getDate()
        )
      }
    } else {
      // Reset calendar when not editing
      selectedDate.value = undefined
    }
  },
  { deep: true }
)
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
              {{ mode === "create" ? "Tambah Laporan Baru" : "Edit Data Laporan" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data laporan baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi laporan. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">

            tes
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
