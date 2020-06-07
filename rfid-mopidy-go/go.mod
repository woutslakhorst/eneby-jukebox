module github.com/woutslakhorst/eneby-jukebox/rfid-mopidy-go

go 1.13

require (
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/kr/pretty v0.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.6.1
	github.com/woutslakhorst/go-rpio v4.2.0+incompatible
	golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.8
)

// see https://github.com/stianeikeland/go-rpio/pull/50
// replace github.com/stianeikeland/go-rpio/v4 => github.com/wfd3/go-rpi/v4 v4.4.0-2019-11282039
