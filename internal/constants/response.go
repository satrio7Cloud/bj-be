package constants

const (
	// Status Code Constants
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFailed  = "fail"

	// General Messages
	MsgSuccess         = "Permintaan berhasil diproses"
	MsgInvalidRequest  = "Format request tidak valid"
	MsgInternalError   = "Terjadi kesalahan pada sistem"
	MsgUnauthorized    = "Sesi tidak valid atau telah berakhir"

	// Auth Messages
	MsgLoginSuccess    = "Login berhasil"
	MsgRegisterSuccess = "Admin berhasil didaftarkan"
	MsgRequiredFields  = "Email dan Password wajib diisi"
)
