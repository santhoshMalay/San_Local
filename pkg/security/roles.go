package security

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Role uint8

const (
	UndefinedRole Role = iota
	Student
	Admin
)

const (
	roleStudent = "student"
	roleAdmin   = "admin"
)

var roles = map[string]Role{
	roleStudent: Student,
	roleAdmin:   Admin,
}

func ParseRole(str string) (Role, error) {
	role, ok := roles[str]
	if !ok {
		return UndefinedRole, fmt.Errorf("undefined role: %q", str)
	}
	return role, nil
}

//goland:noinspection GoMixedReceiverTypes
func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	role, err := ParseRole(s)
	if err != nil {
		return err
	}
	*r = role
	return nil
}

//goland:noinspection GoMixedReceiverTypes
func (r Role) MarshalJSON() ([]byte, error) {
	if err := r.Valid(); err != nil {
		return nil, err
	}
	str := r.String()
	return json.Marshal(str)
}

//goland:noinspection GoMixedReceiverTypes
func (r *Role) Valid() error {
	switch *r {
	case Student, Admin:
		return nil
	}
	return errors.New("undefined role value")
}

//goland:noinspection GoMixedReceiverTypes
func (r *Role) String() string {
	switch *r {
	case Student:
		return roleStudent
	case Admin:
		return roleAdmin
	default:
		return "undefined_role"
	}
}

func ToRoles(from []uint8) []Role {
	
	to := make([]Role, 0, len(from))
	for _, v := range from {
		to = append(to, Role(v))
	}
	return to
}
