# Dental Clinic Web Service in Go

This repository contains the source code for a web server built in Go, specifically tailored for a Dental Clinic application. 
The server uses the Gin framework and interacts with a MySQL database.

## Installation and Usage

### Prerequisites
- Go version 1.21.0 or higher
- MySQL server
- Git

### Installation Steps

1. Clone the repository:

`git clone https://github.com/10Daniel10/web-server-go-ExamenFinal.git`

`cd web-server-go-ExamenFinal`

2. Install dependencies:

`go mod tidy`

3. Run the MySQL database using Docker Compose:

`docker-compose up`

4. Run the server:

`go run main.go`

## Project Structure

The project is organized as follows:

- **cmd/server**
  - **config**: Contains configurations for the server setup.
  - **external/database**: Contains code related to external database connections.
  - **handler**: Contains handlers for various API endpoints.
  - **middleware**: Contains middleware for authentication and other purposes.

- **docs**
  
    Contains Swagger API documentation.

- **internal**

  - **patient**: Contains models and services related to patients.
  - **dentist**: Contains models and services related to dentists.
  - **appointment**: Contains models and services related to appointments.

## Available Methods

### Model: Dentist

- Create: Creates a new dentist.
- Get All: Retrieves all dentists.
- Get by ID: Retrieves a dentist by ID.
- Get by License: Retrieves a dentist by License.
- Update:
  - Put: Updates an existing dentist using the PUT method.
  - Patch: Partially updates an existing dentist using the PATCH method.
- Delete: Deletes a dentist.

### Model: Patient

- Create: Creates a new patient.
- Get All: Retrieves all patients.
- Get by ID: Retrieves a patient by ID.
- Get by DNI: Retrieves a patient by DNI.
- Update:
  - Put: Updates an existing patient using the PUT method.
  - Patch: Partially updates an existing patient using the PATCH method.
- Delete: Deletes a patient.

### Model: Appointment

- Create: Creates a new appointment (by providing the patient's DNI and the dentist's license number).
- Get All: Retrieves all appointments.
- Get by ID: Retrieves an appointment by ID.
- Get by DNI: Retrieves a appointment by patient DNI.
- Update:
  - Put: Updates an appointment patient using the PUT method.
  - Patch: Partially updates an existing appointment using the PATCH method.
- Delete: Deletes an appointment.
