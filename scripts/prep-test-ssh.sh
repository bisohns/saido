# Creates a docker container with ssh and a config file for connecting to it
SSH_KEY_NAME=${SSH_KEY_NAME:-ci-test-key}
SSH_KEY_PATH=$(pwd)/ssh-key/ci
SSH_USER=ci-dev
CONFIG_OUTPUT_PATH=${CONFIG_OUTPUT_PATH:config-test.yaml}
rm -rf $SSH_KEY_PATH
mkdir -p $SSH_KEY_PATH
ssh-keygen -f "$SSH_KEY_PATH/${SSH_KEY_NAME}" -N ""
docker kill saido-linux-sshserver | true
docker run -d -p 2222:2222 -e USER_NAME=$SSH_USER --name saido-linux-sshserver -v ${SSH_KEY_PATH}/${SSH_KEY_NAME}.pub:/config/.ssh/authorized_keys linuxserver/openssh-server
cat <<EOF > $CONFIG_OUTPUT_PATH
hosts:
  children:
    "$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" saido-linux-sshserver)":
      connection:
        type: ssh
        port: 2222
        username: ${SSH_USER}
        private_key_path: $(pwd)/${SSH_KEY_NAME}
    "127.0.0.1":
      connection:
        type: local

metrics:
- memory
- cpu
EOF
