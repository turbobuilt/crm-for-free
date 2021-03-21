package main

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Tenant struct {
	Id         string `json:Id`
	AdminEmail string `json:AdminEmail`
}

func NewTenant(c *gin.Context) {
	tenant := Tenant{}
	c.ShouldBind(&tenant)

	if tenant.AdminEmail == "" {
		c.JSON(400, gin.H{"error": "AdminEmail is required"})
		return
	}

	// check for email first
	userKey, _ := db.Get(ctx, "email:"+tenant.AdminEmail).Result()
	if userKey != "" {
		c.JSON(400, gin.H{"error": "That email is already registered"})
		return
	}

	idNumber, err := db.Incr(ctx, "TenantId").Result()
	if err != nil {
		panic(err)
	}
	id := strconv.FormatInt(idNumber, 10)

	tenantJSON := ""
	json.Marshal(&tenantJSON)
	db.Set(ctx, "tenant:"+id, string(tenantJSON), 0)

}
