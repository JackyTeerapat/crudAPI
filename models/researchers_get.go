package models

type Researcher_get struct {
	ProfileID   int              `json:"profile_id"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	University  string           `json:"university"`
	AddressHome string           `json:"address_home"`
	AddressWork string           `json:"address_work"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Degree      []TempDegree_get     `gorm:"-" json:"degree"`
	Position    []TempPosition_get   `gorm:"-" json:"position"`
	Program     []TempProgram_get    `gorm:"-" json:"program"`
	Experience  []TempExperience_get `gorm:"-" json:"experience"`
	Attach      []TempAttach_get     `gorm:"-" json:"attach"`
	Explore     []TempExplore_get    `gorm:"-" json:"explore"`
}
type TempDegree_get struct {
	DegreeID         int    `json:"id"`
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
	Activated        bool   `json:"activated"`
}
type TempPosition_get struct {
	PositionID   int    `json:"position_id"`
	PositionName string `json:"position_name"`
}
type TempProgram_get struct {
	ProgramID   int    `json:"program_id"`
	ProgramName string `json:"program_name"`
	Activated   bool   `json:"activated"`
}

type TempExperience_get struct {
	ExperienceID         int    `json:"experience_id"`
	ExperienceType       string `json:"experience_type"`
	ExperienceStart      string `json:"experience_start"`
	ExperienceEnd        string `json:"experience_end"`
	ExperienceUniversity string `json:"experience_university"`
	ExperienceRemark     string `json:"experience_remark"`
	Activated            bool   `json:"activated"`
}

type TempAttach_get struct {
	FileID      int    `json:"file_id"`
	FileName    string `json:"file_name"`
	FileAction  string `json:"file_action"`
	FileStorage string `json:"file_storage"`
	Activated   bool   `json:"activated"`
}

type TempExplore_get struct {
	ExploreID     int    `json:"explore_id"`
	ExploreName   string `json:"explore_name"`
	ExploreYear   string `json:"explore_year"`
	ExploreDetail string `json:"explore_detail"`
	Activated     bool   `json:"activated"`
}
