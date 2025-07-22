import { loadCirebonBoundary } from "./modules/geojson.js"

let map = L.map("map").setView([-6.726168577920489, 108.53918387877482], 14)

L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
  maxZoom: 19,
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
}).addTo(map)

async function getStuntingData() {
  try {
    const response = await fetch("/api/geojson/get")
    const responseJson = await response.json()
    const stuntingGeoJsonArray = responseJson.data

    const layers = {}

    stuntingGeoJsonArray.forEach((geojson, idx) => {
      const layer = L.geoJSON(geojson, {
        onEachFeature: function (feature, layer) {
          let props = feature.properties
          let tableRows = Object.entries(props)
            .map(([key, value]) => `<tr><th>${key}</th><td>${value}</td></tr>`)
            .join("")
          let popupContent = `<table class="table" border="1">${tableRows}</table>`
          layer.bindPopup(popupContent)
        },
      })

      const layerName = geojson.name || `Layer ${idx + 1}`
      layers[layerName] = layer
    })

    L.control.layers(null, layers, { collapsed: false }).addTo(map)

    // Add the first layer to the map by default
    const firstLayer = Object.values(layers)[0]
    if (firstLayer) firstLayer.addTo(map)
  } catch (error) {
    console.error("Error loading geojson:", error)
  }
}

getStuntingData()
loadCirebonBoundary(L, map)
