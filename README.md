# jit-ssh

## Use-case

You want to SSH into a machine that is unreachable. This program will SSH into a remote machine (middle man) and expose it's own SSH port, so that others can SSH into the middle man and then SSH into the original machine.

## Without this program

### Spin up an EC2 instance

t3.nano has "up to 5Gb" network, which might be enough.
