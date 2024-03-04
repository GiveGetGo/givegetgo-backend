package verification

type GeneralVerificationResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

// enum for the verification event, register and reset password
const (
	RegisterEvent      = "register"
	ResetPasswordEvent = "reset-password"
)
