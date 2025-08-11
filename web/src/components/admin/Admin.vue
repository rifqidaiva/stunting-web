<script setup lang="ts">
import { SidebarProvider, SidebarTrigger, SidebarInset } from "@/components/ui/sidebar"
import AdminSidebar from "@/components/admin/AdminSidebar.vue"
</script>

<template>
  <SidebarProvider>
    <AdminSidebar />
    <SidebarInset>
      <main class="p-4">
        <SidebarTrigger class="my-2" />
        <div class="relative overflow-hidden">
          <router-view v-slot="{ Component, route }">
            <transition
              name="fade"
              mode="out-in"
              appear
              :duration="150">
              <component
                :is="Component"
                :key="route.path"
                class="w-full" />
            </transition>
          </router-view>
        </div>
      </main>
    </SidebarInset>
  </SidebarProvider>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
