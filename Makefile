all:
	gopherjs build -o Client/client.js
	go build -o server