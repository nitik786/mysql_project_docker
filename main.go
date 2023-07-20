package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Project struct {
	ID           int64  `json:"id"`
	ProjectName  string `json:"project_name"`
	ProjectOwner string `json:"project_owner"`
}

var db *sql.DB

func main() {
	// Connect to MySQL database
	var err error
	db, err = sql.Open("mysql", "your_mysql_user:your_mysql_password@tcp(mysql:3306)/your_mysql_db_name")
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	// Set up Gin router
	r := gin.Default()

	// Routes
	r.GET("/projects", GetProjects)
	r.POST("/projects", CreateProject)
	r.GET("/projects/:id", GetProjectByID)
	r.PUT("/projects/:id", UpdateProject)
	r.DELETE("/projects/:id", DeleteProject)

	// Run the server
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Error running the server: ", err)
	}
}

// GetProjects retrieves all projects from the database
func GetProjects(c *gin.Context) {
	var projects []Project

	rows, err := db.Query("SELECT ID, ProjectName, ProjectOwner FROM Project")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving projects"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var project Project
		err := rows.Scan(&project.ID, &project.ProjectName, &project.ProjectOwner)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning projects"})
			return
		}
		projects = append(projects, project)
	}

	c.JSON(http.StatusOK, projects)
}

// CreateProject creates a new project in the database
func CreateProject(c *gin.Context) {
	var project Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	result, err := db.Exec("INSERT INTO Project (ProjectName, ProjectOwner) VALUES (?, ?)",
		project.ProjectName, project.ProjectOwner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating project"})
		return
	}

	project.ID, _ = result.LastInsertId()
	c.JSON(http.StatusCreated, project)
}

// GetProjectByID retrieves a project by its ID
func GetProjectByID(c *gin.Context) {
	id := c.Param("id")

	var project Project
	err := db.QueryRow("SELECT ID, ProjectName, ProjectOwner FROM Project WHERE ID = ?", id).Scan(&project.ID, &project.ProjectName, &project.ProjectOwner)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject updates a project by its ID
func UpdateProject(c *gin.Context) {
	id := c.Param("id")

	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	_, err := db.Exec("UPDATE Project SET ProjectName = ?, ProjectOwner = ? WHERE ID = ?",
		project.ProjectName, project.ProjectOwner, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProject deletes a project by its ID
func DeleteProject(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM Project WHERE ID = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting project"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
