package models

type AssessmentRequest struct {
	ProfileID              int             `json:"profile_id"`
	AssessmentStart        string          `json:"assessment_start"`
	AssessmentEnd          string          `json:"assessment_end"`
	Assessment_file_name   string          `json:"assessment_file_name"`
	Assessment_file_action string          `json:"assessment_file_action"`
	Project                ProjectDetails  `json:"assessment_project"`
	Progress               ProgressDetails `json:"assessment_progress"`
	Report                 ReportDetails   `json:"assessment_report"`
	Article                ArticleDetails  `json:"assessment_article"`
}

type ProjectDetails struct {
	ProjectYear            string `json:"project_year"`
	ProjectTitle           string `json:"project_title"`
	ProjectPoint           int    `json:"project_point"`
	ProjectEstimate        bool   `json:"project_estimate"`
	ProjectRecommend       bool   `json:"project_recommend"`
	ProjectPeriod          bool   `json:"project_period"`
	ProjectFileName        string `json:"file_name"`
	ProjectTitleFileAction string `json:"file_action"`
}

type ProgressDetails struct {
	ProgressYear            string `json:"progress_year"`
	ProgressTitle           string `json:"progress_title"`
	ProgressEstimate        bool   `json:"progress_estimate"`
	ProgressRecommend       bool   `json:"progress_recommend"`
	ProgressPeriod          bool   `json:"progress_period"`
	ProgressFileName        string `json:"file_name"`
	ProgressTitleFileAction string `json:"file_action"`
}

type ReportDetails struct {
	ReportYear            string `json:"report_year"`
	ReportTitle           string `json:"report_title"`
	ReportEstimate        bool   `json:"report_estimate"`
	ReportRecommend       bool   `json:"report_recommend"`
	ReportPeriod          bool   `json:"report_period"`
	ReportFileName        string `json:"file_name"`
	ReportTitleFileAction string `json:"file_action"`
}

type ArticleDetails struct {
	ArticleYear            string `json:"article_year"`
	ArticleTitle           string `json:"article_title"`
	ArticleEstimate        bool   `json:"article_estimate"`
	ArticleRecommend       bool   `json:"article_recommend"`
	ArticlePeriod          bool   `json:"article_period"`
	ArticleFileName        string `json:"file_name"`
	ArticleTitleFileAction string `json:"file_action"`
}
