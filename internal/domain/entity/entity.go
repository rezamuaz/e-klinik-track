package entity

type Thumbnail struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

type MenuNode struct {
	ID          int32       `json:"id"`
	Label       string      `json:"label"`
	ResourceKey string      `json:"resource_key"`
	Action      string      `json:"action"`
	View        *string     `json:"view"`
	Data        *string     `json:"data"`
	Level       *int16      `json:"level"`
	Path        *string     `json:"path"`
	Children    []*MenuNode `json:"children,omitempty"`
}
