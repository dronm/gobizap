module gobizap

go 1.21.3

replace github.com/dronm/gobizapp/config => ../config
replace github.com/dronm/sqlmigr => /home/andrey/go/sqlmigr

require (
	github.com/dronm/gobizapp/config v0.0.0
	github.com/hoisie/mustache v0.0.0-20160804235033-6375acf62c69
)

require (
	github.com/dchest/jsmin v0.0.0-20220218165748-59f39799265f // indirect
	github.com/dronm/sqlmigr v0.0.0-20240802045545-3d4b879e264d // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
