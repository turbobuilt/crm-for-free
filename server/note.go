package main

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Note struct {
	CreatedBy string `json:CreatedBy`
	UpdatedAt string `json:UpdatedAt`
	Note      string `json:Note`
	Customer  string `json:Customer`
}

func SaveNewNote(c *gin.Context) {
	note := Note{}
	c.ShouldBind(&note)
	note.UpdatedAt = GetUnixMillisecondsString()

	tenantId := "1"

	cmd := db.Incr(ctx, "UserId:"+tenantId)
	id, _ := cmd.Result()

	key := tenantId + ":note:" + strconv.FormatInt(id, 10)
	bodyData, _ := json.Marshal(note)
	db.Set(ctx, key, bodyData, 0)
}

func UpdateNote(c *gin.Context) {

}
