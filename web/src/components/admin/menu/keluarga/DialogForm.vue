<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from "vue"
import { toast } from "vue-sonner"
import L from "leaflet"
import "leaflet/dist/leaflet.css"

import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Save, X, MapPin, Navigation } from "lucide-vue-next"
import type { Keluarga } from "./columns"

// ✅ Fix Leaflet default marker icons
delete (L.Icon.Default.prototype as any)._getIconUrl
L.Icon.Default.mergeOptions({
  iconRetinaUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl: "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
})

interface Props {
  show: boolean
  mode: "create" | "edit"
  keluarga: Keluarga | null
}

interface Emits {
  (e: "close"): void
  (e: "save", keluarga: Keluarga): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Form data
const formData = ref<Partial<Keluarga>>({
  nomor_kk: "",
  nama_ayah: "",
  nama_ibu: "",
  nik_ayah: "",
  nik_ibu: "",
  alamat: "",
  rt: "",
  rw: "",
  id_kelurahan: "",
  koordinat: [-6.7064, 108.5492], // Default Cirebon coordinates
})

// Validation errors
const errors = ref<Record<string, string>>({})

// ✅ Map variables
const mapContainer = ref<HTMLElement>()
let map: L.Map | null = null
let marker: L.Marker | null = null
const mapInitialized = ref(false)

// Dummy kelurahan options (sesuai dengan data yang ada di main.go)
const kelurahanOptions = [
  { id: "1", label: "Kesambi - Kesambi" },
  { id: "2", label: "Tuparev - Kedawung" },
  { id: "3", label: "Perjuangan - Kejaksan" },
  { id: "4", label: "Argasunya - Harjamukti" },
  { id: "5", label: "Lemahwungkuk - Lemahwungkuk" },
  { id: "6", label: "Kesenden - Kesenden" },
]

// Helper functions untuk input restriction
const restrictToNumbers = (value: string) => value.replace(/\D/g, "")
const restrictToNumbers3Digit = (value: string) => value.replace(/\D/g, "").substring(0, 3)

// ✅ Map functions
const initMap = async () => {
  console.log("🗺️ Initializing map...")

  if (!mapContainer.value) {
    console.error("❌ Map container not found")
    return
  }

  if (mapInitialized.value) {
    console.log("⚠️ Map already initialized")
    updateMapWithCurrentCoordinates()
    return
  }

  try {
    // Destroy existing map first
    destroyMap()

    await nextTick()

    // Check container dimensions
    const containerRect = mapContainer.value.getBoundingClientRect()
    console.log("📏 Container dimensions:", containerRect)

    if (containerRect.width === 0 || containerRect.height === 0) {
      console.error("❌ Container has no dimensions, retrying...")
      setTimeout(() => initMap(), 500)
      return
    }

    console.log("🏗️ Creating map instance...")

    // Get current coordinates
    const currentCoords = formData.value.koordinat as [number, number]
    console.log("📍 Initial coordinates:", currentCoords)

    // Initialize map
    map = L.map(mapContainer.value, {
      center: currentCoords,
      zoom: 15,
      zoomControl: true,
      attributionControl: true,
    })

    // Add tile layer
    L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
      attribution: "© OpenStreetMap contributors",
      maxZoom: 19,
    }).addTo(map)

    // Add marker
    marker = L.marker(currentCoords, {
      draggable: true,
    }).addTo(map)

    // Bind marker events
    marker.on("dragend", (e) => {
      const position = e.target.getLatLng()
      formData.value.koordinat = [position.lat, position.lng]
      console.log("📍 Marker dragged to:", position)
    })

    // Bind map click events
    map.on("click", (e) => {
      const { lat, lng } = e.latlng
      formData.value.koordinat = [lat, lng]
      if (marker) {
        marker.setLatLng([lat, lng])
      }
      console.log("🖱️ Map clicked at:", { lat, lng })
    })

    // Handle map ready
    map.whenReady(() => {
      console.log("🗺️ Map is ready")
      setTimeout(() => {
        if (map) {
          map.invalidateSize()
          console.log("🔄 Map resized")
        }
      }, 100)
    })

    mapInitialized.value = true
    console.log("✅ Map initialized successfully")
  } catch (error) {
    console.error("❌ Error initializing map:", error)
    mapInitialized.value = false
  }
}

const destroyMap = () => {
  if (map) {
    try {
      console.log("🗑️ Destroying map...")
      map.remove()
      map = null
      marker = null
      mapInitialized.value = false
      console.log("✅ Map destroyed")
    } catch (error) {
      console.error("❌ Error destroying map:", error)
    }
  }
}

