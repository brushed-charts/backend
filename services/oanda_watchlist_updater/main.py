import os
import time
import parser
from db import DatabaseUpdater

HOST = os.getenv("MONGODB_HOST")
PORT = os.getenv("MONGODB_PORT")
DATABASE = os.getenv("MONGODB_OANDA_DBNAME")
COLLECTION = os.getenv("MONGODB_WATCHLIST_COLLECTION")
WATCHLIST_PATH = os.getenv('OANDA_WATCHLIST_PATH')
REFRESH_RATE = 60*5  # In seconds

if __name__ == "__main__":
    database = DatabaseUpdater(HOST, PORT, DATABASE, COLLECTION)
    while True:
        instruments = parser.get_instruments_from_watchlist(WATCHLIST_PATH)
        database.update(instruments)
        time.sleep(REFRESH_RATE)
