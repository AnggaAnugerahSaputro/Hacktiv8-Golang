package helpers

import (
	"strconv"
	"assignment_2/pkg/errs"
	"github.com/gin-gonic/gin"
)

func GetParamId(ctx *gin.Context, key string) (int, errs.MessageErr) {
	value := ctx.Param(key)

	param, err := strconv.Atoi(value)
	if err != nil {
		return 0, errs.NewBadRequest("invalid uri parameter")
	}

	return param, nil
}
