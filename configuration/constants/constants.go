package constants

/*
	config files constants
*/

// use for get the HOST from the .env file
const HOST string = "HOST"

// use for get the USER from the .env file
const USER string = "USER"

// use for get the PASSWORD from the .env file
const PASSWORD string = "PASSWORD"

// use for get the DATABASE from the .env file
const DATABASE string = "DATABASE"

// use for get the PORT from the .env file
const PORT string = "PORT"

// use for get the SSLMODE from the .env file
const SSLMODE string = "SSLMODE"

// use for get the JWTKEY from the .env file
const JWTKEY string = "JWTKEY"

/*
	status constants
*/

// use for enable status in the .env file
const ENABLE string = "enable"

// use for disable status in the .env file
const DISABLE string = "disable"

/*
	error response constants
*/

// use for error loading response
const ERRLoading string = "Error loading .env file"

// use for port error
const ERRPort string = "failed to get the port"

// use for jwt error
const ERRJWT string = "JWT key is missing"
