HOST_PORT=8002
sudo firewall-cmd --zone=public --add-port=$HOST_PORT/tcp

podman run -d \
  --name health-check \
  -p $HOST_PORT:8080 \
  --security-opt seccomp=/usr/share/containers/seccomp.json \
  --security-opt no-new-privileges \
  --cap-drop all \
  --tmpfs /tmp:mode=1777 \
  --pids-limit 1024 \
  --health-cmd 'wget -q -O /dev/null http://127.0.0.1:8080/check || exit 1' \
  --health-interval 30s \
  --health-retries 3 \
  --health-start-period 10s \
  --health-timeout 5s \
  --health-on-failure restart \
  --restart on-failure:3 \
  ghcr.io/zhassulan1/health-check:latest