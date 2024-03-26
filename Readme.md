
# SSH2EC2 - internal SSHIM helper

## What is SSM
Session Manager (SSM) provides more security over an SSH connection.
With SSM, a port isn’t exposed for SSH traffic, and it avoids any risk with users sharing keys.
SSM occurs within the AWS console, and it is tied to only one IAM user.


## What is this doing
The SSH shim reads from the `.aws/credentials` file to get the details of all your AWS accounts.
Using AWS SDK, checks where ec2 instance is located, it's sending user default ssh public key and opens tunnel to it.


## Installation
1. Download the latest release from the releases page.
2. Place the downloaded file into your directory of choice.


## Configuration
1. Open your `~/.ssh/config` file for editing.
2. Add the following configuration, replacing `_PATH_TO_APP_` with the actual path to where you placed the downloaded file:

```shell
Host i-*
  ProxyCommand sh -c "_PATH_TO_APP_ %h"
  StrictHostKeyChecking no
  User root
```

## Usage
Once you have performed the installation and configuration steps, you can use the SSH2EC2 as follows:

- To SSH into a host, you can use:

```shell
ssh i-xxxxxxxx
```

- SSH2EC2 can also be used with rsync, for instance:

```shell
rsync -avz -P file i-xxxxx:/tmp/test
```

## References
Big thanks to `github.com/mmmorris1975/ssm-session-client` for ssm protocol implementation in golang. 