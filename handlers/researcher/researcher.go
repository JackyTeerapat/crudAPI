package researcher

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResearcherHandler struct {
	db *gorm.DB
}

func NewResearcherHandler(db *gorm.DB) *ResearcherHandler {
	return &ResearcherHandler{db: db}
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
	var researcher models.Researcher_get

	// Execute the query and scan the results into the researcher struct
	result := h.db.Raw("SELECT id as profile_id, first_name, last_name, university, address_home, address_work, email,phone_number FROM profile WHERE id = ?", id).Scan(&researcher)

	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching the researcher data from profile: %v", result.Error)})

		return
	}

	// Start getData from Degree
	var degrees []models.TempDegree_get

	degreeRows, err := h.db.Raw("SELECT id, degree_type, degree_program, degree_university,activated FROM degree WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching degree dataL %v", err.Error())})
		return
	}

	defer degreeRows.Close()

	for degreeRows.Next() {
		var degree models.TempDegree_get
		if err := degreeRows.Scan(&degree.DegreeID, &degree.DegreeType, &degree.DegreeProgram, &degree.DegreeUniversity, &degree.Activated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning degree data"})
			return
		}

		degrees = append(degrees, degree)
	}
	// END Data from Degree

	// Add the degrees data to the researcher struct
	researcher.Degree = degrees

	//start
	// Fetch and add TempPosition data
	var positionID int
	errPositionID := h.db.Raw("SELECT position_id FROM profile WHERE id = ? ", id).Scan(&positionID)

	if errPositionID.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching the researcher data from profile: %v", result.Error)})

		return
	}
	var positions []models.TempPosition_get
	positionRows, err := h.db.Raw("SELECT id, position_name FROM position WHERE id = ?", positionID).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching position data: %v", err.Error())})
		return
	}

	defer positionRows.Close()

	for positionRows.Next() {
		var position models.TempPosition_get
		if err := positionRows.Scan(&position.PositionID, &position.PositionName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning position data"})
			return
		}

		positions = append(positions, position)
	}
	researcher.Position = positions

	// Fetch and add TempProgram data
	var programs []models.TempProgram_get
	programRows, err := h.db.Raw("SELECT id, program_name,activated FROM program WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching program data: %v", err.Error())})
		return
	}

	defer programRows.Close()

	for programRows.Next() {
		var program models.TempProgram_get
		if err := programRows.Scan(&program.ProgramID, &program.ProgramName, &program.Activated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning program data"})
			return
		}

		programs = append(programs, program)
	}
	researcher.Program = programs

	// Fetch and add TempExperience data
	var experiences []models.TempExperience_get
	experienceRows, err := h.db.Raw("SELECT id, experience_type, experience_start, experience_end, experience_university, experience_remark,activated FROM experience WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching experience data: %v", err.Error())})
		return
	}

	defer experienceRows.Close()

	for experienceRows.Next() {
		var experience models.TempExperience_get
		if err := experienceRows.Scan(&experience.ExperienceID, &experience.ExperienceType, &experience.ExperienceStart, &experience.ExperienceEnd, &experience.ExperienceUniversity, &experience.ExperienceRemark, &experience.Activated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning experience data"})
			return
		}

		experiences = append(experiences, experience)
	}
	researcher.Experience = experiences

	// Fetch and add TempAttach data
	var attaches []models.TempAttach_get
	attachRows, err := h.db.Raw("SELECT id, file_name, file_action,activated FROM profile_attach WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching attach data: %v", err.Error())})
		return
	}

	defer attachRows.Close()

	for attachRows.Next() {
		var attach models.TempAttach_get
		if err := attachRows.Scan(&attach.FileID, &attach.FileName, &attach.FileAction, &attach.Activated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning attach data"})
			return
		}

		attaches = append(attaches, attach)
	}
	researcher.Attach = attaches

	// Fetch and add TempExplore data
	var explores []models.TempExplore_get
	exploreRows, err := h.db.Raw("SELECT id, explore_name, explore_year, explore_detail,activated FROM exploration WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching explore data: %v", err.Error())})
		return
	}

	defer exploreRows.Close()

	for exploreRows.Next() {
		var explore models.TempExplore_get
		if err := exploreRows.Scan(&explore.ExploreID, &explore.ExploreName, &explore.ExploreYear, &explore.ExploreDetail, &explore.Activated); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning explore data"})
			return
		}

		explores = append(explores, explore)
	}
	researcher.Explore = explores

	res := api.ResponseApi(http.StatusOK, researcher, nil)
	// Return the researcher data as JSON
	c.JSON(http.StatusOK, res)
}

