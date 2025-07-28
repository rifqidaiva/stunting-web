<script setup lang="ts">
import { onMounted, onBeforeUnmount } from "vue"
import L from "leaflet"
import "leaflet/dist/leaflet.css"

let map: L.Map | null = null


onMounted(() => {
  map = L.map("map").setView([-7.797068, 110.370529], 13)

  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    attribution:
      '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
  }).addTo(map)

  L.marker([-7.797068, 110.370529]).addTo(map).bindPopup("Lokasi Penderita").openPopup()
})

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<template>
  <div>
    <h1
      class="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0">
      Tambah Penderita
    </h1>
    <section class="grid grid-cols-1 gap-4 md:grid-cols-5 my-2">
      <div
        id="map"
        style="height: 400px"
        class="z-0 rounded-md md:col-span-3"></div>
      <div>
        Lorem ipsum, dolor sit amet consectetur adipisicing elit. A corrupti officiis labore dolor
        molestias nemo!
      </div>
    </section>
  </div>
</template>
