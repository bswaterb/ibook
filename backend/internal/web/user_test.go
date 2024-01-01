package web

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"ibook/internal/service"
	usvcmocks "ibook/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserSignup(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "测试1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			usersvc := usvcmocks.NewMockUserService(ctrl)
			h := NewUserHandler(usersvc)
			h.RegisterRoutesV1(server)

			usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&service.User{
				Id:          0,
				Email:       "781201402@qq.com",
				PhoneNumber: "",
				PassWord:    "123@123#abc",
			}, nil)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(`
{
	"email": "781201402@qq.com",
	"password": "123@123#abc",
	"confirmPassword": "123@123#abc"
}
`)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, http.StatusOK, resp.Code)
			// assert.Equal(t, nil, resp.Body.String())
		})
	}
}
