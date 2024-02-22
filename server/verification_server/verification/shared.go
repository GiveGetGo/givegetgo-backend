package verification

type GeneralVerificationResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// enum for the verification event, register and reset password
const (
	RegisterEvent = "register"
	ResetEvent    = "reset"
)
