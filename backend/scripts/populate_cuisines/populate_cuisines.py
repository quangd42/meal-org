import csv
import os
import uuid
from datetime import datetime

import psycopg2
from dotenv import load_dotenv
from psycopg2 import sql, connect


def get_parent_id(cursor, parent_name, table_name):
    query = sql.SQL("SELECT id FROM {table} WHERE name = %s").format(
        table=sql.Identifier(table_name)
    )
    cursor.execute(query, (parent_name,))
    result = cursor.fetchone()
    return result[0] if result else None


def import_csv_to_postgres(csv_file_path, db_config, table_name):
    # Connect to PostgreSQL database
    conn = connect(
        dbname=db_config["dbname"],
        user=db_config["user"],
        password=db_config["password"],
        host=db_config["host"],
        port=db_config["port"],
    )
    cursor = conn.cursor()

    # Open CSV file
    with open(csv_file_path, mode="r") as file:
        reader = csv.DictReader(file)

        # Iterate over CSV rows and insert into database
        for row in reader:
            row_id = str(uuid.uuid4())  # Generate a unique UUID
            name = row["name"]
            parent_name = row["parent"]
            created_at = datetime.now()
            updated_at = created_at

            parent_id = None
            if parent_name:
                parent_id = get_parent_id(cursor, parent_name, table_name)

            insert_query = sql.SQL(
                """
                INSERT INTO {table} (id, name, parent_id, created_at, updated_at)
                VALUES (%s, %s, %s, %s, %s)
            """
            ).format(table=sql.Identifier(table_name))
            cursor.execute(
                insert_query, (row_id, name, parent_id, created_at, updated_at)
            )

    # Commit changes and close the connection
    conn.commit()
    cursor.close()
    conn.close()


if __name__ == "__main__":
    load_dotenv()  # Load environment variables from .env file

    db_config = {
        "dbname": os.getenv("DB_NAME"),
        "user": os.getenv("DB_USER"),
        "password": os.getenv("DB_PASSWORD"),
        "host": os.getenv("DB_HOST"),
        "port": os.getenv("DB_PORT"),
    }

    csv_file_path = "scripts/data/cuisines.csv"
    table_name = "cuisines"

    import_csv_to_postgres(csv_file_path, db_config, table_name)
