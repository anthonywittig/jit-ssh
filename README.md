# jit-ssh

## Use-case

You want to SSH into a machine that is unreachable. This program will SSH into a remote machine (middle man) and expose it's own SSH port, so that others can SSH into the middle man and then SSH into the original machine.

## Without this program

### Spin up an EC2 instance ("middle man")

* t3.nano has "up to 5Gb" network, which might be enough
* use Ubuntu
* enable SSH (default)
* pick the right key pair (create if you don't already have one)

### From the machine that you want to make accessible ("remote machine")

* copy the private key over (if not already present) and `chmod 400` it
* remote port forward our SSH port (22) to some other port (8765 in this example) with something like
  * `ssh -i PATH_TO_KEY -R 8765:localhost:22 ubuntu@ec2-...compute.amazonaws.com`
  * (you can get the remote address from the EC2 colsole by clicking on the "connect" option)

### From the machine that you want to connect from ("my machine")

* copy the private key over (if not already present) and `chmod 400` it
* local port forwarded the port on the remote machine (8765) to some other local port (8901 in this example)
  * `ssh -i PATH_TO_KEY -L 8901:localhost:8765 ubuntu@ec2...compute.amazonaws.com`
* then connect to the forwarded port
  * `ssh -p 8901 USER_NAME@127.0.0.1`
