-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.animal_store
(
    animal_id integer,
    animal character varying(50) COLLATE pg_catalog."default",
    price integer,
	CONSTRAINT animal_store_pk PRIMARY KEY (animal_id)
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE IF EXISTS public.animal_store OWNER to postgres;
-- +goose StatementEnd
-- +goose StatementBegin
GRANT ALL ON TABLE public.animal_store TO public;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.animal_store;
-- +goose StatementEnd
