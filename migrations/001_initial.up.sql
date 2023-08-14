-- +goose Up
CREATE TABLE atm_message (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  mti VARCHAR(4) NOT NULL,
  "transaction" VARCHAR(20) NOT NULL,
  switch VARCHAR(50) NOT NULL,
  primary_account_number VARCHAR(19),
  transaction_amount VARCHAR(12),
  acquiring_institution_code VARCHAR(11) NOT NULL,
  receiving_institution_code VARCHAR(11) NOT NULL,
  terminal_name_location VARCHAR(99) NOT NULL,
  currency_code VARCHAR(3) NOT NULL,
  terminal_id VARCHAR(8) NOT NULL,
  source_account VARCHAR(28),
  destination_account VARCHAR(28),
  channel VARCHAR(20) NOT NULL,
  device VARCHAR(20) NOT NULL,
  target_bank VARCHAR(20),
  rrn VARCHAR(12) NOT NULL,
  trace_number VARCHAR(6) NOT NULL,
  transmission_date_time VARCHAR(10) NOT NULL,
  local_transaction_date_time VARCHAR(12) NOT NULL,
  original_data_elements VARCHAR(35),
  process_code VARCHAR(6) NOT NULL
);

CREATE TABLE config (
  "key" VARCHAR(256) NOT NULL PRIMARY KEY,
  "value" TEXT
);

INSERT INTO config ("key", "value") VALUES 
('HOST', 'localhost'),
('PORT', '5013'),
('BASTION_HOST', NULL),
('BASTION_PORT', NULL),
('TARGET_HOST', NULL),
('SSH_USERNAME', NULL),
('SSH_LOCAL_PORT', NULL),
('SSH_REMOTE_PORT', NULL),
('SSH_KEY', NULL),
('SSH_PASSPHRASE', NULL);