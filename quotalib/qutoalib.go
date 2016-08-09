//Package quotalib will implement all the functions that should be used to contact masters

package quotalib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//Packages of this project
	"common/types"
)

//Declare the structure which will give quota information
type QInfos struct {
	Infos []QuotaInfo `json:"infos"`
}

type QuotaInfo struct {
	Role      string          `json:"role"`      //quota is applied with this role
	Guarantee []GuaranteeInfo `json:"guarantee"` //allocate guaranteed quota to a role associted

}

type GuaranteeInfo struct {
	Name   string     `json:"name"` // name of the resource, ex:cpu
	Role   string     `json:"role"`
	Scalar ScalarInfo `json:"scalar"` //scalar resources info
	Type   string     `json:"type"`   //type of resource ex:scalar
}

type ScalarInfo struct {
	Value float64 `json:"value"` //value of resource,ex:cpu='2'
}

//SetQuota will read a json file the local disk and performe a SET Quota HTTP api call
//role : Simply the role for whcih the Quota to be set
//inputPath : Path from where quota json file should be read
func SetQuota(dc typ.DC, role string, inputPath string) error {
	//Function implementation

	buf, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Printf("Unable to read file = %v", err)
		return err
	}

	body := bytes.NewBuffer(buf)

	resp, err := http.Post(fmt.Sprintf("%s/quota/%s", dc.Endpoint, role), "text/json", body)
	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned error %v", resp)
		return fmt.Errorf("Server returned error respone %v", resp)
	}

	//All okay the Quota has been set

	return nil
}

//DelQuota Will delete the Quota set on the given role
//role : Simply the role name , usually Federation
func DelQuota(dc typ.DC, role string) error {

	//Function implementation

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/quota/%s", dc.Endpoint, role), nil)

	// handle err
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned an error %v", resp)
		return fmt.Errorf("Server returned error respone %v", resp)
	}

	return nil
}

//GetQuota Will get a json representing Quota for the given role
//role : Simply the role name , usually Federation
func GetQuota(dc typ.DC, role string) (*QuotaInfo, error) {
	var data QInfos

	//Function implementation
	resp, err := http.Get(fmt.Sprintf("%s/quota/%s", dc.Endpoint, role))

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	err = json.Unmarshal([]byte(body), &data)

	if err != nil {
		log.Printf("Json Unmarshall error = %v", err)
		return nil, err
	}

	for _, quota := range data.Infos {
		if quota.Role == role {
			return &quota, nil
		}
	}

	//None found matching our role name
	return nil, fmt.Errorf("No Quota found for role=%s", role)
}

//RemainingResources
//role: Role name, usually Federation
//This should automatically resolve the Masters Endpoint from the configuration
//Retruns how much of the gaurentee has been satisfied for a given role
//If a role has a quota of 50 CPUs and out of that 35 have been used
func RemainingResource(dc typ.DC, role string) (float64, float64, float64, error) {

	var CPU, MEM, DISK float64
	var uCPU, uMEM, uDISK float64
	var qCPU, qMEM, qDISK float64

	var mState typ.MasterState

	//Step 1: This should first call GetQuota() and interpret gaurantee
	quotaInfo, err := GetQuota(dc, role)
	if err != nil {
		log.Printf("Error: GetQuota(%s, %s)", dc.Endpoint, role)
		return 0.0, 0.0, 0.0, err
	}

	log.Println("the quota is ", quotaInfo)

	for _, g := range quotaInfo.Guarantee {
		if g.Name == "cpus" {
			qCPU = g.Scalar.Value
		}
		if g.Name == "mem" {
			qMEM = g.Scalar.Value
		}
		if g.Name == "disk" {
			qDISK = g.Scalar.Value
		}
	}

	//Step 2: This should next call the /state-summary endpoint from the master and figure out all the remaining resources
	resp, err := http.Get(fmt.Sprintf("%s/state.json", dc.Endpoint))

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return 0.0, 0.0, 0.0, err
	}

	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	err = json.Unmarshal(body, &mState)
	if err != nil {
		log.Printf("Json Unmarshall error = %v", err)
		return 0.0, 0.0, 0.0, err
	}

	for _, fw := range mState.Frameworks {
		if fw.Role == role {
			usedR := fw.Used_Resources.(map[string]interface{})
			uCPU = usedR["cpus"].(float64)
			uMEM = usedR["mem"].(float64)
			uDISK = usedR["disk"].(float64)
		}
	}

	if qCPU >= uCPU && qMEM >= uMEM && qDISK >= uDISK {

		CPU = qCPU - uCPU
		MEM = qMEM - uMEM
		DISK = qDISK - uDISK
	}

	fmt.Println("remaining resources are ", CPU, MEM, DISK)

	return CPU, MEM, DISK, nil
}
