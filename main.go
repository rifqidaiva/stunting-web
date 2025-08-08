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
	http.HandleFunc("/api/auth/profile", auth.GetUserProfile)

	/* ===================
	   Admin API Endpoints
	====================== */

	// SKPD Management
	http.HandleFunc("/api/admin/skpd/get", admin.AdminSkpdGet)
	http.HandleFunc("/api/admin/skpd/insert", admin.AdminSkpdInsert)
	http.HandleFunc("/api/admin/skpd/update", admin.AdminSkpdUpdate)
	http.HandleFunc("/api/admin/skpd/delete", admin.AdminSkpdDelete)
	http.HandleFunc("/api/admin/skpd/restore", admin.AdminSkpdRestore)

	// Petugas Kesehatan Management
	http.HandleFunc("/api/admin/petugas-kesehatan/get", admin.AdminPetugasKesehatanGet)
	http.HandleFunc("/api/admin/petugas-kesehatan/insert", admin.AdminPetugasKesehatanInsert)
	http.HandleFunc("/api/admin/petugas-kesehatan/update", admin.AdminPetugasKesehatanUpdate)
	http.HandleFunc("/api/admin/petugas-kesehatan/delete", admin.AdminPetugasKesehatanDelete)
	http.HandleFunc("/api/admin/petugas-kesehatan/restore", admin.AdminPetugasKesehatanRestore)

	// Keluarga Management
	http.HandleFunc("/api/admin/keluarga/get", admin.AdminKeluargaGet)
	http.HandleFunc("/api/admin/keluarga/insert", admin.AdminKeluargaInsert)
	http.HandleFunc("/api/admin/keluarga/update", admin.AdminKeluargaUpdate)
	http.HandleFunc("/api/admin/keluarga/delete", admin.AdminKeluargaDelete)
	http.HandleFunc("/api/admin/keluarga/restore", admin.AdminKeluargaRestore)

	// Balita Management
	http.HandleFunc("/api/admin/balita/get", admin.AdminBalitaGet)
	http.HandleFunc("/api/admin/balita/insert", admin.AdminBalitaInsert)
	http.HandleFunc("/api/admin/balita/update", admin.AdminBalitaUpdate)
	http.HandleFunc("/api/admin/balita/delete", admin.AdminBalitaDelete)
	http.HandleFunc("/api/admin/balita/restore", admin.AdminBalitaRestore)

	// Laporan Masyarakat Management
	http.HandleFunc("/api/admin/laporan-masyarakat/get", admin.AdminLaporanMasyarakatGet)
	http.HandleFunc("/api/admin/laporan-masyarakat/insert", admin.AdminLaporanMasyarakatInsert)
	http.HandleFunc("/api/admin/laporan-masyarakat/update", admin.AdminLaporanMasyarakatUpdate)
	http.HandleFunc("/api/admin/laporan-masyarakat/delete", admin.AdminLaporanMasyarakatDelete)
	http.HandleFunc("/api/admin/laporan-masyarakat/restore", admin.AdminLaporanMasyarakatRestore)

	// Intervensi Management (Admin)
	// http.HandleFunc("/api/admin/intervensi/get", admin.AdminIntervensiGet)
	// http.HandleFunc("/api/admin/intervensi/insert", admin.AdminIntervensiInsert)
	// http.HandleFunc("/api/admin/intervensi/update", admin.AdminIntervensiUpdate)
	// http.HandleFunc("/api/admin/intervensi/delete", admin.AdminIntervensiDelete)
	// http.HandleFunc("/api/admin/intervensi/restore", admin.AdminIntervensiRestore)

	// Riwayat Pemeriksaan Management (Admin)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/get", admin.AdminRiwayatPemeriksaanGet)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/insert", admin.AdminRiwayatPemeriksaanInsert)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/update", admin.AdminRiwayatPemeriksaanUpdate)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/delete", admin.AdminRiwayatPemeriksaanDelete)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/restore", admin.AdminRiwayatPemeriksaanRestore)

	// Intervensi Petugas (Junction Table)
	// http.HandleFunc("/api/admin/intervensi-petugas/get", admin.AdminIntervensiPetugasGet)
	// http.HandleFunc("/api/admin/intervensi-petugas/assign", admin.AdminIntervensiPetugasAssign)
	// http.HandleFunc("/api/admin/intervensi-petugas/remove", admin.AdminIntervensiPetugasRemove)

	// Master Data Management
	// http.HandleFunc("/api/admin/status-laporan/get", admin.AdminStatusLaporanGet)
	// http.HandleFunc("/api/admin/masyarakat/get", admin.AdminMasyarakatGet)
	// http.HandleFunc("/api/admin/kecamatan/get", admin.AdminKecamatanGet)
	// http.HandleFunc("/api/admin/kelurahan/get", admin.AdminKelurahanGet)

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
