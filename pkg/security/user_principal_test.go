package security

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUserPrincipal_HasRole(t *testing.T) {
	cases := map[string]struct {
		data     *UserPrincipal
		in       Role
		expected bool
	}{
		"student_has_student": {
			data: &UserPrincipal{
				Roles: []Role{Student},
			},
			in:       Student,
			expected: true,
		},
		"student_does_not_have_admin": {
			data: &UserPrincipal{
				Roles: []Role{Student},
			},
			in:       Admin,
			expected: false,
		},
		"multi_has_student": {
			data: &UserPrincipal{
				Roles: []Role{Student, Admin},
			},
			in:       Student,
			expected: true,
		},
		"multi_has_admin": {
			data: &UserPrincipal{
				Roles: []Role{Student, Admin},
			},
			in:       Admin,
			expected: true,
		},
		"empty_does_not_have_student": {
			data:     &UserPrincipal{},
			in:       Student,
			expected: false,
		},
		"empty_does_not_have_admin": {
			data:     &UserPrincipal{},
			in:       Admin,
			expected: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := tc.data.HasRole(tc.in)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUserPrincipal_IsAdmin(t *testing.T) {
	cases := map[string]struct {
		data     *UserPrincipal
		expected bool
	}{
		"student_is_not_admin": {
			data: &UserPrincipal{
				Roles: []Role{Student},
			},
			expected: false,
		},
		"admin_is_admin": {
			data: &UserPrincipal{
				Roles: []Role{Admin},
			},
			expected: true,
		},
		"multi_is_admin": {
			data: &UserPrincipal{
				Roles: []Role{Student, Admin},
			},
			expected: true,
		},
		"empty_is_not_admin": {
			data:     &UserPrincipal{},
			expected: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result := tc.data.IsAdmin()
			require.Equal(t, tc.expected, result)
		})
	}
}