const updateMapWithCurrentCoordinates = () => {
  if (map && marker && formData.value.koordinat) {
    const [lat, lng] = formData.value.koordinat
    console.log("🔄 Updating map to coordinates:", [lat, lng])

    marker.setLatLng([lat, lng])
    map.setView([lat, lng], map.getZoom())
  }
}

// ✅ Get current location
const getCurrentLocation = () => {
  if ("geolocation" in navigator) {
    toast.info("Mendapatkan lokasi...")

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords
        formData.value.koordinat = [latitude, longitude]

        if (map && marker) {
          marker.setLatLng([latitude, longitude])
          map.setView([latitude, longitude], 16)
        }

        toast.success("Lokasi berhasil didapatkan")
      },
      (error) => {
        console.error("Geolocation error:", error)
        toast.error("Gagal mendapatkan lokasi. Pastikan GPS/lokasi diaktifkan.")
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 60000,
      }
    )
  } else {
    toast.error("Browser tidak mendukung geolocation")
  }
}

// Form validation
const validateForm = (): boolean => {
  errors.value = {}

  // Nomor KK validation
  if (!formData.value.nomor_kk) {
    errors.value.nomor_kk = "Nomor KK harus diisi"
  } else if (!/^\d{16}$/.test(formData.value.nomor_kk)) {
    errors.value.nomor_kk = "Nomor KK harus 16 digit angka"
  }

  // Nama Ayah validation
  if (!formData.value.nama_ayah) {
    errors.value.nama_ayah = "Nama ayah harus diisi"
  } else if (formData.value.nama_ayah.length < 2) {
    errors.value.nama_ayah = "Nama ayah minimal 2 karakter"
  }

  // Nama Ibu validation
  if (!formData.value.nama_ibu) {
    errors.value.nama_ibu = "Nama ibu harus diisi"
  } else if (formData.value.nama_ibu.length < 2) {
    errors.value.nama_ibu = "Nama ibu minimal 2 karakter"
  }

  // NIK Ayah validation
  if (!formData.value.nik_ayah) {
    errors.value.nik_ayah = "NIK ayah harus diisi"
  } else if (!/^\d{16}$/.test(formData.value.nik_ayah)) {
    errors.value.nik_ayah = "NIK ayah harus 16 digit angka"
  }

  // NIK Ibu validation
  if (!formData.value.nik_ibu) {
    errors.value.nik_ibu = "NIK ibu harus diisi"
  } else if (!/^\d{16}$/.test(formData.value.nik_ibu)) {
    errors.value.nik_ibu = "NIK ibu harus 16 digit angka"
  }

  // Alamat validation
  if (!formData.value.alamat) {
    errors.value.alamat = "Alamat harus diisi"
  } else if (formData.value.alamat.length < 5) {
    errors.value.alamat = "Alamat minimal 5 karakter"
  }

  // RT validation
  if (!formData.value.rt) {
    errors.value.rt = "RT harus diisi"
  } else if (!/^\d{1,3}$/.test(formData.value.rt)) {
    errors.value.rt = "RT harus berupa angka 1-3 digit"
  }

  // RW validation
  if (!formData.value.rw) {
    errors.value.rw = "RW harus diisi"
  } else if (!/^\d{1,3}$/.test(formData.value.rw)) {
    errors.value.rw = "RW harus berupa angka 1-3 digit"
  }

  // Kelurahan validation
  if (!formData.value.id_kelurahan) {
    errors.value.id_kelurahan = "Kelurahan harus dipilih"
  }

  // ✅ SIMPLIFIED: Coordinate validation
  if (!formData.value.koordinat || !Array.isArray(formData.value.koordinat)) {
    errors.value.koordinat = "Koordinat harus diset dengan mengklik peta"
  } else {
    const [lat, lng] = formData.value.koordinat
    
    // Only check for obviously invalid coordinates
    if ((lat === 0 && lng === 0) || Math.abs(lat) > 90 || Math.abs(lng) > 180) {
      errors.value.koordinat = "Koordinat tidak valid"
    } else if (props.mode === "create" && lat === -6.7064 && lng === 108.5492) {
      // Only reject default coords for CREATE mode
      errors.value.koordinat = "Silakan tentukan lokasi spesifik dengan mengklik peta"
    }
    // ✅ For EDIT mode, accept any valid coordinates (including default)
  }

  return Object.keys(errors.value).length === 0
}

