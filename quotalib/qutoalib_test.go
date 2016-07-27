package quotalib

import "testing"

func Test_SetQuota(T *testing.T) {
	e := SetQuota("test","/home/ubuntu/quota.json")

	if e != nil {
		T.Error("error", e)
	}
}

func Test_SetQuotaFile(T *testing.T) {
	e := SetQuota("test","quota.json")

	if e != nil {
		T.Error("error", e)
	}
}

func Test_SetQuotaRole(T *testing.T) {
	e := SetQuota("federation","/home/ubuntu/quota1.json")

	if e != nil {
		T.Error("error", e)
	}
}

func Test_GetQuota(T *testing.T) {
        resp,e := GetQuota("test")

        if e != nil {
                T.Error("error", e)
        }

	T.Log("the json is ",resp)
}

func Test_GetQuotaRole(T *testing.T) {
	resp,e := GetQuota("federation")

	if e != nil {
		T.Error("error", e)
	}
	T.Log("the json is ",resp)
}

func Test_DelQuota(T *testing.T) {
        e := DelQuota("test")

        if e != nil {
                T.Error("error", e)
        }
}

func Test_DelQuotaRole(T *testing.T) {
        e := DelQuota("federation")

        if e != nil {
                T.Error("error", e)
        }
}

