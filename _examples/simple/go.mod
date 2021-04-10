module simple

go 1.16

require (
	github.com/aws/aws-lambda-go v1.23.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-chi/chi v1.5.4
	github.com/maxstanley/tango v0.0.0-20210403165208-2dd50600f7a6
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.26.0
)

replace github.com/maxstanley/tango => ../../../tango