// Reset form
const resetForm = () => {
  formData.value = {
    nomor_kk: "",
    nama_ayah: "",
    nama_ibu: "",
    nik_ayah: "",
    nik_ibu: "",
    alamat: "",
    rt: "",
    rw: "",
    id_kelurahan: "",
    koordinat: [-6.7064, 108.5492],
  }
  errors.value = {}
}

// Load form data for edit mode
const loadFormData = (keluarga: Keluarga) => {
  formData.value = {
    id: keluarga.id,
    nomor_kk: keluarga.nomor_kk,
    nama_ayah: keluarga.nama_ayah,
    nama_ibu: keluarga.nama_ibu,
    nik_ayah: keluarga.nik_ayah,
    nik_ibu: keluarga.nik_ibu,
    alamat: keluarga.alamat,
    rt: keluarga.rt,
    rw: keluarga.rw,
    id_kelurahan: keluarga.id_kelurahan,
    kelurahan: keluarga.kelurahan,
    kecamatan: keluarga.kecamatan,
    koordinat:
      keluarga.koordinat && Array.isArray(keluarga.koordinat)
        ? [...keluarga.koordinat]
        : [-6.7064, 108.5492],
    created_date: keluarga.created_date,
    updated_date: keluarga.updated_date,
  }
  errors.value = {}
}

// Handle save
const handleSave = () => {
  if (!validateForm()) {
    toast.error("Mohon periksa kembali form yang Anda isi")

    // Scroll to first error
    const firstErrorElement = document.querySelector(".border-red-500")
    if (firstErrorElement) {
      firstErrorElement.scrollIntoView({ behavior: "smooth", block: "center" })
    }

    return
  }

  // Auto-populate kelurahan dan kecamatan berdasarkan id_kelurahan
  const selectedKelurahan = kelurahanOptions.find((k) => k.id === formData.value.id_kelurahan)
  if (selectedKelurahan) {
    const [kelurahan, kecamatan] = selectedKelurahan.label.split(" - ")
    formData.value.kelurahan = kelurahan
    formData.value.kecamatan = kecamatan
  }

  emit("save", formData.value as Keluarga)
  resetForm()
}

// Handle close
const handleClose = () => {
  resetForm()
  destroyMap()
  emit("close")
}

// ✅ Watch for dialog visibility and mode changes
watch(
  () => props.show,
  async (newVal) => {
    console.log("👀 Dialog visibility changed:", newVal)

    if (newVal) {
      if (props.keluarga && props.mode === "edit") {
        loadFormData(props.keluarga)
      } else {
        resetForm()
      }

      // Wait for dialog animation to complete, then init map
      setTimeout(async () => {
        await initMap()
        // Force update map with loaded coordinates for edit mode
        if (props.mode === "edit" && formData.value.koordinat) {
          setTimeout(() => {
            updateMapWithCurrentCoordinates()
          }, 500)
        }
      }, 700)
    } else {
      destroyMap()
    }
  }
)

// ✅ Watch for keluarga prop changes
watch(
  () => props.keluarga,
  (newKeluarga) => {
    if (newKeluarga && props.mode === "edit" && props.show) {
      loadFormData(newKeluarga)

      // Update map if it's already initialized
      if (mapInitialized.value) {
        setTimeout(() => {
          updateMapWithCurrentCoordinates()
        }, 100)
      }
    }
  },
  { deep: true }
)

// ✅ Watch coordinate changes and update marker
watch(
  () => formData.value.koordinat,
  (newCoords) => {
    if (newCoords && Array.isArray(newCoords) && marker && map) {
      const [lat, lng] = newCoords
      if (lat !== 0 || lng !== 0) {
        marker.setLatLng([lat, lng])
        map.setView([lat, lng], map.getZoom())
      }
    }
  },
  { deep: true }
)

// Real-time validation
watch(
  () => formData.value.nomor_kk,
  (newVal) => {
    if (newVal && newVal.length === 16 && /^\d{16}$/.test(newVal)) {
      delete errors.value.nomor_kk
    }
  }
)

watch(
  () => formData.value.nik_ayah,
  (newVal) => {
    if (newVal && newVal.length === 16 && /^\d{16}$/.test(newVal)) {
      delete errors.value.nik_ayah
    }
  }
)

watch(
  () => formData.value.nik_ibu,
  (newVal) => {
    if (newVal && newVal.length === 16 && /^\d{16}$/.test(newVal)) {
      delete errors.value.nik_ibu
    }
  }
)

