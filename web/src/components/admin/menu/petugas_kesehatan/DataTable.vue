<script setup lang="ts">
import { computed, ref } from "vue";
import { type PetugasKesehatan, columns } from "./columns"
import { FlexRender, getCoreRowModel, getFilteredRowModel, getPaginationRowModel, getSortedRowModel, useVueTable, type ColumnFiltersState, type SortingState, type VisibilityState } from "@tanstack/vue-table";
import { Check, ChevronDown, Eye, EyeOff, Search } from "lucide-vue-next";
import Input from "@/components/ui/input/Input.vue";
import DropdownMenu from "@/components/ui/dropdown-menu/DropdownMenu.vue";
import DropdownMenuTrigger from "@/components/ui/dropdown-menu/DropdownMenuTrigger.vue";
import Button from "@/components/ui/button/Button.vue";
import DropdownMenuContent from "@/components/ui/dropdown-menu/DropdownMenuContent.vue";
import DropdownMenuCheckboxItem from "@/components/ui/dropdown-menu/DropdownMenuCheckboxItem.vue";
import Table from "@/components/ui/table/Table.vue";
import TableHeader from "@/components/ui/table/TableHeader.vue";
import TableRow from "@/components/ui/table/TableRow.vue";
import TableHead from "@/components/ui/table/TableHead.vue";
import TableBody from "@/components/ui/table/TableBody.vue";
import TableCell from "@/components/ui/table/TableCell.vue";

interface Props {
  data: PetugasKesehatan[]
}

const props = defineProps<Props>()

const sorting = ref<SortingState>([])
const columnFilters = ref<ColumnFiltersState>([])
const columnVisibility = ref<VisibilityState>({})
const rowSelection = ref({})

const table = useVueTable({
  get data() {
    return props.data
  },
  columns,
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
  getFilteredRowModel: getFilteredRowModel(),
  onSortingChange: (updaterOrValue) => {
    sorting.value =
      typeof updaterOrValue === "function" ? updaterOrValue(sorting.value) : updaterOrValue
  },
  onColumnFiltersChange: (updaterOrValue) => {
    columnFilters.value =
      typeof updaterOrValue === "function" ? updaterOrValue(columnFilters.value) : updaterOrValue
  },
  onColumnVisibilityChange: (updaterOrValue) => {
    columnVisibility.value =
      typeof updaterOrValue === "function" ? updaterOrValue(columnVisibility.value) : updaterOrValue
  },
  onRowSelectionChange: (updaterOrValue) => {
    rowSelection.value =
      typeof updaterOrValue === "function" ? updaterOrValue(rowSelection.value) : updaterOrValue
  },
  state: {
    get sorting() {
      return sorting.value
    },
    get columnFilters() {
      return columnFilters.value
    },
    get columnVisibility() {
      return columnVisibility.value
    },
    get rowSelection() {
      return rowSelection.value
    },
  },
})

const searchValue = computed({
  get: () => (table.getColumn("nama")?.getFilterValue() as string) ?? "",
  set: (value) => table.getColumn("nama")?.setFilterValue(value),
})

// Helper function to get display name
const getColumnDisplayName = (columnId: string) => {
  const column = table.getColumn(columnId)
  const meta = column?.columnDef.meta as { displayName?: string } | undefined
  return meta?.displayName || columnId
}

// Toggle column visibility handler
const toggleColumnVisibility = (columnId: string) => {
  const column = table.getColumn(columnId)
  if (column) {
    column.toggleVisibility()
  }
}

// Check if column is visible
const isColumnVisible = (columnId: string) => {
  const column = table.getColumn(columnId)
  return column ? column.getIsVisible() : true
}

// Get visible columns count
const visibleColumnsCount = computed(() => {
  const hidableColumns = table.getAllColumns().filter((column) => column.getCanHide())
  return hidableColumns.filter((column) => column.getIsVisible()).length
})

const totalColumnsCount = computed(() => {
  return table.getAllColumns().filter((column) => column.getCanHide()).length
})
</script>

