package models

type RegisterForm struct {
	FullName        string `json:"full_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ForgotPasswordForm struct {
	Email string `json:"email"`
}

type ResetPasswordForm struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UpdateUserInfoForm struct {
	FullName string `json:"full_name"`
	Location string `json:"location"`
	Bio      string `json:"bio"`
	Web      string `json:"web"`
}

type ChangePasswordForm struct {
	PasswordCurrent string `json:"password_current"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type ResendVerificationForm struct {
	Type      string `json:"type"`
	Recipient string `json:"recipient"`
}

type ActivateTFAForm struct {
	Secret string `json:"secret"`
	Code   string `json:"code"`
}

type RemoveTFAForm struct {
	Password string `json:"password"`
}

type AuthenticationForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyTFAForm struct {
	Code string `json:"code"`
}