// ✅ Cleanup on unmount
onUnmounted(() => {
  destroyMap()
})
</script>

<template>
  <Dialog
    :open="show"
    @update:open="handleClose">
    <DialogContent class="w-[98vw] max-w-6xl max-h-[98vh] overflow-hidden p-0 map-dialog">
      <div class="flex flex-col h-full max-h-[98vh]">
        <!-- Header -->
        <div class="p-4 sm:p-6 border-b bg-background">
          <DialogHeader class="space-y-2">
            <DialogTitle class="text-lg sm:text-xl">
              {{ mode === "create" ? "Tambah Keluarga Baru" : "Edit Data Keluarga" }}
            </DialogTitle>
            <DialogDescription class="text-sm">
              {{
                mode === "create"
                  ? "Isi form di bawah untuk menambah data keluarga baru. Semua field dengan tanda (*) wajib diisi."
                  : "Perbarui informasi keluarga. Semua field dengan tanda (*) wajib diisi."
              }}
            </DialogDescription>
          </DialogHeader>
        </div>

        <!-- Content dengan scroll -->
        <div class="flex-1 overflow-y-auto p-4 sm:p-6">
          <div class="grid grid-cols-1 gap-6">
            <!-- Form Section -->
            <div class="space-y-6">
              <div class="space-y-4">
                <!-- Nomor KK -->
                <div class="space-y-2">
                  <Label
                    for="nomor_kk"
                    class="text-sm font-medium flex items-center gap-2">
                    Nomor KK *
                    <span class="text-xs text-muted-foreground">(16 digit angka)</span>
                  </Label>
                  <Input
                    id="nomor_kk"
                    v-model="formData.nomor_kk"
                    @input="(e: Event) => {
                      const target = e.target as HTMLInputElement
                      formData.nomor_kk = restrictToNumbers(target.value).substring(0, 16)
                    }"
                    placeholder="3209012345678901"
                    maxlength="16"
                    inputmode="numeric"
                    :class="errors.nomor_kk ? 'border-red-500' : ''"
                    class="font-mono" />
                  <p
                    v-if="errors.nomor_kk"
                    class="text-sm text-red-500">
                    {{ errors.nomor_kk }}
                  </p>
                </div>

                <!-- Nama Ayah & Ibu -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <Label
                      for="nama_ayah"
                      class="text-sm font-medium flex items-center gap-2">
                      Nama Ayah *
                      <span class="text-xs text-muted-foreground">(min. 2 karakter)</span>
                    </Label>
                    <Input
                      id="nama_ayah"
                      v-model="formData.nama_ayah"
                      placeholder="Budi Santoso"
                      :class="errors.nama_ayah ? 'border-red-500' : ''" />
                    <p
                      v-if="errors.nama_ayah"
                      class="text-sm text-red-500">
                      {{ errors.nama_ayah }}
                    </p>
                  </div>

                  <div class="space-y-2">
                    <Label
                      for="nama_ibu"
                      class="text-sm font-medium flex items-center gap-2">
                      Nama Ibu *
                      <span class="text-xs text-muted-foreground">(min. 2 karakter)</span>
                    </Label>
                    <Input
                      id="nama_ibu"
                      v-model="formData.nama_ibu"
                      placeholder="Siti Rahayu"
                      :class="errors.nama_ibu ? 'border-red-500' : ''" />
                    <p
                      v-if="errors.nama_ibu"
                      class="text-sm text-red-500">
                      {{ errors.nama_ibu }}
                    </p>
                  </div>
                </div>

                <!-- NIK Ayah & Ibu -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <Label
                      for="nik_ayah"
                      class="text-sm font-medium flex items-center gap-2">
                      NIK Ayah *
                      <span class="text-xs text-muted-foreground">(16 digit angka)</span>
                    </Label>
                    <Input
                      id="nik_ayah"
                      v-model="formData.nik_ayah"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.nik_ayah = restrictToNumbers(target.value).substring(0, 16)
                      }"
                      placeholder="3209012345678901"
                      maxlength="16"
                      inputmode="numeric"
                      :class="errors.nik_ayah ? 'border-red-500' : ''"
                      class="font-mono" />
                    <p
                      v-if="errors.nik_ayah"
                      class="text-sm text-red-500">
                      {{ errors.nik_ayah }}
                    </p>
                  </div>

                  <div class="space-y-2">
                    <Label
                      for="nik_ibu"
                      class="text-sm font-medium flex items-center gap-2">
                      NIK Ibu *
                      <span class="text-xs text-muted-foreground">(16 digit angka)</span>
                    </Label>
                    <Input
                      id="nik_ibu"
                      v-model="formData.nik_ibu"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.nik_ibu = restrictToNumbers(target.value).substring(0, 16)
                      }"
                      placeholder="3209012345678902"
                      maxlength="16"
                      inputmode="numeric"
                      :class="errors.nik_ibu ? 'border-red-500' : ''"
                      class="font-mono" />
                    <p
                      v-if="errors.nik_ibu"
                      class="text-sm text-red-500">
                      {{ errors.nik_ibu }}
                    </p>
                  </div>
                </div>

                <!-- Alamat -->
                <div class="space-y-2">
                  <Label
                    for="alamat"
                    class="text-sm font-medium flex items-center gap-2">
                    Alamat *
                    <span class="text-xs text-muted-foreground">(min. 5 karakter)</span>
                  </Label>
                  <Textarea
                    id="alamat"
                    v-model="formData.alamat"
                    placeholder="Jl. Kesambi Raya No. 123"
                    :class="errors.alamat ? 'border-red-500' : ''"
                    class="min-h-[80px] resize-none" />
                  <p
                    v-if="errors.alamat"
                    class="text-sm text-red-500">
                    {{ errors.alamat }}
                  </p>
                </div>

                <!-- RT/RW -->
                <div class="grid grid-cols-2 gap-4">
                  <div class="space-y-2">
                    <Label
                      for="rt"
                      class="text-sm font-medium flex items-center gap-2">
                      RT *
                      <span class="text-xs text-muted-foreground">(1-3 digit angka)</span>
                    </Label>
                    <Input
                      id="rt"
                      v-model="formData.rt"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.rt = restrictToNumbers3Digit(target.value)
                      }"
                      placeholder="001"
                      maxlength="3"
                      inputmode="numeric"
                      :class="errors.rt ? 'border-red-500' : ''"
                      class="font-mono" />
                    <p
                      v-if="errors.rt"
                      class="text-sm text-red-500">
                      {{ errors.rt }}
                    </p>
                  </div>
                  <div class="space-y-2">
                    <Label
                      for="rw"
                      class="text-sm font-medium flex items-center gap-2">
                      RW *
                      <span class="text-xs text-muted-foreground">(1-3 digit angka)</span>
                    </Label>
                    <Input
                      id="rw"
                      v-model="formData.rw"
                      @input="(e: Event) => {
                        const target = e.target as HTMLInputElement
                        formData.rw = restrictToNumbers3Digit(target.value)
                      }"
                      placeholder="002"
                      maxlength="3"
                      inputmode="numeric"
                      :class="errors.rw ? 'border-red-500' : ''"
                      class="font-mono" />
                    <p
                      v-if="errors.rw"
                      class="text-sm text-red-500">
                      {{ errors.rw }}
                    </p>
                  </div>
                </div>

                <!-- Kelurahan -->
                <div class="space-y-2">
                  <Label
                    for="kelurahan"
                    class="text-sm font-medium"
                    >Kelurahan *</Label
                  >
                  <Select v-model="formData.id_kelurahan">
                    <SelectTrigger
                      id="kelurahan"
                      :class="errors.id_kelurahan ? 'border-red-500' : ''">
                      <SelectValue placeholder="Pilih kelurahan" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem
                        v-for="kelurahan in kelurahanOptions"
                        :key="kelurahan.id"
                        :value="kelurahan.id">
                        {{ kelurahan.label }}
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  <p
                    v-if="errors.id_kelurahan"
                    class="text-sm text-red-500">
                    {{ errors.id_kelurahan }}
                  </p>
                </div>
              </div>
            </div>

            <!-- ✅ Map Section dengan Leaflet -->
            <div class="space-y-4">
              <Card
                class="h-fit"
                :class="errors.koordinat ? 'border-red-500' : ''">
                <CardHeader>
                  <CardTitle class="flex items-center gap-2 text-base">
                    <MapPin class="h-5 w-5" />
                    Lokasi Koordinat *
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      @click="getCurrentLocation"
                      class="ml-auto">
                      <Navigation class="h-4 w-4 mr-1" />
                      GPS
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent class="p-0">
                  <!-- ✅ Leaflet Map Container -->
                  <div
                    ref="mapContainer"
                    class="leaflet-map-container"
                    style="
                      height: 400px;
                      width: 100%;
                      position: relative;
                      z-index: 1;
                      display: block;
                      background-color: #f3f4f6;
                      border-radius: 0 0 0.5rem 0.5rem;
                    ">
                    <!-- Loading indicator -->
                    <div
                      v-if="!mapInitialized"
                      class="absolute inset-0 flex items-center justify-center bg-gray-50 rounded-b-lg z-10">
                      <div class="text-center">
                        <div
                          class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-2"></div>
                        <p class="text-sm text-muted-foreground">Memuat peta...</p>
                      </div>
                    </div>

                    <!-- Debug info -->
                    <div
                      v-if="mapInitialized"
                      class="absolute top-2 left-2 bg-black bg-opacity-70 text-white text-xs p-1 rounded z-20">
                      Map: ✅ | Mode: {{ mode }}
                    </div>
                  </div>
                </CardContent>
              </Card>

              <!-- ✅ Coordinate Display dengan styling yang sudah bekerja -->
              <div class="space-y-3">
                <h4 class="text-sm font-medium text-gray-700">Koordinat Lokasi *</h4>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div class="space-y-1">
                    <div class="text-xs text-gray-500 uppercase tracking-wide">Latitude</div>
                    <div class="px-3 py-2 bg-gray-50 border border-gray-200 rounded-md">
                      <span class="font-mono text-sm">
                        {{ formData.koordinat?.[0]?.toFixed(6) || "Belum diset" }}
                      </span>
                    </div>
                  </div>
                  <div class="space-y-1">
                    <div class="text-xs text-gray-500 uppercase tracking-wide">Longitude</div>
                    <div class="px-3 py-2 bg-gray-50 border border-gray-200 rounded-md">
                      <span class="font-mono text-sm">
                        {{ formData.koordinat?.[1]?.toFixed(6) || "Belum diset" }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Coordinate Error -->
              <p
                v-if="errors.koordinat"
                class="text-sm text-red-500">
                {{ errors.koordinat }}
              </p>

              <!-- ✅ Enhanced Info Box -->
              <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
                <p class="text-sm text-blue-800 flex items-start gap-2">
                  <MapPin class="h-4 w-4 mt-0.5 flex-shrink-0" />
                  <span>
                    <strong>Cara menggunakan peta:</strong><br />
                    • Klik pada peta untuk mengatur lokasi yang tepat<br />
                    • Atau drag marker merah untuk memindahkan posisi<br />
                    • Gunakan tombol GPS untuk lokasi otomatis<br />
                    • Koordinat akan otomatis terupdate dan wajib diisi<br />
                    • Pastikan lokasi sesuai dengan alamat yang dimasukkan
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="p-4 sm:p-6 border-t bg-background">
          <DialogFooter class="flex flex-col-reverse sm:flex-row gap-3">
            <Button
              variant="outline"
              @click="handleClose"
              class="w-full sm:w-auto">
              <X class="h-4 w-4 mr-2" />
              Batal
            </Button>
            <Button
              @click="handleSave"
              class="w-full sm:w-auto">
              <Save class="h-4 w-4 mr-2" />
              {{ mode === "create" ? "Simpan" : "Perbarui" }}
            </Button>
          </DialogFooter>
        </div>
      </div>
    </DialogContent>
  </Dialog>
</template>

<style scoped>
/* ✅ Leaflet CSS dengan higher specificity untuk dialog */
.map-dialog :deep(.leaflet-container) {
  height: 400px !important;
  width: 100% !important;
  z-index: 1 !important;
  position: relative !important;
  display: block !important;
  background-color: #f3f4f6 !important;
}

.leaflet-map-container {
  min-height: 400px !important;
  background: #f3f4f6 !important;
}

.map-dialog :deep(.leaflet-container) {
  border-radius: 0 0 0.5rem 0.5rem !important;
}

.map-dialog :deep(.leaflet-map-pane) {
  z-index: 1 !important;
}

.map-dialog :deep(.leaflet-tile-pane) {
  z-index: 1 !important;
}

.map-dialog :deep(.leaflet-overlay-pane) {
  z-index: 2 !important;
}

.map-dialog :deep(.leaflet-marker-pane) {
  z-index: 600 !important;
}

.map-dialog :deep(.leaflet-control-container) {
  z-index: 1000 !important;
}

.map-dialog :deep(.leaflet-tile) {
  opacity: 1 !important;
  visibility: visible !important;
}

.map-dialog :deep(.leaflet-control-attribution) {
  background: rgba(255, 255, 255, 0.8) !important;
  font-size: 10px !important;
}

/* Loading animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
