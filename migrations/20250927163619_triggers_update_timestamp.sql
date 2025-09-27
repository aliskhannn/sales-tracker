-- +goose Up
-- +goose StatementBegin
-- Trigger function to update updated_at timestamp automatically.

CREATE OR REPLACE FUNCTION trg_set_updated_at()
    RETURNS TRIGGER
    LANGUAGE plpgsql AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_items_updated_at
    BEFORE UPDATE
    ON items
    FOR EACH ROW
EXECUTE FUNCTION trg_set_updated_at();

CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE
    ON categories
    FOR EACH ROW
EXECUTE FUNCTION trg_set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_items_updated_at ON items;
DROP TRIGGER IF EXISTS trg_categories_updated_at ON categories;
DROP FUNCTION IF EXISTS trg_set_updated_at();
-- +goose StatementEnd
