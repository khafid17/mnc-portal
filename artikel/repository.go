package artikel

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Artikel, error)
	FindByUserId(user int) ([]Artikel, error)
	FindById(ID int) (Artikel, error)
	Save(artikel Artikel) (Artikel, error)
	Update(artikel Artikel) (Artikel, error)
	CreateImage(Artikelmage ArtikelImage) (ArtikelImage, error)
	MarkAllImageAsNonPrimary(artikelID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindAll() ([]Artikel, error) {
	var artikel []Artikel

	err := r.db.Preload("ArtikelImages", "artikel_images.is_primary = 1").Find(&artikel).Error
	if  err != nil {
		return artikel, err
	}

	return artikel, nil
}

func (r *repository) FindByUserId(userID int) ([]Artikel, error) {
	var artikel []Artikel

	err := r.db.Where("user_id = ?", userID).Preload("ArtikelImages", "artikel_images.is_primary = 1").Find(&artikel).Error
	if err != nil {
		return artikel, err
	}

	return artikel, nil
}

func (r *repository) FindById(ID int) (Artikel, error) {
	var artikel Artikel

	err := r.db.Preload("User").Preload("ArtikelImages").Where("id = ?", ID).Find(&artikel).Error
	
	if err != nil {
		return artikel, err
	}
	
	return artikel, nil
}

func (r *repository) Save(artikel Artikel) (Artikel, error) {
	err := r.db.Create(&artikel).Error
	
	if err != nil {
		return artikel, err
	}

	return artikel, nil
} 

func (r *repository) Update(artikel Artikel) (Artikel, error) {
	err := r.db.Save(&artikel).Error
	
	if err != nil {
		return artikel, err
	}

	return artikel, nil
}

func (r *repository) CreateImage(artikelImage ArtikelImage) (ArtikelImage, error) {
	err := r.db.Create(&artikelImage).Error
	
	if err != nil {
		return artikelImage, err
	}

	return artikelImage, nil
}

func (r *repository) MarkAllImageAsNonPrimary(ArtikelID int) (bool, error){
	err := r.db.Model(&ArtikelImage{}).Where("artikel_id = ?", ArtikelID).Update("is_primary", false).Error
	if  err != nil {
		return false, err
	}
	return true, err
}