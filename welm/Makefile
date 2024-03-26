build:
	elm make src/Main.elm --output=main.js
	tailwindcss -i tailwind.css -o style.css --minify
	go build -o ./bin/web .

css:
	tailwindcss -i tailwind.css -o style.css --watch

css-minify:
	tailwindcss -i tailwind.css -o style.css --minify --watch

dev:
	elm reactor

