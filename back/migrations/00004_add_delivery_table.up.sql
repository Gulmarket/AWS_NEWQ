CREATE TABLE delivery (
      id serial primary key,
      farm_box varchar(100) not null,
      box_size varchar(10) not null,
      mixed bool,
      species varchar(100) not null,
      product varchar(100) not null,
      color varchar(50) not null,
      length varchar(100) not null,
      price float not null,
      boxes int not null
)
