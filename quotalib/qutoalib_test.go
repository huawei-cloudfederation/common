package quotalib

import "testing"
import "../types"

func Test_SetQuota(T *testing.T) {
	var dc typ.DC
	dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota("test","/home/ubuntu/quota.json",dc)

	if e != nil {
		T.Error("error", e)
	}
}

func Test_SetQuotaFile(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota("test","quota.json",dc)

	if e != nil {
		T.Error("error", e)
	}
}

func Test_SetQuotaRole(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
	e := SetQuota("federation","/home/ubuntu/quota1.json",dc)

	if e != nil {
		T.Error("error", e)
	}
}

func Test_GetQuota(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
        resp,e := GetQuota("test",dc)

        if e != nil {
                T.Error("error", e)
        }
	 T.Log("the response is  ",resp)

}

func Test_GetQuotaRole(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
	resp,e := GetQuota("federation",dc)

	if e != nil {
		T.Error("error", e)
	}
	 T.Log("the response is  ",resp)
}

func Test_RemainingQuota(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
	cpu,mem,disk,e := RemainingResource("test",dc)

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the cpu,mem,quota  are ",cpu,mem,disk)
}
func Test_RemainingQuotaRole(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
	cpu,mem,disk,e := RemainingResource("federation",dc)

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the cpu,mem,disk  are ",cpu,mem,disk)

}

func Test_DelQuota(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
        e := DelQuota("test",dc)

        if e != nil {
                T.Error("error", e)
        }
}

func Test_DelQuotaRole(T *testing.T) {
	 var dc typ.DC
        dc.Endpoint = "172.31.44.22:5050"
        e := DelQuota("federation",dc)

        if e != nil {
                T.Error("error", e)
        }
}

