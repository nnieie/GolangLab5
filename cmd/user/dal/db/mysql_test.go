// file: golanglab5/pkg/db/mysql_test.go
package db

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nnieie/golanglab5/pkg/errno"
)

// TestCreateUser 测试创建用户
func TestCreateUser(t *testing.T) {
	// 定义testCase
	type testCase struct {
		name           string
		mockError      error // 模拟 Create 方法返回的错误
		inputUser      *User // 要创建的用户
		expectedID     int64 // 期望返回的用户 ID
		expectingError bool  // 是否期望返回错误
	}

	testCases := []testCase{
		{
			name:           "CreateUser success",
			mockError:      nil,
			inputUser:      &User{UserName: "test"},
			expectedID:     1,
			expectingError: false,
		},
		{
			name:           "CreateUser db err",
			mockError:      errors.New("db err"),
			inputUser:      &User{UserName: "test"},
			expectedID:     0,
			expectingError: true,
		},
	}

	defer mockey.UnPatchAll()
	// 测试
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			// 创建一个 mockDB
			mockDB := new(gorm.DB)

			// (*gorm.DB).WithContext 返回 mockDB
			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockDB
			}).Build()

			// Mock (*gorm.DB).Create
			mockey.Mock((*gorm.DB).Create).To(func(value interface{}) *gorm.DB {
				// 如果测试用例需要模拟错误返回错误
				if tc.mockError != nil {
					mockDB.Error = tc.mockError
					return mockDB
				}
				// 如果成功，模拟 GORM 的行为：把 ID 写回传入的 user 对象
				user, ok := value.(*User)
				if ok {
					user.ID = uint(tc.expectedID)
				}
				mockDB.Error = nil
				return mockDB
			}).Build()

			// 调用要测试的函数 CreateUser
			id, err := CreateUser(context.Background(), tc.inputUser)

			// 检查结果
			if tc.expectingError {
				assert.Error(t, err)
				assert.Equal(t, tc.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

// TestQueryUserByName 测试根据用户名查询用户
func TestQueryUserByName(t *testing.T) {
	type testCase struct {
		name           string
		username       string
		mockResult     *User
		mockError      error
		expectedResult *User
		expectingError bool
		expectedErrno  error
	}

	testCases := []testCase{
		{
			name:           "QueryUserByName success",
			username:       "user",
			mockResult:     &User{UserName: "user", Password: "password"},
			mockError:      nil,
			expectedResult: &User{UserName: "user", Password: ""},
			expectingError: false,
		},
		{
			name:           "QueryUserByName RecordNotFound",
			username:       "uuu",
			mockResult:     nil,
			mockError:      gorm.ErrRecordNotFound,
			expectedResult: nil,
			expectingError: true,
			expectedErrno:  errno.UserIsNotExistErr,
		},
		{
			name:           "QueryUserByName db err",
			username:       "any_user",
			mockResult:     nil,
			mockError:      errors.New("db err"),
			expectedResult: nil,
			expectingError: true,
			expectedErrno:  errors.New("db err"),
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockDB := new(gorm.DB)

			// WithContext 和 Where 方法返回 mockDB 本身
			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB { return mockDB }).Build()
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB { return mockDB }).Build()

			// First 返回 预定的结果
			mockey.Mock((*gorm.DB).First).To(func(dest interface{}, conds ...interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockDB.Error = tc.mockError
					return mockDB
				}
				// 模拟 GORM 找到了数据，并把数据写入到 dest
				user, ok := dest.(*User)
				if ok && tc.mockResult != nil {
					*user = *tc.mockResult
				}
				mockDB.Error = nil
				return mockDB
			}).Build()

			// 调用被测试函数
			user, err := QueryUserByName(context.Background(), tc.username)

			// 检查结果
			if tc.expectingError {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErrno, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, user)
			}
		})
	}
}

func TestUpdateAvatar(t *testing.T) {
	type testCase struct {
		name            string
		userID          int64
		avatarURL       string
		mockError       error
		expectingResult *User
		expectingError  bool
	}

	testCases := []testCase{
		{
			name:            "UpdateAvatar success",
			userID:          1,
			avatarURL:       "http://example.com/avatar.png",
			mockError:       nil,
			expectingResult: &User{UserName: "test", Avatar: "http://example.com/avatar.png", Model: gorm.Model{ID: 1}},
			expectingError:  false,
		},
		{
			name:            "UpdateAvatar db update err",
			userID:          1,
			avatarURL:       "http://example.com/avatar.png",
			mockError:       errors.New("db update err"),
			expectingResult: nil,
			expectingError:  true,
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockDB := new(gorm.DB)
			var capturedUser *User

			// Mock WithContext - 返回 mockDB
			mockey.Mock((*gorm.DB).WithContext).To(func(ctx context.Context) *gorm.DB {
				return mockDB
			}).Build()

			// Mock Model - 捕获传入的 user 指针并返回 mockDB
			mockey.Mock((*gorm.DB).Model).To(func(value interface{}) *gorm.DB {
				// 捕获传入的 user 指针
				if user, ok := value.(*User); ok {
					capturedUser = user
				}
				return mockDB
			}).Build()

			// Mock Clauses - 返回 mockDB
			mockey.Mock((*gorm.DB).Clauses).To(func(conds ...clause.Expression) *gorm.DB {
				return mockDB
			}).Build()

			// Mock Where - 返回 mockDB
			mockey.Mock((*gorm.DB).Where).To(func(query interface{}, args ...interface{}) *gorm.DB {
				return mockDB
			}).Build()

			// Mock Update - 模拟 Returning 效果
			mockey.Mock((*gorm.DB).Update).To(func(column string, value interface{}) *gorm.DB {
				if tc.mockError != nil {
					mockDB.Error = tc.mockError
					return mockDB
				}

				// 模拟 Returning 效果：把预期结果写入 capturedUser
				*capturedUser = *tc.expectingResult
				return mockDB
			}).Build()

			// 调用被测试函数
			user, err := UpdateAvatar(context.Background(), tc.userID, tc.avatarURL)

			// 检查结果
			if tc.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectingResult.ID, user.ID)
				assert.Equal(t, tc.expectingResult.UserName, user.UserName)
				assert.Equal(t, tc.expectingResult.Avatar, user.Avatar)
			}
		})
	}
}
