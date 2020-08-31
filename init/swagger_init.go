package init

import "post-api/docs"

func Swagger() {
	docs.SwaggerInfo.Title = "Swagger POST API"
	docs.SwaggerInfo.Description = "This is Gola POST API Server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
}
