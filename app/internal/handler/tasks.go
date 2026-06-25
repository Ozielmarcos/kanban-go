package handler

import (
	"net/http"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/internal/repository"
	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	var tasks []model.Task
	if err := c.BindJSON(&tasks); err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Não foi encontrado tarefas!"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTasksByUser(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	tasks, err := repository.GetTasksByUser(userId)

	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message": "Não foi encontrado tarefas!", "Erro": err.Error()},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTaskById(c *gin.Context) {

	userId := c.MustGet("user_id").(string)
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	task, err := repository.GetTaskById(id, userId)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message": "Tarefa não encontrada!", "Erro": err.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task model.Task
	c.BindJSON(&task)

	if task.StoryId == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Necessário ID da story!"})
		return
	}

	if task.Title == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Título é obrigatório!"})
		return
	}

	response, err := repository.CreateTask(task)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Erro ao criar tarefa!", "Error": err.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusCreated, response)
}

func UpdateTask(c *gin.Context) {
	var task model.Task
	c.BindJSON(&task)

	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	err := repository.UpdateTask(id, task)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Erro ao atualizar tarefa!", "Erro": err.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Tarefa atualizada com sucesso!"})
}

func UpdateTaskStatus(c *gin.Context) {
	var body struct {
		Status model.TaskStatus `json:"status"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Payload inválido!"})
		return
	}

	id := c.Param("id")

	err := repository.UpdateTaskStatus(id, body.Status)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar a tarefa!", "erro": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Status atualizado com sucesso."})
}

func StarTimerHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	err := repository.StarTimer(id)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Erro ao iniciar timer!", "Erro": err.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Timer iniciado com sucesso!"})
}

func PauseTimerHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID não informado!"})
	}

	err := repository.PauseTimer(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Timer pausado com sucesso!"})
}

func ResumeTimerHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	err := repository.ResumeTimer(id)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Timer continua"})
}
