<script setup lang="ts">
import { ref } from "vue";
import type { Intervensi, RiwayatPemeriksaan } from "./columns"
import { Activity, Plus } from "lucide-vue-next";
import Button from "@/components/ui/button/Button.vue";
import Card from "@/components/ui/card/Card.vue";
import CardHeader from "@/components/ui/card/CardHeader.vue";
import CardTitle from "@/components/ui/card/CardTitle.vue";
import CardContent from "@/components/ui/card/CardContent.vue";
import CardDescription from "@/components/ui/card/CardDescription.vue";
import DataTableIntervensi from "./DataTableIntervensi.vue";

// Dummy data
const intervensiData = ref<Intervensi[]>([
  {
    id: "1",
    id_balita: "1",
    nama_balita: "Andi Pratama",
    jenis: "Nutrisi",
    tanggal: "2024-01-24",
    deskripsi: "Pemberian suplemen vitamin A",
    hasil: "Baik",
    petugas_count: 2,
    riwayat_count: 1,
    created_date: "2024-01-19",
    updated_date: "2024-01-20",
    created_by: "admin"
  },
  {
    id: "2",
    id_balita: "2",
    nama_balita: "Sari Cantika",
    jenis: "Kesehatan",
    tanggal: "2024-01-22",
    deskripsi: "Pemeriksaan kesehatan rutin",
    hasil: "Baik",
    petugas_count: 1,
    riwayat_count: 1,
    created_date: "2024-01-21",
  },
])

const riwayatPemeriksaanData = ref<RiwayatPemeriksaan[]>([
  {
    id: "1",
    id_balita: "1",
    nama_balita: "Andi Pratama",
    umur_balita: "5",
    jenis_kelamin: "L",
    nama_ayah: "Budi Pratama",
    nama_ibu: "Siti Aminah",
    nomor_kk: "1234567890123",
    id_intervensi: "1",
    jenis_intervensi: "Nutrisi",
    tanggal_intervensi: "2024-01-20",
    id_laporan_masyarakat: "1",
    status_laporan: "Diterima",
    tanggal_laporan: "2024-01-19",
    jenis_laporan: "masyarakat",
    tanggal: "2024-01-20",
    berat_badan: "15",
    tinggi_badan: "100",
    status_gizi: "normal",
    keterangan: "Tidak ada keluhan",
    kelurahan: "Kelurahan 1",
    kecamatan: "Kecamatan 1",
    created_date: "2024-01-19",
    updated_date: "2024-01-20",
    created_by: "admin",
    updated_by: "admin",
  },
  {
    id: "2",
    id_balita: "2",
    nama_balita: "Sari Cantika",
    umur_balita: "4",
    jenis_kelamin: "P",
    nama_ayah: "Joko Santoso",
    nama_ibu: "Dewi Sartika",
    nomor_kk: "1234567890124",
    id_intervensi: "2",
    jenis_intervensi: "Kesehatan",
    tanggal_intervensi: "2024-01-22",
    id_laporan_masyarakat: "2",
    status_laporan: "Diterima",
    tanggal_laporan: "2024-01-21",
    jenis_laporan: "masyarakat",
    tanggal: "2024-01-22",
    berat_badan: "14",
    tinggi_badan: "95",
    status_gizi: "stunting",
    keterangan: "Perlu perhatian lebih",
    kelurahan: "Kelurahan 2",
    kecamatan: "Kecamatan 2",
    created_date: "2024-01-21",
    updated_date: "2024-01-22",
    created_by: "admin",
    updated_by: "admin",
  },
])

const totalIntervensi = ref(intervensiData.value.length)
const totalRiwayatPemeriksaan = ref(riwayatPemeriksaanData.value.length)

const updateStatistics = () => {
  totalIntervensi.value = intervensiData.value.length
  totalRiwayatPemeriksaan.value = riwayatPemeriksaanData.value.length
}
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
          <Activity class="h-8 w-8 text-blue-600" />
          Riwayat & Intervensi
        </h1>
        <p class="text-gray-600">Kelola data riwayat dan intervensi balita stunting</p>
      </div>
      <Button @click="" class="gap-2">
        <Plus class="h-4 w-4" />
        Tambah Intervensi
      </Button>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Intervensi</CardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalIntervensi }}</div>
          <p class="text-xs text-muted-foreground">Banyak intervensi yang dilakukan</p>
        </CardContent>
      </Card>
      <Card>
        <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
          <CardTitle class="text-sm font-medium">Total Riwayat Pemeriksaan</CardTitle>
          <Activity class="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div class="text-2xl font-bold">{{ totalRiwayatPemeriksaan }}</div>
          <p class="text-xs text-muted-foreground">Banyak riwayat pemeriksaan yang dilakukan</p>
        </CardContent>
      </Card>
    </div>

    <!-- Data Table -->
    <Card>
      <CardHeader>
        <CardTitle class="flex items-center gap-2">
          <Activity class="h-5 w-5" />
          Daftar Intervensi
        </CardTitle>
        <CardDescription>
          Data intervensi yang dilakukan terhadap balita stunting.
        </CardDescription>
      </CardHeader>
      <CardContent class="overflow-auto">
        <DataTableIntervensi :data="intervensiData" />
      </CardContent>
    </Card>

  </div>
</template>
