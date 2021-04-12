# GrackleDB
Unoptimized database written for fun and to learn Go.

It is extremely unlikely to be a wise choice for a database.

#### Supports:
* Queries
  * select * from z
  * select x from z
  * select x, y from z
  * select x from z where x = 'y'
  * select x, y from z where a = 'b', b = 2
  * select a from b | select c from d
* Commands
  * insert into z(x,y) values(1,'two')

#### Lacks:
* Persistent storage
* Most things you'd expect from a database
