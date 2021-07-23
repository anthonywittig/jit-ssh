buildrpi:
	go test ./...
	env GOOS=linux GOARCH=arm GOARM=5 go build

# Example commands that might be helpful:
#buildrpi2: buildrpi
#	scp jit-ssh .env.json jit-ssh.service pi@192.168.86.142:/usr/local/jit-ssh/2021_06_22/
#copyaws:
#	scp some.creds.aws.credentials pi@192.168.86.123:/home/pi/.aws/credentials
