package gql

import (
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
)

func GraphQLAuthMiddleware(c *gin.Context) {
	currentUser := c.Value("user").(*models.User)
	if currentUser.ID == "" {
		err := NewGqlError(common.ErrorUnauthorised, "unauthorised")
		c.JSON(http.StatusOK, gin.H{"errors": gqlerrors.FormattedErrors{
			gqlerrors.FormatError(gqlerrors.Error{
				Message:       err.Error(),
				Locations:     []location.SourceLocation{},
				OriginalError: err,
			}),
		}})
		c.Abort()
		return
	}
}
