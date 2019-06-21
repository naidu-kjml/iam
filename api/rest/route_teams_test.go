package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTeamsGetter struct {
	mock.Mock
}

func (tg *mockTeamsGetter) GetTeams() (map[string]int, error) {
	args := tg.Called()
	return args.Get(0).(map[string]int), args.Error(1)
}

func TestInternalError(t *testing.T) {
	errMessage := "internal error that shouldn't be exposed"
	tg := &mockTeamsGetter{}
	tg.On("GetTeams").Return(map[string]int{}, errors.New(errMessage))

	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/", getTeams(tg))
	router.ServeHTTP(response, request)

	assert.Equal(t, 500, response.Code)
	assert.NotEqual(t, errMessage, response.Body.String())
	tg.AssertExpectations(t)
}

func TestGetTeams(t *testing.T) {
	tg := &mockTeamsGetter{}
	tg.On("GetTeams").Return(map[string]int{"team1": 3, "team2": 1, "team3": 1}, nil)

	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/", getTeams(tg))
	router.ServeHTTP(response, request)

	var expected = "{\"team1\":3,\"team2\":1,\"team3\":1}\n"
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, expected, response.Body.String())
	tg.AssertExpectations(t)
}
