package main

import (
	"fmt"
	"net/http"

	_ "github.com/rifqidaiva/stunting-web/docs" // Import for Swagger documentation
	"github.com/rifqidaiva/stunting-web/internal/api/admin"
	"github.com/rifqidaiva/stunting-web/internal/api/auth"
	"github.com/rifqidaiva/stunting-web/internal/object"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Stunting Web API
// @version 0.0.2
// @description API for managing stunting data
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func main() {
	// Authentication
	http.HandleFunc("/api/auth/login", auth.Login)
	http.HandleFunc("/api/auth/register", auth.Register)
	http.HandleFunc("/api/auth/register_admin", auth.RegisterAdmin)
	http.HandleFunc("/api/auth/profile", auth.UserProfileGet)

	/* ===================
	   Admin API Endpoints
	====================== */

	// SKPD Management
	http.HandleFunc("/api/admin/skpd/get", admin.SKPDGet)
	http.HandleFunc("/api/admin/skpd/insert", admin.SKPDInsert)
	http.HandleFunc("/api/admin/skpd/update", admin.SKPDUpdate)
	http.HandleFunc("/api/admin/skpd/delete", admin.SKPDDelete)
	http.HandleFunc("/api/admin/skpd/restore", admin.SKPDRestore)

	// Petugas Kesehatan Management
	http.HandleFunc("/api/admin/petugas-kesehatan/get", admin.PetugasKesehatanGet)
	http.HandleFunc("/api/admin/petugas-kesehatan/insert", admin.PetugasKesehatanInsert)
	http.HandleFunc("/api/admin/petugas-kesehatan/update", admin.PetugasKesehatanUpdate)
	http.HandleFunc("/api/admin/petugas-kesehatan/delete", admin.PetugasKesehatanDelete)
	http.HandleFunc("/api/admin/petugas-kesehatan/restore", admin.PetugasKesehatanRestore)

	// Keluarga Management
	http.HandleFunc("/api/admin/keluarga/get", admin.KeluargaGet)
	http.HandleFunc("/api/admin/keluarga/insert", admin.KeluargaInsert)
	http.HandleFunc("/api/admin/keluarga/update", admin.KeluargaUpdate)
	http.HandleFunc("/api/admin/keluarga/delete", admin.KeluargaDelete)
	http.HandleFunc("/api/admin/keluarga/restore", admin.KeluargaRestore)

	// Balita Management
	http.HandleFunc("/api/admin/balita/get", admin.BalitaGet)
	http.HandleFunc("/api/admin/balita/insert", admin.BalitaInsert)
	http.HandleFunc("/api/admin/balita/update", admin.BalitaUpdate)
	http.HandleFunc("/api/admin/balita/delete", admin.BalitaDelete)
	http.HandleFunc("/api/admin/balita/restore", admin.BalitaRestore)

	// Laporan Masyarakat Management
	http.HandleFunc("/api/admin/laporan-masyarakat/get", admin.LaporanMasyarakatGet)
	http.HandleFunc("/api/admin/laporan-masyarakat/insert", admin.LaporanMasyarakatInsert)
	http.HandleFunc("/api/admin/laporan-masyarakat/update", admin.LaporanMasyarakatUpdate)
	http.HandleFunc("/api/admin/laporan-masyarakat/delete", admin.LaporanMasyarakatDelete)
	http.HandleFunc("/api/admin/laporan-masyarakat/restore", admin.LaporanMasyarakatRestore)

	// Intervensi Management
	http.HandleFunc("/api/admin/intervensi/get", admin.IntervensiGet)
	http.HandleFunc("/api/admin/intervensi/insert", admin.IntervensiInsert)
	http.HandleFunc("/api/admin/intervensi/update", admin.IntervensiUpdate)
	http.HandleFunc("/api/admin/intervensi/delete", admin.IntervensiDelete)
	http.HandleFunc("/api/admin/intervensi/restore", admin.IntervensiRestore)

	// Riwayat Pemeriksaan Management
	http.HandleFunc("/api/admin/riwayat-pemeriksaan/get", admin.RiwayatPemeriksaanGet)
	http.HandleFunc("/api/admin/riwayat-pemeriksaan/insert", admin.RiwayatPemeriksaanInsert)
	http.HandleFunc("/api/admin/riwayat-pemeriksaan/update", admin.RiwayatPemeriksaanUpdate)
	http.HandleFunc("/api/admin/riwayat-pemeriksaan/delete", admin.RiwayatPemeriksaanDelete)
	http.HandleFunc("/api/admin/riwayat-pemeriksaan/restore", admin.RiwayatPemeriksaanRestore)

	// Intervensi Petugas (Junction Table)
	http.HandleFunc("/api/admin/intervensi-petugas/get", admin.IntervensiPetugasGet)
	http.HandleFunc("/api/admin/intervensi-petugas/assign", admin.IntervensiPetugasAssign)
	http.HandleFunc("/api/admin/intervensi-petugas/remove", admin.IntervensiPetugasRemove)

	// Master Data Management
	// http.HandleFunc("/api/admin/status-laporan/get", admin.StatusLaporanGet)
	// http.HandleFunc("/api/admin/masyarakat/get", admin.MasyarakatGet)
	// http.HandleFunc("/api/admin/kecamatan/get", admin.KecamatanGet)
	// http.HandleFunc("/api/admin/kelurahan/get", admin.KelurahanGet)

	/* ========================
	   Masyarakat API Endpoints
	=========================== */

	// // Masyarakat - Keluarga Management (untuk input data keluarga balita yang dilaporkan)
	// http.HandleFunc("/api/masyarakat/keluarga/insert", admin.MasyarakatKeluargaInsert)
	// http.HandleFunc("/api/masyarakat/keluarga/get", admin.MasyarakatKeluargaGet)
	// http.HandleFunc("/api/masyarakat/keluarga/update", admin.MasyarakatKeluargaUpdate)

	// // Masyarakat - Balita Management (untuk input data balita yang dilaporkan)
	// http.HandleFunc("/api/masyarakat/balita/insert", admin.MasyarakatBalitaInsert)
	// http.HandleFunc("/api/masyarakat/balita/get", admin.MasyarakatBalitaGet)
	// http.HandleFunc("/api/masyarakat/balita/update", admin.MasyarakatBalitaUpdate)

	// // Masyarakat - Laporan Management (untuk melaporkan balita)
	// http.HandleFunc("/api/masyarakat/laporan/insert", admin.MasyarakatLaporanInsert)
	// http.HandleFunc("/api/masyarakat/laporan/get", admin.MasyarakatLaporanGet)

	// // Masyarakat - Master Data (untuk dropdown/reference)
	// http.HandleFunc("/api/masyarakat/kelurahan/get", admin.MasyarakatKelurahanGet)
	// http.HandleFunc("/api/masyarakat/kecamatan/get", admin.MasyarakatKecamatanGet)
	// http.HandleFunc("/api/masyarakat/status-laporan/get", admin.MasyarakatStatusLaporanGet)

	// API test endpoint
	http.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		response := object.NewResponse(http.StatusOK, "Test API is working", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Swagger documentation endpoint
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("/swagger/doc.json", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe("localhost:8080", nil)
}
