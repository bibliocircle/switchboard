package auth

import (
	"net/http"
	"switchboard/internal/common"
	"switchboard/internal/gql"
	"switchboard/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/location"
)

func GraphQLAuthMiddleware(c *gin.Context) {
	currentUser := c.Value("user").(*user.User)
	if currentUser.ID == "" {
		err := gql.NewGqlError(common.ErrorUnauthorised, "unauthorised")
		c.JSON(http.StatusOK, gin.H{
			"errors": gqlerrors.FormattedErrors{
				gqlerrors.FormatError(gqlerrors.Error{
					Message:       err.Error(),
					Locations:     []location.SourceLocation{},
					OriginalError: err,
				}),
			},
		})
		c.Abort()
		return
	}
}
