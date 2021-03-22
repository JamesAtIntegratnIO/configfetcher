package main

import (
	"fmt"

	"github.com/jamesattensure/configfetcher"
	"github.com/mitchellh/mapstructure"
)

// Struct that is used to pass to the configfetcher.ConfigStruct
type yamlconfig struct {
	Data struct {
		NestOneOne string `yaml:"nest_one_one"`
		NestOneTwo string `yaml:"nest_one_two"`
		NestTwo    struct {
			NestTwoOne  string   `yaml:"nest_two_one"`
			NestTwoList []string `yaml:"nest_two_list"`
		}
	}
}

// Struct that is used to convert the returned map from configfetcher.ConfigStruct with github.com/mitchellh/mapstructure
type Conf struct {
	Data struct {
		NestOneOne string `mapstructure:"nest_one_one"`
		NestOneTwo string `mapstructure:"nest_one_two"`
		NestTwo    struct {
			NestTwoOne  string   `mapstructure:"nest_two_one"`
			NestTwoList []string `mapstructure:"nest_two_list"`
		}
	}
}

func main() {

	//
	//  Example using reading a file
	//

	cf := configfetcher.Data{}
	_, err := cf.ReadFile("/Volumes/CaseSensitive/configfetcher/cmd/data.yaml")
	if err != nil {
		fmt.Println(err)
	}
	cf.ConfigStruct = yamlconfig{}
	if err := cf.GetConfigStruct("yaml"); err != nil {
		fmt.Println(err)
	}

	var conf Conf
	if err := mapstructure.Decode(cf.ConfigStruct, &conf); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", conf)

	//
	// Example using reading from GCP Secret Manager
	//

	cf2 := configfetcher.Data{}
	_, err = cf2.ReadGCPSMSecret()
	if err != nil {
		fmt.Println(err)
	}
	if err := cf2.GetConfigStruct("yaml"); err != nil {
		fmt.Println(err)
	}
	var conf2 Conf
	if err := mapstructure.Decode(cf2.ConfigStruct, &conf2); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", conf2)

}
