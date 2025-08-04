<script setup lang="ts">
import { ref } from "vue"
import { SidebarProvider, SidebarTrigger, SidebarInset } from "@/components/ui/sidebar"
import AdminSidebar from "@/components/admin/AdminSidebar.vue"
import AdminAddStunting from "./menu/AdminAddStunting.vue"
import AdminCommunityReport from "./menu/AdminCommunityReport.vue"
import AdminStatistics from "./menu/AdminStatistics.vue"

const activeMenu = ref("Tambah Balita Stunting")
function handleMenuChange(menu: string) {
  activeMenu.value = menu
}
</script>

<template>
  <SidebarProvider>
    <AdminSidebar
      :activeMenu="activeMenu"
      @menu-change="handleMenuChange" />
    <SidebarInset>
      <main class="p-4">
        <SidebarTrigger class="my-2" />
        <div v-if="activeMenu === 'Laporan Masyarakat'">
          <AdminCommunityReport />
        </div>
        <div v-else-if="activeMenu === 'Tambah Balita Stunting'">
          <AdminAddStunting />
        </div>
        <div v-else-if="activeMenu === 'Statistik'">
          <AdminStatistics />
        </div>
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
