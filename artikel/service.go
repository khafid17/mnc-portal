package artikel

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetArtikels(userID int) ([]Artikel, error)
	GetArtikelByID(input GetArtikelDetailInput) (Artikel, error)
	CreateArtikel(input CreateArtikelInput) (Artikel, error)
	UpdateArtikel(inputID GetArtikelDetailInput, inputData CreateArtikelInput) (Artikel, error)
	SaveArtikelImage(input CreateArtikelImageInput, fileLocation string) (ArtikelImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetArtikels(userID int) ([]Artikel, error) {
	if userID != 0 {
		Artikels, err := s.repository.FindByUserId(userID)
		if err != nil {
			return Artikels, err
		}

		return Artikels, nil
	}
	Artikels, err := s.repository.FindAll()
	if err != nil {
		return Artikels, err
	}

	return Artikels, nil
}

func (s *service) GetArtikelByID(input GetArtikelDetailInput) (Artikel, error) {
	artikel, err := s.repository.FindById(input.ID)

	if err != nil {
		return artikel, err
	}

	return artikel, nil
}

func (s *service) CreateArtikel(input CreateArtikelInput) (Artikel, error) {
	artikel := Artikel{}
	artikel.Judul = input.Judul
	artikel.Diskripsi = input.Diskripsi
	artikel.Kategori = input.Kategori
	artikel.Tags = input.Tags
	artikel.UserId = input.User.ID

	slugCandidate := fmt.Sprintf("%s %d", input.Judul, input.User.ID)
	artikel.Slug = slug.Make(slugCandidate)

	newArtikel, err := s.repository.Save(artikel)
	if err != nil {
		return newArtikel, err
	}

	return newArtikel, nil
}

func (s *service) UpdateArtikel(inputID GetArtikelDetailInput, inputData CreateArtikelInput) (Artikel, error) {
	artikel, err := s.repository.FindById(inputID.ID)
	if err != nil {
		return artikel, err
	}

	if artikel.UserId != inputData.User.ID {
		return artikel, errors.New("Not an awner of the artikel")
	}

	artikel.Judul = inputData.Judul
	artikel.Diskripsi = inputData.Diskripsi
	artikel.Kategori = inputData.Kategori
	artikel.Tags = inputData.Tags

	updatedArtikel, err := s.repository.Update(artikel)
	if err != nil {
		return updatedArtikel, err
	}

	return updatedArtikel, nil
}

func (s *service) SaveArtikelImage(input CreateArtikelImageInput, fileLocation string) (ArtikelImage, error) {
	artikel, err := s.repository.FindById(input.ArtikelID)
	if err != nil {
		return ArtikelImage{}, err
	}

	if artikel.UserId != input.User.ID {
		return ArtikelImage{}, errors.New("Not an awner of the artikel")
	}

	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImageAsNonPrimary(input.ArtikelID)
		if err != nil {
			return ArtikelImage{}, err
		}
	}

	artikelImage := ArtikelImage{}
	artikelImage.ArtikelID = input.ArtikelID
	artikelImage.IsPrimary = isPrimary
	artikelImage.FileName = fileLocation

	newArtikelImage, err := s.repository.CreateImage(artikelImage)
	if err != nil {
		return newArtikelImage, err
	}

	return newArtikelImage, nil
}
