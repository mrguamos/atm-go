-- +goose Up
CREATE TABLE atm_message (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  mti VARCHAR(4) NOT NULL CHECK (mti <> ''),
  "transaction" VARCHAR(20) NOT NULL CHECK ("transaction" <> ''),
  switch VARCHAR(50) NOT NULL CHECK (switch <> ''),
  primary_account_number VARCHAR(19) NOT NULL,
  transaction_amount VARCHAR(12) NOT NULL CHECK (transaction_amount <> ''),
  acquiring_institution_code VARCHAR(11) NOT NULL CHECK (acquiring_institution_code <> ''),
  receiving_institution_code VARCHAR(11) NOT NULL,
  terminal_name_location VARCHAR(99) NOT NULL,
  currency_code VARCHAR(3) NOT NULL CHECK (currency_code <> ''),
  terminal_id VARCHAR(8) NOT NULL CHECK (terminal_id <> ''),
  source_account VARCHAR(28) NOT NULL,
  destination_account VARCHAR(28) NOT NULL,
  channel VARCHAR(20) NOT NULL CHECK (channel <> ''),
  device VARCHAR(20) NOT NULL CHECK (device <> ''),
  target_bank VARCHAR(20) NOT NULL,
  rrn VARCHAR(12) NOT NULL,
  trace_number VARCHAR(6) NOT NULL CHECK (trace_number <> ''),
  transmission_date_time VARCHAR(10) NOT NULL CHECK (transmission_date_time <> ''),
  local_transaction_date_time VARCHAR(12) NOT NULL CHECK (local_transaction_date_time <> ''),
  original_data_elements VARCHAR(35) NOT NULL,
  process_code VARCHAR(6) NOT NULL CHECK (process_code <> '')
);

CREATE TABLE config (
  "key" VARCHAR(256) NOT NULL PRIMARY KEY CHECK ("key" <> ''),
  "value" TEXT NOT NULL
);

INSERT INTO config ("key", "value") VALUES 
('HOST', ''),
('PORT', ''),
('BASTION_HOST', ''),
('BASTION_PORT', ''),
('TARGET_HOST', ''),
('SSH_USERNAME', ''),
('SSH_LOCAL_PORT', ''),
('SSH_REMOTE_PORT', ''),
('SSH_KEY', ''),
('SSH_PASSPHRASE', '');