# Table: businesses

This table shows data for Businesses.

The primary key for this table is **_cq_id**.

## Columns

| Name          | Type          |
| ------------- | ------------- |
|_cq_id (PK)|`uuid`|
|_cq_parent_id|`uuid`|
|id|`utf8`|
|external_id|`utf8`|
|created_at|`timestamp[us, tz=UTC]`|
|status|`utf8`|
|formation_entity_type|`utf8`|
|registrations_entity_types|`list<item: utf8, nullable>`|