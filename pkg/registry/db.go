package registry

type Device struct {
	ID                              []byte `gorm:"primary_key"`
	Name                            string `gorm:"NOT NULL;unique_index:device_name_namespace_uq"`
	Enabled                         bool
	Certificate                     string
	CertificateType                 string
	CertificateFingerprint          []byte `gorm:"index"`
	CertificateFingerprintAlgorithm string
}
