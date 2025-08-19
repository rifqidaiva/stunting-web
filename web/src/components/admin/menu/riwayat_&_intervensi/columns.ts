import type { ColumnDef } from "@tanstack/vue-table"
import {
  ArrowUpDown,
  MoreHorizontal,
  Pencil,
  Trash2,
  Users,
  FileText,
  Stethoscope,
  Calendar,
  Baby,
  UserCheck,
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

// ===============================
// RIWAYAT PEMERIKSAAN INTERFACE
// ===============================
export interface RiwayatPemeriksaan {
  id: string
  id_balita: string
  nama_balita: string
  umur_balita: string
  jenis_kelamin: string
  nama_ayah: string
  nama_ibu: string
  nomor_kk: string
  id_intervensi: string
  jenis_intervensi: string
  tanggal_intervensi: string
  id_laporan_masyarakat: string
  status_laporan: string
  tanggal_laporan: string
  jenis_laporan: string // "masyarakat" | "admin"
  tanggal: string
  berat_badan: string
  tinggi_badan: string
  status_gizi: string // "normal" | "stunting" | "gizi buruk"
  keterangan: string
  kelurahan: string
  kecamatan: string
  created_date: string
  updated_date?: string
  created_by?: string
  updated_by?: string
}

// ===============================
// INTERVENSI INTERFACE
// ===============================
export interface Intervensi {
  id: string
  id_balita: string
  nama_balita: string
  jenis: string // "gizi" | "kesehatan" | "sosial"
  tanggal: string
  deskripsi: string
  hasil: string
  petugas_count: number
  riwayat_count: number
  created_date: string
  updated_date?: string
  created_by?: string
  updated_by?: string
}

// ===============================
// HELPER FUNCTIONS
// ===============================

// Status gizi badge variants
const getStatusGiziBadgeVariant = (
  status: string
): "default" | "secondary" | "destructive" | "outline" => {
  switch (status.toLowerCase()) {
    case "normal":
      return "default"
    case "stunting":
      return "destructive"
    case "gizi buruk":
      return "destructive"
    default:
      return "outline"
  }
}

// Status gizi color classes
const getStatusGiziColorClass = (status: string): string => {
  switch (status.toLowerCase()) {
    case "normal":
      return "bg-green-100 text-green-800 border-green-200"
    case "stunting":
      return "bg-red-100 text-red-800 border-red-200"
    case "gizi buruk":
      return "bg-red-100 text-red-800 border-red-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

// Jenis intervensi badge variants
const getJenisIntervensiBadgeVariant = (
  jenis: string
): "default" | "secondary" | "destructive" | "outline" => {
  switch (jenis.toLowerCase()) {
    case "gizi":
      return "default"
    case "kesehatan":
      return "secondary"
    case "sosial":
      return "outline"
    default:
      return "outline"
  }
}

// Jenis intervensi color classes
const getJenisIntervensiColorClass = (jenis: string): string => {
  switch (jenis.toLowerCase()) {
    case "gizi":
      return "bg-green-100 text-green-800 border-green-200"
    case "kesehatan":
      return "bg-blue-100 text-blue-800 border-blue-200"
    case "sosial":
      return "bg-purple-100 text-purple-800 border-purple-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

// Jenis laporan color classes
const getJenisLaporanColorClass = (jenis: string): string => {
  switch (jenis.toLowerCase()) {
    case "masyarakat":
      return "bg-blue-100 text-blue-800 border-blue-200"
    case "admin":
      return "bg-orange-100 text-orange-800 border-orange-200"
    default:
      return "bg-gray-100 text-gray-800 border-gray-200"
  }
}

// Status intervensi berdasarkan petugas count
const getStatusIntervensi = (petugasCount: number) => {
  if (petugasCount === 0) {
    return {
      text: "Belum Ada Petugas",
      variant: "destructive" as const,
      class: "bg-red-100 text-red-800 border-red-200",
    }
  } else if (petugasCount >= 1 && petugasCount <= 2) {
    return {
      text: "Aktif",
      variant: "default" as const,
      class: "bg-green-100 text-green-800 border-green-200",
    }
  } else {
    return {
      text: "Tim Lengkap",
      variant: "secondary" as const,
      class: "bg-blue-100 text-blue-800 border-blue-200",
    }
  }
}

// =================================
// MARK: RIWAYAT PEMERIKSAAN COLUMNS
// =================================
export const riwayatPemeriksaanColumns: ColumnDef<RiwayatPemeriksaan>[] = [
  {
    id: "balita_info",
    header: "Data Balita",
    cell: ({ row }) => {
      const namaBalita = row.original.nama_balita
      const umur = row.original.umur_balita
      const jenisKelamin = row.original.jenis_kelamin
      const nomorKk = row.original.nomor_kk

      return h("div", { class: "max-w-[200px]" }, [
        h("div", { class: "font-medium text-sm flex items-center gap-2" }, [
          h(Baby, { class: "h-4 w-4 text-blue-600" }),
          h("span", {}, namaBalita),
        ]),
        h("div", { class: "text-xs text-muted-foreground mt-1" }, [
          h("div", {}, `${umur} • ${jenisKelamin === "L" ? "Laki-laki" : "Perempuan"}`),
          h("div", {}, `KK: ${nomorKk.substring(0, 12)}...`),
        ]),
      ])
    },
    meta: {
      displayName: "Data Balita",
    },
    accessorFn: (row) => `${row.nama_balita} ${row.umur_balita} ${row.nomor_kk}`,
  },
  {
    id: "tanggal",
    accessorKey: "tanggal",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Tanggal Periksa", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const tanggal = row.getValue("tanggal") as string
      const tanggalIntervensi = row.original.tanggal_intervensi

      return h("div", { class: "text-sm" }, [
        h("div", { class: "font-medium flex items-center gap-1" }, [
          h(Calendar, { class: "h-3 w-3" }),
          h("span", {}, new Date(tanggal).toLocaleDateString("id-ID")),
        ]),
        h(
          "div",
          { class: "text-xs text-muted-foreground" },
          `Intervensi: ${new Date(tanggalIntervensi).toLocaleDateString("id-ID")}`
        ),
      ])
    },
    meta: {
      displayName: "Tanggal Pemeriksaan",
    },
  },
  {
    id: "hasil_pemeriksaan",
    header: "Hasil Pemeriksaan",
    cell: ({ row }) => {
      const beratBadan = row.original.berat_badan
      const tinggiBadan = row.original.tinggi_badan
      const statusGizi = row.original.status_gizi

      return h("div", { class: "space-y-2" }, [
        h("div", { class: "text-sm" }, [
          h("div", {}, `BB: ${beratBadan} kg`),
          h("div", {}, `TB: ${tinggiBadan} cm`),
        ]),
        h(
          Badge,
          {
            variant: getStatusGiziBadgeVariant(statusGizi),
            class: `text-xs ${getStatusGiziColorClass(statusGizi)}`,
          },
          () => statusGizi.charAt(0).toUpperCase() + statusGizi.slice(1)
        ),
      ])
    },
    meta: {
      displayName: "Hasil Pemeriksaan",
    },
    accessorFn: (row) => `${row.berat_badan} ${row.tinggi_badan} ${row.status_gizi}`,
  },
  {
    id: "intervensi_info",
    header: "Info Intervensi",
    cell: ({ row }) => {
      const jenisIntervensi = row.original.jenis_intervensi

      return h("div", { class: "space-y-1" }, [
        h(
          Badge,
          {
            variant: getJenisIntervensiBadgeVariant(jenisIntervensi),
            class: `text-xs ${getJenisIntervensiColorClass(jenisIntervensi)}`,
          },
          () => jenisIntervensi.charAt(0).toUpperCase() + jenisIntervensi.slice(1)
        ),
        h("div", { class: "text-xs text-muted-foreground" }, `ID: ${row.original.id_intervensi}`),
      ])
    },
    meta: {
      displayName: "Info Intervensi",
    },
    accessorFn: (row) => row.jenis_intervensi,
  },
  {
    id: "laporan_info",
    header: "Info Laporan",
    cell: ({ row }) => {
      const jenisLaporan = row.original.jenis_laporan
      const statusLaporan = row.original.status_laporan
      const tanggalLaporan = row.original.tanggal_laporan

      return h("div", { class: "space-y-1" }, [
        h(
          Badge,
          {
            variant: "outline",
            class: `text-xs ${getJenisLaporanColorClass(jenisLaporan)}`,
          },
          () => jenisLaporan.charAt(0).toUpperCase() + jenisLaporan.slice(1)
        ),
        h("div", { class: "text-xs text-muted-foreground" }, [
          h("div", {}, statusLaporan),
          h("div", {}, new Date(tanggalLaporan).toLocaleDateString("id-ID")),
        ]),
      ])
    },
    meta: {
      displayName: "Info Laporan",
    },
    accessorFn: (row) => `${row.jenis_laporan} ${row.status_laporan}`,
  },
  {
    id: "wilayah",
    header: "Wilayah",
    cell: ({ row }) => {
      const kelurahan = row.original.kelurahan
      const kecamatan = row.original.kecamatan

      return h("div", { class: "text-sm" }, [
        h("div", { class: "font-medium" }, kelurahan),
        h("div", { class: "text-xs text-muted-foreground" }, kecamatan),
      ])
    },
    meta: {
      displayName: "Wilayah",
    },
    accessorFn: (row) => `${row.kelurahan} ${row.kecamatan}`,
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
      const createdBy = row.original.created_by

      return h("div", { class: "text-sm space-y-1" }, [
        h("div", {}, new Date(createdDate).toLocaleDateString("id-ID")),
        updatedDate &&
          h(
            "div",
            { class: "text-xs text-muted-foreground" },
            `Update: ${new Date(updatedDate).toLocaleDateString("id-ID")}`
          ),
        createdBy && h("div", { class: "text-xs text-blue-600" }, `by ${createdBy}`),
      ])
    },
    meta: {
      displayName: "Info Pembuatan",
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const riwayat = row.original

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
                { align: "end", class: "w-[250px]" },
                {
                  default: () => [
                    h(DropdownMenuLabel, {}, { default: () => "Actions" }),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(riwayat.id)
                        },
                      },
                      { default: () => "Copy Riwayat ID" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("edit-riwayat", { detail: riwayat })
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
                            new CustomEvent("delete-riwayat", { detail: riwayat })
                          )
                        },
                        class: "text-red-600",
                      },
                      { default: () => [h(Trash2, { class: "mr-2 h-4 w-4" }), "Delete"] }
                    ),
                    h(
                      "div",
                      { class: "px-2 py-1 border-t" },
                      h(
                        "p",
                        { class: "text-xs text-muted-foreground" },
                        `Balita: ${riwayat.nama_balita} • ${riwayat.status_gizi}`
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

// ========================
// MARK: INTERVENSI COLUMNS
// ========================
export const intervensiColumns: ColumnDef<Intervensi>[] = [
  {
    id: "balita_info",
    header: "Data Balita",
    cell: ({ row }) => {
      const namaBalita = row.original.nama_balita
      const idBalita = row.original.id_balita

      return h("div", { class: "max-w-[180px]" }, [
        h("div", { class: "font-medium text-sm flex items-center gap-2" }, [
          h(Baby, { class: "h-4 w-4 text-blue-600" }),
          h("span", {}, namaBalita),
        ]),
        h("div", { class: "text-xs text-muted-foreground mt-1" }, `ID: ${idBalita}`),
      ])
    },
    meta: {
      displayName: "Data Balita",
    },
    accessorFn: (row) => row.nama_balita,
  },
  {
    id: "jenis_tanggal",
    header: "Jenis & Tanggal",
    cell: ({ row }) => {
      const jenis = row.original.jenis
      const tanggal = row.original.tanggal

      return h("div", { class: "space-y-2" }, [
        h(
          Badge,
          {
            variant: getJenisIntervensiBadgeVariant(jenis),
            class: `text-xs ${getJenisIntervensiColorClass(jenis)}`,
          },
          () => jenis.charAt(0).toUpperCase() + jenis.slice(1)
        ),
        h("div", { class: "text-sm flex items-center gap-1" }, [
          h(Calendar, { class: "h-3 w-3 text-muted-foreground" }),
          h("span", {}, new Date(tanggal).toLocaleDateString("id-ID")),
        ]),
      ])
    },
    meta: {
      displayName: "Jenis & Tanggal",
    },
    accessorFn: (row) => `${row.jenis} ${row.tanggal}`,
  },
  {
    id: "deskripsi",
    accessorKey: "deskripsi",
    header: "Deskripsi",
    cell: ({ row }) => {
      const deskripsi = row.getValue("deskripsi") as string
      const hasil = row.original.hasil

      return h("div", { class: "max-w-[250px] space-y-1" }, [
        h(
          "div",
          { class: "text-sm font-medium" },
          deskripsi.length > 60 ? `${deskripsi.substring(0, 60)}...` : deskripsi
        ),
        h(
          "div",
          { class: "text-xs text-muted-foreground" },
          `Hasil: ${hasil.length > 40 ? `${hasil.substring(0, 40)}...` : hasil}`
        ),
      ])
    },
    meta: {
      displayName: "Deskripsi & Hasil",
    },
  },
  {
    id: "petugas_status",
    header: "Status Petugas",
    cell: ({ row }) => {
      const petugasCount = row.original.petugas_count
      const statusInfo = getStatusIntervensi(petugasCount)

      return h("div", { class: "text-center space-y-1" }, [
        h("div", { class: "flex items-center justify-center gap-1" }, [
          h(Users, { class: "h-4 w-4 text-muted-foreground" }),
          h("span", { class: "font-bold" }, petugasCount.toString()),
        ]),
        h(
          Badge,
          {
            variant: statusInfo.variant,
            class: `text-xs ${statusInfo.class}`,
          },
          () => statusInfo.text
        ),
      ])
    },
    meta: {
      displayName: "Status Petugas",
    },
    accessorFn: (row) => {
      const petugasCount = row.petugas_count
      if (petugasCount === 0) return "belum_ada_petugas"
      if (petugasCount >= 1 && petugasCount <= 2) return "aktif"
      return "tim_lengkap"
    },
  },
  {
    id: "riwayat_count",
    accessorKey: "riwayat_count",
    header: ({ column }) => {
      return h(
        Button,
        {
          variant: "ghost",
          onClick: () => column.toggleSorting(column.getIsSorted() === "asc"),
        },
        () => ["Riwayat", h(ArrowUpDown, { class: "ml-2 h-4 w-4" })]
      )
    },
    cell: ({ row }) => {
      const riwayatCount = row.getValue("riwayat_count") as number

      return h("div", { class: "text-center" }, [
        h("div", { class: "flex items-center justify-center gap-1" }, [
          h(Stethoscope, { class: "h-4 w-4 text-muted-foreground" }),
          h("span", { class: "font-bold text-lg" }, riwayatCount.toString()),
        ]),
        h("div", { class: "text-xs text-muted-foreground" }, "pemeriksaan"),
      ])
    },
    meta: {
      displayName: "Jumlah Riwayat",
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
      const createdBy = row.original.created_by

      return h("div", { class: "text-sm space-y-1" }, [
        h("div", {}, new Date(createdDate).toLocaleDateString("id-ID")),
        updatedDate &&
          h(
            "div",
            { class: "text-xs text-muted-foreground" },
            `Update: ${new Date(updatedDate).toLocaleDateString("id-ID")}`
          ),
        createdBy && h("div", { class: "text-xs text-blue-600" }, `by ${createdBy}`),
      ])
    },
    meta: {
      displayName: "Info Pembuatan",
    },
  },
  {
    id: "actions",
    enableHiding: false,
    cell: ({ row }) => {
      const intervensi = row.original

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
                { align: "end", class: "w-[250px]" },
                {
                  default: () => [
                    h(DropdownMenuLabel, {}, { default: () => "Actions" }),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          navigator.clipboard.writeText(intervensi.id)
                        },
                      },
                      { default: () => "Copy Intervensi ID" }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          // Manage petugas assignments
                          console.log("Manage petugas for:", intervensi.id)
                        },
                      },
                      { default: () => [h(UserCheck, { class: "mr-2 h-4 w-4" }), "Kelola Petugas"] }
                    ),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          // View riwayat pemeriksaan
                          document.dispatchEvent(
                            new CustomEvent("view-riwayat-intervensi", { detail: intervensi })
                          )
                        },
                      },
                      { default: () => [h(FileText, { class: "mr-2 h-4 w-4" }), "Lihat Riwayat"] }
                    ),
                    h(DropdownMenuSeparator),
                    h(
                      DropdownMenuItem,
                      {
                        onClick: () => {
                          document.dispatchEvent(
                            new CustomEvent("edit-intervensi", { detail: intervensi })
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
                            new CustomEvent("delete-intervensi", { detail: intervensi })
                          )
                        },
                        class:
                          intervensi.riwayat_count > 0 || intervensi.petugas_count > 0
                            ? "text-red-400 cursor-not-allowed"
                            : "text-red-600",
                        disabled: intervensi.riwayat_count > 0 || intervensi.petugas_count > 0,
                      },
                      {
                        default: () => [
                          h(Trash2, { class: "mr-2 h-4 w-4" }),
                          intervensi.riwayat_count > 0 || intervensi.petugas_count > 0
                            ? "Tidak Bisa Dihapus"
                            : "Delete",
                        ],
                      }
                    ),
                    (intervensi.riwayat_count > 0 || intervensi.petugas_count > 0) &&
                      h(
                        "div",
                        { class: "px-2 py-1 border-t" },
                        h(
                          "p",
                          { class: "text-xs text-red-600 font-medium" },
                          `⚠️ Ada ${intervensi.riwayat_count} riwayat & ${intervensi.petugas_count} petugas terkait`
                        )
                      ),
                    h(
                      "div",
                      { class: "px-2 py-1 border-t" },
                      h(
                        "p",
                        { class: "text-xs text-muted-foreground" },
                        `Balita: ${intervensi.nama_balita} • ${intervensi.jenis}`
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
