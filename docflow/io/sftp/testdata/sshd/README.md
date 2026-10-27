# An OpenSSH sshd for TestSFTPReal

An SSH server with `internal-sftp`, the user `tester` (password
`s3cret`, by password and keyboard-interactive), whose home has `in/`
and `done/`, for `TestSFTPReal` (real_test.go), which is skipped without
the environment below. The test fetches the host key from the server
first, as the editor's **Ask the server** does.

```sh
docker build -t docflow-sshd docflow/io/sftp/testdata/sshd
docker run -d --name docflow-sshd -p 2222:22 docflow-sshd
BDF_SFTP_REAL_URL='sftp://tester@127.0.0.1:2222/~/' BDF_SFTP_REAL_PASSWORD=s3cret \
  go test -v -run TestSFTPReal ./docflow/io/sftp/
docker rm -f docflow-sshd
```

`TestClientOpenSSH` (client_test.go) needs no server: it speaks the
protocol to the machine's own `sftp-server` binary over its standard
streams, where there is one (macOS, Linux with openssh-server).
