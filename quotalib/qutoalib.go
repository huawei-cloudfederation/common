//Package quotalib will implement all the functions that should be used to contact masters

package quotalib

//SetQuota will read a json file the local disk and performe a SET Quota HTTP api call
//role : Sipply the role for whcih the Quota to be set
//inputPath : Path from where quota json file should be read
func SetQuota(role string, inputPath string) error {

	//Function implementation

	return nil
}

//DelQuota Will delete the Quota set on the given role
//role : Sipply the role name , usually Federation
func DelQuota(role string) error {

	//Function implementation

	return nil

}

//GetQuota Will get a json representing Quota for the given role
//role : Sipply the role name , usually Federation
func GetQuota(role string) ([]byte, error) {

	var json_data []byte

	//Function implementation

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
