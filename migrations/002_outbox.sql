CREATE TABLE IF NOT EXISTS message_outbox (
 id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
 message_id BIGINT UNSIGNED NOT NULL,
 published_at TIMESTAMP NULL,
 created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
 UNIQUE KEY uq_outbox_message (message_id),
 INDEX idx_outbox_pending (published_at,id),
 CONSTRAINT fk_outbox_message FOREIGN KEY(message_id) REFERENCES messages(id)
) ENGINE=InnoDB;
