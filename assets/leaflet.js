import { loadCirebonBoundary } from "./modules/geojson.js"

let map = L.map("map").setView([-6.726168577920489, 108.53918387877482], 14)

L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map)

async function getStuntingData() {
  try {
    const response = await fetch("static/geojson.geojson")
    const stuntingGeoJson = await response.json()

    L.geoJSON(stuntingGeoJson, {
      onEachFeature: function (feature, layer) {
        let props = feature.properties
        let tableRows = Object.entries(props)
          .map(([key, value]) => `<tr><th>${key}</th><td>${value}</td></tr>`)
          .join("")
        let popupContent = `<table class="table" border="1">${tableRows}</table>`
        layer.bindPopup(popupContent)
      },
    }).addTo(map)
  } catch (error) {
    console.error("Error loading geojson:", error)
  }
}

getStuntingData()
loadCirebonBoundary(L, map)
