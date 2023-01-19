package artikel

import "mnc-portal/user"

type GetArtikelDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateArtikelInput struct {
	Judul     string `json:"judul" binding:"required"`
	Diskripsi string `json:"diskripsi" binding:"required"`
	Kategori  string `json:"kategori" binding:"required"`
	Tags      string `json:"tags" binding:"required"`
	User      user.User
}

type CreateArtikelImageInput struct {
	ArtikelID int  `form:"artikel_id" binding:"required"`
	IsPrimary bool `form:"is_primary"`
	User      user.User
}
