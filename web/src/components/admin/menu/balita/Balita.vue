<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Plus, Baby, Users, Calendar } from "lucide-vue-next"
import type { Balita } from "./columns"

// Dummy data
const balitaData = ref<Balita[]>([
  {
    id: "1",
    id_keluarga: "1",
    nomor_kk: "3209012345678901",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nama: "Andi Pratama",
    tanggal_lahir: "2022-03-15",
    jenis_kelamin: "L",
    berat_lahir: "3200",
    tinggi_lahir: "48",
    umur: "2 tahun 5 bulan",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    created_date: "2024-01-15",
    updated_date: "2024-01-20",
  },
  {
    id: "2",
    id_keluarga: "1",
    nomor_kk: "3209012345678901",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Rahayu",
    nama: "Sari Cantika",
    tanggal_lahir: "2021-06-20",
    jenis_kelamin: "P",
    berat_lahir: "2800",
    tinggi_lahir: "46",
    umur: "3 tahun 2 bulan",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    created_date: "2024-01-15",
  },
  {
    id: "3",
    id_keluarga: "2",
    nomor_kk: "3209012345678903",
    nama_ayah: "Ahmad Wijaya",
    nama_ibu: "Dewi Sartika",
    nama: "Rafi Ahmad",
    tanggal_lahir: "2023-01-10",
    jenis_kelamin: "L",
    berat_lahir: "3500",
    tinggi_lahir: "50",
    umur: "1 tahun 7 bulan",
    kelurahan: "Tuparev",
    kecamatan: "Kedawung",
    created_date: "2024-01-16",
  },
  {
    id: "4",
    id_keluarga: "3",
    nomor_kk: "3209012345678905",
    nama_ayah: "Rizki Pratama",
    nama_ibu: "Maya Sari",
    nama: "Dinda Permata",
    tanggal_lahir: "2022-11-05",
    jenis_kelamin: "P",
    berat_lahir: "3100",
    tinggi_lahir: "47",
    umur: "2 tahun 1 bulan",
    kelurahan: "Perjuangan",
    kecamatan: "Kejaksan",
    created_date: "2024-01-17",
  },
  {
    id: "5",
    id_keluarga: "4",
    nomor_kk: "3209012345678907",
    nama_ayah: "Dedi Kurniawan",
    nama_ibu: "Rina Melati",
    nama: "Bayu Saputra",
    tanggal_lahir: "2020-08-25",
    jenis_kelamin: "L",
    berat_lahir: "3800",
    tinggi_lahir: "52",
    umur: "4 tahun 0 bulan",
    kelurahan: "Argasunya",
    kecamatan: "Harjamukti",
    created_date: "2024-01-18",
  },
  {
    id: "6",
    id_keluarga: "5",
    nomor_kk: "3209012345678909",
    nama_ayah: "Eko Prasetyo",
    nama_ibu: "Lilis Suryani",
    nama: "Citra Dewi",
    tanggal_lahir: "2023-04-12",
    jenis_kelamin: "P",
    berat_lahir: "2900",
    tinggi_lahir: "45",
    umur: "1 tahun 4 bulan",
    kelurahan: "Lemahwungkuk",
    kecamatan: "Lemahwungkuk",
    created_date: "2024-01-19",
  },
  {
    id: "7",
    id_keluarga: "6",
    nomor_kk: "3209012345678911",
    nama_ayah: "Fajar Setiawan",
    nama_ibu: "Sari Wulandari",
    nama: "Kevin Ardiansyah",
    tanggal_lahir: "2022-12-08",
    jenis_kelamin: "L",
    berat_lahir: "3400",
    tinggi_lahir: "49",
    umur: "1 tahun 8 bulan",
    kelurahan: "Kesenden",
    kecamatan: "Kesenden",
    created_date: "2024-01-20",
  },
  {
    id: "8",
    id_keluarga: "2",
    nomor_kk: "3209012345678903",
    nama_ayah: "Ahmad Wijaya",
    nama_ibu: "Dewi Sartika",
    nama: "Luna Safira",
    tanggal_lahir: "2021-09-14",
    jenis_kelamin: "P",
    berat_lahir: "3000",
    tinggi_lahir: "47",
    umur: "2 tahun 11 bulan",
    kelurahan: "Tuparev",
    kecamatan: "Kedawung",
    created_date: "2024-01-21",
  },
])

// Dialog states
const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedBalita = ref<Balita | null>(null)

