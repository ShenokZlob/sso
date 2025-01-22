/* В реальных проектах так не делаем. Это для примера! Это небезопасно и просто плохо*/
ALTER TABLE users
    ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;