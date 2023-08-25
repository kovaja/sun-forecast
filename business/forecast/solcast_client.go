package forecast

import (
	"fmt"
	"kovaja/sun-forecast/utils"
	"kovaja/sun-forecast/utils/httpclient"
	"kovaja/sun-forecast/utils/logger"
)

func getUrl() (string, error) {
	siteId, err := utils.GetEnvVariable("SOLCAST_SITE_ID")
	url := fmt.Sprintf("https://api.solcast.com.au/rooftop_sites/%s/forecasts?format=json", siteId)

	return utils.ReturnStringResultOrError(url, err)
}

func fetchForecasts() (*SolcastApiForcastResponse, error) {
	apiKey, err := utils.GetEnvVariable("SOLCAST_API_KEY")
	dev, err := utils.GetEnvVariable("DEV")

	url, err := getUrl()
	if err != nil {
		return nil, utils.CustomError("Failed to build solcast api url", err)
	}

	var body SolcastApiForcastResponse
	err = httpclient.GetJsonWithAuth(url, apiKey, &body)

	if err != nil {
		if dev == "1" {
			logger.Log("Call to solcast api failed in dev, using mock data, %v", err)

			err = utils.ReadJson("./forecast/mock.json", &body)
			if err != nil {
				return nil, utils.CustomError("Failed to read mock data", err)
			}
		} else {
			return nil, utils.CustomError("Failed to call solcast api", err)
		}
	}

	return &body, nil
}
