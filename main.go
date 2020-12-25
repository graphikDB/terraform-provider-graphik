package main

import (
	"context"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/empty"
	apipb "github.com/graphikDB/graphik/gen/grpc/go"
	"github.com/graphikDB/graphik/graphik-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "unique name of the index",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"gtype": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"expression": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"target_docs": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
						"target_connections": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
					},
					Create: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetIndexes().GetIndexes()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Index{
									Name:        data.Get("name").(string),
									Gtype:       data.Get("gtype").(string),
									Expression:  data.Get("expression").(string),
									Docs:        data.Get("target_docs").(bool),
									Connections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Index{
								Name:        data.Get("name").(string),
								Gtype:       data.Get("gtype").(string),
								Expression:  data.Get("expression").(string),
								Docs:        data.Get("target_docs").(bool),
								Connections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetIndexes(ctx, &apipb.Indexes{Indexes: values}); err != nil {
							return err
						}
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						id := data.Id()
						for _, a := range scheme.GetIndexes().GetIndexes() {
							if a.GetName() == id {
								if err := data.Set("name", a.GetName()); err != nil {
									return err
								}
								if err := data.Set("gtype", a.GetGtype()); err != nil {
									return err
								}
								if err := data.Set("expression", a.GetExpression()); err != nil {
									return err
								}
								if err := data.Set("target_connections", a.GetConnections()); err != nil {
									return err
								}
								if err := data.Set("target_docs", a.GetDocs()); err != nil {
									return err
								}
							}
						}
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetIndexes().GetIndexes()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Index{
									Name:        data.Get("name").(string),
									Gtype:       data.Get("gtype").(string),
									Expression:  data.Get("expression").(string),
									Docs:        data.Get("target_docs").(bool),
									Connections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Index{
								Name:        data.Get("name").(string),
								Gtype:       data.Get("gtype").(string),
								Expression:  data.Get("expression").(string),
								Docs:        data.Get("target_docs").(bool),
								Connections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetIndexes(ctx, &apipb.Indexes{Indexes: values}); err != nil {
							return err
						}
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						indexes := scheme.GetIndexes().GetIndexes()
						var index = -1
						for i, a := range indexes {
							if a.GetName() == data.Get("name") {
								index = i
								indexes[i] = &apipb.Index{
									Name:        data.Get("name").(string),
									Gtype:       data.Get("gtype").(string),
									Expression:  data.Get("expression").(string),
									Docs:        data.Get("target_docs").(bool),
									Connections: data.Get("target_connections").(bool),
								}
							}
						}
						if index >= 0 {
							removeIndex(index, indexes)
							if err := client.SetIndexes(ctx, &apipb.Indexes{Indexes: indexes}); err != nil {
								return err
							}
						}
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return false, err
						}
						values := scheme.GetIndexes()
						var has = false
						for _, a := range values.GetIndexes() {
							if a.GetName() == data.Id() {
								has = true
							}
						}
						return has, nil
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
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "unique name of the index",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"gtype": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"expression": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"trigger": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"target_docs": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
						"target_connections": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
					},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetTriggers().GetTriggers()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Trigger{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									Trigger:           data.Get("trigger").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Trigger{
								Name:              data.Get("name").(string),
								Gtype:             data.Get("gtype").(string),
								Expression:        data.Get("expression").(string),
								Trigger:           data.Get("trigger").(string),
								TargetDocs:        data.Get("target_docs").(bool),
								TargetConnections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetTriggers(ctx, &apipb.Triggers{Triggers: values}); err != nil {
							return err
						}
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						id := data.Id()
						for _, a := range scheme.GetTriggers().GetTriggers() {
							if a.GetName() == id {
								if err := data.Set("name", a.GetName()); err != nil {
									return err
								}
								if err := data.Set("gtype", a.GetGtype()); err != nil {
									return err
								}
								if err := data.Set("expression", a.GetExpression()); err != nil {
									return err
								}
								if err := data.Set("trigger", a.GetTrigger()); err != nil {
									return err
								}
								if err := data.Set("target_connections", a.GetTargetConnections()); err != nil {
									return err
								}
								if err := data.Set("target_docs", a.GetTargetDocs()); err != nil {
									return err
								}
							}
						}
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetTriggers().GetTriggers()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Trigger{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									Trigger:           data.Get("trigger").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Trigger{
								Name:              data.Get("name").(string),
								Gtype:             data.Get("gtype").(string),
								Expression:        data.Get("expression").(string),
								Trigger:           data.Get("trigger").(string),
								TargetDocs:        data.Get("target_docs").(bool),
								TargetConnections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetTriggers(ctx, &apipb.Triggers{Triggers: values}); err != nil {
							return err
						}
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						triggers := scheme.GetTriggers().GetTriggers()
						var index = -1
						for i, a := range triggers {
							if a.GetName() == data.Get("name") {
								index = i
								triggers[i] = &apipb.Trigger{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									Trigger:           data.Get("trigger").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}
						if index >= 0 {
							removeTrigger(index, triggers)
							if err := client.SetTriggers(ctx, &apipb.Triggers{Triggers: triggers}); err != nil {
								return err
							}
						}
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return false, err
						}
						values := scheme.GetTriggers()
						var has = false
						for _, a := range values.GetTriggers() {
							if a.GetName() == data.Id() {
								has = true
							}
						}
						return has, nil
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
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "unique name of the index",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"gtype": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"expression": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"target_docs": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
						"target_connections": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
					},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetConstraints().GetConstraints()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Constraint{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Constraint{
								Name:              data.Get("name").(string),
								Gtype:             data.Get("gtype").(string),
								Expression:        data.Get("expression").(string),
								TargetDocs:        data.Get("target_docs").(bool),
								TargetConnections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetConstraints(ctx, &apipb.Constraints{Constraints: values}); err != nil {
							return err
						}
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						id := data.Id()
						for _, a := range scheme.GetConstraints().GetConstraints() {
							if a.GetName() == id {
								if err := data.Set("name", a.GetName()); err != nil {
									return err
								}
								if err := data.Set("gtype", a.GetGtype()); err != nil {
									return err
								}
								if err := data.Set("expression", a.GetExpression()); err != nil {
									return err
								}
								if err := data.Set("target_connections", a.GetTargetConnections()); err != nil {
									return err
								}
								if err := data.Set("target_docs", a.GetTargetDocs()); err != nil {
									return err
								}
							}
						}
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetConstraints().GetConstraints()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								has = true
								values[i] = &apipb.Constraint{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}

						if !has {
							values = append(values, &apipb.Constraint{
								Name:              data.Get("name").(string),
								Gtype:             data.Get("gtype").(string),
								Expression:        data.Get("expression").(string),
								TargetDocs:        data.Get("target_docs").(bool),
								TargetConnections: data.Get("target_connections").(bool),
							})
						}
						if err := client.SetConstraints(ctx, &apipb.Constraints{Constraints: values}); err != nil {
							return err
						}
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						values := scheme.GetConstraints().GetConstraints()
						var index = -1
						for i, a := range values {
							if a.GetName() == data.Get("name") {
								index = i
								values[i] = &apipb.Constraint{
									Name:              data.Get("name").(string),
									Gtype:             data.Get("gtype").(string),
									Expression:        data.Get("expression").(string),
									TargetDocs:        data.Get("target_docs").(bool),
									TargetConnections: data.Get("target_connections").(bool),
								}
							}
						}
						if index >= 0 {
							removeConstraint(index, values)
							if err := client.SetConstraints(ctx, &apipb.Constraints{Constraints: values}); err != nil {
								return err
							}
						}
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return false, err
						}
						values := scheme.GetConstraints()
						var has = false
						for _, a := range values.GetConstraints() {
							if a.GetName() == data.Id() {
								has = true
							}
						}
						return has, nil
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
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "unique name of the index",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"method": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"expression": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "replace me",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"target_requests": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
						"target_responses": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "replace me",
						},
					},
					SchemaVersion:  0,
					MigrateState:   nil,
					StateUpgraders: nil,
					Create: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						authorizers := scheme.GetAuthorizers()
						data.SetId(data.Get("name").(string))
						var has = false
						for i, a := range authorizers.GetAuthorizers() {
							if a.GetName() == data.Get("name") {
								has = true
								authorizers.Authorizers[i] = &apipb.Authorizer{
									Name:            data.Get("name").(string),
									Method:          data.Get("method").(string),
									Expression:      data.Get("expression").(string),
									TargetRequests:  data.Get("target_requests").(bool),
									TargetResponses: data.Get("target_responses").(bool),
								}
							}
						}

						if !has {
							authorizers.Authorizers = append(authorizers.Authorizers, &apipb.Authorizer{
								Name:            data.Id(),
								Method:          data.Get("method").(string),
								Expression:      data.Get("expression").(string),
								TargetRequests:  data.Get("target_requests").(bool),
								TargetResponses: data.Get("target_responses").(bool),
							})
						}
						if err := client.SetAuthorizers(ctx, authorizers); err != nil {
							return err
						}
						return nil
					},
					Read: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						id := data.Id()
						for _, a := range scheme.GetAuthorizers().GetAuthorizers() {
							if a.GetName() == id {
								if err := data.Set("name", a.GetName()); err != nil {
									return err
								}
								if err := data.Set("expression", a.GetExpression()); err != nil {
									return err
								}
								if err := data.Set("method", a.GetMethod()); err != nil {
									return err
								}
								if err := data.Set("target_requests", a.GetTargetRequests()); err != nil {
									return err
								}
								if err := data.Set("target_responses", a.GetTargetResponses()); err != nil {
									return err
								}
							}
						}
						return nil
					},
					Update: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						authorizers := scheme.GetAuthorizers()
						var has = false
						for i, a := range authorizers.GetAuthorizers() {
							if a.GetName() == data.Get("name") {
								has = true
								authorizers.Authorizers[i] = &apipb.Authorizer{
									Name:            data.Get("name").(string),
									Method:          data.Get("method").(string),
									Expression:      data.Get("expression").(string),
									TargetRequests:  data.Get("target_requests").(bool),
									TargetResponses: data.Get("target_responses").(bool),
								}
							}
						}
						data.SetId(data.Get("name").(string))
						if !has {
							authorizers.Authorizers = append(authorizers.Authorizers, &apipb.Authorizer{
								Name:            data.Id(),
								Method:          data.Get("method").(string),
								Expression:      data.Get("expression").(string),
								TargetRequests:  data.Get("target_requests").(bool),
								TargetResponses: data.Get("target_responses").(bool),
							})
						}
						if err := client.SetAuthorizers(ctx, authorizers); err != nil {
							return err
						}
						return nil
					},
					Delete: func(data *schema.ResourceData, i interface{}) error {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return err
						}
						authorizers := scheme.GetAuthorizers().GetAuthorizers()
						var index = -1
						for i, a := range authorizers {
							if a.GetName() == data.Get("name") {
								index = i
								authorizers[i] = &apipb.Authorizer{
									Name:            data.Get("name").(string),
									Method:          data.Get("method").(string),
									Expression:      data.Get("expression").(string),
									TargetRequests:  data.Get("target_requests").(bool),
									TargetResponses: data.Get("target_responses").(bool),
								}
							}
						}
						if index >= 0 {
							removeAuthorizer(index, authorizers)
							if err := client.SetAuthorizers(ctx, &apipb.Authorizers{Authorizers: authorizers}); err != nil {
								return err
							}
						}
						return nil
					},
					Exists: func(data *schema.ResourceData, i interface{}) (bool, error) {
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						defer cancel()
						client := i.(*graphik.Client)
						scheme, err := client.GetSchema(ctx, &empty.Empty{})
						if err != nil {
							return false, err
						}
						authorizers := scheme.GetAuthorizers()
						var has = false
						for _, a := range authorizers.GetAuthorizers() {
							if a.GetName() == data.Id() {
								has = true
							}
						}
						return has, nil
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
				clientSecret := data.Get("client_secret").(string)
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
					ClientID:     clientId,
					ClientSecret: clientSecret,
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

func removeAuthorizer(i int, values []*apipb.Authorizer) {
	values[i] = values[len(values)-1]
	values[len(values)-1] = nil
	values = values[:len(values)-1]
}

func removeIndex(i int, values []*apipb.Index) {
	values[i] = values[len(values)-1]
	values[len(values)-1] = nil
	values = values[:len(values)-1]
}

func removeConstraint(i int, values []*apipb.Constraint) {
	values[i] = values[len(values)-1]
	values[len(values)-1] = nil
	values = values[:len(values)-1]
}

func removeTrigger(i int, values []*apipb.Trigger) {
	values[i] = values[len(values)-1]
	values[len(values)-1] = nil
	values = values[:len(values)-1]
}
