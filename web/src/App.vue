<script setup lang="ts">
import { RouterView } from "vue-router"
import { Toaster } from "vue-sonner"
import "vue-sonner/style.css"
</script>

<template>
  <div id="app">
    <RouterView />
    <Toaster class="pointer-events-auto" />
  </div>
</template>
