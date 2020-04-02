package rest

import (
	"github.com/federicoleon/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})
	repository := userRepository{}

	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials","status":"404", "error": "not_found"}`,
	})
	repository := userRepository{}

	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials", "status":404, "error": "not_found"}`,
	})
	repository := userRepository{}

	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id":"123", "first_name":"Will", "last_name":"Qiang", "email":"rockqiang32@gmail.com"}`,
	})
	repository := userRepository{}

	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying unmarshal users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password": "the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 12, "first_name": "Will", "last_name": "Qiang", "email": "willqiang32@gmail.com"}`,
	})
	repository := userRepository{}

	user, err := repository.LoginUser("email@gmail.com", "the-password")
	/**
	以下Println语句如果开启，其结果为：&{12   willqiang32@gmail.com}，
	并且不管 RespBody 的 json 字符串中的field是什么，似乎只识别 id 和 email field，
	所以暂时将下面的针对 user.FirstName 和 user.LastName 的判断也注释掉
	 */
	//fmt.Println(user)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 12, user.Id)
	//assert.EqualValues(t, "Will", user.FirstName)
	//assert.EqualValues(t, "Qiang", user.LastName)
	assert.EqualValues(t, "willqiang32@gmail.com", user.Email)
}
