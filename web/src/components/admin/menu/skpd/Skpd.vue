<script setup lang="ts">
import { ref } from "vue"
import type { Skpd } from "./columns"
import { toast } from "vue-sonner"
import { Building, Plus } from "lucide-vue-next"
import Button from "@/components/ui/button/Button.vue"
import Card from "@/components/ui/card/Card.vue"
import CardHeader from "@/components/ui/card/CardHeader.vue"
import CardTitle from "@/components/ui/card/CardTitle.vue"
import CardContent from "@/components/ui/card/CardContent.vue"
import CardDescription from "@/components/ui/card/CardDescription.vue"
import DataTable from "./DataTable.vue"

// Dummy data
const skpdData = ref<Skpd[]>([
  {
    id: "1",
    skpd: "Dinas Kesehatan tes",
    jenis: "puskesmas",
    petugas_count: 2,
    created_date: "2023-01-01",
  },
  {
    id: "2",
    skpd: "Dinas Pendidikan",
    jenis: "kelurahan",
    petugas_count: 5,
    created_date: "2023-01-02",
  },
  {
    id: "3",
    skpd: "Dinas Pekerjaan Umum",
    jenis: "skpd",
    petugas_count: 0,
    created_date: "2023-01-03",
  },
  {
    id: "4",
    skpd: "Dinas Perumahan",
    jenis: "skpd",
    petugas_count: 1,
    created_date: "2023-01-04",
  },
])

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedSkpd = ref<Skpd | null>(null)

// Statistics
const totalSkpd = ref(skpdData.value.length)

// Update Statistics
const updateStatistics = () => {
  totalSkpd.value = skpdData.value.length
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedSkpd.value = null
  showDialog.value = true
}

const handleEdit = (skpd: Skpd) => {
  dialogMode.value = "edit"
  selectedSkpd.value = { ...skpd }
  showDialog.value = true
}

const handleDelete = (skpd: Skpd) => {
  if (confirm(`Apakah Anda yakin ingin menghapus data SKPD ${skpd.skpd}?`)) {
    const index = skpdData.value.findIndex((b) => b.id === skpd.id)
    if (index > -1) {
      skpdData.value.splice(index, 1)
      updateStatistics()
      toast.success("Data SKPD berhasil dihapus")
    }
  }
}

const handleSave = (skpd: Skpd) => {
  if (dialogMode.value === "create") {
    toast.success("Data SKPD berhasil ditambahkan")
  } else if (dialogMode.value === "edit") {
    toast.success("Data SKPD berhasil diperbarui")
  }

  updateStatistics()
  showDialog.value = false
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Building class="h-8 w-8 text-slate-600" />
          Data SKPD
        </h1>
        <p class="text-gray-600">Kelola data SKPD di Kota Cirebon</p>
      </div>
      <Button
        @click="handleCreate"
        class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah SKPD
      </Button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total SKPD</CardTitle>
          <Building class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalSkpd }}</div>
          <p class="text-xs text-muted-foreground">SKPD terdaftar</p>
        </CardContent>
      </Card>
    </div>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <Building class="h-5 w-5" />
          Daftar SKPD
        </CardTitle>
        <CardDescription>
          Kelola data SKPD yang ada di Kota Cirebon.
        </CardDescription>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTable :data="skpdData" />
      </CardContent>
    </Card>
  </div>
</template>
