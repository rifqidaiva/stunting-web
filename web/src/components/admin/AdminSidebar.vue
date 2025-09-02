<script setup lang="ts">
import { ref, onMounted } from "vue"
import { Users, Baby, FileText, Activity, Stethoscope, Building } from "lucide-vue-next"

import {
  Sidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar"

import AdminSidebarUser from "@/components/admin/AdminSidebarUser.vue"
import { authUtils } from "@/lib/utils"

const menuItems = [
  {
    category: "Manajemen Data",
    items: [
      {
        title: "Keluarga",
        icon: Users,
        route: "/admin/keluarga",
      },
      {
        title: "Balita",
        icon: Baby,
        route: "/admin/balita",
      },
      {
        title: "Data Pelapor (Masyarakat)",
        icon: FileText,
        route: "/admin/pelapor"
      },
      {
        title: "Laporan Masyarakat",
        icon: FileText,
        route: "/admin/laporan-masyarakat",
      },
      {
        title: "Intervensi",
        icon: Activity,
        route: "/admin/intervensi",
      },
      {
        title: "Riwayat Pemeriksaan",
        icon: FileText,
        route: "/admin/riwayat-pemeriksaan",
      },
    ],
  },
  {
    category: "SKPD",
    items: [
      {
        title: "SKPD",
        icon: Building,
        route: "/admin/skpd",
      },
      {
        title: "Petugas Kesehatan",
        icon: Stethoscope,
        route: "/admin/petugas-kesehatan",
      },
    ],
  },
]

// User data state
const currentUser = ref({
  name: "Admin DKIS",
  email: "admin@dkis-cirebon.go.id",
  avatar: "/api/placeholder/32/32",
  role: "admin",
})

// Load user data from localStorage using authUtils
const loadUserData = () => {
  try {
    const userData = authUtils.getUserData()
    if (userData) {
      currentUser.value = {
        name: userData.data?.nama || userData.nama || "Admin DKIS",
        email: userData.email || "admin@dkis-cirebon.go.id",
        avatar: "/api/placeholder/32/32", // Default avatar
        role: userData.role || "admin",
      }
      console.log("User data loaded from authUtils:", currentUser.value)
    } else {
      console.log("No user data found in authUtils, using defaults")
    }
  } catch (error) {
    console.error("Error loading user data from authUtils:", error)
  }
}

// Load user data on component mount
onMounted(() => {
  loadUserData()
})
</script>

<template>
  <Sidebar>
    <SidebarHeader>
      <div class="flex items-center gap-2 my-2">
        <img
          src="@/assets/logo_dkis.png"
          alt="Logo"
          class="h-12" />
        <div class="flex-1">
          <p class="text-lg font-bold text-primary">DKIS Kota Cirebon</p>
          <p class="text-xs text-muted-foreground">
            Dinas Komunikasi, Informatika, dan Statistik Kota Cirebon
          </p>
        </div>
      </div>
      <div class="p-3">
        <p class="text-sm font-semibold text-foreground">Dashboard Admin</p>
        <p class="text-xs text-muted-foreground leading-relaxed">
          Sistem pemantauan dan analisis data stunting di Kota Cirebon
        </p>
      </div>
    </SidebarHeader>

    <SidebarContent class="px-2">
      <div
        v-for="menuGroup in menuItems"
        :key="menuGroup.category">
        <SidebarGroup>
          <SidebarGroupLabel class="text-xs font-medium text-muted-foreground">
            {{ menuGroup.category }}
          </SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu class="">
              <SidebarMenuItem
                v-for="item in menuGroup.items"
                :key="item.route">
                <SidebarMenuButton asChild>
                  <router-link
                    :to="item.route"
                    class="w-full justify-start gap-3 px-3 py-2 text-sm transition-all hover:bg-accent hover:text-accent-foreground rounded-md flex items-center"
                    active-class="bg-accent text-accent-foreground font-medium"
                    exact-active-class="bg-accent text-accent-foreground font-medium">
                    <component
                      :is="item.icon"
                      class="h-4 w-4 shrink-0" />
                    <span class="truncate">{{ item.title }}</span>
                  </router-link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </div>
    </SidebarContent>

    <SidebarFooter>
      <div class="p-2 border-t">
        <AdminSidebarUser :user="currentUser" />
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
