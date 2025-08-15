<script setup lang="ts">
import { ref } from "vue";
import type { PetugasKesehatan } from "./columns"
import { toast } from "vue-sonner";
import { Plus, Stethoscope } from "lucide-vue-next";
import Button from "@/components/ui/button/Button.vue";
import Card from "@/components/ui/card/Card.vue";
import CardHeader from "@/components/ui/card/CardHeader.vue";
import CardTitle from "@/components/ui/card/CardTitle.vue";
import CardContent from "@/components/ui/card/CardContent.vue";
import CardDescription from "@/components/ui/card/CardDescription.vue";
import DataTable from "./DataTable.vue";

const petugasData = ref<PetugasKesehatan[]>([
  {
    id: "1",
    id_pengguna: "1",
    id_skpd: "1",
    email: "petugas1@example.com",
    nama: "Petugas 1",
    skpd: "Puskesmas 1",
    jenis_skpd: "puskesmas",
    intervensi_count: 3,
    created_date: "2023-01-01",
    updated_date: "2023-01-02",
  },
  {
    id: "2",
    id_pengguna: "2",
    id_skpd: "2",
    email: "petugas2@example.com",
    nama: "Petugas 2",
    skpd: "Kelurahan 1",
    jenis_skpd: "kelurahan",
    intervensi_count: 0,
    created_date: "2023-01-01",
    updated_date: "2023-01-02",
  },
  {
    id: "3",
    id_pengguna: "3",
    id_skpd: "3",
    email: "petugas3@example.com",
    nama: "Petugas 3",
    skpd: "SKPD 1",
    jenis_skpd: "skpd",
    intervensi_count: 5,
    created_date: "2023-01-01",
    updated_date: "2023-01-02",
  },
])

const showDialog = ref(false)
const dialogMode = ref<"create" | "edit">("create")
const selectedPetugas = ref<PetugasKesehatan | null>(null)

const totalPetugas = ref(petugasData.value.length)

const updateStatistics = () => {
  totalPetugas.value = petugasData.value.length
}

const handleCreate= () => {
  dialogMode.value = "create"
  selectedPetugas.value = null
  showDialog.value = true
}

const handleEdit = (petugas: PetugasKesehatan) => {
  dialogMode.value = "edit"
  selectedPetugas.value = { ...petugas }
  showDialog.value = true
}

const handleDelete = (petugas: PetugasKesehatan) => {
  if (confirm(`Apakah Anda yakin ingin menghapus petugas ${petugas.nama}?`)) {
    const index = petugasData.value.findIndex((p) => p.id === petugas.id)
    if (index > -1) {
      petugasData.value.splice(index, 1)
      updateStatistics()
      toast.success("Data petugas berhasil dihapus")
    }
  }
}

const handleSave = (petugas: PetugasKesehatan) => {
  if (dialogMode.value === "create") {
    // Generate new ID
    const newId = (Math.max(...petugasData.value.map((p) => parseInt(p.id))) + 1).toString()
    const newPetugas = {
      ...petugas,
      id: newId,
      created_date: new Date().toISOString().split("T")[0],
    }
    petugasData.value.unshift(newPetugas)
    toast.success("Data petugas berhasil ditambahkan")
  } else {
    // Update existing
    const index = petugasData.value.findIndex((p) => p.id === petugas.id)
    if (index > -1) {
      petugasData.value[index] = {
        ...petugas,
        updated_date: new Date().toISOString().split("T")[0],
      }
      toast.success("Data petugas berhasil diperbarui")
    }
  }

  updateStatistics()
  showDialog.value = false
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Stethoscope class="h-8 w-8 text-blue-600" />
          Data Petugas Kesehatan
        </h1>
        <p class="text-gray-600">Kelola data petugas kesehatan di Kota Cirebon</p>
      </div>
      <Button @click="handleCreate" class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah Petugas Kesehatan
      </Button>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Petugas</CardTitle>
          <Stethoscope class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalPetugas }}</div>
          <p class="text-xs text-muted-foreground">Petugas terdaftar</p>
        </CardContent>
      </Card>
    </div>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <Stethoscope class="h-5 w-5" />
          Daftar Petugas
        </CardTitle>
        <CardDescription>
          Data petugas yang terdaftar dalam sistem pemantauan stunting. Termasuk informasi kontak, jabatan, dan lokasi.
        </CardDescription>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTable :data="petugasData" />
      </CardContent>
    </Card>
  </div>
</template>
