package security

// UserPrincipal represents core user identification data (claims) passed via security tokens
type UserPrincipal struct {
	UserId string
	Roles  []Role
}

func (up *UserPrincipal) HasRole(role Role) bool {
	for _, r := range up.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (up *UserPrincipal) IsAdmin() bool {
	return up.HasRole(Admin)
}
