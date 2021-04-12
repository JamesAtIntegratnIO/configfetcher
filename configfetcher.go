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

// GcloudVars contains the data needed to connect to GCP Secret Manager
type GcloudVars struct {
	ProjectID     string
	SecretName    string
	SecretVersion string
}

// Data contains everything needed to interact with this abstraction
type Data struct {
	ConfigData []byte
	GcloudVars
	ConfigStruct interface{}
}

// ReadGCPSMSecret returns the Data struct and error after attempting to fetch the environment variables:
//	`PROJECT_ID`      required
//	`SECRET_NAME`     required
//	`SECRET_VERSION`  optional
func (d *Data) ReadGCPSMSecret() (*Data, error) {
	d.setGcloudVars()
	return d.getSecretFromGSM()
}

// ReadFile returns the Data struct after attempting to read the data from a file
func (d *Data) ReadFile(filePath string) (*Data, error) {
	var err error
	d.ConfigData, err = ioutil.ReadFile(filePath)
	return d, err
}

// GetConfigStruct populates the Data.ConfigStruct interface with the data from Data.ConfigData
// It expects either `ReadGCPSMSecret()` or `ReadFile()` has been ran to populate Data.ConfigData
func (d *Data) GetConfigStruct(configType string) error {
	if d.ConfigData == nil {
		return errors.New("no config data provided\npopulate with ReadGCPSMSecret() or ReadFile()")
	}
	switch configType {
	case "yaml":
		err := yaml.Unmarshal([]byte(d.ConfigData), &d.ConfigStruct)
		if err != nil {
			return err
		}
	case "json":
		err := json.Unmarshal([]byte(d.ConfigData), d.ConfigStruct)
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid `configType` provided. expects a string of `yaml` or `json` to match the data type that is being unmarshaled")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func (d *Data) setGcloudVars() *Data {
	d.ProjectID = getEnv("PROJECT_ID", "")
	d.SecretName = getEnv("SECRET_NAME", "")
	d.SecretVersion = getEnv("SECRET_VERSION", "latest")
	return d
}

func (d *Data) getSecretFromGSM() (*Data, error) {
	if d.ProjectID == "" || d.SecretName == "" {
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
			d.ProjectID,
			d.SecretName,
			d.SecretVersion),
	}
	resp, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	d.ConfigData = resp.Payload.Data
	return d, nil
}
