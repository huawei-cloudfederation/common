package quotalib

import "testing"
import "common/types"

//setquota with correct input
func Test_SetQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota(dc, "federation", "quota.json")

	if e != nil {
		T.Error("error", e)
	}
}

//setquota with incorrect IP
func Test_SetQuota_Ip(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "54.187.130.7:5050"
	e := SetQuota(dc,"federation", "quota.json")

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
	e := SetQuota(dc, "federation", "quota_j.json")

	if e != nil {
		T.Error("error", e)
	}
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
	resp, e := GetQuota(dc,"federation")

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
	cpu, mem, disk, e := RemainingResource(dc,"federation")

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
