package certRepo

type CertCreatorFunc func() Certificate

type certFactory struct {
	certCreators map[string]CertCreatorFunc
}

func NewCertFactory() *certFactory {
	cf := &certFactory{
		certCreators: make(map[string]CertCreatorFunc),
	}
	cf.AddCreatorFunc("RS256", NewRsaCert)
	cf.AddCreatorFunc("RS384", NewRsaCert)
	cf.AddCreatorFunc("RS512", NewRsaCert)
	return cf
}

func (cf *certFactory) AddCreatorFunc(alghorythmName string, f CertCreatorFunc) {
	cf.certCreators[alghorythmName] = f
}
