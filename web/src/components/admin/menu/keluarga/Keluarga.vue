<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Users } from "lucide-vue-next"
import type { Keluarga } from "./columns"

// Dummy data
const keluargaData = ref<Keluarga[]>([
  {
    id: "1",
    nomor_kk: "3209012345678901",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nik_ayah: "3209012345678901",
    nik_ibu: "3209012345678902",
    alamat: "Jl. Kesambi Raya No. 123",
    rt: "001",
    rw: "002",
    id_kelurahan: "1",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    koordinat: [-6.7064, 108.5492],
    created_date: "2024-01-15",
    updated_date: "2024-01-20",
  },
  {
    id: "2",
    nomor_kk: "3209012345678903",
    nama_ayah: "Ahmad Wijaya",
    nama_ibu: "Dewi Sartika",
    nik_ayah: "3209012345678903",
    nik_ibu: "3209012345678904",
    alamat: "Jl. Tuparev No. 45",
    rt: "003",
    rw: "001",
    id_kelurahan: "2",
    kelurahan: "Tuparev",
    kecamatan: "Kedawung",
    koordinat: [-6.7104, 108.5532],
    created_date: "2024-01-16",
    updated_date: "2024-01-21",
  },
  {
    id: "3",
    nomor_kk: "3209012345678905",
    nama_ayah: "Rizki Pratama",
    nama_ibu: "Maya Sari",
    nik_ayah: "3209012345678905",
    nik_ibu: "3209012345678906",
    alamat: "Jl. Perjuangan No. 67",
    rt: "005",
    rw: "003",
    id_kelurahan: "3",
    kelurahan: "Perjuangan",
    kecamatan: "Kejaksan",
    koordinat: [-6.7144, 108.5572],
    created_date: "2024-01-17",
  },
  {
    id: "4",
    nomor_kk: "3209012345678907",
    nama_ayah: "Dedi Kurniawan",
    nama_ibu: "Rina Melati",
    nik_ayah: "3209012345678907",
    nik_ibu: "3209012345678908",
    alamat: "Jl. Brigjen Dharsono No. 89",
    rt: "002",
    rw: "004",
    id_kelurahan: "4",
    kelurahan: "Argasunya",
    kecamatan: "Harjamukti",
    koordinat: [-6.7184, 108.5612],
    created_date: "2024-01-18",
  },
  {
    id: "5",
    nomor_kk: "3209012345678909",
    nama_ayah: "Eko Prasetyo",
    nama_ibu: "Lilis Suryani",
    nik_ayah: "3209012345678909",
    nik_ibu: "3209012345678910",
    alamat: "Jl. Kartini No. 12",
    rt: "004",
    rw: "002",
    id_kelurahan: "5",
    kelurahan: "Lemahwungkuk",
    kecamatan: "Lemahwungkuk",
    koordinat: [-6.7224, 108.5652],
    created_date: "2024-01-19",
  },
  {
    id: "6",
    nomor_kk: "3209012345678911",
    nama_ayah: "Fajar Setiawan",
    nama_ibu: "Sari Wulandari",
    nik_ayah: "3209012345678911",
    nik_ibu: "3209012345678912",
    alamat: "Jl. Cipto Mangunkusumo No. 34",
    rt: "006",
    rw: "005",
    id_kelurahan: "6",
    kelurahan: "Kesenden",
    kecamatan: "Kesenden",
    koordinat: [-6.7250, 108.5700],
    created_date: "2024-01-20",
  },
  {
    id: "7",
    nomor_kk: "3209012345678913",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nik_ayah: "3209012345678913",
    nik_ibu: "3209012345678914",
    alamat: "Jl. Kesambi Raya No. 123",
    rt: "001",
    rw: "002",
    id_kelurahan: "1",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    koordinat: [-6.7074, 108.5502],
    created_date: "2024-01-15",
    updated_date: "2024-01-20",
  },
  {
    id: "8",
    nomor_kk: "3209012345678915",
    nama_ayah: "Fajar Setiawan",
    nama_ibu: "Sari Wulandari",
    nik_ayah: "3209012345678915",
    nik_ibu: "3209012345678916",
    alamat: "Jl. Cipto Mangunkusumo No. 34",
    rt: "006",
    rw: "005",
    id_kelurahan: "6",
    kelurahan: "Kesenden",
    kecamatan: "Kesenden",
    koordinat: [-6.7260, 108.5710],
    created_date: "2024-01-20",
  },
  {
    id: "9",
    nomor_kk: "3209012345678917",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nik_ayah: "3209012345678917",
    nik_ibu: "3209012345678918",
    alamat: "Jl. Kesambi Raya No. 123",
    rt: "001",
    rw: "002",
    id_kelurahan: "1",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    koordinat: [-6.7084, 108.5512],
    created_date: "2024-01-15",
    updated_date: "2024-01-20",
  },
  {
    id: "10",
    nomor_kk: "3209012345678919",
    nama_ayah: "Fajar Setiawan",
    nama_ibu: "Sari Wulandari",
    nik_ayah: "3209012345678919",
    nik_ibu: "3209012345678920",
    alamat: "Jl. Cipto Mangunkusumo No. 34",
    rt: "006",
    rw: "005",
    id_kelurahan: "6",
    kelurahan: "Kesenden",
    kecamatan: "Kesenden",
    koordinat: [-6.7270, 108.5720],
    created_date: "2024-01-20",
  },
  {
    id: "11",
    nomor_kk: "3209012345678921",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nik_ayah: "3209012345678921",
    nik_ibu: "3209012345678922",
    alamat: "Jl. Kesambi Raya No. 123",
    rt: "001",
    rw: "002",
    id_kelurahan: "1",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    koordinat: [-6.7094, 108.5522],
    created_date: "2024-01-15",
    updated_date: "2024-01-20",
  },
  {
    id: "12",
    nomor_kk: "3209012345678923",
    nama_ayah: "Fajar Setiawan",
    nama_ibu: "Sari Wulandari",
    nik_ayah: "3209012345678923",
    nik_ibu: "3209012345678924",
    alamat: "Jl. Cipto Mangunkusumo No. 34",
    rt: "006",
    rw: "005",
    id_kelurahan: "6",
    kelurahan: "Kesenden",
    kecamatan: "Kesenden",
    koordinat: [-6.7280, 108.5730],
    created_date: "2024-01-20",
  },
])

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedKeluarga = ref<Keluarga | null>(null)

