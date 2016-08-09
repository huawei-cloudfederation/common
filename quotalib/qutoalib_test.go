package quotalib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//Packages specific to this project
	"common/types"
)

func CreatJson() {

	quota0 := &QuotaInfo{Role: "federation", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR"}, {"disk", "*", ScalarInfo{2.9}, "SCALAR"}}}
	q, _ := json.MarshalIndent(quota0, " ", "  ")
	err := ioutil.WriteFile("quota0.json", q, 0644)
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

	if !strings.Contains(err.Error(), "Post http://10.11.12.13:5050/quota/federation: dial tcp 10.11.12.13:5050: i/o timeout") {
		//If its some other error then fail
		T.Fail()
	}
}

//SetQuota with invalid file
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

	if !strings.Contains(err.Error(), "open random.json: no such file or directory") {
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
	err := SetQuota(dc, "federation", "quota2.json")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Server returned error respone") {
		//If its some other error then fail
		T.Fail()
	}

}

//GetQuota with valid input
func TestGetQuotaValid(T *testing.T) {
	var dc typ.DC

	quota := `{"infos":[{"guarantee":[{"name":"cpus","role":"*","scalar":{"value":2.0},"type":"SCALAR"},{"name":"mem","role":"*","scalar":{"value":2048.0},"type":"SCALAR"}],"role":"federation"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, quota)
	}))

	defer ts.Close()
	dc.Endpoint = ts.URL

	_, err := GetQuota(dc, "federation")

	if err != nil {
		//no error should occur
		T.Fail()
	}

}

//getquota with mesos master unreachable
func TestGetQuotaBadMaster(T *testing.T) {
	var dc typ.DC

	dc.Endpoint = "http://10.11.12.13:5050"

	_, err := GetQuota(dc, "federation")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Get http://10.11.12.13:5050/quota/federation: dial tcp 10.11.12.13:5050: i/o timeout") {
		//If its some other error then fail
		T.Fail()
	}

}

//Test for invalid json format
//We dont really have to supply an invalid json but we need to just simulate a error response from the server
func TestGetQuotaInvalidFileJsonformat(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "wrong json", http.StatusInternalServerError)
	}))

	defer ts.Close()

	dc.Endpoint = ts.URL
	_, err := GetQuota(dc, "federation")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "invalid character 'w' looking for beginning of value") {
		//If its some other error then fail
		T.Fail()
	}

}

//RemainingResource with valid quota input
func TestRemainingResourceValid(T *testing.T) {
	var dc typ.DC

	quota := `{"infos":[{"guarantee":[{"name":"cpus","role":"*","scalar":{"value":2.0},"type":"SCALAR"},{"name":"mem","role":"*","scalar":{"value":2048.0},"type":"SCALAR"}],"role":"federation"}]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, quota)
	}))

	defer ts.Close()
	dc.Endpoint = ts.URL

	cpu, mem, disk, err := RemainingResource(dc, "federation")

	if err != nil {
		//no error should occur
		T.Fail()
	}
	fmt.Println(cpu, mem, disk)

}

