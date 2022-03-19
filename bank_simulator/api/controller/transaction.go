package controller

import (
	"errors"
	"fmt"
	"net/http"

	creditcard "github.com/durango/go-credit-card"
	"github.com/gin-gonic/gin"
)

const (
	statusApproved string = "approved"
	statusFailed   string = "failed"
	cardType       string = "credit_card"
)

func TransactionControllerFunc(c *gin.Context) {
	var headers Headers
	if err := c.ShouldBindHeader(&headers); err != nil {
		c.Error(err)
		httpError := MapErrors[2]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		httpError := MapErrors[9]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	var (
		card creditcard.Card = creditcard.Card{
			Number: req.Number,
			Cvv:    req.Cvv,
			Month:  req.ExpiryMonth,
			Year:   req.ExpiryYear,
		}
	)

	var testCard = TestCards[card.Number]
	if testCard > 0 {
		httpError := MapErrors[testCard]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	if err := card.ValidateExpiration(); err != nil {
		c.Error(err)
		httpError := MapErrors[8]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	lastFour, err := card.LastFour()
	if err != nil {
		c.Error(err)
		httpError := MapErrors[7]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	if !card.ValidateNumber() {
		c.Error(errors.New("invalid number"))
		httpError := MapErrors[7]
		c.AbortWithStatusJSON(httpError.StatusCode, gin.H{
			"code":    httpError.Code,
			"message": httpError.Message,
			"status":  statusFailed,
		})
		return
	}

	c.JSON(http.StatusCreated, &Response{
		Amount:   req.Amount,
		Currency: req.Currency,
		Code:     "00",
		Message:  "ok",
		Status:   statusApproved,
		Source: Source{
			ExpiryMonth: req.ExpiryMonth,
			ExpiryYear:  req.ExpiryYear,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Number:      fmt.Sprintf("XXXX-XXXX-XXXX-%s", lastFour),
			Bin:         req.Number[0:5],
			CardType:    cardType,
		},
	})
}