// Statistics
const totalBalita = ref(balitaData.value.length)
const totalLakiLaki = ref(balitaData.value.filter((b) => b.jenis_kelamin === "L").length)
const totalPerempuan = ref(balitaData.value.filter((b) => b.jenis_kelamin === "P").length)

// Helper function untuk update statistics
const updateStatistics = () => {
  totalBalita.value = balitaData.value.length
  totalLakiLaki.value = balitaData.value.filter((b) => b.jenis_kelamin === "L").length
  totalPerempuan.value = balitaData.value.filter((b) => b.jenis_kelamin === "P").length
}

// Event handlers
const handleCreate = () => {
  dialogMode.value = "create"
  selectedBalita.value = null
  showDialog.value = true
}

const handleEdit = (balita: Balita) => {
  dialogMode.value = "edit"
  selectedBalita.value = { ...balita }
  showDialog.value = true
}

const handleDelete = (balita: Balita) => {
  if (confirm(`Apakah Anda yakin ingin menghapus data balita ${balita.nama}?`)) {
    const index = balitaData.value.findIndex((b) => b.id === balita.id)
    if (index > -1) {
      balitaData.value.splice(index, 1)
      updateStatistics()
      toast.success("Data balita berhasil dihapus")
    }
  }
}

const handleSave = (balita: Balita) => {
  if (dialogMode.value === "create") {
    // Generate new ID
    const newId = (Math.max(...balitaData.value.map((b) => parseInt(b.id))) + 1).toString()
    const newBalita = {
      ...balita,
      id: newId,
      created_date: new Date().toISOString().split("T")[0],
    }
    balitaData.value.unshift(newBalita)
    toast.success("Data balita berhasil ditambahkan")
  } else {
    // Update existing
    const index = balitaData.value.findIndex((b) => b.id === balita.id)
    if (index > -1) {
      balitaData.value[index] = {
        ...balita,
        updated_date: new Date().toISOString().split("T")[0],
      }
      toast.success("Data balita berhasil diperbarui")
    }
  }

  updateStatistics()
  showDialog.value = false
}

const handleCustomEvents = (event: Event) => {
  const customEvent = event as CustomEvent

  if (event.type === "edit-balita") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-balita") {
    handleDelete(customEvent.detail)
  }
}

onMounted(() => {
  document.addEventListener("edit-balita", handleCustomEvents)
  document.addEventListener("delete-balita", handleCustomEvents)
})

onUnmounted(() => {
  document.removeEventListener("edit-balita", handleCustomEvents)
  document.removeEventListener("delete-balita", handleCustomEvents)
})
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Baby class="h-8 w-8 text-blue-600" />
          Data Balita
        </h1>
        <p class="text-gray-600">Kelola data balita (anak usia di bawah 5 tahun) di Kota Cirebon</p>
      </div>
      <Button @click="handleCreate" class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah Balita
      </Button>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Balita</CardTitle>
          <Baby class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalBalita }}</div>
          <p class="text-xs text-muted-foreground">Balita terdaftar</p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Laki-laki</CardTitle>
          <Users class="h-4 w-4 text-blue-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-blue-600">{{ totalLakiLaki }}</div>
          <p class="text-xs text-muted-foreground">
            {{ totalBalita > 0 ? Math.round((totalLakiLaki / totalBalita) * 100) : 0 }}% dari total
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Perempuan</CardTitle>
          <Users class="h-4 w-4 text-pink-500" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold text-pink-600">{{ totalPerempuan }}</div>
          <p class="text-xs text-muted-foreground">
            {{ totalBalita > 0 ? Math.round((totalPerempuan / totalBalita) * 100) : 0 }}% dari total
          </p>
        </CardContent>
      </Card>

      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Usia Rata-rata</CardTitle>
          <Calendar class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">2.1</div>
          <p class="text-xs text-muted-foreground">Tahun (estimasi)</p>
        </CardContent>
      </Card>
    </div>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <Baby class="h-5 w-5" />
          Daftar Balita
        </CardTitle>
        <CardDescription>
          Data balita yang terdaftar dalam sistem pemantauan stunting. Termasuk informasi orang tua, data lahir, dan lokasi.
        </CardDescription>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTable :data="balitaData" />
      </CardContent>
    </Card>

    <!-- Dialog Form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :balita="selectedBalita"
      @close="showDialog = false"
      @save="handleSave" />
  </div>
</template>