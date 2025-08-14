<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue"
import { toast } from "vue-sonner"
import { Plus, FileText } from "lucide-vue-next"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import type { LaporanMasyarakat } from "./columns"
import DataTable from "./DataTable.vue"
import DialogForm from "./DialogForm.vue"

// Dummy data - Replace with real API calls
const laporanData: LaporanMasyarakat[] = [
  {
    id: "1",
    id_masyarakat: "1",
    id_balita: "1",
    id_status_laporan: "1",
    tanggal_laporan: "2024-08-10",
    hubungan_dengan_balita: "Ibu",
    nomor_hp_pelapor: "081234567890",
    nomor_hp_keluarga_balita: "081234567891",
    created_date: "2024-08-10T10:00:00Z",
    updated_date: "2024-08-11T15:30:00Z",
    // Extended fields
    nama_pelapor: "Siti Nurhaliza",
    email_pelapor: "siti@example.com",
    nama_balita: "Ahmad Fauzi",
    nama_ayah: "Budi Santoso",
    nama_ibu: "Siti Nurhaliza",
    nomor_kk: "1234567890123456",
    alamat: "Jl. Merdeka No. 123",
    kelurahan: "Kejaksan",
    kecamatan: "Kejaksan",
    status_laporan: "Belum Diproses",
    jenis_laporan: "Gizi Buruk",
  },
  {
    id: "2",
    id_masyarakat: "2",
    id_balita: "2",
    id_status_laporan: "3",
    tanggal_laporan: "2024-08-09",
    hubungan_dengan_balita: "Tetangga",
    nomor_hp_pelapor: "081234567892",
    nomor_hp_keluarga_balita: "081234567893",
    created_date: "2024-08-09T14:00:00Z",
    updated_date: "2024-08-12T09:15:00Z",
    // Extended fields
    nama_pelapor: "Ali Pelapor",
    email_pelapor: "ali@example.com",
    nama_balita: "Fatimah Zahra",
    nama_ayah: "Ali Rahman",
    nama_ibu: "Khadijah",
    nomor_kk: "9876543210987654",
    alamat: "Jl. Sudirman No. 456",
    kelurahan: "Pekalangan",
    kecamatan: "Pekalangan",
    status_laporan: "Diproses dan Data Sesuai",
    jenis_laporan: "Stunting",
  },
  {
    id: "3",
    id_masyarakat: "3",
    id_balita: "3",
    id_status_laporan: "2",
    tanggal_laporan: "2024-08-08",
    hubungan_dengan_balita: "Ayah",
    nomor_hp_pelapor: "081234567894",
    nomor_hp_keluarga_balita: "081234567895",
    created_date: "2024-08-08T09:00:00Z",
    updated_date: "2024-08-10T16:45:00Z",
    // Extended fields
    nama_pelapor: "Rahman Hidayat",
    email_pelapor: "rahman@example.com",
    nama_balita: "Aisyah Putri",
    nama_ayah: "Rahman Hidayat",
    nama_ibu: "Dewi Sari",
    nomor_kk: "5555666677778888",
    alamat: "Jl. Diponegoro No. 789",
    kelurahan: "Harjamukti",
    kecamatan: "Harjamukti",
    status_laporan: "Sedang Diproses",
    jenis_laporan: "Gizi Kurang",
  },
  {
    id: "4",
    id_masyarakat: "4",
    id_balita: "4",
    id_status_laporan: "1",
    tanggal_laporan: "2024-08-07",
    hubungan_dengan_balita: "Nenek",
    nomor_hp_pelapor: "081234567896",
    nomor_hp_keluarga_balita: "081234567897",
    created_date: "2024-08-07T11:30:00Z",
    updated_date: "2024-08-07T11:30:00Z",
    // Extended fields
    nama_pelapor: "Nurhayati",
    email_pelapor: "nurhayati@example.com",
    nama_balita: "Muhammad Rizki",
    nama_ayah: "Agus Setiawan",
    nama_ibu: "Maya Sari",
    nomor_kk: "1111222233334444",
    alamat: "Jl. Ahmad Yani No. 321",
    kelurahan: "Lemahwungkuk",
    kecamatan: "Lemahwungkuk",
    status_laporan: "Belum Diproses",
    jenis_laporan: "Stunting",
  },
  {
    id: "5",
    id_masyarakat: "5",
    id_balita: "5",
    id_status_laporan: "3",
    tanggal_laporan: "2024-08-06",
    hubungan_dengan_balita: "Ibu",
    nomor_hp_pelapor: "081234567898",
    nomor_hp_keluarga_balita: "081234567899",
    created_date: "2024-08-06T13:15:00Z",
    updated_date: "2024-08-13T10:20:00Z",
    // Extended fields
    nama_pelapor: "Rina Marlina",
    email_pelapor: "rina@example.com",
    nama_balita: "Zahra Amelia",
    nama_ayah: "Dedi Kurniawan",
    nama_ibu: "Rina Marlina",
    nomor_kk: "7777888899990000",
    alamat: "Jl. Kartini No. 654",
    kelurahan: "Kecapi",
    kecamatan: "Harjamukti",
    status_laporan: "Diproses dan Data Sesuai",
    jenis_laporan: "Gizi Buruk",
  },
  {
    id: "6",
    id_masyarakat: "6",
    id_balita: "6",
    id_status_laporan: "2",
    tanggal_laporan: "2024-08-05",
    hubungan_dengan_balita: "Paman",
    nomor_hp_pelapor: "081234567800",
    nomor_hp_keluarga_balita: "081234567801",
    created_date: "2024-08-05T15:45:00Z",
    updated_date: "2024-08-11T14:30:00Z",
    // Extended fields
    nama_pelapor: "Bambang Sutrisno",
    email_pelapor: "bambang@example.com",
    nama_balita: "Kevin Pratama",
    nama_ayah: "Andi Wijaya",
    nama_ibu: "Sari Dewi",
    nomor_kk: "2222333344445555",
    alamat: "Jl. Pemuda No. 987",
    kelurahan: "Pegambiran",
    kecamatan: "Lemahwungkuk",
    status_laporan: "Sedang Diproses",
    jenis_laporan: "Gizi Kurang",
  },
  {
    id: "7",
    id_masyarakat: "7",
    id_balita: "7",
    id_status_laporan: "1",
    tanggal_laporan: "2024-08-04",
    hubungan_dengan_balita: "Ibu",
    nomor_hp_pelapor: "081234567802",
    nomor_hp_keluarga_balita: "081234567803",
    created_date: "2024-08-04T08:00:00Z",
    updated_date: "2024-08-04T08:00:00Z",
    // Extended fields
    nama_pelapor: "Indah Permata",
    email_pelapor: "indah@example.com",
    nama_balita: "Nabila Azzahra",
    nama_ayah: "Rudi Hartono",
    nama_ibu: "Indah Permata",
    nomor_kk: "6666777788889999",
    alamat: "Jl. Gajah Mada No. 159",
    kelurahan: "Kesambi",
    kecamatan: "Kesambi",
    status_laporan: "Belum Diproses",
    jenis_laporan: "Stunting",
  },
  {
    id: "8",
    id_masyarakat: "8",
    id_balita: "8",
    id_status_laporan: "3",
    tanggal_laporan: "2024-08-03",
    hubungan_dengan_balita: "Ayah",
    nomor_hp_pelapor: "081234567804",
    nomor_hp_keluarga_balita: "081234567805",
    created_date: "2024-08-03T12:30:00Z",
    updated_date: "2024-08-14T11:45:00Z",
    // Extended fields
    nama_pelapor: "Teguh Raharjo",
    email_pelapor: "teguh@example.com",
    nama_balita: "Aditya Prasetyo",
    nama_ayah: "Teguh Raharjo",
    nama_ibu: "Lestari Wulan",
    nomor_kk: "3333444455556666",
    alamat: "Jl. Veteran No. 753",
    kelurahan: "Pulasaren",
    kecamatan: "Pekalangan",
    status_laporan: "Diproses dan Data Sesuai",
    jenis_laporan: "Gizi Buruk",
  },
  {
    id: "9",
    id_masyarakat: "9",
    id_balita: "9",
    id_status_laporan: "2",
    tanggal_laporan: "2024-08-02",
    hubungan_dengan_balita: "Bibi",
    nomor_hp_pelapor: "081234567806",
    nomor_hp_keluarga_balita: "081234567807",
    created_date: "2024-08-02T16:20:00Z",
    updated_date: "2024-08-09T13:15:00Z",
    // Extended fields
    nama_pelapor: "Yuni Astuti",
    email_pelapor: "yuni@example.com",
    nama_balita: "Putri Cantika",
    nama_ayah: "Wawan Setiadi",
    nama_ibu: "Fitri Handayani",
    nomor_kk: "8888999900001111",
    alamat: "Jl. Sukarno Hatta No. 258",
    kelurahan: "Larangan",
    kecamatan: "Harjamukti",
    status_laporan: "Sedang Diproses",
    jenis_laporan: "Gizi Kurang",
  },
  {
    id: "10",
    id_masyarakat: "10",
    id_balita: "10",
    id_status_laporan: "1",
    tanggal_laporan: "2024-08-01",
    hubungan_dengan_balita: "Ibu",
    nomor_hp_pelapor: "081234567808",
    nomor_hp_keluarga_balita: "081234567809",
    created_date: "2024-08-01T07:45:00Z",
    updated_date: "2024-08-01T07:45:00Z",
    // Extended fields
    nama_pelapor: "Tri Wahyuni",
    email_pelapor: "tri@example.com",
    nama_balita: "Bayu Anggara",
    nama_ayah: "Eko Prasetya",
    nama_ibu: "Tri Wahyuni",
    nomor_kk: "4444555566667777",
    alamat: "Jl. Dewi Sartika No. 147",
    kelurahan: "Argasunya",
    kecamatan: "Harjamukti",
    status_laporan: "Belum Diproses",
    jenis_laporan: "Stunting",
  },
  {
    id: "11",
    id_masyarakat: "11",
    id_balita: "11",
    id_status_laporan: "1",
    tanggal_laporan: "2024-08-01",
    hubungan_dengan_balita: "Ibu",
    nomor_hp_pelapor: "081234567810",
    nomor_hp_keluarga_balita: "081234567811",
    created_date: "2024-08-01T07:45:00Z",
    updated_date: "2024-08-01T07:45:00Z",
    // Extended fields
    nama_pelapor: "Tri Wahyuni",
    email_pelapor: "tri@example.com",
    nama_balita: "Bayu Anggara",
    nama_ayah: "Eko Prasetya",
    nama_ibu: "Tri Wahyuni",
    nomor_kk: "4444555566667777",
    alamat: "Jl. Dewi Sartika No. 147",
    kelurahan: "Argasunya",
    kecamatan: "Harjamukti",
    status_laporan: "Belum Diproses",
    jenis_laporan: "Stunting",
  },
]

