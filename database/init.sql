CREATE DATABASE IF NOT EXISTS fleetify_db;
USE fleetify_db;

-- 1. Tabel users
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    role ENUM('SA', 'APPROVAL') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 2. Tabel vehicles
CREATE TABLE vehicles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    license_plate VARCHAR(20) NOT NULL UNIQUE,
    model VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 3. Tabel master_items
CREATE TABLE master_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    item_name VARCHAR(100) NOT NULL,
    type ENUM('PART', 'SERVICE') NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- 4. Tabel maintenance_reports
CREATE TABLE maintenance_reports (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vehicle_id INT NOT NULL,
    created_by INT NOT NULL,
    odometer INT NOT NULL,
    complaint TEXT NOT NULL,
    status ENUM('PENDING_APPROVAL', 'APPROVED', 'COMPLETED') DEFAULT 'PENDING_APPROVAL',
    initial_photo VARCHAR(255),
    proof_photo VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
    FOREIGN KEY (created_by) REFERENCES users(id)
) ENGINE=InnoDB;

-- 5. Tabel report_items
CREATE TABLE report_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    report_id INT NOT NULL,
    item_id INT NOT NULL,
    quantity INT NOT NULL,
    price_snapshot DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (report_id) REFERENCES maintenance_reports(id) ON DELETE CASCADE,
    FOREIGN KEY (item_id) REFERENCES master_items(id)
) ENGINE=InnoDB;

-- ==========================================
-- DATA SEEDER
-- ==========================================

-- Minimal 2 user
INSERT INTO users (username, role) VALUES 
('sajon', 'SA'),
('bos_approval', 'APPROVAL');

-- Minimal 3 kendaraan
INSERT INTO vehicles (license_plate, model) VALUES 
('B 1234 ABC', 'Toyota Avanza'),
('B 5678 DEF', 'Honda Innova'),
('D 9101 GHI', 'Suzuki Ertiga');

-- Minimal 5 master items
INSERT INTO master_items (item_name, type, price) VALUES 
('Oli Mesin 10W-40', 'PART', 150000.00),
('Filter Oli', 'PART', 35000.00),
('Kampas Rem Depan', 'PART', 250000.00),
('Jasa Ganti Oli', 'SERVICE', 50000.00),
('Jasa Servis Rem', 'SERVICE', 100000.00);