package httpresponse

import (
	"net/http"
	"octlink/ovs/utils/merrors"
	"octlink/ovs/utils/uuid"

	"github.com/gin-gonic/gin"
)

func BuildErrorObj(ctx *gin.Context, code int, errlog interface{},
	data interface{}) map[string]interface{} {

	return gin.H{
		"errorObj": gin.H{
			"errorNo":  code,
			"errorLog": errlog,
			"errorMsg": merrors.GetMsg(code),
		},
		"apiId": uuid.Generate().Simple(),
		"data":  data,
	}
}

// RHttprespnse retrun none error code 200
func Ok(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, BuildErrorObj(ctx, merrors.ErrSuccess, nil, data))
	return
}

// RHttprespnse retrun none error code 200
func Error(ctx *gin.Context, err int, errlog interface{}) {
	ctx.JSON(http.StatusOK, BuildErrorObj(ctx, err, errlog, nil))
	return
}

// RHttprespnse retrun none error code 201
func Create(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusCreated, gin.H{"code": merrors.ErrSuccess, "data": nil})
	return
}

// RHttprespnse retrun none error code 204
func Delete(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusNoContent, gin.H{"code": merrors.ErrSuccess, "data": data})
	return
}

// RHttprespnse retrun none error code 202
func Update(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusAccepted, gin.H{"code": merrors.ErrSuccess, "data": data})
	return
}
