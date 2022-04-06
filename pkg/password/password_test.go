package password

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func newStringPtr(s string) *string {
	return &s
}

var examplePassword = "MuDz9QX2cU6e67⌘bdbd\\\\xb2=\\\\B2jKGiJJrix.TjLG.GVupBizaBmV*wre_mkGCgN7rqRg!njsDqcvJsF9UsNW8bKPvpmvc7VCMz3Aofbo2yp*"

func benchmarkHashing(hashingParams HashingParams, b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = HashPassword(examplePassword, hashingParams)
	}
}

func BenchmarkDefaultHashingParams(b *testing.B) {
	benchmarkHashing(GetHashingParams(), b)
}

func TestHashPassword(t *testing.T) {
	type expected struct {
		hashingPasswordErr   error
		verifyingPasswordErr error
		valid                bool
	}

	type testCase struct {
		name                 string
		password             string
		verificationPassword string
		hashingParams        HashingParams
		expected             expected
	}

	testCases := []testCase{
		{
			name:                 "success",
			password:             examplePassword,
			verificationPassword: examplePassword,
			hashingParams:        GetHashingParams(),
			expected: expected{
				valid: true,
			},
		},
		{
			name:          "nil password",
			hashingParams: GetHashingParams(),
			expected: expected{
				hashingPasswordErr:   ErrEmptyPassword,
				verifyingPasswordErr: ErrInvalidHash,
			},
		},
		{
			name:                 "invalid password",
			hashingParams:        GetHashingParams(),
			password:             "the_first_password",
			verificationPassword: "wrong_password_doesnt_match_the_first_one",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tC.password, tC.hashingParams)
			assert.Equal(t, tC.expected.hashingPasswordErr, err)

			valid, err := VerifyPassword(tC.verificationPassword, hashedPassword)
			assert.Equal(t, tC.expected.verifyingPasswordErr, err)
			assert.Equal(t, tC.expected.valid, valid)

		})
	}
}
