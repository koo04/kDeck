module kdeck-client

go 1.16

require (
	github.com/koo04/kdeck-server v0.0.0-20210317033807-62d5d5d944ab
	github.com/leaanthony/mewn v0.10.7
	github.com/micmonay/keybd_event v1.1.0 // indirect
	github.com/wailsapp/wails v1.11.0
	golang.org/x/net v0.0.0-20210315170653-34ac3e1c2000 // indirect
	golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005 // indirect
	google.golang.org/genproto v0.0.0-20210315173758-2651cd453018 // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
)

replace github.com/koo04/kdeck-server => ../kdeck-server
