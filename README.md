# combuuid

![status tests](https://github.com/pckilgore/combuuid/actions/workflows/pr.yml/badge.svg)

Generates v4 `github.com/google/uuid`s optimized for use as primary keys in
relational databases like PostgreSQL, MySQL, MariaDB, and SQLite.

These combuuids sacrifice two bytes of collision resistance to improve insertion
speed and avoid WAL-amplification by adding a 2 bytes of sequential data as the 
most significant bits of the uuid. These bits increase monotonically every 60
seconds (default) until wrapping over after 65535 periods have elapsed, which
greatly improves insertion speed into an index b-tree over random ids.

See [Sequential UUID Generators](https://www.2ndquadrant.com/en/blog/sequential-uuid-generators/)
by Tomas Vondra for more details, explaination, and benchmarks of insertion
performance and WAL sizes. See also [The Cost of GUIDs as Primary Keys](https://www.informit.com/articles/article.aspx?p=25862)
by Jimmy Nilsson from which the name `COMB`s originates.

Small, optimized, fuzz tested and benchmarked for minimal overhead over `uuid`
(only a few ns overhead beyond retreiving system time):

![combuuid benchmarks vs. google/uuid](/bench.jpg)

## Stop. Do you need this?

This is an optimization that incurs a small time cost over generating regular v4
uuids. You should use regular uuids if:
 - You're using a read-heavy workload.
 - You're not using relational database.
 - You're not storing UUIDs as your primary key in that relational database.
 - You absolutely need to rely on collision resistance of *universally* unique
   ID generation, not just *application* unique ID generation.
    - For example, you're interacting with another UUID-generating system that
      you do not have control over, and your application cannot reasonably
      recover from errors due to UNIQUE contraint violations when coordinating
      between systems that expect to share UUIDs.
    - In general:
      - Monolith:   Probably fine.
      - Serverless: Avoid (and you're probably using a NoSQL database anyways).

# Contributing

Contributions welcome. So is just copying the whole damn thing into your
codebase and tweaking. I only would say that I appreciate a star or comment in
thanks.
