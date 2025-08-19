<script setup lang="ts">
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import Button from "@/components/ui/button/Button.vue"
import type { RiwayatPemeriksaan } from "./columns"
import { ref, watch } from "vue"
import type { DateValue } from "reka-ui"
import { toast } from "vue-sonner"
import { format } from "date-fns"
import { CalendarDate } from "@internationalized/date"
import { FileText, Save, X } from "lucide-vue-next"

interface Props {
  show: boolean
  mode: "create" | "edit"
  riwayatPemeriksaan: RiwayatPemeriksaan | null
}

interface Emits {
  (e: "close"): void
  (e: "save", riwayatPemeriksaan: RiwayatPemeriksaan): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formData = ref<Partial<RiwayatPemeriksaan>>({
  id: "",
  id_balita: "",
  id_intervensi: "",
  id_laporan_masyarakat: "",
  keterangan: "",
  status_gizi: "",
  tanggal: "",
  tinggi_badan: "",
  berat_badan: "",
})

const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

const errors = ref<Record<string, string>>({})

const validateForm = () => {
  errors.value = {}
  if (!formData.value.id_balita) {
    errors.value.id_balita = "Balita is required."
  }
  if (!formData.value.keterangan) {
    errors.value.keterangan = "Keterangan is required."
  }
  if (!formData.value.status_gizi) {
    errors.value.status_gizi = "Status Gizi is required."
  }
  if (!formData.value.tanggal) {
    errors.value.tanggal = "Tanggal is required."
  }
  if (!formData.value.tinggi_badan) {
    errors.value.tinggi_badan = "Tinggi Badan is required."
  }
  if (!formData.value.berat_badan) {
    errors.value.berat_badan = "Berat Badan is required."
  }

  return Object.keys(errors.value).length === 0
}

const resetForm = () => {
  formData.value = {
    id: "",
    id_balita: "",
    id_intervensi: "",
    id_laporan_masyarakat: "",
    keterangan: "",
    status_gizi: "",
    tanggal: "",
    tinggi_badan: "",
    berat_badan: "",
  }
  selectedDate.value = undefined
  isCalendarOpen.value = false
  errors.value = {}
}

const loadFormData = (riwayatPemeriksaan: RiwayatPemeriksaan) => {
  formData.value = {
    ...riwayatPemeriksaan,
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

  emit("save", formData.value as RiwayatPemeriksaan)
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
    formData.value.tanggal = format(jsDate, "yyyy-MM-dd")
    isCalendarOpen.value = false
  }
}

// Watch for dialog visibility and mode changes
watch(
  () => props.show,
  (newVal) => {
    if (newVal) {
      if (props.riwayatPemeriksaan && props.mode === "edit") {
        loadFormData(props.riwayatPemeriksaan)
      } else {
        resetForm()
      }
    }
  }
)

// Watch for riwayatPemeriksaan prop changes
watch(
  () => props.riwayatPemeriksaan,
  (newRiwayatPemeriksaan) => {
    if (newRiwayatPemeriksaan && props.mode === "edit" && props.show) {
      loadFormData(newRiwayatPemeriksaan)
      // Initialize date for calendar
      if (newRiwayatPemeriksaan.tanggal) {
        const jsDate = new Date(newRiwayatPemeriksaan.tanggal)
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
  }
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
              <FileText class="h-5 w-5" />
              {{ mode === "create" ? "Tambah Riwayat Baru" : "Edit Data Riwayat" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data riwayat baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi riwayat. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="space-y-6">tes</div>
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
