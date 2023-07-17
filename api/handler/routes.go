package handler

import (
	"biling-nats/api/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HandlerApi struct {
	Nats   *nats.Conn
	Logger *log.Logger
}

func NewHandlerApi(l *log.Logger, nats *nats.Conn) *HandlerApi {
	return &HandlerApi{Logger: l, Nats: nats}
}

//----------------------------------------------------------------------------

func (h *HandlerApi) SignUp(c *gin.Context) {
	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid JSON provided",
		})
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "failed to marshal the data",
		})
		h.Logger.Println("failed to marshal the data. Error is: ", err.Error())
		return
	}

	request, err := h.Nats.Request("user.create", bytes, 5*time.Second)
	if err != nil {
		h.Logger.Println("failed to send a request. Error is: ", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": string(request.Data),
	})
}

func (h *HandlerApi) Login(c *gin.Context) {
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]any{
			"error": "invalid JSON provided",
		})
		return
	}

	bytes, err := json.Marshal(user)
	if err != nil {
		h.Logger.Println(err)
		return
	}

	request, err := h.Nats.Request("user.login", bytes, 5*time.Second)
	if err != nil {
		h.Logger.Println("failed to send a request. Error is: ", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": string(request.Data),
	})
}

//----------------------------------------------------------------------------

func (h *HandlerApi) CreateAccount(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": err.Error(),
		})
		return
	}

	var account *model.Account
	if err := c.BindJSON(&account); err != nil {
		c.JSON(400, map[string]any{
			"error": "invalid JSON provided",
		})
		return
	}

	if account.Id == "" || account.Name == "" {
		c.JSON(400, map[string]any{
			"error": "invalid id or name provided",
		})
		return
	}
	account.UserId = userId
	account.Balance = 0

	//fmt.Println(account.UserId)

	bytes, err := json.Marshal(account)
	if err != nil {
		c.JSON(500, map[string]any{
			"error": "couldn't marshal the data",
		})
		h.Logger.Println(err)
		return
	}

	request, err := h.Nats.Request("account.create", bytes, 5*time.Second)
	if err != nil {
		c.JSON(500, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send a request. Error is:", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": string(request.Data),
	})
}

func (h *HandlerApi) GetAccounts(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": err.Error(),
		})
		return
	}

	//mapId := make(map[string]int)
	//mapId["userId"] = userId
	idStr := strconv.Itoa(userId)

	//bytes, err := json.Marshal(idStr)
	//if err != nil {
	//	c.JSON(400, map[string]any{
	//		"error": "failed to marshal the data",
	//	})
	//	return
	//}
	//fmt.Println(string(bytes))

	request, err := h.Nats.Request("account.get", []byte(idStr), 5*time.Second)
	if err != nil {
		c.JSON(500, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send a request. Error is:", err.Error())
		return
	}

	var accounts model.Accounts
	if err := json.Unmarshal(request.Data, &accounts); err != nil {
		c.JSON(400, map[string]any{
			"error": "could not unmarshal the data",
		})
		h.Logger.Println("error while unmarshalling. error is: ", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"accounts": accounts,
	})
}

func (h *HandlerApi) Transaction(c *gin.Context) {
	var transaction model.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON provided",
		})
		return
	}

	bytes, err := json.Marshal(transaction)
	if err != nil {
		c.JSON(500, map[string]any{
			"error": "couldn't marshal the data",
		})
		h.Logger.Println(err)
		return
	}

	request, err := h.Nats.Request("account.transfer", bytes, 5*time.Second)
	if err != nil {
		c.JSON(500, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send a request. Error is:", err.Error())
		return
	}

	id, err := strconv.Atoi(string(request.Data))
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "could not convert to integer",
		})
		h.Logger.Println("failed while converting to integer. Error is:", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": id,
	})
}

//----------------------------------------------------------------------------

func (h *HandlerApi) AddBook(c *gin.Context) {
	userRole, err := GetUserRole(c)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": err.Error(),
		})
		return
	}

	if userRole != "seller" {
		c.JSON(400, map[string]any{
			"error": "sign up as a seller",
		})
		h.Logger.Println("you are not allowed to add a book, sign up as a seller")
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": err.Error(),
		})
		return
	}

	var book model.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(400, map[string]any{
			"error": "could not unmarshal the data",
		})
		return
	}
	book.SellerId = userId

	bytes, err := json.Marshal(book)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "could not marshal the data",
		})
		h.Logger.Println(err.Error())
		return
	}

	request, err := h.Nats.Request("book.add", bytes, 5*time.Second)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send the request. Error is:", err.Error())
		return
	}

	id, err := strconv.Atoi(string(request.Data))
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "could not convert to integer",
		})
		h.Logger.Println("could not convert to integer. Error is:", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": id,
	})
}

type BuyBook struct {
	Id      int
	BuyerId int
}

func (h *HandlerApi) BuyBook(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": err.Error(),
		})
		return
	}

	bookId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "invalid param provided",
		})
		return
	}

	var book = BuyBook{
		Id:      bookId,
		BuyerId: userId,
	}

	bytes, err := json.Marshal(book)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "couldn't marshal the data",
		})
		return
	}

	request, err := h.Nats.Request("book.buy", bytes, 10*time.Second)
	if err != nil {
		c.JSON(400, map[string]any{
			"error": "failed to send the request",
		})
		h.Logger.Println("failed to send the request. Error is:", err.Error())
		return
	}

	c.JSON(200, map[string]any{
		"message": string(request.Data),
	})
}
