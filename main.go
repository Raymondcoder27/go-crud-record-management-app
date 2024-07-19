package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type record struct{
	ID string `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
	DocumentNumber string `json:"document_number"`
	Quantity int `json:"quantity"`
}

var records = []record{
	{ID: "1", Title: "Title", Body: "Body", DocumentNumber: "Document number",Quantity: 3},
	{ID: "2", Title: "Title 2", Body: "Body 2", DocumentNumber: "Document number 2",Quantity: 4},
	{ID: "3", Title: "Title 3", Body: "Body 6", DocumentNumber: "Document number 3",Quantity: 8},
}


func getAllRecords(c *gin.Context){
	c.IndentedJSON(http.StatusOK, records)
}

func createNewRecord(c *gin.Context){
	//create a new record variable of type record
	var newRecord record

	//bind the body json to the variable
	if err := c.BindJSON(&newRecord); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "record not created"})
		return
	}

	//add the record to the records slice 
	records = append(records, newRecord)

	//return the response
	c.IndentedJSON(http.StatusCreated, newRecord)
}

//helper function to retrieve a record by the id
func getRecordById(id string)(*record, error){
	//loop through the available records to retrieve the id at the current position
	for i, record := range records {
		if record.ID == id {
			return &records[i], nil
		}
	}
	return nil, errors.New("the record was not found")
}

//retrieve the record by the id using the helper function
func recordById(c *gin.Context){
	//create a variable for the id
	id := c.Param("id")

	//retrieve the record
	record, err := getRecordById(id)
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"The record was not found."})
	}

	c.IndentedJSON(http.StatusOK, record)
}

//delete a record (removing a record from the record slice)
func deleteRecord(c *gin.Context){
	//get the id of the record
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Wrong or missing querry parameter."})
		return
	} 

	//retrieve the record
	record, err := getRecordById(id)
	if err != nil{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"The id was not found"})
	}

	//check if the quantity is more than 0
	if record.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"no records found."})
		return
	}

	//delete the record from the slice 
	record.Quantity -= 1

	//return with the new record object
	c.IndentedJSON(http.StatusOK, record)
}

//return a record
func returnRecord(c *gin.Context){
	//get the id of the record from the query parameter
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message":"The id was not found."})
	}

	//retrieve the record
	record, err := getRecordById(id)
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message":"The record was not found"})
	}

	//add the record to the records slice
	record.Quantity += 1

	//return with the new record object
	c.IndentedJSON(http.StatusOK, record)
}

func main() {
	r := gin.Default()
	r.GET("/records", getAllRecords)
	r.POST("/records", createNewRecord)
	r.GET("/records/:id", recordById)
	r.PATCH("/returnRecord", returnRecord)
	r.PATCH("/deleteRecord", deleteRecord)
	r.Run("localhost:8080")
}
