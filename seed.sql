-- ============================================
-- Hotel Arca — Seed Data
-- Run against PostgreSQL database: db_hotel_arca
-- Each service uses its own schema, so set search_path per section.
-- ============================================

-- ==========================================
-- Schema: auth (Auth Service)
-- ==========================================
SET search_path TO auth;

-- Owner (email: owner@arca.com / password: password123)
INSERT INTO "owner" (name, email, password) VALUES
('Pak Budi', 'owner@arca.com', '$2a$10$JtHnkOT2r7i6sIFDt082jOuNlpvUZYqXv7epdw7RC2dgmU8DAHeTW');

-- Staff (email: staff@arca.com / password: password123)
INSERT INTO staff (name, email, password) VALUES
('Rina Staff', 'staff@arca.com', '$2a$10$Y559wh9uSWXuGxX7dggu7OLCQFLM7sNZ.7MEd6UZ.aD/nEKRZP/AO');

-- Customer (email: customer@arca.com / password: password123)
INSERT INTO customer (name, email, password, phone_number) VALUES
('Andi Tamu', 'customer@arca.com', '$2a$10$0Eib4uuc1FEzF5TTKQunvOzuRzFfSWe7Qa29K9Tpi69yEUAchyUcO', '081234567890');


-- ==========================================
-- Schema: room (Room Service)
-- ==========================================
SET search_path TO room;

INSERT INTO room_type (name, price, description) VALUES
('Standard', 350000, 'Kamar nyaman dengan single bed, AC, TV. Cocok untuk solo traveler.'),
('Deluxe', 750000, 'King bed, city view, bathtub, minibar. Cocok untuk pasangan.'),
('Suite', 1500000, 'Kamar mewah dengan living room, kolam renang pribadi, pelayanan 24 jam.');

INSERT INTO room (room_number, id_room_type, availability) VALUES
('101', 1, true),
('102', 1, true),
('201', 2, true),
('202', 2, true),
('301', 3, true);


-- ==========================================
-- Schema: booking (Booking Service)
-- ==========================================
SET search_path TO booking;


-- ==========================================
-- Schema: payment (Payment Service)
-- ==========================================
SET search_path TO payment;


-- ==========================================
-- Schema: chat (Chat Service)
-- ==========================================
SET search_path TO chat;


-- ==========================================
-- Schema: review (Review Service)
-- ==========================================
SET search_path TO review;


-- ==========================================
-- Schema: report (Report Service)
-- ==========================================
SET search_path TO report;
