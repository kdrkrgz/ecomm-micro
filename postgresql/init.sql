psql -U ecommdb_root -tc "SELECT 1 FROM pg_database WHERE datname = 'orders'" | grep -q 1 || psql -U ecommdb_root -c "CREATE DATABASE orders"
psql -U ecommdb_root -tc "SELECT 1 FROM pg_database WHERE datname = 'payments'" | grep -q 1 || psql -U ecommdb_root -c "CREATE DATABASE payments"
