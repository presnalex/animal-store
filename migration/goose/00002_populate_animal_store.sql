-- +goose Up
-- +goose StatementBegin
insert into animal_store values (1, 'Zebra', 5000);
insert into animal_store values (2, 'Elephant', 70000);
insert into animal_store values (3, 'Crocodile', 3000);
insert into animal_store values (4, 'Lion', 8000);
insert into animal_store values (5, 'Rhino', 45000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from animal_store;
-- +goose StatementEnd
