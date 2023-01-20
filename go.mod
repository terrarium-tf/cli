module github.com/terrarium-tf/cli

go 1.18

require (
	github.com/hashicorp/terraform-exec v0.16.0
	github.com/ojizero/gofindup v1.1.3
	github.com/spf13/cobra v1.3.0
)

require (
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/hashicorp/terraform-json v0.13.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/hashicorp/terraform-exec v0.16.0 => github.com/digitalkaoz/terraform-exec v0.16.1-0.20220222225145-2509737b3247