const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedLaporan = ref<LaporanMasyarakat | null>(null)

// Statistic
const totalLaporan = ref(0)

const updateStatistics = () => {
  totalLaporan.value = laporanData.length
}

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

const handleDelete = (laporan: LaporanMasyarakat) => {
  if (confirm(`Apakah Anda yakin ingin menghapus laporan dengan ID ${laporan.id}?`)) {
    // Logic to delete the report

    toast.success("Laporan berhasil dihapus")
  }
}

const handleSave = (laporan: LaporanMasyarakat) => {
  if (dialogMode.value === "create") {
    // Logic to create a new report
    toast.success("Laporan berhasil dibuat")
  } else if (dialogMode.value === "edit" && selectedLaporan.value) {
    // Logic to update an existing report
    toast.success("Laporan berhasil diperbarui")
  }
  updateStatistics()
  showDialog.value = false
}

const handleCustomEvent = (event: Event) => {
  const customEvent = event as CustomEvent

  if (event.type === "edit-laporan") {
    handleEdit(customEvent.detail)
  } else if (event.type === "delete-laporan") {
    handleDelete(customEvent.detail)
  }
}

onMounted(() => {
  document.addEventListener("edit-laporan", handleCustomEvent)
  document.addEventListener("delete-laporan", handleCustomEvent)
})

