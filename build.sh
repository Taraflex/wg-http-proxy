DATE_SUFFIX=$(date --utc +"%y%m%d")
envs=( "windows/amd64/.exe" "windows/arm/.exe" "windows/arm64/.exe" "android/arm64" "darwin/amd64" "darwin/arm64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "freebsd/arm64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "linux/mips" "linux/mips64" "linux/mips64le" "linux/mipsle" "linux/ppc64" "linux/ppc64le" "linux/riscv64" "linux/s390x" "openbsd/amd64" "openbsd/arm64" "openbsd/mips64" "windows/386/.exe" )
 
for env in "${envs[@]}"
do
    echo $env
    export GOOS=$(echo $env | cut -f1 -d/)
    export GOARCH=$(echo $env | cut -f2 -d/)
    export SUFFIX=$(echo $env | cut -f3 -d/)
    go build -ldflags="-w -s" -o wg-http-proxy$SUFFIX . && tar -czf wg-http-proxy.$GOOS-$GOARCH.$DATE_SUFFIX-$GITHUB_RUN_NUMBER.tar.gz wg-http-proxy$SUFFIX README.md LICENSE.txt
done