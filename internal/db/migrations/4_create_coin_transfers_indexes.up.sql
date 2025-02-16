CREATE INDEX IF NOT EXISTS from_user_id_index ON coin_transfers(from_user_id);
CREATE INDEX IF NOT EXISTS to_user_id_index ON coin_transfers(to_user_id);