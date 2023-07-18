package clickup

import (
	"net/http"
	"reflect"
	"testing"
)

func TestClient_ToJson(t *testing.T) {
	type fields struct {
		token      string
		httpClient *http.Client
		apiUrl     string
		logger     Logger
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				token:      tt.fields.token,
				httpClient: tt.fields.httpClient,
				apiUrl:     tt.fields.apiUrl,
				logger:     tt.fields.logger,
			}
			if got := c.ToJson(tt.args.data); got != tt.want {
				t.Errorf("Client.ToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ToJsonByte(t *testing.T) {
	type fields struct {
		token      string
		httpClient *http.Client
		apiUrl     string
		logger     Logger
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				token:      tt.fields.token,
				httpClient: tt.fields.httpClient,
				apiUrl:     tt.fields.apiUrl,
				logger:     tt.fields.logger,
			}
			if got := c.ToJsonByte(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ToJsonByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDefaultClient(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefaultClient(tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefaultClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		token  string
		apiUrl string
		logger Logger
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.token, tt.args.apiUrl, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithLogger(t *testing.T) {
	type args struct {
		token  string
		apiUrl string
		logger Logger
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClientWithLogger(tt.args.token, tt.args.apiUrl, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClientWithLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_requestGet(t *testing.T) {
	type fields struct {
		token      string
		httpClient *http.Client
		apiUrl     string
		logger     Logger
	}
	type args struct {
		endpoint string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				token:      tt.fields.token,
				httpClient: tt.fields.httpClient,
				apiUrl:     tt.fields.apiUrl,
				logger:     tt.fields.logger,
			}
			got, err := c.requestGet(tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.requestGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.requestGet() = %v, want %v", got, tt.want)
			}
		})
	}
}
