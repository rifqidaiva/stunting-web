package main

import (
	"fmt"
	"net/http"

	_ "github.com/rifqidaiva/stunting-web/docs" // Import for Swagger documentation
	"github.com/rifqidaiva/stunting-web/internal/api"
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
	// http.HandleFunc("/api/auth/petugas/login", auth.PetugasLogin)
	// http.HandleFunc("/api/auth/petugas/profile", auth.GetPetugasProfile)

	// SKPD Management
	http.HandleFunc("/api/admin/skpd/get", api.AdminSkpdGet)
	http.HandleFunc("/api/admin/skpd/insert", api.AdminSkpdInsert)
	http.HandleFunc("/api/admin/skpd/update", api.AdminSkpdUpdate)
	http.HandleFunc("/api/admin/skpd/delete", api.AdminSkpdDelete)
	http.HandleFunc("/api/admin/skpd/restore", api.AdminSkpdRestore)

	// Petugas Kesehatan Management
	http.HandleFunc("/api/admin/petugas-kesehatan/get", api.AdminPetugasKesehatanGet)
	http.HandleFunc("/api/admin/petugas-kesehatan/insert", api.AdminPetugasKesehatanInsert)
	http.HandleFunc("/api/admin/petugas-kesehatan/update", api.AdminPetugasKesehatanUpdate)
	http.HandleFunc("/api/admin/petugas-kesehatan/delete", api.AdminPetugasKesehatanDelete)
	http.HandleFunc("/api/admin/petugas-kesehatan/restore", api.AdminPetugasKesehatanRestore)

	// Keluarga Management
	http.HandleFunc("/api/admin/keluarga/get", api.AdminKeluargaGet)
	http.HandleFunc("/api/admin/keluarga/insert", api.AdminKeluargaInsert)
	http.HandleFunc("/api/admin/keluarga/update", api.AdminKeluargaUpdate)
	http.HandleFunc("/api/admin/keluarga/delete", api.AdminKeluargaDelete)
	http.HandleFunc("/api/admin/keluarga/restore", api.AdminKeluargaRestore)

	// Balita Management
	http.HandleFunc("/api/admin/balita/get", api.AdminBalitaGet)
	http.HandleFunc("/api/admin/balita/insert", api.AdminBalitaInsert)
	http.HandleFunc("/api/admin/balita/update", api.AdminBalitaUpdate)
	http.HandleFunc("/api/admin/balita/delete", api.AdminBalitaDelete)
	http.HandleFunc("/api/admin/balita/restore", api.AdminBalitaRestore)

	// Laporan Masyarakat Management
	http.HandleFunc("/api/admin/laporan-masyarakat/get", api.AdminLaporanMasyarakatGet)
	http.HandleFunc("/api/admin/laporan-masyarakat/insert", api.AdminLaporanMasyarakatInsert)
	http.HandleFunc("/api/admin/laporan-masyarakat/update", api.AdminLaporanMasyarakatUpdate)
	http.HandleFunc("/api/admin/laporan-masyarakat/delete", api.AdminLaporanMasyarakatDelete)
	http.HandleFunc("/api/admin/laporan-masyarakat/restore", api.AdminLaporanMasyarakatRestore)

	// Intervensi Management (Admin)
	// http.HandleFunc("/api/admin/intervensi/get", api.AdminIntervensiGet)
	// http.HandleFunc("/api/admin/intervensi/insert", api.AdminIntervensiInsert)
	// http.HandleFunc("/api/admin/intervensi/update", api.AdminIntervensiUpdate)
	// http.HandleFunc("/api/admin/intervensi/delete", api.AdminIntervensiDelete)
	// http.HandleFunc("/api/admin/intervensi/restore", api.AdminIntervensiRestore)

	// Riwayat Pemeriksaan Management (Admin)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/get", api.AdminRiwayatPemeriksaanGet)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/insert", api.AdminRiwayatPemeriksaanInsert)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/update", api.AdminRiwayatPemeriksaanUpdate)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/delete", api.AdminRiwayatPemeriksaanDelete)
	// http.HandleFunc("/api/admin/riwayat-pemeriksaan/restore", api.AdminRiwayatPemeriksaanRestore)

	// Intervensi Petugas (Junction Table)
	// http.HandleFunc("/api/admin/intervensi-petugas/get", api.AdminIntervensiPetugasGet)
	// http.HandleFunc("/api/admin/intervensi-petugas/assign", api.AdminIntervensiPetugasAssign)
	// http.HandleFunc("/api/admin/intervensi-petugas/remove", api.AdminIntervensiPetugasRemove)

	// Master Data Management
	// http.HandleFunc("/api/admin/status-laporan/get", api.AdminStatusLaporanGet)
	// http.HandleFunc("/api/admin/masyarakat/get", api.AdminMasyarakatGet)
	// http.HandleFunc("/api/admin/kecamatan/get", api.AdminKecamatanGet)
	// http.HandleFunc("/api/admin/kelurahan/get", api.AdminKelurahanGet)

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
