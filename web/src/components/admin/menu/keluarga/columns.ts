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

export interface Keluarga {
  id: string
  nomor_kk: string
  nama_ayah: string
  nama_ibu: string
  nik_ayah: string
  nik_ibu: string
  alamat: string
  rt: string
  rw: string
  id_kelurahan: string
  kelurahan: string
  kecamatan: string
  koordinat: [number, number]
  created_date: string
  updated_date?: string
}

export const columns: ColumnDef<Keluarga>[] = [
  {
    id: "nomor_kk",
    accessorKey: "nomor_kk",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Nomor KK", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const nomorKk = row.getValue("nomor_kk") as string
      return h("div", { class: "font-mono text-sm" }, nomorKk)
    },
    meta: {
      displayName: "Nomor KK",
    },
  },
  {
    id: "nama_ayah",
    accessorKey: "nama_ayah",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Nama Ayah", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const namaAyah = row.getValue("nama_ayah") as string
      return h("div", { class: "font-medium" }, namaAyah)
    },
    meta: {
      displayName: "Nama Ayah",
    },
  },
  {
    id: "nama_ibu",
    accessorKey: "nama_ibu",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Nama Ibu", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const namaIbu = row.getValue("nama_ibu") as string
      return h("div", { class: "font-medium" }, namaIbu)
    },
    meta: {
      displayName: "Nama Ibu",
    },
  },
  {
    id: "alamat",
    accessorKey: "alamat",
    header: "Alamat",
    cell: ({ row }) => {
      const alamat = row.getValue("alamat") as string
      const rt = row.original.rt
      const rw = row.original.rw
      return h("div", { class: "max-w-[200px]" }, [
        h("div", { class: "truncate text-sm" }, alamat),
        h("div", { class: "text-xs text-muted-foreground" }, `RT ${rt}/RW ${rw}`),
      ])
    },
    meta: {
      displayName: "Alamat",
    },
  },
  {
    id: "kelurahan",
    accessorKey: "kelurahan",
    header: "Wilayah",
    cell: ({ row }) => {
      const kelurahan = row.getValue("kelurahan") as string
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
    id: "koordinat",
    accessorKey: "koordinat",
    header: "Koordinat",
    cell: ({ row }) => {
      const koordinat = row.getValue("koordinat") as [number, number]
      return h("div", { class: "font-mono text-xs" }, [
        h("div", {}, `${koordinat[1].toFixed(6)}`), // Longitude
        h("div", {}, `${koordinat[0].toFixed(6)}`), // Latitude
      ])
    },
    meta: {
      displayName: "Koordinat",
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
      const keluarga = row.original

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
                          navigator.clipboard.writeText(keluarga.id)
                        },
                      },
                      { default: () => "Copy ID" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("edit-keluarga", { detail: keluarga })
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
                            new CustomEvent("delete-keluarga", { detail: keluarga })
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
