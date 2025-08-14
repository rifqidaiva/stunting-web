<script setup lang="ts">
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-4">
    <Card class="w-full max-w-sm">
      <CardHeader>
        <CardTitle class="text-2xl">Masuk</CardTitle>
        <CardDescription>
          Masukkan email Anda di bawah ini untuk masuk ke akun Anda
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div class="grid gap-4">
          <div class="grid gap-2">
            <Label for="email">Email</Label>
            <Input
              id="email"
              type="email"
              placeholder="m@example.com"
              required />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center">
              <Label for="password">Password</Label>
            </div>
            <Input
              id="password"
              type="password"
              required />
          </div>
          <div class="grid gap-2">
            <Label for="address">Alamat</Label>
            <Input
              id="address"
              type="text"
              placeholder="Masukkan alamat Anda"
              required />
          </div>
          <Button
            type="submit"
            class="w-full">
            Masuk
          </Button>
        </div>
        <div class="mt-4 text-center text-sm">
          Belum punya akun?
          <router-link
            to="/auth/register"
            class="underline">
            Daftar
          </router-link>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
