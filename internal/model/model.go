package model

type Family struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Group struct {
	ID       int    `json:"id"`
	FamilyID int    `json:"family_id"`
	Name     string `json:"name"`
}

type Image struct {
	ID         int      `json:"id"`
	GroupID    int      `json:"group_id"`
	Name       string   `json:"name"`
	FilePath   string   `json:"file_path"`
	UsageCount int      `json:"usage_count"`
	MetaTags   []string `json:"meta_tags"`
}
