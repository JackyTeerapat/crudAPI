package researcher

import (
	"CRUD-API/api"
	"CRUD-API/models"
	"fmt"
	"log"
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
	result := h.db.Raw("SELECT id as profile_id, profile_status, first_name, last_name, university, address_home, address_work, email,phone_number FROM profile WHERE id = ?", id).Scan(&researcher)

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
	var positions []models.Position
	positionRows, err := h.db.Raw("SELECT id, position_name FROM position WHERE id = ?", positionID).Rows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching position data: %v", err.Error())})
		return
	}

	defer positionRows.Close()

	for positionRows.Next() {
		var position models.Position
		if err := positionRows.Scan(&position.ID, &position.Position_name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while scanning position data"})
			return
		}

		positions = append(positions, position)
	}
	researcher.PrefixName = positions[0].Position_name
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

	res := api.ResponseApiWithDescription(http.StatusOK, researcher, "SUCCESS", nil)
	c.JSON(http.StatusOK, res)
}

func (h *ResearcherHandler) ListResearcherbyID(id int) models.Researcher_get {

	// Start getDate from Profile
	var researcher models.Researcher_get

	// Execute the query and scan the results into the researcher struct
	result := h.db.Raw("SELECT id as profile_id, profile_status, first_name, last_name, university, address_home, address_work, email,phone_number FROM profile WHERE id = ?", id).Scan(&researcher)

	if result.Error != nil {
		// Handle the error if the query fails
		log.Printf("An error occurred while fetching the researcher data from profile: %v", result.Error)
	}

	// Start getData from Degree
	var degrees []models.TempDegree_get

	degreeRows, err := h.db.Raw("SELECT id, degree_type, degree_program, degree_university FROM degree WHERE profile_id = ?", id).Rows()
	if err != nil {
		log.Printf("An error occurred while fetching degree dataL %v", err.Error())
	}

	defer degreeRows.Close()

	for degreeRows.Next() {
		var degree models.TempDegree_get
		if err := degreeRows.Scan(&degree.DegreeID, &degree.DegreeType, &degree.DegreeProgram, &degree.DegreeUniversity); err != nil {
			log.Printf("error: An error occurred while scanning degree data")
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
		log.Printf("An error occurred while fetching the researcher data from profile: %v", result.Error)

	}
	var positions []models.Position
	positionRows, err := h.db.Raw("SELECT id, position_name FROM position WHERE id = ?", positionID).Rows()
	if err != nil {
		log.Printf("An error occurred while fetching position data: %v", err.Error())
	}

	defer positionRows.Close()

	for positionRows.Next() {
		var position models.Position
		if err := positionRows.Scan(&position.ID, &position.Position_name); err != nil {
			log.Printf("error: An error occurred while scanning position data")
		}

		positions = append(positions, position)
	}
	researcher.PrefixName = positions[0].Position_name

	// Fetch and add TempProgram data
	var programs []models.TempProgram_get
	programRows, err := h.db.Raw("SELECT id, program_name FROM program WHERE profile_id = ?", id).Rows()
	if err != nil {
		log.Printf("An error occurred while fetching program data: %v", err.Error())
	}

	defer programRows.Close()

	for programRows.Next() {
		var program models.TempProgram_get
		if err := programRows.Scan(&program.ProgramID, &program.ProgramName); err != nil {
			log.Printf("error: An error occurred while scanning program data")
		}

		programs = append(programs, program)
	}
	researcher.Program = programs

	if result.Error != nil || err != nil || errPositionID.Error != nil {
		return models.Researcher_get{}
	}

	return researcher
}
func (h *ResearcherHandler) CreateResearcher(c *gin.Context) {
	var researcher models.ResearcherRequest

	if err := c.ShouldBindJSON(&researcher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !validateDegreeOrder(researcher.Degree) {
		res := api.ResponseApi(http.StatusBadRequest, researcher.Degree, fmt.Errorf("invalid degree order"))
		c.JSON(http.StatusBadRequest, res)
		return
	}

	// Create createdBy and updatedBy variables
	createdBy := "Champlnwza007"
	updatedBy := "Champlnwza007"
	profileStatus := true

	// Find position name
	var positionID int
	var position models.Position

	// First, try to get the position from the database.
	result := h.db.Raw("SELECT * FROM position WHERE position_name = ?", researcher.PrefixName).Scan(&position)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Database query error position": result.Error.Error()})
		return
	}

	// If the ID is zero, it means no position was found, so create a new one.
	if position.ID == 0 {
		position.Created_by = createdBy
		position.Updated_by = updatedBy
		position.Position_name = researcher.PrefixName

		r := h.db.Create(&position)
		if err := r.Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Create position error": err.Error()})
			return
		}
		positionID = position.ID
	} else {
		positionID = position.ID
	}

	// Update the INSERT statement for the profile table
	insertResult := h.db.Exec("INSERT INTO profile (first_name, last_name, university, address_home, address_work, email, phone_number, position_id, created_by, updated_by, profile_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id",
		researcher.FirstName, researcher.LastName, researcher.University, researcher.AddressHome, researcher.AddressWork, researcher.Email, researcher.PhoneNumber, positionID, createdBy, updatedBy, profileStatus)

	if insertResult.Error != nil {
		res := api.ResponseApi(http.StatusBadRequest, researcher.Degree, fmt.Errorf("position name does not exist in the database"))
		c.JSON(http.StatusBadRequest, res)
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

	res := api.ResponseApiWithDescription(http.StatusCreated, h.ListResearcherbyID(profileID), "CREATED SUCCESS", nil)
	c.JSON(http.StatusCreated, res)
}

func (h *ResearcherHandler) VSdeleteResearcher(c *gin.Context) {
	profileID := c.Param("id")
	updatedBy := "Champlnwza007"
	profile_status := false

	tablesToUpdate := []string{"profile"}

	for _, tableName := range tablesToUpdate {
		if err := h.db.Exec(fmt.Sprintf("UPDATE %s SET updated_by = ?, profile_status = ? WHERE id = ?", tableName), updatedBy, profile_status, profileID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating the profile_status in the %s table: %v", tableName, err)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       200,
		"description":  "SUCCESS",
		"errorMessage": nil,
		"data": gin.H{
			"profile_status": profile_status,
		},
	})

}

// func (h *ResearcherHandler) UpdateResearcher(c *gin.Context) {
// 	var researcher models.ResearcherRequest
// 	profileID := c.Param("id")
// 	if err := c.ShouldBindJSON(&researcher); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if !validateDegreeOrder(researcher.Degree) {
// 		res := api.ResponseApi(http.StatusBadRequest, researcher.Degree, fmt.Errorf("invalid degree order"))
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	intProfileID, err := strconv.Atoi(profileID)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, fmt.Errorf("error to convert profile ID"))
// 		return
// 	}
// 	researcher.ProfileID = intProfileID
// 	//find position name
// 	var positionID int
// 	err = h.db.Raw("SELECT id FROM position WHERE position_name = ?", researcher.PositionName).Scan(&positionID).Error

