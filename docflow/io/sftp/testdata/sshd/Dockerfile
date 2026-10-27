FROM alpine:3.21
RUN apk add --no-cache openssh-server && \
    adduser -D -h /home/tester -s /bin/sh tester && echo "tester:s3cret" | chpasswd && \
    mkdir -p /home/tester/in /home/tester/done && chown -R tester:tester /home/tester && \
    ssh-keygen -A && \
    sed -i 's/^#\?PasswordAuthentication .*/PasswordAuthentication yes/; s/^#\?KbdInteractiveAuthentication .*/KbdInteractiveAuthentication yes/' /etc/ssh/sshd_config && \
    echo "Subsystem sftp internal-sftp" >> /etc/ssh/sshd_config
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D", "-e"]
