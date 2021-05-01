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

## With this program

### AWS setup

* create a private S3 bucket and upload the remote config there (see the `example.env.remote.json`)
* create a new IAM (group and/or just) user that can access the remote config
  * copy the user credentials to the remote machine (in `~/.aws/credentails`)
* create an EC2 Key Pair
  * you'll need the .pem file on the local and remote machines since they'll both use it to connect to the middle man

### Remote machine setup

* configure jit-ssh
  * create an `.env.json` file (based off of the `.example.env.json)
* install jit-ssh as a service that will restart if it crashes
  * it's probably a good idea to namespace jit-ssh so that you can upgrade without killing the service, e.g. `/usr/local/jit-ssh/[DATE]/`
  * install and run as a service - see `jit-ssh.service` for more information

### When ready to connect

* spin up EC2 instance and update the remote config to have the connection string (e.g. `ubuntu@ec2-...compute.amazonaws.com`)
* on your local machine
  * `ssh -i PATH_TO_KEY -L 8901:localhost:8765 ubuntu@ec2...compute.amazonaws.com`
    * the first port number (8901) can be anything you want - it's the local port you'll use in a minute
    * the second port number (8765) must match the `env.remote.json`
    * this command can appear to work until you run the next, at which point it might say something like `channel 3: open failed: connect failed: Connection refused` - that probably means you're using the wrong port or the remote machine hasn't remote port forwarded
  * with the above command running, open a second termainal and run `ssh -p 8901 USER_NAME@127.0.0.1`

Debugging:
* on the middleman, install netstate `sudo apt install net-tools` and run `sudo netstat -lntu`. If the remote machine has connected you'll see:
  * `tcp        0      0 127.0.0.1:8765          0.0.0.0:*               LISTEN`
  * `tcp6       0      0 ::1:8765                :::*                    LISTEN`
