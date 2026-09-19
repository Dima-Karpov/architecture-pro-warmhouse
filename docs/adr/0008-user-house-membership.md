# ADR-008: User — House многие-ко-многим через HouseMember

- Status: Accepted
- Date: 2026-09-07
- Deciders: команда разработки «Тёплый дом»

## Context

Пример курса для задания 3: User — House **1:N**, у дома один `owner_id`. Домом владеют несколько человек или пользуется семья. Токен с одним `house_id` после этого не работает: у жильца много домов.

## Decision

User — House **N:N** через `HouseMember` в `user-house`. Хозяин — `role = owner`, не колонка на `House`. Несколько `owner` можно. Дом без `owner` нельзя: создание = `House` + `HouseMember(owner)` атомарно.

Доступ — членство, не «хозяин». Токен несёт только `user_id`. Клиент шлёт `house_id` (header/body); Gateway проверяет членство; `command` сверяет `device.house_id` ([ADR-007](0007-auth-at-gateway.md)).

Это сознательное отклонение от примера курса.

## Alternatives considered

- Как в задании: `House.owner_id`, 1:N — проще, но семья и второй хозяин не живут в модели.
- `owner_id` плюс отдельная таблица жильцов — два источника правды, кто хозяин.

## Consequences

- Positive: семья и несколько хозяев в одной связи; Gateway отличает resident от чужого жильца.
- Negative: инвайт, смена `role`, выход последнего `owner` — отдельные операции `user-house`.
- Follow-up: `command` проверяет «устройство этого дома»; членство уже проверил Gateway.
