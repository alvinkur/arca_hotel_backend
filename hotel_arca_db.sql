BEGIN;

-- ============================================================
-- Hotel Arca — Database Structure (GORM-compatible)
-- 8 schemas, 11 tables, bigint PKs, no custom ID formats
-- ============================================================

-- --------------------------------------------------------
-- 1. SCHEMA: auth (Auth Service — port 8001)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS auth;
SET search_path TO auth;

CREATE TABLE customer (
  id_customer bigserial NOT NULL,
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  phone_number text,
  PRIMARY KEY (id_customer),
  CONSTRAINT idx_customer_email UNIQUE (email)
);

CREATE TABLE owner (
  id_owner bigserial NOT NULL,
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  PRIMARY KEY (id_owner),
  CONSTRAINT idx_owner_email UNIQUE (email)
);

CREATE TABLE staff (
  id_staff bigserial NOT NULL,
  name text NOT NULL,
  email text NOT NULL,
  password text NOT NULL,
  PRIMARY KEY (id_staff),
  CONSTRAINT idx_staff_email UNIQUE (email)
);

-- --------------------------------------------------------
-- 2. SCHEMA: room (Room Service — port 8002)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS room;
SET search_path TO room;

CREATE TABLE room_type (
  id_room_type bigserial NOT NULL,
  name text NOT NULL,
  price numeric NOT NULL,
  description text,
  PRIMARY KEY (id_room_type),
  CONSTRAINT idx_room_type_name UNIQUE (name)
);

CREATE TABLE room (
  id_room bigserial NOT NULL,
  room_number text NOT NULL,
  id_room_type bigint NOT NULL,
  availability boolean DEFAULT true,
  PRIMARY KEY (id_room),
  CONSTRAINT idx_room_room_number UNIQUE (room_number),
  CONSTRAINT fk_room_room_type FOREIGN KEY (id_room_type) REFERENCES room_type (id_room_type)
);

-- --------------------------------------------------------
-- 3. SCHEMA: booking (Booking Service — port 8003)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS booking;
SET search_path TO booking;

CREATE TABLE booking (
  id_booking bigserial NOT NULL,
  id_customer bigint NOT NULL,
  id_room bigint NOT NULL,
  date_in timestamptz NOT NULL,
  date_out timestamptz NOT NULL,
  total_payment numeric,
  status_payment text DEFAULT 'pending',
  PRIMARY KEY (id_booking)
);

-- --------------------------------------------------------
-- 4. SCHEMA: payment (Payment Service — port 8004)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS payment;
SET search_path TO payment;

CREATE TABLE payment (
  id_payment bigserial NOT NULL,
  id_booking bigint NOT NULL,
  total_payment numeric NOT NULL,
  method text NOT NULL,
  date timestamptz,
  status text DEFAULT 'pending',
  PRIMARY KEY (id_payment)
);

-- --------------------------------------------------------
-- 5. SCHEMA: chat (Chat Service — port 8005)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS chat;
SET search_path TO chat;

CREATE TABLE chat (
  id_chat bigserial NOT NULL,
  id_customer bigint NOT NULL,
  id_staff bigint NOT NULL,
  date timestamptz DEFAULT now(),
  PRIMARY KEY (id_chat)
);

CREATE TABLE chat_message (
  id_chat_message bigserial NOT NULL,
  id_chat bigint NOT NULL,
  sender_type text NOT NULL,
  message text NOT NULL,
  date timestamptz DEFAULT now(),
  PRIMARY KEY (id_chat_message)
);

-- --------------------------------------------------------
-- 6. SCHEMA: review (Review Service — port 8006)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS review;
SET search_path TO review;

CREATE TABLE review (
  id_review bigserial NOT NULL,
  id_customer bigint NOT NULL,
  id_room bigint NOT NULL,
  rating bigint NOT NULL CHECK (rating >= 1 AND rating <= 5),
  comment text,
  PRIMARY KEY (id_review)
);

-- --------------------------------------------------------
-- 7. SCHEMA: report (Report Service — port 8007)
-- --------------------------------------------------------
CREATE SCHEMA IF NOT EXISTS report;
SET search_path TO report;

CREATE TABLE revenue_report (
  id_revenue bigserial NOT NULL,
  period text NOT NULL,
  total_revenue numeric,
  total_booking bigint,
  total_review bigint,
  detail_income text,
  PRIMARY KEY (id_revenue)
);

-- --------------------------------------------------------
-- 8. INDEXES
-- --------------------------------------------------------
CREATE INDEX IF NOT EXISTS idx_booking_id_customer ON booking.booking (id_customer);
CREATE INDEX IF NOT EXISTS idx_booking_id_room ON booking.booking (id_room);
CREATE INDEX IF NOT EXISTS idx_chat_id_customer ON chat.chat (id_customer);
CREATE INDEX IF NOT EXISTS idx_chat_id_staff ON chat.chat (id_staff);
CREATE INDEX IF NOT EXISTS idx_chat_message_id_chat ON chat.chat_message (id_chat);
CREATE INDEX IF NOT EXISTS idx_payment_id_booking ON payment.payment (id_booking);
CREATE INDEX IF NOT EXISTS idx_review_id_customer ON review.review (id_customer);
CREATE INDEX IF NOT EXISTS idx_review_id_room ON review.review (id_room);
CREATE INDEX IF NOT EXISTS idx_room_id_room_type ON room.room (id_room_type);

-- --------------------------------------------------------
-- 9. VIEWS
-- --------------------------------------------------------

CREATE OR REPLACE VIEW booking.view_booking_total_payment AS
SELECT
    b.id_booking,
    b.id_customer,
    b.id_room,
    b.date_in,
    b.date_out,
    (b.date_out::date - b.date_in::date) AS total_night,
    rt.price,
    (b.date_out::date - b.date_in::date) * rt.price AS total_payment
FROM booking.booking b
JOIN room.room r ON b.id_room = r.id_room
JOIN room.room_type rt ON r.id_room_type = rt.id_room_type;

CREATE OR REPLACE VIEW booking.view_monthly_booking AS
SELECT
    to_char(booking.date_in, 'YYYY-MM') AS period,
    count(*) AS total_booking
FROM booking.booking
GROUP BY to_char(booking.date_in, 'YYYY-MM');

CREATE OR REPLACE VIEW booking.view_monthly_revenue AS
SELECT
    to_char(v.date_in, 'YYYY-MM') AS period,
    v.total_booking,
    v.total_revenue
FROM (
    SELECT date_in, count(*) AS total_booking, sum(total_payment) AS total_revenue
    FROM booking.view_booking_total_payment
    GROUP BY date_in
) v;

CREATE OR REPLACE VIEW review.view_total_review AS
SELECT count(*) AS total_review
FROM review.review;

COMMIT;
