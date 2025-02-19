package HTTPServer

import "github.com/aerosystems/common-service/presenters/httpserver"

type Config struct {
	httpserver.Config
	Mode string
}
