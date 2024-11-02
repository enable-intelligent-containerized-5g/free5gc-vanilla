package analyticsinfo

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
	group := engine.Group("/nnwdaf-analyticsinfo/v1")

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

	// subPatternShort := "/subscription-data/:ueId"
	// group.Any(subPatternShort, subMsgShortDispatchHandlerFunc)
	//
	// subPattern := "/subscription-data/:ueId/:servingPlmnId"
	// group.Any(subPattern, subMsgDispatchHandlerFunc)
	//
	// eePatternShort := "/subscription-data/:ueId/:servingPlmnId/ee-subscriptions"
	// group.Any(eePatternShort, eeMsgShortDispatchHandlerFunc)
	//
	// eePattern := "/subscription-data/:ueId/:servingPlmnId/ee-subscriptions/:subsId"
	// group.Any(eePattern, eeMsgDispatchHandlerFunc)

	/*
	 * GIN wildcard issue:
	 * '/application-data/influenceData/:influenceId' and
	 * '/application-data/influenceData/subs-to-notify' patterns will be conflicted.
	 * Only can use '/application-data/influenceData/:influenceId' pattern and
	 * use a dispatch handler to distinguish "subs-to-notify" from ":influenceId".
	 */
	// appInfluDataPattern := "/application-data/influenceData/:influenceId"
	// group.Any(appInfluDataPattern, appInfluDataMsgDispatchHandlerFunc)
	//
	// expoPatternShort := "/exposure-data/:ueId/:subId"
	// group.Any(expoPatternShort, expoMsgDispatchHandlerFunc)
	//
	// expoPattern := "/exposure-data/:ueId/:subId/:pduSessionId"
	// group.Any(expoPattern, expoMsgDispatchHandlerFunc)

	return group
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"",
		Index,
	},

	{
		"NwdafAnalyticsInfoRequest",
		strings.ToUpper("Post"),
		"analyticsinfo/request",
		HTTPNwdafAnalyticsInfoRequest,
	},
}
