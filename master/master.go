package master

import (
	"encoding/json"
	"github.com/Pdh362/Exp1/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------
var CollatedList = make(map[string][]string)

//----------------------------------------------------------------------------------------------------------------------
// Update:
//
// Triggered when a request for the collated list is made.
// Simply report it back as a json response
//
func Results(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"results": CollatedList,
	})
	c.Next()
}

//----------------------------------------------------------------------------------------------------------------------
// Update:
//
// Triggered when the master receives an update from a watcher.
// Reform the posted data into a useful type, and send it over to UpdateFolder
//
// Not too happy with the data process below: given time, I'd optimise the data
// sent by the watcher, as well as how I process the data here.
//
func Update(c *gin.Context) {

	results, err := c.GetRawData()
	if err != nil {
		log.Standard.Error("Master: Failed to get gin context raw data")
		return
	}

	// Unmarshal the binary blob into a json data object
	var jsonData map[string]interface{}
	err = json.Unmarshal(results, &jsonData)
	if err != nil {
		log.Standard.Error("Master: Failed to unmarshal raw data results:" + err.Error())
		return
	}

	// Not happy about having to remap like this : find alternative approach?
	resinterface := jsonData["results"].([]interface{})
	finalres := make([]string, len(resinterface))
	for i, v := range resinterface {
		finalres[i] = v.(string)
	}

	UpdateFolder(jsonData["path"].(string), finalres)
}

//----------------------------------------------------------------------------------------------------------------------
// UpdateFolder:
//
// Updates the global map that holds the collated information
//
func UpdateFolder(path string, results []string) {

	// log.Standard.Printf("Path = %s, results = %s", path, results)

	CollatedList[path] = results
}
