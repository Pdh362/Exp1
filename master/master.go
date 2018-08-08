package master

import (
	"encoding/json"
	"github.com/Pdh362/Exp1/log"
	"github.com/gin-gonic/gin"
)

func Results(c *gin.Context) {

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
	aString := make([]string, len(resinterface))
	for i, v := range resinterface {
		aString[i] = v.(string)
	}

	UpdateFolder(jsonData["path"].(string), aString)
}

func UpdateFolder(path string, results []string) {

	log.Standard.Printf("Path = %s, results = %s", path, results)
}
