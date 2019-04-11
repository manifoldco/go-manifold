module github.com/manifoldco/go-manifold

require (
	github.com/alecthomas/gometalinter v3.0.0+incompatible
	github.com/asaskevich/govalidator v0.0.0-20180720115003-f9ffefc3facf
	github.com/client9/misspell v0.3.4
	github.com/dchest/blake2b v1.0.0
	github.com/go-openapi/errors v0.0.0-20170104180542-fc3f73a22449 // indirect
	github.com/go-openapi/runtime v0.0.0-20170303002511-e66a4c440602
	github.com/go-openapi/strfmt v0.0.0-20170319025125-93a31ef21ac2
	github.com/go-openapi/swag v0.0.0-20170129222639-d5f8ebc3b1c5 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/mailru/easyjson v0.0.0-20170328210357-2af9a745a611 // indirect
	github.com/manifoldco/go-base32 v1.0.3
	github.com/manifoldco/go-base64 v1.0.2
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	golang.org/x/crypto v0.0.0-20190123085648-057139ce5d2b
	golang.org/x/lint v0.0.0-20181217174547-8f45f776aaf1
	golang.org/x/net v0.0.0-20190125091013-d26f9f9a57f3 // indirect
	golang.org/x/tools v0.0.0-20190124215303-cc6a436ffe6b // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
