//Package usersrv ...
package usersrv

import (
	"testing"
)

func TestMockOauth2UserService_UpdateUser(t *testing.T) {
	var c MockOauth2UserService

	var us UserResponse
	us.Success = true
	c.MockUpdateUserResponse = &us

	s := c.GetNew()

	var u UpdateUser
	res := s.UpdateUser(u)
	if !res.Success {
		t.Fail()
	}
}

func TestMockOauth2UserService_GetUser(t *testing.T) {
	var c MockOauth2UserService
	var us User
	c.MockUser = &us
	c.MockUserCode = 200

	s := c.GetNew()
	res, cd := s.GetUser("test", "344")
	if res == nil && cd == 0 {
		t.Fail()
	}
}

func TestMockOauth2UserService_SetToken(t *testing.T) {
	var c MockOauth2UserService
	s := c.GetNew()
	s.SetToken("123")
}
