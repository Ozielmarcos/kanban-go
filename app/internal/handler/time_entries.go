package handler

import (
	"net/http"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/internal/repository"
	"github.com/gin-gonic/gin"
)

func GetTimeEntriesByStoryHandler(c *gin.Context) {
	storyId := c.Param("story_id")

	if storyId == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID do projeto não informado!"})
		return
	}

	entries, err := repository.GetTimeEntriesByStory(storyId)

	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Não foi encontrado apontamentos!", "erro": err},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, entries)
}

func UpdateEntryHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "ID não informado!"},
		)
		return
	}

	var entry model.TimeEntry

	if err := c.ShouldBindJSON(&entry); err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Faltando body do apontamento!", "erro": err.Error()},
		)
		return
	}

	if err := repository.UpdateEntry(id, entry); err != nil {
		if err.Error() == "apontamento não encontrado!" {
			c.IndentedJSON(
				http.StatusNotFound,
				gin.H{"message": err.Error()},
			)
			return
		}

		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Erro ao atualizar apontamento!", "erro": err.Error()},
		)
		return
	}

	c.IndentedJSON(
		http.StatusOK,
		gin.H{"message": "Apontamento atualizado com sucesso."},
	)
}

func DeleteEntryHandler(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "ID não informado!"})
		return
	}

	if err := repository.RemoveEntry(id); err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Erro ao excluir apontamento!", "erro": err.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Apontamento excluido com sucesso!"})
}
