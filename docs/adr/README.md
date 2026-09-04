# Architecture Decision Records

Канон: Michael Nygard. Шаблон: Context / Decision / Alternatives / Consequences.

| ADR | Решение | Status |
| --- | --- | --- |
| [ADR-001](0001-as-is-monolith.md) | As-Is описываем как один монолит | Accepted |
| [ADR-002](0002-sync-server-to-sensor.md) | Сервер синхронно опрашивает датчик | Accepted |
| [ADR-003](0003-postgresql.md) | PostgreSQL — единственное хранилище As-Is | Accepted |
| [ADR-004](0004-db-per-service.md) | Своя PostgreSQL на каждый сервис (To-Be) | Accepted |
| [ADR-005](0005-device-type-not-service.md) | Новые приборы — DeviceType, не сервис | Accepted |
| [ADR-006](0006-sync-http-and-broker.md) | HTTP синхронно, брокер для телеметрии сценариев | Accepted |
| [ADR-007](0007-auth-at-gateway.md) | Личность в user-house, доступ проверяет Gateway | Accepted |
