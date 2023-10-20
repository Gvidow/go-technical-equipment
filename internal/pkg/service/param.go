package service

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrMissingID = errors.New("required id")

func FetchIdFromURLPath(c *gin.Context) (int, error) {
	strID := c.Param("id")
	if strID == "" {
		return 0, ErrMissingID
	}
	int64ID, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse id from url path: %w", err)
	}
	return int(int64ID), nil
}
