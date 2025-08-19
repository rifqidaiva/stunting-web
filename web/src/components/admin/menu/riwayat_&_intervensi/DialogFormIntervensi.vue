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
import type { Intervensi } from "./columns"
import { ref, watch } from "vue"
import type { DateValue } from "reka-ui"
import { toast } from "vue-sonner"
import { format } from "date-fns"
import { CalendarDate } from "@internationalized/date"
import { FileText, Save, X } from "lucide-vue-next"

interface Props {
  show: boolean
  mode: "create" | "edit"
  intervensi: Intervensi | null
}

interface Emits {
  (e: "close"): void
  (e: "save", intervensi: Intervensi): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formData = ref<Partial<Intervensi>>({
  id: "",
  id_balita: "",
  deskripsi: "",
  hasil: "",
  jenis: "",
  tanggal: ""
})

const selectedDate = ref<DateValue>()
const isCalendarOpen = ref(false)

const balitaOptions = [
  { id: "1", name: "Balita 1" },
  { id: "2", name: "Balita 2" },
  { id: "3", name: "Balita 3" },
]

const errors = ref<Record<string, string>>({})

const validateForm = () => {
  errors.value = {}
  if (!formData.value.id_balita) {
    errors.value.id_balita = "Balita is required."
  }
  if (!formData.value.deskripsi) {
    errors.value.deskripsi = "Deskripsi is required."
  }
  if (!formData.value.hasil) {
    errors.value.hasil = "Hasil is required."
  }
  if (!formData.value.jenis) {
    errors.value.jenis = "Jenis is required."
  }
  if (!formData.value.tanggal) {
    errors.value.tanggal = "Tanggal is required."
  }

  return Object.keys(errors.value).length === 0
}

const resetForm = () => {
  formData.value = {
    id: "",
    id_balita: "",
    deskripsi: "",
    hasil: "",
    jenis: "",
    tanggal: ""
  }
  selectedDate.value = undefined
  isCalendarOpen.value = false
  errors.value = {}
}

const loadFormData = (intervensi: Intervensi) => {
  formData.value = {
    ...intervensi
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

  emit("save", formData.value as Intervensi)
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
      if (props.intervensi && props.mode === "edit") {
        loadFormData(props.intervensi)
      } else {
        resetForm()
      }
    }
  }
)

// Watch for intervensi prop changes
watch(
  () => props.intervensi,
  (newIntervensi) => {
    if (newIntervensi && props.mode === "edit" && props.show) {
      loadFormData(newIntervensi)
      // Initialize date for calendar
      if (newIntervensi.tanggal) {
        const jsDate = new Date(newIntervensi.tanggal)
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
              {{ mode === "create" ? "Tambah Intervensi Baru" : "Edit Data Intervensi" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data intervensi baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi intervensi. Semua field dengan tanda (*) wajib diisi."
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