onUnmounted(() => {
  document.removeEventListener("edit-laporan", handleCustomEvent)
  document.removeEventListener("delete-laporan", handleCustomEvent)
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <FileText class="h-8 w-8 text-slate-600" />
          Data Laporan
        </h1>
        <p class="text-gray-600">
          Kelola data laporan balita (anak usia di bawah 5 tahun) di Kota Cirebon
        </p>
      </div>
      <Button
        @click="handleCreate"
        class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah Laporan
      </Button>
    </div>

    <!-- Statistics card -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Laporan</CardTitle>
          <FileText class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalLaporan }}</div>
          <p class="text-xs text-muted-foreground">Laporan terdaftar</p>
        </CardContent>
      </Card>
    </div>

    <!-- Data table -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <FileText class="h-5 w-5" />
          Daftar Laporan
        </CardTitle>
        <CardDescription>
          Data laporan yang terdaftar dalam sistem pemantauan stunting. Termasuk informasi orang
          tua, data lahir, dan lokasi.
        </CardDescription>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTable :data="laporanData" />
      </CardContent>
    </Card>

    <!-- Dialog form -->
    <DialogForm
      :show="showDialog"
      :mode="dialogMode"
      :laporan="selectedLaporan"
      @close="showDialog = false"
      @save="handleSave" />
  </div>
</template>
