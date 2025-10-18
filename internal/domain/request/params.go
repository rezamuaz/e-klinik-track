package request

import "github.com/gofrs/uuid/v5"

type SearchParams struct {
	Query  string
	Type   []string
	Genre  []string
	Status []string
	Sort   string
	Limit  int64
	Page   int64
}

type CreatePost struct {
	Title            string   `json:"title"`
	AlternativeTitle []string `json:"alternative_title"`
	ArtistIDs        []string `json:"artist_ids"`
	GenreIDs         []string `json:"genre_ids"`
	Cover            []string `json:"cover"`
	Thumbnails       string   `json:"thumbnails"`
	Category         []string `json:"category"`
	Tags             []string `json:"tags"`
	Status           string   `json:"status"`
	ExtraInfo        string   `json:"extra_info"`
	Description      string   `json:"description"`
	CreatedBy        string   `json:"created_by"`
}

type UpdatePost struct {
	ID string `json:"id"`
}

type GetKeyAndValue struct {
	AppID    string `form:"app_id" binding:"required"`
	TenantID string `form:"tenant_id" binding:"required"`
	Name     string `form:"name" binding:"required"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Register struct {
	Username  string   `json:"username"`
	Nama      string   `json:"nama"`
	Password  string   `json:"password"`
	Role      []Option `json:"role"`
	CreatedBy *string  `json:"created_by"`
}

type SearchMataKuliah struct {
	Page       int32   `form:"page" json:"page"`
	MataKuliah *string `form:"mata_kuliah" json:"mata_kuliah"`
	IsActive   *bool   `form:"is_active" json:"is_active"`
	OrderBy    *string `form:"order_by" json:"order_by"`
	Sort       *string `form:"sort" json:"sort"`
	Offset     int32   `form:"offset" json:"offset"`
	Limit      int32   `form:"limit" json:"limit"`
}

type SearchFasilitasKesehatan struct {
	Page     int32   `form:"page" json:"page"`
	Nama     *string `form:"nama" json:"nama"`
	Propinsi *string `form:"propinsi" json:"propinsi"`
	Kab      *string `form:"kab" json:"kab"`
	KabID    *string `form:"kab_id" json:"kab_id"`
	Tipe     *string `form:"tipe" json:"tipe"`
	IsActive *bool   `form:"is_active" json:"is_active"`
	OrderBy  *string `form:"order_by" json:"order_by"`
	Sort     *string `form:"sort" json:"sort"`
	Offset   int32   `form:"offset" json:"offset"`
	Limit    int32   `form:"limit" json:"limit"`
}
type SearchKontrak struct {
	Page              int32   `form:"page" json:"page"`
	Nama              *string `form:"nama" json:"nama"`
	FasilitasNama     *string `form:"fasilitas_nama" json:"fasilitas_nama"`
	FasilitasKab      *string `form:"fasilitas_kab" json:"fasilitas_kab"`
	FasilitasPropinsi *string `form:"fasilitas_propinsi" json:"fasilitas_propinsi"`
	PeriodeMulai      *string `form:"periode_mulai" json:"periode_mulai"`
	PeriodeSelesai    *string `form:"periode_selesai" json:"periode_selesai"`
	IsActive          *bool   `form:"is_active" json:"is_active"`
	OrderBy           *string `form:"order_by" json:"order_by"`
	Sort              *string `form:"sort" json:"sort"`
	Offset            int32   `form:"offset" json:"offset"`
	Limit             int32   `form:"limit" json:"limit"`
}

type SearchAktifKontrak struct {
	FasilitasNama *string `form:"fasilitas_nama" json:"fasilitas_nama"`
}

type CreateKontrak struct {
	FasilitasID    uuid.UUID `json:"fasilitas_id"`
	Nama           string    `json:"nama"`
	PeriodeMulai   *string   `json:"periode_mulai"`
	PeriodeSelesai *string   `json:"periode_selesai"`
	Durasi         *string   `json:"durasi"`
	Deskripsi      *string   `json:"deskripsi"`
	CreatedBy      *string   `json:"created_by"`
}

type SearchRuangan struct {
	Page        int32   `form:"page" json:"page"`
	NamaRuangan *string `form:"nama_ruangan" json:"nama_ruangan"`
	Fasilitas   *string `form:"fasilitas" json:"fasilitas"`
	Kontrak     *string `form:"kontrak" json:"kontrak"`
	IsActive    *bool   `form:"is_active" json:"is_active"`
	OrderBy     *string `form:"order_by" json:"order_by"`
	Sort        *string `form:"sort" json:"sort"`
	Offset      int32   `form:"offset" json:"offset"`
	Limit       int32   `form:"limit" json:"limit"`
}

type SearchRuanganByKontrak struct {
	FasilitasID string `form:"fasilitas_id" json:"fasilitas_id"`
	KontrakID   string `form:"kontrak_id" json:"kontrak_id"`
}

type SearchKehadiran struct {
	Page         int32   `form:"page" json:"page"`
	JadwalDinas  *string `json:"jadwal_dinas"`
	FasilitasID  *string `json:"fasilitas_id"`
	UserID       *string `json:"user_id"`
	KontrakID    *string `json:"kontrak_id"`
	PembimbingID *string `json:"pembimbing_id"`
	IsActive     *bool   `form:"is_active" json:"is_active"`
	OrderBy      *string `form:"order_by" json:"order_by"`
	Sort         *string `form:"sort" json:"sort"`
	Offset       int32   `form:"offset" json:"offset"`
	Limit        int32   `form:"limit" json:"limit"`
}
type SearchKehadiranSkp struct {
	Page            int32   `form:"page" json:"page"`
	Status          *string `json:"status"`
	KehadiranID     *string `json:"kehadiran_id"`
	SkpIntervensiID *string `json:"skp_intervensi_id"`
	IsActive        *bool   `form:"is_active" json:"is_active"`
	OrderBy         *string `form:"order_by" json:"order_by"`
	Sort            *string `form:"sort" json:"sort"`
	Offset          int32   `form:"offset" json:"offset"`
	Limit           int32   `form:"limit" json:"limit"`
}

type SearchPropinsi struct {
	Propinsi *string `form:"propinsi" json:"propinsi"`
}
type SearchKabupaten struct {
	Kab        *string `form:"kab" json:"kab"`
	PropinsiID *string `form:"propinsi_id" json:"propinsi_id"`
}

type CreateRoleForUser struct {
	UserID    uuid.UUID `json:"user_id"`
	RoleID    uuid.UUID `json:"role_id"`
	CreatedBy *string   `json:"created_by"`
}

type SearchMenu struct {
	Label *string `form:"label" json:"label"`
}

type SearchUser struct {
	Page     int32   `form:"page" json:"page"`
	Nama     *string `form:"nama" json:"nama"`
	Username *string `form:"username" json:"username"`
	IsActive *bool   `form:"is_active" json:"is_active"`
	OrderBy  *string `form:"order_by" json:"order_by"`
	Sort     *string `form:"sort" json:"sort"`
	Offset   int32   `form:"offset" json:"offset"`
	Limit    int32   `form:"limit" json:"limit"`
}
type UpdateUser struct {
	Nama        *string   `json:"nama"`
	IsActive    *bool     `json:"is_active"`
	Password    *string   `json:"password"`
	Refresh     *string   `json:"refresh"`
	UpdatedNote *string   `json:"updated_note"`
	UpdatedBy   *string   `json:"updated_by"`
	ID          uuid.UUID `json:"id"`
	Role        []Option  `json:"role"`
}

type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type UpdateRolePolicy struct {
	Policy    []int32 `json:"policy"`
	RoleID    int32   `json:"role_id"`
	CreatedBy *string `json:"created_by"`
}
type Refresh struct {
	RefreshToken string `json:"refreshToken"`
}

type SearchRekapKehadiranMahasiswa struct {
	UserID   string `form:"user_id" json:"user_id"`
	TglAwal  string `form:"tgl_awal" json:"tgl_awal"`
	TglAkhir string `form:"tgl_akhir" json:"tgl_akhir"`
}
