package seeder

import (
	"log"
	"time"

	"service/internal/pkg/config"
	"service/internal/pkg/model"
)

type ReportSeeder struct{}

// ── Text variants ─────────────────────────────────────────────────────────────
// Setiap field memiliki beberapa varian konten yang dirotasi per user × hari.
// Teks sengaja mengandung: newline, list bernomor/huruf, dan karakter unik (emoji, tanda baca khusus).

var completedYesterdayVariants = []string{
	"Menyelesaikan review & merge PR #38 — fitur notifikasi push.\n\nDetail pekerjaan:\n1. Review kode dari 2 kontributor\n2. Fix konflik merge di branch `feature/notification`\n3. Update unit test yang gagal akibat perubahan interface\n4. Deploy ke environment staging ✅",

	"Progress kemarin:\na. Implementasi endpoint `GET /reports/export` lengkap dengan groupBy\nb. Tambah validasi input (fromDate ≤ toDate) agar tidak error 500\nc. Tulis integrasi test — semua 12 kasus pass 🎉\nd. Dokumentasi Swagger sudah diupdate",

	"• Debugging issue koneksi Redis yang timeout secara intermittent\n• Root cause ditemukan: pool size terlalu kecil (default 10 → naik ke 50)\n• Hotfix sudah di-deploy ke production pukul 14:30 WIB 🚀\n• Monitoring 1 jam setelahnya — tidak ada error lagi",

	"Kemarin fokus pada refactor modul autentikasi:\n1) Pisahkan logic JWT ke package tersendiri (`pkg/auth`)\n2) Hilangkan duplikasi kode di 3 handler (Login, Refresh, Logout)\n3) Tambahkan middleware rate-limiter untuk endpoint `/login` (maks 5 req/menit)\n\nTotal baris kode berkurang ~120 baris 💪",

	"Meeting & koordinasi:\n- Sesi planning sprint Q2 bersama Product & Design (2 jam)\n- Breakdown task menjadi 14 story points\n- Estimasi selesai: akhir minggu ke-2\n\nSelain itu sempat fix bug #201: tanggal \"créé le\" tidak tampil di PDF export karena karakter è tidak di-escape dengan benar.",

	"Pekerjaan selesai kemarin:\na) Setup CI/CD pipeline baru menggunakan GitHub Actions\nb) Konfigurasi environment variables via GitHub Secrets\nc) Tes pipeline dari branch `develop` → berhasil build & test otomatis ✅\nd) Buat dokumentasi onboarding untuk developer baru",

	"• Analisis performa query SQL pada tabel `reports` (≈ 2 juta baris)\n• Tambahkan composite index pada kolom (`userId`, `reportDate`)\n• Query time turun dari ~1.2 s → 48 ms 📉\n• Buat laporan benchmark dan bagikan ke tim via Slack",
}

var planTodayVariants = []string{
	"Rencana hari ini:\n1. Lanjutkan implementasi fitur export CSV\n2. Fix bug #217: karakter spesial (', \", &) tidak di-encode di response JSON\n3. Code review PR dari Dian & Budi\n4. Update changelog versi 2.4.0",

	"Target hari ini:\na. Selesaikan migrasi database tabel `report_sessions` — tambah kolom `completedAt`\nb. Pastikan zero-downtime migration dengan strategi expand/contract\nc. Test di staging sebelum deploy ke production\nd. Buat runbook untuk proses rollback jika gagal",

	"• Mulai develop modul dashboard — wireframe sudah ada dari Design\n• Fokus dulu pada komponen tabel dengan fitur sort & filter\n• Koordinasi dengan tim FE soal format response API (nested vs flat)\n• Target: skeleton komponen selesai sebelum EOD 🎯",

	"Hari ini akan:\n1) Implementasi pagination cursor-based untuk endpoint `/reports` (gantikan offset)\n2) Benchmark perbandingan offset vs cursor pada 1 juta baris\n3) Update dokumentasi API\n4) PR siap untuk review sebelum pukul 17:00 WIB",

	"Fokus hari ini:\n- Selesaikan unit test untuk `ReportExportService` (target coverage 80%)\n- Tambah test case untuk karakter unicode & multi-line input\n- Fix linter warning yang tersisa (3 issues di `pkg/docx`)\n- Sync dengan tim QA untuk jadwal UAT minggu depan",

	"Plan:\na) Review & respond komentar di 2 PR yang masih open\nb) Implementasi soft-delete pada endpoint `DELETE /users/:id`\nc) Pastikan data tidak benar-benar terhapus — hanya `deletedAt` di-set\nd) Tambah migration dan seed test data\ne) EOD check-in dengan lead",

	"• Lanjutkan spike riset WebSocket untuk fitur real-time notification\n• Buat PoC sederhana — server → klien push update laporan baru\n• Dokumentasikan temuan dalam ADR (Architecture Decision Record)\n• Estimasi effort jika dilanjutkan ke production",
}

