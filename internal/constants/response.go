package constants

const (
	// Status Code Constants
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFailed  = "fail"

	// Database Error Messages
	MsgDbConnected     = "Koneksi database berhasil"
	MsgDbDisconnected  = "Koneksi database gagal"
	MsgMinioInitFailed = "Gagal menginisialisasi MinIO client: %v"
	MsgRedisInitFailed = "Gagal menginisialisasi Redis client: %v"
	
	// Storage Messages
	MsgFileOpenFailed    = "gagal membuka file: %w"
	MsgMinioUploadFailed = "gagal mengunggah file ke MinIO: %w"
	MsgMinioURLFailed    = "gagal membuat URL file: %w"
	MsgMinioDeleteFailed = "gagal menghapus file dari MinIO: %w"

	// General Messages
	MsgSuccess         = "Permintaan berhasil diproses"
	MsgInvalidRequest  = "Format request tidak valid"
	MsgInternalError   = "Terjadi kesalahan pada sistem"
	MsgUnauthorized    = "Sesi tidak valid atau telah berakhir"
	MsgFileNotFound    = "File tidak ditemukan"
	MsgFetchFileFailed = "Gagal mengambil file dari server"

	// Auth Messages
	MsgLoginSuccess    = "Login berhasil"
	MsgRegisterSuccess = "Admin berhasil didaftarkan"
	MsgRequiredFields  = "Email dan Password wajib diisi"

	// Biteship Messages
	MsgDestinationPostalCodeRequired = "destination_postal_code is required"
	MsgTrackingIDRequired            = "tracking_id is required"
	MsgSuccessFetchRates             = "Success fetch rates"
	MsgSuccessFetchCouriers          = "Success fetch couriers"
	MsgSuccessCreateOrder            = "Success create order"
	MsgSuccessFetchTracking          = "Success fetch tracking"
)
