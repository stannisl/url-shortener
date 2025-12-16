package handlers

import (
	"github.com/wb-go/wbf/ginext"
)

func SetError(c *ginext.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
