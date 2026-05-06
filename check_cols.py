#!/usr/bin/env python3
import sqlite3
conn = sqlite3.connect('/www/orange/data/data.db')
cursor = conn.cursor()
cursor.execute("PRAGMA table_info(team_requests)")
for row in cursor.fetchall():
    print(f"  {row[1]}")
conn.close()
