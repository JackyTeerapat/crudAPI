package researcher

import (
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

	degreeRows, err := h.db.Raw("SELECT id, degree_type, degree_program, degree_university FROM degree WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching degree dataL %v", err.Error())})
		return
	}

	defer degreeRows.Close()

	for degreeRows.Next() {
		var degree models.TempDegree_get
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
	programRows, err := h.db.Raw("SELECT id, program_name FROM program WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching program data: %v", err.Error())})
		return
	}

	defer programRows.Close()

	for programRows.Next() {
		var program models.TempProgram_get
		if err := programRows.Scan(&program.ProgramID, &program.ProgramName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning program data"})
			return
		}

		programs = append(programs, program)
	}
	researcher.Program = programs

	// Fetch and add TempExperience data
	var experiences []models.TempExperience_get
	experienceRows, err := h.db.Raw("SELECT id, experience_type, experience_start, experience_end, experience_university, experience_remark FROM experience WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching experience data: %v", err.Error())})
		return
	}

	defer experienceRows.Close()

	for experienceRows.Next() {
		var experience models.TempExperience_get
		if err := experienceRows.Scan(&experience.ExperienceID, &experience.ExperienceType, &experience.ExperienceStart, &experience.ExperienceEnd, &experience.ExperienceUniversity, &experience.ExperienceRemark); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning experience data"})
			return
		}

		experiences = append(experiences, experience)
	}
	researcher.Experience = experiences

	// Fetch and add TempAttach data
	var attaches []models.TempAttach_get
	attachRows, err := h.db.Raw("SELECT id, file_name, file_action, file_storage FROM profile_attach WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching attach data: %v", err.Error())})
		return
	}

	defer attachRows.Close()

	for attachRows.Next() {
		var attach models.TempAttach_get
		if err := attachRows.Scan(&attach.FileID, &attach.FileName, &attach.FileAction, &attach.FileStorage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning attach data"})
			return
		}

		attaches = append(attaches, attach)
	}
	researcher.Attach = attaches

	// Fetch and add TempExplore data
	var explores []models.TempExplore_get
	exploreRows, err := h.db.Raw("SELECT id, explore_name, explore_year, explore_detail FROM exploration WHERE profile_id = ?", id).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching explore data: %v", err.Error())})
		return
	}

	defer exploreRows.Close()

	for exploreRows.Next() {
		var explore models.TempExplore_get
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

func (h *ResearcherHandler) CreateResearcher(c *gin.Context) {
	var researcher models.ResearcherRequest

	if err := c.ShouldBindJSON(&researcher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create createdBy and updatedBy variables
	createdBy := "Champlnwza007"
	updatedBy := "Champlnwza007"
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
		if err := h.db.Exec("INSERT INTO degree (profile_id, degree_type, degree_program, degree_university, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
			profileID, degree.DegreeType, degree.DegreeProgram, degree.DegreeUniversity, createdBy, updatedBy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting degree data: %v", err.Error())})
			return
		}
	}

	// Save program data
	for _, program := range researcher.Program {
		if err := h.db.Exec("INSERT INTO program (profile_id, program_name, created_by, updated_by) VALUES (?, ?, ?, ?)",
			profileID, program.ProgramName, createdBy, updatedBy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting program data: %v", err.Error())})
			return
		}
	}

	// Save experience data
	for _, experience := range researcher.Experience {
		if err := h.db.Exec("INSERT INTO experience (profile_id, experience_type, experience_start, experience_end, experience_university, experience_remark, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			profileID, experience.ExperienceType, experience.ExperienceStart, experience.ExperienceEnd, experience.ExperienceUniversity, experience.ExperienceRemark, createdBy, updatedBy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting experience data: %v", err.Error())})
			return
		}
	}

	// Save explore data
	for _, explore := range researcher.Explore {
		if err := h.db.Exec("INSERT INTO exploration (profile_id, explore_name, explore_year, explore_detail, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)",
			profileID, explore.ExploreName, explore.ExploreYear, explore.ExploreDetail, createdBy, updatedBy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting explore data: %v", err.Error())})
			return
		}
	}

	// Save attach data
	for _, attach := range researcher.Attach {
		if err := h.db.Exec("INSERT INTO profile_attach (profile_id, file_name, file_action,file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?,?)",
			profileID, attach.FileName, attach.FileAction, attach.File_storage, createdBy, updatedBy).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while inserting attach data: %v", err.Error())})
			return
		}
	}
	c.JSON(http.StatusCreated, gin.H{"Suscess": fmt.Sprintf("Profile ID : %v Created", profileID)})
}


