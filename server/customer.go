package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Customer struct {
	Id          string `json:Id`
	Name        string `json:Name`
	Notes       string `json:Notes`
	NotesSearch string `json:NotesSearch`
}

func CreateCustomer(c *gin.Context) {
	var customer Customer
	c.ShouldBindJSON(&customer)

	tenantId := c.GetString("tenantId")

	customerIdInt, _ := db.Incr(ctx, "CustomerId:"+tenantId).Result()
	customer.Id = strconv.FormatInt(customerIdInt, 10)

	customerKey := "cust:" + tenantId + ":" + customer.Id
	db.HSet(ctx, customerKey, map[string]interface{}{
		"Id":          customer.Id,
		"Name":        customer.Name,
		"Notes":       customer.Notes,
		"NotesSearch": customer.NotesSearch,
	})
	db.ZAdd(ctx, "custIndex:"+tenantId, &redis.Z{
		Score:  math.Round(float64(time.Now().UnixNano() / 1000 / 1000)),
		Member: customerKey,
	})

	c.JSON(200, gin.H{"id": customer.Id, "Customer": customer})
}

func UpdateCustomer(c *gin.Context) {
	var customer Customer
	c.ShouldBindJSON(&customer)
	tenantId := c.GetString("tenantId")
	customerId, _ := c.Params.Get("id")
	// JSON, _ := json.Marshal(customer)
	customerKey := "cust:" + tenantId + ":" + customerId
	db.HSet(ctx, customerKey, map[string]interface{}{
		"Id":          customer.Id,
		"Name":        customer.Name,
		"Notes":       customer.Notes,
		"NotesSearch": customer.NotesSearch,
	})
	db.ZAdd(ctx, "custIndex:"+tenantId, &redis.Z{
		Score:  math.Round(float64(time.Now().UnixNano() / 1000 / 1000)),
		Member: customerKey,
	})
	c.JSON(200, gin.H{"id": customerId, "Customer": customer})
}

func GetCustomersList(c *gin.Context) {
	tenantId := c.GetString("tenantId")
	customerIds, _ := db.ZRevRange(ctx, "custIndex:"+tenantId, 0, 100).Result()
	searchText := c.Query("SearchText")

	if searchText != "" {
		fmt.Println("Search Text Is", searchText)


		sdb := redisearch.NewClient("localhost:6379", "cust:"+tenantId+":Index")
		results, _, _ := sdb.Search(redisearch.NewQuery(searchText).Summarize("NotesSearch")) //.Highlight([]string{"NotesSearch"},"<b>","</b>")
		customers := make([]interface{},len(results),len(results))
		for i, item := range results {
			customers[i] = item.Properties;
		}
		c.JSON(200, customers)
		return;
	}

	vals := make([]*redis.StringStringMapCmd, 0, 10)
	pipe := db.TxPipeline()
	for _, key := range customerIds {
		vals = append(vals, pipe.HGetAll(ctx, key))
	}
	pipe.Exec(ctx)

	customers := make([]interface{}, 0, 10)
	for _, val := range vals {
		customers = append(customers, val.Val())
	}

	c.JSON(200, customers)
}

func GetCustomer(c *gin.Context) {
	tenantId := c.GetString("tenantId")
	customerId, _ := c.Params.Get("id")

	customer, _ := db.HGetAll(ctx, "cust:"+tenantId+":"+customerId).Result()
	c.JSON(200, customer)
}
