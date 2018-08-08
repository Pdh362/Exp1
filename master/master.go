package master

import (
	"encoding/json"
	"github.com/Pdh362/Exp1/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------

var MasterList map[string][]string = make(map[string][]string)

//----------------------------------------------------------------------------------------------------------------------

func Results(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"results": MasterList,
	})
	c.Next()
}

//----------------------------------------------------------------------------------------------------------------------
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

func UpdateFolder(path string, results []string) {

	// log.Standard.Printf("Path = %s, results = %s", path, results)

	MasterList[path] = results
}
