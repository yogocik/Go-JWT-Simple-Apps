package presentation

type UserForm struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"user_password" binding:"required"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Username string `json:"user_name" binding:"required"`
	Password string `json:"user_password" binding:"required"`
}
