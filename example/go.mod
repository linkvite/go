module example

go 1.24.5

require (
	github.com/joho/godotenv v1.5.1
	// github.com/linkvite/go is in ../
	github.com/linkvite/go v0.0.0
)

replace github.com/linkvite/go => ../
