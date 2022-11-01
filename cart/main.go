package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

func main() {

	r := gin.New()
	r.POST("/api/cart", func(c *gin.Context) {

		var inventory Inventory

		if err := c.ShouldBindJSON(&inventory); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "400",
				"error":  err.Error(),
			})
			return
		}

		postInventory(inventory)

		c.JSON(http.StatusCreated, gin.H{
			"status":  "201",
			"message": "add to cart",
		})
	})

	r.Run(":8650")
}

const topic = "PRODUCT.order"

type Inventory struct {
	ProductNumber string  `json:"productNumber"`
	Name          string  `json:"name"`
	Price         float32 `json:"price"`
}

func postInventory(i Inventory) {

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	maps := make(map[string]interface{})

	maps["productNumber"] = i.ProductNumber
	maps["name"] = i.Name
	maps["price"] = i.Price

	marshal, err := json.Marshal(maps)

	if err != nil {
		log.Fatal(err)
	}

	nc.Publish(topic, []byte(marshal))
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", topic, marshal)
	}

}
