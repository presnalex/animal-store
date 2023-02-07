package storage

const (
	getAnimal = ` 
	select s.animal_id, s.animal, s.price
       from animal_store s
       where s.animal_id = $1;`
)
