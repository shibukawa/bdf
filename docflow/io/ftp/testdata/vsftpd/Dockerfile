FROM alpine:3.21
RUN apk add --no-cache vsftpd openssl && \
    adduser -D -h /home/tester tester && echo "tester:s3cret" | chpasswd && \
    mkdir -p /home/tester/in /home/tester/done && chown -R tester:tester /home/tester && \
    openssl req -x509 -nodes -newkey rsa:2048 -days 30 -subj "/CN=127.0.0.1" -keyout /etc/vsftpd/key.pem -out /etc/vsftpd/cert.pem 2>/dev/null && \
    cat /etc/vsftpd/cert.pem /etc/vsftpd/key.pem > /etc/vsftpd/vsftpd.pem && \
    openssl x509 -in /etc/vsftpd/cert.pem -noout -fingerprint -sha256 | sed 's/.*=//' > /etc/vsftpd/fingerprint && \
    mkdir -p /var/run/vsftpd /var/log/vsftpd
RUN printf '%s\n' \
  'listen=YES' 'listen_ipv6=NO' 'background=NO' \
  'anonymous_enable=NO' 'local_enable=YES' 'write_enable=YES' 'local_umask=022' \
  'chroot_local_user=YES' 'allow_writeable_chroot=YES' \
  'pasv_enable=YES' 'pasv_min_port=21100' 'pasv_max_port=21110' 'pasv_address=127.0.0.1' 'pasv_addr_resolve=NO' \
  'ssl_enable=YES' 'rsa_cert_file=/etc/vsftpd/vsftpd.pem' 'rsa_private_key_file=/etc/vsftpd/vsftpd.pem' \
  'force_local_logins_ssl=YES' 'force_local_data_ssl=YES' 'require_ssl_reuse=YES' \
  'seccomp_sandbox=NO' 'xferlog_enable=YES' 'xferlog_std_format=NO' 'log_ftp_protocol=YES' 'vsftpd_log_file=/var/log/vsftpd/vsftpd.log' 'ftpd_banner=vsftpd for docflow' \
  > /etc/vsftpd/vsftpd.conf
EXPOSE 21 21100-21110
CMD ["sh", "-c", "cat /etc/vsftpd/fingerprint; exec vsftpd /etc/vsftpd/vsftpd.conf"]
