<script setup lang="ts">
import { PlusCircle, Database, BarChart } from "lucide-vue-next"
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

import AdminSidebarUser from "./AdminSidebarUser.vue"

const items = [
  {
    title: "Laporan Masyarakat",
    icon: Database,
  },
  {
    title: "Tambah Balita Stunting",
    icon: PlusCircle,
  },
  {
    title: "Statistik",
    icon: BarChart,
  },
]

const props = defineProps<{ activeMenu: string }>()

const emit = defineEmits<{
  (e: "menu-change", menu: string): void
}>()

function selectMenu(menu: string) {
  emit("menu-change", menu)
}
</script>

<template>
  <Sidebar>
    <SidebarHeader>
      <div class="flex items-center gap-2 my-2">
        <img
          src="@/assets/logo_dkis.png"
          alt="Logo"
          class="h-12" />
        <p class="text-xl font-bold">DKIS Kota Cirebon</p>
      </div>
      <div>
        <p class="text-lg font-semibold">Dashboard Admin</p>
        <p class="text-sm text-muted-foreground">
          Sistem informasi pemantauan dan analisis data stunting di Kota Cirebon.
        </p>
      </div>
    </SidebarHeader>
    <SidebarContent>
      <SidebarGroup>
        <SidebarGroupLabel>Menu Utama</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem
              v-for="item in items"
              :key="item.title">
              <SidebarMenuButton
                asChild
                :isActive="props.activeMenu === item.title">
                <button @click.prevent="selectMenu(item.title)">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </button>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>
    <SidebarFooter>
      <AdminSidebarUser
        :user="{
          name: 'Rifqi Daiva Tri Nandhika',
          email: 'rifqi.daiva@example.com',
          avatar: 'https://via.placeholder.com/150',
        }" />
    </SidebarFooter>
  </Sidebar>
</template>
