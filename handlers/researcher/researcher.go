package researcher

import (
	"net/http"

	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResearcherHandler struct {
	db *gorm.DB
}

func NewResearcherHandler(db *gorm.DB) *ResearcherHandler {
	return &ResearcherHandler{db: db}
}

type Researcher struct {
	ProfileID   int              `json:"profile_id"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	University  string           `json:"university"`
	AddressHome string           `json:"address_home"`
	AddressWork string           `json:"address_work"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Degree      []tempDegree     `gorm:"-" json:"degree"`
	Position    []tempPosition   `gorm:"-" json:"position"`
	Program     []tempProgram    `gorm:"-" json:"program"`
	Experience  []tempExperience `gorm:"-" json:"experience"`
	Attach      []tempAttach     `gorm:"-" json:"attach"`
	Explore     []tempExplore    `gorm:"-" json:"explore"`
}
type tempDegree struct {
	DegreeID         int    `json:"id"`
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
}
type tempPosition struct {
	PositionID   int    `json:"position_id"`
	PositionName string `json:"position_name"`
}
type tempProgram struct {
	ProgramID   int    `json:"program_id"`
	ProgramName string `json:"program_name"`
}

type tempExperience struct {
	ExperienceID         int    `json:"experience_id"`
	ExperienceType       string `json:"experience_type"`
	ExperienceStart      string `json:"experience_start"`
	ExperienceEnd        string `json:"experience_end"`
	ExperienceUniversity string `json:"experience_university"`
	ExperienceRemark     string `json:"experience_remark"`
}

type tempAttach struct {
	FileID      int    `json:"file_id"`
	FileName    string `json:"file_name"`
	FileAction  string `json:"file_action"`
	FileStorage string `json:"file_storage"`
}

type tempExplore struct {
	ExploreID     int    `json:"explore_id"`
	ExploreName   string `json:"explore_name"`
	ExploreYear   string `json:"explore_year"`
	ExploreDetail string `json:"explore_detail"`
}

// researcher/profile_detail godoc
// @Summary Get a Profile_detail
// @Description Get a data profile_detail from database.
// @Tags Profile detail
// @Produce  application/json
// @Param id path int true "researcher ProfileID"
// @Success 200 {object} Researcher{}
// @Router /researcher/profile_detail/{id} [get]
func (h *ResearcherHandler) ListResearcher(c *gin.Context) {
	id := c.Param("id")

	// Start getDate from Profile
	var researcher Researcher

	// Execute the query and scan the results into the researcher struct
	result := h.db.Raw("SELECT id as profile_id, first_name, last_name, university, address_home, address_work, email,phone_number FROM profile WHERE id = ?", id).Scan(&researcher)

	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching the researcher data from profile: %v", result.Error)})

		return
	}

	// Start getData from Degree
	var degrees []tempDegree

	degreeRows, err := h.db.Raw("SELECT id, degree_type, degree_program, degree_university FROM degree WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching degree dataL %v", err.Error())})
		return
	}

	defer degreeRows.Close()

	for degreeRows.Next() {
		var degree tempDegree
		if err := degreeRows.Scan(&degree.DegreeID, &degree.DegreeType, &degree.DegreeProgram, &degree.DegreeUniversity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning degree data"})
			return
		}

		degrees = append(degrees, degree)
	}
	// END Data from Degree

	// Add the degrees data to the researcher struct
	researcher.Degree = degrees

	//start
	// Fetch and add tempPosition data
	var positions []tempPosition
	positionRows, err := h.db.Raw("SELECT id, position_name FROM position WHERE id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching position data: %v", err.Error())})
		return
	}

	defer positionRows.Close()

	for positionRows.Next() {
		var position tempPosition
		if err := positionRows.Scan(&position.PositionID, &position.PositionName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning position data"})
			return
		}

		positions = append(positions, position)
	}
	researcher.Position = positions

	// Fetch and add tempProgram data
	var programs []tempProgram
	programRows, err := h.db.Raw("SELECT id, program_name FROM program WHERE id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching program data: %v", err.Error())})
		return
	}

	defer programRows.Close()

	for programRows.Next() {
		var program tempProgram
		if err := programRows.Scan(&program.ProgramID, &program.ProgramName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning program data"})
			return
		}

		programs = append(programs, program)
	}
	researcher.Program = programs

	// Fetch and add tempExperience data
	var experiences []tempExperience
	experienceRows, err := h.db.Raw("SELECT id, experience_type, experience_start, experience_end, experience_university, experience_remark FROM experience WHERE id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching experience data: %v", err.Error())})
		return
	}

	defer experienceRows.Close()

	for experienceRows.Next() {
		var experience tempExperience
		if err := experienceRows.Scan(&experience.ExperienceID, &experience.ExperienceType, &experience.ExperienceStart, &experience.ExperienceEnd, &experience.ExperienceUniversity, &experience.ExperienceRemark); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning experience data"})
			return
		}

		experiences = append(experiences, experience)
	}
	researcher.Experience = experiences

	// Fetch and add tempAttach data
	var attaches []tempAttach
	attachRows, err := h.db.Raw("SELECT id, file_name, file_action, file_storage FROM profile_attach WHERE id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching attach data: %v", err.Error())})
		return
	}

	defer attachRows.Close()

	for attachRows.Next() {
		var attach tempAttach
		if err := attachRows.Scan(&attach.FileID, &attach.FileName, &attach.FileAction, &attach.FileStorage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning attach data"})
			return
		}

		attaches = append(attaches, attach)
	}
	researcher.Attach = attaches

	// Fetch and add tempExplore data
	var explores []tempExplore
	exploreRows, err := h.db.Raw("SELECT id, explore_name, explore_year, explore_detail FROM exploration WHERE id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching explore data: %v", err.Error())})
		return
	}

	defer exploreRows.Close()

	for exploreRows.Next() {
		var explore tempExplore
		if err := exploreRows.Scan(&explore.ExploreID, &explore.ExploreName, &explore.ExploreYear, &explore.ExploreDetail); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning explore data"})
			return
		}

		explores = append(explores, explore)
	}
	researcher.Explore = explores

	// Return the researcher data as JSON
	c.JSON(http.StatusOK, researcher)
}
