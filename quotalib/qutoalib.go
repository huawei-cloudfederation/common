//Package quotalib will implement all the functions that should be used to contact masters

package quotalib

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

//SetQuota will read a json file the local disk and performe a SET Quota HTTP api call
//role : Sipply the role for whcih the Quota to be set
//inputPath : Path from where quota json file should be read
func SetQuota(role string, inputPath string) error {

	//Function implementation

	buf,err :=  ioutil.ReadFile(inputPath) 
	if err != nil {
                log.Printf("Unable to read file = %v", err)
		return err
        }
	body := bytes.NewBuffer(buf)

	resp ,err := http.Post(fmt.Sprintf("http://172.31.44.22:5050/quota/%s", role),"text/json",body)
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
func DelQuota(role string) error {

	//Function implementation
 
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://172.31.44.22:5050/quota/%s", role), nil)
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
func GetQuota(role string) ([]byte, error) {

	var json_data []byte

	var data map[string]interface{}
	//Function implementation
	 resp, err := http.Get(fmt.Sprintf("http://172.31.44.22:5050/quota/%s",role))

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
        log.Println("Get: ",data["infos"])

	return json_data, nil
}

//RemainingResources
//role: Role name, usually Federation
//This should automatically resolve the Masters Endpoint from the configuration
//Retruns how much of the gaurentee has been satisfied for a given role
//If a role has a quota of 50 CPUs and out of that 35 have been used
func RemainingResource(role string) (float64, float64, float64, error) {

	var CPU, MEM, DISK float64

	//Step 1: This should first call GetQuota() and interpret gaurantee

	//Step 2: This should next call the /state-summary endpoint from the master and figure out all the remaining resources

	return CPU, MEM, DISK, nil
}
