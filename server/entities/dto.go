package entities

type X509CertificateRequest struct {
	Certificate string `json:"certificate" validate:"required"`
}

type X509CertificateResponse struct {
	X509Certificate string `json:"x509_certificate"`
}
