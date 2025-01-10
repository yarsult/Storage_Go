# Storage_Go #

Storage_Go is a lightweight, file-based storage server built with Go. It offers RESTful APIs to manage scalar values, slices, and maps, with support for persistence and periodic cleanup.

### 🚀 Features ### 
	•	Scalar, map, and slice storage with key-value pairs
	•	Persistent storage using JSON
	•	Periodic cleanup of expired entries
	•	RESTful APIs for easy integration
	•	Configurable via environment variables
	•	Logging with zap

### 🛠️ Installation ###

## Prerequisites ##
	•	Go version 1.18 or higher
	•	PostgreSQL database (optional for advanced features)

## Steps ##
1. Clone the repository:

``` git clone https://github.com/yarsult/Storage_Go.git ```
``` cd Storage_Go ```

2. Run with Docker Compose

```docker-compose up --build```



### 🔧 Configuration ###

You can configure the application using the following environment variables:
	•	STORAGE_FILE_PATH: Path to the JSON file for storage (default: slice_storage.json).
	•	BASIC_SERVER_PORT: Port for the server to run (default: 8090).
	•	POSTGRES: PostgreSQL connection string (optional for database integration).


### 📚 API Endpoints ###

# Health Check #

Endpoint: /health
Method: GET
Description: Check if the server is running.

# Scalar Operations #
**Set Value:**
POST /scalar/set/:key/:value
Sets a scalar value.

**Get Value:**
GET /scalar/get/:key
Retrieves the scalar value for the given key.

# Slice Operations #
**Push Value**
POST /slice/lpush/:key
Pushes a value to a slice.

# Map Operations #
**Set Field:**
POST /map/hset/:key
Sets a field in a map.

**Get Field:**
GET /map/hget/:key/:field
Retrieves a field value from the map.

### 🛡️ Security ###
	•	Use HTTPS in production.
	•	Regularly clean expired data using the built-in periodic cleaner.

 ￼
