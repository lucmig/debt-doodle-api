package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Sample - debt sample point
type Sample struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float32   `json:"value"`
}

// Debt - a debt
type Debt struct {
	ID          string    `bson:"_id" json:"_id,omitempty"`
	Description string    `json:"description"`
	Samples     []Sample  `json:"samples"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// database instance
var collection *mongo.Collection

// DebtCollection - mongo db collection for debts
func DebtCollection(c *mongo.Database) {
	collection = c.Collection("debts")
}

// GetAllDebts - all debts
func GetAllDebts(c *gin.Context) {
	debts := []Debt{}
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		log.Printf("Error while getting all debts, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	// Iterate through the returned cursor.
	for cursor.Next(context.TODO()) {
		var debt Debt
		cursor.Decode(&debt)
		debts = append(debts, debt)
	}

	c.JSON(http.StatusOK, debts)
	return
}

// AddDebt - add a debt
func AddDebt(c *gin.Context) {
	var debt Debt
	c.BindJSON(&debt)
	description := debt.Description
	samples := []Sample{}
	id := debt.ID

	newDebt := Debt{
		ID:          id,
		Description: description,
		Samples:     samples,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := collection.InsertOne(context.TODO(), newDebt)

	if err != nil {
		log.Printf("Error while inserting new debt into db, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "Debt created Successfully",
	})
	return
}

// GetDebt - get a bedt by id
func GetDebt(c *gin.Context) {
	debtID := c.Param("debtId")

	debt := Debt{}
	err := collection.FindOne(context.TODO(), bson.M{"_id": debtID}).Decode(&debt)
	if err != nil {
		log.Printf("Error while getting a single debt, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Debt not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Debt",
		"data":    debt,
	})
	return
}

// EditDebt - edit a single debt selected by id
func EditDebt(c *gin.Context) {
	debtID := c.Param("debtId")
	var debt Debt
	c.BindJSON(&debt)

	newData := bson.M{
		"$set": bson.M{
			"description": debt.Description,
			"updated_at":  time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": debtID}, newData)
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Debt Edited Successfully",
	})
	return
}

// DeleteDebt - delete a single debt selected by id
func DeleteDebt(c *gin.Context) {
	debtID := c.Param("debtId")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": debtID})
	if err != nil {
		log.Printf("Error while deleting a single debt, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Debt deleted successfully",
	})
	return
}

// UpdateSample - add or update debt sample point
func UpdateSample(c *gin.Context) {
	debtID := c.Param("debtId")
	var sample Sample
	c.BindJSON(&sample)

	oldData := bson.M{
		"$pull": bson.M{
			"samples": bson.M{
				"timestamp": sample.Timestamp,
			},
		},
	}

	samples := []Sample{
		Sample{
			Timestamp: sample.Timestamp,
			Value:     sample.Value,
		},
	}

	newData := bson.M{
		"$push": bson.M{
			"samples": bson.M{
				"$each": samples,
				"$sort": bson.M{
					"timestamp": -1,
				},
			},
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": debtID}, oldData)
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": err,
		})
		return
	}

	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": debtID}, newData)
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": res,
	})
	return
}

// DeleteSample - delete debt sample point
func DeleteSample(c *gin.Context) {
	debtID := c.Param("debtId")
	timestamp, _ := time.Parse(time.RFC3339, c.Param("timestamp"))

	data := bson.M{
		"$pull": bson.M{
			"samples": bson.M{
				"timestamp": timestamp,
			},
		},
	}

	res, err := collection.UpdateOne(context.TODO(), bson.M{"_id": debtID}, data)
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": res,
	})
	return
}
