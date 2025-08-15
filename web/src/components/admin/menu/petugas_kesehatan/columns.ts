import type { ColumnDef } from "@tanstack/vue-table"
import {
  ArrowUpDown,
  MoreHorizontal,
  Pencil,
  Trash2,
  Users,
  Mail,
  Shield,
  Activity,
} from "lucide-vue-next"
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

export interface PetugasKesehatan {
  id: string
  id_pengguna: string
  id_skpd: string
  email: string
  nama: string
  skpd: string
  jenis_skpd: string // "puskesmas", "kelurahan", "skpd"
  intervensi_count: number
  created_date: string
  updated_date?: string
}

// Helper function untuk get jenis SKPD badge variant
const getJenisSkpdBadgeVariant = (
  jenis: string
): "default" | "secondary" | "destructive" | "outline" => {
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

// Helper function untuk get jenis SKPD color class
const getJenisSkpdColorClass = (jenis: string): string => {
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

// Helper function untuk get status berdasarkan intervensi count
const getStatusInfo = (intervensiCount: number) => {
  if (intervensiCount === 0) {
    return {
      text: "Belum Ada Intervensi",
      variant: "outline" as const,
      class: "bg-gray-100 text-gray-600 border-gray-300",
    }
  } else if (intervensiCount >= 1 && intervensiCount <= 5) {
    return {
      text: "Aktif",
      variant: "default" as const,
      class: "bg-green-100 text-green-800 border-green-200",
    }
  } else {
    return {
      text: "Sangat Aktif",
      variant: "secondary" as const,
      class: "bg-blue-100 text-blue-800 border-blue-200",
    }
  }
}

export const columns: ColumnDef<PetugasKesehatan>[] = [
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
        () => ["Nama Petugas", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const nama = row.getValue("nama") as string
      const email = row.original.email

      return h("div", { class: "max-w-[250px]" }, [
        h("div", { class: "font-medium text-sm flex items-center gap-2" }, [
          h(Users, { class: "h-4 w-4 text-blue-600" }),
          h("span", {}, nama),
        ]),
        h("div", { class: "text-xs text-muted-foreground mt-1 flex items-center gap-1" }, [
          h(Mail, { class: "h-3 w-3" }),
          h("span", {}, email),
        ]),
        h("div", { class: "text-xs text-blue-600 mt-1" }, `ID: ${row.original.id}`),
      ])
    },
    meta: {
      displayName: "Data Petugas",
    },
  },
  {
    id: "skpd_info",
    // Virtual column yang menggabungkan info SKPD
    header: "SKPD",
    cell: ({ row }) => {
      const skpd = row.original.skpd
      const jenisSkpd = row.original.jenis_skpd

      return h("div", { class: "max-w-[200px]" }, [
        h("div", { class: "flex items-center gap-2 mb-2" }, [
          h(
            Badge,
            {
              variant: getJenisSkpdBadgeVariant(jenisSkpd),
              class: `text-xs ${getJenisSkpdColorClass(jenisSkpd)}`,
            },
            () => [
              h("span", {}, jenisSkpd.charAt(0).toUpperCase() + jenisSkpd.slice(1)),
            ]
          ),
        ]),
        h("div", { class: "font-medium text-sm" }, skpd),
        h("div", { class: "text-xs text-muted-foreground" }, `SKPD ID: ${row.original.id_skpd}`),
      ])
    },
    meta: {
      displayName: "Info SKPD",
    },
    // Accessor function untuk sorting/filtering
    accessorFn: (row) => `${row.jenis_skpd} ${row.skpd}`,
  },
  {
    id: "intervensi_count",
    accessorKey: "intervensi_count",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Intervensi", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const intervensiCount = row.getValue("intervensi_count") as number

      return h("div", { class: "text-center" }, [
        h("div", { class: "flex items-center justify-center gap-1 mb-1" }, [
          h(Activity, { class: "h-4 w-4 text-muted-foreground" }),
          h("span", { class: "font-bold text-lg" }, intervensiCount.toString()),
        ]),
        h("div", { class: "text-xs text-muted-foreground" }, "intervensi"),
      ])
    },
    meta: {
      displayName: "Jumlah Intervensi",
    },
  },
  {
    id: "status",
    // Virtual column untuk status berdasarkan aktivitas
    header: "Status Aktivitas",
    cell: ({ row }) => {
      const intervensiCount = row.original.intervensi_count
      const statusInfo = getStatusInfo(intervensiCount)

      return h(
        Badge,
        {
          variant: statusInfo.variant,
          class: `text-xs ${statusInfo.class}`,
        },
        () => statusInfo.text
      )
    },
    meta: {
      displayName: "Status Aktivitas",
    },
    // Accessor function untuk sorting/filtering
    accessorFn: (row) => {
      const intervensiCount = row.intervensi_count
      if (intervensiCount === 0) return "belum_ada_intervensi"
      if (intervensiCount >= 1 && intervensiCount <= 5) return "aktif"
      return "sangat_aktif"
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
        h("div", { class: "flex items-center gap-1" }, [
          h("span", { class: "text-xs text-muted-foreground" }, "Dibuat:"),
          h("span", {}, new Date(createdDate).toLocaleDateString("id-ID")),
        ]),
        updatedDate &&
          h("div", { class: "flex items-center gap-1" }, [
            h("span", { class: "text-xs text-muted-foreground" }, "Update:"),
            h("span", { class: "text-xs" }, new Date(updatedDate).toLocaleDateString("id-ID")),
          ]),
      ])
    },
    meta: {
      displayName: "Tanggal Dibuat",
    },
  },
  {
    id: "account_info",
    // Virtual column untuk info akun
    header: "Info Akun",
    cell: ({ row }) => {
      const idPengguna = row.original.id_pengguna
      // const email = row.original.email

      return h("div", { class: "text-sm space-y-1" }, [
        h("div", { class: "flex items-center gap-1" }, [
          h(Shield, { class: "h-3 w-3 text-green-600" }),
          h("span", { class: "text-xs text-green-600 font-medium" }, "Petugas Kesehatan"),
        ]),
        h("div", { class: "text-xs text-muted-foreground" }, `User ID: ${idPengguna}`),
        h("div", { class: "text-xs text-muted-foreground flex items-center gap-1" }, [
          h("span", {}, "✓ Active Account"),
        ]),
      ])
    },
    meta: {
      displayName: "Info Akun",
    },
    // Accessor function untuk filtering
    accessorFn: (row) => `${row.email} ${row.id_pengguna}`,
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const petugas = row.original

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
                { align: "end", class: "w-[220px]" },
                {
                  default: () => [
                    h(DropdownMenuLabel, {}, { default: () => "Actions" }),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(petugas.id)
                        },
                      },
                      { default: () => "Copy Petugas ID" }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(petugas.email)
                        },
                      },
                      { default: () => "Copy Email" }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(petugas.nama)
                        },
                      },
                      { default: () => "Copy Nama" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          // Navigate to detail view atau show detail modal
                          console.log("View details for:", petugas.nama)
                        },
                      },
                      { default: () => [h(Users, { class: "mr-2 h-4 w-4" }), "View Details"] }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("edit-petugas", { detail: petugas })
                          )
                        },
                      },
                      { default: () => [h(Pencil, { class: "mr-2 h-4 w-4" }), "Edit"] }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("delete-petugas", { detail: petugas })
                          )
                        },
                        class:
                          petugas.intervensi_count > 0
                            ? "text-red-400 cursor-not-allowed"
                            : "text-red-600",
                        disabled: petugas.intervensi_count > 0,
                      },
                      {
                        default: () => [
                          h(Trash2, { class: "mr-2 h-4 w-4" }),
                          petugas.intervensi_count > 0 ? "Tidak Bisa Dihapus" : "Delete",
                        ],
                      }
                    ),
                    petugas.intervensi_count > 0 &&
                      h(
                        "div",
                        { class: "px-2 py-1 border-t mt-1" },
                        h(
                          "p",
                          { class: "text-xs text-red-600 font-medium" },
                          `⚠️ Ada ${petugas.intervensi_count} intervensi terkait`
                        )
                      ),
                    h(
                      "div",
                      { class: "px-2 py-1 border-t" },
                      h(
                        "p",
                        { class: "text-xs text-muted-foreground" },
                        `SKPD: ${petugas.jenis_skpd} - ${petugas.skpd}`
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
