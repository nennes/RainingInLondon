package config

import (
	"github.com/nennes/RainingInLondon/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/nennes/RainingInLondon/utils"
	"io/ioutil"
)

type Codes struct {
	Weather    map[string]string
	Visibility map[string]string
}

type ForecastLongTerm struct {
	RegionalFcst struct {
		FcstPeriods struct {
			Period []struct {
				Paragraph []struct {
					Text  string `json:"$"`
					Title string `json:"title"`
				} `json:"Paragraph"`
				ID string `json:"id"`
			} `json:"Period"`
		} `json:"FcstPeriods"`
		CreatedOn string `json:"createdOn"`
		IssuedAt  string `json:"issuedAt"`
		RegionID  string `json:"regionId"`
	} `json:"RegionalFcst"`
}

var (
	WeatherCodes = &Codes{}
)

func init() {
	configYaml, readErr := ioutil.ReadFile("config/codes.yaml")
	utils.ErrorPanic(readErr)

	yamlErr := yaml.Unmarshal(configYaml, WeatherCodes)
	utils.ErrorPanic(yamlErr)

}
