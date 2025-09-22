package enum

type PostStatus int16

const (
	Publish PostStatus = iota + 1
	Pending
	Draf
	AutoDraf
	Future
	Private
	Trash
)

func (s PostStatus) String() string {
	return [...]string{"Generation One", "Generation Two"}[s]
}

func (s PostStatus) ValueToString(val int16) string {
	return [...]string{"", "Generation One", "Generation Two"}[val]
}

type StatusListResponse struct {
	Name  string     `json:"name"`
	Value PostStatus `json:"value"`
	Url   string     `json:"url"`
}

// var SourceName = map[PostStatus]StatusInfo {
// 	MyAnimeList: {Name: "My Anime List", Source: MyAnimeList, Url: "https://api.jikan.moe/v4"},
// 	MangaUpdate: {Name: "pvsave Update", Source: MangaUpdate, Url: "https://api.mangaupdates.com"},
// }

// type StatusInfo struct {
// 	Name   string
// 	Source PostStatus
// 	Url    string
// }
