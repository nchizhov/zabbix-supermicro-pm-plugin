module zbplugin

go 1.26.7

replace ipmi => ./ipmi

replace support => ./support

require (
	golang.zabbix.com/sdk v1.2.2-0.20260828131423-e396d1d91436
	ipmi v0.0.0-00010101000000-000000000000
)

require (
	github.com/Microsoft/go-winio v0.6.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	support v0.0.0-00010101000000-000000000000 // indirect
)
