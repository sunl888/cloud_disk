package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wq1019/cloud_disk/errors"
	"github.com/zm-dev/gerrors"
)

func NewHandleErrorMiddleware(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // execute all the handlers

		// at this point, all the handlers finished. Let's read the errors!
		// in this example we only will use the **last error typed as public**
		// but you could iterate over all them since c.Errors is a slice!
		errorToPrint := c.Errors.Last()
		if errorToPrint != nil {
			var ge *gerrors.GlobalError

			switch errorToPrint.Err {
			case gorm.ErrRecordNotFound:
				ge = errors.NotFound(errorToPrint.Err.Error()).(*gerrors.GlobalError)
			default:
				ge = &gerrors.GlobalError{}
				if json.Unmarshal([]byte(errorToPrint.Err.Error()), ge) != nil {
					ge = errors.InternalServerError(errorToPrint.Err.Error(), errorToPrint.Err).(*gerrors.GlobalError)
				}
			}

			if ge.ServiceName == "" {
				ge.ServiceName = serviceName
			}
			c.JSON(ge.StatusCode, gin.H{
				"code":    ge.Code,
				"message": ge.Message,
			})
		}

	}
}
