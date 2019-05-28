package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

var testUser = okta.User{
	FirstName:   "Test",
	LastName:    "Tester",
	Position:    "QA Tester",
	Permissions: []string{"action:read"},
}

type userService struct {
	mock.Mock
}

func (u *userService) GetUser(email string) (okta.User, error) {
	argsToReturn := u.Called(email)
	return argsToReturn.Get(0).(okta.User), argsToReturn.Error(1)
}

func (u *userService) AddPermissions(user *okta.User, service string) error {
	argsToReturn := u.Called(user, service)
	return argsToReturn.Error(0)
}

func createFakeRouter() (*httprouter.Router, *userService) {
	s := &userService{}

	router := httprouter.New()
	router.GET("/", getOktaUserByEmail(s))
	return router, s
}

func TestMissingQuery(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	router, s := createFakeRouter()

	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when entering wrong email")

	responseBody := response.Body.String()
	assert.Equal(t, "Missing email\n", responseBody, "Returns correct body")
	s.AssertNotCalled(t, "GetUser")
	s.AssertNotCalled(t, "AddPermissions")
}

func TestWrongEmail(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=testtest", nil)
	response := httptest.NewRecorder()
	router, s := createFakeRouter()
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when entering wrong email")

	responseBody := response.Body.String()
	assert.Equal(t, "Invalid email\n", responseBody, "Returns correct body")
	s.AssertNotCalled(t, "GetUser")
	s.AssertNotCalled(t, "AddPermissions")
}

func TestMissingUserAgent(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=test@test.com", nil)
	response := httptest.NewRecorder()
	router, s := createFakeRouter()
	router.ServeHTTP(response, request)
	assert.Equal(t, 400, response.Code, "Returns 400 when a user agent header is missing")

	responseBody := response.Body.String()
	assert.Equal(t, "Invalid user agent\n", responseBody, "Returns correct body")
	s.AssertNotCalled(t, "GetUser")
	s.AssertNotCalled(t, "AddPermissions")
}

func TestHappyPath(t *testing.T) {
	// Success response
	request, _ := http.NewRequest("GET", "/?email=test@test.com", nil)
	request.Header.Set("User-Agent", "testservice")
	response := httptest.NewRecorder()
	router, s := createFakeRouter()
	s.On("GetUser", "test@test.com").Return(testUser, nil)
	s.On("AddPermissions", &testUser, "TESTSERVICE").Return(nil)

	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "Returns 200 on success")
	jsonUser, _ := json.Marshal(testUser)
	responseJSON := response.Body.String()

	// For some reason response adds a extra line break
	assert.Equal(t, string(jsonUser)+"\n", responseJSON, "Returns correct body")
	s.AssertNumberOfCalls(t, "GetUser", 1)
	s.AssertNumberOfCalls(t, "AddPermissions", 1)
}

func TestControllerFailurePath(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=bs@test.com", nil)
	request.Header.Set("User-Agent", "testservice")
	response := httptest.NewRecorder()
	router, s := createFakeRouter()

	s.On("GetUser", "bs@test.com").Return(okta.User{}, errors.New("boom"))
	router.ServeHTTP(response, request)
	assert.Equal(t, 500, response.Code, "Returns 500 on controller failure")

	responseBody := response.Body.String()
	assert.Equal(t, "Service unavailable\n", responseBody, "Returns error correct body")
	s.AssertNumberOfCalls(t, "GetUser", 1)
	s.AssertNotCalled(t, "AddPermissions")
}

func TestNotFoundPath(t *testing.T) {
	request, _ := http.NewRequest("GET", "/?email=notfound@test.com", nil)
	request.Header.Set("User-Agent", "testservice")
	response := httptest.NewRecorder()
	router, s := createFakeRouter()

	s.On("GetUser", "notfound@test.com").Return(okta.User{}, okta.ErrUserNotFound)
	router.ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code, "Returns 404 on user not found")

	responseBody := response.Body.String()
	assert.Equal(t, "User notfound@test.com not found\n", responseBody, "Returns correct body")
	s.AssertNumberOfCalls(t, "GetUser", 1)
	s.AssertNotCalled(t, "AddPermissions")
}
