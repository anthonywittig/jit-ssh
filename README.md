# jit-ssh

## Use-case

You want to SSH into a machine that is unreachable. This program will SSH into a remote machine (middle man) and expose it's own SSH port, so that others can SSH into the middle man and then SSH into the original machine.

## Without this program

1. spin up an EC2 instance ("middle man")
  * t3.nano has "up to 5Gb" network, which might be enough
  * use Ubuntu
  * enable SSH (default)
  * pick the right key pair (create if you don't already have one)

2. from the machine that you want to make accessible ("remote machine")
  * copy the private key over (if not already present) and `chmod 400` it
  * connect with something like `ssh -i PATH_TO_KEY ubuntu@ec2-....compute.amazonaws.com` (you can get the address from the EC2 console by clicking on the "connect" option)
