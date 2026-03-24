sync-env:
	grep -o '^[^=]*' .env > .env.sample