func (h *ResearcherHandler) UpdateResearcher(c *gin.Context) {
	id := c.Param("id")

	createdBy := "Champlnwza007"
	updatedBy := "Champlnwza007"

	// Parse the JSON payload from the request body into the ResearcherRequest struct
	var req models.ResearcherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start updating the profile table
	result := h.db.Exec("UPDATE profile SET first_name = ?, last_name = ?, university = ?, address_home = ?, address_work = ?, email = ?, phone_number = ?, position_id = ? WHERE id = ?", req.FirstName, req.LastName, req.University, req.AddressHome, req.AddressWork, req.Email, req.PhoneNumber, req.PositionID, id)

	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating the researcher data in profile: %v", result.Error)})
		return
	}

	// Start updating degree data
	h.db.Exec("DELETE FROM degree WHERE profile_id = ?", id)
	for _, degree := range req.Degree {
		h.db.Exec("INSERT INTO degree (profile_id, degree_type, degree_program, degree_university, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", id, degree.DegreeType, degree.DegreeProgram, degree.DegreeUniversity, createdBy, updatedBy)
	}

	// Start updating program data
	h.db.Exec("DELETE FROM program WHERE profile_id = ?", id)
	for _, program := range req.Program {
		h.db.Exec("INSERT INTO program (profile_id, program_name, created_by, updated_by) VALUES (?, ?, ?, ?)", id, program.ProgramName, createdBy, updatedBy)
	}

	// Start updating experience data
	h.db.Exec("DELETE FROM experience WHERE profile_id = ?", id)
	for _, experience := range req.Experience {
		h.db.Exec("INSERT INTO experience (profile_id, experience_type, experience_start, experience_end, experience_university, experience_remark, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", id, experience.ExperienceType, experience.ExperienceStart, experience.ExperienceEnd, experience.ExperienceUniversity, experience.ExperienceRemark, createdBy, updatedBy)
	}

	// Start updating explore data
	h.db.Exec("DELETE FROM exploration WHERE profile_id = ?", id)
	for _, explore := range req.Explore {
		h.db.Exec("INSERT INTO exploration (profile_id, explore_name, explore_year, explore_detail, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", id, explore.ExploreName, explore.ExploreYear, explore.ExploreDetail, createdBy, updatedBy)
	}

	// Start updating attach data
	h.db.Exec("DELETE FROM profile_attach WHERE profile_id = ?", id)
	for _, attach := range req.Attach {
		h.db.Exec("INSERT INTO profile_attach (profile_id, file_name, file_action, file_storage, created_by, updated_by) VALUES (?, ?, ?, ?, ?, ?)", id, attach.FileName, attach.FileAction, attach.File_storage, createdBy, updatedBy)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Researcher profile updated successfully."})
}

func (h *ResearcherHandler) DeleteResearcher(c *gin.Context) {
	id := c.Param("id")

	// Start deleting degree data
	result := h.db.Exec("DELETE FROM degree WHERE profile_id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting degree data for researcher with ID %s: %v", id, result.Error)})
		return
	}

	// Start deleting program data
	result = h.db.Exec("DELETE FROM program WHERE profile_id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting program data for researcher with ID %s: %v", id, result.Error)})
		return
	}

	// Start deleting experience data
	result = h.db.Exec("DELETE FROM experience WHERE profile_id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting experience data for researcher with ID %s: %v", id, result.Error)})
		return
	}

	// Start deleting exploration data
	result = h.db.Exec("DELETE FROM exploration WHERE profile_id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting exploration data for researcher with ID %s: %v", id, result.Error)})
		return
	}

	// Start deleting attachment data
	result = h.db.Exec("DELETE FROM profile_attach WHERE profile_id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting attachment data for researcher with ID %s: %v", id, result.Error)})
		return
	}

	// Finally, start deleting the profile
	result = h.db.Exec("DELETE FROM profile WHERE id = ?", id)
	if result.Error != nil {
		// Handle the error if the query fails
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while deleting the researcher profile: %v", result.Error)})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Researcher with ID %s has been deleted successfully.", id)})
}