// Statistics
const totalKeluarga = ref(keluargaData.value.length)
const totalWithCoordinates = ref(
  keluargaData.value.filter((k) => k.koordinat[0] !== 0 && k.koordinat[1] !== 0).length
)

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

const handleDelete = (keluarga: Keluarga) => {
  if (confirm(`Apakah Anda yakin ingin menghapus data keluarga ${keluarga.nama_ayah}?`)) {
    const index = keluargaData.value.findIndex((k) => k.id === keluarga.id)
    if (index > -1) {
      keluargaData.value.splice(index, 1)
      totalKeluarga.value = keluargaData.value.length
      totalWithCoordinates.value = keluargaData.value.filter(
        (k) => k.koordinat[0] !== 0 && k.koordinat[1] !== 0
      ).length
      toast.success("Data keluarga berhasil dihapus")
    }
  }
}

const handleSave = (keluarga: Keluarga) => {
  if (dialogMode.value === "create") {
    // Generate new ID
    const newId = (Math.max(...keluargaData.value.map((k) => parseInt(k.id))) + 1).toString()
    const newKeluarga = {
      ...keluarga,
      id: newId,
      created_date: new Date().toISOString().split("T")[0],
    }
    keluargaData.value.unshift(newKeluarga)
    toast.success("Data keluarga berhasil ditambahkan")
  } else {
    // Update existing
    const index = keluargaData.value.findIndex((k) => k.id === keluarga.id)
    if (index > -1) {
      keluargaData.value[index] = {
        ...keluarga,
        updated_date: new Date().toISOString().split("T")[0],
      }
      toast.success("Data keluarga berhasil diperbarui")
    }
  }

  totalKeluarga.value = keluargaData.value.length
  totalWithCoordinates.value = keluargaData.value.filter(
    (k) => k.koordinat[0] !== 0 && k.koordinat[1] !== 0
  ).length
  showDialog.value = false
}

const handleCustomEvents = (event: Event) => {
  const customEvent = event as CustomEvent

  if (event.type === "edit-keluarga") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-keluarga") {
    handleDelete(customEvent.detail)
  }
}

onMounted(() => {
  document.addEventListener("edit-keluarga", handleCustomEvents)
  document.addEventListener("delete-keluarga", handleCustomEvents)
})

onUnmounted(() => {
  document.removeEventListener("edit-keluarga", handleCustomEvents)
  document.removeEventListener("delete-keluarga", handleCustomEvents)
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Data Keluarga</h1>
        <p class="text-gray-600">Kelola data keluarga balita di Kota Cirebon</p>
      </div>
      <Button
        @click="handleCreate"
        class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah Keluarga
      </Button>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Keluarga</CardTitle>
          <Users class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalKeluarga }}</div>
          <p class="text-xs text-muted-foreground">Keluarga terdaftar</p>
        </CardContent>
      </Card>

      <!-- <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Dengan Koordinat</CardTitle>
          <MapPin class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalWithCoordinates }}</div>
          <p class="text-xs text-muted-foreground">Memiliki koordinat lokasi</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Persentase Koordinat</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">
            {{ totalKeluarga > 0 ? Math.round((totalWithCoordinates / totalKeluarga) * 100) : 0 }}%
          </div>
          <p class="text-xs text-muted-foreground">Kelengkapan data lokasi</p>
        </CardContent>
      </Card> -->
    </div>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle>Daftar Keluarga</CardTitle>
        <CardDescription>
          Data keluarga yang terdaftar dalam sistem pemantauan stunting
        </CardDescription>
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
      @close="showDialog = false"
      @save="handleSave" />
  </div>
</template>
