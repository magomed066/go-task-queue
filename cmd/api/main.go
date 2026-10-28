package main

import "os"

func main() {
	cfg := config{
		db: dbConfig{
			dns: "db-string",
		},
	}

	app := application{
		addr:   ":8080",
		config: cfg,
	}

	err := app.run(app.mount())
	if err != nil {
		os.Exit(1)
	}

}