// 	if err != nil {
// 		res := api.ResponseApi(http.StatusBadRequest, researcher.Degree, err)
// 		c.JSON(http.StatusBadRequest, res)
// 		return
// 	}

// 	// Update createdBy and updatedBy variables
// 	updatedBy := "Champlnwza007"

// 	// Update the researcher's profile data
// 	if err := h.db.Exec("UPDATE profile SET first_name = ?, last_name = ?, university = ?, address_home = ?, address_work = ?, email = ?, phone_number = ?, position_id = ?, updated_by = ? WHERE id = ?",
// 		researcher.FirstName, researcher.LastName, researcher.University, researcher.AddressHome, researcher.AddressWork, researcher.Email, researcher.PhoneNumber, positionID, updatedBy, profileID).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating researcher data in the profile table: %v", err)})
// 		return
// 	}

// 	// Fetch the IDs of all degree records associated with the profile_id
// 	var degreeIDs []int
// 	if err := h.db.Raw("SELECT id FROM degree WHERE profile_id = ?", profileID).Pluck("id", &degreeIDs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching degree ids: %v", err)})
// 		return
// 	}

// 	// Update degree data
// 	for i, degree := range researcher.Degree {
// 		// Use the degree ID fetched from the database
// 		degreeID := degreeIDs[i]
// 		if err := h.db.Exec("UPDATE degree SET degree_type = ?, degree_program = ?, degree_university = ?, updated_by = ? WHERE id = ?",
// 			degree.DegreeType, degree.DegreeProgram, degree.DegreeUniversity, updatedBy, degreeID).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating degree data: %v", err)})
// 			return
// 		}
// 	}

// 	// Fetch the IDs of all program records associated with the profile_id
// 	var programIDs []int
// 	if err := h.db.Raw("SELECT id FROM program WHERE profile_id = ?", profileID).Pluck("id", &programIDs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching program ids: %v", err)})
// 		return
// 	}

// 	// Update program data
// 	for i, program := range researcher.Program {
// 		// Use the program ID fetched from the database
// 		programID := programIDs[i]
// 		if err := h.db.Exec("UPDATE program SET program_name = ?, updated_by = ? WHERE id = ?",
// 			program.ProgramName, updatedBy, programID).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while updating program data: %v", err)})
// 			return
// 		}
// 	}

// 	// Fetch the IDs of all experience records associated with the profile_id
// 	var experienceIDs []int
// 	if err := h.db.Raw("SELECT id FROM experience WHERE profile_id = ?", profileID).Pluck("id", &experienceIDs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching experience ids: %v", err)})
// 		return
// 	}

// 	// Fetch the IDs of all explore records associated with the profile_id
// 	var exploreIDs []int
// 	if err := h.db.Raw("SELECT id FROM exploration WHERE profile_id = ?", profileID).Pluck("id", &exploreIDs).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("An error occurred while fetching explore ids: %v", err)})
// 		return
// 	}

// 	intProfileID, errProfileID := strconv.Atoi(profileID)

// 	if errProfileID != nil {
// 		fmt.Println("Error converting string to integer:", err)
// 		return
// 	}
// 	res := api.ResponseApiWithDescription(http.StatusCreated, h.ListResearcherbyID(intProfileID), "CREATED SUCCESS", nil)
// 	c.JSON(http.StatusCreated, res)

// }

func validateDegreeOrder(degrees []models.TempDegree_create) bool {
	hasBachelor := false
	hasMaster := false
	hasDoctor := false

	for _, degree := range degrees {
		switch degree.DegreeType {
		case "bachelor":
			hasBachelor = true
		case "master":
			hasMaster = true
		case "doctor":
			hasDoctor = true
		}
	}

	if hasDoctor && !hasMaster {
		return false
	}

	if hasMaster && !hasBachelor {
		return false
	}

	return true
}
