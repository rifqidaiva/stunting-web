<script setup lang="ts">
import { ref } from "vue"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"

// Form data
const formData = ref({
  email: "",
  name: "",
  alamat: "",
  password: "",
  confirmPassword: "",
})

// Validation errors
const errors = ref<Record<string, string>>({})

// Validate passwords match
const validatePasswords = () => {
  if (formData.value.password !== formData.value.confirmPassword) {
    errors.value.confirmPassword = "Password tidak cocok"
  } else {
    delete errors.value.confirmPassword
  }
}

// Handle form submission
const handleSubmit = () => {
  errors.value = {}

  // Basic validation
  if (!formData.value.email) {
    errors.value.email = "Email harus diisi"
  }
  if (!formData.value.name) {
    errors.value.name = "Nama harus diisi"
  }
  if (!formData.value.alamat) {
    errors.value.alamat = "Alamat harus diisi"
  }
  if (!formData.value.password) {
    errors.value.password = "Password harus diisi"
  }
  if (!formData.value.confirmPassword) {
    errors.value.confirmPassword = "Konfirmasi password harus diisi"
  }

  validatePasswords()

  // If no errors, proceed with registration
  if (Object.keys(errors.value).length === 0) {
    console.log("Registration data:", {
      email: formData.value.email,
      name: formData.value.name,
      alamat: formData.value.alamat,
      password: formData.value.password,
    })
    // TODO: Implement actual registration logic
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <Card class="w-full max-w-md">
      <CardHeader>
        <CardTitle class="text-2xl">Daftar</CardTitle>
        <CardDescription>
          Masukkan informasi Anda di bawah ini untuk membuat akun baru
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form
          @submit.prevent="handleSubmit"
          class="grid gap-4">
          <!-- Email -->
          <div class="grid gap-2">
            <Label for="email">Email *</Label>
            <Input
              id="email"
              v-model="formData.email"
              type="email"
              placeholder="m@example.com"
              :class="errors.email ? 'border-red-500' : ''"
              required />
            <p
              v-if="errors.email"
              class="text-sm text-red-500">
              {{ errors.email }}
            </p>
          </div>

          <!-- Nama -->
          <div class="grid gap-2">
            <Label for="name">Nama Lengkap *</Label>
            <Input
              id="name"
              v-model="formData.name"
              type="text"
              placeholder="Masukkan nama lengkap Anda"
              :class="errors.name ? 'border-red-500' : ''"
              required />
            <p
              v-if="errors.name"
              class="text-sm text-red-500">
              {{ errors.name }}
            </p>
          </div>

          <!-- Alamat -->
          <div class="grid gap-2">
            <Label for="alamat">Alamat *</Label>
            <Textarea
              id="alamat"
              v-model="formData.alamat"
              placeholder="Masukkan alamat lengkap Anda"
              :class="errors.alamat ? 'border-red-500' : ''"
              class="min-h-[80px] resize-none"
              required />
            <p
              v-if="errors.alamat"
              class="text-sm text-red-500">
              {{ errors.alamat }}
            </p>
          </div>

          <!-- Password -->
          <div class="grid gap-2">
            <Label for="password">Password *</Label>
            <Input
              id="password"
              v-model="formData.password"
              type="password"
              placeholder="Masukkan password"
              :class="errors.password ? 'border-red-500' : ''"
              required />
            <p
              v-if="errors.password"
              class="text-sm text-red-500">
              {{ errors.password }}
            </p>
          </div>

          <!-- Konfirmasi Password -->
          <div class="grid gap-2">
            <Label for="confirmPassword">Konfirmasi Password *</Label>
            <Input
              id="confirmPassword"
              v-model="formData.confirmPassword"
              @blur="validatePasswords"
              type="password"
              placeholder="Ulangi password"
              :class="errors.confirmPassword ? 'border-red-500' : ''"
              required />
            <p
              v-if="errors.confirmPassword"
              class="text-sm text-red-500">
              {{ errors.confirmPassword }}
            </p>
          </div>

          <!-- Submit Button -->
          <Button
            type="submit"
            class="w-full">
            Daftar
          </Button>
        </form>

        <!-- Login Link -->
        <div class="mt-4 text-center text-sm">
          Sudah punya akun?
          <router-link
            to="/auth/login"
            class="underline hover:text-primary">
            Masuk
          </router-link>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
