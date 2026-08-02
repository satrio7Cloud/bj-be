package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Fullname    string `json:"fullname" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  AdminData  `json:"user"`
}

// AdminData digunakan agar password dan data sensitif lain tidak ikut terkirim ke frontend
type AdminData struct {
	ID       string `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
