agent:
	go run . agent

login:
	go run . login

logout:
	go run . logout

gogog:
	go run .

kill:
	lsof -ti :8080 | xargs kill
	
build:
	templ generate
	go build -o ./tmp/gogog .

watch:
	air

codegen:
	templ generate
