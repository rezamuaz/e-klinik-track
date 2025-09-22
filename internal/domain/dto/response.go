package dto

type PostedCreated struct {
	ID string `json:"id"`
}

type PostedDelete struct {
	PostID string `json:"post_id"`
}

type PostCreatedResponse struct {
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	AlternativeTitle []string `json:"alternative_title"`
	Artists          []string `json:"artists"`
	ArtistIDs        []string `json:"artist_ids"`
	ArtistKanji      []string `json:"artist_kanji"`
	Genres           []string `json:"genres"`
	GenreIds         []string `json:"genre_ids"`
	Thumbnails       string   `json:"thumbnails"`
	Category         []string `json:"category"`
	Tags             []string `json:"tags"`
	Status           string   `json:"status"`
	Views            int64    `json:"views"`
	PublishedAt      int64    `json:"published_at"`
	PublishedBy      string   `json:"published_by"`
	CreatedAt        int64    `json:"created_at"`
	CreatedBy        string   `json:"created_by"`
	UpdatedAt        int64    `json:"updated_at"`
	UpdatedBy        string   `json:"updated_by"`
}

type Taxonomy struct {
	Name string `json:"name"`
}

// func ToPostIndexCreatedResp(data pg.GetPostForIndexByIDRow) (PostCreatedResponse, error) {

// 	getStringArray := func(val interface{}) ([]string, error) {
// 		arr, ok := val.([]string)
// 		if !ok {
// 			return nil, fmt.Errorf("cannot convert %v to []string", val)
// 		}
// 		return arr, nil
// 	}

// 	artists, err := getStringArray(data.Artists)
// 	if err != nil {
// 		return PostCreatedResponse{}, fmt.Errorf("artists: %w", err)
// 	}

// 	artistIDs, err := getStringArray(data.ArtistIds)
// 	if err != nil {
// 		return PostCreatedResponse{}, fmt.Errorf("artist_ids: %w", err)
// 	}

// 	artistKanji, err := getStringArray(data.ArtistKanji)
// 	if err != nil {
// 		return PostCreatedResponse{}, fmt.Errorf("artist_kanji: %w", err)
// 	}

// 	genres, err := getStringArray(data.Genres)
// 	if err != nil {
// 		return PostCreatedResponse{}, fmt.Errorf("genres: %w", err)
// 	}

// 	genreIds, err := getStringArray(data.GenreIds)
// 	if err != nil {
// 		return PostCreatedResponse{}, fmt.Errorf("genre_ids: %w", err)
// 	}
// 	// Set all other fields
// 	return PostCreatedResponse{
// 		ID:               data.ID,
// 		Title:            data.Title,
// 		AlternativeTitle: data.AlternativeTitle,
// 		Artists:          artists,
// 		ArtistIDs:        artistIDs,
// 		ArtistKanji:      artistKanji,
// 		Genres:           genres,
// 		GenreIds:         genreIds,
// 		Thumbnails:       data.Thumbnails.String,
// 		Category:         data.Category,
// 		Tags:             data.Tags,
// 		Status:           data.Status.String,
// 		Views:            data.Views,
// 		PublishedAt:      data.PublishedAt.Int64,
// 		PublishedBy:      data.PublishedBy.String,
// 		CreatedAt:        data.CreatedAt,
// 		CreatedBy:        data.CreatedBy.String,
// 		UpdatedAt:        data.UpdatedAt,
// 		UpdatedBy:        data.UpdatedBy.String,
// 	}, nil
// }
