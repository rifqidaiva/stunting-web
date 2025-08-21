<script setup lang="ts">
import { ref } from "vue"
import { toast } from "vue-sonner"
import { ChevronsUpDown, LogOut, User, Shield } from "lucide-vue-next"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"
import { authUtils } from "@/lib/utils"

const props = defineProps<{
  user: {
    name: string
    email: string
    avatar: string
    role?: string
  }
}>()

const { isMobile } = useSidebar()
const isLoggingOut = ref(false)

// Get actual user data from localStorage using authUtils
const getUserData = () => {
  try {
    const userData = authUtils.getUserData()
    if (userData) {
      return {
        name: userData.data?.nama || userData.nama || props.user.name,
        email: userData.email || props.user.email,
        avatar: props.user.avatar,
        role: userData.role || props.user.role || "admin",
      }
    }
  } catch (error) {
    console.error("Error getting user data:", error)
  }
  return props.user
}

const currentUser = getUserData()

// Handle logout using authUtils
const handleLogout = async (): Promise<void> => {
  if (isLoggingOut.value) return

  isLoggingOut.value = true

  try {
    console.log("Starting logout process...")

    // Show loading toast
    toast.info("Sedang logout...", { duration: 1000 })

    // Use authUtils.logout which handles both API call and cleanup
    await authUtils.logout("/auth/login")

    // Success message (this might not show if redirect happens immediately)
    toast.success("Logout berhasil. Sampai jumpa!")
  } catch (error) {
    console.error("Logout error:", error)

    // Fallback: force logout if regular logout fails
    authUtils.forceLogout("Logout berhasil (dengan peringatan)")
  } finally {
    isLoggingOut.value = false
  }
}

// Handle account settings (optional)
const handleAccount = (): void => {
  toast.info("Fitur pengaturan akun akan segera tersedia")
  // router.push("/admin/profile") // Uncomment when profile page is ready
}

// Get initials for avatar fallback
const getInitials = (name: string): string => {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2)
}

// Get role badge color
const getRoleColor = (role: string): string => {
  switch (role?.toLowerCase()) {
    case "admin":
      return "text-blue-600"
    case "petugas_kesehatan":
      return "text-green-600"
    case "masyarakat":
      return "text-orange-600"
    default:
      return "text-gray-600"
  }
}
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground">
            <Avatar class="h-8 w-8 rounded-lg">
              <AvatarImage
                :src="currentUser.avatar"
                :alt="currentUser.name" />
              <AvatarFallback class="rounded-lg bg-primary text-primary-foreground">
                {{ getInitials(currentUser.name) }}
              </AvatarFallback>
            </Avatar>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ currentUser.name }}</span>
              <span class="truncate text-xs">{{ currentUser.email }}</span>
            </div>
            <ChevronsUpDown class="ml-auto size-4" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
          :side="isMobile ? 'bottom' : 'right'"
          align="end"
          :side-offset="4">
          <DropdownMenuLabel class="p-0 font-normal">
            <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
              <Avatar class="h-8 w-8 rounded-lg">
                <AvatarImage
                  :src="currentUser.avatar"
                  :alt="currentUser.name" />
                <AvatarFallback class="rounded-lg bg-primary text-primary-foreground">
                  {{ getInitials(currentUser.name) }}
                </AvatarFallback>
              </Avatar>
              <div class="grid flex-1 text-left text-sm leading-tight">
                <span class="truncate font-semibold">{{ currentUser.name }}</span>
                <span class="truncate text-xs">{{ currentUser.email }}</span>
                <div class="flex items-center gap-1 mt-1">
                  <Shield
                    class="h-3 w-3"
                    :class="getRoleColor(currentUser.role || 'admin')" />
                  <span
                    class="text-xs font-medium"
                    :class="getRoleColor(currentUser.role || 'admin')">
                    {{ currentUser.role === "admin" ? "Administrator" : currentUser.role }}
                  </span>
                </div>
              </div>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuGroup>
            <DropdownMenuItem
              @click="handleAccount"
              class="cursor-pointer">
              <User class="mr-2 h-4 w-4" />
              <span>Pengaturan Akun</span>
            </DropdownMenuItem>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            @click="handleLogout"
            :disabled="isLoggingOut"
            class="cursor-pointer text-red-600 focus:text-red-600 focus:bg-red-50">
            <LogOut class="mr-2 h-4 w-4" />
            <span v-if="isLoggingOut">Logging out...</span>
            <span v-else>Logout</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </SidebarMenuItem>
  </SidebarMenu>
</template>