//RemainingResource with valid state json input
func TestRemainingResourceValidStateJson(T *testing.T) {
	var dc typ.DC

	var flag int

	quota := `{"infos":[{"guarantee":[{"name":"cpus","role":"*","scalar":{"value":2.0},"type":"SCALAR"},{"name":"mem","role":"*","scalar":{"value":2048.0},"type":"SCALAR"}],"role":"federation"}]}`

	state := `{"version":"1.1.0","git_sha":"77ddbb62dd2ab4faaa22de8355f4766e7bbe0f2d","git_branch":"refs\/heads\/master","build_date":"2016-07-19 13:25:21","build_time":1468934721.0,"build_user":"root","start_time":1470646716.76334,"elected_time":1470646716.7739,"id":"69364140-fbfc-4c2b-b532-415123e6ea7b","pid":"master@172.31.44.22:5050","hostname":"ip-172-31-44-22","activated_slaves":2.0,"deactivated_slaves":0.0,"leader":"master@172.31.44.22:5050","log_dir":"\/home\/ubuntu\/demo\/masterLogs","flags":{"agent_ping_timeout":"15secs","agent_reregister_timeout":"10mins","allocation_interval":"1secs","allocator":"HierarchicalDRF","authenticate_agents":"false","authenticate_frameworks":"false","authenticate_http":"false","authenticate_http_frameworks":"false","authenticators":"crammd5","authorizers":"local","framework_sorter":"drf","help":"false","hostname_lookup":"true","http_authenticators":"basic","initialize_driver_logging":"true","ip":"172.31.44.22","log_auto_initialize":"true","log_dir":"\/home\/ubuntu\/demo\/masterLogs","logbufsecs":"0","logging_level":"INFO","max_agent_ping_timeouts":"5","max_completed_frameworks":"50","max_completed_tasks_per_framework":"1000","port":"5050","quiet":"false","recovery_agent_removal_limit":"100%","registry":"replicated_log","registry_fetch_timeout":"1mins","registry_store_timeout":"20secs","registry_strict":"false","roles":"federation","root_submissions":"true","user_sorter":"drf","version":"false","webui_dir":"\/home\/ubuntu\/mesos\/build\/..\/src\/webui","work_dir":"\/home\/ubuntu","zk_session_timeout":"10secs"},"slaves":[{"id":"44f2627d-59ed-4cae-97c7-e2464933c903-S0","pid":"slave(1)@172.31.7.100:5051","hostname":"ip-172-31-7-100.us-west-2.compute.internal","registered_time":1470646722.80892,"reregistered_time":1470646726.15501,"resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"},"used_resources":{"disk":0.0,"mem":1024.0,"gpus":0.0,"cpus":1.0},"offered_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"reserved_resources":{},"unreserved_resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"},"attributes":{},"active":true,"version":"1.1.0"},{"id":"44f2627d-59ed-4cae-97c7-e2464933c903-S1","pid":"slave(1)@172.31.44.22:5051","hostname":"ip-172-31-44-22","registered_time":1470646721.9035,"reregistered_time":1470646722.72788,"resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"},"used_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"offered_resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"},"reserved_resources":{},"unreserved_resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"},"attributes":{},"active":true,"version":"1.1.0"}],"frameworks":[{"id":"1222-5555-999-300-18","name":"Test Framework (Go)","pid":"scheduler(1)@54.149.214.54:54970","used_resources":{"disk":0.0,"mem":1024.0,"gpus":0.0,"cpus":1.0},"offered_resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0},"capabilities":[],"hostname":"ip-172-31-44-22","webui_url":"","active":true,"user":"ubuntu","failover_timeout":0.0,"checkpoint":false,"role":"federation","registered_time":1470653632.26932,"unregistered_time":0.0,"resources":{"disk":14896.0,"mem":3952.0,"gpus":0.0,"cpus":3.0},"tasks":[{"id":"1","name":"go-task-1","framework_id":"1222-5555-999-300-18","executor_id":"default","slave_id":"44f2627d-59ed-4cae-97c7-e2464933c903-S0","state":"TASK_STARTING","resources":{"disk":0.0,"mem":1024.0,"gpus":0.0,"cpus":1.0},"statuses":[{"state":"TASK_RUNNING","timestamp":1470653611.0,"container_status":{"network_infos":[{"ip_addresses":[{"ip_address":"172.31.7.100"}]}]}},{"state":"TASK_STARTING","timestamp":1470653611.0,"container_status":{"network_infos":[{"ip_addresses":[{"ip_address":"172.31.7.100"}]}]}}]}],"completed_tasks":[],"offers":[{"id":"69364140-fbfc-4c2b-b532-415123e6ea7b-O3","framework_id":"1222-5555-999-300-18","slave_id":"44f2627d-59ed-4cae-97c7-e2464933c903-S1","resources":{"disk":14896.0,"mem":2928.0,"gpus":0.0,"cpus":2.0,"ports":"[31000-32000]"}}],"executors":[{"executor_id":"default","name":"Test Executor (Go)","framework_id":"1222-5555-999-300-18","command":{"value":".\/executor -logtostderr=true -v=0 -slow_tasks=true","argv":[],"uris":[{"value":"http:\/\/54.149.214.54:12345\/executor","executable":true}]},"resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"slave_id":"44f2627d-59ed-4cae-97c7-e2464933c903-S0"}]}],"completed_frameworks":[{"id":"1222-5555-999-300-17","name":"Test Framework (Go)","pid":"scheduler(1)@54.149.214.54:54582","used_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"offered_resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"capabilities":[],"hostname":"ip-172-31-44-22","webui_url":"","active":false,"user":"ubuntu","failover_timeout":0.0,"checkpoint":false,"role":"federation","registered_time":1470653571.10144,"unregistered_time":1470653629.24112,"resources":{"disk":0.0,"mem":0.0,"gpus":0.0,"cpus":0.0},"tasks":[],"completed_tasks":[],"offers":[],"executors":[]}],"orphan_tasks":[],"unregistered_frameworks":["1222-5555-999-300-17","1222-5555-999-300-17"]}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if flag == 0 {
			fmt.Fprintln(w, quota)
			flag++
		} else {
			fmt.Fprintln(w, state)
		}
	}))

	defer ts.Close()
	dc.Endpoint = ts.URL

	cpu, mem, disk, err := RemainingResource(dc, "federation")

	if err != nil {
		//no error should occur
		T.Fail()
	}
	fmt.Println(cpu, mem, disk)

}

//RemainingResource  with bad mesos master
func TestRemainingResourceBadMaster(T *testing.T) {
	var dc typ.DC

	dc.Endpoint = "http://10.11.12.13:5050"

	_, _, _, err := RemainingResource(dc, "federation")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Get http://10.11.12.13:5050/quota/federation: dial tcp 10.11.12.13:5050: i/o timeout") {
		//If its some other error then fail
		T.Fail()
	}

}

//Test for invalid json format
//We dont really have to supply an invalid json but we need to just simulate a error response from the server
func TestRemainingResourceInvalidJsonformat(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "wrong json", http.StatusInternalServerError)
	}))

	defer ts.Close()

	dc.Endpoint = ts.URL
	_, _, _, err := RemainingResource(dc, "federation")

	if err == nil {
		//Error cannot be nil
		T.Fail()
	}

	if !strings.Contains(err.Error(), "invalid character 'w' looking for beginning of value") {
		//If its some other error then fail
		T.Fail()
	}

}

//delete quota with valid input
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

//delete quota with wrong quota name
func TestDelQuotaInvalidResponse(T *testing.T) {
	var dc typ.DC

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Wrong quota name", http.StatusInternalServerError)
	}))
	defer ts.Close()

	dc.Endpoint = ts.URL

	err := DelQuota(dc, "test")

	if err == nil {
		//Should return an error for invalid server response
		T.Fail()
	}

	if !strings.Contains(err.Error(), "Server returned error respone") {
		//If its some other error then fail
		T.Fail()
	}
}

//delete quota with mesos master not reachable
func TestDelQuotaInvalidMaster(T *testing.T) {
	var dc typ.DC

	dc.Endpoint = "https://10.11.12.13:5050"
	err := DelQuota(dc, "federation")

	if err == nil {
		//Should return an error for invalid server response
		T.Fail()
	}

	if !strings.Contains(err.Error(), "dial tcp 10.11.12.13:5050: i/o timeout") {
		//If its some other error then fail
		T.Fail()
	}
}
