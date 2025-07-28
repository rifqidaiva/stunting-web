<script setup lang="ts">
import { onMounted, ref } from "vue"

const data = ref(null)

onMounted(async () => {
  try {
    const response = await fetch("/api/test")
    data.value = await response.json()
    console.log(data)
  } catch (error) {
    console.error("API request failed:", error)
  }
})
</script>

<template>
  <h1>Home Page</h1>
  <pre>{{ data }}</pre>
</template>
