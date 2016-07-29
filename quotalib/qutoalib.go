//Package quotalib will implement all the functions that should be used to contact masters

package quotalib

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"../types"
)

type QInfos struct{
	Infos []subInfo "json:'info'"
	
}

type subInfo struct{
	Guarantee []subsubInfo "json:'guarantee'"
	Role string "json:'role'"
	
}


type subsubInfo struct{
		Name string "json:'name'"
		Role string "json:'role'"
		Scalar subScalar "json:'scalar'"
		Type string			   "json:'type'"
}


type subScalar struct{
	Value float64 "json:'value'"
}

//SetQuota will read a json file the local disk and performe a SET Quota HTTP api call
//role : Sipply the role for whcih the Quota to be set
//inputPath : Path from where quota json file should be read
func SetQuota(role string, inputPath string,dc typ.DC) error {
	log.Printf("master Endpoint is",dc.Endpoint)

	//Function implementation

	buf,err :=  ioutil.ReadFile(inputPath) 
	if err != nil {
                log.Printf("Unable to read file = %v", err)
		return err
        }
	body := bytes.NewBuffer(buf)

	resp ,err := http.Post(fmt.Sprintf("http://%s/quota/%s",dc.Endpoint, role),"text/json",body)
	 if err != nil {
                log.Printf("Unable to reach the Master error = %v", err)
		return err
        }

        defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(response))

	return nil
}

//DelQuota Will delete the Quota set on the given role
//role : Sipply the role name , usually Federation
func DelQuota(role string,dc typ.DC) error {

	//Function implementation
 
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/quota/%s",dc.Endpoint, role), nil)
	// handle err
	resp, err := http.DefaultClient.Do(req)

	 if err != nil {

                log.Printf("Unable to reach the Master error = %v", err)
                //return
        }

	 log.Printf("the response from the master = %v", resp)	
	return nil

}

//GetQuota Will get a json representing Quota for the given role
//role : Sipply the role name , usually Federation
func GetQuota(role string,dc typ.DC) (subInfo, error) {
	var json_data []byte
	var index int
	var data QInfos

	//Function implementation
	 resp, err := http.Get(fmt.Sprintf("http://%s/quota/%s",dc.Endpoint,role))

        if err != nil {

                log.Printf("Unable to reach the Master error = %v", err)
                //return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)

        if err != nil {
                log.Printf("Unable to read the body error = %v", err)
                //return
        }
	json_data = []byte(body)

	 json.Unmarshal(json_data,&data)

        if err != nil {
                log.Printf("Json Unmarshall error = %v", err)
                //return
        }

	for index = 0 ; index < len(data.Infos) ; index++ {
		if data.Infos[index].Role == role {
		  break
        	}
	}
	log.Println("the quota for role is ",data.Infos[index].Guarantee)

	return data.Infos[index],  nil
}

//RemainingResources
//role: Role name, usually Federation
//This should automatically resolve the Masters Endpoint from the configuration
//Retruns how much of the gaurentee has been satisfied for a given role
//If a role has a quota of 50 CPUs and out of that 35 have been used
func RemainingResource(role string,dc typ.DC) (float64, float64, float64, error) {

	var CPU, MEM, DISK float64
	var rCPU, rMEM, rDISK float64
	var qCPU, qMEM, qDISK float64

	var data map[string]interface{}
	
	var flag int
	//Step 1: This should first call GetQuota() and interpret gaurantee
	quota , err := GetQuota(role,dc)
        if err != nil {
                log.Printf("Unable to GetQuota = %v", err)
        }

	for index := 0 ; index < len(quota.Guarantee) ; index++{
		      if quota.Guarantee[index].Name == "cpus" {
			    qCPU = quota.Guarantee[index].Scalar.Value
			}
		      if quota.Guarantee[index].Name == "mem" {
			    qMEM = quota.Guarantee[index].Scalar.Value
			}
		      if quota.Guarantee[index].Name == "disk" {
			    qDISK = quota.Guarantee[index].Scalar.Value
			}
	}

	//Step 2: This should next call the /state-summary endpoint from the master and figure out all the remaining resources
	 resp, err := http.Get(fmt.Sprintf("http://%s/state/summary",dc.Endpoint))

        if err != nil {

                log.Printf("Unable to reach the Master error = %v", err)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)

        if err != nil {
                log.Printf("Unable to read the body error = %v", err)
        }

	 err = json.Unmarshal(body, &data)
        if err != nil {
                log.Printf("Json Unmarshall error = %v", err)
        }
	

	for k ,v := range data {
          switch vv := v.(type) {
            case []interface{}:
		if k == "frameworks" {
                for j, u := range vv {
                 switch aa :=  u.(type) {
                    case map[string]interface{}:
			for n, m := range aa {
        	         switch aaa := m.(type) {
		            case string:
			     if n == "role" {
			     	if m == role {
				flag = 1
			       }else {
				 flag = 0
				}
			     }
                	    case map[string]interface{}:
			     if n == "used_resources" && flag == 1 {
                	        for l, p := range aaa {
			         if l == "cpus"{
				    rCPU = p.(float64)
				  }
			         if l == "mem"{
				    rMEM = p.(float64)
				}
			         if l == "disk"{
				    rDISK = p.(float64)
			      	}
			     } 
			   }

         	           default:
                	        fmt.Println(n, "is of a type I don't know how to handle")
                    	}
 
                          }
                    default:
                        fmt.Println(j, "is of a type I don't know how to handle")
   		     }
		    }
		    }
            default:
                fmt.Println(k, "is of a type I don't know how to handle")
            }
         }

	CPU = qCPU - rCPU
	MEM = qMEM - rMEM
	DISK = qDISK - rDISK
                	           fmt.Println("remaining resources are ",CPU,MEM,DISK)

	return CPU, MEM, DISK, nil
}
