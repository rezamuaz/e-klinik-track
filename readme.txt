migrate create -ext sql -dir migrations/sql -seq user
migrate -path migrations/sql -database "postgresql://destiny:qn8prVZ6Cr75@localhost:5435/e-klinik?sslmode=disable" -verbose up

ALTER TABLE your_table ALTER COLUMN id RESTART WITH 1;
SELECT setval(pg_get_serial_sequence('lesson', 'lesson_id'), COALESCE((SELECT MAX(lesson_id)+1 FROM lesson), 1), false);

ssh -L 5432:localhost:5432 root@109.123.238.119
ssh -L 6379:localhost:6379 root@109.123.238.119


docker run --rm -v $(pwd):/src -w /src sqlc/sqlc generate