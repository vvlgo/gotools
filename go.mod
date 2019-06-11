module github.com/vvlgo/gotools

go 1.12

require (
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/pkg/errors v0.8.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.1
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/tebeka/strftime v0.0.0-20140926081919-3f9c7761e312 // indirect
	github.com/weekface/mgorus v0.0.0-20181029072001-239539fe10e4
	golang.org/x/net v0.0.0-00010101000000-000000000000 // indirect
	google.golang.org/appengine v0.0.0-00010101000000-000000000000 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	golang.org/x/net => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190403152447-81d4e9dc473e
	google.golang.org/appengine => github.com/golang/appengine v1.1.0
)
