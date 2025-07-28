<script setup lang="ts">
import { ref } from "vue"
import { SidebarProvider, SidebarTrigger, SidebarInset } from "@/components/ui/sidebar"
import AdminSidebar from "@/components/admin/AdminSidebar.vue"
import AdminAddSufferer from "./menu/AdminAddSufferer.vue"
import AdminData from "./menu/AdminData.vue"
import AdminStatistics from "./menu/AdminStatistics.vue"

const activeMenu = ref("Tambah Penderita")
function handleMenuChange(menu: string) {
  activeMenu.value = menu
}
</script>

<template>
  <SidebarProvider>
    <AdminSidebar :activeMenu="activeMenu" @menu-change="handleMenuChange" />
    <SidebarInset>
      <main class="p-4">
        <SidebarTrigger class="my-2" />
        <div v-if="activeMenu === 'Tambah Penderita'">
          <AdminAddSufferer />
        </div>
        <div v-else-if="activeMenu === 'Data'">
          <AdminData />
        </div>
        <div v-else-if="activeMenu === 'Statistik'">
          <AdminStatistics />
        </div>
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>
