import type { ColumnDef } from "@tanstack/vue-table"
import { ArrowUpDown, MoreHorizontal, Pencil, Trash2, Users } from "lucide-vue-next"
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

export interface Skpd {
  id: string
  skpd: string
  jenis: string // "puskesmas", "kelurahan", "skpd"
  petugas_count: number
  created_date: string
  updated_date?: string
}

// Helper function untuk get jenis badge variant
const getJenisBadgeVariant = (jenis: string): "default" | "secondary" | "destructive" | "outline" => {
  switch (jenis.toLowerCase()) {
    case "puskesmas":
      return "default"
    case "kelurahan":
      return "secondary"
    case "skpd":
      return "outline"
    default:
      return "outline"
  }
}

// Helper function untuk get jenis badge color class
const getJenisColorClass = (jenis: string): string => {
  switch (jenis.toLowerCase()) {
    case "puskesmas":
      return "bg-blue-100 text-blue-800 border-blue-200"
    case "kelurahan":
      return "bg-green-100 text-green-800 border-green-200"
    case "skpd":
      return "bg-purple-100 text-purple-800 border-purple-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

export const columns: ColumnDef<Skpd>[] = [
  {
    id: "skpd",
    accessorKey: "skpd",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Nama SKPD", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const namaSkpd = row.getValue("skpd") as string

      return h("div", { class: "max-w-[300px]" }, [
        h("div", { class: "font-medium text-sm flex items-center gap-2" }, [
          h("span", {}, namaSkpd),
        ]),
        h("div", { class: "text-xs text-muted-foreground mt-1" }, `ID: ${row.original.id}`),
      ])
    },
    meta: {
      displayName: "Nama SKPD",
    },
  },
  {
    id: "jenis",
    accessorKey: "jenis",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Jenis", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const jenis = row.getValue("jenis") as string
      const jenisCapitalized = jenis.charAt(0).toUpperCase() + jenis.slice(1)

      return h(
        Badge,
        {
          variant: getJenisBadgeVariant(jenis),
          class: `text-xs ${getJenisColorClass(jenis)}`,
        },
        () => jenisCapitalized
      )
    },
    meta: {
      displayName: "Jenis SKPD",
    },
  },
  {
    id: "petugas_count",
    accessorKey: "petugas_count",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Petugas", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const petugasCount = row.getValue("petugas_count") as number

      return h("div", { class: "text-center" }, [
        h("div", { class: "flex items-center justify-center gap-1" }, [
          h(Users, { class: "h-4 w-4 text-muted-foreground" }),
          h("span", { class: "font-medium" }, petugasCount.toString()),
        ]),
        h("div", { class: "text-xs text-muted-foreground" }, "petugas"),
      ])
    },
    meta: {
      displayName: "Jumlah Petugas",
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
      const createdDate = row.getValue("created_date") as string
      const updatedDate = row.original.updated_date

      return h("div", { class: "text-sm space-y-1" }, [
        h("div", {}, new Date(createdDate).toLocaleDateString("id-ID")),
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
    id: "status",
    // Virtual column untuk menampilkan status berdasarkan jumlah petugas
    header: "Status",
    cell: ({ row }) => {
      const petugasCount = row.original.petugas_count

      let statusText = ""
      let statusVariant: "default" | "secondary" | "destructive" | "outline" = "outline"
      let statusClass = ""

      if (petugasCount === 0) {
        statusText = "Belum Ada Petugas"
        statusVariant = "destructive"
        statusClass = "bg-red-100 text-red-800 border-red-200"
      } else if (petugasCount >= 1 && petugasCount <= 3) {
        statusText = "Aktif"
        statusVariant = "default"
        statusClass = "bg-green-100 text-green-800 border-green-200"
      } else {
        statusText = "Aktif (Banyak Petugas)"
        statusVariant = "secondary"
        statusClass = "bg-blue-100 text-blue-800 border-blue-200"
      }

      return h(
        Badge,
        {
          variant: statusVariant,
          class: `text-xs ${statusClass}`,
        },
        () => statusText
      )
    },
    meta: {
      displayName: "Status",
    },
    // Tambahkan accessor function untuk sorting/filtering
    accessorFn: (row) => {
      const petugasCount = row.petugas_count
      if (petugasCount === 0) return "belum_ada_petugas"
      if (petugasCount >= 1 && petugasCount <= 3) return "aktif"
      return "aktif_banyak"
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const skpd = row.original

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
                { align: "end", class: "w-[200px]" },
                {
                  default: () => [
                    h(DropdownMenuLabel, {}, { default: () => "Actions" }),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(skpd.id)
                        },
                      },
                      { default: () => "Copy ID" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(new CustomEvent("edit-skpd", { detail: skpd }))
                        },
                      },
                      { default: () => [h(Pencil, { class: "mr-2 h-4 w-4" }), "Edit"] }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(new CustomEvent("delete-skpd", { detail: skpd }))
                        },
                        class: skpd.petugas_count > 0 ? "text-red-400 cursor-not-allowed" : "text-red-600",
                        disabled: skpd.petugas_count > 0,
                      },
                      {
                        default: () => [
                          h(Trash2, { class: "mr-2 h-4 w-4" }),
                          skpd.petugas_count > 0 ? "Tidak Bisa Dihapus" : "Delete",
                        ],
                      }
                    ),
                    skpd.petugas_count > 0 &&
                      h(
                        "div",
                        { class: "px-2 py-1" },
                        h(
                          "p",
                          { class: "text-xs text-muted-foreground" },
                          `Ada ${skpd.petugas_count} petugas terkait`
                        )
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