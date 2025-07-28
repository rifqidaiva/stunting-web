// Data example

// [
//     {
//         "id": "3a11a99c-66af-11f0-a701-2811a8eb1247",
//         "name": "Intan Kedua",
//         "nik": "1234567890123457",
//         "date_of_birth": "2018-05-21",
//         "coordinates": [
//             110.124,
//             -7.123
//         ],
//         "status": "Belum diproses",
//         "reported_by_id": "b0429037-66ad-11f0-a701-2811a8eb1247"
//     },
//     {
//         "id": "a4640e72-66bc-11f0-a701-2811a8eb1247",
//         "name": "Budi",
//         "nik": "1234567890123450",
//         "date_of_birth": "2018-05-21",
//         "coordinates": [
//             110.125,
//             -7.1245
//         ],
//         "status": "Belum diproses",
//         "reported_by_id": "b0429037-66ad-11f0-a701-2811a8eb1247"
//     },
//     {
//         "id": "ba24d2bf-66ae-11f0-a701-2811a8eb1247",
//         "name": "Budi",
//         "nik": "1234567890123456",
//         "date_of_birth": "2018-05-21",
//         "coordinates": [
//             110.123,
//             -7.123
//         ],
//         "status": "Belum diproses",
//         "reported_by_id": "b0429037-66ad-11f0-a701-2811a8eb1247"
//     }
// ]

import { h } from "vue"
import type { ColumnDef } from "@tanstack/vue-table"

interface Sufferer {
  id: string
  name: string
  nik: string
  date_of_birth: string
  coordinates: [number, number]
  status:
    | "Belum diproses"
    | "Diproses dan data tidak sesuai"
    | "Diproses dan data sesuai"
    | "Belum ditindaklanjuti"
    | "Sudah ditindaklanjuti"
    | "Sudah perbaikan gizi"
  reported_by_id: string
}

const columns: ColumnDef<Sufferer>[] = [
  {
    accessorKey: "name",
    header: () => h("span", "Name"),
  },
  {
    accessorKey: "nik",
    header: () => h("span", "NIK"),
  },
  {
    accessorKey: "date_of_birth",
    header: () => h("span", "Date of Birth"),
  },
  {
    accessorKey: "coordinates",
    header: () => h("span", "Coordinates"),
  },
  {
    accessorKey: "status",
    header: () => h("span", "Status"),
  },
  {
    accessorKey: "reported_by_id",
    header: () => h("span", "Reported By"),
  },
]

export default columns
