// Package configfetcher provides a simple abstraction to fetch config values
// from a file or GCP Secrets Manager to populate a struct.
package configfetcher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	sm "cloud.google.com/go/secretmanager/apiv1beta1"
	pb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"gopkg.in/yaml.v2"
)

// GetConfig is the entry point ineterface to either get the raw data or to populate a struct
type GetConfig interface {
	Data(useGCPSecrets bool, filePath string) ([]byte, error)
	Struct(useGCPSecrets bool, configType string, filePath string, configStruct interface{}) error
}

// Data
// useGCPSecrets (bool): If set it is expected that the following environment variables are set:
//		PROJECT_ID
//		SECRET_NAME
//		SECRET_VERSION
// filePath (string): If passing a yaml or json file directly this points to the file. If not used just give empty "".
func Data(useGCPSecrets bool, filePath string) ([]byte, error) {
	if useGCPSecrets {
		gcloudVars := setGcloudVars()
		data, err := gcloudVars.getSecretFromGSM()
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		return data, nil

	}
}

// Struct populates a struct to provide a config
// useGCPSecrets (bool): If set it is expected that the following environment variables are set:
//		PROJECT_ID
//		SECRET_NAME
//		SECRET_VERSION
// configType (string): expects `yaml` or `json` to properly unmarshal the data into the struct
// filePath (string): If passing a yaml or json file directly this points to the file. If not used just give empty "".
// config (interface{}): expects a struct that has either yaml or json mappings and matches the data that is used to populate it.
func Struct(useGCPSecrets bool,
	configType string,
	filePath string, configStruct interface{}) error {

	data, err := Data(useGCPSecrets, filePath)
	if err != nil {
		return err
	}
	if err := getYamlOrJsonData(configType, data, configStruct); err != nil {
		return err
	}
	return nil
}

func getYamlOrJsonData(configType string, data []byte, configStruct interface{}) error {
	switch configType {
	case "yaml":
		err := yaml.Unmarshal([]byte(data), &configStruct)
		if err != nil {
			return err
		}
	case "json":
		err := json.Unmarshal([]byte(data), &configStruct)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid `configType` provided. expects a string of yaml or json to match the data type that is being unmarshaled")
	}
	return nil
}

type gcloudVars struct {
	projectID     string
	secretName    string
	secretVersion string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setGcloudVars() gcloudVars {
	return gcloudVars{
		getEnv("PROJECT_ID", ""),
		getEnv("SECRET_NAME", ""),
		getEnv("SECRET_VERSION", "latest"),
	}
}

func (g gcloudVars) getSecretFromGSM() ([]byte, error) {
	if g.projectID == "" || g.secretName == "" {
		return nil, errors.New(`
		environment variables for gcp are not set
		please set 'PROJECT_ID', 'SECRET_NAME', and 'SECRET_VERSION'
		`)
	}
	ctx := context.Background()
	c, err := sm.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &pb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s",
			g.projectID,
			g.secretName,
			g.secretVersion),
	}
	resp, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Data, nil
}
