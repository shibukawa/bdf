# A vsftpd for TestFTPReal

An FTPS server (AUTH TLS, a self-signed certificate, `require_ssl_reuse`
as vsftpd has it by default, no MLSD) with the user `tester` (password
`s3cret`), chrooted to a home with `in/` and `done/`, for
`TestFTPReal` (real_test.go), which is skipped without the environment
below. vsftpd prints the certificate's fingerprint first.

```sh
docker build -t docflow-vsftpd docflow/io/ftp/testdata/vsftpd
docker run -d --name docflow-vsftpd -p 2121:21 -p 21100-21110:21100-21110 docflow-vsftpd \
  sh -c 'cat /etc/vsftpd/fingerprint; while true; do vsftpd /etc/vsftpd/vsftpd.conf; sleep 0.2; done'
FP=$(docker logs docflow-vsftpd 2>&1 | head -1)
BDF_FTP_REAL_URL=ftp://tester@127.0.0.1:2121/ BDF_FTP_REAL_PASSWORD=s3cret BDF_FTP_REAL_FINGERPRINT=$FP \
  go test -v -run TestFTPReal ./docflow/io/ftp/
docker rm -f docflow-vsftpd
```

The loop is for Alpine's vsftpd, whose listener was seen to die (139)
after a session on an arm64 Mac; `/var/log/vsftpd/vsftpd.log` in the
container has every command and reply.
