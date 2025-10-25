# compile outside and run
  - `GOOS=android GOARCH=arm64 go build -o hello-termux main.go`
  - sftp to termux
  - run inside termux `./hello-termux`

# installation
- download from external go binary (arm64) and add in path.
```bash
curl -LO https://go.dev/dl/go1.25.3.linux-arm64.tar.gz

~/storage/external-1 $ sha256sum go1.25.3.linux-arm64.tar.gz | grep 1d42ebc84999b5e2069f5e31b67d6fc5d67308adad3e178d5a2ee2c9ff2001f5
1d42ebc84999b5e2069f5e31b67d6fc5d67308adad3e178d5a2ee2c9ff2001f5  go1.25.3.linux-arm64.tar.gz

rm -rf $PREFIX/usr/local/go 

tar -C $PREFIX/local -xzf go1.25.3.linux-arm64.tar.gz

~ $ vim ~/.bashrc 
# add new path

~ $ cat ~/.bashrc
export PATH="$PATH:$PREFIX/local/go/bin"

source ~/.bashrc
# or exit and re-ssh

go version

```

# compile and running
facing problem, simple hello world program does not run, seems like there is separate build for android applications.
```
~/code/hello $ go run main.go
SIGSYS: bad system call
PC=0x15b80 m=4 sigcode=1

goroutine 23 gp=0x400008b180 m=4 mp=0x4000088008 [syscall]:
syscall.Syscall6(0x1b7, 0xffffffffffffff9c, 0x400026e780, 0x1, 0x200, 0x0, 0x0)
	syscall/syscall_linux.go:96 +0x2c fp=0x4000182a40 sp=0x40001829e0 pc=0xaf4dc
syscall.faccessat2(0xffffffffffffff9c, {0x400026e5f0?, 0x400026e730?}, 0x1, 0x200)
	syscall/zsyscall_linux_arm64.go:33 +0x84 fp=0x4000182aa0 sp=0x4000182a40 pc=0xabfa4
syscall.Faccessat(0xffffffffffffff9c, {0x400026e5f0, 0x45}, 0x1, 0x200)
	syscall/syscall_linux.go:167 +0x3c fp=0x4000182b80 sp=0x4000182aa0 pc=0xa95ac
internal/syscall/unix.Eaccess(...)

```