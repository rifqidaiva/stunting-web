import { printGeoJsonSelectedFeature, printGeoJsonTable } from "./modules/geojson.js"

var map = L.map("map").setView([-6.726168577920489, 108.53918387877482], 14)
let geoJson

L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map)

async function loadStunting() {
  try {
    const response = await fetch("static/geojson.geojson")
    geoJson = await response.json()

    L.geoJSON(geoJson, {
      onEachFeature: function (feature, layer) {
        let props = feature.properties
        let tableRows = Object.entries(props)
          .map(([key, value]) => `<tr><th>${key}</th><td>${value}</td></tr>`)
          .join("")
        let popupContent = `<table class="table" border="1" cellpadding="4" cellspacing="0">${tableRows}</table>`
        layer.bindPopup(popupContent)
        layer.on("click", function () {
          printGeoJsonSelectedFeature(feature, "selected-feature")
        })
      },
    }).addTo(map)
    printGeoJsonTable(geoJson, "geojson-table")
  } catch (error) {
    console.error("Error loading geojson:", error)
  }
}

async function loadCirebonBoundary() {
  try {
    const response = await fetch("static/cirebon_boundary.geojson")
    const data = await response.json()

    L.geoJSON(data, {
      style: {
        color: "red",
        weight: 3,
        fill: false,
      },
    }).addTo(map)
  } catch (error) {
    console.error("Error loading Cirebon boundary:", error)
  }
}

loadStunting()
loadCirebonBoundary()

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

map.on("pm:create", function (e) {
  if (e.shape === "Marker") {
    const marker = e.layer
    const latLng = marker.getLatLng()
    const popupContent = `<table class="table" border="1" cellpadding="4" cellspacing="0">
      <tr><th>Latitude</th><td>${latLng.lat}</td></tr>
      <tr><th>Longitude</th><td>${latLng.lng}</td></tr>
    </table>`
    marker.bindPopup(popupContent).openPopup()

    // Add the new feature to the geoJson object
    if (geoJson && geoJson.type === "FeatureCollection") {
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
      geoJson.features.push(newFeature)
      printGeoJsonTable(geoJson, "geojson-table")
    }
  }
})
