function printGeoJsonTable(geoJson, containerId) {
  if (!geoJson || !geoJson.features || geoJson.features.length === 0) {
    document.getElementById(containerId).innerHTML = "<p>No data available</p>"
    return
  }

  // Take all unique property keys from the features
  const allProps = new Set()
  geoJson.features.forEach((feature) => {
    Object.keys(feature.properties).forEach((key) => allProps.add(key))
  })
  const headers = Array.from(allProps)
  const headerRow = headers.map((key) => `<th>${key}</th>`).join("")
  const dataRows = geoJson.features
    .map((feature) => {
      return (
        `<tr>` +
        headers.map((key) => `<td>${feature.properties[key] ?? ""}</td>`).join("") +
        `</tr>`
      )
    })
    .join("")

  let tableContent = `
    <table class="table" border="1" cellpadding="4" cellspacing="0">
      <thead><tr>${headerRow}</tr></thead>
      <tbody>${dataRows}</tbody>
    </table>
  `
  document.getElementById(containerId).innerHTML = tableContent
}

function printGeoJsonSelectedFeature(feature, containerId) {
  if (!feature || !feature.properties) {
    document.getElementById(containerId).innerHTML = "<p>No feature selected</p>"
    return
  }

  let props = feature.properties
  let tableRows = Object.entries(props)
    .map(([key, value]) => `<tr><th>${key}</th><td>${value}</td></tr>`)
    .join("")
  let tableContent = `<table class="table" border="1" cellpadding="4" cellspacing="0">${tableRows}</table>`

  document.getElementById(containerId).innerHTML = tableContent
}

async function loadCirebonBoundary(L, map) {
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

export { printGeoJsonTable, printGeoJsonSelectedFeature, loadCirebonBoundary }
