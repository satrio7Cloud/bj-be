package usecase

import (
	"errors"
	"os"
	"time"

	"bj/internal/domain/user/dto"
	"bj/internal/domain/user/entity"
	"bj/internal/domain/user/repository"
	"bj/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AdminUsecase interface {
	Login(req dto.LoginRequest) (dto.LoginResponse, error)
	Register(req dto.RegisterRequest) error
}

type adminUsecase struct {
	adminRepo repository.AdminRepository
}

func NewAdminUsecase(repo repository.AdminRepository) AdminUsecase {
	return &adminUsecase{repo}
}

func (u *adminUsecase) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	var resp dto.LoginResponse

	// 1. Cari Admin berdasarkan Email
	admin, err := u.adminRepo.FindByEmail(req.Email)
	if err != nil {
		return resp, errors.New("terjadi kesalahan pada sistem")
	}
	if admin == nil {
		return resp, errors.New("email atau password salah") // Jangan beritahu secara spesifik email tidak ditemukan demi keamanan
	}

	// 2. Cek apakah Admin aktif
	if !admin.IsActive {
		return resp, errors.New("akun anda telah dinonaktifkan")
	}

	// 3. Cocokkan Password (menggunakan utils Argon2)
	isMatch, err := utils.ComparePassword(req.Password, admin.Password)
	if err != nil || !isMatch {
		return resp, errors.New("email atau password salah")
	}

	// 4. Generate JWT Token
	// Kita ambil Secret Key dari .env — WAJIB diset, tidak ada fallback (keamanan).
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return resp, errors.New("konfigurasi server tidak lengkap: JWT_SECRET belum diset")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.ID,
		"email": admin.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return resp, errors.New("gagal membuat sesi login")
	}

	// 5. Susun Response
	resp.Token = tokenString
	resp.User = dto.AdminData{
		ID:       admin.ID,
		Fullname: admin.Fullname,
		Email:    admin.Email,
		Role:     "admin",
	}

	return resp, nil
}

func (u *adminUsecase) Register(req dto.RegisterRequest) error {
	// 1. Cek email apakah sudah ada
	existingAdmin, err := u.adminRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("terjadi kesalahan pada sistem saat mengecek email")
	}
	if existingAdmin != nil {
		return errors.New("email sudah terdaftar")
	}

	// 2. Hash password menggunakan Argon2
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("gagal mengamankan kata sandi")
	}

	// 3. Generate UUID untuk ID
	adminID := uuid.New().String()

	// 4. Susun Entity Admin
	newAdmin := &entity.Admin{
		ID:          adminID,
		Fullname:    req.Fullname,
		Email:       req.Email,
		Password:    hashedPassword,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true, // Langsung aktif saat register (atau sesuaikan dengan kebutuhan)
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 5. Simpan ke database
	err = u.adminRepo.Create(newAdmin)
	if err != nil {
		return errors.New("gagal menyimpan data admin baru")
	}

	return nil
}
