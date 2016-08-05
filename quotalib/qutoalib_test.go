package quotalib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	//	"os/exec"
	"testing"

	//Packages specific to this project
	"common/types"
)

func CreatJson() {

	quota0 := &QuotaInfo{Role: "federation", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR"}}}
	q, _ := json.MarshalIndent(quota0, " ", "  ")
	err := ioutil.WriteFile("quota0.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %v", err)
	}
	quota1 := &QuotaInfo{Role: "role1", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR"}}}
	q, _ = json.MarshalIndent(quota1, " ", "  ")
	err = ioutil.WriteFile("quota1.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %v", err)
	}
	quota2 := &QuotaInfo{Role: "federation", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR,"}}}
	q, _ = json.MarshalIndent(quota2, " ", "  ")
	err = ioutil.WriteFile("quota2.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %v", err)
	}

}

func Test_InitPlain(T *testing.T) {
	CreatJson()
}

//setquota with correct input
func TestSetQuotaValid(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))

	defer ts.Close()

	dc.Endpoint = ts.URL
	e := SetQuota(dc, "federation", "quota0.json")

	if e != nil {
		T.Fail()
	}
}

//SetQuota with Mesos Master not reachable
func TestSetQuotaBadMaster(T *testing.T) {
	var dc typ.DC

	dc.Endpoint = "http://10.11.12.13:5050"
	err := SetQuota(dc, "federation", "quota0.json")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "unreachable host") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestSetQuotaInvalidFilePath(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))

	defer ts.Close()

	dc.Endpoint = ts.URL
	err := SetQuota(dc, "federation", "random.json")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "The system cannot find the file") {
		//If its some other error then fail
		T.Fail()
	}
}

//Test for invalid json format
//We dont really have to supply an invalid json but we need to just simulate a error response from the server
func TestSetQuotaInvalidFileJsonformat(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Wrong json", http.StatusInternalServerError)
	}))

	defer ts.Close()

	dc.Endpoint = ts.URL
	err := SetQuota(dc, "federation", "quota0.json")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Server returned error respone") {
		//If its some other error then fail
		T.Fail()
	}

}

func TestDelQuotaValid(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))

	defer ts.Close()
	dc.Endpoint = ts.URL

	err := DelQuota(dc, "federation")

	if err != nil {
		//no error should occur
		T.Fail()
	}
}

func TestDelQuotaInvalidResponse(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Wrong quota name", http.StatusInternalServerError)
	}))
	defer ts.Close()

	dc.Endpoint = ts.URL

	err := DelQuota(dc, "federation")

	if err == nil {
		//Should return an error for invalid server response
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Server returned error respone") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestDelQuotaInvalidMaster(T *testing.T) {
	var dc typ.DC

	dc.Endpoint = "https://10.11.12.13:5050"
	err := DelQuota(dc, "federation")

	if err == nil {
		//Should return an error for invalid server response
		T.Fail()
	}

	if !strings.Contains(err.Error(), "unreachable host") {
		//If its some other error then fail
		T.Fail()
	}
}

func TestGetQuotaValid(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-OK")
	}))

	defer ts.Close()
	dc.Endpoint = ts.URL

	_, err := GetQuota(dc, "federation")

	if err != nil {
		//no error should occur
		T.Fail()
	}

}
