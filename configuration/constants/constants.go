package constants

/*
	config files constants
*/

// use for get the HOST from the .env file
const POSTGRES_HOST string = "POSTGRES_HOST"

// use for get the USER from the .env file
const POSTGRES_USER string = "POSTGRES_USER"

// use for get the PASSWORD from the .env file
const POSTGRES_PASSWORD string = "POSTGRES_PASSWORD"

// use for get the DATABASE from the .env file
const POSTGRES_DATABASE string = "POSTGRES_DATABASE"

// use for get the PORT from the .env file
const POSTGRES_PORT string = "POSTGRES_PORT"

// use for get the SSLMODE from the .env file
const POSTGRES_SSLMODE string = "POSTGRES_SSLMODE"

// use for get the REDISHOST from the .env file
const REDISHOST string = "REDIS_HOST"
const REDISPORT string = "REDIS_PORT"
const REDISDB string = "REDIS_DB"

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
