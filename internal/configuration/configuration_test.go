package configuration

import (
	errs "app/internal/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserIdentifierRequest_Validate(t *testing.T) {
	type args struct {
		usr UserIdentifierRequest
	}
	type exp struct {
		errors []errs.ValidationError
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				usr: UserIdentifierRequest{
					UserId: "bc5bfa3b-8270-4aaf-b80b-f51836268747",
				},
			},
			exp: exp{
				errors: nil,
			},
		},
		{
			name: "invalid values format request body",
			args: args{
				usr: UserIdentifierRequest{
					UserId: "bc5bfa3b_8270_4aaf_b80b_f51836268747",
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID format",
					},
				},
			},
		},
		{
			name: "invalid min values request body",
			args: args{
				usr: UserIdentifierRequest{
					UserId: "bc5bfa3b-8270-4aaf-b80b-f5183626874",
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID length: 35",
					},
				},
			},
		},
		{
			name: "invalid max values request body",
			args: args{
				usr: UserIdentifierRequest{
					UserId: string(make([]byte, 37)),
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID length: 37",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp.errors, tc.args.usr.Validate())
		})
	}
}

func TestUserRequest_Validate(t *testing.T) {
	type args struct {
		usr UserRequest
	}
	type exp struct {
		errors []errs.ValidationError
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				usr: UserRequest{
					UserId: "bc5bfa3b-8270-4aaf-b80b-f51836268747",
					Name:   "John Doe",
					Age:    41,
				},
			},
			exp: exp{
				errors: nil,
			},
		},
		{
			name: "invalid values format request body",
			args: args{
				usr: UserRequest{
					UserId: "bc5bfa3b_8270_4aaf_b80b_f51836268747",
					Name:   "John Doe",
					Age:    41,
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID format",
					},
				},
			},
		},
		{
			name: "invalid min values request body",
			args: args{
				usr: UserRequest{
					UserId: "bc5bfa3b-8270-4aaf-b80b-f5183626874",
					Name:   "usr",
					Age:    0,
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID length: 35",
					},
					{
						Path:    "name",
						Message: "invalid length exceeded - (4-100)",
					},
					{
						Path:    "age",
						Message: "invalid min value exceeded - (1)",
					},
				},
			},
		},
		{
			name: "invalid max values request body",
			args: args{
				usr: UserRequest{
					UserId: string(make([]byte, 37)),
					Name:   string(make([]byte, 101)),
					Age:    1,
				},
			},
			exp: exp{
				errors: []errs.ValidationError{
					{
						Path:    "user_id",
						Message: "invalid UUID length: 37",
					},
					{
						Path:    "name",
						Message: "invalid length exceeded - (4-100)",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.exp.errors, tc.args.usr.Validate())
		})
	}
}
