package main

import (
	"context"
	"encoding/json"
	"github.com/graphikDB/graphik/graphik-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

func init() {
	godotenv.Load()
}

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: func() terraform.ResourceProvider {
		return &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAPHIK_HOST", nil),
				},
				"client_id": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAPHIK_CLIENT_ID", nil),
				},
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAPHIK_USERNAME", nil),
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAPHIK_PASSWORD", nil),
				},
				"oidc_metadata_uri": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("GRAPHIK_OPEN_ID_METADATA_URI", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"graphik_index": {
					Schema: map[string]*schema.Schema{},
					Create: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						return false, nil
					},
					CustomizeDiff: nil,
					Importer: &schema.ResourceImporter{
						State: schema.ImportStatePassthrough,
					},
					DeprecationMessage: "",
					Timeouts:           nil,
					Description:        "a graph primitive used for fast lookups of docs/connections that pass a boolean CEL expression",
				},
				"graphik_trigger": {
					Schema:         map[string]*schema.Schema{},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						return false, nil
					},
					CustomizeDiff: nil,
					Importer: &schema.ResourceImporter{
						State: schema.ImportStatePassthrough,
					},
					DeprecationMessage: "",
					Timeouts:           nil,
					Description:        "used to automatically mutate the attributes of documents/connections before they are commited to the database",
				},
				"graphik_constraint": {
					Schema:         map[string]*schema.Schema{},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						return false, nil
					},
					CustomizeDiff: nil,
					Importer: &schema.ResourceImporter{
						State: schema.ImportStatePassthrough,
					},
					DeprecationMessage: "",
					Timeouts:           nil,
					Description:        "a graph primitive used to validate custom doc/connection constraints",
				},
				"graphik_authorizer": {
					Schema:         map[string]*schema.Schema{},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						_ = i.(*graphik.Client)
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						_ = i.(*graphik.Client)
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						_ = i.(*graphik.Client)
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						_ = i.(*graphik.Client)
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						_ = i.(*graphik.Client)
						return false, nil
					},
					Importer: &schema.ResourceImporter{
						State: schema.ImportStatePassthrough,
					},
					DeprecationMessage: "",
					Timeouts:           nil,
					Description:        "a graph primitive used for authorizing inbound requests and/or responses(see AuthTarget)",
				},
			},
			ConfigureFunc: func(data *schema.ResourceData) (interface{}, error) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				host := data.Get("host").(string)
				clientId := data.Get("client_id").(string)
				username := data.Get("username").(string)
				password := data.Get("password").(string)
				metadataUri := data.Get("oidc_metadata_uri").(string)
				metadata := map[string]interface{}{}
				resp, err := http.Get(metadataUri)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get oidc metadata")
				}
				defer resp.Body.Close()
				if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
					return nil, errors.Wrap(err, "failed to get oidc metadata")
				}
				cfg := &oauth2.Config{
					ClientID: clientId,
					//ClientSecret: data.Get("client_secret").(string),
					Endpoint: oauth2.Endpoint{
						AuthURL:  metadata["authorization_endpoint"].(string),
						TokenURL: metadata["token_endpoint"].(string),
					},
					//RedirectURL: p.config.RedirectURL,
					Scopes: []string{"openid", "email", "profile"},
				}
				token, err := cfg.PasswordCredentialsToken(ctx, username, password)
				if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
					return nil, errors.Wrap(err, "failed to get token")
				}

				client, err := graphik.NewClient(ctx, host,
					graphik.WithTokenSource(oauth2.StaticTokenSource(token)),
					graphik.WithRetry(2),
				)
				if err != nil {
					return nil, errors.Wrap(err, "failed to create graphik client")
				}
				return client, nil
			},
		}
	}})
}
