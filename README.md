# First-Time Setup/New Migrations

1. `cd docker`
2. `docker-compose --file docker-compose.yml build`
3. `docker-compose --file docker-compose.db.yml build`
4. `docker-compose --file docker-compose.yml up`
5. `docker-compose --file docker-compose.db.yml up`
6. `cd ../SeedDB`
7. `python3 migrate_to_postgres.py` or `python` if 3 is default for you

# Running the Application

1. `cd docker`
2. `docker-compose build`
3. `docker-compose up`

# Development

1. `cd docker`
2. `docker-compose up postgres`
3. `cd ..`
4. `buffalo dev`

The application will be available at `http://localhost:3000`
