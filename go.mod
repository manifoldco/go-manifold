module github.com/manifoldco/go-manifold

require (
	github.com/alecthomas/gometalinter v2.0.12+incompatible
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/client9/misspell v0.3.4
	github.com/dchest/blake2b v1.0.0
	github.com/go-openapi/analysis v0.18.0 // indirect
	github.com/go-openapi/errors v0.18.0 // indirect
	github.com/go-openapi/jsonpointer v0.18.0 // indirect
	github.com/go-openapi/jsonreference v0.18.0 // indirect
	github.com/go-openapi/loads v0.18.0 // indirect
	github.com/go-openapi/runtime v0.18.0
	github.com/go-openapi/spec v0.18.0 // indirect
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0 // indirect
	github.com/go-openapi/validate v0.18.0 // indirect
	github.com/golang/lint v0.0.0-20181217174547-8f45f776aaf1
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/kr/pty v1.1.3 // indirect
	github.com/manifoldco/go-base32 v1.0.2
	github.com/manifoldco/go-base64 v1.0.1
	github.com/pkg/errors v0.8.0
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	golang.org/x/crypto v0.0.0-20190102171810-8d7daa0c54b3
	golang.org/x/lint v0.0.0-20181217174547-8f45f776aaf1 // indirect
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3 // indirect
	golang.org/x/tools v0.0.0-20190102200130-52ae6dee2324 // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce // indirect
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
