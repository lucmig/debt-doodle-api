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
	Date  string  `json:"date" time_format:"2006-01-02" time_utc:"1"`
	Value float32 `json:"value"`
}

//Debt Json request payload is as follows,
//{
//  "_id": "1",
//  "description": "Blue Credit Card",
//  "samples": [],
//  "created_at": "2020-04-01T00:00:00Z",
//  "updated_at": "2020-04-01T00:00:00Z"
//}

// A ValidationError is an error that is used when the required input fails validation.
// swagger:response validationError
type ValidationError struct {
	// The error message
	// in: body
	Body struct {
		// The validation message
		//
		// Required: true
		// Example: Expected type int
		Message string
		// An optional field name to which this validation applies
		FieldName string
	}
}

// Debt is a debt item
//
// A debt
//
// swagger:parameters AddDebt
type Debt struct {
	// in:body
	// Required: true
	// Example: Expected type string
	ID          string    `bson:"_id" json:"_id,omitempty"`
	Description string    `json:"description"`
	Samples     []Sample  `json:"samples"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DebtValue - a debt value
type DebtValue struct {
	ID    string `bson:"_id" json:"_id,omitempty"`
	Value float32
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

	c.JSON(http.StatusOK, debt)
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
				"date": sample.Date,
			},
		},
	}

	samples := []Sample{
		{
			Date:  sample.Date,
			Value: sample.Value,
		},
	}

	newData := bson.M{
		"$push": bson.M{
			"samples": bson.M{
				"$each": samples,
				"$sort": bson.M{
					"date": -1,
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
	layout := "2006-01-02"
	sampleDate, err := time.Parse(layout, c.Param("date"))
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": err,
		})
		return
	}

	data := bson.M{
		"$pull": bson.M{
			"samples": bson.M{
				"date": sampleDate.Format("2006-01-02"),
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

// GetValuesDate - get debts value at a date
func GetValuesDate(c *gin.Context) {
	date := c.Param("date")

	debtValues := []DebtValue{}
	matchStage := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "samples.date", Value: bson.D{
				{Key: "$lte", Value: date},
			}},
		}},
	}
	setStage := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "value", Value: bson.D{
				{Key: "$max", Value: bson.D{
					{Key: "$filter", Value: bson.D{
						{Key: "input", Value: "$samples"},
						{Key: "as", Value: "smp"},
						{Key: "cond", Value: bson.D{
							{Key: "$lte", Value: bson.A{"$$smp.date", date}},
						}},
					}},
				}},
			}},
		}},
	}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "value", Value: "$value.value"},
		}},
	}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, setStage, projectStage})

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
		var debt DebtValue
		cursor.Decode(&debt)
		debtValues = append(debtValues, debt)
	}

	c.JSON(http.StatusOK, debtValues)
	return
}
