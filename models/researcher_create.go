package models

type ResearcherRequest struct {
	PrefixName string `json:"prefix_name"`
	// PositionID   int                  `json:"position_id"`
	PositionName string               `json:"position_name"`
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	AddressHome  string               `json:"address_home"`
	Degree       []TempDegree_create  `json:"degree"`
	Program      []TempProgram_create `json:"program"`
	University   string               `json:"university"`
	AddressWork  string               `json:"address_work"`
	Email        string               `json:"email"`
	PhoneNumber  string               `json:"phone_number"`
}

type TempDegree_create struct {
	DegreeType       string `json:"degree_type"`
	DegreeProgram    string `json:"degree_program"`
	DegreeUniversity string `json:"degree_university"`
}

type TempProgram_create struct {
	ProgramName string `json:"program_name"`
}

type TempExperience_create struct {
	ExperienceType       string  `json:"experience_type"`
	ExperienceStart      string  `json:"experience_start"`
	ExperienceEnd        *string `json:"experience_end"`
	ExperienceUniversity string  `json:"experience_university"`
	ExperienceRemark     string  `json:"experience_remark"`
}

type TempExplore_create struct {
	ExploreName   string `json:"explore_name"`
	ExploreYear   string `json:"explore_year"`
	ExploreDetail string `json:"explore_detail"`
}

type TempAttach_create struct {
	FileName    string `json:"file_name"`
	FileAction  string `json:"file_action"`
	FileStorage string `json:"file_storage"`
}
