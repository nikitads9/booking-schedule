-- name: CheckAvailibility :one
SELECT NOT EXISTS (
    SELECT 1 FROM bookings b
    WHERE ((b.suite_id = $2 AND (
        (b.start_date >= sqlc.arg(start_date)::text AND b.start_date <= sqlc.arg(end_date)::text) OR (b.end_date >= sqlc.arg(start_date)::text AND b.end_date <= sqlc.arg(end_date)::text)
        ))) ) as availible, 
        (SELECT EXISTS (
            SELECT 1 FROM bookings b
            WHERE (((b.suite_id = $2 AND b.user_id = $3) OR b.id = $1) 
            AND 
            ((b.start_date > sqlc.arg(start_date)::text AND b.start_date < sqlc.arg(end_date)::text) OR (b.end_date > sqlc.arg(start_date)::text AND b.end_date < sqlc.arg(end_date)::text)) ))) as occupied_by_client;

-- name: AddBooking :exec
INSERT INTO bookings 
(id, user_id, suite_id, start_date, end_date, created_at, notify_at)
VALUES ($1,$2,$3,$4,$5,$6, notify_at = sqlc.narg('notify_at'));

-- name: GetBookingListByDate :many
SELECT id, suite_id, start_date, end_date, notify_at, created_at, updated_at, user_id 
FROM bookings 
WHERE 
((start_date > $1 AND start_date <= sqlc.arg(EndDate)::text) 
OR 
(start_date-notify_at > $1 AND start_date-notify_at <= sqlc.arg(EndDate)::text));

-- name: DeleteBookingsBeforeDate :exec
DELETE 
FROM bookings 
WHERE end_date < $1;

-- name: DeleteBooking :exec
DELETE 
FROM bookings 
WHERE (id = $1 AND user_id = $2);

-- name: GetBooking :one
SELECT id, suite_id, start_date, end_date, notify_at, created_at, updated_at, user_id 
FROM bookings 
WHERE (id = $1 AND user_id = $2);

-- name: GetBookings :many
SELECT id, suite_id, start_date, end_date, notify_at, created_at, updated_at, user_id 
FROM bookings 
WHERE (user_id = $1 AND ((start_date >= sqlc.arg(StartDate)::text AND start_date <= sqlc.arg(EndDate)::text) 
OR 
(end_date >= sqlc.arg(StartDate)::text AND end_date <= sqlc.arg(EndDate)::text)));


-- name: GetOccupiedIntervals :many
SELECT start_date as start, end_date as end 
FROM bookings 
WHERE (suite_id = $1 AND (end_date > sqlc.arg(Now)::text AND start_date < sqlc.arg(Month)::text)) 
ORDER BY start_date;


-- name: GetVacantRooms :many
SELECT DISTINCT rooms.id AS suite_id, name, capacity 
FROM rooms 
WHERE NOT EXISTS (
    SELECT 1 FROM bookings AS e 
    WHERE (e.suite_id=rooms.id AND ((e.start_date < sqlc.arg(start_date)::text AND e.end_date > sqlc.arg(end_date)::text) 
    OR 
    (e.start_date < sqlc.arg(end_date)::text AND e.end_date > sqlc.arg(start_date)::text)))) 
    OR 
    NOT EXISTS (SELECT DISTINCT suite_id FROM bookings);

-- name: UpdateBooking :exec
UPDATE bookings 
SET suite_id = $3, start_date = $4, end_date = $5, notify_at = sqlc.narg('notify_at'), updated_at = $6
WHERE (id = $1 AND user_id = $2);
