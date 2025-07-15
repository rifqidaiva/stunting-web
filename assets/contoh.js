var map = L.map("map", {
  fullscreenControl: true,
}).setView([-6.743458, 108.5550681], 13)

// ------------------------- Map Tiles ------------------------- //
// Streets
var streets = L.tileLayer(
  "http://1.base.maps.cit.api.here.com/maptile/2.1/maptile/newest/{style}/{z}/{x}/{y}/256/png8?app_id={app_id}&app_code={app_code}&lg=eng",
  {
    attribution: 'Imagery &copy; <a href="https://www.here.com">HERE</a>',
    style: "normal.day",
    app_id: "eAdkWGYRoc4RfxVo0Z4B",
    app_code: "TrLJuXVK62IQk0vuXFzaig",
  }
).addTo(map)
// Hybrid
var hybrid = L.tileLayer(
  "http://1.aerial.maps.cit.api.here.com/maptile/2.1/maptile/newest/{style}/{z}/{x}/{y}/256/png8?app_id={app_id}&app_code={app_code}&lg=eng",
  {
    attribution: 'Imagery &copy; <a href="https://www.here.com">HERE</a>',
    style: "hybrid.day",
    app_id: "eAdkWGYRoc4RfxVo0Z4B",
    app_code: "TrLJuXVK62IQk0vuXFzaig",
  }
)

// ------------------------- Feature ------------------------- //
// Feature - Batas
function feature_batas(feature, layer) {
  var out = []
  var header = ""
  var body = ""
  if (feature.properties) {
    console.log(feature.properties["Kecamatan"])
    if (feature.properties["Kelurahan"]) {
      header = "Kelurahan " + feature.properties["Kelurahan"]
    } else if (feature.properties["Kecamatan"]) {
      header = "Kecamatan " + feature.properties["Kecamatan"]
    } else {
      header = feature.properties["Kabupaten/Kota"]
    }
    body += '<table class="table table-sm table-striped">'
    for (key in feature.properties) {
      body += "<tr><td>" + key + "</td><td>:</td><td>" + feature.properties[key] + "</td></tr>"
    }
    body += "</table>"
    out.push(
      '<div class="card"><div class="card-header text-center fw-bold">' +
        header +
        '</div><div class="card"><div class="card-body">' +
        body +
        "</div></div>"
    )
    layer.bindTooltip(header).bindPopup(out.join(), {
      maxWidth: 400,
    })
  }
}
// Feature - Banjir
function feature_banjir(feature, layer) {
  var out = []
  var body = ""
  if (feature.properties) {
    // console.log(feature.properties);
    if (feature.properties["Foto"]) {
      body =
        '<img src="assets/images/bpbd/banjir/' +
        feature.properties["Foto"] +
        '" class="img-fluid" alt="Responsive image"><br><br>'
    }
    body +=
      "<strong>Lokasi:</strong><br>" +
      feature.properties["Kampung"] +
      ", " +
      feature.properties["RW"] +
      " " +
      feature.properties["Kelurahan"] +
      "<br><br> <strong>Keterangan:</strong>" +
      '<table class="table table-sm table-striped">'
    if (feature.properties["Dampak"] != "-") {
      body += "<tr><td>Dampak</td><td>:</td><td>" + feature.properties["Dampak"] + "</td></tr>"
    }
    if (feature.properties["Durasi"] != "-") {
      body += "<tr><td>Durasi</td><td>:</td><td>" + feature.properties["Durasi"] + "</td></tr>"
    }
    if (feature.properties["Kedalaman"] != "-") {
      body +=
        "<tr><td>Kedalaman</td><td>:</td><td>" + feature.properties["Kedalaman"] + "</td></tr>"
    }
    if (feature.properties["Kerugian"] != "-") {
      body += "<tr><td>Kerugian</td><td>:</td><td>" + feature.properties["Kerugian"] + "</td></tr>"
    }
    if (feature.properties["Penyebab"] != "-") {
      body += "<tr><td>Penyebab</td><td>:</td><td>" + feature.properties["Penyebab"] + "</td></tr>"
    }
    body +=
      "<tr><td>Sumber</td><td>:</td><td>" +
      feature.properties["Sumber"] +
      "</td></tr>" +
      "<tr><td>Tahun</td><td>:</td><td>" +
      feature.properties["Tahun"] +
      "</td></tr>"
    body += "</table>"
    out.push(
      '<div class="card"><div class="card-header text-center fw-bold">' +
        feature.properties["Kampung"] +
        '</div><div class="card-body">' +
        body +
        "</div></div>"
    )
    layer.bindTooltip(feature.properties["Kampung"]).bindPopup(out.join(), {
      maxWidth: 350,
    })
  }
}
// Feature - Titik Evakuasi
function feature_titik_evakuasi(feature, layer) {
  var out = []
  var body = ""
  if (feature.properties) {
    // console.log(feature.properties);
    body =
      "<strong>Lokasi:</strong><br>" +
      feature.properties["Alamat"] +
      ", " +
      feature.properties["Kelurahan"] +
      "<br><br> <strong>Keterangan:</strong>" +
      '<table class="table table-sm table-striped">' +
      "<tr><td>Luas</td><td>:</td><td>" +
      feature.properties["Luas"] +
      "</td></tr>" +
      "<tr><td>Akses Air Bersih</td><td>:</td><td>" +
      feature.properties["Akses Air Bersih"] +
      "</td></tr>" +
      "<tr><td>Jenis</td><td>:</td><td>" +
      feature.properties["Jenis"] +
      "</td></tr>" +
      "<tr><td>Keterangan</td><td>:</td><td>" +
      feature.properties["Keterangan"] +
      "</td></tr>" +
      "</table>"
    out.push(
      '<div class="card"><div class="card-header text-center fw-bold">' +
        feature.properties["Titik Evakuasi"] +
        '</div><div class="card-body">' +
        body +
        "</div></div>"
    )
    layer.bindTooltip(feature.properties["Titik Evakuasi"]).bindPopup(out.join(), {
      maxWidth: 350,
    })
  }
}

