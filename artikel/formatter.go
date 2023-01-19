package artikel

import (
	"strings"
)

type ArtikelFormatter struct {
	ID        int    `json:"id"`
	UserId    int    `json:"user_id"`
	Judul     string `json:"judul"`
	Diskripsi string `json:"diskripsi"`
	Kategori  string `json:"kategori"`
	ImageURL  string `json:"image_url"`
	Slug      string `json:"slug"`
}

func FormatterArtikel(artikel Artikel) ArtikelFormatter {
	artikelFormatter := ArtikelFormatter{}
	artikelFormatter.ID = artikel.ID
	artikelFormatter.UserId = artikel.UserId
	artikelFormatter.Judul = artikel.Judul
	artikelFormatter.Diskripsi = artikel.Diskripsi
	artikelFormatter.Kategori = artikel.Kategori
	artikelFormatter.Slug = artikel.Slug
	artikelFormatter.ImageURL = ""

	if len(artikel.ArtikelImages) > 0 {
		artikelFormatter.ImageURL = artikel.ArtikelImages[0].FileName
	}

	return artikelFormatter
}

func FormatterArtikels(artikels []Artikel) []ArtikelFormatter {
	artikelsFormatter := []ArtikelFormatter{}

	for _, artikel := range artikels {
		artikelFormatter := FormatterArtikel(artikel)
		artikelsFormatter = append(artikelsFormatter, artikelFormatter)
	}

	return artikelsFormatter
}

type ArtikelDetailFormatter struct {
	ID        int                     `json:"id"`
	Judul     string                  `json:"name"`
	Diskripsi string                  `json:"diskripsi"`
	Kategori  string                  `json:"kategori"`
	ImagesURL string                  `json:"image_url"`
	UserID    int                     `json:"iser_id"`
	Slug      string                  `json:"slug"`
	Tags      []string                `json:"tags"`
	User      ArtikelUserFormatter    `json:"user"`
	Images    []ArtikelImageFormatter `json:"images"`
}

type ArtikelUserFormatter struct {
	Name      string `json:"name"`
	ImagesUrl string `json:"image_url"`
}

type ArtikelImageFormatter struct {
	ImagesUrl string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

func FormatArtikelDetail(artikel Artikel) ArtikelDetailFormatter {
	artikelDetailFormatter := ArtikelDetailFormatter{}
	artikelDetailFormatter.ID = artikel.ID
	artikelDetailFormatter.Judul = artikel.Judul
	artikelDetailFormatter.Diskripsi = artikel.Diskripsi
	artikelDetailFormatter.Kategori = artikel.Kategori
	artikelDetailFormatter.UserID = artikel.UserId
	artikelDetailFormatter.Slug = artikel.Slug
	artikelDetailFormatter.ImagesURL = ""

	if len(artikel.ArtikelImages) > 0 {
		artikelDetailFormatter.ImagesURL = artikel.ArtikelImages[0].FileName
	}

	var tags []string

	for _, tag := range strings.Split(artikel.Tags, ",") {
		tags = append(tags, tag)
	}

	artikelDetailFormatter.Tags = tags

	user := artikel.User

	artikelUserFormatter := ArtikelUserFormatter{}
	artikelUserFormatter.Name = user.Name
	artikelUserFormatter.ImagesUrl = user.AvatarFileName

	artikelDetailFormatter.User = artikelUserFormatter

	images := []ArtikelImageFormatter{}

	for _, image := range artikel.ArtikelImages {
		artikelImageFormatter := ArtikelImageFormatter{}
		artikelImageFormatter.ImagesUrl = image.FileName

		isPrimary := false

		if image.IsPrimary == 1 {
			isPrimary = true
		}

		artikelImageFormatter.IsPrimary = isPrimary

		images = append(images, artikelImageFormatter)

	}

	artikelDetailFormatter.Images = images

	return artikelDetailFormatter
}
