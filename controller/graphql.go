package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"net/http"

	"fmt"
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/graphql"
	"github.com/qb0C80aE/clay/extension"
)

type graphqlController struct {
	BaseController
}

func newGraphqlController() *graphqlController {
	return CreateController(&graphqlController{}, model.NewGraphql()).(*graphqlController)
}

func (receiver *graphqlController) DoBeforeRouterSetup(r *gin.Engine) error {
	graphqlTypes := extension.RegisteredGraphqlTypes()

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

func (receiver *graphqlController) GetResourceSingleURL() string {
	return "graphql"
}

func (receiver *graphqlController) Bind(c *gin.Context, container interface{}) error {
	buffer := make([]byte, c.Request.ContentLength, c.Request.ContentLength)
	_, err := c.Request.Body.Read(buffer)

	if err != nil {
		logging.Logger().Debug(err.Error())
		//return err
	}

	container.(*model.Graphql).Query = string(buffer)
	container.(*model.Graphql).Mutation = string(buffer)

	return nil
}

func (receiver *graphqlController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodPost: {
			extension.URLSingle: receiver.Execute,
		},
	}
	return routeMap
}

func (receiver *graphqlController) Execute(c *gin.Context) {
	graphqlModel := &model.Graphql{}

	if err := receiver.Bind(c, graphqlModel); err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	a := graphql.UniqueDesignGraphqlType()
	b := a.BuiltTypeObject()

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

	c.JSON(http.StatusOK, result)
}

func (controller *graphqlController) Mutation(c *gin.Context) {
}

var uniqueGraphqlController = newGraphqlController()

func init() {
	extension.RegisterController(uniqueGraphqlController)
	extension.RegisterInitializer(uniqueGraphqlController)
}
