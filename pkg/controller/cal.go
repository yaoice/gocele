package controller

import (
	"github.com/RichardKnop/machinery/v1"
	mconfig "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/gin-gonic/gin"
	"github.com/yaoice/gocele/pkg/config"
	"github.com/yaoice/gocele/pkg/log"
	"github.com/yaoice/gocele/pkg/util"
	"net/http"
	"sync"
)

var (
	once   sync.Once
	calC   *calController
)

type CalInterface interface {
	Add(c *gin.Context)
	Mul(c *gin.Context)
	GetTask(c *gin.Context)
	RegisterTasks(map[string]interface{})
	CreateWorker(string) *machinery.Worker
}

type calController struct {
	server *machinery.Server
}

type numbers struct {
	Numbers string `json:"numbers" form:"numbers" query:"numbers"`
}

type uuid struct {
	UUID string `json:"uuid" form:"uuid" query:"uuid"`
}

func NewCalController() *calController {
	once.Do(func() {
		prefix := "machinery"
		broker := config.GetString(prefix + ".broker")
		resultBackend := config.GetString(prefix + ".result_backend")
		exchange := config.GetString(prefix + ".exchange")
		exchangeType := config.GetString(prefix + ".exchange_type")
		defaultQueue := config.GetString(prefix + ".default_queue")
		bindingKey := config.GetString(prefix + ".binding_key")

		cnf := mconfig.Config{
			Broker:        broker,
			ResultBackend: resultBackend,
			AMQP: &mconfig.AMQPConfig{
				Exchange:     exchange,
				ExchangeType: exchangeType,
				BindingKey:   bindingKey,
			},
			DefaultQueue: defaultQueue,
		}

		server, err := machinery.NewServer(&cnf)
		if err != nil {
			log.Fatalf("Could not initialize server %v", err.Error())
		}

		calC = &calController{
			server: server,
		}
	})
	return calC
}

// @Summary Add
// @Description Add测试
// @Param   body   body    numbers  true     "The numbers"
// @Success 200 {string} json "" OK
// @Failure 500 {string} json "" Internal Server Error
// @router /add [post]
func (this *calController) Add(c *gin.Context) {
	var message interface{}
	nums := numbers{}
	//var args []signatures.TaskArg

	if err := c.Bind(&nums); err != nil {
		return
	}

	nbrs := utils.SAtoI(nums.Numbers)

	args := []tasks.Arg{}

	for _, v := range nbrs {
		args = append(args, tasks.Arg{Type: "int64", Value: v})
	}

	signature := tasks.Signature{
		Name: "add",
		Args: args,
	}

	asyncResult, err := this.server.SendTask(&signature)
	if err != nil {
		log.Error("Could not send task", err)
	}

	result, err := asyncResult.GetWithTimeout(5000000000, 1)
	if err != nil { // Handle errors reading the config file
		taskState := asyncResult.GetState()
		c.String(http.StatusInternalServerError, "Defered! %s", taskState.TaskUUID)
		return
	}
	for _, r := range result {
		message = r.Int()
	}
	c.String(http.StatusOK, "Result: %v", message)
}

// @Summary Mul
// @Description Mul测试
// @Param   body   body    numbers  true     "The numbers"
// @Success 200 {string} json "" OK
// @Failure 500 {string} json "" Internal Server Error
// @router /mul [post]
func (this *calController) Mul(c *gin.Context) {
	var message interface{}
	nums := numbers{}
	//var args []signatures.TaskArg

	if err := c.Bind(&nums); err != nil {
		return
	}

	nbrs := utils.SAtoI(nums.Numbers)

	args := []tasks.Arg{}

	for _, v := range nbrs {
		args = append(args, tasks.Arg{Type: "int64", Value: v})
	}

	signature := tasks.Signature{
		Name: "multiply",
		Args: args,
	}

	asyncResult, err := this.server.SendTask(&signature)
	if err != nil {
		log.Error("Could not send task", err)
	}

	result, err := asyncResult.GetWithTimeout(5000000000, 1)
	if err != nil { // Handle errors reading the config file
		taskState := asyncResult.GetState()
		c.String(http.StatusInternalServerError, "Defered! %s", taskState.TaskUUID)
		return
	}
	for _, r := range result {
		message = r.Int()
	}
	c.String(http.StatusOK, "Result: %v", message)
}

// @Summary Tasks
// @Description Tasks测试
// @Param   body   body    uuid  true     "The uuid"
// @Success 200 {string} json "" OK
// @Failure 500 {string} json "" Internal Server Error
// @router /tasks [post]
func (this *calController) GetTask(c *gin.Context) {
	var message interface{}
	u := uuid{}

	if err := c.Bind(&u); err != nil {
		return
	}

	task := u.UUID
	taskState, err := this.server.GetBackend().GetState(task)

	if err != nil {
		c.String(http.StatusBadRequest, "Error: Task not found")

	}

	if taskState.State == "PENDING" {
		c.String(http.StatusOK, "Status: %s", tasks.StatePending)
		return
	}

	if taskState.State == "RECEIVED" {
		c.String(http.StatusOK, "Status: %s", tasks.StateReceived)
		return
	}

	if taskState.State == "STARTED" {
		c.String(http.StatusOK, "Status: %s",tasks.StateStarted)
		return
	}

	if taskState.State == "FAILURE" {
		c.String(http.StatusOK, "Status: ",tasks.StateFailure)
		return
	}

	if taskState.State == "SUCCESS" {
		for _, r := range taskState.Results {
			message = r.Value
		}
		c.String(http.StatusOK, "Status: %s\nResult: %v", tasks.StateSuccess, message)
		return
	}
	c.String(http.StatusBadRequest, "ERROR: Something broken")
}

func (this *calController) RegisterTasks(tasks map[string]interface{}) {
	this.server.RegisterTasks(tasks)
}

func (this *calController) CreateWorker(ctag string) *machinery.Worker {
	return this.server.NewWorker(ctag, 1)
}
