module github.com/terrarium-tf/cli

go 1.18

require (
	github.com/hashicorp/terraform-exec v0.17.3
	github.com/ojizero/gofindup v1.1.3
	github.com/spf13/cobra v1.6.1
)

require (
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/terraform-json v0.14.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	golang.org/x/text v0.6.0 // indirect
)

replace github.com/hashicorp/terraform-exec v0.17.3 => github.com/digitalkaoz/terraform-exec v0.17.4-0.20230120122230-3d5d15e21122
