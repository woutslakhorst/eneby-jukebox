module github.com/woutslakhorst/eneby-jukebox/rfid-mopidy-go

go 1.13

require (
	github.com/google/gousb v0.0.0-20190812193832-18f4c1d8a750
	github.com/jpoirier/gousb v0.0.0-20160821211425-38255c2cef15 // indirect
	github.com/karalabe/hid v1.0.0
	github.com/signal11/hidapi v0.0.0-20160920034012-a6a622ffb680 // indirect
	github.com/spf13/cobra v0.0.5 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.6.2 // indirect
	github.com/woutslakhorst/go-rpio v4.2.0+incompatible
	gopkg.in/yaml.v2 v2.2.8
)

// see https://github.com/stianeikeland/go-rpio/pull/50
// replace github.com/stianeikeland/go-rpio/v4 => github.com/wfd3/go-rpi/v4 v4.4.0-2019-11282039
