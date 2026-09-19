package main

import (
	"os"
	"strings"
)

// swag v2 путает body в oneOf и не пишет блок examples (требование задания).
func main() {
	path := os.Args[1]
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	s := string(raw)
	s = strings.ReplaceAll(s, emptyExternalDocs, "")
	s = strings.ReplaceAll(s, oneOf("SendCommand", "Команда", "request"), ref("SendCommand"))
	s = strings.ReplaceAll(s, oneOf("DeviceStatusUpdate", "Новый статус on/off", "request"), ref("DeviceStatusUpdate"))

	for _, p := range patches {
		s = strings.ReplaceAll(s, p.old, p.new)
	}

	s = flipPatchDevice(s)

	if err = os.WriteFile(path, []byte(s), 0o644); err != nil {
		panic(err)
	}
}

func oneOf(schema, desc, summary string) string {
	return "schema:\n              oneOf:\n              - type: object\n              - $ref: '#/components/schemas/" +
		schema + "'\n                description: " + desc + "\n                summary: " + summary
}

func ref(schema string) string {
	return "schema:\n              $ref: '#/components/schemas/" + schema + "'"
}

func req(schema string) string {
	return "            schema:\n              $ref: '#/components/schemas/" + schema + "'"
}

func res(schema string) string {
	return "              schema:\n                $ref: '#/components/schemas/" + schema + "'"
}

const emptyExternalDocs = `externalDocs:
  description: ""
  url: ""
`

type patch struct{ old, new string }

var patches = []patch{
	{req("SendCommand"), req("SendCommand") + `
            examples:
              command:
                value:
                  device_id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f
                  action: "on"`},
	{req("DeviceStatusUpdate"), req("DeviceStatusUpdate") + `
            examples:
              status:
                value:
                  status: "off"`},
	{res("Command"), res("Command") + `
              examples:
                command:
                  value:
                    id: 01932c4e-8a1e-7333-8444-4d7e8f901234
                    device_id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f
                    action: "on"
                    source: user
                    result: ok
                    created_at: "2026-09-05T10:00:00Z"`},
	{res("Device"), res("Device") + deviceExample("on")},
	{res("NotFoundResponse"), res("NotFoundResponse") + `
              examples:
                not_found:
                  value:
                    errors:
                      - type: NotFoundError
                        message: Device not found
                        context:
                          id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f`},
	{res("ConflictResponse"), res("ConflictResponse") + `
              examples:
                inactive:
                  value:
                    errors:
                      - type: DeviceInactiveError
                        message: Device is inactive
                        context:
                          id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f`},
	{res("InternalErrorResponse"), res("InternalErrorResponse") + `
              examples:
                internal:
                  value:
                    errors:
                      - type: InternalServerError
                        message: Something went wrong`},
	{`              schema:
                items:
                  $ref: '#/components/schemas/TelemetryData'
                type: array`, `              schema:
                items:
                  $ref: '#/components/schemas/TelemetryData'
                type: array
              examples:
                readings:
                  value:
                    - id: 01932c4e-8a1f-7444-8555-5e8f90123456
                      device_id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f
                      value: 21.4
                      unit: °C
                      recorded_at: "2026-09-05T10:00:00Z"`},
	{`      - description: ID устройств через запятую
        example: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f,01932c4e-8a1c-7111-8222-2b5c6d7e8f90
        in: query
        name: device_ids`, `      - description: ID устройств через запятую
        examples:
          ids:
            value: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f,01932c4e-8a1c-7111-8222-2b5c6d7e8f90
        in: query
        name: device_ids`},
}

func deviceExample(status string) string {
	return `
              examples:
                device:
                  value:
                    id: 01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f
                    type_id: 01932c4e-8a1c-7111-8222-2b5c6d7e8f90
                    house_id: 01932c4e-8a1d-7222-8333-3c6d7e8f9012
                    type_code: heating
                    serial_number: RL-22
                    address: "192.168.10.4:47808"
                    status: "` + status + `"`
}

func flipPatchDevice(s string) string {
	i := strings.Index(s, "/api/v1/devices/{id}/status:")
	if i < 0 {
		return s
	}

	const on = `                    address: "192.168.10.4:47808"
                    status: "on"`
	const off = `                    address: "192.168.10.4:47808"
                    status: "off"`

	return s[:i] + strings.Replace(s[i:], on, off, 1)
}
