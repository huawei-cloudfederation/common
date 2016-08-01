package typ

//Master_State Complete definition of the /state.json from the master
//Details we do not need at the moment are left as interface
//Framework ahs a full definition which we need
//Used for Unmarshalling the state.json endpoint to a structure
//For example if you need to process all the frameworks in the mesos then
// var json_state MasterState
// err := json.Unmarshal(date, &jsonState)
// for _, fw := range(jsonState.Frameworks) {
//      fmt.Printf("%s\n", fw.Role)
//  }
//

type MasterState struct {
	Version                string        `json:"version"`
	GitSha                 string        `json:"git_sha"`
	GitBranch              string        `json:"git_branch"`
	BuildDate              string        `json:"build_date"`
	BuildTime              float64       `json:"build_time"`
	BuildUser              string        `json:"build_user"`
	StartTime              float64       `json:"start_time"`
	ElectedTime            float64       `json:"elected_time"`
	ID                     string        `json:"id"`
	Pid                    string        `json:"pid"`
	Hostname               string        `json:"hostname"`
	ActivatedSlaves        float64       `json:"activated_slaves"`
	DeactivatedSlaves      float64       `json:"deactivated_slaves"`
	Leader                 string        `json:"leader"`
	Frameworks             []FW          `json:"frameworks"`
	Slaves                 []interface{} `json:"slaves"`
	Flags                  interface{}   `json:"flags"`
	CompletedFrameworks    []interface{} `json:"completed_frameworks"`
	OrphanTasks            []interface{} `json:"orphan_tasks"`
	UnregisteredFrameworks []interface{} `json:"unregistered_frameworks"`
}

//FW detailed structure definition of Mesos framework representation
//Used for unmarrshalling state.json to a structure
type FW struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	Pid               string `json:"pid"`
	used_resources    interface{}
	offered_resources interface{}
	Capabilities      []interface{} `json:"capabilities"`
	Hostname          string        `json:"hostname"`
	WebuiURL          string        `json:"webui_url"`
	Active            bool          `json:"active"`
	User              string        `json:"user"`
	FailoverTimeout   float64       `json:"failover_timeout"`
	Checkpoint        bool          `json:"checkpoint"`
	Role              string        `json:"role"`
	RegisteredTime    float64       `json:"registered_time"`
	UnregisteredTime  float64       `json:"unregistered_time"`
	Resources         struct {
		Disk  float64 `json:"disk"`
		Mem   float64 `json:"mem"`
		Gpus  float64 `json:"gpus"`
		Cpus  float64 `json:"cpus"`
		Ports string  `json:"ports"`
	} `json:"resources"`
	Tasks          []interface{} `json:"tasks"`
	CompletedTasks []interface{} `json:"completed_tasks"`
	offers         []interface{}
	executors      []interface{}
}