var finishEstimationVariants = []string{
	"Semua task utama target selesai pukul 17:00 WIB.\nJika ada blocker dari pihak DevOps (deployment), mungkin mundur ke besok pagi.",

	"Estimasi selesai EOD hari ini — sekitar pukul 16:30.\nBagian migrasi DB mungkin butuh waktu ekstra ±1 jam jika ada issue di staging.",

	"Target: PR siap review pukul 15:00, merge & deploy ke staging pukul 17:00.\nKalau review dari senior butuh revisi, paling lambat besok siang.",

	"Pekerjaan pengembangan selesai hari ini.\nUntuk deploy ke production menunggu approval tim Ops — estimasi besok atau lusa.",

	"Unit test & PR selesai hari ini pukul 16:00.\nJadwal UAT dengan QA: Senin minggu depan 🗓️",

	"Estimasi kelar sebelum standup sore (15:30).\nKomponen FE butuh approval dari Design dulu — mungkin ada sedikit revisi.",

	"PoC selesai hari ini, ADR draf selesai besok.\nKeputusan final soal implementasi production dibahas di sprint planning Senin.",
}

var blockersVariants = []string{
	"Tidak ada blocker saat ini. Semua dependency sudah tersedia ✅",

	"Ada sedikit hambatan:\n- Akses ke environment staging terbatas hari ini (maintenance terjadwal 10:00–12:00)\n- Akan bekerja di task lain dulu selama window tersebut",

	"Menunggu feedback dari tim Design untuk komponen modal — belum ada respons sejak kemarin sore.\nSudah ping ulang via Slack, tunggu balasan.",

	"Tidak ada blocker teknis.\nHanya perlu konfirmasi dari Product soal edge case: bagaimana jika user submit laporan di luar jam kerja? 🤔",

	"• PR #44 masih menunggu approval dari 1 reviewer (sudah 2 hari)\n• Tanpa merge PR tersebut, task hari ini tidak bisa dimulai\n• Sudah escalate ke lead",

	"Tidak ada hambatan. Semua lancar 🚀",

	"Koneksi VPN kantor tidak stabil sejak pagi — kadang putus sendiri.\nSudah lapor ke IT, katanya sedang diperbaiki. Sementara menggunakan hotspot.",
}

var moodVariants = []string{
	"Semangat! Tidur cukup kemarin, siap produktif hari ini 💪😊",
	"Baik. Sedikit lelah karena meeting marathon kemarin, tapi masih fokus ✅",
	"Sangat baik! Senang karena bug kritis akhirnya ketemu root cause-nya 🎉",
	"Oke. Mood standar — tidak ada yang spesial, tapi kerja tetap jalan 🙂",
	"Agak kurang fit, sedikit flu. Tapi bisa WFH jadi masih bisa produktif 🤧💊",
	"Excited! Hari ini mulai fitur baru yang sudah ditunggu-tunggu sejak lama 🚀✨",
	"Baik sekali. Kopi pagi sudah, musik lo-fi on, siap coding 🎵☕️",
}

// ── Seeder ────────────────────────────────────────────────────────────────────

func (s *ReportSeeder) Seed() {
	db := config.PgSQL

	workdays := lastWorkdays(7)
	if len(workdays) == 0 {
		log.Println("[ReportSeeder] No workdays found, skipping.")
		return
	}

	created := 0
	skipped := 0

	for userID := uint(1); userID <= 14; userID++ {
		for dayIdx, date := range workdays {
			// Cek apakah session sudah ada (hindari duplikat)
			var existing model.ReportSession
			err := db.
				Where(`"userId" = ? AND "reportDate" = ?`, userID, date).
				First(&existing).Error
			if err == nil {
				skipped++
				continue
			}

			// Pilih varian konten berdasarkan rotasi (userID + dayIdx)
			pick := func(variants []string) string {
				idx := (int(userID-1) + dayIdx) % len(variants)
				return variants[idx]
			}

			completedAt := date.Add(17*time.Hour + time.Duration(userID)*time.Minute)

			// Buat session (completed)
			session := model.ReportSession{
				UserID:      userID,
				ReportDate:  date,
				CurrentStep: 5,
				IsCompleted: true,
				StartedAt:   date.Add(9 * time.Hour),
				CompletedAt: &completedAt,
			}
			if err := db.Create(&session).Error; err != nil {
				log.Printf("[ReportSeeder] Failed to create session user=%d date=%s: %v", userID, date.Format("2006-01-02"), err)
				continue
			}

			// Buat report
			report := model.Report{
				UserID:             userID,
				ReportDate:         date,
				CompletedYesterday: pick(completedYesterdayVariants),
				PlanToday:          pick(planTodayVariants),
				FinishEstimation:   pick(finishEstimationVariants),
				Blockers:           pick(blockersVariants),
				Mood:               pick(moodVariants),
			}
			if err := db.Create(&report).Error; err != nil {
				log.Printf("[ReportSeeder] Failed to create report user=%d date=%s: %v", userID, date.Format("2006-01-02"), err)
				// Rollback session yang sudah dibuat
				db.Delete(&session)
				continue
			}

			created++
		}
	}

	log.Printf("[ReportSeeder] Done — created: %d, skipped (already exist): %d", created, skipped)
}

// lastWorkdays mengembalikan slice hari kerja (Senin–Jumat) dari `n` hari
// kalender terakhir, diurutkan dari yang paling lama ke yang paling baru.
func lastWorkdays(n int) []time.Time {
	var days []time.Time
	now := time.Now()

	for i := 1; i <= n*2 && len(days) < n; i++ {
		d := now.AddDate(0, 0, -i)
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			continue
		}
		// Normalisasi ke midnight (00:00:00) zona lokal
		days = append([]time.Time{
			time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()),
		}, days...)
	}

	return days
}
