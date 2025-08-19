<script setup lang="ts">
import { computed } from "vue"
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { 
  Activity, 
  FileText, 
  Plus, 
  Stethoscope,
} from "lucide-vue-next"
import type { Intervensi, RiwayatPemeriksaan } from "./columns"
import DataTableRiwayatPemeriksaan from "./DataTableRiwayatPemeriksaan.vue"

interface Props {
  show: boolean
  intervensi: Intervensi | null
  riwayatData: RiwayatPemeriksaan[]
}

interface Emits {
  (e: "close"): void
  (e: "add-riwayat", intervensiId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Filter riwayat berdasarkan intervensi yang dipilih
const filteredRiwayatData = computed(() => {
  if (!props.intervensi) return []
  
  return props.riwayatData.filter(
    (riwayat) => riwayat.id_intervensi === props.intervensi?.id
  )
})

// Event handlers
const handleClose = () => {
  emit("close")
}

const handleAddRiwayat = () => {
  if (props.intervensi) {
    emit("add-riwayat", props.intervensi.id)
  }
}
</script>

<template>
  <Sheet :open="show" @update:open="(open) => !open && handleClose()">
    <SheetContent class="w-full sm:max-w-7xl overflow-hidden p-0 flex flex-col">
      <!-- Header -->
      <div class="p-6 border-b bg-background">
        <SheetHeader class="space-y-3">
          <div class="flex items-center justify-between">
            <SheetTitle class="text-xl font-semibold flex items-center gap-3">
              <Activity class="h-6 w-6 text-blue-600" />
              Riwayat Pemeriksaan Intervensi
            </SheetTitle>
          </div>
          
          <SheetDescription class="text-sm text-muted-foreground">
            Detail riwayat pemeriksaan untuk intervensi yang dipilih
          </SheetDescription>
        </SheetHeader>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-hidden flex flex-col">
        <!-- Actions -->
        <div class="px-6 pb-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold flex items-center gap-2">
              <FileText class="h-5 w-5" />
              Daftar Riwayat Pemeriksaan
            </h3>
            <Button @click="handleAddRiwayat" class="gap-2">
              <Plus class="h-4 w-4" />
              Tambah Riwayat
            </Button>
          </div>
        </div>

        <!-- Content Body -->
        <div class="flex-1 overflow-hidden">
          <!-- Data Table -->
          <div v-if="filteredRiwayatData.length > 0" class="h-full flex flex-col">
            <ScrollArea class="flex-1 px-6">
              <DataTableRiwayatPemeriksaan :data="filteredRiwayatData" />
            </ScrollArea>
          </div>

          <!-- Empty State -->
          <div v-else class="flex-1 flex items-center justify-center p-6">
            <div class="text-center space-y-4">
              <div class="p-4 bg-gray-100 rounded-full w-16 h-16 mx-auto flex items-center justify-center">
                <Stethoscope class="h-8 w-8 text-gray-400" />
              </div>
              <div class="space-y-2">
                <h3 class="text-lg font-semibold text-gray-900">Belum Ada Riwayat Pemeriksaan</h3>
                <p class="text-sm text-gray-500 max-w-sm mx-auto">
                  Intervensi ini belum memiliki riwayat pemeriksaan. Tambahkan pemeriksaan pertama untuk mulai melacak perkembangan balita.
                </p>
              </div>
              <Button @click="handleAddRiwayat" class="gap-2">
                <Plus class="h-4 w-4" />
                Tambah Riwayat Pertama
              </Button>
            </div>
          </div>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>