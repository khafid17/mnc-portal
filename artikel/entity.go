package artikel

import (
	"mnc-portal/user"
	"time"
)

type Artikel struct {
	ID            int
	Judul         string
	Diskripsi     string
	Kategori      string
	UserId        int
	Slug          string
	Tags          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ArtikelImages []ArtikelImage
	User          user.User
}

type ArtikelImage struct {
	ID        int
	ArtikelID int
	FileName  string
	IsPrimary int
	CreatedAt time.Time
	UpdatedAt time.Time
}
