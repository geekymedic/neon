package types

const (
	StateName          = "neon.bff.State"
	ResponseStatusCode = "neon.bff.response.statuscode"
	ResponseErr        = "neon.bff.response.msg"
	NeonSession        = "neon.bff.session"
)

// resolve cycle import
const (
	CodeSuccess     = 0
	CodeServerError = 1006
)
