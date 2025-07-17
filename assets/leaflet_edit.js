import {
  loadCirebonBoundary,
  printGeoJsonSelectedFeature,
  printGeoJsonTable,
} from "./modules/geojson.js"

let map = L.map("map").setView([-6.726168577920489, 108.53918387877482], 14)
let stuntingGeoJson

L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map)

async function getStuntingData() {
  try {
    const response = await fetch("static/geojson.geojson")
    stuntingGeoJson = await response.json()

    L.geoJSON(stuntingGeoJson, {
      onEachFeature: function (feature, layer) {
        layer.on("click", function () {
          printGeoJsonSelectedFeature(feature, "selected-feature")
        })
      },
    }).addTo(map)

    printGeoJsonTable(stuntingGeoJson, "geojson-table")
  } catch (error) {
    console.error("Error loading geojson:", error)
  }
}

getStuntingData()
loadCirebonBoundary(L, map)

map.pm.addControls({
  position: "topleft",
  drawMarker: true,
  drawCircleMarker: false,
  drawPolyline: false,
  drawRectangle: false,
  drawPolygon: false,
  drawCircle: false,
  drawText: false,
  cutPolygon: false,
  rotateMode: false,
  editMode: false,
})

// MARK: Create Marker
// Handle the creation of new markers
map.on("pm:create", function (e) {
  if (e.shape === "Marker") {
    const marker = e.layer
    const latLng = marker.getLatLng()
    const popupContent = `<table class="table" border="1">
      <tr><th>Latitude</th><td>${latLng.lat}</td></tr>
      <tr><th>Longitude</th><td>${latLng.lng}</td></tr>
    </table>`

    marker.bindPopup(popupContent).openPopup()

    // Add the new feature to the geoJson object
    if (stuntingGeoJson && stuntingGeoJson.type === "FeatureCollection") {
      const newFeature = {
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [latLng.lng, latLng.lat],
        },
        properties: {
          RW: "-",
          Kampung: "-",
          Kelurahan: "-",
          Kedalaman: "-",
          Durasi: "-",
          Dampak: "-",
          Penyebab: "-",
          Kerugian: "-",
          Tahun: "-",
          Sumber: "-",
          Foto: "-",
        },
      }
      stuntingGeoJson.features.push(newFeature)
      printGeoJsonTable(stuntingGeoJson, "geojson-table")
    }
  }
})

// MARK: Remove Marker
// Handle the deletion of markers
map.on("pm:remove", function (e) {
  if (e.layer instanceof L.Marker) {
    const marker = e.layer
    const latLng = marker.getLatLng()

    // Find and remove the feature from the geoJson object
    if (stuntingGeoJson && stuntingGeoJson.type === "FeatureCollection") {
      stuntingGeoJson.features = stuntingGeoJson.features.filter(
        (feature) =>
          feature.geometry.type !== "Point" ||
          feature.geometry.coordinates[0] !== latLng.lng ||
          feature.geometry.coordinates[1] !== latLng.lat
      )
      printGeoJsonTable(stuntingGeoJson, "geojson-table")
    }
  }
})
