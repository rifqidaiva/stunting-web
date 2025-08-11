import { createApp } from "vue"
import { createRouter, createWebHistory } from "vue-router"

import "./style.css"
import App from "./App.vue"

import Home from "./components/Home.vue"
import Admin from "./components/admin/Admin.vue"
import Community from "./components/community/Community.vue"

import Keluarga from "./components/admin/menu/keluarga/Keluarga.vue"
import Balita from "./components/admin/menu/balita/Balita.vue"
import LaporanMasyarakat from "./components/admin/menu/laporan_masyarakat/LaporanMasyarakat.vue"
import RiwayatAndIntervensi from "./components/admin/menu/riwayat_&_intervensi/Riwayat&Intervensi.vue"
import PetugasKesehatan from "./components/admin/menu/petugas_kesehatan/PetugasKesehatan.vue"
import Skpd from "./components/admin/menu/skpd/Skpd.vue"

const routes = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/admin",
    component: Admin,
    children: [
      {
        path: "",
        redirect: "/admin/keluarga",
      },
      {
        path: "keluarga",
        component: Keluarga,
        meta: { title: "Admin - Keluarga" },
      },
      {
        path: "balita",
        component: Balita,
        meta: { title: "Admin - Balita" },
      },
      {
        path: "laporan-masyarakat",
        component: LaporanMasyarakat,
        meta: { title: "Admin - Laporan Masyarakat" },
      },
      {
        path: "riwayat-dan-intervensi",
        component: RiwayatAndIntervensi,
        meta: { title: "Admin - Riwayat dan Intervensi" },
      },
      {
        path: "petugas-kesehatan",
        component: PetugasKesehatan,
        meta: { title: "Admin - Petugas Kesehatan" },
      },
      {
        path: "skpd",
        component: Skpd,
        meta: { title: "Admin - SKPD" },
      },
    ],
  },
  {
    path: "/community",
    component: Community,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

createApp(App).use(router).mount("#app")
