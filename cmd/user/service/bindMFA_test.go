package service

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/nnieie/golanglab5/cmd/user/dal/cache"
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/utils"
)

func TestUserService_BindMFA(t *testing.T) {
	type testCase struct {
		name   string
		code   string
		secret string
		userID int64

		mockCacheResult   string
		mockCacheError    error
		mockDecryptResult string
		mockDecryptError  error
		mockCheckResult   bool
		mockDBError       error

		expectedOK     bool
		expectingError bool
		expectedErrno  error
	}

	testCases := []testCase{
		{
			name:              "bind mfa success",
			code:              "114514",
			secret:            "secret",
			userID:            1,
			mockCacheResult:   "encryptedSecret",
			mockCacheError:    nil,
			mockDecryptResult: "secret",
			mockDecryptError:  nil,
			mockCheckResult:   true,
			mockDBError:       nil,
			expectedOK:        true,
			expectingError:    false,
		},
		{
			name:            "no totp in cache",
			userID:          1,
			mockCacheResult: "",
			mockCacheError:  nil,
			expectedOK:      false,
			expectingError:  true,
			expectedErrno:   errno.NotGenerateTotpErr,
		},
		{
			name:              "secret not match",
			secret:            "wrong",
			userID:            1,
			mockCacheResult:   "encryptedSecret",
			mockCacheError:    nil,
			mockDecryptResult: "secret",
			mockDecryptError:  nil,
			expectedOK:        false,
			expectingError:    true,
			expectedErrno:     errno.MFAInvalidCodeErr,
		},
		{
			name:              "invalid code",
			code:              "114514",
			secret:            "secret",
			userID:            1,
			mockCacheResult:   "encryptedSecret",
			mockCacheError:    nil,
			mockDecryptResult: "secret",
			mockDecryptError:  nil,
			mockCheckResult:   false,
			expectedOK:        false,
			expectingError:    true,
			expectedErrno:     errno.MFAInvalidCodeErr,
		},
		{
			name:              "db err",
			code:              "114514",
			secret:            "secret",
			userID:            1,
			mockCacheResult:   "encryptedSecret",
			mockCacheError:    nil,
			mockDecryptResult: "secret",
			mockDecryptError:  nil,
			mockCheckResult:   true,
			mockDBError:       errors.New("db err"),
			expectedOK:        false,
			expectingError:    true,
			expectedErrno:     errors.New("db err"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(cache.GetTOTPSecret).To(func(ctx context.Context, userID int64) (string, error) {
				return tc.mockCacheResult, tc.mockCacheError
			}).Build()

			mockey.Mock(utils.Decrypt).To(func(encoded string) (string, error) {
				return tc.mockDecryptResult, tc.mockDecryptError
			}).Build()

			mockey.Mock(utils.CheckTotp).To(func(code, secret string) bool {
				return tc.mockCheckResult
			}).Build()

			mockey.Mock(db.UpdateMFA).To(func(ctx context.Context, secret string, userID int64) error {
				return tc.mockDBError
			}).Build()

			s := &UserService{ctx: context.Background()}

			ok, err := s.BindMFA(tc.code, tc.secret, tc.userID)

			assert.Equal(t, tc.expectedOK, ok)

			if tc.expectingError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrno, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
