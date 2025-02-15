package response

import (
	"encoding/json"
	"net/http"
)

func bytes(resp interface{}) []byte {
	data, _ := json.Marshal(resp)
	return data

}

func sendResponse(w http.ResponseWriter, r *http.Request, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	if data != nil {
		_, _ = w.Write(bytes(data))
	}

}

// 200
func StatusOk(w http.ResponseWriter, r *http.Request, data interface{}) {
	sendResponse(w, r, http.StatusOK, data)
}

// 201
func StatusCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	sendResponse(w, r, http.StatusCreated, nil)
}

// 204
func StatusNoContent(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, http.StatusNoContent, nil)

}

// 400
func StatusBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusBadRequest, data)
}

// 401
func StatusNotAuthorized(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusUnauthorized, data)
}

// 404
func StatusNotFound(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusNotFound, data)
}

// 405
func StatusMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, r, http.StatusMethodNotAllowed, nil)
}

// 404
func StatusConflict(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusConflict, data)
}

// 500
func StatusInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusInternalServerError, data)
}

// 422
func StatusUnprocessableEntity(w http.ResponseWriter, r *http.Request, err error) {
	data := map[string]interface{}{"error": err.Error()}
	sendResponse(w, r, http.StatusUnprocessableEntity, data)
}
