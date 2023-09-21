# Project Name Omnibus Price Check

This repository contains a collection of tools and utilities designed to streamline and optimize the process of data management. Key features include data ingestion, reporting, a webhook interface, and a fake data generator for testing purposes.

## Features

### 1. Data Ingestion
- **Description**: This module reads and processes external data, transforming it into a format suitable for insertion into our database.
- **Usage**: [Provide command or code snippet to run the data ingestion]

### 2. Report
- **Description**: Generates insightful reports based on the ingested data. Allows for a better understanding of the processed information.
- **Usage**: 
```bash
go run r.go
```
### 3. Webhook
- **Description**: A webhook interface that allows external services to send data directly to our system. Especially useful for real-time data integration, can be linked to price update events. Easy to scale if needed.
- **Endpoint**: `/webhook`
- **Parameters**: `sku`, `list_price`, and `final_price`
- **Method**: `GET`
- **Usage**: 
```bash
curl "http://yourserver:8082/webhook?sku=SKU12345&list_price=100.50&final_price=90.50"
```

### 4. Fake Data Generator
- **Description**: Generates random but coherent test data, useful for simulating real-world scenarios in a controlled testing environment.
- **Usage**:
```bash
go run fake.go
```

## Getting Started

### Prerequisites

- Go v1.21.1 or later
- MySQL (or your choice of database)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/florinel-chis/omnibus-pricecheck
cd omnibus-pricecheck
```

2. Install the required Go packages:
```bash
go get -v ./...
```

3. Adjust configuration settings as necessary in the `config.yaml` file.

4. Run the application:
```bash
go run main.go
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Ensure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
