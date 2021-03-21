// Package configfetcher provides a simple abstraction to fetch config values
// from a file or GCP Secrets Manager to populate a struct.
package configfetcher

import (
	"os"
	"reflect"
	"testing"
)

func TestGetConfigStructFromGCPSM(t *testing.T) {
	type args struct {
		configType   string
		configStruct interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigStructFromGCPSM(tt.args.configType, tt.args.configStruct); (err != nil) != tt.wantErr {
				t.Errorf("GetConfigStructFromGCPSM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetConfigStructFromFile(t *testing.T) {
	type args struct {
		configType   string
		filePath     string
		configStruct interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetConfigStructFromFile(tt.args.configType, tt.args.filePath, tt.args.configStruct); (err != nil) != tt.wantErr {
				t.Errorf("GetConfigStructFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getYamlOrJsonData(t *testing.T) {
	type args struct {
		configType   string
		data         []byte
		configStruct interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getYamlOrJsonData(tt.args.configType, tt.args.data, tt.args.configStruct); (err != nil) != tt.wantErr {
				t.Errorf("getYamlOrJsonData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gcloudVars_getSecretFromGSM(t *testing.T) {
	type fields struct {
		projectID     string
		secretName    string
		secretVersion string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := gcloudVars{
				projectID:     tt.fields.projectID,
				secretName:    tt.fields.secretName,
				secretVersion: tt.fields.secretVersion,
			}
			got, err := g.getSecretFromGSM()
			if (err != nil) != tt.wantErr {
				t.Errorf("gcloudVars.getSecretFromGSM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("gcloudVars.getSecretFromGSM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setGcloudVars(t *testing.T) {
	tests := []struct {
		name string
		want gcloudVars
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setGcloudVars(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setGcloudVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "mykey")
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test getEnv",
			args: args{
				key:      "TEST_KEY",
				fallback: "myfallback",
			},
			want: "mykey",
		},
		{
			name: "test getEnv",
			args: args{
				key:      "TEST_KEY2",
				fallback: "myfallback",
			},
			want: "myfallback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Unsetenv("TEST_KEY")
}
