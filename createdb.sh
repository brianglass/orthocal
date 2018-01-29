rm oca_calendar.db
sqlite3 oca_calendar.db < sql/days.sql
sqlite3 oca_calendar.db < sql/readings.sql
sqlite3 oca_calendar.db < sql/zachalos.sql
sqlite3 oca_calendar.db < sql/xceptions.sql
