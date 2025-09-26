-- name: GetAvailableSlot :one
SELECT available_slots.available_quantity FROM available_slots WHERE id = @slot_id;

-- name: ReserveSlotForOrder :exec
INSERT INTO slot_reservations (order_id, slot_id) VALUES ($1, $2) ON CONFLICT (order_id) DO UPDATE SET slot_id = $2;

-- name: UnreserveSlotForOrder :exec
DELETE FROM slot_reservations WHERE order_id = @order_id;

-- name: DecreaseAvailableSlot :exec
UPDATE available_slots SET available_quantity = available_quantity - 1 WHERE id = @slot_id;

-- name: IncreaseAvailableSlot :exec
UPDATE available_slots SET available_quantity = available_quantity + 1 WHERE id = @slot_id;

-- name: CheckOrderDelivery :one
SELECT * FROM slot_reservations WHERE order_id = @order_id;
