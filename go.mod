module github.com/manifoldco/go-manifold

require (
	github.com/alecthomas/gometalinter v2.0.11+incompatible
	github.com/asaskevich/govalidator v0.0.0-20161001163130-7b3beb6df3c4
	github.com/client9/misspell v0.3.4
	github.com/dchest/blake2b v1.0.0
	github.com/go-openapi/errors v0.0.0-20170104180542-fc3f73a22449 // indirect
	github.com/go-openapi/runtime v0.0.0-20170303002511-e66a4c440602
	github.com/go-openapi/strfmt v0.0.0-20170319025125-93a31ef21ac2
	github.com/go-openapi/swag v0.0.0-20170129222639-d5f8ebc3b1c5 // indirect
	github.com/golang/lint v0.0.0-20181026193005-c67002cb31c3
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/mailru/easyjson v0.0.0-20170328210357-2af9a745a611 // indirect
	github.com/manifoldco/go-base32 v1.0.2
	github.com/manifoldco/go-base64 v1.0.1
	github.com/mitchellh/mapstructure v0.0.0-20170307201123-53818660ed49 // indirect
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	golang.org/x/crypto v0.0.0-20170930174604-9419663f5a44
	golang.org/x/net v0.0.0-20170406210907-d1e1b351919c // indirect
	golang.org/x/tools v0.0.0-20181120200622-9c8bd463e3ac // indirect
	gopkg.in/mgo.v2 v2.0.0-20160818020120-3f83fa500528 // indirect
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
