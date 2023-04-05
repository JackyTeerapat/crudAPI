package models

type Researcher struct {
	ProfileID   int              `json:"profile_id"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	University  string           `json:"university"`
	AddressHome string           `json:"address_home"`
	AddressWork string           `json:"address_work"`
	Email       string           `json:"email"`
	PhoneNumber string           `json:"phone_number"`
	Degree      []TempDegree     `gorm:"-" json:"degree"`
	Position    []TempPosition   `gorm:"-" json:"position"`
	Program     []TempProgram    `gorm:"-" json:"program"`
	Experience  []TempExperience `gorm:"-" json:"experience"`
	Attach      []TempAttach     `gorm:"-" json:"attach"`
	Explore     []TempExplore    `gorm:"-" json:"explore"`
}
type TempDegree struct {
	DegreeID         int    `json:"id"`
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
	Activated        bool   `json:"activated"`
}
type TempPosition struct {
	PositionID   int    `json:"position_id"`
	PositionName string `json:"position_name"`
}
type TempProgram struct {
	ProgramID   int    `json:"program_id"`
	ProgramName string `json:"program_name"`
	Activated   bool   `json:"activated"`
}

type TempExperience struct {
	ExperienceID         int    `json:"experience_id"`
	ExperienceType       string `json:"experience_type"`
	ExperienceStart      string `json:"experience_start"`
	ExperienceEnd        string `json:"experience_end"`
	ExperienceUniversity string `json:"experience_university"`
	ExperienceRemark     string `json:"experience_remark"`
	Activated            bool   `json:"activated"`
}

type TempAttach struct {
	FileID      int    `json:"file_id"`
	FileName    string `json:"file_name"`
	FileAction  string `json:"file_action"`
	FileStorage string `json:"file_storage"`
	Activated   bool   `json:"activated"`
}

type TempExplore struct {
	ExploreID     int    `json:"explore_id"`
	ExploreName   string `json:"explore_name"`
	ExploreYear   string `json:"explore_year"`
	ExploreDetail string `json:"explore_detail"`
	Activated     bool   `json:"activated"`
}
