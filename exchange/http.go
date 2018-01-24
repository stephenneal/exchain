package exchange

import (
    "encoding/json"
    "net/http"
    //"strconv"
    //"time"

    "github.com/romana/rlog"
)

//type JavaTime time.Time

func GetJson(url string, respType interface{}) error {
    return httpJsonResp("GET", url, respType)
}

func PostJson(url string, respType interface{}) error {
    return httpJsonResp("POST", url, respType)
}

func httpJsonResp(httpMethod string, url string, respType interface{}) error {
    rlog.Debug("URL: ", url)

    // Build the request
    req, err := http.NewRequest(httpMethod, url, nil)
    if err != nil {
        rlog.Critical("NewRequest: ", err)
        return err
    }

    // Create an HTTP Client for control over HTTP client headers, redirect policy, and other settings.
    client := &http.Client{}

    // Send the request via a client
    resp, err := client.Do(req)
    if err != nil {
        rlog.Critical("Do: ", err)
        return err
    }

    // Callers should close resp.Body when done reading from it
    // Defer the closing of the body
    defer resp.Body.Close()

    return decodeJson(resp, respType)
}

func decodeJson(resp *http.Response, respType interface{}) error {
    // Fill the response type with the data from the JSON
    // Use json.Decode for reading streams of JSON data
    return json.NewDecoder(resp.Body).Decode(&respType)
}

/*
func (j *JavaTime) UnmarshalJSON(data []byte) error {
    millis, err := strconv.ParseInt(string(data), 10, 64)
    if err != nil {
        return err
    }
    *j = JavaTime(time.Unix(0, millis * int64(time.Millisecond)))
    return nil
}

func (b JavaTime) String() string {
    return fmt.Sprintf("%b", b)
}
*/