// ------------------------- Style ------------------------- //
// Style - Batas Administrasi
function style_batas() {
  return {
    color: "#656565",
    fillOpacity: 0,
  }
}

// ------------------------- Icon ------------------------- //
// Icon - Banjir
var banjir_icon = L.icon({
  iconUrl: "assets/markers/bpbd/banjir_icon.png",
  iconSize: [29, 43], // Icon Size
  iconAnchor: [14, 40], // Icon Anchor
  popupAnchor: [1, -37], // Pop Up Anchor
  tooltipAnchor: [15, -25], // Tooltip Anchor
})
// Icon - Titik Evakuasi
var titik_evakuasi_icon = L.icon({
  iconUrl: "assets/markers/bpbd/titik_evakuasi_icon.png",
  iconSize: [29, 43], // Icon Size
  iconAnchor: [14, 40], // Icon Anchor
  popupAnchor: [1, -37], // Pop Up Anchor
  tooltipAnchor: [15, -25], // Tooltip Anchor
})

// ------------------------- GeoJSON ------------------------- //
// Batas Kota
var batas_kot = new L.GeoJSON.AJAX(["geojson/bpbd/kota_batas.geojson"], {
  style: style_batas,
  onEachFeature: feature_batas,
}).addTo(map)
// Batas Kecamatan
var batas_kec = new L.GeoJSON.AJAX(["geojson/bpbd/kecamatan_batas.geojson"], {
  style: style_batas,
  onEachFeature: feature_batas,
})
// Batas Kelurahan
var batas_kel = new L.GeoJSON.AJAX(["geojson/bpbd/kelurahan_batas.geojson"], {
  style: style_batas,
  onEachFeature: feature_batas,
})
// Sebaran Banjir
var banjir_sebaran = new L.GeoJSON.AJAX(["geojson/bpbd/banjir_sebaran.geojson"], {
  pointToLayer: function (geoJsonPoint, latlng) {
    return L.marker(latlng, { icon: banjir_icon })
  },
  onEachFeature: feature_banjir,
}).addTo(map)
// Sebaran Banjir 2023
var banjir_sebaran_2023 = new L.GeoJSON.AJAX(["geojson/bpbd/banjir_sebaran_2023.geojson"], {
  pointToLayer: function (geoJsonPoint, latlng) {
    return L.marker(latlng, { icon: banjir_icon })
  },
  onEachFeature: feature_banjir,
}).addTo(map)
// Sebaran Titik Evakuasi
var titik_evakuasi_sebaran = new L.GeoJSON.AJAX(["geojson/bpbd/titik_evakuasi_sebaran.geojson"], {
  pointToLayer: function (geoJsonPoint, latlng) {
    return L.marker(latlng, { icon: titik_evakuasi_icon })
  },
  onEachFeature: feature_titik_evakuasi,
}).addTo(map)

// ------------------------- Layer Control ------------------------- //
// Maps
var baseLayers = {
  Streets: streets,
  Hybrid: hybrid,
}
// Layers
var overlays = {
  "Batas Kota": batas_kot,
  "Batas Kecamatan": batas_kec,
  "Batas Kelurahan": batas_kel,
  "Sebaran Banjir 2017 - 2022": banjir_sebaran,
  "Sebaran Banjir 2023": banjir_sebaran_2023,
  "Sebaran Titik Evakuasi": titik_evakuasi_sebaran,
}

L.control.layers(baseLayers, overlays).addTo(map)
