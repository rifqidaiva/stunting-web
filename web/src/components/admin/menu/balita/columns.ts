import type { ColumnDef } from "@tanstack/vue-table"
import { ArrowUpDown, MoreHorizontal, Pencil, Trash2 } from "lucide-vue-next"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { h } from "vue"

export interface Balita {
  id: string
  id_keluarga: string
  nomor_kk: string
  nama_ayah: string
  nama_ibu: string
  nama: string
  tanggal_lahir: string
  jenis_kelamin: "L" | "P"
  berat_lahir: string
  tinggi_lahir: string
  umur: string
  kelurahan: string
  kecamatan: string
  created_date: string
  updated_date?: string
}

export const columns: ColumnDef<Balita>[] = [
  {
    id: "nama",
    accessorKey: "nama",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Nama Balita", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const nama = row.getValue("nama") as string
      const jenisKelamin = row.original.jenis_kelamin
      return h("div", [
        h("div", { class: "font-medium" }, nama),
        h(
          "div",
          { class: "text-xs text-muted-foreground" },
          jenisKelamin === "L" ? "Laki-laki" : "Perempuan"
        ),
      ])
    },
    meta: {
      displayName: "Nama Balita",
    },
  },
  {
    id: "tanggal_lahir",
    accessorKey: "tanggal_lahir",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Tanggal Lahir", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const tanggalLahir = row.getValue("tanggal_lahir") as string
      const umur = row.original.umur
      return h("div", [
        h("div", { class: "text-sm" }, new Date(tanggalLahir).toLocaleDateString("id-ID")),
        h("div", { class: "text-xs text-muted-foreground" }, umur),
      ])
    },
    meta: {
      displayName: "Tanggal Lahir",
    },
  },
  {
    id: "keluarga",
    // Tidak menggunakan accessorKey karena ini adalah kolom virtual yang menggabungkan beberapa field
    header: "Orang Tua",
    cell: ({ row }) => {
      const namaAyah = row.original.nama_ayah
      const namaIbu = row.original.nama_ibu
      const nomorKk = row.original.nomor_kk
      return h("div", { class: "max-w-[200px]" }, [
        h("div", { class: "text-sm font-medium" }, `Ayah: ${namaAyah}`),
        h("div", { class: "text-sm font-medium" }, `Ibu: ${namaIbu}`),
        h("div", { class: "text-xs text-muted-foreground font-mono" }, nomorKk),
      ])
    },
    meta: {
      displayName: "Orang Tua",
    },
    // Tambahkan accessor function untuk sorting/filtering
    accessorFn: (row) => `${row.nama_ayah} ${row.nama_ibu} ${row.nomor_kk}`,
  },
  {
    id: "berat_tinggi_lahir",
    // Tidak menggunakan accessorKey karena ini kolom virtual
    header: "Data Lahir",
    cell: ({ row }) => {
      const beratLahir = row.original.berat_lahir
      const tinggiLahir = row.original.tinggi_lahir
      return h("div", [
        h("div", { class: "text-sm" }, `${beratLahir} gram`),
        h("div", { class: "text-sm" }, `${tinggiLahir} cm`),
      ])
    },
    meta: {
      displayName: "Data Lahir",
    },
    // Tambahkan accessor function untuk sorting/filtering
    accessorFn: (row) => `${row.berat_lahir} ${row.tinggi_lahir}`,
  },
  {
    id: "wilayah",
    accessorKey: "kelurahan",
    header: "Wilayah",
    cell: ({ row }) => {
      const kelurahan = row.getValue("wilayah") as string
      const kecamatan = row.original.kecamatan
      return h("div", { class: "text-sm" }, [
        h("div", { class: "font-medium" }, kelurahan),
        h("div", { class: "text-xs text-muted-foreground" }, kecamatan),
      ])
    },
    meta: {
      displayName: "Wilayah",
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
        () => ["Tanggal Dibuat", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const date = row.getValue("created_date") as string
      return h("div", { class: "text-sm" }, new Date(date).toLocaleDateString("id-ID"))
    },
    meta: {
      displayName: "Tanggal Dibuat",
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const balita = row.original

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
                          navigator.clipboard.writeText(balita.id)
                        },
                      },
                      { default: () => "Copy ID" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(new CustomEvent("edit-balita", { detail: balita }))
                        },
                      },
                      { default: () => [h(Pencil, { class: "mr-2 h-4 w-4" }), "Edit"] }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("delete-balita", { detail: balita })
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
