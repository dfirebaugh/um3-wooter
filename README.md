# um3-wooter
woots about new ultimaker 3 printjobs in a slack channel


## .env file
a .env file must be created in the root of this repo that stores the environmental variables.
> SLACK_HOOK=https://< url to slack hook >

> UM3_URI=https://< url to ultimaker api>/um3/print_job

## How to start
Build the docker container
> docker-compose build

Stand up the docker container
> docker-compose up