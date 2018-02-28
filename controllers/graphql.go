package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
	"net/http"

	"fmt"
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/graphql"
)

type graphqlController struct {
	*BaseController
}

func newGraphqlController() *graphqlController {
	controller := &graphqlController{
		BaseController: NewBaseController(
			models.SharedGraphqlModel(),
			logics.UniqueGraphqlLogic(),
		),
	}
	controller.SetBinder(controller)
	return controller
}

func (controller *graphqlController) InitializeEarly(r *gin.Engine) error {
	graphqlTypes := extensions.RegisteredGraphqlTypes()

	for _, graphqlType := range graphqlTypes {
		fields := graphqlType.BuildTypeFieldsEarly()
		inputObjectConfigFieldMap := graphqlType.BuildInputTypeFieldsEarly()
		fieldConfigArguments := graphqlType.BuildTypeArguments()
		fmt.Println(fields)
		graphqlType.BuildTypeObject(fields, inputObjectConfigFieldMap, fieldConfigArguments)
	}

	for _, graphqlType := range graphqlTypes {
		fields := graphqlType.BuildTypeFieldsLate()
		inputObjectConfigFieldMap := graphqlType.BuildInputTypeFieldsLate()
		fmt.Println(fields)
		graphqlType.UpdateTypeObject(fields, inputObjectConfigFieldMap)
	}

	return nil
}

func (controller *graphqlController) InitializeLate(r *gin.Engine) error {
	return nil
}

func (controller *graphqlController) ResourceSingleURL() string {
	return "graphql"
}

func (controller *graphqlController) Bind(c *gin.Context, container interface{}) error {
	buffer := make([]byte, c.Request.ContentLength, c.Request.ContentLength)
	_, err := c.Request.Body.Read(buffer)

	if err != nil {
		logging.Logger().Debug(err.Error())
		//return err
	}

	container.(*models.Graphql).Query = string(buffer)
	container.(*models.Graphql).Mutation = string(buffer)

	return nil
}

func (controller *graphqlController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodPost: {
			extensions.URLSingle: controller.Execute,
		},
	}
	return routeMap
}

func (controller *graphqlController) Execute(c *gin.Context) {
	graphqlModel := &models.Graphql{}

	if err := controller.Bind(c, graphqlModel); err != nil {
		logging.Logger().Debug(err.Error())
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	a := graphql.UniqueDesignGraphqlType()
	b := a.BuiltTypeObject()

	fmt.Println("A1")
	fmt.Println(a)
	fmt.Println("A2")
	fmt.Println(b)

	var schema, _ = graphqlGo.NewSchema(
		graphqlGo.SchemaConfig{
			Query:    b,
			Mutation: b,
		},
	)

	result := graphqlGo.Do(
		graphqlGo.Params{
			Schema:        schema,
			RequestString: graphqlModel.Query,
			Context:       c,
		},
	)

	fmt.Println("A3")
	fmt.Println(result)

	c.JSON(http.StatusOK, result)
}

func (controller *graphqlController) Mutation(c *gin.Context) {
}

var uniqueGraphqlController = newGraphqlController()

func init() {
	extensions.RegisterController(uniqueGraphqlController)
	extensions.RegisterRouterInitializer(uniqueGraphqlController)
}
