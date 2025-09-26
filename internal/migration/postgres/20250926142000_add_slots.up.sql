CREATE TABLE delivery_time_slots(
    id SERIAL PRIMARY KEY,
    time_slot string NOT NULL
);

INSERT INTO delivery_time_slots (time_slot)
VALUES ("09-12"),
    ("12-15"),
    ("15-18"),
    ("18-21");
    
CREATE TABLE available_slots(
    id int PRIMARY KEY,
    available_quantity int NOT NULL
);

CREATE TABLE slot_reservations(
    id SERIAL PRIMARY KEY,
    order_id string NOT NULL,
    slot_id int NOT NULL
);

COMMENT ON COLUMN delivery_time_slots.id IS 'id строки';
COMMENT ON COLUMN delivery_time_slots.time_slot IS 'Диапазон времени';
COMMENT ON COLUMN available_slots.id IS 'id диапазона';
COMMENT ON COLUMN available_slots.available_quantity IS 'Количество свободных слотов';
COMMENT ON COLUMN slot_reservations.id IS 'id строки';
COMMENT ON COLUMN slot_reservations.order_id IS 'идентификатор заказа';
COMMENT ON COLUMN slot_reservations.slot_id IS 'идентификатор диапазона';

CREATE UNIQUE INDEX idx_order_id_slot_reservations ON slot_reservations (order_id);
