package registry

type Device struct {
	ID                              []byte `gorm:"primary_key"`
	Namespace                       string `gorm:"type: text CHECK(length(namespace)>2);NOT NULL;unique_index:device_name_namespace_uq"`
	Name                            string `gorm:"NOT NULL;unique_index:device_name_namespace_uq"`
	Enabled                         bool
	Certificate                     string
	CertificateType                 string
	CertificateFingerprint          []byte `gorm:"index"`
	CertificateFingerprintAlgorithm string
}
