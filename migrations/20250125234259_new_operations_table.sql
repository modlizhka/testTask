-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS operations (
  id SERIAL PRIMARY KEY,
  type VARCHAR(255) NOT NULL,
  recipient INT,
  sender INT,
  amount DECIMAL(50, 3),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS operations;
-- +goose StatementEnd


