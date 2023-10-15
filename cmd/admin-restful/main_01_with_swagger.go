//go:build swagger_docs_enabled

package main

import (
	"fmt"

	_ "github.com/josestg/swe-be-mono/cmd/admin-restful/swagger-docs"
	"github.com/josestg/swe-be-mono/internal/app/adminrestful"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/swaggo/swag"
)

// init initializes the swagger documentation.
//
//	@title						Swagger Documentation for Admin REST API.
//	@version					1.0
//	@description				This is the swagger documentation for Admin REST API.
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.url				http://www.swagger.io/support
//	@license.name				Apache 2.0
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func init() {
	fmt.Println("swagger_docs_enabled: started")
	defer fmt.Println("swagger_docs_enabled: finished")
	spec, ok := swag.GetSwagger("enduser").(*swag.Spec)
	if !ok {
		panic("forgotten to import admin swagger docs?")
	}
	spec.Version = buildVersion
	spec.BasePath = adminrestful.BasePath
	adminrestful.SetDocHandler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.Layout(httpSwagger.BaseLayout),
		httpSwagger.InstanceName("admin"),
	))
}
