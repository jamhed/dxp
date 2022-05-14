# app

## Local

Populate `public/config.js` with:

```json
{
	"domain": "leapx-accelerator.us.auth0.com",
	"client_id": "...",
	"audience": "vueapp"
}
```

And then:
```sh
npm run serve
```

## Docker

```sh
docker build . -t vueapp
docker run -p 8080:80 \
	-e AUTH0_DOMAIN=leapx-accelerator.us.auth0.com \
	-e AUTH0_CLIENT_ID=... \
	-it vueapp
```
