-- Add device_mac and target_ssid columns to auth_logs table
ALTER TABLE auth_logs ADD COLUMN device_mac VARCHAR(255) DEFAULT '' AFTER user_agent;
ALTER TABLE auth_logs ADD COLUMN target_ssid VARCHAR(255) DEFAULT '' AFTER device_mac;
