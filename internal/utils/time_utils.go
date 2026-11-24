package utils

import (
	"time"
)

// WIBLocation merepresentasikan zona waktu resmi Asia/Jakarta (UTC+7 / WIB).
var WIBLocation = time.FixedZone("WIB", 7*3600)

// ToWIB mengunci dan mengonversi waktu ke zona waktu Asia/Jakarta (+07:00).
//
// Mengapa ini penting?
// Driver database SQL Server membaca tipe data DATETIME/DATETIME2 tanpa offset zona waktu
// sebagai UTC. Jika dikirim langsung ke frontend sebagai UTC ("...Z"), browser di Indonesia (+7)
// akan otomatis menambahkan 7 jam lagi. Contoh: 20:53 WIB dianggap 20:53 UTC oleh browser,
// lalu ditambah 7 jam menjadi 03:53 keesokan harinya (anomali tanggal loncat 1 hari).
//
// ToWIB mempertahankan angka jam dinding (wall-clock) asli dari database dan secara eksplisit
// menempelkan offset zona waktu WIB (+07:00), sehingga frontend menerima format ISO 8601
// yang presisi: "YYYY-MM-DDTHH:mm:ss+07:00".
func ToWIB(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	// Jika waktu terbaca sebagai UTC dari database driver
	if t.Location() == time.UTC {
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), WIBLocation)
	}
	return t.In(WIBLocation)
}

// ToWIBPtr mengonversi pointer *time.Time ke zona waktu WIB (+07:00).
// Jika pointer nil, fungsi ini dengan aman mengembalikan nil.
func ToWIBPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	res := ToWIB(*t)
	return &res
}

// NowWIB mengembalikan waktu saat ini yang terkunci pada zona waktu Asia/Jakarta (WIB).
func NowWIB() time.Time {
	return time.Now().In(WIBLocation)
}

