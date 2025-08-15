import type { ColumnDef } from "@tanstack/vue-table"
import { ArrowUpDown, MoreHorizontal, Pencil, Trash2 } from "lucide-vue-next"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { h } from "vue"

export interface LaporanMasyarakat {
  id: string
  id_masyarakat: string
  id_balita: string
  id_status_laporan: string
  tanggal_laporan: string
  hubungan_dengan_balita: string
  nomor_hp_pelapor: string
  nomor_hp_keluarga_balita: string
  created_date: string
  updated_date?: string

  // Extended fields from JOIN
  nama_pelapor: string
  email_pelapor: string
  nama_balita: string
  nama_ayah: string
  nama_ibu: string
  nomor_kk: string
  alamat: string
  kelurahan: string
  kecamatan: string
  status_laporan: string
  jenis_laporan: string
}

// Status badge variant mapping
const getStatusVariant = (status: string): "default" | "secondary" | "destructive" | "outline" => {
  switch (status.toLowerCase()) {
    case "belum diproses":
      return "outline"
    case "diproses dan data tidak sesuai":
      return "destructive"
    case "diproses dan data sesuai":
      return "secondary"
    case "belum ditindaklanjuti":
      return "outline"
    case "sudah ditindaklanjuti":
      return "default"
    case "sudah perbaikan gizi":
      return "secondary"
    default:
      return "outline"
  }
}

