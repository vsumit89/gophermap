# Gopher Map - Basic Key-Value Store

Gopher Map is a simple key-value store project that allows you to store and retrieve string mappings of keys to values. I implemented as a learning exercise in persistence using file logging. With Gopher Map, you can easily store and retrieve data with minimal configuration and setup.

## Features

- Store key-value pairs in a persistent manner using file logging.
- Retrieve values associated keys in O(1) as it's store in hashmap.
- Learning project for understanding the fundamentals of persistence and data storage.

## Getting Started

To get started with Gopher Map, follow these simple steps:

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/vsumit89/gopher-map.git
   ```

2. Navigate to the project directory:

   ```bash
   cd gopher-map/cmd
   ```

3. Build and compile the project:

   ```bash
    go run main.go
   ```

4. Run the project:

   ```bash
   ./main
   ```

## Usage

Gopher Map provides http api to store and retrieve key-value pairs. The following commands are supported:

- `GET /map/{key}`: Retrieves the value associated with the given key.
- `POST /map/{key}/{value}`: Stores the given key-value pair.
- `DELETE /map/{key}`: Deletes the given key-value pair.

## Data Persistence

Gopher Map achieves data persistence by writing key-value pairs to a log file. This log file stores the key-value pairs in a structured format, allowing data to be restored upon program restart. It's written to the `data` directory in the project root.

## Future Improvements

- Add support for multiple data types and data structures.
- Add support for database persistence.
