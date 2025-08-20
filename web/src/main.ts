import { createApp } from "vue"
import { createRouter, createWebHistory } from "vue-router"

import "./style.css"
import App from "./App.vue"

import Home from "./components/Home.vue"
import Auth from "./components/auth/Auth.vue"
import Admin from "./components/admin/Admin.vue"
import Community from "./components/community/Community.vue"

import Login from "./components/auth/Login.vue"
import Register from "./components/auth/Register.vue"

import Keluarga from "./components/admin/menu/keluarga/Keluarga.vue"
import Balita from "./components/admin/menu/balita/Balita.vue"
import LaporanMasyarakat from "./components/admin/menu/laporan_masyarakat/LaporanMasyarakat.vue"
import RiwayatAndIntervensi from "./components/admin/menu/riwayat_&_intervensi/Riwayat&Intervensi.vue"
import PetugasKesehatan from "./components/admin/menu/petugas_kesehatan/PetugasKesehatan.vue"
import Skpd from "./components/admin/menu/skpd/Skpd.vue"

import CommunityKeluarga from "./components/community/menu/keluarga/Keluarga.vue"
import CommunityBalita from "./components/community/menu/balita/Balita.vue"
import CommunityLaporanMasyarakat from "./components/community/menu/laporan_masyarakat/LaporanMasyarakat.vue"


const routes = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/auth",
    component: Auth,
    children: [
      {
        path: "",
        redirect: "/auth/login",
      },
      {
        path: "login",
        component: Login,
        meta: { title: "Login" },
      },
      {
        path: "register",
        component: Register,
        meta: { title: "Register" },
      },
    ],
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
    children: [
      {
        path: "keluarga",
        component: CommunityKeluarga,
        meta: { title: "Community - Keluarga" },
      },
      {
        path: "balita",
        component: CommunityBalita,
        meta: { title: "Community - Balita" },
      },
      {
        path: "laporan-masyarakat",
        component: CommunityLaporanMasyarakat,
        meta: { title: "Community - Laporan Masyarakat" },
      },
    ]
  },
]

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior() {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve({ left: 0, top: 0 })
      }, 150)
    })
  },
  routes,
})

// Title handling
router.beforeEach((to) => {
  document.title = (to.meta.title as string) || "Stunting Web"
})

const app = createApp(App)
app.use(router)
app.mount("#app")
