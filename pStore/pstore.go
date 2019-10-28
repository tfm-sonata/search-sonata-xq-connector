package fraud

import (
	"encoding/json"
	"fmt"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/awshelper"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/model"
	"git-codecommit.eu-central-1.amazonaws.com/search-sonata-xq-connector/util"
	"log"
)

type ParameterStoreConfig struct {
	ConfigProvider string
	Servicename    string
}

var storeConfigs map[string]map[string]*model.Config

// LoadServiceConfigFromParameterStore gets the configuration for the service from the parameterStore
// for the specified Service
func LoadServiceConfigFromParameterStore(configProvider string, service string) error {
	log.Println("loading config store for : " + configProvider + " service " + service)
	result, err := awshelper.InvokeLambdaResponseFunction(configProvider, []byte(fmt.Sprintf(`{"SERVICENAME":"%s"}`, service)))
	var configs map[string][]*model.Config
	// convert the string into a map with configurations
	err = json.Unmarshal(result, &configs)
	if err != nil {
		log.Fatalln("config strore umarshalling error : ", err)
		return err
	}

	log.Println("config store unmarshalled : ", configs)

	storeConfigs = make(map[string]map[string]*model.Config, len(configs))

	for product, confs := range configs {
		fmt.Println(product, len(confs))
		storeConfigs[product] = make(map[string]*model.Config, len(confs))
		for _, conf := range confs {
			storeConfigs[product][conf.Id] = conf
		}

	}
	log.Println("storeConfigs filled : ", storeConfigs)

	return nil
}

// Product and Channel delivered
func ServiceConfigForProductChannel(configProvider string, service string, product string, channel string) (*model.Config, error) {
	if storeConfigs == nil {
		err := LoadServiceConfigFromParameterStore(configProvider, service)
		if err != nil {
			return nil, err
		}
	}
	log.Println("Configs: ", util.MapToString(storeConfigs))

	// absolute fallback
	resultConfig := storeConfigs["default"]["default"]

	// if porduct exist ...
	if p := storeConfigs[product]; p != nil {
		// search for channel ...
		if c := p[channel]; c != nil {
			resultConfig = c
		} else {
			resultConfig = p["default"]
		}
	}

	model.RequestConfig = resultConfig

	return resultConfig, nil

}
