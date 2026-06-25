package handler

import (
	"fmt"
	"net/http"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/internal/repository"
	"github.com/gin-gonic/gin"
)

func CreateStory(c *gin.Context) {
	var story model.Story
	c.BindJSON(&story)

	userId := c.MustGet("user_id").(string)
	story.UserId = userId

	if story.Title == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Título é obrigatório!"})
		return
	}

	response, err := repository.CreateStory(story)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao criar história!"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "História criada com sucesso!", "data": response})
}

func GetStoriesByUser(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	fmt.Print(userId)
	stories, err := repository.GetStoriesByUser(userId)
	fmt.Print(stories)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Não foi encontrado histórias!"})
		return
	}
	c.IndentedJSON(http.StatusOK, stories)
}

func GetTasksByStory(c *gin.Context) {
	storyId := c.Param("story_id")

	tasks, err := repository.GetStoriesTask(storyId)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Não foram encontradas tarefas vinculadas à esse projeto!"})
		return
	}
	c.IndentedJSON(http.StatusOK, tasks)
}

func UpdateStory(c *gin.Context) {
	var story model.Story
	c.BindJSON(&story)

	userId := c.MustGet("user_id").(string)
	story.UserId = userId

	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	err := repository.UpdateStory(id, story)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao atualizar história!"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "História atualizada com sucesso!"})
}

func RemoveStory(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID inválido!"})
		return
	}

	err := repository.RemoveStory(id, userId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Erro ao remover história!"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "História removida com sucesso!"})
}
