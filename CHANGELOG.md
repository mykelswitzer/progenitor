# Changelog


## [2.0.1]

### Updated
- Merged pkg/scaffold and pkg/template as the only files that used the template were in the scaffold folder

## [2.0.0]

### Updated
- Moved scaffold structural logic to pkg level to be used by remote scaffolds
- Created constants for all config setting keys
- Removed the go-grpc scaffold file and templates, and moved them to progenitor-tmpl-go-grpc repo

## [1.0.4]

### Updated
- Terraform plan in go-grpc project template to use Terraform `0.15.4` and latest version of `tf-modules`
