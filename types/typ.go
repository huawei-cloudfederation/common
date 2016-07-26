//Package typ is a package in common library, which should define all the datastructures that will be used both by gossiper and policy engine
//Things like json strucutre that will be passed around should be defined here
package typ

//Declare some structure that will eb common for both Anonymous and Gossiper modulesv
type DC struct {
	OutOfResource bool
	Name          string
	City          string
	Country       string
	Endpoint      string
	CPU           float64
	MEM           float64
	DISK          float64
	Ucpu          float64 //Remaining CPU
	Umem          float64 //Remaining Memory
	Udisk         float64 //Remaining Disk
	LastUpdate    int64   //Time stamp of current DC status
	LastOOR       int64   //Time stamp of when was the last OOR Happpend
	IsActiveDC    bool
}

func init() {

	//All the initialization code should go here
}
