package a

// This empty file is needed here to avoid following error:
//   Cannot find module providing package github.com/citihub/probr-sdk: invalid github.com/ import path "github.com/citihub"

// To replicate issue, remove this file and run TestReadStaticFile test func:
// > go clean --testcache
// > go test ./utils -run TestReadStaticFile
// Ref:
//   https://github.com/markbates/pkger/issues/49
//   https://github.com/markbates/pkger/issues/44
