<script setup lang="ts">
import {
  Users,
  Baby,
  FileText,
  MapPin,
  Activity,
  Stethoscope,
  Building,
  Map,
  BarChart3,
} from "lucide-vue-next"

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
        title: "Laporan Masyarakat",
        icon: FileText,
        route: "/admin/laporan-masyarakat",
      },
      {
        title: "Riwayat & Intervensi",
        icon: Activity,
        route: "/admin/riwayat-dan-intervensi",
      },
    ],
  },
  {
    category: "Data Geografis",
    items: [
      {
        title: "Peta Stunting",
        icon: Map,
        route: "/admin/peta-stunting",
      },
      {
        title: "Area Wilayah",
        icon: MapPin,
        route: "/admin/area-wilayah",
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
  {
    category: "Analisis",
    items: [
      {
        title: "Statistik & Laporan",
        icon: BarChart3,
        route: "/admin/statistik",
      },
    ],
  },
]
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
        <AdminSidebarUser
          :user="{
            name: 'Admin DKIS',
            email: 'admin@dkis-cirebon.go.id',
            avatar: '/api/placeholder/32/32',
          }" />
      </div>
    </SidebarFooter>
  </Sidebar>
</template>