<template>
  
  <div class="w-full">
    <div class="flex items-center py-4 gap-4">
      <!-- Search -->
      <div class="relative flex-1 max-w-sm">
        <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
        <Input
          v-model="searchValue"
          placeholder="Cari nama petugas..."
          class="pl-10" />
      </div>

      <!-- Column Visibility -->
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button
            variant="outline"
            class="gap-2">
            <ChevronDown class="h-4 w-4" />
            Kolom
            <span class="ml-1 text-xs text-muted-foreground">
              ({{ visibleColumnsCount }}/{{ totalColumnsCount }})
            </span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent
          align="end"
          class="w-[220px]">
          <div class="p-2 text-xs text-muted-foreground border-b">
            Pilih kolom yang ingin ditampilkan
          </div>

          <DropdownMenuCheckboxItem
            v-for="column in table.getAllColumns().filter((column) => column.getCanHide())"
            :key="column.id"
            :checked="isColumnVisible(column.id)"
            @select="(e: Event) => e.preventDefault()"
            @click="toggleColumnVisibility(column.id)"
            class="cursor-pointer">
            <div class="flex items-center space-x-2 w-full">
              <div class="flex items-center justify-center w-4 h-4">
                <Check
                  v-if="isColumnVisible(column.id)"
                  class="h-3 w-3 text-primary" />
              </div>
              <span class="flex-1">{{ getColumnDisplayName(column.id) }}</span>
              <div class="flex items-center">
                <Eye
                  v-if="isColumnVisible(column.id)"
                  class="h-3 w-3 text-green-500" />
                <EyeOff
                  v-else
                  class="h-3 w-3 text-gray-400" />
              </div>
            </div>
          </DropdownMenuCheckboxItem>
          <div class="p-2 pt-2 border-t mt-1">
            <div class="flex justify-between text-xs text-muted-foreground">
              <span>Tampil: {{ visibleColumnsCount }}</span>
              <span>Total: {{ totalColumnsCount }}</span>
            </div>
          </div>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <!-- Table -->
    <div class="rounded-md border">
      <Table class="table-auto">
        <TableHeader>
          <TableRow
            v-for="headerGroup in table.getHeaderGroups()"
            :key="headerGroup.id">
            <TableHead
              v-for="header in headerGroup.headers"
              :key="header.id">
              <FlexRender
                v-if="!header.isPlaceholder"
                :render="header.column.columnDef.header"
                :props="header.getContext()" />
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <template v-if="table.getRowModel().rows?.length">
            <TableRow
              v-for="row in table.getRowModel().rows"
              :key="row.id"
              :data-state="row.getIsSelected() && 'selected'">
              <TableCell
                v-for="cell in row.getVisibleCells()"
                :key="cell.id">
                <FlexRender
                  :render="cell.column.columnDef.cell"
                  :props="cell.getContext()" />
              </TableCell>
            </TableRow>
          </template>

          <TableRow v-else>
            <TableCell
              :colspan="columns.length"
              class="h-24 text-center">
              Tidak ada data petugas.
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>

    <!-- Pagination -->
    <div
      class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between md:space-x-2 py-4">
      <div class="flex flex-col gap-3 md:flex-row md:items-center md:space-x-6 lg:space-x-8">
        <div class="flex items-center justify-center space-x-2">
          <p class="text-sm font-medium">Baris per halaman</p>
          <select
            :value="table.getState().pagination.pageSize"
            @change="table.setPageSize(Number(($event.target as HTMLSelectElement).value))"
            class="h-8 w-[70px] rounded border border-input bg-background px-3 py-1 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2">
            <option
              v-for="pageSize in [10, 20, 30, 40, 50]"
              :key="pageSize"
              :value="pageSize">
              {{ pageSize }}
            </option>
          </select>
        </div>
        <div class="flex w-full md:w-[100px] items-center justify-center text-sm font-medium">
          Hal {{ table.getState().pagination.pageIndex + 1 }} dari {{ table.getPageCount() }}
        </div>
        <div class="flex items-center justify-center space-x-2">
          <Button
            variant="outline"
            size="sm"
            :disabled="!table.getCanPreviousPage()"
            @click="table.previousPage()">
            Sebelumnya
          </Button>
          <Button
            variant="outline"
            size="sm"
            :disabled="!table.getCanNextPage()"
            @click="table.nextPage()">
            Selanjutnya
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>