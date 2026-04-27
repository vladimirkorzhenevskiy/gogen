package dto

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestApp_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     App
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "success",
			filename: "testdata/gogen.yaml",
			want: App{
				Metadata: Metadata{
					Name:   "example",
					Module: "example",
				},
				Spec: Spec{
					Drivers: Drivers{
						{
							Name: "postgres",
							Kind: "postgres",
						},
						{
							Name: "redis",
							Kind: "redis",
						},
					},
					Controllers: Controllers{
						{
							Component: Component{
								Name:      "consumer",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "updater",
										FQDN: "usecases.leadership.update",
									},
								},
							},
							Kind:     "consumer",
							Protocol: "kafka",
						},
						{
							Component: Component{
								Name:      "scheduler",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "snapshotter",
										FQDN: "usecases.leadership.snapshot",
									},
								},
							},
							Kind:     "scheduler",
							Protocol: "cron",
						},
						{
							Component: Component{
								Name:      "server",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "getter",
										FQDN: "usecase.leadership.get",
									},
								},
							},
							Kind:     "server",
							Protocol: "http",
						},
					},
					UseCases: UseCases{
						{
							Component: Component{
								Name:      "get",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "cache",
										FQDN: "repository.leadership.redis",
									},
									{
										Name: "repository",
										FQDN: "repository.leadership.postgres",
									},
								},
							},
						},
						{
							Component: Component{
								Name:      "snapshot",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "cache",
										FQDN: "repository.leadership.redis",
									},
									{
										Name: "repository",
										FQDN: "repository.leadership.postgres",
									},
								},
							},
						},
						{
							Component: Component{
								Name:      "update",
								Namespace: "leadership",
								Dependencies: Dependencies{
									{
										Name: "cache",
										FQDN: "repository.leadership.redis",
									},
									{
										Name: "repository",
										FQDN: "repository.leadership.postgres",
									},
								},
							},
						},
					},
					Repositories: Repositories{
						{
							Component: Component{
								Name:      "postgres",
								Namespace: "leadership",
								Config: Config{
									{
										Name:     "timeout",
										Type:     "duration",
										Required: false,
										Default:  "10s",
									},
								},
								Dependencies: Dependencies{
									{
										Name: "postgres",
										FQDN: "drivers.postgres",
									},
								},
							},
						},
						{
							Component: Component{
								Name:      "redis",
								Namespace: "leadership",
								Config: Config{
									{
										Name:     "timeout",
										Type:     "duration",
										Required: false,
										Default:  "10s",
									},
								},
								Dependencies: Dependencies{
									{
										Name: "redis",
										FQDN: "drivers.redis",
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Fatal(err)
			}

			{
				var app App

				tt.wantErr(t, yaml.Unmarshal(data, &app))
				assert.Equal(t, tt.want.Metadata, app.Metadata)
				assert.Equal(t, tt.want.Spec.Entrypoints, app.Spec.Entrypoints)
				assert.Equal(t, tt.want.Spec.Drivers, app.Spec.Drivers)
				assert.Equal(t, tt.want.Spec.Controllers, app.Spec.Controllers)
				assert.Equal(t, tt.want.Spec.Middlewares, app.Spec.Middlewares)
				assert.Equal(t, tt.want.Spec.UseCases, app.Spec.UseCases)
				assert.Equal(t, tt.want.Spec.Repositories, app.Spec.Repositories)
				assert.Equal(t, tt.want.Spec.Adapters, app.Spec.Adapters)

				model, err := app.ToModel()
				if err != nil {
					t.Fatal(err)
				}

				_ = model
			}
		})
	}
}
