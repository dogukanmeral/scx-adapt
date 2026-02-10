module helper

go 1.25.6

require internal/checks v1.0.0

require internal/errs v1.0.0

replace internal/checks => ../checks

replace internal/errs => ../errs

require (
	github.com/go-playground/validator/v10 v10.30.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
