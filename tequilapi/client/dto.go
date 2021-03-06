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

package client

import "math/big"

// Fees represents the transactor fee
type Fees struct {
	Registration  *big.Int `json:"registration"`
	Settlement    *big.Int `json:"settlement"`
	Hermes        uint16   `json:"hermes"`
	DecreaseStake *big.Int `json:"decreaseStake"`
}

// RegistrationDataDTO holds input data required to register new myst identity on blockchain smart contract
type RegistrationDataDTO struct {
	Status     string `json:"status"`
	Registered bool   `json:"registered"`
}

// SettleRequest represents the request to settle accountant promises
type SettleRequest struct {
	HermesID   string `json:"hermes_id"`
	ProviderID string `json:"provider_id"`
}

// SettleWithBeneficiaryRequest represent the request to settle with new beneficiary address.
type SettleWithBeneficiaryRequest struct {
	SettleRequest
	Beneficiary string `json:"beneficiary"`
}

// DecreaseStakeRequest represents the decrease stake request.
type DecreaseStakeRequest struct {
	ID            string   `json:"id,omitempty"`
	Amount        *big.Int `json:"amount,omitempty"`
	TransactorFee *big.Int `json:"transactor_fee,omitempty"`
}