func (h *ResearcherHandler) CreateResearcher(c *gin.Context) {
	var researcher models.ResearcherRequest

	if err := c.ShouldBindJSON(&researcher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create createdBy and updatedBy variables
	createdBy := "Champlnwza007"
	updatedBy := "Champlnwza007"
	activated := true
	// Update the INSERT statement for the profile table
	result := h.db.Exec("INSERT INTO profile (first_name, last_name, university, address_home, address_work, email, phone_number, position_id, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id",
		researcher.FirstName, researcher.LastName, researcher.University, researcher.AddressHome, researcher.AddressWork, researcher.Email, researcher.PhoneNumber, researcher.PositionID, createdBy, updatedBy)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting researcher data into the profile table: %v", result.Error)})
		return
	}
	var profile models.Profile
	h.db.Last(&profile)
	profileID := profile.ID

	// The rest of the code remains the same until the INSERT statements for the other tables

	// Save degree data
	for _, degree := range researcher.Degree {
		if err := h.db.Exec("INSERT INTO degree (profile_id, degree_type, degree_program, degree_university, created_by, updated_by,activated) VALUES (?, ?, ?, ?, ?, ?,?)",
			profileID, degree.DegreeType, degree.DegreeProgram, degree.DegreeUniversity, createdBy, updatedBy, activated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting degree data: %v", err.Error())})
			return
		}
	}

	// Save program data
	for _, program := range researcher.Program {
		if err := h.db.Exec("INSERT INTO program (profile_id, program_name, created_by, updated_by,activated) VALUES (?, ?, ?, ?,?)",
			profileID, program.ProgramName, createdBy, updatedBy, activated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting program data: %v", err.Error())})
			return
		}
	}

	// Save experience data
	for _, experience := range researcher.Experience {
		if err := h.db.Exec("INSERT INTO experience (profile_id, experience_type, experience_start, experience_end, experience_university, experience_remark, created_by, updated_by,activated) VALUES (?, ?, ?, ?, ?, ?, ?, ?,?)",
			profileID, experience.ExperienceType, experience.ExperienceStart, experience.ExperienceEnd, experience.ExperienceUniversity, experience.ExperienceRemark, createdBy, updatedBy, activated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting experience data: %v", err.Error())})
			return
		}
	}

	// Save explore data
	for _, explore := range researcher.Explore {
		if err := h.db.Exec("INSERT INTO exploration (profile_id, explore_name, explore_year, explore_detail, created_by, updated_by,activated) VALUES (?, ?, ?, ?, ?, ?,?)",
			profileID, explore.ExploreName, explore.ExploreYear, explore.ExploreDetail, createdBy, updatedBy, activated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting explore data: %v", err.Error())})
			return
		}
	}

	// Save attach data
	for _, attach := range researcher.Attach {
		if err := h.db.Exec("INSERT INTO profile_attach (profile_id, file_name, file_action,file_storage, created_by, updated_by, activated) VALUES (?, ?, ?, ?, ?,?,?)",
			profileID, attach.FileName, attach.FileAction, attach.File_storage, createdBy, updatedBy, activated).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting attach data: %v", err.Error())})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"Suscess": fmt.Sprintf("Profile ID : %v Created", profileID)})
}
func (h *ResearcherHandler) VSdeleteResearcher(c *gin.Context) {
	profileID := c.Param("id")
	updatedBy := "Champlnwza007"
	activated := false

	tablesToUpdate := []string{"degree", "program", "experience", "exploration", "profile_attach"}

	for _, tableName := range tablesToUpdate {
		if err := h.db.Exec(fmt.Sprintf("UPDATE %s SET updated_by = ?, activated = ? WHERE profile_id = ?", tableName), updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating the activated status in the %s table: %v", tableName, err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("Profile ID : %v deactivated", profileID)})

}
func (h *ResearcherHandler) UpdateResearcher(c *gin.Context) {
	var researcher models.ResearcherRequest
	profileID := c.Param("id")

	if err := c.ShouldBindJSON(&researcher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update createdBy and updatedBy variables
	updatedBy := "Champlnwza007"
	activated := true

	// Update the researcher's profile data
	if err := h.db.Exec("UPDATE profile SET first_name = ?, last_name = ?, university = ?, address_home = ?, address_work = ?, email = ?, phone_number = ?, position_id = ?, updated_by = ? WHERE id = ?",
		researcher.FirstName, researcher.LastName, researcher.University, researcher.AddressHome, researcher.AddressWork, researcher.Email, researcher.PhoneNumber, researcher.PositionID, updatedBy, profileID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating researcher data in the profile table: %v", err)})
		return
	}

	// Use the given profile ID for the rest of the update operations

	// Update degree data
	for _, degree := range researcher.Degree {
		if err := h.db.Exec("UPDATE degree SET degree_type = ?, degree_program = ?, degree_university = ?, updated_by = ?, activated = ? WHERE profile_id = ?",
			degree.DegreeType, degree.DegreeProgram, degree.DegreeUniversity, updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating degree data: %v", err)})
			return
		}
	}

	// Update program data
	for _, program := range researcher.Program {
		if err := h.db.Exec("UPDATE program SET program_name = ?, updated_by = ?, activated = ? WHERE profile_id = ?",
			program.ProgramName, updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating program data: %v", err)})
			return
		}
	}

	// Update experience data
	for _, experience := range researcher.Experience {
		if err := h.db.Exec("UPDATE experience SET experience_type = ?, experience_start = ?, experience_end = ?, experience_university = ?, experience_remark = ?, updated_by = ?, activated = ? WHERE profile_id = ?",
			experience.ExperienceType, experience.ExperienceStart, experience.ExperienceEnd, experience.ExperienceUniversity, experience.ExperienceRemark, updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating experience data: %v", err)})
			return
		}
	}

	// Update explore data
	for _, explore := range researcher.Explore {
		if err := h.db.Exec("UPDATE exploration SET explore_name = ?, explore_year = ?, explore_detail = ?, updated_by = ?, activated = ? WHERE profile_id = ?",
			explore.ExploreName, explore.ExploreYear, explore.ExploreDetail, updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating explore data: %v", err)})
			return
		}
	}

	// Update attach data
	for _, attach := range researcher.Attach {
		if err := h.db.Exec("UPDATE profile_attach SET file_name = ?, file_action = ?, file_storage = ?, updated_by = ?, activated = ? WHERE profile_id = ?",
			attach.FileName, attach.FileAction, attach.File_storage, updatedBy, activated, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating attach data: %v", err)})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"Success": fmt.Sprintf("Profile ID : %v Updated", profileID)})
}
