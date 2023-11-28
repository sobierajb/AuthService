package certRepo

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type Certificate interface {
	GetAlghorythm() AlghorythmType
	GetThumbprint() string
	GetSigningMethod() (jwt.SigningMethod, error)
	GetRawData() []byte
	Generate(AlghorythmType) error
	GetPublicKey() (interface{}, error)
	GetPrivateKey() (interface{}, error)
	GetPublicStruct() (*PublicCertResponse, error)
}

type CertRepo interface {
	GetSigningKeyAndMethod() (string, interface{})
	GetCryptKeyAndMethod() (string, interface{})
	CreateCertificate(Certificate) (Certificate, error)
	GetCertificateById(id string) (Certificate, error)
	GetCertificateByClientId(clientId string) (Certificate, error)
	DeleteCertificate(id string) error
}

type certRepo struct {
	db *gorm.DB
	cf *certFactory
}

func NewCertRepo(db *gorm.DB) *certRepo {
	return &certRepo{
		db: db,
	}
}

func (cr *certRepo) CreateCertificate(kt KeyType) (Certificate, error) {

	// Get constructor function
	constFunc := cr.cf.certCreators[string(kt)]

	// Call constructor to create proper certificate object
	cert := constFunc()

	// generate new certifcate key
	err := cert.Generate(kt)
	if err != nil {
		return nil, err
	}

	// get certificate pointer where orm will write results from DB
	certData := cert.GetRawData()

	// call orm create function
	result := cr.db.Create(&certData)
	if result.Error != nil {
		return nil, result.Error
	}
	return cert, nil
}

// Read certificateData from DB and create Certificate object
func (cr *certRepo) GetCertificateById(id string) (Certificate, error) {
	// Create empty cert data with id
	certData := &CertificateData{
		Id: id,
	}
	// Fill certData with data from DB
	result := cr.db.Find(&certData)
	if result.RowsAffected == 0 {
		return nil, ErrCertNotFound
	}
	// Get constructor to create certificate object based on Type
	constFunc := cr.cf.certCreators[string(certData.Type)]
	// Create certificate object
	cert := constFunc()
	// Set raw Data and return
	err := cert.SetRawData(certData)
	if err != nil {
		return nil, err
	}
	return cert, err

}

func (cr *certRepo) DeleteCertificate(id string) error {
	certData := &CertificateData{Id: id}
	result := cr.db.Delete(&certData)
	if result.RowsAffected == 0 {
		return ErrCertNotFound
	}
	return result.Error
}
