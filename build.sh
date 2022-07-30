envs=( "android/386" "android/amd64" "android/arm" "android/arm64" "darwin/amd64" "darwin/arm64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "illumos/amd64" "ios/amd64" "ios/arm64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "netbsd/386" "netbsd/amd64" "netbsd/arm" "netbsd/arm64" "openbsd/386" "openbsd/amd64" "openbsd/arm" "openbsd/arm64" "openbsd/mips64" "plan9/386" "plan9/amd64" "plan9/arm" "solaris/amd64" "windows/386/.exe" "windows/amd64/.exe" "windows/arm/.exe" "windows/arm64/.exe" )
 
for env in "${envs[@]}"
do
    echo $env
    export GOOS=$(echo $env | cut -f1 -d/)
    export GOARCH=$(echo $env | cut -f2 -d/)
    export SUFFIX=$(echo $env | cut -f3 -d/)
    go build -ldflags="-w -s" -o wg-http-proxy$SUFFIX . && tar -czf $GOOS-$GOARCH-wg-http-proxy.tar.gz wg-http-proxy$SUFFIX README.md LICENSE.txt
done