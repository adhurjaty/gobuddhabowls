# First-Time Setup/New Migrations

1. `docker-compose --file docker-compose.yml build`
2. `docker-compose --file docker-compose.db.yml build`
3. `docker-compose --file docker-compose.yml up`
4. `docker-compose --file docker-compose.db.yml up`
5. `cd app/SeedDB`
6. `python3 migrate_to_postgres.py` or `python` if 3 is default for you

# Running the Application

1. `docker-compose build`
2. `docker-compose up`

The application will be available at `http://localhost:3000`
