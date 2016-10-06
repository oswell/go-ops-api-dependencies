package util

import (
    "net/http"
    "encoding/json"
)

type HttpStatus struct {
    Code    int    `json:code`
    Message string `json:message`
}

type HttpResponse struct {
    Status   HttpStatus  `json:status`
    Data   []interface{} `json:data,omitempty`
}

// Convert_Request_To_Object converts the client request object into an actual object.
// It returns an error if the request fails.
// Example: util.Convert_Request_To_Object(req, &ReqObject)
func Convert_Request_To_Object(req *http.Request, obj interface{})(error) {
    if req.ContentLength != 0 {
        decoder := json.NewDecoder(req.Body)
        err := decoder.Decode(obj) ; if err != nil {
            return err
        }
    }
    return nil
}

// Create_Http_Response generates an HTTP response body object.
// This meathod is meant to ensure a consistent response object to all calls.
// Returns a HttpResponse object.
func Create_Http_Response(status int, message string, objlist []interface{})(*HttpResponse) {
    httpResponse := &HttpResponse{}
    httpResponse.Status = HttpStatus{Code:status, Message:message}
    if objlist != nil {
        httpResponse.Data = objlist
    }

    // logger.Infof("Creating response, code: %d, message: %s", status, message)
    return httpResponse
}

// HttpRespond_with_content_type writes the response object to the http client.
// It sets the desired content type by writing out the Content-Type header.
func HttpRespond_with_content_type(w http.ResponseWriter, response *HttpResponse, content_type string)() {
    w.Header().Set("Content-Type", content_type)
    HttpRespond(w, response)
}

// HttpRespond_with_gzip writes the response object to the http client.
func HttpRespond(w http.ResponseWriter, response *HttpResponse)(error) {
    json_result, err := json.Marshal(response) ; if err != nil {
        return err
    }

    w.WriteHeader(response.Status.Code)
    w.Write(json_result)
    return nil
}
