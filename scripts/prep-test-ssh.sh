set -e
# Creates a docker container with ssh and a config file for connecting to it
SSH_KEY_NAME=${SSH_KEY_NAME:-ci-test-key}
SSH_KEY_PATH=$(pwd)/ssh-key/ci
SSH_USER=ci-dev
if [ -f "${SSH_KEY_PATH}/${SSH_KEY_NAME}" ]; then
      echo "${SSH_KEY_PATH}/${SSH_KEY_NAME} exists."
else
  rm -rf $SSH_KEY_PATH
  mkdir -p $SSH_KEY_PATH
  ssh-keygen -f "$SSH_KEY_PATH/${SSH_KEY_NAME}" -N ""
fi
docker kill saido-linux-sshserver | true
docker rm saido-linux-sshserver || echo 'could not remove container'
docker run -d -p 2222:2222 -e USER_NAME=$SSH_USER --name saido-linux-sshserver -v ${SSH_KEY_PATH}/${SSH_KEY_NAME}.pub:/config/.ssh/authorized_keys linuxserver/openssh-server
SSH_HOST="$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" saido-linux-sshserver)"
cat <<EOF > config-test.yaml
hosts:
  children:
    "0.0.0.0":
      connection:
        type: ssh
        username: ${SSH_USER}
        port: 2222
        private_key_path: "$SSH_KEY_PATH/${SSH_KEY_NAME}"
metrics:
  memory:
  disk:
  tcp:
  docker:
poll-interval: 10
EOF
