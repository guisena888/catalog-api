package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/mytheresa/go-hiring-challenge/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ResponseSuite struct {
	suite.Suite
}

func TestResponseSuite(t *testing.T) {
	suite.Run(t, new(ResponseSuite))
}

func (s *ResponseSuite) TestOKResponse() {
	type sampleResponse struct {
		Message string `json:"message"`
	}

	recorder := httptest.NewRecorder()
	OKResponse(recorder, sampleResponse{Message: "Success"})

	s.Equal(http.StatusOK, recorder.Code)
	s.Equal("application/json", recorder.Header().Get("Content-Type"))
	s.JSONEq(`{"message":"Success"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestCreatedResponse() {
	type sampleResponse struct {
		Name string `json:"name"`
	}

	recorder := httptest.NewRecorder()
	CreatedResponse(recorder, sampleResponse{Name: "Hats"})

	s.Equal(http.StatusCreated, recorder.Code)
	s.Equal("application/json", recorder.Header().Get("Content-Type"))
	s.JSONEq(`{"name":"Hats"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestErrorResponse() {
	recorder := httptest.NewRecorder()
	ErrorResponse(recorder, http.StatusInternalServerError, "Some error occurred")

	s.Equal(http.StatusInternalServerError, recorder.Code)
	s.Equal("application/json", recorder.Header().Get("Content-Type"))
	s.JSONEq(`{"error":"Some error occurred"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestHandleError_InvalidInput() {
	recorder := httptest.NewRecorder()
	HandleError(recorder, apperrors.ErrInvalidLimit)

	s.Equal(http.StatusBadRequest, recorder.Code)
	s.JSONEq(`{"error":"invalid limit"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestHandleError_NotFound() {
	recorder := httptest.NewRecorder()
	HandleError(recorder, apperrors.ErrProductNotFound)

	s.Equal(http.StatusNotFound, recorder.Code)
	s.JSONEq(`{"error":"product not found"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestHandleError_AlreadyExists() {
	recorder := httptest.NewRecorder()
	HandleError(recorder, apperrors.ErrCategoryAlreadyExists)

	s.Equal(http.StatusConflict, recorder.Code)
	s.JSONEq(`{"error":"category already exists"}`, recorder.Body.String())
}

func (s *ResponseSuite) TestHandleError_UnknownError() {
	recorder := httptest.NewRecorder()
	HandleError(recorder, errors.New("unexpected failure"))

	s.Equal(http.StatusInternalServerError, recorder.Code)
	s.JSONEq(`{"error":"internal server error"}`, recorder.Body.String())
}

func TestOKResponse(t *testing.T) {
	type sampleResponse struct {
		Message string `json:"message"`
	}

	sample := sampleResponse{Message: "Success"}

	t.Run("succesful http200 json response", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		OKResponse(recorder, sample)

		assert.Equal(t, http.StatusOK, recorder.Code, "Expected status code 200 OK")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"message":"Success"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}

func TestErrorResponse(t *testing.T) {
	t.Run("json response for a given http status code", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		ErrorResponse(recorder, http.StatusInternalServerError, "Some error occurred")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code, "Expected status code 500 Internal Server Error")
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"), "Expected Content-Type to be application/json")

		expected := `{"error":"Some error occurred"}`
		assert.JSONEq(t, expected, recorder.Body.String(), "Response body does not match expected")
	})
}
