package mlmodelinfo

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nwdaf/internal/logger"
	logger_util "github.com/free5gc/util/logger"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := logger_util.NewGinWithLogrus(logger.GinLog)
	AddService(router)
	return router
}

func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("/nnwdaf-mlmodelinfo/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "PATCH":
			group.PATCH(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

// Index is the index handler.
func IndexLoader(c *gin.Context) {
	// Devolver el texto como un archivo .txt
	c.Header("Content-Disposition", "attachment; filename=loader.txt")
	c.Data(http.StatusOK, "text/plain", []byte("loaderio-f0107ec5e785d516e4effa15f7e96699"))
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/",
		Index,
	},

	{
		"Index",
		"GET",
		"/loaderio-f0107ec5e785d516e4effa15f7e96699.txt/",
		IndexLoader,
	},

	{
		"NwdafMlModelInfoRequest",
		strings.ToUpper("Get"),
		"mlmodelinfo/request",
		HTTPNwdafMlModelInfoRequest,
	},
}
