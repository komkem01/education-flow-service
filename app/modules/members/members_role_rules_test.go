package members

import (
	"testing"

	"education-flow/app/modules/entities/ent"
)

func TestValidateStudentRoleExclusivity(t *testing.T) {
	testCases := []struct {
		name      string
		existing  []ent.MemberRole
		target    ent.MemberRole
		wantError bool
	}{
		{
			name:      "add teacher when no roles",
			existing:  nil,
			target:    ent.MemberRoleTeacher,
			wantError: false,
		},
		{
			name:      "keep student only",
			existing:  []ent.MemberRole{ent.MemberRoleStudent},
			target:    ent.MemberRoleStudent,
			wantError: false,
		},
		{
			name:      "add student to multi role member",
			existing:  []ent.MemberRole{ent.MemberRoleTeacher, ent.MemberRoleAdmin},
			target:    ent.MemberRoleStudent,
			wantError: true,
		},
		{
			name:      "add teacher when member is student",
			existing:  []ent.MemberRole{ent.MemberRoleStudent},
			target:    ent.MemberRoleTeacher,
			wantError: true,
		},
		{
			name:      "add teacher for multi-role non-student member",
			existing:  []ent.MemberRole{ent.MemberRoleAdmin, ent.MemberRoleStaff},
			target:    ent.MemberRoleTeacher,
			wantError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateStudentRoleExclusivity(testCase.existing, testCase.target)
			if testCase.wantError && err == nil {
				t.Fatalf("validateStudentRoleExclusivity() expected error, got nil")
			}
			if !testCase.wantError && err != nil {
				t.Fatalf("validateStudentRoleExclusivity() expected nil, got %v", err)
			}
		})
	}
}
