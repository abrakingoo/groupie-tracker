# Groupie Trackers

Groupie Trackers is a web application that allows users to explore and visualize information about various artists and bands, including their concert dates, locations, and other relevant details. The application fetches data from a provided API, manipulates it, and presents it in an engaging and user-friendly manner.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)
- [Usage](#usage)
- [API Structure](#api-structure)
- [Event Handling](#event-handling)
- [Contributing](#contributing)
- [License](#license)

## Features

- Display detailed information about artists and bands, including:
  - Artist/band names
  - Images
  - Year of activity
  - Date of first album
  - Band members
- List upcoming and past concert locations and dates.
- Interactive data visualizations such as:
  - Cards displaying artist information
  - Tables for concert dates and locations
- Responsive design for mobile and desktop users.
- search bar functionality that filter searches based on this criteria
  - artist/band name
  - members
  - locations
  - first album date
  - creation date
- search bar has typing suggestions as you key in your query
## Technologies Used

- **Frontend:**
  - HTML, CSS
  - javascript

- **Backend:**
  - Go (Golang)
  - RESTful API (provided)

- **Other:**
  - Git for version control

## Installation

To get started with Groupie Trackers, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://learn.zone01kisumu.ke/git/abrakingoo/groupie-tracker

2. Navigate to the project directory:
      ```bash
   cd groupie-tracker

3. Start the Go backend server:
    ```bash
    go run cmd/main.go

4. Open your browser and navigate to:
    ```bash
     http://localhost:8080


## Usage
Once the application is running, you can explore the following features:

    - Browse Artists: View all artists and their information.
    - Concert Locations: Check the latest and upcoming concert locations.
    - Concert Dates: View the scheduled dates for concerts.
    - Search: search for any specific queries on your page or run general search which will return the artist card related artist     cards or card related to your search

## API Structure
The application interacts with a provided API, which consists of four main parts:

### Artists:

- **Information about bands/artists, including:**
    - Name
    - Image URL
    - Year of activity
    - Date of first album
    - Members

- **Locations:**

    - Concert locations associated with each artist.

- **Dates:**

    - Concert dates associated with each artist.
- **Relation:**

    -Links between artists, concert dates, and locations.

## Contributing
Contributions are welcome! If you'd like to contribute to the project, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or fix:
    ```bash
    git checkout -b feature/YourFeatureName
3. Commit your changes:
    ```bash
    git commit -m "Add your message here"
4. Push to the branch:
    ```bash
    git push origin feature/YourFeatureName
5. Create a pull request.

## License
This project is licensed under the MIT License.