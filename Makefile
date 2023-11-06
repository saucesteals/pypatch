DLL_OUT = "./out/patch.dll"
INJECTOR_OUT = "./out/injector.exe"

zigdll:
	CC="zig cc -target x86_64-windows" CXX="zig c++ -target x86_64-windows" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o $(DLL_OUT) -buildmode=c-shared .

dll:
	go build -o $(DLL_OUT) -buildmode=c-shared .

injector:
	go build -o $(INJECTOR_OUT) ./cmd/injector