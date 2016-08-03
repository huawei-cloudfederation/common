//Package quotalib will implement all the functions that should be used to contact masters

package quotalib

import (
	"common/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//Declare the structure which will give quota information
type QInfos struct {
	Infos []QuotaInfo "json:'info'"
}

type QuotaInfo struct {
	Guarantee []GuaranteeInfo "json:'guarantee'" //allocate guaranteed quota to a role associted
	Role      string          "json:'role'"      //quota is applied with this role

}

type GuaranteeInfo struct {
	Name   string     "json:'name'" // name of the resource, ex:cpu
	Role   string     "json:'role'"
	Scalar ScalarInfo "json:'scalar'" //scalar resources info
	Type   string     "json:'type'"   //type of resource ex:scalar
}

type ScalarInfo struct {
	Value float64 "json:'value'" //value of resource,ex:cpu='2'
}

//SetQuota will read a json file the local disk and performe a SET Quota HTTP api call
//role : Sipply the role for whcih the Quota to be set
//inputPath : Path from where quota json file should be read
func SetQuota(dc typ.DC, role string, inputPath string) error {
	log.Printf("master Endpoint is", dc.Endpoint)

	//Function implementation

	buf, err := ioutil.ReadFile(inputPath)
	if err != nil {
		log.Printf("Unable to read file = %v", err)
	}

	body := bytes.NewBuffer(buf)

	resp, err := http.Post(fmt.Sprintf("http://%s/quota/%s", dc.Endpoint, role), "text/json", body)
	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return nil
	}

	response,_ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	log.Println(string(response))

	return nil
}

//DelQuota Will delete the Quota set on the given role
//role : Sipply the role name , usually Federation
func DelQuota(dc typ.DC, role string) error {

	//Function implementation

	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/quota/%s", dc.Endpoint, role), nil)

	// handle err
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
	}

	log.Printf("the response from the master = %v", resp)
	return nil

}

//GetQuota Will get a json representing Quota for the given role
//role : Sipply the role name , usually Federation
func GetQuota(dc typ.DC, role string) ([]GuaranteeInfo, error) {
	var data QInfos
	var guarante []GuaranteeInfo

	//Function implementation
	resp, err := http.Get(fmt.Sprintf("http://%s/quota/%s", dc.Endpoint, role))

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return guarante,nil	
	}

	body,_ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()


	error := json.Unmarshal([]byte(body), &data)

	if error != nil {
		log.Printf("Json Unmarshall error = %v", err)
	}

	for key, _ := range data.Infos {
		if data.Infos[key].Role == role {
			guarante = data.Infos[key].Guarantee
			log.Println(guarante)
			break
		} else if data.Infos[key].Role != role {
			log.Println("the quota for the", role, " doesn't exist", guarante)
		}
	}

	return guarante, nil
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
	quota, _ := GetQuota(dc, role)

	log.Println("the quota is ", quota, "\n")

	for index, _ := range quota {
		if quota[index].Name == "cpus" {
			qCPU = quota[index].Scalar.Value
		}
		if quota[index].Name == "mem" {
			qMEM = quota[index].Scalar.Value
		}
		if quota[index].Name == "disk" {
			qDISK = quota[index].Scalar.Value
		}
	}

	//Step 2: This should next call the /state-summary endpoint from the master and figure out all the remaining resources
	resp, err := http.Get(fmt.Sprintf("http://%s/state.json", dc.Endpoint))

	if err != nil {
		log.Printf("Unable to reach the Master error = %v", err)
		return CPU,MEM,DISK,nil
	}


	body, _ := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	err = json.Unmarshal(body, &mState)
	if err != nil {
		log.Printf("Json Unmarshall error = %v", err)
	}

	for key, _ := range mState.Frameworks {
		if mState.Frameworks[key].Role == role {
			usedR := mState.Frameworks[key].Used_Resources.(map[string]interface{})
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
