package quotalib

import (
	"common/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"testing"
)

func CreatJson() {

	quota0 := &QuotaInfo{Role: "federation", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR"}}}
	q, _ := json.MarshalIndent(quota0, " ", "  ")
	err := ioutil.WriteFile("quota0.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %tv", err)
	}
	quota1 := &QuotaInfo{Role: "role1", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR"}}}
	q, _ = json.MarshalIndent(quota1, " ", "  ")
	err = ioutil.WriteFile("quota1.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %tv", err)
	}
	quota2 := &QuotaInfo{Role: "federation", Guarantee: []GuaranteeInfo{{"cpus", "*", ScalarInfo{2}, "SCALAR"}, {"mem", "*", ScalarInfo{2048}, "SCALAR,"}}}
	q, _ = json.MarshalIndent(quota2, " ", "  ")
	err = ioutil.WriteFile("quota2.json", q, 0644)
	if err != nil {
		fmt.Printf("WriteFile json Error: %tv", err)
	}

}

func Test_InitPlain(T *testing.T) {
	CreatJson()
}

//setquota with correct input
func Test_SetQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota(dc, "federation", "quota0.json")

	if e != nil {
		T.Error("error", e)
	}
}

//setquota with incorrect IP
func Test_SetQuota_Ip(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "54.187.130.7:5050"
	e := SetQuota(dc, "federation", "quota0.json")

	if e != nil {
		T.Error("error", e)
	}
}

//setquota with the file which is not exist
func Test_SetQuotaFile(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota(dc, "federation", "q.json")

	if e != nil {
		T.Error("error", e)
	}
}

//setquota with role which is not there in masterlist
func Test_SetQuotaRole(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota(dc, "role1", "quota1.json")

	if e != nil {
		T.Error("error", e)
	}
}

//setquota with incorrect json
func Test_SetQuotaJ(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota(dc, "federation", "quota2.json")

	if e != nil {
		T.Error("error", e)
	}

	removejson()

}

//getquota with correct input
func Test_GetQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	resp, e := GetQuota(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the response is  ", resp)

}

func Test_GetQuota_Ip(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.10:5050"
	resp, e := GetQuota(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the response is  ", resp)

}

//getquota with role which is not there in masterlist
func Test_GetQuotaRole(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	resp, e := GetQuota(dc, "role1")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the response is  ", resp)
}

//remaining resource with correct input
func Test_RemainingQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	cpu, mem, disk, e := RemainingResource(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the cpu,mem,quota  are ", cpu, mem, disk)
}

func Test_RemainingQuota_Ip(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.10:5050"
	cpu, mem, disk, e := RemainingResource(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the cpu,mem,quota  are ", cpu, mem, disk)
}

//remaining resource with role which is not there in masterlist
func Test_RemainingQuotaRole(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	cpu, mem, disk, e := RemainingResource(dc, "role1")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the cpu,mem,disk  are ", cpu, mem, disk)

}

//delete quota with correct input
func Test_DelQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := DelQuota(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
}

//delete quota with role which is not there in masterlist
func Test_DelQuotaRole(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := DelQuota(dc, "role1")

	if e != nil {
		T.Error("error", e)
	}
}

//delete quota with incorrect master_ip
func Test_DelQuota_Ip(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.10:5050"
	e := DelQuota(dc, "federation")

	if e != nil {
		T.Error("error", e)
	}
}

func removejson(){
	_, err := exec.Command("rm", "quota0.json", "quota1.json", "quota2.json").Output()
	if err != nil {
		fmt.Println(err)
	}
}
