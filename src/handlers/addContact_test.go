package main

import (
//    "fmt"
    "os"
    "context"
//    "io/ioutil"
//    "io"
//    "net/http"
//    "net/http/httptest"
    "testing"
    
    "gotest.tools/assert"
    "github.com/aws/aws-lambda-go/events"
//    is "gotest.tools/assert/cmp"
)

func TestAddContact(t *testing.T){
    /*tests := []struct{
        request events.APIGatewayProxyRequest
        response string
        err error
    }{
        {
        request: events.APIGatewayProxyRequest{Body: "Paul"},
        expect: "Hello Paul",
        err: nil,
        },
        {
        request: events.APIGatewayProxyRequest{Body: ""},
        expect: "",
        err: nil,
        },
    }
    for _, test := range tests {
        response, err := test.request
        assert.IsType(t, test.err, err)
        assert.Equal(t, test.expect, response.Body)
    }*/
    /*err := os.Open("contact-form-crud/serverless.yml")
    for _, test := range err{
        if test != nil {
            assert.Error(t, test, "No file found")
        }
    }*/
    os.Getenv("FUNCTION")
    handler := func(ctx context.Context, request events.APIGatewayProxyRequest){
        events.APIGatewayProxyResponse{StatusCode: 200}
    }
    
//    req := 
    /*handler := func(w http.ResponseWriter, r *http.Request) {
        io.WriteString(w, "{\"status\": \"good\"}")
    }
    req := httptest.NewRequest("POST","https://gfq99349qg.execute-api.ap-southeast-1.amazonaws.com/dev/contacts",nil)
    w := httptest.NewRecorder()
    handler(w, req)
    
    resp := w.Result()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println(string(body))
    if 200 != resp.StatusCode {
        t.Fatal("Status Code Not OK")
    }*/
//    assert.Equal(t, resp.StatusCode, "The response is not 200")
}