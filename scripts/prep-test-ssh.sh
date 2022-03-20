# Creates a docker container with ssh and a config file for connecting to it
SSH_KEY_NAME=${SSH_KEY_NAME:-ci-test-key}
SSH_KEY_PATH=ssh-key/ci
SSH_USER=ci-dev
rm -rf $SSH_KEY_PATH
mkdir -p $SSH_KEY_PATH
ssh-keygen -f "$SSH_KEY_PATH/${SSH_KEY_NAME}" -N ""
docker kill saido-linux-sshserver | true
docker run -d -p 2222:2222 -e USER_NAME=$SSH_USER --name saido-linux-sshserver -v $(pwd)/ci.pub:/config/.ssh/authorized_keys linuxserver/openssh-server
cat <<EOF > config-ci.yaml
hosts:
	connection: 
		type: ssh
		username: ${SSH_USER}
    port: 2222
		private_key_path: $(pwd)/${SSH_KEY_NAME}.pub
	children:
		"$(docker inspect -f "{{ .NetworkSettings.IPAddress }}" saido-linux-sshserver)":
		"127.0.0.1":
      connection: 
        type: local

metrics:
	- memory
	- cpu
EOF
