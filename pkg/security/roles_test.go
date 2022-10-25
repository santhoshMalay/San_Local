package security

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseRole(t *testing.T) {
	cases := map[string]struct {
		str     string
		success bool
		role    Role
	}{
		"student": {
			str:     "student",
			success: true,
			role:    Student,
		},
		"admin": {
			str:     "admin",
			success: true,
			role:    Admin,
		},
		"empty string": {
			str:     "",
			success: false,
		},
		"unknown value": {
			str:     "SomeOtherRole",
			success: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			role, err := ParseRole(tc.str)
			if tc.success {
				require.Equal(t, tc.role, role)
				require.NoError(t, err)
			} else {
				require.Equal(t, UndefinedRole, role)
				require.Error(t, err)
			}
		})
	}
}

func TestRole_UnmarshalSingle(t *testing.T) {
	type container struct {
		MyRole Role `json:"my_role"`
	}
	cases := map[string]struct {
		payload string
		success bool
		role    Role
	}{
		"success": {
			payload: `{"my_role":"student"}`,
			success: true,
			role:    Student,
		},
		"failure_unknown_role": {
			payload: `{"my_role":"SomeOtherRole"}`,
			success: false,
		},
		"failure_empty_value": {
			payload: `{}`,
			success: true,
			role:    UndefinedRole,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var result container
			err := json.Unmarshal([]byte(tc.payload), &result)
			if tc.success {
				require.NoError(t, err)
				expected := container{
					MyRole: tc.role,
				}
				require.Equal(t, expected, result)
			} else {
				require.Error(t, err)
				//t.Log(err)
			}
		})
	}
}

func TestRole_UnmarshalArray(t *testing.T) {
	type container struct {
		MyRoles []Role `json:"my_roles"`
	}
	cases := map[string]struct {
		payload string
		success bool
		roles   []Role
	}{
		"success_one_item": {
			payload: `{"my_roles":["student"]}`,
			success: true,
			roles:   []Role{Student},
		},
		"success_multiple_items": {
			payload: `{"my_roles":["student","admin"]}`,
			success: true,
			roles:   []Role{Student, Admin},
		},
		"success_empty": {
			payload: `{"my_roles":[]}`,
			success: true,
			roles:   []Role{},
		},
		"failure_single_unknown": {
			payload: `{"my_roles":["SomeOtherRole"]}`,
			success: false,
		},
		"failure_known_and_unknown": {
			payload: `{"my_roles":["student","SomeOtherRole"]}`,
			success: false,
		},
		"success_missing_array": {
			payload: `{}`,
			success: true,
			roles:   nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var result container
			err := json.Unmarshal([]byte(tc.payload), &result)
			if tc.success {
				require.NoError(t, err)
				require.ElementsMatch(t, tc.roles, result.MyRoles)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestRole_MarshalSingle(t *testing.T) {
	type container struct {
		MyRole Role `json:"my_role"`
	}
	cases := map[string]struct {
		data     container
		success  bool
		expected string
	}{
		"success_student": {
			data:     container{Student},
			success:  true,
			expected: `{"my_role":"student"}`,
		},
		"success_admin": {
			data:     container{Admin},
			success:  true,
			expected: `{"my_role":"admin"}`,
		},
		"failure_undefined": {
			data:    container{UndefinedRole},
			success: false,
		},
		"failure_arbitrary": {
			data:    container{42},
			success: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.data)
			if tc.success {
				require.NoError(t, err)
				require.Equal(t, tc.expected, string(bytes))
			} else {
				require.Error(t, err)
				//t.Log(err)
			}
		})
	}
}

func TestRole_MarshalArray(t *testing.T) {
	type container struct {
		MyRoles []Role `json:"my_roles"`
	}
	cases := map[string]struct {
		data     container
		success  bool
		expected string
	}{
		"success_single_item": {
			data:     container{[]Role{Student}},
			success:  true,
			expected: `{"my_roles":["student"]}`,
		},
		"success_multiple_items": {
			data:     container{[]Role{Student, Admin}},
			success:  true,
			expected: `{"my_roles":["student","admin"]}`,
		},
		"success_empty": {
			data:     container{[]Role{}},
			success:  true,
			expected: `{"my_roles":[]}`,
		},
		"success_missing_array": {
			data:     container{},
			success:  true,
			expected: `{"my_roles":null}`,
		},
		"failure_single_unknown": {
			data:    container{[]Role{UndefinedRole}},
			success: false,
		},
		"failure_known_and_unknown": {
			data:    container{[]Role{Student, UndefinedRole}},
			success: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			bytes, err := json.Marshal(tc.data)
			if tc.success {
				require.NoError(t, err)
				require.Equal(t, tc.expected, string(bytes))
			} else {
				require.Error(t, err)
				//t.Log(err)
			}
		})
	}
}

func TestRole_Valid(t *testing.T) {
	cases := map[string]struct {
		role        Role
		mustBeValid bool
	}{
		"student":   {Student, true},
		"admin":     {Admin, true},
		"undefined": {UndefinedRole, false},
		"arbitrary": {42, false},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.role.Valid()
			if tc.mustBeValid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
