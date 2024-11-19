### Go Workspaces

[Go Workspaces](https://go.dev/doc/tutorial/workspaces)
Go to each folder and run

1. go mod init "github/handle/current_module_OR_Folder

In the main folder run

2. go work init \* or ./moduleName

### Hot Reloading - For Air Run Command

1. alias air='$(go env GOPATH)/bin/air'
2. air init
3. air --> creates a toml file, which contains App Configuarion

Sometimes there are issues reloading as the Port is Still active
Consider adding below in the Toml File, Under "post_cmd"

post_cmd = ["lsof -ti:8080 | xargs kill"]

### Environment Setup: Direnv and gotoenv

Set up Enviroments for App

1. https://github.com/joho/godotenv
2. https://direnv.net/

### General: Kill Process at a Port

1. kill $(lsof -t -i:8080)
   or more violently:
2. kill -9 $(lsof -t -i:8080)

### General: Useful: Search through code

1. fn + F12
2. Crtl + '-' // For Forward and backward move

## Start Up

1. Run all services and get them Up, using Air
2. Allow Webhooks Listening: Run command

stripe listen --forward-to localhost:8081/webhook
