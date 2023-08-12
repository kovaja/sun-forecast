package forecast

import (
	"fmt"
	"kovaja/sun-forecast/httpClient"
	"kovaja/sun-forecast/utils"
)

func getUrl() (string, error) {
	siteId, err := utils.GetEnvVariable("SOLCAST_SITE_ID")
	url := fmt.Sprintf("https://api.solcast.com.au/rooftop_sites/%s/forecasts?format=json", siteId)

	return utils.ReturnStringResultOrError(url, err)
}

func fetchForecasts() (*SolcastApiForcastResponse, error) {
	apiKey, err := utils.GetEnvVariable("SOLCAST_API_KEY")

	url, err := getUrl()
	if err != nil {
		return nil, utils.CustomError("Failed to build solcast api url", err)
	}

	var body SolcastApiForcastResponse
	err = httpClient.GetJsonWithAuth(url, apiKey, &body)

	if err != nil {
		return nil, utils.CustomError("Failed to call solcast api", err)
	}

	return &body, nil
}