export const columns: ColumnDef<LaporanMasyarakat>[] = [
  {
    id: "tanggal_laporan",
    accessorKey: "tanggal_laporan",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Tanggal Laporan", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const date = row.getValue("tanggal_laporan") as string
      return h("div", { class: "text-sm" }, new Date(date).toLocaleDateString("id-ID"))
    },
    meta: {
      displayName: "Tanggal Laporan",
    },
  },
  {
    id: "nama_pelapor",
    accessorKey: "nama_pelapor",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Pelapor", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const namaPelapor = row.getValue("nama_pelapor") as string
      const nomorHp = row.original.nomor_hp_pelapor
      const hubungan = row.original.hubungan_dengan_balita

      return h("div", { class: "space-y-1" }, [
        h("div", { class: "font-medium text-sm" }, namaPelapor),
        h("div", { class: "text-xs text-muted-foreground" }, nomorHp),
        h("div", { class: "text-xs text-blue-600" }, `(${hubungan})`),
      ])
    },
    meta: {
      displayName: "Data Pelapor",
    },
  },
  {
    id: "nama_balita",
    accessorKey: "nama_balita",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Balita", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const namaBalita = row.getValue("nama_balita") as string
      const namaAyah = row.original.nama_ayah
      const namaIbu = row.original.nama_ibu

      return h("div", { class: "space-y-1" }, [
        h("div", { class: "font-medium text-sm" }, namaBalita),
        h("div", { class: "text-xs text-muted-foreground" }, `Ayah: ${namaAyah}`),
        h("div", { class: "text-xs text-muted-foreground" }, `Ibu: ${namaIbu}`),
      ])
    },
    meta: {
      displayName: "Data Balita",
    },
  },
  {
    id: "keluarga_info",
    // Hapus accessorKey karena kita akan menggunakan custom cell
    header: "Info Keluarga",
    cell: ({ row }) => {
      const nomorKk = row.original.nomor_kk // Gunakan row.original
      const alamat = row.original.alamat
      const nomorHpKeluarga = row.original.nomor_hp_keluarga_balita

      return h("div", { class: "space-y-1 max-w-[200px]" }, [
        h("div", { class: "font-mono text-xs" }, nomorKk),
        h("div", { class: "text-xs text-muted-foreground truncate" }, alamat),
        h("div", { class: "text-xs text-blue-600" }, nomorHpKeluarga),
      ])
    },
    // Tambahkan filterFn untuk pencarian
    filterFn: (row, value) => {
      const nomorKk = row.original.nomor_kk?.toLowerCase() || ""
      const alamat = row.original.alamat?.toLowerCase() || ""
      return nomorKk.includes(value.toLowerCase()) || alamat.includes(value.toLowerCase())
    },
    meta: {
      displayName: "Info Keluarga",
    },
  },
  {
    id: "wilayah",
    // Hapus accessorKey karena kita akan menggunakan custom cell
    header: "Wilayah",
    cell: ({ row }) => {
      const kelurahan = row.original.kelurahan // Gunakan row.original
      const kecamatan = row.original.kecamatan

      return h("div", { class: "text-sm" }, [
        h("div", { class: "font-medium" }, kelurahan),
        h("div", { class: "text-xs text-muted-foreground" }, kecamatan),
      ])
    },
    // Tambahkan filterFn untuk pencarian
    filterFn: (row, value) => {
      const kelurahan = row.original.kelurahan?.toLowerCase() || ""
      const kecamatan = row.original.kecamatan?.toLowerCase() || ""
      return kelurahan.includes(value.toLowerCase()) || kecamatan.includes(value.toLowerCase())
    },
    meta: {
      displayName: "Wilayah",
    },
  },
  {
    id: "jenis_laporan",
    accessorKey: "jenis_laporan",
    header: "Jenis",
    cell: ({ row }) => {
      const jenisLaporan = row.getValue("jenis_laporan") as string

      return h(
        Badge,
        {
          variant: "outline",
          class: "text-xs",
        },
        () => jenisLaporan
      )
    },
    meta: {
      displayName: "Jenis Laporan",
    },
  },
  {
    id: "status_laporan",
    accessorKey: "status_laporan",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Status", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const statusLaporan = row.getValue("status_laporan") as string

      return h(
        Badge,
        {
          variant: getStatusVariant(statusLaporan),
          class: "text-xs whitespace-nowrap",
        },
        () => statusLaporan
      )
    },
    meta: {
      displayName: "Status Laporan",
    },
  },
  {
    id: "created_date",
    accessorKey: "created_date",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Dibuat", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const date = row.getValue("created_date") as string
      const updatedDate = row.original.updated_date

      return h("div", { class: "text-sm space-y-1" }, [
        h("div", {}, new Date(date).toLocaleDateString("id-ID")),
        updatedDate &&
          h(
            "div",
            { class: "text-xs text-muted-foreground" },
            `Diupdate: ${new Date(updatedDate).toLocaleDateString("id-ID")}`
          ),
      ])
    },
    meta: {
      displayName: "Tanggal Dibuat",
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const laporan = row.original

      return h(
        "div",
        { class: "relative" },
        h(
          DropdownMenu,
          {},
          {
            default: () => [
              h(
                DropdownMenuTrigger,
                {},
                {
                  default: () =>
                    h(
                      Button,
                      { variant: "ghost", class: "h-8 w-8 p-0" },
                      {
                        default: () => [
                          h("span", { class: "sr-only" }, "Open menu"),
                          h(MoreHorizontal, { class: "h-4 w-4" }),
                        ],
                      }
                    ),
                }
              ),
              h(
                DropdownMenuContent,
                { align: "end" },
                {
                  default: () => [
                    h(DropdownMenuLabel, {}, { default: () => "Actions" }),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(laporan.id)
                        },
                      },
                      { default: () => "Copy ID" }
                    ),
                    h(DropdownMenuSeparator),
                    // h(
                    //   DropdownMenuItem,
                    //   {
                    //     onClick: () => {
                    //       document.dispatchEvent(
                    //         new CustomEvent("view-laporan", { detail: laporan })
                    //       )
                    //     },
                    //   },
                    //   { default: () => [h(Eye, { class: "mr-2 h-4 w-4" }), "Lihat Detail"] }
                    // ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("edit-laporan", { detail: laporan })
                          )
                        },
                      },
                      { default: () => [h(Pencil, { class: "mr-2 h-4 w-4" }), "Edit"] }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("delete-laporan", { detail: laporan })
                          )
                        },
                        class: "text-red-600",
                      },
                      { default: () => [h(Trash2, { class: "mr-2 h-4 w-4" }), "Delete"] }
                    ),
                  ],
                }
              ),
            ],
          }
        )
      )
    },
  },
]
