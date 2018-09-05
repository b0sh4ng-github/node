/*
 * Copyright (C) 2017 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package discovery

import (
	"encoding/json"

	dto_discovery "github.com/mysteriumnetwork/node/service_discovery/dto"
)

// Bootstrap loads NATS discovery package into the overall system
func Bootstrap() {
	dto_discovery.RegisterContactDefinitionUnserializer(
		TypeContactNATSV1,
		func(rawDefinition *json.RawMessage) (dto_discovery.ContactDefinition, error) {
			var contact ContactNATSV1
			err := json.Unmarshal(*rawDefinition, &contact)

			return contact, err
		},
	)
}
